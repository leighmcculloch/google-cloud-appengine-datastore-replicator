PROJECT_IDS ?= crossregion-a crossregion-b crossregion-c

run:
	cd app && dev_appserver.py --port=8080 app.yaml

deploy:
	go install ./...
	cd app && $(foreach PID,$(PROJECT_IDS), \
		gcloud --project $(PID) app deploy -q && \
	) true

delete-old-versions:
	$(foreach PID,$(PROJECT_IDS), \
		gcloud --project $(PID) app versions list --filter="traffic_split=0" --format json \
		| jq -r ".[].id" \
		| xargs gcloud --project $(PID) app versions delete -q ; \
	)

deploy-indexes:
	cd app && $(foreach PID,$(PROJECT_IDS), \
		gcloud --project $(PID) datastore create-indexes index.yaml -q && \
	) true

deploy-replicator:
	$(MAKE) -C $(GOPATH)/src/4d63.com/google-cloud-appengine-datastore-replicator/service PROJECT_IDS="$(PROJECT_IDS)" deploy

create-pubsub:
	$(MAKE) -C $(GOPATH)/src/4d63.com/google-cloud-appengine-datastore-replicator/service PROJECT_IDS="$(PROJECT_IDS)" create-pubsub

delete-pubsub:
	$(MAKE) -C $(GOPATH)/src/4d63.com/google-cloud-appengine-datastore-replicator/service PROJECT_IDS="$(PROJECT_IDS)" delete-pubsub
