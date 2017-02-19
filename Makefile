test:
	staticcheck ./bank ./protoscan ./protostruct ./wirer
	go vet ./bank ./protostruct ./wirer # ./protoscan ## unsafe use due to sym-scan
	go test -race -cpu 4 -cover ./bank ./protoscan ./protostruct ./wirer

deps:
	@echo "## The following dependencies have been tested & are certified to"
	@echo "## work correcly with TuyauDB's packages:\n"
	@for d in $(shell go  list -f '{{join .Deps "\n"}}' ./kv ./pipe ./client ./service | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | grep -v github.com/znly/tuyauDB); do \
		cd ${GOPATH}/src/$$d && git log --format='%H' | head -1 | tr -d '\n' ; \
		/bin/echo " $$d" ; \
	done

toc:
	docker run --rm -it -v ${PWD}:/usr/src jorgeandrada/doctoc --github

