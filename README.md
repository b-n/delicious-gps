# Delicious GPS

A data scrubber that uses the running GPSD server to scrape the GPS data and throw it into a local sqlite database

## Pre-requistes

### Running

The system running the output binary should be running gpsd `sudo apt install gpsd`

### Building

You'll need docker running at least. The build process is housed in there. You'll always want build kit running with armv7 builds enabled.

TODO: Add some more instructions here on the above

## Development

### Local testing

As simple as `make` or `make build`

And run it: `./bin/delcious-gps`

NOTE: you may need additional dependencies, and GPIO things won't work

### Raspberry pi binar

To get a binary: `make build-rpi`

Copy the binary to your raspberry pi: `scp bin/delcious-gps pi@<ip>:~/delicious-gps`

SSH into it: `ssh pi@<ip address>`

And run it: `sudo ./delicious-gps`

## Database

The database will output to data.db and is a simple sqlite3 database.

Data is not cleared between runs, so expect to see old data in there unless you
manuall clear it.

## TODO

- [x] Basic GPSD listener and sqlite writer
- [x] Configuration (command line can support args)
- [x] "Proper" logging/verbosity
- [x] Some level of "ready" state
- [ ] Upstartd script/installer
- [x] Log location values to the database
- [x] GPIO interfacing
  - [x] Buttons
  - [x] neopixel (for feedback)
- [ ] Make it all work
