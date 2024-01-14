package apachepulsar

//go:generate paramgen -output=paramgen_dest.go DestinationConfig

import (
	"context"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Destination struct {
	sdk.UnimplementedDestination

	producer pulsar.Producer

	config DestinationConfig
}

type DestinationConfig struct {
	URL   string `json:"URL" validate:"required"`
	Topic string `json:"topic" validate:"required"`
}

func NewDestination() sdk.Destination {
	return sdk.DestinationWithMiddleware(&Destination{}, sdk.DefaultDestinationMiddleware()...)
}

func (d *Destination) Parameters() map[string]sdk.Parameter {
	return d.config.Parameters()
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) error {
	sdk.Logger(ctx).Info().Msg("Configuring Destination...")

	validate := newConfigValidator(cfg)
	validatedConfig := DestinationConfig{
		URL:   validate.Required("URL"),
		Topic: validate.Required("topic"),
	}
	d.config = validatedConfig
	if err := validate.Error(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	return nil
}

func (d *Destination) Open(ctx context.Context) error {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: d.config.URL,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	d.producer, err = client.CreateProducer(pulsar.ProducerOptions{
		Topic: d.config.Topic,
	})
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}

	return nil
}

func (d *Destination) Write(ctx context.Context, records []sdk.Record) (int, error) {
	var written int
	for _, record := range records {
		sdk.Logger(ctx).Debug().Msgf("Writing record")
		key := string(record.Key.Bytes())

		_, err := d.producer.Send(ctx, &pulsar.ProducerMessage{
			Payload: record.Payload.After.Bytes(),
			Key:     key,
		})
		if err != nil {
			return written, fmt.Errorf("failed to send message: %w", err)
		}

		sdk.Logger(ctx).Debug().Msgf("Successfully wrote record")
		written++
	}

	return written, nil
}

func (d *Destination) Teardown(ctx context.Context) error {
	if d.producer != nil {
		d.producer.Close()
	}

	return nil
}
