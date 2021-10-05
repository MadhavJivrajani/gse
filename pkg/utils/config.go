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

	"github.com/MadhavJivrajani/gse/pkg/core"
	"gopkg.in/yaml.v2"
)

// ReadConfig reads the provided config into the Config
// type and returns an error if any.
func ReadConfig(path string) (*core.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &core.Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
