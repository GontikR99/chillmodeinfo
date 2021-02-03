.PHONY: all clean start package server

$(shell mkdir -p bin electron/src >/dev/null 2>&1 || true)

all: server package;

server: bin/chillmodeinfo.exe

start: electron/.electron
	cd electron && npm start

package: electron/.electron
	cd electron && npm run make

electron/.electron: web/exe/exe.js web/exe/preload.js web/static/data/app.wasm bin/exe.wasm
	cp bin/exe.wasm electron/src
	cp web/exe/exe.js electron/src
	cp web/exe/preload.js electron/src
	cp -r web/static/data/* electron/src
	touch $@

bin/chillmodeinfo.exe: web/static/staticfiles_vfsdata.go $(shell find cmd/chillmodeinfo -type f) $(shell find internal -type f)
	go build -o $@ ./cmd/chillmodeinfo

bin/exe.wasm: $(shell find web/exe -name \*.go) $(shell find internal -type f)
	GOOS=js GOARCH=wasm go build -tags electron -o $@ ./web/exe

web/static/data/app.wasm: $(shell find web/app -type f) $(shell find internal -type f)
	go run -mod=vendor github.com/vugu/vugu/cmd/vugugen -s -r -skip-go-mod -skip-main web/app
	GOOS=js GOARCH=wasm go build -tags web -o $@ ./web/app

web/static/staticfiles_vfsdata.go: $(shell find web/static/data -type f) web/static/data/app.wasm
	go generate -tags=dev ./web/static

clean:
	rm -rf \
		out \
		bin/* \
		build/* \
		electron/.electron \
		electron/src/* \
		electron/out \
		web/static/staticfiles_vfsdata.go \
		web/static/data/app.wasm \
		$(shell find web/app -name 0_components_vgen.go)

