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

package pulsar

import (
	"testing"
	"time"

	"github.com/conduitio-labs/conduit-connector-apache-pulsar/test"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/matryer/is"
	"go.uber.org/goleak"
)

func init() {
	// Uncomment this up a logger for tests to use. By default sdk.Logger log calls won't
	// output anything
	// log := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// log = log.Level(zerolog.InfoLevel)
	// zerolog.DefaultContextLogger = &log
}

func TestAcceptance(t *testing.T) {
	is := is.New(t)

	topic := test.SetupTopicName(t, is)
	cfg := map[string]string{
		"url":              test.PulsarURL,
		"topic":            topic,
		"subscriptionName": "test-subscription",
	}

	sdk.AcceptanceTest(t, sdk.ConfigurableAcceptanceTestDriver{
		Config: sdk.ConfigurableAcceptanceTestDriverConfig{
			Connector:         Connector,
			GoleakOptions:     []goleak.Option{goleak.IgnoreCurrent()},
			SourceConfig:      cfg,
			DestinationConfig: cfg,
			BeforeTest: func(t *testing.T) {
				topic := test.SetupTopicName(t, is)
				cfg["topic"] = topic
			},
			Skip: []string{
				"TestSource_Configure_RequiredParams",
				"TestDestination_Configure_RequiredParams",
			},
			WriteTimeout: 500 * time.Millisecond,
			ReadTimeout:  500 * time.Millisecond,
		},
	})
}
