# Delicious GPS

A data scrubber that uses the running GPSD server to scrape the GPS data and throw it into a local sqlite database

Designed and tested on a Raspberry Pi Zero W. Your mileage may vary on other devices/platforms

## Pre-requistes

The system running the output binary needs to be running gpsd

```sh
sudo apt install gpsd
```

The build system needs:

- Docker [Desktop]
- Docker buildx running with `linux/arm/v6` in `docker buildx ls`
- [Optional] cmake, golang, etc for building and debugging locally

On Linux (tested), you if you don't see `linux/arm/v6` with `docker build ls` then this should hopefully help:

```sh
docker run --privileged --rm tonistiigi/binfmt --install linux/arm/v6
```

Note: this is setting up QEMU and the likes on your host, so Best Careful and Know What You're Doing :tm:

## Building

To speed up building, it is faster to build the binary "locally" (i.e. in a docker container), and then send it to the raspberry pi. Building on the Raspberry Pi will probably work, but with the docker build you should be able to avoid the "It works on my machine" problem.

The process requires:

- Generating a docker image to build - Once or when deps are updated
- Running the make process against the docker image - Every time you need a build

This is as simple as:

```sh
make docker-builder   # For the docker image (run once)
make build-rpi        # For the raspberry pi image (run for each build)
```

After this, a new binary should be available in `./bin/delicious-gps`

## Installing & Running

After build, you should be able to copy the binary to the Rasperry Pi

```sh
scp ./bin/delicious-gps pi@<ip-address>:~/delicious-gps
```

And then run via a ssh shell

```sh
ssh pi@<ip-address>
sudo ./delicious-gps
```

## Development

### Local testing

For compiling and running locally (LED + Button won't work of course)

```sh
go run cmd/delicious-gps/main.go
```

### Raspberry pi binary

See Building above. You can in theory run this with the same command as local testing, but you'd need the build chain available on the raspberry pi (discouraged)

Need to test lots of things on the Pi? Set `PI_IP_ADDRESS` environmen variable e.g. `export PI_IP_ADDRESS=<ip address>` then `make build-and-upload`

## Database

The database will output to data.db and is a simple sqlite3 database.

Data is not cleared between runs, so expect to see old data in there unless you
manuall clear it.

## Tips

**Easily find your pi**

`nmap -p 22 --open 192.168.x.0/24`, but making sure x is the same as your current network

## TODO

- [x] Basic GPSD listener and sqlite writer
- [x] Configuration (command line can support args)
- [x] "Proper" logging/verbosity
- [x] Some level of "ready" state
- [ ] Upstartd script/installer
- [x] Log location values to the database
- [x] GPIO interfacing
  - [ ] Buttons
  - [x] neopixel (for feedback)
- [ ] Make it all work
- [ ] Buffer the database write. Every now and then it takes some time, and blocks the main thread
