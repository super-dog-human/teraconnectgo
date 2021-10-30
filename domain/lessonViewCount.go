package domain

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// LessonViewCountKeyは、Redis内で使用されるLesson参照回数保持用のキーのプリフィックスをstringで返します。
func LessonViewCountKeyPrefix(currentTime time.Time) string {
	return fmt.Sprintf("lessonViewCount_%s", currentTime.Format("20060102"))
}

// LessonViewCountKeyは、Redis内で使用されるLesson参照回数保持用のキーをstringで返します。
func lessonViewCountKey(lessonID int64, currentTime time.Time) string {
	return fmt.Sprintf("%s_%d", LessonViewCountKeyPrefix(currentTime), lessonID)
}

// IncrementLessonViewCountは、lessonIDのLessonのViewCountを1つ増分します。
// 増分は即座に行われず、Redisに格納されます。その後、定時バッチでLesson.ViewCountとUser.TotalLessonViewCountに反映されます。
func IncrementLessonViewCount(ctx context.Context, lessonID int64) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ENDPOINT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	currentJSTTime := time.Now().In(jst)
	key := lessonViewCountKey(lessonID, currentJSTTime)

	_, err = rdb.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}
