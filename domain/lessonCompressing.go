package domain

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type LessonCompressingErrorCode uint

const (
	AlreadyCompressed  LessonCompressingErrorCode = 1
	AnotherTaskWillRun LessonCompressingErrorCode = 2
)

func (e LessonCompressingErrorCode) Error() string {
	switch e {
	case AlreadyCompressed:
		return "lesson compressing is already executed."
	case AnotherTaskWillRun:
		return "created another task after this."
	default:
		return "unknown error"
	}
}

type LessonMaterialForCompressing struct {
	LessonMaterial
	IsCompressing bool
}

func GetLessonMaterialForCompressing(ctx context.Context, id string) (LessonMaterialForCompressing, error) {
	key := datastore.NameKey("LessonMaterialForCompressing", id, nil)
	lessonMaterial := new(LessonMaterialForCompressing)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *lessonMaterial, err
	}

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err := tx.Get(key, lessonMaterial); err != nil {
			return err
		}

		if lessonMaterial.IsCompressing {
			return AlreadyCompressed
		}

		lessonMaterial.IsCompressing = true
		if _, err := tx.Put(key, lessonMaterial); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return *lessonMaterial, err
	}

	return *lessonMaterial, nil
}

func UpdateLessonAfterCompressing(ctx context.Context, id int64, durationSec float32, updated time.Time) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		key := datastore.IDKey("Lesson", id, nil)
		lesson := new(Lesson)
		if err := tx.Get(key, lesson); err != nil {
			return err
		}

		if lesson.Updated.After(updated) {
			return AnotherTaskWillRun
		}

		lesson.DurationSec = durationSec
		lesson.Published = updated
		if _, err := tx.Put(key, lesson); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func DeleteLessonMaterialForCompress(ctx context.Context, id string) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	key := datastore.NameKey("LessonMaterialForCompressing", id, nil)
	if err := client.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}

func CompressLesson(ctx context.Context, lesson *Lesson, taskName string, lessonMaterial *LessonMaterialForCompressing) error {
	workingDir, err := ioutil.TempDir(os.TempDir(), taskName)
	if err != nil {
		return err
	}

	defer func() {
		_ = os.RemoveAll(workingDir) // 一時ファイルの削除に失敗しても実害はないのでエラーは関知しない
	}()

	if err := downloadBGMFiles(ctx, lesson.ID, workingDir, &lessonMaterial.Musics); err != nil {
		return err
	}

	if err := downloadOrCreateVoiceFiles(ctx, lesson.ID, workingDir, lessonMaterial); err != nil {
		return err
	}

	if err := mixAllSounds(ctx, lesson, workingDir, lessonMaterial.DurationSec, &lessonMaterial.Musics, &lessonMaterial.Speeches); err != nil {
		return err
	}

	if err := createCompressedMaterialToCloudStorage(ctx, lessonMaterial); err != nil {
		return err
	}

	return nil
}

func createLessonMaterialForCompressing(ctx context.Context, id string, lessonMaterial *LessonMaterial) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	key := datastore.NameKey("LessonMaterialForCompressing", id, nil)
	if _, err := client.Put(ctx, key, lessonMaterial); err != nil {
		return err
	}

	return nil
}

func createCompressingTask(ctx context.Context, taskName string, materialID int64, currentTime time.Time) error {
	taskEta := currentTime.Add(5 * time.Minute)
	// タスクに必要な情報はtaskNameで事足りるのでmessageは空文字でよい
	if _, err := infrastructure.CreateTask(ctx, taskName, taskEta, ""); err != nil {
		return err
	}

	return nil
}

