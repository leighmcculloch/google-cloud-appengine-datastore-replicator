package taskstates

type Value string

const (
	ValueOpen Value = "open"
	ValueDone Value = "done"
)

type TaskState struct {
	TaskID      int64 `datastore:"-" json:"task_id"`
	TaskStateID int64 `datastore:"-" json:"task_state_id"`
	Value       Value `datastore:"value" json:"value"`
}
