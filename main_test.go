package apachepulsar_test

import (
	"errors"
	"net/http"
	"os"
	"testing"
	"time"
)

func healthcheck() error {
	for counter := 0; counter < 100; counter++ {
		res, err := http.Get("http://localhost:8080/admin/v2/brokers/health")
		if err != nil {
			return err
		}

		if res.StatusCode == 200 {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	return errors.New("healthcheck failed")
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
