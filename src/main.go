package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"app/src/drivers"

	"github.com/joho/godotenv"
)

// loadEnv loads .env into the process environment. Missing credentials are
// reported by ConnectTS, which names the variables that were not set.
func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Fail loading .env: %s (copy .env.example to .env first)", err)
	}
}

// initServices brings the dependencies up before the first reading. It returns
// an error rather than logging and continuing: a reader without a writable time
// series store would loop forever and report nothing.
func initServices() error {
	return drivers.ConnectTS()
}

func readIIO(path string) (float64, bool) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, false
	}

	raw := strings.TrimSpace(string(b))
	if raw == "" {
		return 0, false
	}

	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}

	return float64(v) / 1000.0, true
}

func main() {
	loadEnv()
	if err := initServices(); err != nil {
		log.Fatalf("⛔ %s", err)
	}

	const (
		tempPath = "/sys/bus/iio/devices/iio:device0/in_temp_input"
		humPath  = "/sys/bus/iio/devices/iio:device0/in_humidityrelative_input"
	)

	for {
		t, tempOK := readIIO(tempPath)
		h, humOK := readIIO(humPath)
		if !tempOK || !humOK {
			fmt.Println("driver is not ready, retry...")
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Printf("T = %.2f °C, RH = %.1f %%\n", t, h)

		if err := drivers.TSWrite("dht11", map[string]string{"type": "dht11"}, map[string]any{"temperature": t, "humidity": h}); err != nil {
			fmt.Println("Failed to write to TS:", err)
		}
		time.Sleep(60 * time.Second)
	}
}
