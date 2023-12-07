package apachepulsar_test

import (
	"context"
	"testing"

	apachepulsar "github.com/alarbada/conduit-connector-apachepulsar"
	"github.com/matryer/is"
)

func TestTeardownSource_NoOpen(t *testing.T) {
	is := is.New(t)
	con := apachepulsar.NewSource()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}
