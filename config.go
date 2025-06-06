// Copyright © 2024 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pulsar

import (
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Config struct {
	// The pulsar instance to connect to.
	URL string `json:"url" validate:"required"`

	// The Pulsar topic used by the connector.
	Topic string `json:"topic" validate:"required"`

	// The duration for which the client will
	// attempt to establish a connection before timing out.
	ConnectionTimeout time.Duration `json:"connectionTimeout"`

	// The duration after which an operation is considered
	// to have timed out.
	OperationTimeout time.Duration `json:"operationTimeout"`

	// Limits the number of connections to each broker.
	MaxConnectionsPerBroker int `json:"maxConnectionsPerBroker"`

	// Sets the memory limit for the client in bytes.
	// If the limit is exceeded, the client may start to block or fail operations.
	MemoryLimitBytes int64 `json:"memoryLimitBytes"`

	// Determines if the client should support transactions.
	EnableTransaction bool `json:"enableTransaction"`

	// The path to the TLS key file
	TLSKeyFilePath string `json:"tlsKeyFilePath"`

	// The path to the TLS certificate file
	TLSCertificateFile string `json:"tlsCertificateFile"`

	// The path to the trusted TLS certificate file
	TLSTrustCertsFilePath string `json:"tlsTrustCertsFilePath"`

	// Whether the internal Pulsar client accepts untrusted TLS certificate from broker (default: false)
	TLSAllowInsecureConnection bool `json:"tlsAllowInsecureConnection"`

	// Whether the pulsar client verifies the validity of the host name from broker (default: false)
	TLSValidateHostname bool `json:"tlsValidateHostname"`

	// Disables pulsar client logs
	DisableLogging bool `json:"disableLogging"`
}

type SourceConfig struct {
	sdk.DefaultSourceMiddleware
	Config

	// The name of the subscription to be used for consuming messages.
	SubscriptionName string `json:"subscriptionName"`
}

type DestinationConfig struct {
	sdk.DefaultDestinationMiddleware
	Config
}
