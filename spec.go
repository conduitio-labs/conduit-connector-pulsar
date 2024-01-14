package apachepulsar

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
)

// version is set during the build process with ldflags (see Makefile).
// Default version matches default from runtime/debug.
var version = "(devel)"

// Specification returns the connector's specification.
func Specification() sdk.Specification {
	return sdk.Specification{
		Name:        "apachepulsar",
		Summary:     "An Apache Pulsar source and destination plugin for Conduit, written in Go.",
		Description: "An Apache Pulsar source and destination plugin for Conduit, written in Go.",
		Version:     version,
		Author:      "Meroxa, Inc.",
	}
}
