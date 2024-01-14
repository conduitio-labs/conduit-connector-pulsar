package apachepulsar_test

import (
	"context"
	"fmt"
	"sync"
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

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		connectorDestinationWrite(is)
	}()

	go func() {
		defer wg.Done()

		pulsarGoClientRead(is)
	}()

	wg.Wait()
}

func connectorDestinationWrite(is *is.I) {
	con := apachepulsar.NewDestination()

	ctx := context.Background()

	err := con.Configure(ctx, map[string]string{
		"URL":          "pulsar://localhost:6650",
		"topic":        "destination_test",
		"subscription": "destination_test",
	})
	is.NoErr(err)

	err = con.Open(ctx)
	is.NoErr(err)

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
}

func pulsarGoClientRead(is *is.I) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	is.NoErr(err)
	defer client.Close()

	fmt.Println("creating consumer")

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            "destination_test",
		SubscriptionName: "destination_test",
	})
	is.NoErr(err)

	fmt.Println("receiving message")

	msg, err := consumer.Receive(context.Background())
	is.NoErr(err)

	fmt.Println(string(msg.Payload()))
}
