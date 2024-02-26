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

//go:generate paramgen -output=paramgen_src.go SourceConfig

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
)

type Source struct {
	sdk.UnimplementedSource

	client   pulsar.Client
	consumer pulsar.Consumer
	config   SourceConfig
}

func NewSource() sdk.Source {
	return sdk.SourceWithMiddleware(&Source{}, sdk.DefaultSourceMiddleware()...)
}

func (s *Source) Parameters() map[string]sdk.Parameter {
	return s.config.Parameters()
}

func (s *Source) Configure(ctx context.Context, cfg map[string]string) error {
	if err := sdk.Util.ParseConfig(cfg, &s.config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	sdk.Logger(ctx).Info().Str("topic", s.config.Topic).Msg("configured source")

	return nil
}

func (s *Source) Open(ctx context.Context, pos sdk.Position) (err error) {
	var logger log.Logger
	if s.config.DisableLogging {
		logger = log.DefaultNopLogger()
	}

	s.client, err = pulsar.NewClient(pulsar.ClientOptions{
		URL:                        s.config.URL,
		ConnectionTimeout:          s.config.ConnectionTimeout,
		OperationTimeout:           s.config.OperationTimeout,
		MaxConnectionsPerBroker:    s.config.MaxConnectionsPerBroker,
		MemoryLimitBytes:           s.config.MemoryLimitBytes,
		EnableTransaction:          s.config.EnableTransaction,
		TLSKeyFilePath:             s.config.TLSKeyFilePath,
		TLSCertificateFile:         s.config.TLSCertificateFile,
		TLSTrustCertsFilePath:      s.config.TLSTrustCertsFilePath,
		TLSAllowInsecureConnection: s.config.TLSAllowInsecureConnection,
		TLSValidateHostname:        s.config.TLSValidateHostname,

		Logger: logger,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	sdk.Logger(ctx).Debug().Msg("Created Pulsar client")

	s.consumer, err = s.client.Subscribe(pulsar.ConsumerOptions{
		Topic:                       s.config.Topic,
		SubscriptionName:            s.config.SubscriptionName,
		Type:                        pulsar.Exclusive,
		SubscriptionInitialPosition: pulsar.SubscriptionPositionEarliest,
	})
	if err != nil {
		s.client.Close()
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	sdk.Logger(ctx).Debug().Msg("created pulsar consumer")

	if pos != nil {
		p, err := parsePosition(pos)
		if err != nil {
			return err
		}

		if s.config.SubscriptionName != "" && s.config.SubscriptionName != p.SubscriptionName {
			return fmt.Errorf("the old position contains a different subscription name than the connector configuration (%q vs %q), please check if the configured subscription name changed since the last run", p.SubscriptionName, s.config.SubscriptionName)
		}

		s.config.SubscriptionName = p.SubscriptionName

		sdk.Logger(ctx).Info().Str("subscriptionName", s.config.SubscriptionName).Msg("resuming from position")
	}

	if s.config.SubscriptionName == "" {
		// this must be the first run of the connector, create a new group ID
		s.config.SubscriptionName = uuid.NewString()
		sdk.Logger(ctx).Info().Str("subscriptionName", s.config.SubscriptionName).Msg("assigning source to new subscription")
	}

	return nil
}

func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
	msg, err := s.consumer.Receive(ctx)
	if err != nil {
		return sdk.Record{}, fmt.Errorf("failed to receive message: %w", err)
	}

	position := Position{
		MessageID:        msg.ID().Serialize(),
		SubscriptionName: s.config.SubscriptionName,
	}
	sdkPos := position.ToSDKPosition()

	metadata := sdk.Metadata{"pulsar.topic": msg.Topic()}
	metadata.SetCreatedAt(msg.EventTime())

	key := sdk.RawData(msg.Key())
	payload := sdk.RawData(msg.Payload())

	newRecord := sdk.Util.Source.NewRecordCreate(sdkPos, metadata, key, payload)

	sdk.Logger(ctx).Debug().Msg("received message")

	return newRecord, nil
}

func (s *Source) Ack(ctx context.Context, position sdk.Position) error {
	parsed, err := parsePosition(position)
	if err != nil {
		return err
	}

	msgID, err := pulsar.DeserializeMessageID(parsed.MessageID)
	if err != nil {
		return fmt.Errorf("failed to deserialize message ID: %w", err)
	}

	sdk.Logger(ctx).Debug().Str("MessageID", msgID.String()).Msg("acked message")

	return s.consumer.AckID(msgID)
}

func (s *Source) Teardown(ctx context.Context) error {
	if s.consumer != nil {
		s.consumer.Close()
	}

	if s.client != nil {
		s.client.Close()
	}

	sdk.Logger(ctx).Debug().Msg("source teardown complete")

	return nil
}

type Position struct {
	MessageID        []byte `json:"messageID"`
	SubscriptionName string `json:"subscriptionName"`
}

func parsePosition(pos sdk.Position) (Position, error) {
	var p Position
	err := json.Unmarshal(pos, &p)

	return p, err
}

func (p Position) ToSDKPosition() sdk.Position {
	bs, err := json.Marshal(p)
	if err != nil {
		// this error should not be possible
		panic(fmt.Errorf("error marshaling position to JSON: %w", err))
	}

	return sdk.Position(bs)
}
