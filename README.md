# google-cloud-appengine-datastore-replicator

A package and AppEngine Standard service that can be used together to replicate Datastore Puts from one AppEngine project to another. The primary use case is to replicate data across regions.

Replication retains namespace and full parent keys.

There is no coordination to determine latest write when race conditions occur. Recommended only for use with Datastore entities that are write-once.

## Status

Experimental

## Usage

Use the package to publish Datastore Puts, and deploy the service to all regions where the data should be replicated. Use the Makefile to setup the cross region PubSub topics and subscriptions.

### 1. Publish Datastore Puts to other regions

```go
import "4d63.com/google-cloud-appengine-datastore-replicator"
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

### 2. Subscribe to Datastore Puts and replicate them

#### 2a. Use the replicator service

```
cd service
```

##### Create the PubSub topics and subscriptions

```
make create-pubsub
```

##### Deploy replicator service

```
make deploy
```

##### Delete the PubSub topics and subscriptions

```
make delete-pubsub
```

#### 2b. Use the package

Take a look at [service/app.go](service/app.go) for how to consume the PubSub request, unpack it, then put it to the local Datastore.

## Example

See [example](example).
