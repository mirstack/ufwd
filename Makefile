VERSION=$(shell ./ufwd -version 2>&1 | sed -e 's/.*v//')
MACHINE=$(shell uname -sp | tr '[A-Z]' '[a-z]' | sed -e 's/\s/-/')

pack:
	go build .
	mkdir -p tmp/bin pkg
	cp ufwd tmp/bin
	cd tmp && zip -r ../pkg/ufwd-$(VERSION)-$(MACHINE).zip bin
	rm -r tmp
