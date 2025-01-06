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

release:
	@if [ ! -f VERSION ]; then \
		echo "Error: VERSION file not found"; \
	else \
		VERSION=$$(cat VERSION); \
	fi; \
	if [ -z "$$VERSION" ]; then \
		VERSION="v0.0.0"; \
	fi; \
	if ! echo "$$VERSION" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+$$" > /dev/null; then \
		echo "Error: Version format must be vX.X.X"; \
		exit 1; \
	else \
		MAJOR=$$(echo "$$VERSION" | cut -d. -f1) && \
		MINOR=$$(echo "$$VERSION" | cut -d. -f2) && \
		PATCH=$$(echo "$$VERSION" | cut -d. -f3) && \
		NEW_PATCH=$$((PATCH + 1)) && \
		NEW_VERSION="$$MAJOR.$$MINOR.$$NEW_PATCH" && \
		rm -f ./VERSION &&\
		echo "$$NEW_VERSION" > ./VERSION &&\
		echo "add tag $$NEW_VERSION"; \
		git add . && \
		git commit -m "release version $$NEW_VERSION" && \
		git tag -a "$$NEW_VERSION" -m "version $$NEW_VERSION" && \
		git push &&\
		git push --tags && \
		echo "release version"; \
		echo ""; \
		echo "$$NEW_VERSION"; \
		echo ""; \
	fi

get.next.version:
	@if [ ! -f VERSION ]; then \
		echo "Error: VERSION file not found"; \
	else \
		VERSION=$$(cat VERSION); \
	fi; \
	if [ -z "$$VERSION" ]; then \
		VERSION="v0.0.0"; \
	fi; \
	if ! echo "$$VERSION" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+$$" > /dev/null; then \
		echo "Error: Version format must be vX.X.X"; \
		exit 1; \
	else \
		MAJOR=$$(echo "$$VERSION" | cut -d. -f1) && \
		MINOR=$$(echo "$$VERSION" | cut -d. -f2) && \
		PATCH=$$(echo "$$VERSION" | cut -d. -f3) && \
		NEW_PATCH=$$((PATCH + 1)) && \
		NEW_VERSION="$$MAJOR.$$MINOR.$$NEW_PATCH" && \
		echo "next version will be"; \
		echo ""; \
		echo "$$NEW_VERSION"; \
		echo ""; \
	fi