test:
	staticcheck ./bank ./protoscan ./protostruct ./wirer
	go vet ./bank ./protostruct ./wirer # ./protoscan ## unsafe use due to sym-scan
	go test -race -cpu 4 -cover ./bank ./protoscan ./protostruct ./wirer

toc:
	docker run --rm -it -v ${PWD}:/usr/src jorgeandrada/doctoc --github

