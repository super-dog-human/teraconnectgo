package infrastructure

import (
	"context"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/golang/protobuf/ptypes"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

func CreateTask(ctx context.Context, queueID, name string, eta time.Time, message string) (*taskspb.Task, error) {
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloud tasks NewClient: %v", err)
	}
	defer client.Close()

	scheduleTime, err := ptypes.TimestampProto(eta)
	if err != nil {
		return nil, err
	}

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", ProjectID(), LocationID(), queueID)
	taskName := fmt.Sprintf("%s/tasks/%s", queuePath, name)
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
					HttpMethod:  taskspb.HttpMethod_POST,
					RelativeUri: "/zip_lesson",
				},
			},
			Name:         taskName,
			ScheduleTime: scheduleTime,
		},
	}

	req.Task.GetAppEngineHttpRequest().Body = []byte(message)

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloud tasks CreateTask: %v", err)
	}

	return createdTask, nil
}
