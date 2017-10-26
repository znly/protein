.PHONY: toc test deps deps_external

toc:
	docker run --rm -it -v ${PWD}:/usr/src jorgeandrada/doctoc --github

test:
	staticcheck . ./protoscan
	PROT_CQL_ADDRS="localhost:9042" PROT_REDIS_URI="redis://localhost:6379/0" PROT_MEMCACHED_ADDRS="localhost:11211" go test -v -cpu 1,4 -run=. -bench=xxx . ./protoscan

deps_external:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/znly/linkname-gen
	go get -u golang.org/x/tools/cmd/stringer

deps: deps_external
	dep ensure -v
	go generate .
	go tool fix -force context -r context . || true
