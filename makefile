PHONY:gomod_init gotest gotest_debug

gomod_init:
	go mod init github.com/donnol/tools

tooldebug=
gotest:
	TOOLDEBUG=$(tooldebug) gotest -v ./...

gotest_debug:
	TOOLDEBUG=true gotest -v ./...

tbc_build_test:
	cd cmd/tbc && \
	go build && \
	./tbc interface -p github.com/donnol/tools/importpath  -r

tbc_install:
	cd cmd/tbc && \
	go install
