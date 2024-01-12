package apachepulsar

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

	config SourceConfig
}

type SourceConfig struct {
	URL              string           `json:"URL" validate:"required"`
	Topic            string           `json:"topic" validate:"required"`
	SubscriptionName string           `json:"subscriptionName" validate:"required"`
	SubscriptionType SubscriptionType `json:"subscriptionType"`
}

type SubscriptionType string

const (
	// Exclusive there can be only 1 consumer on the same topic with the same subscription name
	Exclusive SubscriptionType = "exclusive"

	// Shared subscription mode, multiple consumer will be able to use the same subscription name
	// and the messages will be dispatched according to
	// a round-robin rotation between the connected consumers
	Shared SubscriptionType = "shared"

	// Failover subscription mode, multiple consumer will be able to use the same subscription name
	// but only 1 consumer will receive the messages.
	// If that consumer disconnects, one of the other connected consumers will start receiving messages.
	Failover SubscriptionType = "failover"

	// KeyShared subscription mode, multiple consumer will be able to use the same
	// subscription and all messages with the same key will be dispatched to only one consumer
	KeyShared SubscriptionType = "key_shared"
)

func ParseSubscriptionType(s string) (SubscriptionType, bool) {
	switch SubscriptionType(s) {
	case Exclusive:
		return Exclusive, true
	case Shared:
		return Shared, true
	case Failover:
		return Failover, true
	case KeyShared:
		return KeyShared, true
	default:
		return "", false
	}
}

func (s SubscriptionType) PulsarType() pulsar.SubscriptionType {
	switch s {
	case Exclusive:
		return pulsar.Exclusive
	case Shared:
		return pulsar.Shared
	case Failover:
		return pulsar.Failover
	case KeyShared:
		return pulsar.KeyShared
	default:
		return pulsar.Exclusive
	}
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
	err := sdk.Util.ParseConfig(cfg, &s.config)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	if stype, ok := cfg["subscriptionType"]; ok {
		subscriptionType, ok := ParseSubscriptionType(stype)
		if !ok {
			return fmt.Errorf("invalid subscriptionType: %s", stype)
		}

		s.config.SubscriptionType = subscriptionType
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

	subscriptionType := s.config.SubscriptionType.PulsarType()

	s.consumer, err = client.Subscribe(pulsar.ConsumerOptions{
		Topic:            s.config.Topic,
		SubscriptionName: s.config.SubscriptionName,
		Type:             subscriptionType,
	})
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to create consumer: %w", err)
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

func (s *Source) Teardown(ctx context.Context) error {
	if s.consumer != nil {
		s.consumer.Close()
	}
	return nil
}
