package apachepulsar_test

import (
	"context"
	"testing"

	apachepulsar "github.com/alarbada/conduit-connector-apachepulsar"
	"github.com/apache/pulsar-client-go/pulsar"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/matryer/is"
)

func TestTeardown_NoOpen(t *testing.T) {
	is := is.New(t)
	con := apachepulsar.NewDestination()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}

func TestDestination_Integration(t *testing.T) {
	is := is.New(t)

	con := apachepulsar.NewDestination()
	ctx := context.Background()

	err := con.Configure(ctx, map[string]string{
		"URL":          "pulsar://localhost:6650",
		"topic":        "test",
		"subscription": "test",
	})
	is.NoErr(err)

	err = con.Open(ctx)
	is.NoErr(err)

	defer func() {
		err := con.Teardown(ctx)
		is.NoErr(err)
	}()

	_, err = con.Write(ctx, []sdk.Record{
		{
			Position:  []byte{},
			Operation: 0,
			Metadata:  map[string]string{},
			Key:       nil,
			Payload: sdk.Change{
				Before: nil,
				After:  sdk.RawData("example message"),
			},
		},
	})
	is.NoErr(err)

	{ // consume the test message
		client, err := pulsar.NewClient(pulsar.ClientOptions{
			URL: "pulsar://localhost:6650",
		})
		is.NoErr(err)
		defer client.Close()

		consumer, err := client.Subscribe(pulsar.ConsumerOptions{
			Topic:            "test",
			SubscriptionName: "test",
		})
		is.NoErr(err)

		msg, err := consumer.Receive(context.Background())
		is.NoErr(err)

		is.Equal(msg.Payload(), []byte("example message"))
	}
}
