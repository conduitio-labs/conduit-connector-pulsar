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
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/conduitio-labs/conduit-connector-pulsar/test"
	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestTeardown_NoOpen(t *testing.T) {
	is := is.New(t)
	con := NewDestination()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}

func TestDestination_Integration(t *testing.T) {
	is := is.New(t)

	topic := test.SetupTopicName(t, is)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		connectorDestinationWrite(is, topic)
	}()

	go func() {
		defer wg.Done()

		pulsarGoClientRead(is, topic)
	}()

	wg.Wait()
}

var exampleMessage = "example message"

func connectorDestinationWrite(is *is.I, topic string) {
	con := NewDestination()

	ctx := context.Background()

	cfgMap := map[string]string{
		"url":   test.PulsarURL,
		"topic": topic,
	}

	err := sdk.Util.ParseConfig(ctx, cfgMap, con.Config(), Connector.NewSpecification().DestinationParams)
	is.NoErr(err)

	is.NoErr(con.Open(ctx))

	rec := sdk.Util.Source.NewRecordCreate(
		[]byte(uuid.NewString()),
		opencdc.Metadata{"pulsar.topic": topic},
		opencdc.RawData("test-key"),
		opencdc.RawData(exampleMessage),
	)

	written, err := con.Write(ctx, []opencdc.Record{rec})
	is.NoErr(err)
	is.Equal(written, 1)
}

func pulsarGoClientRead(is *is.I, topic string) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: test.PulsarURL,
	})
	is.NoErr(err)
	defer client.Close()

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: topic,
	})
	is.NoErr(err)

	msg, err := consumer.Receive(context.Background())
	is.NoErr(err)

	var received struct {
		Payload struct {
			After opencdc.RawData `json:"after"`
		} `json:"payload"`
	}
	err = json.Unmarshal(msg.Payload(), &received)
	is.NoErr(err)

	receivedMsg := string(received.Payload.After.Bytes())
	is.Equal(receivedMsg, exampleMessage)
}
