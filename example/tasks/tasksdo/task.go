package tasksdo

type Task struct {
	TaskID string `datastore:"-"`
	Title  string `datastore:"title"`
}
