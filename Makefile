LDF="-X main.version=v1.0.1 \
-X main.revision=$(git rev-parse --short HEAD) \
-X main.date=$(date +%FT%T%z)"

example1:
	go run -race -ldflags=${LDF} ./.examples/1_basic/main.go

example2:
	go get -v -t -d ./.examples/2_collector/...
	go run -race -ldflags=${LDF} ./.examples/2_collector/main.go

example3:
	go run -race -ldflags=${LDF} ./.examples/3_cli_flags/main.go -v
