package tasksdo

import (
	"context"

	replicator "4d63.com/google-cloud-appengine-datastore-replicator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const replicatorTopic = "datastore-puts"
const taskKind = "task"

func Count(c context.Context) (int, error) {
	taskKeys, err := datastore.NewQuery(taskKind).KeysOnly().GetAll(c, nil)
	if err != nil {
		return 0, err
	}
	return len(taskKeys), nil
}

func GetAll(c context.Context) ([]*Task, error) {
	var tasks []*Task
	taskKeys, err := datastore.NewQuery(taskKind).GetAll(c, &tasks)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(tasks); i++ {
		tasks[i].TaskID = taskKeys[i].StringID()
	}
	return tasks, nil
}

func Get(c context.Context, taskID string) (*Task, error) {
	taskKey := datastore.NewKey(c, taskKind, taskID, 0, nil)
	var task Task
	err := datastore.Get(c, taskKey, &task)
	if err != nil {
		return nil, err
	}
	task.TaskID = taskKey.StringID()
	return &task, nil
}

func Store(c context.Context, task *Task) error {
	taskKey := datastore.NewKey(c, taskKind, task.TaskID, 0, nil)
	log.Infof(c, "Persisting %v %#v", taskKey, task)
	taskKey, err := datastore.Put(c, taskKey, task)
	if err != nil {
		return err
	}
	task.TaskID = taskKey.StringID()
	log.Infof(c, "Persisted %v %#v", taskKey, task)

	if appengine.IsDevAppServer() {
		log.Warningf(c, "Skipping publishing because dev app server")
		return nil
	}

	err = replicator.Publish(c, replicatorTopic, taskKey, task)
	if err != nil {
		return err
	}
	log.Infof(c, "Published %v %#v", taskKey, task)

	return nil
}
