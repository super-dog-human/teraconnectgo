package domain

type Position2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// VoiceSynthesisConfig is synthesis voice settings. used from SynthesisVoice params and LessonSpeech.
// https://github.com/googleapis/go-genproto/blob/master/googleapis/cloud/texttospeech/v1beta1/cloud_tts.pb.go#L663
type VoiceSynthesisConfig struct {
	LanguageCode string  `json:"languageCode"` // ja-JP/en-US
	Name         string  `json:"name"`         // ja-JP-Wavenet-A~D/en-US-Wavenet-A~J
	SpeakingRate float64 `json:"speakingRate"` // 0.25 ~ 4.0
	Pitch        float64 `json:"pitch"`        // -20.0 ~ 20.0
	VolumeGainDb float64 `json:"volumeGainDb"` // -5.0 ~ 10.0
}
