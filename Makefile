test:
	go test -cover ./...

deps:
	@echo "## The following dependencies have been tested & are certified to"
	@echo "## work correcly with Protein's packages:\n"
	@for d in $(shell go  list -f '{{join .Deps "\n"}}' ./bank ./protoscan ./protostruct ./wirer | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | grep -v github.com/znly/protein); do \
		cd ${GOPATH}/src/$$d && git log --format='%H' | head -1 | tr -d '\n' ; \
		/bin/echo " $$d" ; \
	done
