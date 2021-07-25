# Build ws2811 c library

FROM debian AS ws281x_builder

WORKDIR /foundry

RUN apt-get update -y && apt-get install -y \
  build-essential \
  cmake \
  git

RUN git clone https://github.com/jgarff/rpi_ws281x.git \
  && cd rpi_ws281x \
  && mkdir build \
  && cd build \
  && cmake -D BUILD_SHARED=OFF -D BUILD_TEST=OFF .. \
  && cmake --build . \
  && make install

# Setup buikd env

FROM golang:1.15

ARG CACHE_LOCATION=/.cache

COPY --from=ws281x_builder /usr/local/lib/libws2811.a /usr/local/lib/
COPY --from=ws281x_builder /usr/local/include/ws2811 /usr/local/include/ws2811

RUN apt-get update -y && apt-get install -y \
  build-essential \
  cmake

RUN groupadd -r -g 1000 user && useradd --no-log-init --no-create-home -d / -r -g user -u 1000 user
RUN mkdir -p /cache-build
RUN mkdir -p $CACHE_LOCATION
RUN chown -R user:user /cache-build $CACHE_LOCATION

USER user

WORKDIR /cache-build
COPY . /cache-build
RUN env GCO_ENABLED=1 GOCACHE=$CACHE_LOCATION go build -v -o cache-build cmd/delicious-gps/main.go
