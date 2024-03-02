.PHONY:

open:
	open ${HOME}/go/bin

ls:
	ls ${HOME}/go/bin

test:
	make install &&\
	go generate -v ./...

install:
	go install ${CURDIR}/cmd/inspector &&\
	go install ${CURDIR}/cmd/esc-gen-model &&\
	go install ${CURDIR}/cmd/esc-gen-domain

remove:
	rm -rf ${HOME}/go/bin/inspector;\
	rm -rf ${HOME}/go/bin/esc-gen-model;\
	rm -rf ${HOME}/go/bin/esc-gen-domain