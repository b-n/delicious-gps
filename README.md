# Delicious GPS

A data scrubber that uses the running GPSD server to scrape the GPS data and throw it into a local sqlite database

## Pre-requistes

### Running

The system running the output binary should be running gpsd `sudo apt install gpsd`

### Building

You may need the following packages installed:

    gcc-arm-linux-gnueabi
    make

## Development

### Local testing

As simple as `make` or `make build`

And run it: `./bin/delcious-gps`

### Raspberry pi binar

To get a binary: `make build-rpi`

Copy the binary to your raspberry pi: `scp bin/delcious-gps pi@<ip>:delicious-gps`

SSH into it: `ssh pi@<ip address>`

And run it: `./delicious-gps`

## Database

The database will output to data.db and is a simple sqlite3 database.

Data is not cleared between runs, so expect to see old data in there unless you
manuall clear it.

### Schema

```
position_data:
  lat: real
  lon: real
  alt: real
  velocity: real
  satelliteCount: integer
  time: datetime
```

## TODO

- [x] Basic GPSD listener and sqlite writer
- [x] Configuration (command line can support args)
- [x] "Proper" logging/verbosity
- [x] Some level of "ready" state
- [ ] Upstartd script/installer
- [ ] Log location values to the database
- [ ] GPIO interfacing
  - [ ] Buttons
  - [ ] neopixel (for feedback)

## Design thought

- Modules
  - Location
    - Has a channel `positions` which the GPSD will fire GPSD update events into
      - Event details, lon/lat/alt + errors + device name + metadata from sats (e.g. all data that makes up that position)
  - Status (GPIO)
    - has a channel `statuses` which fires events when the buttons are pressed
    - Event details: uint_16 notifying the combination of values
    - Keep strack of the current buttons that are pressed, and highlights the correct keys etc
  - Persistence
    - Methods to easily store location data in sqlite db?
- Main
  - Plumb channels together (i.e. keep current state emited by status, listen to location module )
  - Loads config
- Config
  -
