.PHONY: all clean start package server

$(shell mkdir -p bin electron/src >/dev/null 2>&1 || true)

all: server package;

server: bin/chillmodeinfo.exe bin/chillmodeinfo.linux bin/cmidb.linux

start: electron/.electron
	cd electron && npm start

package: electron/.electron
	cd electron && npm run make
	find electron/out/make -name \*.exe -exec cp \{\} bin \;

WASMS=bin/webapp.wasm bin/overlay_position.wasm bin/overlay_bids.wasm

electron/.electron: $(WASMS) bin/electronmain.wasm cmd/electronmain/electronmain.js cmd/electronmain/preload.js $(shell find web/static/data -type f)
	cp -r web/static/data/* electron/src
	cp bin/electronmain.wasm cmd/electronmain/electronmain.js cmd/electronmain/preload.js electron/src
	mkdir -p electron/src/bin
	cp $(WASMS) electron/src/bin
	touch $@

bin/chillmodeinfo.linux: web/static/staticfiles_vfsdata.go web/bin/binfiles_vfsdata.go $(shell find cmd/chillmodeinfo -type f) $(shell find internal -type f) $(shell find pkg -type f)
	GOOS=linux GOARCH=amd64 go build -tags server -o $@ ./cmd/chillmodeinfo

bin/cmidb.linux: $(shell find cmd/cmidb -type f) $(shell find internal -type f) $(shell find pkg -type f)
	GOOS=linux GOARCH=amd64 go build -tags server -o $@ ./cmd/cmidb

bin/chillmodeinfo.exe: web/static/staticfiles_vfsdata.go web/bin/binfiles_vfsdata.go $(shell find cmd/chillmodeinfo -type f) $(shell find internal -type f) $(shell find pkg -type f)
	go build -tags server -o $@ ./cmd/chillmodeinfo

bin/electronmain.wasm: $(shell find cmd/electronmain -name \*.go) $(shell find internal -type f) $(shell find pkg -type f)
	GOOS=js GOARCH=wasm go build -tags electron -o $@ ./cmd/electronmain

bin/webapp.wasm: $(shell find cmd/webapp -type f) $(shell find internal -type f) $(shell find pkg -type f)
	go run -mod=vendor github.com/vugu/vugu/cmd/vugugen -s -r -skip-go-mod -skip-main cmd/webapp
	GOOS=js GOARCH=wasm go build -tags web -o $@ ./cmd/webapp

bin/overlay_position.wasm: $(shell find cmd/overlay_position -type f) $(shell find internal -type f) $(shell find pkg -type f)
	go run -mod=vendor github.com/vugu/vugu/cmd/vugugen -s -r -skip-go-mod -skip-main cmd/overlay_position
	GOOS=js GOARCH=wasm go build -tags web -o $@ ./cmd/overlay_position

bin/overlay_bids.wasm: $(shell find cmd/overlay_bids -type f) $(shell find internal -type f) $(shell find pkg -type f)
	go run -mod=vendor github.com/vugu/vugu/cmd/vugugen -s -r -skip-go-mod -skip-main cmd/overlay_bids
	GOOS=js GOARCH=wasm go build -tags web -o $@ ./cmd/overlay_bids

web/bin/data/webapp.wasm: bin/webapp.wasm
	cp $< $@

web/static/staticfiles_vfsdata.go: $(shell find web/static/data -type f)
	go generate -tags=dev ./web/static

web/bin/binfiles_vfsdata.go: $(shell find web/bin/data -type f) web/bin/data/webapp.wasm
	go generate -tags=dev ./web/bin

clean:
	rm -rf \
		out \
		bin/* \
		build/* \
		electron/.electron \
		electron/src/* \
		electron/out \
		web/static/staticfiles_vfsdata.go \
		web/bin/binfiles_vfsdata.go \
		web/bin/data/webapp.wasm \
		$(shell find cmd -name 0_components_vgen.go)

