.PHONY:

test:
	make install &&\
	go generate -v ./...

install:
	go install ${CURDIR}/cmd/inspector &&\
	go install ${CURDIR}/cmd/goentity