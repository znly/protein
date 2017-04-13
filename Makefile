test:
	staticcheck . ./protoscan
	go vet . # ./protoscan ## unsafe use due to sym-scan
	go test -race -cpu 1,2,4 -cover . ./protoscan

toc:
	docker run --rm -it -v ${PWD}:/usr/src jorgeandrada/doctoc --github

