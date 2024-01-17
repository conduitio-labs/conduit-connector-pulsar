package pulsar

//go:generate paramgen -output=paramgen_src.go SourceConfig

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/pulsar-client-go/pulsar"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Source struct {
	sdk.UnimplementedSource

	consumer pulsar.Consumer
	received map[string]pulsar.Message
	mx       *sync.Mutex
	config   SourceConfig
}

func NewSource() sdk.Source {
	source := &Source{
		mx:       &sync.Mutex{},
		received: make(map[string]pulsar.Message),
	}

	return sdk.SourceWithMiddleware(source, sdk.DefaultSourceMiddleware()...)
}

func (s *Source) Parameters() map[string]sdk.Parameter {
	return s.config.Parameters()
}

func (s *Source) Configure(ctx context.Context, cfg map[string]string) error {
	sdk.Logger(ctx).Info().Msg("Configuring Source...")

	if err := sdk.Util.ParseConfig(cfg, &s.config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}

func (s *Source) Open(_ context.Context, pos sdk.Position) error {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:                     s.config.URL,
		ConnectionTimeout:       s.config.ConnectionTimeout,
		OperationTimeout:        s.config.OperationTimeout,
		MaxConnectionsPerBroker: s.config.MaxConnectionsPerBroker,
		MemoryLimitBytes:        s.config.MemoryLimitBytes,
		EnableTransaction:       s.config.EnableTransaction,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	subscriptionType := parseSubscriptionType(s.config.SubscriptionType)

	s.consumer, err = client.Subscribe(pulsar.ConsumerOptions{
		Topic:            s.config.Topic,
		SubscriptionName: s.config.SubscriptionName,
		Type:             subscriptionType,
	})
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	if pos != nil {
		messageID, err := pulsar.DeserializeMessageID(pos)
		if err != nil {
			return fmt.Errorf(
				"failed to deserialize message ID while trying to seek the given position %s: %w",
				string(pos), err,
			)
		}
		if err := s.consumer.Seek(messageID); err != nil {
			return fmt.Errorf(
				"failed to seek to the given position %s: %w",
				string(pos), err,
			)
		}
	}

	return nil
}

const (
	// MetadataPulsarTopic is the metadata key for storing the pulsar topic
	MetadataPulsarTopic = "pulsar.topic"
)

func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
	sdk.Logger(ctx).Debug().Msg("reading message")
	msg, err := s.consumer.Receive(ctx)
	if err != nil {
		return sdk.Record{}, fmt.Errorf("failed to receive message: %w", err)
	}

	s.mx.Lock()
	s.received[msg.ID().String()] = msg
	s.mx.Unlock()

	position := sdk.Position(msg.ID().Serialize())

	sdk.Logger(ctx).Debug().Str("MessageID", string(position)).Msg("Setting position for message")

	metadata := sdk.Metadata{MetadataPulsarTopic: msg.Topic()}
	metadata.SetCreatedAt(msg.EventTime())

	key := sdk.RawData(msg.Key())
	payload := sdk.RawData(msg.Payload())

	return sdk.Util.Source.NewRecordCreate(
		position,
		metadata,
		key,
		payload,
	), nil
}

func (s *Source) Ack(ctx context.Context, position sdk.Position) error {
	sdk.Logger(ctx).Debug().Str("MessageID", string(position)).Msg("Attempting to ack message")

	msgID, err := pulsar.DeserializeMessageID(position)
	if err != nil {
		return fmt.Errorf("failed to deserialize message ID: %w", err)
	}

	s.mx.Lock()
	defer s.mx.Unlock()
	msg, ok := s.received[msgID.String()]
	if ok {
		delete(s.received, msgID.String())
		return s.consumer.Ack(msg)
	}

	return fmt.Errorf("message not found for position: %s", string(position))
}

func (s *Source) Teardown(_ context.Context) error {
	if s.consumer != nil {
		s.consumer.Close()
	}
	return nil
}
