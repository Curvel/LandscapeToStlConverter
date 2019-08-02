#!/bin/sh

echo "1dfd4e36aa74dbdd7cbe0427747380f6 /go/src/golang.org/x/image/tiff/reader.go" | md5sum -c - 
if [ $? -eq 0 ]; then
    echo "Patching reader.go ..."
	rm /go/src/golang.org/x/image/tiff/reader.go
	cp /app/patches/reader.go /go/src/golang.org/x/image/tiff/reader.go
	echo "reader.go patched!"
else
    echo "Unable to patch reader.go"
	exit 1
fi

# all fine
exit 0
