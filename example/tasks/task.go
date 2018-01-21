package tasks

type Task struct {
	TaskID int64  `datastore:"-" json:"task_id"`
	Title  string `datastore:"title" json:"title"`
}
