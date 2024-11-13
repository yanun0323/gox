.PHONY:

CURDIR = $(shell printf "%q\n" "$(PWD)")

open:
	open /usr/local/bin/

ls:
	ls /usr/local/bin/

run:
	make install &&\
	make clean.test &&\
	go generate ./...

run.debug:
	make install &&\
	make clean.test &&\
	go generate -v ./...

clean.test:
	rm -rf ./example_output ;\
	rm -rf ./example/output ;\
	rm -f ./example/same_folder_file.go


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