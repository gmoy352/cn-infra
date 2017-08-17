// Copyright (c) 2017 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugin

import (
	"github.com/ligato/cn-infra/datasync"
	"github.com/ligato/cn-infra/datasync/adapters"
	"github.com/ligato/cn-infra/datasync/persisted/dbsync"
	"github.com/ligato/cn-infra/db/keyval"
	"github.com/ligato/cn-infra/db/keyval/kvproto"
	"github.com/ligato/cn-infra/logging"
	"github.com/ligato/cn-infra/servicelabel"
	"github.com/ligato/cn-infra/utils/safeclose"
)

// Connection defines an access to a particular key-value data store implementation.
type Connection interface {
	keyval.CoreBrokerWatcher
}

// Skeleton of a KV plugin is a generic part of KV plugin.
type Skeleton struct {
	Transport    *adapters.TransportAggregator
	Logger       logging.Logger
	serviceLabel *servicelabel.Plugin
	name         string
	logFactory   logging.LogFactory
	conn         Connection
	protoWrapper *kvproto.ProtoWrapper
	connect      func(logger logging.Logger) (Connection, error)
}

// NewSkeleton creates a new instance of the Skeleton with the given connector.
// The connection is established in AfterInit phase.
func NewSkeleton(name string, transport *adapters.TransportAggregator, factory logging.LogFactory, serviceLabel *servicelabel.Plugin,
	connector func(log logging.Logger) (Connection, error)) *Skeleton {
	return &Skeleton{Transport: transport, serviceLabel: serviceLabel, name: name, logFactory: factory, connect: connector}
}

// Init is called on plugin startup
func (plugin *Skeleton) Init() (err error) {
	plugin.Logger, err = plugin.logFactory.NewLogger(plugin.name)
	if err != nil {
		return err
	}
	plugin.conn, err = plugin.connect(plugin.Logger)
	if err != nil {
		return err
	}
	plugin.protoWrapper = kvproto.NewProtoWrapperWithSerializer(plugin.conn, &keyval.SerializerJSON{})

	plugin.Transport.InitTransport(plugin.conn, plugin.serviceLabel, plugin.name)
	datasync.RegisterTransportOfDifferentAgent(func(microserviceLabel string) datasync.TransportAdapter {
		dbOfDifferentAgent := plugin.conn.NewBroker(plugin.serviceLabel.GetDifferentAgentPrefix(microserviceLabel))
		dbWOfDifferentAgent := plugin.conn.NewWatcher(plugin.serviceLabel.GetDifferentAgentPrefix(microserviceLabel))
		return dbsync.NewAdapter(microserviceLabel, dbOfDifferentAgent, dbWOfDifferentAgent)
	})

	return err
}

// AfterInit is called once all plugin have been initialized. The connection to datastore
// is established in this phase.
func (plugin *Skeleton) AfterInit() (err error) {
	return nil
}

// Close cleans up the resources
func (plugin *Skeleton) Close() error {
	return safeclose.Close(plugin.conn)
}

// NewBroker creates new instance of prefixed broker that provides API with arguments of type proto.Message
func (plugin *Skeleton) NewBroker(keyPrefix string) keyval.ProtoBroker {
	return plugin.protoWrapper.NewBroker(keyPrefix)
}

// NewWatcher creates new instance of prefixed broker that provides API with arguments of type proto.Message
func (plugin *Skeleton) NewWatcher(keyPrefix string) keyval.ProtoWatcher {
	return plugin.protoWrapper.NewWatcher(keyPrefix)
}
