# example

A demonstration of how to use Google Cloud PubSub across multiple projects, and to replicate data across regional products like Google Cloud Datastore.

Deployed at:

* https://crossregion-a.appspot.com/ (us-east-1)
* https://crossregion-b.appspot.com/ (australia-southeast1)
* https://crossregion-c.appspot.com/ (europe-west2)

## Status

WIP

## Usage

1. Create two or more projects, and activate appengine, selecting different regions for each.

2. Setup PubSub topics and subscriptions

   ```
   make PROJECT_IDS="<project-id-1> <project-id-2> ..." create-pubsub
   ```

3. Deploy the demo app and the replicator service to the projects.

   ```
   make PROJECT_IDS="<project-id-1> <project-id-2> ..." deploy deploy-indexes deploy-replicator
   ```

