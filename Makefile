.PHONY: toc test deps

toc:
	docker run --rm -it -v ${PWD}:/usr/src jorgeandrada/doctoc --github

test:
	staticcheck . ./protoscan
	go test -v -cpu 1,4 -run=. -bench=xxx . ./protoscan

deps:
	dep ensure
	go generate .
	go tool fix -force context -r context . || true
