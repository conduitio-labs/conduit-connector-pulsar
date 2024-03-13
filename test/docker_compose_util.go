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
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

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
	topic = strings.ReplaceAll(topic, "/", "_")
	DeletePulsarTopic(is, topic)

	return topic
}

func DeletePulsarTopic(is *is.I, topic string) {
	// Use a context with timeout to avoid hanging requests.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel() // Ensure resources are freed even in case of early return.
	url := fmt.Sprintf(
		"http://127.0.0.1:8080/admin/v2/persistent/public/default/%s?force=true",
		topic)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	is.NoErr(err)

	res, err := http.DefaultClient.Do(req)
	is.NoErr(err)
	defer res.Body.Close()

	// Read and discard the response body to ensure the connection can be reused.
	_, err = io.Copy(io.Discard, res.Body)
	is.NoErr(err)

	// topic not found, nothing to delete
	if res.StatusCode == http.StatusNotFound {
		return
	}

	is.Equal(res.StatusCode, http.StatusNoContent)
}
