package replicator // import "4d63.com/google-cloud-appengine-datastore-replicator"

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"golang.org/x/oauth2/google"
	pubsub "google.golang.org/api/pubsub/v1"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const (
	attributeKeySourceAppID          = "source-app-id"
	attributeKeySourceDatacenterID   = "source-datacenter-id"
	attributeKeySourceInstanceID     = "source-instance-id"
	attributeKeySourceServiceName    = "source-service-name"
	attributeKeySourceServiceVersion = "source-service-version"
	attributeKeySourceRequestID      = "source-request-id"
)

func init() {
	gob.Register(time.Time{})
}

type Source struct {
	AppID          string
	DatacenterID   string
	InstanceID     string
	ServiceName    string
	ServiceVersion string
	RequestID      string
}

func encodeEntity(c context.Context, key *datastore.Key, src interface{}) ([]byte, error) {
	log.Infof(c, "Encoding %s %#v", key, src)

	var entity datastore.Entity
	var err error
	entity.Key = key
	if saver, ok := src.(datastore.PropertyLoadSaver); ok {
		entity.Properties, err = saver.Save()
	} else {
		entity.Properties, err = datastore.SaveStruct(src)
	}
	if err != nil {
		return nil, err
	}

	log.Infof(c, "Encoded entity %#v %#v", entity.Key, entity.Properties)

	buf := bytes.Buffer{}
	err = gob.NewEncoder(&buf).Encode(&entity)
	if err != nil {
		return nil, err
	}
	data := buf.Bytes()

	return data, nil
}

func decodeEntity(c context.Context, data []byte) (*datastore.Key, *datastore.PropertyList, error) {
	var entity datastore.Entity
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(&entity)
	if err != nil {
		return nil, nil, err
	}

	key := entity.Key
	props := datastore.PropertyList(entity.Properties)

	log.Infof(c, "Decoded entity %#v %#v", key, &props)

	return key, &props, nil
}

func Publish(c context.Context, topic string, key *datastore.Key, src interface{}) error {
	log.Infof(c, "Publishing to %s entity %#v %#v", topic, key, src)

	data, err := encodeEntity(c, key, src)
	if err != nil {
		return err
	}

	hc, err := google.DefaultClient(c, pubsub.PubsubScope)
	if err != nil {
		return err
	}

	ps, err := pubsub.New(hc)
	if err != nil {
		return err
	}

	fullTopic := fmt.Sprintf("projects/%s/topics/%s", appengine.AppID(c), topic)

	request := &pubsub.PublishRequest{
		Messages: []*pubsub.PubsubMessage{
			{
				Attributes: map[string]string{
					attributeKeySourceAppID:          appengine.AppID(c),
					attributeKeySourceDatacenterID:   appengine.Datacenter(c),
					attributeKeySourceInstanceID:     appengine.InstanceID(),
					attributeKeySourceServiceName:    appengine.ModuleName(c),
					attributeKeySourceServiceVersion: appengine.VersionID(c),
					attributeKeySourceRequestID:      appengine.RequestID(c),
				},
				Data: base64.StdEncoding.EncodeToString(data),
			},
		},
	}

	log.Infof(c, "Publishing %s %#v", fullTopic, request)

	result, err := ps.Projects.Topics.Publish(fullTopic, request).Do()
	if err != nil {
		return err
	}

	log.Infof(c, "Published %s %s %#v", fullTopic, result.MessageIds[0], request)
	return nil
}

func Unpack(c context.Context, r io.Reader) (source *Source, key *datastore.Key, props *datastore.PropertyList, err error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil, nil, err
	}

	log.Infof(c, "Unpacking raw: %s", raw)

	var body struct {
		Message struct {
			Attributes map[string]string
			Data       string
		}
	}

	err = json.Unmarshal(raw, &body)
	if err != nil {
		return nil, nil, nil, err
	}

	log.Infof(c, "Unpacked attributes: %v", body.Message.Attributes)

	source = &Source{
		AppID:          body.Message.Attributes[attributeKeySourceAppID],
		DatacenterID:   body.Message.Attributes[attributeKeySourceDatacenterID],
		InstanceID:     body.Message.Attributes[attributeKeySourceInstanceID],
		ServiceName:    body.Message.Attributes[attributeKeySourceServiceName],
		ServiceVersion: body.Message.Attributes[attributeKeySourceServiceVersion],
		RequestID:      body.Message.Attributes[attributeKeySourceRequestID],
	}

	log.Infof(c, "Unpacked source: %#v", source)

	log.Infof(c, "Unpacked data: %s", body.Message.Data)

	data, err := base64.StdEncoding.DecodeString(body.Message.Data)
	if err != nil {
		return nil, nil, nil, err
	}

	key, props, err = decodeEntity(c, data)

	log.Infof(c, "Unpacked entity %#v %#v", key, props)

	localKey, err := localizeKey(c, key)
	if err != nil {
		return nil, nil, nil, err
	}

	log.Infof(c, "Unpacked local entity %#v %#v", localKey, props)

	return source, localKey, props, err
}

func localizeKey(c context.Context, key *datastore.Key) (*datastore.Key, error) {
	var localParent *datastore.Key
	if parent := key.Parent(); parent != nil {
		var err error
		localParent, err = localizeKey(c, parent)
		if err != nil {
			return nil, err
		}
	}

	namespacedCtx, err := appengine.Namespace(c, key.Namespace())
	if err != nil {
		return nil, err
	}

	return datastore.NewKey(namespacedCtx, key.Kind(), key.StringID(), key.IntID(), localParent), nil
}
