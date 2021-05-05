#
# Your Cross compiles make easy
#
# 'make [build]'    compiles using your current go environments build
# 'make build-rpi'  compiles a binary for use on a raspberry pi zero w (armv6)
# 'make clean'      removes any of the output binaries

OS = linux
ARCH = arm
ARCH_VERSION = 6
COMPILER = gcc
CC = arm-linux-gnueabi-gcc
FLAGS = GCO_ENABLED=1

BUILD_DIR = bin
TARGET = $(BUILD_DIR)/delicious-gps
ENTRY_POINT = cmd/delicious-gps/main.go

default: build

build: clean
	env $(FLAGS) go build -o $(TARGET) $(ENTRY_POINT)

build-rpi: clean
	env GOOS=$(OS) GOARCH=$(ARCH) GOARM=$(ARCH_VERSION) CC_FOR_$(OS)_$(ARCH)=$(COMPILER) CC=$(CC) $(FLAGS) go build -o $(TARGET) $(ENTRY_POINT)

clean:
	$(RM) -r $(BUILD_DIR)
	mkdir $(BUILD_DIR)
