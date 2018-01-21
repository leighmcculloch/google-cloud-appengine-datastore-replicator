package app

import (
	"context"
	"net/http"

	replicator "4d63.com/google-cloud-appengine-datastore-replicator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func init() {
	http.Handle("/_ah/push-handlers/datastore-puts", appHandler(datastorePutsHandler))
}

type appHandler func(context.Context, http.ResponseWriter, *http.Request) error

func (h appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	err := h(c, w, r)
	if err != nil {
		log.Criticalf(c, "Error: %#v", err)
		http.Error(w, "There was an error. Please try again.", 500)
	}
}

func datastorePutsHandler(c context.Context, w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	source, key, props, err := replicator.Unpack(c, r.Body)
	if err != nil {
		return err
	}

	if source.AppID == appengine.AppID(c) {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	log.Infof(c, "Storing %s %#v %#v", source, key, props)

	storedKey, err := datastore.Put(c, key, props)
	if err != nil {
		return err
	}

	log.Infof(c, "Stored %s %#v %#v", source, storedKey, props)

	w.WriteHeader(http.StatusCreated)
	return nil
}
