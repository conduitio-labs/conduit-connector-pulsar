package apachepulsar_test

import (
	"context"
	"sync"
	"testing"

	apachepulsar "github.com/alarbada/conduit-connector-apachepulsar"
	"github.com/apache/pulsar-client-go/pulsar"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
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

	topic := "destination_test"

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
	con := apachepulsar.NewDestination()

	ctx := context.Background()

	err := con.Configure(ctx, map[string]string{
		"URL":          "pulsar://localhost:6650",
		"topic":        topic,
		"subscription": topic,
	})
	is.NoErr(err)

	err = con.Open(ctx)
	is.NoErr(err)

	rec := sdk.Util.Source.NewRecordCreate(
		[]byte(uuid.NewString()),
		sdk.Metadata{"pulsar.topic": topic},
		sdk.RawData("test-key"),
		sdk.RawData(exampleMessage),
	)
	written, err := con.Write(ctx, []sdk.Record{rec})
	is.NoErr(err)
	is.Equal(written, 1)
}

func pulsarGoClientRead(is *is.I, topic string) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
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

	msgContents := string(msg.Payload())

	is.Equal(msgContents, exampleMessage)
}

func TestDestinationConfiguration(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	{
		con := apachepulsar.NewDestination()
		err := con.Configure(ctx, map[string]string{})
		is.True(err != nil)
	}

	{
		con := apachepulsar.NewDestination()
		err := con.Configure(ctx, map[string]string{
			"URL":              "pulsar://localhost:6650",
			"topic":            "",
		})
		is.True(err != nil)
	}

	{
		con := apachepulsar.NewDestination()
		err := con.Configure(ctx, map[string]string{
			"URL":              "pulsar://localhost:6650",
			"topic":            "test",
		})
		is.True(err == nil)
	}
}
