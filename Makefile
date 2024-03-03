.PHONY:

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
	esc-model-gen -h &&\
	esc-domain-gen -h

install:
	go install ${CURDIR}/cmd/inspector &&\
	go install ${CURDIR}/cmd/esc-model-gen &&\
	go install ${CURDIR}/cmd/esc-domain-gen

remove:
	rm -rf ${HOME}/go/bin/inspector;\
	rm -rf ${HOME}/go/bin/esc-model-gen;\
	rm -rf ${HOME}/go/bin/esc-domain-gen