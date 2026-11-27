# poc-raspi-dht-sensor

Reads a DHT11/DHT22 sensor on a Raspberry Pi, read temperature & humidity and forward to [InfluxDB Cloud Serverless](https://us-east-1-1.aws.cloud2.influxdata.com/) every 60 seconds. A technical demonstration of the edge-to-cloud path, not a production agent.

Sensors tested:

- DHT11
- DHT22/AM2302

Full feature tested on:

- Raspberry Pi
  - Gen 5
  - Gen 4
  - Zero 2W
- Tinker Board
  - 2S

## Usage

Connect Raspberry Pi device pins with sensor pins, better when power-offed:

| Raspberry Pi | Sensor |
|------------------|------------|
| GPIO4            | DAT or out |
| GND              | GND or - |
| 3V3              | VCC or + |

If there is a `NC` pin on the DHT22 sensor board, ignore it. It is a physical placeholder required by the protocal of AM2302 that DHT22 follows, no output at all.

> Use [`pinout` cli command](https://kamaslau.wordpress.com/2026/10/31/cli-commands-for-inspecting-raspberry-pi/) to map GPTO pin names with physical (J8) pins through a graphic.

```bash
# Mount the overlay on GPIO 4 and persist it across reboots
sudo sh -c 'echo "dtoverlay=dht11,gpiopin=4" >> /boot/firmware/config.txt'
sudo reboot

# Init .env, then fill in the three TS_* values (see .env.example)
cp .env.example .env
go run ./...
```

The process exits when the time series store is not usable — missing `TS_*` values, a rejected token, or a failed first write. It does not keep collecting into nothing; the error names what is wrong.

## Run in background

```bash
# Build the executable
mkdir -p tmp && go build -o ./tmp/main ./src

cd ~/poc-raspi-dht-sensor && nohup ./tmp/main > output.log 2>&1 &
```

## Run on OS start

```bash
mkdir -p tmp && go build -o ./tmp/main ./src

sudo nano /etc/systemd/system/poc-raspi-dht-sensor.service
```

```ini
[Unit]
Description=poc-raspi-dht-sensor
After=network-online.target
Wants=network-online.target

[Service]
Restart=always
User=pi
WorkingDirectory=/home/pi/Apps/poc-raspi-dht-sensor
ExecStart=/home/pi/Apps/poc-raspi-dht-sensor/tmp/main

[Install]
WantedBy=multi-user.target
```

> The user & path should match **YOUR** situation.
> Run `systemctl daemon-reload` after modify the `.service` file.

`After=network-online.target` rather than `network.target`: the service needs a resolved hostname and a route, not just a configured interface.

Then enable and start it:

```bash
sudo systemctl enable poc-raspi-dht-sensor.service
sudo systemctl start poc-raspi-dht-sensor.service
sudo systemctl status poc-raspi-dht-sensor.service
```

---

Built by [Kamas Lau](https://www.kamaslau.com/works?utm_source=github&utm_medium=repo-readme&utm_campaign=always-on&utm_content=poc) — I do this for foreign companies operating in China, on their own infrastructure.
