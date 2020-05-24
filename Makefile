default:
	PKG_CONFIG_PATH="/usr/local/lib/pkgconfig" \
	GO111MODULE=on GOPROXY=https://goproxy.io go build \
		-ldflags '-w -s' -buildmode="c-shared" \
		-o bin/client.dll -v ./cmd/client
