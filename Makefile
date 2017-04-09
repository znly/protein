test:
	staticcheck . ./protoscan ./protostruct
	go vet . ./protostruct # ./protoscan ## unsafe use due to sym-scan
	go test -race -cpu 4 -cover . ./protoscan ./protostruct

toc:
	docker run --rm -it -v ${PWD}:/usr/src jorgeandrada/doctoc --github

