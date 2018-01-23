SERVICE_NAME ?= datastore-replicator
PROJECT_IDS ?= proj-a proj-b
TOPIC_NAME ?= datastore-puts

run:
	go install ./...
	dev_appserver.py --port=8080 app.yaml

deploy:
	go install ./...
	$(foreach PID,$(PROJECT_IDS), \
		gcloud --project $(PID) app deploy -q && \
	) true

create-pubsub:
	$(foreach PID,$(PROJECT_IDS), \
		gcloud --project $(PID) pubsub topics create $(TOPIC_NAME) && \
	) true
	$(foreach PID_PUB,$(PROJECT_IDS), \
		$(foreach PID_SUB,$(PROJECT_IDS), \
			gcloud --project $(PID_SUB) pubsub subscriptions create pub-$(PID_PUB)-sub-$(PID_SUB) \
				--topic-project=$(PID_PUB) --topic=$(TOPIC_NAME) \
				--push-endpoint=https://$(SERVICE_NAME)-dot-$(PID_SUB).appspot.com/_ah/push-handlers/datastore-puts && \
		) \
	) true

delete-pubsub:
	$(foreach PID,$(PROJECT_IDS), \
		gcloud --project $(PID) pubsub topics delete $(TOPIC_NAME) && \
	) true
	$(foreach PID_PUB,$(PROJECT_IDS), \
		$(foreach PID_SUB,$(PROJECT_IDS), \
			gcloud --project $(PID_SUB) pubsub subscriptions delete pub-$(PID_PUB)-sub-$(PID_SUB) && \
		) \
	) true
