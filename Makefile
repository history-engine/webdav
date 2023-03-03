NAME=webdav
VERSION=$(shell git describe --tags 2> /dev/null || echo "dev-master")
RELEASE_DIR=release
GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'

PLATFORM_LIST = darwin-amd64 darwin-arm64 linux-amd64
WINDOWS_ARCH_LIST = windows-amd64

all: linux-amd64 darwin-amd64 darwin-arm64 windows-amd64

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(RELEASE_DIR)/$(NAME)-$@

darwin-arm64:
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(RELEASE_DIR)/$(NAME)-$@

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(RELEASE_DIR)/$(NAME)-$@

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(RELEASE_DIR)/$(NAME)-$@.exe

gz_releases=$(addsuffix .gz, $(PLATFORM_LIST))
zip_releases=$(addsuffix .zip, $(WINDOWS_ARCH_LIST))

$(gz_releases): %.gz : %
	@echo $(VERSION)
	chmod +x $(RELEASE_DIR)/$(NAME)-$(basename $@)
	zip -m -j $(RELEASE_DIR)/$(NAME)-$(basename $@)-$(VERSION).zip $(RELEASE_DIR)/$(NAME)-$(basename $@)

$(zip_releases): %.zip : %
	@echo $(VERSION)
	zip -m -j $(RELEASE_DIR)/$(NAME)-$(basename $@)-$(VERSION).zip $(RELEASE_DIR)/$(NAME)-$(basename $@).exe

release: $(gz_releases) $(zip_releases)

clean:
	rm $(RELEASE_DIR)/*
