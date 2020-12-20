PHONY:gomod_init gotest gotest_debug

gomod_init:
	go mod init github.com/donnol/tools

tooldebug=
gotest:
	TOOLDEBUG=$(tooldebug) gotest -v ./...

gotest_debug:
	TOOLDEBUG=true gotest -v ./...

tbc_build_test:tbc_install
	tbc interface -p github.com/donnol/tools/importpath  -r

tbc_install:
	cd cmd/tbc && \
	go install

tbc_callgraph_test:tbc_install
	tbc callgraph --path=github.com/donnol/tools/parser --func=UcFirst --ignore=std --depth=3 && \
	tbc callgraph --path=github.com/donnol/tools/parser --func=Func.PrintCallGraph --ignore=std --depth=3
