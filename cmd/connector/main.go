package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"

	pulsar "github.com/alarbada/conduit-connector-apache-pulsar"
)

func main() {
	sdk.Serve(pulsar.Connector)
}
