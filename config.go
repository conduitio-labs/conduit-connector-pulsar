package pulsar

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

type SourceConfig struct {
	Config

	// Topic specifies the Pulsar topic from which the source will consume messages.
	Topic string `json:"topic" validate:"required"`

	// SubscriptionName is the name of the subscription to be used for
	// consuming messages.
	SubscriptionName string `json:"subscriptionName" validate:"required"`

	// SubscriptionType defines the type of subscription to use. This can be
	// either "exclusive", "shared", "failover" or "keyshared".
	//
	// With "exclusive" there can be only 1 consumer on the same topic with the same subscription name
	//
	// With "shared" subscription mode, multiple consumer will be able to use the same subscription name
	// and the messages will be dispatched according to
	// a round-robin rotation between the connected consumers
	//
	// With "failover" subscription mode, multiple consumer will be able to use the same subscription name
	// but only 1 consumer will receive the messages.
	// If that consumer disconnects, one of the other connected consumers will start receiving messages.
	//
	// With "key_shared" subscription mode, multiple consumer will be able to use the same
	// subscription and all messages with the same key will be dispatched to only one consumer
	SubscriptionType string `json:"subscriptionType" validate:"inclusion=exclusive|shared|failover|key_shared"`
}

var subscriptionTypes = map[string]pulsar.SubscriptionType{
	"exclusive":  pulsar.Exclusive,
	"shared":     pulsar.Shared,
	"failover":   pulsar.Failover,
	"key_shared": pulsar.KeyShared,
}

func parseSubscriptionType(s string) (pulsar.SubscriptionType, bool) {
	subscriptionType, ok := subscriptionTypes[s]
	return subscriptionType, ok
}

type DestinationConfig struct {
	Config

	// Topic specifies the Pulsar topic to which the destination will produce messages.
	Topic string `json:"topic" validate:"required"`
}
