.PHONY: all clean start package server

$(shell mkdir -p bin electron/src >/dev/null 2>&1 || true)

all: server package;

server: bin/chillmodeinfo.exe

start: electron/.electron
	cd electron && npm start

package: electron/.electron
	cd electron && npm run make

electron/.electron: bin/electronmain.wasm cmd/electronmain/electronmain.js cmd/electronmain/preload.js bin/webapp.wasm
	cp -r web/static/data/* electron/src
	cp bin/electronmain.wasm cmd/electronmain/electronmain.js cmd/electronmain/preload.js bin/webapp.wasm electron/src
	touch $@

bin/chillmodeinfo.exe: web/static/staticfiles_vfsdata.go $(shell find cmd/chillmodeinfo -type f) $(shell find internal -type f)
	go build -o $@ ./cmd/chillmodeinfo

bin/electronmain.wasm: $(shell find cmd/electronmain -name \*.go) $(shell find internal -type f)
	GOOS=js GOARCH=wasm go build -tags electron -o $@ ./cmd/electronmain

bin/webapp.wasm: $(shell find cmd/webapp -type f) $(shell find internal -type f)
	go run -mod=vendor github.com/vugu/vugu/cmd/vugugen -s -r -skip-go-mod -skip-main cmd/webapp
	GOOS=js GOARCH=wasm go build -tags web -o $@ ./cmd/webapp

web/static/data/app.wasm: bin/app.wasm
	cp $< $@

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
		$(shell find cmd -name 0_components_vgen.go)

