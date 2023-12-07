package apachepulsar_test

import (
	"context"
	"testing"

	apachepulsar "github.com/alarbada/conduit-connector-apachepulsar"
	"github.com/matryer/is"
)

func TestTeardown_NoOpen(t *testing.T) {
	is := is.New(t)
	con := apachepulsar.NewDestination()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}
