.PHONY:

open:
	open ${HOME}/go/bin

ls:
	ls ${HOME}/go/bin

test:
	make install &&\
	go generate -v ./...

debug:
	make install &&\
	esc-gen-model -help

yapi: 
	make install &&\
	yapi ./...

install:
	go install ${CURDIR}/cmd/inspector &&\
	go install ${CURDIR}/cmd/esc-gen-model

remove:
	rm -rf ${HOME}/go/bin/inspector;\
	rm -rf ${HOME}/go/bin/esc-gen-model