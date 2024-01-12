package apachepulsar_test

import (
	"context"
	"testing"
	"time"

	apachepulsar "github.com/alarbada/conduit-connector-apachepulsar"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/matryer/is"
)

func TestTeardownSource_NoOpen(t *testing.T) {
	is := is.New(t)
	con := apachepulsar.NewSource()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}

func TestSource_Integration(t *testing.T) {
	is := is.New(t)

	topic := "source_test"
	con := apachepulsar.NewSource()
	ctx := context.Background()
	err := con.Configure(ctx, map[string]string{
		"URL":              "pulsar://localhost:6650",
		"topic":            topic,
		"subscriptionName": "source_test",
	})
	is.NoErr(err)

	err = con.Open(ctx, nil)
	is.NoErr(err)

	defer func() {
		err := con.Teardown(ctx)
		is.NoErr(err)
	}()

	readEndC := make(chan struct{})

	go func() {
		defer close(readEndC)

		record, err := con.Read(ctx)
		is.NoErr(err)
		is.Equal(string(record.Payload.After.Bytes()), "example message")
		err = con.Ack(ctx, record.Position)
		is.NoErr(err)
	}()

	produceExampleMsg(is, topic)

	select {
	case <-readEndC:
		// Message has been correctly read
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out waiting for message consumption")
	}
}

func produceExampleMsg(is *is.I, topic string) {
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

	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: []byte("example message"),
	})
	is.NoErr(err)
}
