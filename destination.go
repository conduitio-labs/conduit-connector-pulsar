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
	"context"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Destination struct {
	sdk.UnimplementedDestination

	client   pulsar.Client
	producer pulsar.Producer
	config   DestinationConfig
}

func (d *Destination) Config() sdk.DestinationConfig {
	return &d.config
}

func NewDestination() sdk.Destination {
	return sdk.DestinationWithMiddleware(&Destination{})
}

func (d *Destination) Open(ctx context.Context) (err error) {
	var logger log.Logger
	if d.config.DisableLogging {
		logger = log.DefaultNopLogger()
	}

	d.client, err = pulsar.NewClient(pulsar.ClientOptions{
		URL:                        d.config.URL,
		ConnectionTimeout:          d.config.ConnectionTimeout,
		OperationTimeout:           d.config.OperationTimeout,
		MaxConnectionsPerBroker:    d.config.MaxConnectionsPerBroker,
		MemoryLimitBytes:           d.config.MemoryLimitBytes,
		EnableTransaction:          d.config.EnableTransaction,
		TLSKeyFilePath:             d.config.TLSKeyFilePath,
		TLSCertificateFile:         d.config.TLSCertificateFile,
		TLSTrustCertsFilePath:      d.config.TLSTrustCertsFilePath,
		TLSAllowInsecureConnection: d.config.TLSAllowInsecureConnection,
		TLSValidateHostname:        d.config.TLSValidateHostname,

		Logger: logger,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	sdk.Logger(ctx).Info().Msg("created destination client")

	d.producer, err = d.client.CreateProducer(pulsar.ProducerOptions{
		Topic: d.config.Topic,

		// SendTimeout set to -1 disables the timeout to prevent acceptance
		// tests to detect leaking goroutines.
		// TODO: it might be better for this to be configurable (issue #9)
		SendTimeout: -1,
	})
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	sdk.Logger(ctx).Info().Msg("created destination producer")

	return nil
}

func (d *Destination) Write(ctx context.Context, records []opencdc.Record) (int, error) {
	var written int
	for _, record := range records {
		key := string(record.Key.Bytes())
		_, err := d.producer.Send(ctx, &pulsar.ProducerMessage{
			Payload: record.Bytes(),
			Key:     key,
		})
		if err != nil {
			return written, fmt.Errorf("failed to send message: %w", err)
		}

		sdk.Logger(ctx).Trace().
			Str("topic", d.config.Topic).
			Str("key", key).Msg("sent message")
		written++
	}

	sdk.Logger(ctx).Trace().Int("total", written).Msg("wrote messages to destination")
	return written, nil
}

func (d *Destination) Teardown(ctx context.Context) error {
	if d.producer != nil {
		d.producer.Close()
	}

	if d.client != nil {
		d.client.Close()
	}

	sdk.Logger(ctx).Debug().Msg("destination teardown complete")

	return nil
}
