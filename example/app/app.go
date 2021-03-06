package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"4d63.com/google-cloud-appengine-datastore-replicator/example/tasks"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

const maxTasks = 30

func init() {
	http.Handle("/", appHandler(homeHandler))
	http.Handle("/task/new", appHandler(newTaskHandler))
	http.Handle("/task/complete", appHandler(completeTaskHandler))
	http.Handle("/task/incomplete", appHandler(incompleteTaskHandler))
}

type appHandler func(context.Context, *user.User, http.ResponseWriter, *http.Request) error

func (h appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		log.Warningf(c, "Nil user on a user handler")
		http.Error(w, "Access denied.", 403)
		return
	}
	namespacedCtx, err := appengine.Namespace(c, fmt.Sprintf("%s-%s", u.ID, strings.Replace(u.Email, "@", "-", -1)))
	if err != nil {
		log.Errorf(c, "Error: %#v", err)
		http.Error(w, "Access denied.", 403)
		return
	}
	c = namespacedCtx
	err = h(c, u, w, r)
	if err != nil {
		log.Errorf(c, "Error: %#v", err)
		http.Error(w, "There was an error. Please try again.", 500)
	}
}

func homeHandler(c context.Context, u *user.User, w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	allTasks, err := tasks.GetAll(c)
	if err != nil {
		return err
	}
	w.Write([]byte(`<html>`))
	w.Write([]byte(`<head>`))
	w.Write([]byte(`<style>`))
	w.Write([]byte(`body { font-family: sans-serif; }`))
	w.Write([]byte(`</style>`))
	w.Write([]byte(`</head>`))
	w.Write([]byte(`<body>`))
	w.Write([]byte(`<center>`))
	w.Write([]byte(`<h2>Tasks</h3>`))
	w.Write([]byte(`<h3>` + u.Email + `</h3>`))
	for _, t := range allTasks {
		if t.State() == tasks.StateIncomplete {
			w.Write([]byte(`<form action="/task/complete">`))
			w.Write([]byte(`<input type="hidden" name="taskID" value="` + t.ID() + `">`))
			w.Write([]byte(`<input type="checkbox" onchange="submit();">`))
			w.Write([]byte(` ` + t.Title()))
			w.Write([]byte(`</form>`))
		} else {
			w.Write([]byte(`<form action="/task/incomplete">`))
			w.Write([]byte(`<input type="hidden" name="taskID" value="` + t.ID() + `">`))
			w.Write([]byte(`<input type="checkbox" onchange="submit();" checked>`))
			w.Write([]byte(` <strike>` + t.Title() + `</strike>`))
			w.Write([]byte(`</form>`))
		}
	}
	if len(allTasks) < maxTasks {
		w.Write([]byte(`<form action="/task/new">`))
		w.Write([]byte(`<input type="text" name="title" value="" maxlength="30" />`))
		w.Write([]byte(`<input type="submit" value="Add Task" />`))
		w.Write([]byte(`</form>`))
	}
	w.Write([]byte(`</center>`))
	w.Write([]byte(`</body>`))
	w.Write([]byte(`</html>`))
	return nil
}

func newTaskHandler(c context.Context, u *user.User, w http.ResponseWriter, r *http.Request) error {
	taskCount, err := tasks.Count(c)
	if err != nil {
		return err
	}
	if taskCount >= maxTasks {
		http.Error(w, "Task limit reached.", 403)
		return nil
	}

	title := r.FormValue("title")

	if len(title) > 30 {
		title = title[:30]
	}

	_, err = tasks.Create(c, title)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", 302)

	return nil
}

func completeTaskHandler(c context.Context, u *user.User, w http.ResponseWriter, r *http.Request) error {
	taskID := r.FormValue("taskID")
	log.Infof(c, "Request to complete task %s", taskID)

	task, err := tasks.Get(c, taskID)
	if err != nil {
		return err
	}
	if task == nil {
		http.NotFound(w, r)
		return nil
	}

	err = task.Complete(c)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", 302)

	return nil
}

func incompleteTaskHandler(c context.Context, u *user.User, w http.ResponseWriter, r *http.Request) error {
	taskID := r.FormValue("taskID")
	log.Infof(c, "Request to incomplete task %s", taskID)

	task, err := tasks.Get(c, taskID)
	if err != nil {
		return err
	}
	if task == nil {
		http.NotFound(w, r)
		return nil
	}

	err = task.Incomplete(c)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", 302)

	return nil
}
