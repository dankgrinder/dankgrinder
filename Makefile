GO=go
VERSION=0.1.1

.PHONY: build pack clean release

build:
	(cd cmd/dankgrinder; $(GO) build -v -o ../../dankgrinder .)
	(cd cmd/candy; $(GO) build -v -o ../../candy .)

pack:
	mkdir dankgrinder-$(VERSION)-windows-amd64
	mkdir dankgrinder-$(VERSION)-darwin-amd64
	mkdir dankgrinder-$(VERSION)-linux-amd64
	cp config.json dankgrinder-$(VERSION)-windows-amd64
	cp config.json dankgrinder-$(VERSION)-darwin-amd64
	cp config.json dankgrinder-$(VERSION)-linux-amd64
	(cd cmd/dankgrinder; GOOS=windows GOARCH=amd64 $(GO) build -v -o ../../dankgrinder-$(VERSION)-windows-amd64/dankgrinder.exe .)
	(cd cmd/dankgrinder; GOOS=darwin GOARCH=amd64 $(GO) build -v -o ../../dankgrinder-$(VERSION)-darwin-amd64/dankgrinder .)
	(cd cmd/dankgrinder; GOOS=linux GOARCH=amd64 $(GO) build -v -o ../../dankgrinder-$(VERSION)-linux-amd64/dankgrinder .)
	(cd cmd/candy; GOOS=windows GOARCH=amd64 $(GO) build -v -o ../../dankgrinder-$(VERSION)-windows-amd64/candy.exe .)
	(cd cmd/candy; GOOS=darwin GOARCH=amd64 $(GO) build -v -o ../../dankgrinder-$(VERSION)-darwin-amd64/candy .)
	(cd cmd/candy; GOOS=linux GOARCH=amd64 $(GO) build -v -o ../../dankgrinder-$(VERSION)-linux-amd64/candy .)
	zip -r9j dankgrinder-$(VERSION)-windows-amd64.zip dankgrinder-$(VERSION)-windows-amd64/*
	zip -r9j dankgrinder-$(VERSION)-darwin-amd64.zip dankgrinder-$(VERSION)-darwin-amd64/*
	tar -czvf dankgrinder-$(VERSION)-linux-amd64.tar.gz -C dankgrinder-$(VERSION)-linux-amd64 .

release: pack clean

clean:
	$(GO) clean
	rm -rf dankgrinder-$(VERSION)-windows-amd64 dankgrinder-$(VERSION)-darwin-amd64 dankgrinder-$(VERSION)-linux-amd64
