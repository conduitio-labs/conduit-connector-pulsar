// Copyright © 2024 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/matryer/is"
)

// named docker_compose_util.go so it stays close to the compose file, as it
// abstracts the docker-compose defined ports

const (
	PulsarURL    = "pulsar://127.0.0.1:6650"
	PulsarTLSURL = "pulsar+ssl://127.0.0.1:6651"
)

// SetupTopicName creates a new topic name for the test and deletes it if it
// exists, so that the test can start from a clean slate.
func SetupTopicName(t *testing.T, is *is.I) string {
	topic := "pulsar.topic." + t.Name()
	DeletePulsarTopic(is, topic)

	return topic
}

func DeletePulsarTopic(is *is.I, topic string) {
	serverURLs := []string{
		"http://127.0.0.1:8080",
		"http://127.0.0.1:8081",
	}

	// We delete the given topic on both servers, we assume that the topic
	// name is unique to the test itself

	for _, serverURL := range serverURLs {
		url := fmt.Sprintf(
			"%s/admin/v2/persistent/public/default/%s?force=true",
			serverURL, topic,
		)

		req, err := http.NewRequestWithContext(context.Background(), "DELETE", url, nil)
		is.NoErr(err)

		res, err := http.DefaultClient.Do(req)
		is.NoErr(err)
		res.Body.Close()

		// topic not found, nothing to delete
		if res.StatusCode == http.StatusNotFound {
			continue
		}

		is.Equal(res.StatusCode, http.StatusNoContent)
	}
}
