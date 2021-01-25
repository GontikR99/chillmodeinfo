.PHONY: all clean start

$(shell mkdir -p bin build/electron >/dev/null 2>&1 || true)

all: bin/chillmodeinfo.exe ;

start: web/exe/main.js web/static/data/chillmodeinfo.wasm build/electron/main.wasm build/electron/main.wasm
	cp web/exe/main.js build/electron
	cp -r web/static/data/* build/electron
	npm start

bin/chillmodeinfo.exe: web/static/staticfiles_vfsdata.go $(shell find cmd/chillmodeinfo -type f) $(shell find internal -type f)
	go build -o $@ ./cmd/chillmodeinfo

build/electron/main.wasm: $(shell find web/exe -name \*.go) $(shell find internal -type f)
	GOOS=js GOARCH=wasm go build -o $@ ./web/exe

web/static/data/chillmodeinfo.wasm: $(shell find web/app -type f) $(shell find internal -type f)
	go run -mod=vendor github.com/vugu/vugu/cmd/vugugen -s -r -skip-go-mod -skip-main web/app
	#GOOS=js GOARCH=wasm tinygo build -o $@ ./web/app
	GOOS=js GOARCH=wasm go build -o $@ ./web/app

web/static/staticfiles_vfsdata.go: $(shell find web/static/data -type f) web/static/data/chillmodeinfo.wasm
	go generate -tags=dev ./web/static

clean:
	rm -rf \
		out \
		bin/* \
		build/* \
		web/static/staticfiles_vfsdata.go \
		web/static/data/chillmodeinfo.wasm \
		web/exe/main.wasm \
		$(shell find web/app -name 0_components_vgen.go)

