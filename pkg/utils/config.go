/*
Copyright © 2021 Madhav Jivrajani madhav.jiv@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the top level configuration file
// that will be read.
type Config struct {
	// PathToBinary is the absolute path
	// to the binary for which scheduler
	// traces need to be collected.
	PathToBinary string `yaml:"path,omitempty"`
	// SchedTrace is the specific configuration
	// for collecting scheduler traces.
	SchedTrace SchedTraceConfig `yaml:"sched,omitempty"`
	// TODO: add prometheus endpoint related config.
}

// SchedTraceConfig represents the config to configure
// how the scheduler traces will be collected.
type SchedTraceConfig struct {
	// Interval is the frequency at which traces will be
	// collected, and this value should be in milliseconds.
	Interval uint64 `yaml:"interval"`
}

// ReadConfig reads the provided config into the Config
// type and returns an error if any.
func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
