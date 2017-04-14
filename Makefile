.PHONY: toc test bench test-bench

toc:
	docker run --rm -it -v ${PWD}:/usr/src jorgeandrada/doctoc --github

test:
	staticcheck . ./protoscan
	go vet . # ./protoscan ## unsafe use due to sym-scan
	go test -v -race -cpu 1,2,4,8,24 -cover -run=. -bench=xxx . ./protoscan

bench:
	go test -v -cpu 1,2,4,8,24 -cover -run=xxx -bench=. . ./protoscan

test-bench:
	go test -v -cpu 1,2,4,8,24 -cover -run=. -bench=. . ./protoscan
