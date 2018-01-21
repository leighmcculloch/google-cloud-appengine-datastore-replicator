PROJECT_IDS ?= crossregion-a crossregion-b

run:
	dev_appserver.py --port=8080 app.yaml

deploy:
	go install ./...
	$(foreach PID,$(PROJECT_IDS), \
		gcloud --project $(PID) app deploy -q && \
	) true
	$(MAKE) -C $(GOPATH)/src/4d63.com/google-cloud-appengine-datastore-replicator/service PROJECT_IDS="$(PROJECT_IDS)" deploy

create-pubsub:
	$(MAKE) -C $(GOPATH)/src/4d63.com/google-cloud-appengine-datastore-replicator/service PROJECT_IDS="$(PROJECT_IDS)" create-pubsub

delete-pubsub:
	$(MAKE) -C $(GOPATH)/src/4d63.com/google-cloud-appengine-datastore-replicator/service PROJECT_IDS="$(PROJECT_IDS)" delete-pubsub
