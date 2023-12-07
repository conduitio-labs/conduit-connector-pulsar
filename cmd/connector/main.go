package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"

	apachepulsar "github.com/alarbada/conduit-connector-apachepulsar"
)

func main() {
	sdk.Serve(apachepulsar.Connector)
}
