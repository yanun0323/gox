.PHONY:

CURDIR = $(shell printf "%q\n" "$(PWD)")

open:
	open /usr/local/bin/

ls:
	ls /usr/local/bin/

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
	GOBIN=/usr/local/bin/ sudo go install ${CURDIR}/cmd/modelgen &&\
	GOBIN=/usr/local/bin/ sudo go install ${CURDIR}/cmd/domaingen

remove:
	rm -rf ${HOME}/go/bin/modelgen;\
	rm -rf ${HOME}/go/bin/domaingen;\
	rm -rf /usr/local/bin/modelgen;\
	rm -rf /usr/local/bin/domaingen