package infrastructure

import (
	"context"
	"fmt"
	"strconv"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 現状ではLessonMaterialの圧縮にしか使用していないので固定値
const (
	queueID     string = "compressLesson"
	relativeUri string = "/lesson_compressing"
)

func LessonCompressingTaskName(lessonID int64, currentTime time.Time, requestID string) string {
	// シーケンシャルな値を避けるため、ランダム文字列であるリクエストIDを先頭に付与する
	return requestID + "-" + strconv.FormatInt(lessonID, 10) + "-" + strconv.FormatInt(currentTime.UnixNano(), 10)
}

func CreateTask(ctx context.Context, name string, eta time.Time, message string) (*taskspb.Task, error) {
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	scheduleTime := timestamppb.New(eta)
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", ProjectID(), LocationID(), queueID)
	taskName := fmt.Sprintf("%s/tasks/%s", queuePath, name)
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
					HttpMethod:  taskspb.HttpMethod_POST,
					RelativeUri: relativeUri,
				},
			},
			Name:         taskName,
			ScheduleTime: scheduleTime,
		},
	}

	req.Task.GetAppEngineHttpRequest().Body = []byte(message)

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, err
	}

	return createdTask, nil
}
