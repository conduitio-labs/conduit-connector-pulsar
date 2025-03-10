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
	"os"
	"testing"

	"github.com/conduitio-labs/conduit-connector-pulsar/test"
	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func Test_mTLS_Setup(t *testing.T) {
	if os.Getenv("PULSAR_TLS") != "true" {
		t.Skip("Skipping mTLS tests")
	}

	is := is.New(t)

	topic := test.SetupTopicName(t, is)
	ctx := context.Background()
	source := NewSource()
	err := sdk.Util.ParseConfig(ctx, map[string]string{
		"url":                        test.PulsarTLSURL,
		"topic":                      topic,
		"subscriptionName":           topic + "-subscription",
		"tlsAllowInsecureConnection": "false",
		"tlsValidateHostname":        "true",
		"tlsCertificateFile":         "./test/certs/client.cert.pem",
		"tlsKeyFilePath":             "./test/certs/client.key-pk8.pem",
		"tlsTrustCertsFilePath":      "./test/certs/ca.cert.pem",
	}, source.Config(), Connector.NewSpecification().SourceParams)
	is.NoErr(err)

	err = source.Open(ctx, nil)
	is.NoErr(err)

	defer func() {
		err := source.Teardown(ctx)
		is.NoErr(err)
	}()

	destination := NewDestination()
	err = sdk.Util.ParseConfig(ctx, map[string]string{
		"url":                        test.PulsarTLSURL,
		"topic":                      topic,
		"tlsAllowInsecureConnection": "false",
		"tlsValidateHostname":        "true",
		"tlsCertificateFile":         "./test/certs/client.cert.pem",
		"tlsKeyFilePath":             "./test/certs/client.key-pk8.pem",
		"tlsTrustCertsFilePath":      "./test/certs/ca.cert.pem",
	}, destination.Config(), Connector.NewSpecification().DestinationParams)
	is.NoErr(err)

	err = destination.Open(ctx)
	is.NoErr(err)

	defer func() {
		err := destination.Teardown(ctx)
		is.NoErr(err)
	}()

	rec := sdk.Util.Source.NewRecordCreate(
		[]byte(uuid.NewString()),
		opencdc.Metadata{"pulsar.topic": topic},
		opencdc.RawData("test-key"),
		opencdc.RawData(exampleMessage),
	)

	_, err = destination.Write(ctx, []opencdc.Record{rec})
	is.NoErr(err)

	readRec, err := source.Read(ctx)
	is.NoErr(err)

	want := string(rec.Key.Bytes())
	actual := string(readRec.Key.Bytes())
	is.Equal(want, actual)
}
