.PHONY:

test:
	make install &&\
	go generate -v ./...

debug:
	make install &&\
	goentity -help

install:
	go install ${CURDIR}/cmd/inspector &&\
	go install ${CURDIR}/cmd/goentity