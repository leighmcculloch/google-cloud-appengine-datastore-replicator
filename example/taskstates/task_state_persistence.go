package taskstates

import (
	"context"

	replicator "4d63.com/google-cloud-appengine-datastore-replicator"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const replicatorTopic = "datastore-puts"
const taskKind = "task"
const taskStateKind = "taskstate"

func GetAll(c context.Context, taskID int64) ([]*TaskState, error) {
	taskKey := datastore.NewKey(c, taskKind, "", taskID, nil)
	var taskStates []*TaskState
	_, err := datastore.NewQuery(taskStateKind).Ancestor(taskKey).GetAll(c, &taskStates)
	if err != nil {
		return nil, err
	}
	return taskStates, nil
}

func Store(c context.Context, taskState *TaskState) error {
	taskKey := datastore.NewKey(c, taskKind, "", taskState.TaskID, nil)
	taskStateKey := datastore.NewKey(c, taskStateKind, "", taskState.TaskStateID, taskKey)
	log.Infof(c, "Persisting %v %#v", taskStateKey, taskState)
	taskStateKey, err := datastore.Put(c, taskStateKey, taskState)
	if err != nil {
		return err
	}
	taskState.TaskStateID = taskStateKey.IntID()
	log.Infof(c, "Persisted %v %#v", taskStateKey, taskState)
	err = replicator.Publish(c, replicatorTopic, taskStateKey, taskState)
	if err != nil {
		return err
	}
	log.Infof(c, "Published %v %#v", taskStateKey, taskState)
	return nil
}
