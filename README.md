# poc-raspi-dht-sensor

Reads a DHT11 sensor on a Raspberry Pi and writes temperature and humidity to InfluxDB Cloud every 60 seconds. A technical demonstration of the edge-to-cloud path, not a production agent.

## Usage

Connect GPIO4 to the sensor's data output, plus 3.3 V and GND.

```bash
# Mount the overlay on GPIO 4 and persist it across reboots
sudo sh -c 'echo "dtoverlay=dht11,gpiopin=4" >> /boot/firmware/config.txt'
sudo reboot

# Init .env, then fill in the three TS_* values (see .env.example)
cp .env.example .env

# Read
go mod tidy
go run ./...

# Build the executable the unit below runs
mkdir -p tmp && go build -o ./tmp/main ./src
```

The process exits when the time series store is not usable — missing `TS_*` values, a rejected token, or a failed first write. It does not keep collecting into nothing; the error names what is wrong.

## Run in background

```bash
cd ~/poc-raspi-dht-sensor && nohup ./tmp/main > output.log 2>&1 &
```

## Run on OS start

```bash
sudo nano /etc/systemd/system/poc-raspi-dht-sensor.service
```

```ini
[Unit]
Description=IoT client (DHT11 to InfluxDB)
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=/home/pi/poc-raspi-dht-sensor/tmp/main
Restart=always
User=pi
WorkingDirectory=/home/pi/poc-raspi-dht-sensor

[Install]
WantedBy=multi-user.target
```

`After=network-online.target` rather than `network.target`: the service needs a resolved hostname and a route, not just a configured interface.

Then enable and start it:

```bash
sudo systemctl enable poc-raspi-dht-sensor.service
sudo systemctl start poc-raspi-dht-sensor.service
sudo systemctl status poc-raspi-dht-sensor.service
```

---

Built by [Kamas Lau](https://www.kamaslau.com/works?utm_source=github&utm_medium=repo-readme&utm_campaign=always-on&utm_content=poc) — I do this for foreign companies operating in China, on their own infrastructure.
