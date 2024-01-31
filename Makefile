.PHONY:

test:
	make install &&\
	go generate -v ./...

debug:
	make install &&\
	esc-gen-model -help

install:
	go install ${CURDIR}/cmd/inspector &&\
	go install ${CURDIR}/cmd/esc-gen-model