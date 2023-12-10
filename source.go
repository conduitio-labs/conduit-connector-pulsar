package apachepulsar

//go:generate paramgen -output=paramgen_src.go SourceConfig

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/apache/pulsar-client-go/pulsar"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Source struct {
	sdk.UnimplementedSource

	consumer pulsar.Consumer
	received []pulsar.Message
	mx       *sync.Mutex

	config SourceConfig
}

type SourceConfig struct {
	Config

	URL          string `json:"URL" validate:"required"`
	Topic        string `json:"topic" validate:"required"`
	Subscription string `json:"subscription" validate:"required"`
}

func NewSource() sdk.Source {
	return sdk.SourceWithMiddleware(&Source{mx: &sync.Mutex{}}, sdk.DefaultSourceMiddleware()...)
}

func (s *Source) Parameters() map[string]sdk.Parameter {
	return s.config.Parameters()
}

func (s *Source) Configure(ctx context.Context, cfg map[string]string) error {
	sdk.Logger(ctx).Info().Msg("Configuring Source...")
	err := sdk.Util.ParseConfig(cfg, &s.config)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	
	return nil
}

func (s *Source) Open(ctx context.Context, pos sdk.Position) error {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: s.config.URL,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	s.consumer, err = client.Subscribe(pulsar.ConsumerOptions{
		Topic:            s.config.Topic,
		SubscriptionName: s.config.Subscription,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	return nil
}

func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
	msg, err := s.consumer.Receive(ctx)
	if err != nil {
		return sdk.Record{}, fmt.Errorf("failed to receive message: %w", err)
	}

	s.mx.Lock()
	s.received = append(s.received, msg)
	s.mx.Unlock()

	var position sdk.Position
	var metadata sdk.Metadata
	var key sdk.Data
	var payload sdk.Data = sdk.RawData(msg.Payload())

	return sdk.Util.Source.NewRecordCreate(
		position,
		metadata,
		key,
		payload,
	), nil
}

func (s *Source) Ack(ctx context.Context, _ sdk.Position) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	var err error
	for _, msg := range s.received {
		ackErr := s.consumer.Ack(msg)
		if err != nil {
			err = errors.Join(err, fmt.Errorf("failed to ack message: %w", ackErr))
		}
	}

	return err
}

func (s *Source) Teardown(ctx context.Context) error {
	if s.consumer != nil {
		s.consumer.Close()
	}
	return nil
}
