LDF="-X main.version=v1.0.1 \
-X main.buildDate=`$(date +%FT%T%z)` \
-X main.gitBranch=`$(git branch --show-current)` \
-X main.gitCommit=`$(git rev-parse --short HEAD)`"

example1:
	go run -race -ldflags=${LDF} ./.examples/1_basic/main.go

example2:
	go get -v -t -d ./.examples/2_collector
	go run -race -ldflags=${LDF} ./.examples/2_collectormain.go

example3:
	go run -race -ldflags=${LDF} ./.examples/3_cli_flags/main.go -v
