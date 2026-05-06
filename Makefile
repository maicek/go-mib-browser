build:
	go build \
	-a \
	-trimpath \
	-ldflags="-s -w -X main.embedPath=dist" \
	-gcflags="all=-l=4" \
	-tags=netgo \
	-o mib-browser \
	./cmd/mib-browser