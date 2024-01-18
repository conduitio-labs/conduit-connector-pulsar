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
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Destination struct {
	sdk.UnimplementedDestination

	producer pulsar.Producer
	config   DestinationConfig
}

func NewDestination() sdk.Destination {
	return sdk.DestinationWithMiddleware(&Destination{}, sdk.DefaultDestinationMiddleware()...)
}

func (d *Destination) Parameters() map[string]sdk.Parameter {
	return d.config.Parameters()
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) error {
	sdk.Logger(ctx).Info().Msg("Configuring Destination...")

	if err := sdk.Util.ParseConfig(cfg, &d.config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}

func (d *Destination) Open(_ context.Context) error {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
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
			Payload: record.Bytes(),
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

func (d *Destination) Teardown(_ context.Context) error {
	if d.producer != nil {
		d.producer.Close()
	}

	return nil
}
