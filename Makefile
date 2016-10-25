all: lint build

build:
	gox -arch=amd64 -os="linux darwin"

lint: format
	cd -P . && go vet $(PACKAGES)
	for package in $(PACKAGES); do \
		golint -min_confidence .25 $$package ; \
	done

format:
	gofmt -w -s $(NOVENDOR)
