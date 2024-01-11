package apachepulsar_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

// As of writing, the latest version of the pulsar-client-go driver (0.11.1)
// doesn't include the pulsaradmin.
// It is currently on the master branch (v0.12.0, but not yet released, so we do the
// health request ourselves.
// Here the docs: https://pulsar.apache.org/admin-rest-api/#operation/BrokersBase_healthCheck

func checkBrokersHealth() (ok bool, err error) {
	timeout := time.Duration(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx,
		"GET", "http://localhost:8080/admin/v2/brokers/health", nil)
	if err != nil {
		return false, err
	}

	res, err := http.DefaultClient.Do(req)
	if errors.Is(err, context.DeadlineExceeded) {
		// timeout, caller needs to retry
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("error while checking brokers health: %w", err)
	}

	if res.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}

func healthcheck() error {
	for timesChecked := 0; timesChecked < 100; timesChecked++ {
		ok, err := checkBrokersHealth()
		if err != nil {
			fmt.Println(err)
			continue
		} else if !ok {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		return nil
	}

	return errors.New("healthcheck failed, retry limit reached")
}

func TestMain(m *testing.M) {
	for {
		if err := healthcheck(); err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}

	code := m.Run()
	os.Exit(code)
}
