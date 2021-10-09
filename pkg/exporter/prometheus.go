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

package exporter

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MadhavJivrajani/gse/pkg/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ServeMetrics serves promethues metrics on a particular endpoint
// and a particular port.
func ServeMetrics(config utils.PrometheusConfig) {
	log.Printf("serving metrics on endpoint %s and port %d", config.Endpoint, config.Port)
	http.Handle(fmt.Sprintf("/%s", config.Endpoint), promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
