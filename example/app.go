package app

import (
	"context"
	"net/http"

	"4d63.com/google-cloud-appengine-datastore-replicator/example/tasks"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func init() {
	http.Handle("/", userHandler(homeHandler))
	http.Handle("/new-task", userHandler(newTaskHandler))
}

type appHandler func(context.Context, http.ResponseWriter, *http.Request) error

func (h appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	err := h(c, w, r)
	if err != nil {
		log.Errorf(c, "Error: %#v", err)
		http.Error(w, "There was an error. Please try again.", 500)
	}
}

type userHandler func(context.Context, *user.User, http.ResponseWriter, *http.Request) error

func (h userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		log.Warningf(c, "Nil user on a user handler")
		http.Error(w, "Access denied.", 403)
		return
	}
	c, err := appengine.Namespace(c, u.ID)
	if err != nil {
		log.Errorf(c, "Error: %#v", err)
		http.Error(w, "Access denied.", 403)
		return
	}
	err = h(c, u, w, r)
	if err != nil {
		log.Errorf(c, "Error: %#v", err)
		http.Error(w, "There was an error. Please try again.", 500)
	}
}

func homeHandler(c context.Context, u *user.User, w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	tasks, err := tasks.GetAll(c)
	if err != nil {
		return err
	}
	for _, t := range tasks {
		w.Write([]byte(`<form action="/set-task-state">`))
		w.Write([]byte(`<input id="checkBox" type="checkbox">`))
		w.Write([]byte(` ` + t.Title + ` `))
		w.Write([]byte(`<input type="submit" value="Update" />`))
		w.Write([]byte(`</form>`))
	}
	w.Write([]byte(`<form action="/new-task">`))
	w.Write([]byte(`<input type="text" name="title" value="" />`))
	w.Write([]byte(`<input type="submit" value="Add Task" />`))
	w.Write([]byte(`</form>`))
	return nil
}

func newTaskHandler(c context.Context, u *user.User, w http.ResponseWriter, r *http.Request) error {
	title := r.FormValue("title")

	task := tasks.Task{Title: title}

	err := tasks.Store(c, &task)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", 302)

	return nil
}
