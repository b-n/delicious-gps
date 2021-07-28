#
# Your Cross compiles make easy
#
# 'make [build]'    compiles using your current go environments build
# 'make build-rpi'  compiles a binary for use on a raspberry pi zero w (armv6)
# 'make clean'      removes any of the output binaries

PLATFORM=linux/arm/v6 # RPI Zero W
FLAGS = GCO_ENABLED=1
BUILD_OPTS = -v

BUILD_DIR = bin
TARGET = $(BUILD_DIR)/delicious-gps
ENTRY_POINT = cmd/delicious-gps/*.go
DOCKER_GO_CACHE = /.cache

default: build

build: clean
	env $(FLAGS) go build $(BUILD_OPTS) -o $(TARGET) $(ENTRY_POINT)

build-rpi:
	docker run --rm -v "$$PWD":/usr/src/delicious-gps --platform $(PLATFORM) -w /usr/src/delicious-gps -u `id -u $$USER` delicious-builder env GOCACHE=$(DOCKER_GO_CACHE) make build

docker-builder:
	docker buildx build --platform $(PLATFORM) --tag delicious-builder --file Dockerfile --progress=plain --build-arg CACHE_LOCATION=$(DOCKER_GO_CACHE) .

build-and-upload: build-rpi
	scp $(TARGET) pi@${PI_IP_ADDRESS}:~/delicious-gps

clean:
	$(RM) -r $(BUILD_DIR)
	mkdir $(BUILD_DIR)
