FROM golang:latest

# golang
WORKDIR /app

# First all large folders
COPY srtm ./srtm

# Rest
COPY lib ./lib
COPY patches ./patches
COPY patch-tiff-package.sh ./patch-tiff-package.sh
COPY TiffDecoder.go ./TiffDecoder.go

# Download all dependencies. Dependencies will be cached if image not changed
RUN go get -d .

# Patch tiff package
RUN ./patch-tiff-package.sh

# Build the go app
RUN go build -o TiffDecoder .

# keep docker running
#CMD tail -f /dev/null