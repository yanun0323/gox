.PHONY:

open:
	open ${HOME}/go/bin

ls:
	ls ${HOME}/go/bin

test:
	make install &&\
	go generate -v ./...

help:
	make install &&\
	esc-modelgen -h &&\
	esc-domaingen -h

install:
	go install ${CURDIR}/cmd/inspector &&\
	go install ${CURDIR}/cmd/esc-modelgen &&\
	go install ${CURDIR}/cmd/esc-domaingen

remove:
	rm -rf ${HOME}/go/bin/inspector;\
	rm -rf ${HOME}/go/bin/esc-modelgen;\
	rm -rf ${HOME}/go/bin/esc-domaingen