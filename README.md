# google-cloud-appengine-datastore-replicator

A package and AppEngine Standard service that can be used together to replicate Datastore Puts from one AppEngine project to another. The primary use case is to replicate data across regions.

Replication retains namespace and full parent keys.

There is no coordination to determine latest write when race conditions occur. Recommended only for use with Datastore entities that are write-once.

## Status

Experimental

## Requirements

* Google AppEngine (Standard)
* Google Cloud Datastore
* Google Cloud PubSub

## Usage

Use the Makefile to setup the cross region PubSub topics and subscriptions. Use the package to publish Datastore Puts, and deploy the service to all regions where the data should be replicated.

### 1. Create the PubSub topics and subscriptions

```
cd service
make PROJECT_IDS="<project-id-1> <project-id-2> ..." create-pubsub
```

### 2. Publish Datastore Puts to other regions

```go
import replicator "4d63.com/google-cloud-appengine-datastore-replicator"
```

```go
// Put the entity as you normally would.
err := datastore.Put(ctx, key, entity)
if err != nil {
	return err
}

// Publish the key and entity that was put.
err = replicator.Publish(ctx, "datastore-puts", key, entity)
if err != nil {
	return err
}
```

Recommendations and strategies:

- Generate IDs for keys yourself and use complete keys, don't use the Datastore Allocator to generate IDs by putting entities with an incomplete key. If you use the allocator it may generate overlapping IDs in different projects causing loss of existing entities.
- Write an entity once to a key, and write updates as new entities. Resolve conflicts when reading.

### 3. Subscribe to Datastore Puts and replicate them

Use the replicator service to receive pushes from the subscriptions setup in step 1. The replicator service will store the entities into its local project respecting the same namespace and key, including the parent keys.

```
cd service
make PROJECT_IDS="<project-id-1> <project-id-2> ..." deploy
```

Alternatively if you'd like to create your own service to store the entities, take a look at [service/app.go](service/app.go) for how to consume the PubSub request, unpack it, then put it to the local Datastore.

## Example

See [example](example).
