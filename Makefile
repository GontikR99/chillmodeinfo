.PHONY: all clean

all: bin/chillmodeinfo ;

bin/chillmodeinfo: web/static/staticfiles_vfsdata.go $(shell find cmd/chillmodeinfo -type f) $(shell find internal -type f)
	go build -o $@ ./cmd/chillmodeinfo

web/static/data/chillmodeinfo.wasm: $(shell find web/app -type f)
	go run -mod=vendor github.com/vugu/vugu/cmd/vugugen -s -skip-go-mod -skip-main web/app
	#GOOS=js GOARCH=wasm tinygo build -o $@ ./web/app
	GOOS=js GOARCH=wasm go build -o $@ ./web/app

web/static/staticfiles_vfsdata.go: $(shell find web/static/data -type f) web/static/data/chillmodeinfo.wasm
	go generate -tags=dev ./web/static

clean:
	rm -f bin/* \
		web/static/staticfiles_vfsdata.go \
		web/static/data/chillmodeinfo.wasm \
		web/app/0_components_vgen.go 

