export node_bins=$(PWD)/node_modules/.bin
export uglifyjs=$(node_bins)/uglifyjs
export gulp=$(node_bins)/gulp
export GO111MODULE=on

.PHONY: server client imager test

all: server client

client: client_vendor
	$(gulp)

client_deps:
	npm install --progress false --depth 0

watch:
	$(gulp) -w

client_vendor: client_deps
	mkdir -p www/js/vendor
	$(uglifyjs) node_modules/almond/almond.js -o www/js/vendor/almond.js

css:
	$(gulp) css

generate:
	go generate ./...

statik:
	statik -src ./static/src -f -dest ./static

server:
	go build -v

server_full: generate statik
	go build -v

server_force:
	go build -v -a

server_force_full: generate statik
	go build -v -a

client_clean:
	rm -rf www/js www/css/*.css www/css/maps

client_clean_full:
	rm -rf www/js www/css/*.css www/css/maps node_modules

clean: client_clean
	rm -rf .build .ffmpeg .package target meguca-*.zip meguca-*.tar.xz meguca meguca.exe server/pkg

test:
	go test --race ./...

test_no_race:
	go test ./...

test_docker:
	docker-compose build
	docker-compose run --rm -e CI=true meguca make test

