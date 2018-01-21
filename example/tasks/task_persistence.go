package tasks

import (
	"context"

	replicator "4d63.com/google-cloud-appengine-datastore-replicator"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const replicatorTopic = "datastore-puts"
const taskKind = "task"

func GetAll(c context.Context) ([]*Task, error) {
	var tasks []*Task
	_, err := datastore.NewQuery(taskKind).GetAll(c, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func Get(c context.Context, taskID int64) (*Task, error) {
	taskKey := datastore.NewKey(c, taskKind, "", taskID, nil)
	var task Task
	err := datastore.Get(c, taskKey, &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func Store(c context.Context, task *Task) error {
	taskKey := datastore.NewKey(c, taskKind, "", task.TaskID, nil)
	log.Infof(c, "Persisting %v %#v", taskKey, task)
	taskKey, err := datastore.Put(c, taskKey, task)
	if err != nil {
		return err
	}
	task.TaskID = taskKey.IntID()
	log.Infof(c, "Persisted %v %#v", taskKey, task)
	err = replicator.Publish(c, replicatorTopic, taskKey, task)
	if err != nil {
		return err
	}
	log.Infof(c, "Published %v %#v", taskKey, task)
	return nil
}
