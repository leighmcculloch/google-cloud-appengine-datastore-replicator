package taskstatesdo

import "time"

type TaskState struct {
	TaskID      string    `datastore:"-"`
	TaskStateID string    `datastore:"-"`
	Created     time.Time `datastore:"created"`
	State       string    `datastore:"value"`
}
