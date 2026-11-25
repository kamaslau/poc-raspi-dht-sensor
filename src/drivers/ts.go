package drivers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

// TSClient is the process-wide InfluxDB client. It stays nil until ConnectTS
// succeeds, so every caller checks it instead of dereferencing it.
var TSClient *influxdb3.Client

// ConnectTS connects to the time series store. An error means there is no
// usable client: the caller must stop rather than continue, because the first
// write would otherwise panic on a nil pointer.
func ConnectTS() error {
	host := os.Getenv("TS_HOST")
	token := os.Getenv("TS_TOKEN")
	database := os.Getenv("TS_BUCKET")

	var missing []string
	for _, item := range []struct{ name, value string }{
		{"TS_HOST", host},
		{"TS_TOKEN", token},
		{"TS_BUCKET", database},
	} {
		if item.value == "" {
			missing = append(missing, item.name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing configuration: %s (see .env.example)", strings.Join(missing, ", "))
	}

	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     host,
		Token:    token,
		Database: database,
	})
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}
	TSClient = client

	return logTSConnected()
}

// logTSConnected writes a point that proves the database accepts writes. A
// client that can be created but not written to is not a working connection.
func logTSConnected() error {
	point := influxdb3.NewPoint(
		"connection_log",
		map[string]string{"type": "ts"},
		map[string]any{
			"status":  "connected",
			"message": fmt.Sprintf("Time series storage connection succeed at %s", time.Now().Format(time.RFC3339)),
		},
		time.Now(),
	)

	if err := TSClient.WritePoints(context.Background(), []*influxdb3.Point{point}); err != nil {
		return fmt.Errorf("write connection log: %w", err)
	}

	log.Print("👍 TS connected")
	return nil
}

// TSWrite writes a single measurement. It returns an error instead of panicking
// when the client was never connected.
func TSWrite(measurement string, tags map[string]string, fields map[string]any) error {
	if TSClient == nil {
		return fmt.Errorf("time series storage is not connected")
	}

	point := influxdb3.NewPoint(
		measurement,
		tags,
		fields,
		time.Now(),
	)

	if err := TSClient.WritePoints(context.Background(), []*influxdb3.Point{point}); err != nil {
		return err
	}
	return nil
}
