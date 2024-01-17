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
	"fmt"
	"testing"

	"github.com/apache/pulsar-client-go/pulsar"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/matryer/is"
)

func TestTeardownSource_NoOpen(t *testing.T) {
	is := is.New(t)
	con := NewSource()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}

func TestSource_Integration_RestartFull(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	topic := "pulsar.topic." + t.Name()
	cfgMap := map[string]string{
		"URL":              "pulsar://localhost:6650",
		"topic":            topic,
		"subscriptionName": topic + "-subscription",
	}

	recs1 := generatePulsarMsgs(1, 3)
	go producePulsarMsgs(is, topic, recs1)
	lastPosition := testSourceIntegrationRead(is, cfgMap, nil, recs1, false)

	recs2 := generatePulsarMsgs(4, 6)
	go producePulsarMsgs(is, topic, recs2)
	testSourceIntegrationRead(is, cfgMap, lastPosition, recs2, false)
}

func TestSource_Integration_RestartPartial(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	topic := "pulsar.topic." + t.Name()
	cfgMap := map[string]string{
		"URL":              "pulsar://localhost:6650",
		"topic":            topic,
		"subscriptionName": topic + "-subscription",
	}

	recs1 := generatePulsarMsgs(1, 3)
	go producePulsarMsgs(is, topic, recs1)

	lastPosition := testSourceIntegrationRead(is, cfgMap, nil, recs1, true)

	// only first record was acked, produce more records and expect to resume
	// from last acked record
	recs2 := generatePulsarMsgs(4, 6)
	go producePulsarMsgs(is, topic, recs2)

	var wantRecs []*pulsar.ProducerMessage
	wantRecs = append(wantRecs, recs1[1:]...)
	wantRecs = append(wantRecs, recs2...)
	testSourceIntegrationRead(is, cfgMap, lastPosition, wantRecs, false)
}

func generatePulsarMsgs(from, to int) []*pulsar.ProducerMessage {
	var msgs []*pulsar.ProducerMessage

	for i := from; i <= to; i++ {
		msgs = append(msgs, &pulsar.ProducerMessage{
			Key:     fmt.Sprintf("test-key-%d", i),
			Payload: []byte(fmt.Sprintf("test-key-%d", i)),
		})
	}

	return msgs
}

func producePulsarMsgs(is *is.I, topic string, msgs []*pulsar.ProducerMessage) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	is.NoErr(err)
	defer client.Close()

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	is.NoErr(err)
	defer producer.Close()

	for _, msg := range msgs {
		_, err = producer.Send(context.Background(), msg)
		is.NoErr(err)
	}
}

// testSourceIntegrationRead reads and acks messages in range [from,to].
// If ackFirst is true, only the first message will be acknowledged.
// Returns the position of the last message read.
func testSourceIntegrationRead(
	is *is.I,
	cfgMap map[string]string,
	startFrom sdk.Position,
	wantRecords []*pulsar.ProducerMessage,
	ackFirstOnly bool,
) sdk.Position {
	ctx := context.Background()

	underTest := NewSource()
	defer func() {
		err := underTest.Teardown(ctx)
		is.NoErr(err)
	}()

	err := underTest.Configure(ctx, cfgMap)
	is.NoErr(err)
	err = underTest.Open(ctx, startFrom)
	is.NoErr(err)

	var positions []sdk.Position
	for _, wantRecord := range wantRecords {
		rec, err := underTest.Read(ctx)
		is.NoErr(err)
		is.Equal(wantRecord.Key, string(rec.Key.Bytes()))

		positions = append(positions, rec.Position)
	}

	for i, p := range positions {
		if i > 0 && ackFirstOnly {
			break
		}
		err = underTest.Ack(ctx, p)
		is.NoErr(err)
	}

	return positions[len(positions)-1]
}