func downloadBGMFiles(ctx context.Context, lessonID int64, workingDir string, musics *[]LessonMusic) error {
	for _, music := range *musics {
		originalFilePath := fmt.Sprintf("bgm/%d.mp3", music.BackgroundMusicID)

		var bgmFile []byte
		if _, err := os.Stat(originalFilePath); err != nil {
			bucketName := infrastructure.PublicBucketName()
			bgmFile, err = infrastructure.GetFileFromGCS(ctx, bucketName, originalFilePath)
			if err != nil {
				return err
			}
		} else {
			bgmFile, err = ioutil.ReadFile(originalFilePath) // BGMが既にDL済みならそれを使用
			if err != nil {
				return err
			}
		}

		tmpFilePath := fmt.Sprintf("%s/bgm_%d.mp3", workingDir, music.BackgroundMusicID)
		if err := ioutil.WriteFile(tmpFilePath, bgmFile, 0644); err != nil {
			return err
		}
	}
	return nil
}

func downloadOrCreateVoiceFiles(ctx context.Context, lessonID int64, workingDir string, lessonMaterial *LessonMaterialForCompressing) error {
	for _, speech := range lessonMaterial.Speeches {
		if !speech.IsSynthesis {
			continue
		}

		bucketName := infrastructure.MaterialBucketName()
		var voiceFile []byte
		var err error

		if speech.VoiceID == 0 && speech.Subtitle != "" {
			voiceParams := CreateSynthesisVoiceParam{LessonID: lessonID, Text: speech.Subtitle}
			if speech.SynthesisConfig.Name == "" {
				voiceParams.VoiceSynthesisConfig = lessonMaterial.VoiceSynthesisConfig
			} else {
				voiceParams.VoiceSynthesisConfig = speech.SynthesisConfig
			}
			// フロント側から声生成は済んでいるはずなので、ここで都度生成することはほぼない
			voiceFile, err = CreateSynthesizedVoice(ctx, &voiceParams, bucketName, "")
			if err != nil {
				return err
			}
		} else if speech.VoiceID != 0 {
			filePath := CloudStorageVoiceFilePath(lessonID, speech.VoiceID)
			voiceFile, err = infrastructure.GetFileFromGCS(ctx, bucketName, filePath)
			if err != nil {
				return err
			}
		}

		tmpFilePath := fmt.Sprintf("%s/%d.mp3", workingDir, speech.VoiceID)
		if err := ioutil.WriteFile(tmpFilePath, voiceFile, 0644); err != nil {
			return err
		}
	}

	return nil
}

