package apachepulsar_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

// As of writing, the latest version of the pulsar-client-go driver (0.11.1)
// doesn't include the pulsaradmin.
// It is currently on the master branch (v0.12.0, but not yet released, so we do the
// health request ourselves.
// Here the docs: https://pulsar.apache.org/admin-rest-api/#operation/BrokersBase_healthCheck

func areBrokersHealthy() bool {
	timeout := time.Duration(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx,
		"GET", "http://localhost:8080/admin/v2/brokers/health", nil)
	if err != nil {
		return false
	}

	res, err := http.DefaultClient.Do(req)
	if errors.Is(err, context.DeadlineExceeded) {
		return false
	} else if err != nil {
		return false
	}

	return res.StatusCode == 200
}

func healthcheck(t *testing.T) {
	for timesChecked := 0; timesChecked < 100; timesChecked++ {
		if !areBrokersHealthy() {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		return
	}

	t.Errorf("healthcheck retry limit reached")
}
