// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World!")
	})
	http.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		env := "local"
		if value, ok := os.LookupEnv("APP_ENVIRONMENT"); ok {
			env = value
		}
		fmt.Fprint(w, env)
	})
	http.HandleFunc("/release", func(w http.ResponseWriter, r *http.Request) {
		release := "unknown"
		if value, ok := os.LookupEnv("RELEASE_ID"); ok {
			release = value
		}
		fmt.Fprint(w, release)
	})
	http.HandleFunc("/failed", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failed!")
	})
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8080", nil))
}
