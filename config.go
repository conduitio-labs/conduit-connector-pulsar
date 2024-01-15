package apachepulsar

import (
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

//go:generate paramgen -output=paramgen_src.go SourceConfig
//go:generate paramgen -output=paramgen_dest.go DestinationConfig

type Config struct {
	// URL of the Pulsar instance to connect to.
	URL string `json:"URL" validate:"required"`

	// ConnectionTimeout specifies the duration for which the client will
	// attempt to establish a connection before timing out.
	ConnectionTimeout time.Duration `json:"connectionTimeout"`

	// OperationTimeout is the duration after which an operation is considered
	// to have timed out.
	OperationTimeout time.Duration `json:"operationTimeout"`

	// MaxConnectionsPerBroker limits the number of connections to each broker.
	MaxConnectionsPerBroker int `json:"maxConnectionsPerBroker"`

	// MemoryLimitBytes sets the memory limit for the client in bytes.
	// If the limit is exceeded, the client may start to block or fail operations.
	MemoryLimitBytes int `json:"memoryLimitBytes"`

	// EnableTransaction determines if the client should support transactions.
	EnableTransaction bool `json:"enableTransaction"`
}

func parseConfig(validate *configValidator, cfg map[string]string) Config {
	config := Config{
		URL:                     validate.RequiredStr("URL"),
		ConnectionTimeout:       validate.OptionalDuration("connectionTimeout"),
		OperationTimeout:        validate.OptionalDuration("operationTimeout"),
		MaxConnectionsPerBroker: validate.OptionalInt("maxConnectionsPerBroker"),
		MemoryLimitBytes:        validate.OptionalInt("maxConnectionsPerBroker"),
		EnableTransaction:       validate.OptionalBool("enableTransaction"),
	}

	return config
}

type SourceConfig struct {
	Config

	// Topic specifies the Pulsar topic from which the source will consume messages.
	Topic string `json:"topic" validate:"required"`

	// SubscriptionName is the name of the subscription to be used for
	// consuming messages.
	SubscriptionName string `json:"subscriptionName" validate:"required"`

	// SubscriptionType defines the type of subscription to use. This can be
	// either "exclusive", "shared", "failover" or "keyshared".
	SubscriptionType SubscriptionType `json:"subscriptionType"`
}

func parseSourceConfig(cfg map[string]string) (SourceConfig, error) {
	validate := newConfigValidator(cfg)

	sharedConfig := parseConfig(validate, cfg)

	config := SourceConfig{
		Config:           sharedConfig,
		Topic:            validate.RequiredStr("topic"),
		SubscriptionName: validate.RequiredStr("subscriptionName"),
		SubscriptionType: parseSubscriptionType(validate, "subscriptionType"),
	}

	return config, validate.Error()
}

// SubscriptionType is the type of subscription to create. Mapped directly from
// pulsar.SubscriptionType
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

func (s SubscriptionType) ToPulsar() pulsar.SubscriptionType {
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

type DestinationConfig struct {
	Config

	// Topic specifies the Pulsar topic to which the destination will produce messages.
	Topic string `json:"topic" validate:"required"`
}

func parseDestinationConfig(cfg map[string]string) (DestinationConfig, error) {
	validate := newConfigValidator(cfg)

	sharedConfig := parseConfig(validate, cfg)

	config := DestinationConfig{
		Config: sharedConfig,
		Topic:  validate.RequiredStr("topic"),
	}

	return config, validate.Error()
}