func mixAllSounds(ctx context.Context, lesson *Lesson, workingDir string, fullDurationSec float32, musics *[]LessonMusic, speeches *[]LessonSpeech) error {
	commandOptions := []string{"-y"}

	var inputFiles []string
	var filters []string

	fullDurationStr := strconv.FormatFloat(float64(fullDurationSec), 'f', 3, 64)
	silenceFileName := fmt.Sprintf("%s/silence.mp3", workingDir)
	out, err := exec.Command("ffmpeg", "-f", "lavfi", "-i", "anullsrc=channel_layout=mono", "-t", fullDurationStr, "-ab", "32k", silenceFileName).CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", string(out))
		return err
	}
	inputFiles = append(inputFiles, "-i", silenceFileName)

	for i, music := range *musics {
		if music.Action == MusicActionStop {
			continue
		}

		var nextMusic LessonMusic
		var durationSec float32
		if i+1 < len(*musics) {
			nextMusic = (*musics)[i+1]
			durationSec = nextMusic.ElapsedTime - music.ElapsedTime
		} else {
			durationSec = fullDurationSec - music.ElapsedTime
		}

		originalBgmFileName := fmt.Sprintf("%s/bgm_%d.mp3", workingDir, music.BackgroundMusicID)
		currentBgmFileName := fmt.Sprintf("%s/fixed_bgm_%d.mp3", workingDir, i)
		inputFiles = append(inputFiles, "-i", currentBgmFileName)

		musicDuration := strconv.FormatFloat(float64(durationSec), 'f', 3, 64)
		out, err := exec.Command("ffmpeg", "-stream_loop", "-1", "-fflags", "+genpts", "-i", originalBgmFileName, "-t", musicDuration, "-c", "copy", currentBgmFileName).CombinedOutput()
		if err != nil {
			fmt.Printf("%s\n", string(out))
			return err
		}

		fileIndex := len(inputFiles)/2 - 1
		filters = append(filters, musicFilter(fileIndex, durationSec, &music, &nextMusic))
	}

	for _, speech := range *speeches {
		delayMilliSec := speech.ElapsedTime * 1000
		voiceFileName := fmt.Sprintf("%s/%d.mp3", workingDir, speech.VoiceID)
		inputFiles = append(inputFiles, "-i", voiceFileName)
		fileIndex := len(inputFiles)/2 - 1
		filters = append(filters, fmt.Sprintf("[%d:a]adelay=%f|%f[%d];", fileIndex, delayMilliSec, delayMilliSec, fileIndex))
	}

	inputsCount := len(inputFiles) / 2
	allInputs := func() string {
		var inputs string
		for i := 0; i < inputsCount; i++ {
			inputs += fmt.Sprintf("[%d]", i)
		}
		return inputs
	}()
	outputFilePath := fmt.Sprintf("%s/speech.mp3", workingDir)

	commandOptions = append(commandOptions, inputFiles...)
	commandOptions = append(commandOptions, "-filter_complex") // filter_complexの引数のエスケープはexec.Commandが自動でつけてくれる
	commandOptions = append(commandOptions, fmt.Sprintf("%s%samix=dropout_transition=600000:inputs=%d,volume=%d", strings.Join(filters, " "), allInputs, inputsCount, inputsCount))
	commandOptions = append(commandOptions, "-ar", "44100", "-ab", "128k", "-acodec", "libmp3lame", outputFilePath)

	out, err = exec.Command("ffmpeg", commandOptions...).CombinedOutput()
	if err != nil {
		fmt.Printf("ffmpeg %s\n", strings.Join(commandOptions, " "))
		fmt.Printf("%s\n", string(out))
		return err
	}

	var bucketName string
	if lesson.Status == LessonStatusPublic {
		bucketName = infrastructure.PublicBucketName()
	} else {
		bucketName = infrastructure.MaterialBucketName()
	}
	bucketFilePath := fmt.Sprintf("lesson/%d/speech.mp3", lesson.ID)
	outputFile, err := ioutil.ReadFile(outputFilePath)
	if err != nil {
		return err
	}

	if err := infrastructure.CreateFileToGCS(ctx, bucketName, bucketFilePath, "audio/mpeg", outputFile); err != nil {
		return err
	}

	return nil
}

func musicFilter(fileIndex int, durationSec float32, music *LessonMusic, nextMusic *LessonMusic) string {
	delayMilliSec := music.ElapsedTime * 1000

	var fade string
	if nextMusic != nil {
		if music.IsFading && nextMusic.IsFading && nextMusic.Action == MusicActionStop {
			if durationSec < 6.0 {
				fade = "afade=t=in:st=0:d=3,afade=t=out:st=3:d=3"
			} else {
				fade = fmt.Sprintf("afade=t=in:st=0:d=3,afade=t=out:st=%f:d=3", durationSec-3.0)
			}
		} else if music.IsFading {
			fade = "afade=t=in:st=0:d=3"
		} else if nextMusic.IsFading && nextMusic.Action == MusicActionStop {
			if durationSec < 3.0 {
				fade = "afade=t=out:st=0:d=3"
			} else {
				fade = fmt.Sprintf("afade=t=out:st=%f:d=3", durationSec-3.0)
			}

		}
	} else {
		if music.IsFading {
			fade = "afade=t=in:st=0:d=3"
		}
	}

	if fade == "" {
		return fmt.Sprintf("[%d:a]volume=%f,adelay=%f|%f[%d];", fileIndex, music.Volume, delayMilliSec, delayMilliSec, fileIndex)
	} else {
		return fmt.Sprintf("[%d:a]%s,volume=%f,adelay=%f|%f[%d];", fileIndex, fade, music.Volume, delayMilliSec, delayMilliSec, fileIndex)
	}
}

func createCompressedMaterialToCloudStorage(ctx context.Context, lessonMaterial *LessonMaterialForCompressing) error {
	// jsonをzstd圧縮する
	return nil
}
