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

	"github.com/alarbada/conduit-connector-apache-pulsar/test"
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
	err := source.Configure(ctx, map[string]string{
		"url":                        test.PulsarTLSURL,
		"topic":                      topic,
		"subscriptionName":           topic + "-subscription",
		"tlsAllowInsecureConnection": "false",
		"tlsValidateHostname":        "true",
		"tlsCertificateFile":         "./test/client.cert.pem",
		"tlsKeyFilePath":             "./test/client.key-pk8.pem",
		"tlsTrustCertsFilePath":      "./test/ca.cert.pem",
	})
	is.NoErr(err)

	err = source.Open(ctx, nil)
	is.NoErr(err)

	defer func() {
		err := source.Teardown(ctx)
		is.NoErr(err)
	}()

	destination := NewDestination()
	err = destination.Configure(ctx, map[string]string{
		"url":                        test.PulsarTLSURL,
		"topic":                      topic,
		"tlsAllowInsecureConnection": "false",
		"tlsValidateHostname":        "true",
		"tlsCertificateFile":         "./test/client.cert.pem",
		"tlsKeyFilePath":             "./test/client.key-pk8.pem",
		"tlsTrustCertsFilePath":      "./test/ca.cert.pem",
	})
	is.NoErr(err)

	err = destination.Open(ctx)
	is.NoErr(err)

	defer func() {
		err := destination.Teardown(ctx)
		is.NoErr(err)
	}()

	rec := sdk.Util.Source.NewRecordCreate(
		[]byte(uuid.NewString()),
		sdk.Metadata{"pulsar.topic": topic},
		sdk.RawData("test-key"),
		sdk.RawData(exampleMessage),
	)

	_, err = destination.Write(ctx, []sdk.Record{rec})
	is.NoErr(err)

	readRec, err := source.Read(ctx)
	is.NoErr(err)

	want := string(rec.Key.Bytes())
	actual := string(readRec.Key.Bytes())
	is.Equal(want, actual)
}
