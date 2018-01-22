package taskstatesdo

import (
	"context"

	replicator "4d63.com/google-cloud-appengine-datastore-replicator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const replicatorTopic = "datastore-puts"
const taskKind = "task"
const taskStateKind = "taskstate"

func GetAll(c context.Context, taskID string) ([]*TaskState, error) {
	taskKey := datastore.NewKey(c, taskKind, taskID, 0, nil)
	var taskStates []*TaskState
	taskStateKeys, err := datastore.NewQuery(taskStateKind).Ancestor(taskKey).Order("created").GetAll(c, &taskStates)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(taskStates); i++ {
		taskStates[i].TaskID = taskID
		taskStates[i].TaskStateID = taskStateKeys[i].StringID()
	}
	return taskStates, nil
}

func Store(c context.Context, taskState *TaskState) error {
	taskKey := datastore.NewKey(c, taskKind, taskState.TaskID, 0, nil)
	taskStateKey := datastore.NewKey(c, taskStateKind, taskState.TaskStateID, 0, taskKey)
	log.Infof(c, "Persisting %v %#v", taskStateKey, taskState)
	taskStateKey, err := datastore.Put(c, taskStateKey, taskState)
	if err != nil {
		return err
	}
	taskState.TaskStateID = taskStateKey.StringID()
	log.Infof(c, "Persisted %v %#v", taskStateKey, taskState)

	if appengine.IsDevAppServer() {
		log.Warningf(c, "Skipping publishing because dev app server")
		return nil
	}

	err = replicator.Publish(c, replicatorTopic, taskStateKey, taskState)
	if err != nil {
		return err
	}
	log.Infof(c, "Published %v %#v", taskStateKey, taskState)

	return nil
}
