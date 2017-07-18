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

package flavours

import (
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/db/keyval/etcdv3"
	"github.com/ligato/cn-infra/logging/logrus"
	"github.com/ligato/cn-infra/utils/config"
	"github.com/namsral/flag"
)

// Generic is set of common used generic plugins. This flavour can be used as a base
// for different flavours. The plugins are initialized in the same order as they appear
// in the structure.
type Generic struct {
	etcdConfigFile string
	etcdCfg        etcdv3.Config

	Lg   *logrus.Plugin
	Etcd *etcdv3.PluginEtcd
}

// RegisterFlags registers the options that need to be parsed.
func (f *Generic) RegisterFlags() {
	flag.StringVar(&f.etcdConfigFile, "etcdv3-config", "", "Location of the Etcd configuration file; also set via 'ETCDV3_CONFIG' env variable.")
}

// ApplyConfig loads the config and creates the plugins.
func (f *Generic) ApplyConfig() error {
	// config Parsing
	if f.etcdConfigFile != "" {
		err := config.ParseConfigFromYamlFile(f.etcdConfigFile, &f.etcdCfg)
		if err != nil {
			return err
		}
	}

	// call the constructors
	f.Lg = logrus.NewLogrusPlugin()
	f.Etcd = etcdv3.NewEtcdPlugin(&f.etcdCfg)

	return nil
}

// Inject interconnects plugins - inject the dependencies
func (f *Generic) Inject() error {

	f.Etcd.Lg = f.Lg

	return nil
}

// Plugins returns all plugins from the flavour. The set of plugins is supposed to be passed to the agent constructor.
func (f *Generic) Plugins() []*core.NamedPlugin {
	return core.Named(core.ListPluginsInFlavor(f))
}