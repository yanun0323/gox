.PHONY:

CURDIR = $(shell printf "%q\n" "$(PWD)")

open:
	open ${HOME}/go/bin

ls:
	ls ${HOME}/go/bin

run:
	make install &&\
	go generate ./...

run.debug:
	make install &&\
	go generate -v ./...

help:
	make install &&\
	modelgen -h &&\
	domaingen -h

install:
	GOBIN=/usr/local/bin/ sudo go install ${CURDIR}/cmd/inspector &&\
	GOBIN=/usr/local/bin/ sudo go install ${CURDIR}/cmd/modelgen &&\
	GOBIN=/usr/local/bin/ sudo go install ${CURDIR}/cmd/domaingen

remove:
	rm -rf ${HOME}/go/bin/inspector;\
	rm -rf ${HOME}/go/bin/modelgen;\
	rm -rf ${HOME}/go/bin/domaingen