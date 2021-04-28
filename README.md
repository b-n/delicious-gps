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

And run it: `./delcious-gps`

### Raspberry pi output

To get a binary: `make build-rpi`

Copy the binary to your raspberry pi: `scp delcious-gps pi@<ip>:delicious-gps`

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

