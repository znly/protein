.PHONY: toc test deps

toc:
	docker run --rm -it -v ${PWD}:/usr/src jorgeandrada/doctoc --github

test:
	staticcheck . ./protoscan
	PROT_CQL_ADDRS="localhost:9042" PROT_REDIS_URI="redis://localhost:6379/0" PROT_MEMCACHED_ADDRS="localhost:11211" go test -v -cpu 1,4 -run=. -bench=xxx . ./protoscan

deps:
	dep ensure
	go generate .
	go tool fix -force context -r context . || true
