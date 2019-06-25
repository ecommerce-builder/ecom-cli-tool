ODIR=./bin
VERSION=`cat VERSION`
INSTALLDIR=`go env GOPATH`

all: clean compile

build:
	go build -o $(ODIR)/ecom -ldflags "-X main.version=$(VERSION)" ./main.go

compile:
	GOOS=linux   GOARCH=amd64 go build -o $(ODIR)/ecom-$(VERSION)-linux -ldflags "-X main.version=$(VERSION)"
	GOOS=darwin  GOARCH=amd64 go build -o $(ODIR)/ecom-$(VERSION)-darwin -ldflags "-X main.version=$(VERSION)"
	GOOS=windows GOARCH=amd64 go build -o $(ODIR)/ecom-$(VERSION).exe -ldflags "-X main.version=$(VERSION)"

install:
	go build -o $(INSTALLDIR)/bin/ecom -ldflags "-X main.version=$(VERSION)"

run:
	@go run -ldflags "-X main.version=$(VERSION)" ./main.go $(filter-out $@,$(MAKECMDGOALS))

clean:
	-@rm -r $(ODIR)/* 2> /dev/null || true
