package apachepulsar_test

import (
	"context"
	"testing"

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
	healthcheck(t)

	is := is.New(t)

	produceExampleMsg(is)

	con := apachepulsar.NewSource()
	ctx := context.Background()

	err := con.Configure(ctx, map[string]string{
		"URL":              "pulsar://localhost:6650",
		"topic":            "source_test",
		"subscriptionName": "source_test",
	})
	is.NoErr(err)

	err = con.Open(ctx, nil)
	is.NoErr(err)

	defer func() {
		err := con.Teardown(ctx)
		is.NoErr(err)
	}()

	record, err := con.Read(ctx)
	is.NoErr(err)
	is.Equal(string(record.Payload.After.Bytes()), "example message")

	err = con.Ack(ctx, record.Position)
	is.NoErr(err)
}

func produceExampleMsg(is *is.I) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	is.NoErr(err)

	defer client.Close()

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "source_test",
	})
	is.NoErr(err)

	defer producer.Close()

	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: []byte("example message"),
	})
	is.NoErr(err)
}
