TOOLS    = $(CURDIR)/../tools

all:
	$(TOOLS)/go-bindata -o testdata.go --pkg testdata pkg/ && \
		gofmt -w -s testdata.go
