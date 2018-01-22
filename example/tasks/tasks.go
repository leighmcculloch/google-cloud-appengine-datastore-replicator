package tasks

import (
	"context"
	"time"

	"4d63.com/google-cloud-appengine-datastore-replicator/example/randid"
	"4d63.com/google-cloud-appengine-datastore-replicator/example/tasks/tasksdo"
	"4d63.com/google-cloud-appengine-datastore-replicator/example/tasks/taskstatesdo"
)

const (
	StateIncomplete = "incomplete"
	StateComplete   = "complete"
)

type Task struct {
	taskDO       *tasksdo.Task
	taskStateDOs []*taskstatesdo.TaskState
}

func Create(c context.Context, title string) (*Task, error) {
	taskDO := &tasksdo.Task{
		TaskID: randid.Generate(),
		Title:  title,
	}
	err := tasksdo.Store(c, taskDO)
	if err != nil {
		return nil, err
	}
	return &Task{taskDO: taskDO}, nil
}

func Get(c context.Context, taskID string) (*Task, error) {
	taskDO, err := tasksdo.Get(c, taskID)
	if err != nil {
		return nil, err
	}

	taskStateDOs, err := taskstatesdo.GetAll(c, taskID)
	if err != nil {
		return nil, err
	}

	task := &Task{
		taskDO:       taskDO,
		taskStateDOs: taskStateDOs,
	}

	return task, nil
}

func GetAll(c context.Context) ([]*Task, error) {
	taskDOs, err := tasksdo.GetAll(c)
	if err != nil {
		return nil, err
	}

	tasks := make([]*Task, len(taskDOs))

	for i := 0; i < len(taskDOs); i++ {
		taskDO := taskDOs[i]
		taskStateDOs, err := taskstatesdo.GetAll(c, taskDO.TaskID)
		if err != nil {
			return nil, err
		}
		tasks[i] = &Task{
			taskDO:       taskDO,
			taskStateDOs: taskStateDOs,
		}
	}

	return tasks, nil
}

func (t *Task) ID() string {
	return t.taskDO.TaskID
}

func (t *Task) Title() string {
	return t.taskDO.Title
}

func (t *Task) State() string {
	if len(t.taskStateDOs) == 0 {
		return StateIncomplete
	}
	return t.taskStateDOs[len(t.taskStateDOs)-1].State
}

func (t *Task) Complete(c context.Context) error {
	taskStateDO := &taskstatesdo.TaskState{
		TaskID:      t.taskDO.TaskID,
		TaskStateID: randid.Generate(),
		Created:     time.Now(),
		State:       StateComplete,
	}
	return taskstatesdo.Store(c, taskStateDO)
}

func (t *Task) Incomplete(c context.Context) error {
	taskStateDO := &taskstatesdo.TaskState{
		TaskID:      t.taskDO.TaskID,
		TaskStateID: randid.Generate(),
		Created:     time.Now(),
		State:       StateIncomplete,
	}
	return taskstatesdo.Store(c, taskStateDO)
}

func (t *Task) Delete(c context.Context) error {
	return nil
}

func (t *Task) Recover(c context.Context) error {
	return nil
}
