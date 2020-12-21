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
	# 包括func的比对组，ignore的比对组，depth的比对组
	tbc callgraph --path=github.com/donnol/tools/parser --func=UcFirst --depth=3 && \
	tbc callgraph --path=github.com/donnol/tools/parser --func=Func.PrintCallGraph --depth=2 && \
	tbc callgraph --path=github.com/donnol/tools/parser --func=Func.PrintCallGraph --depth=3 && \
	tbc callgraph --path=github.com/donnol/tools/parser --func=Func.PrintCallGraph --ignore=std --depth=3 && \
	tbc callgraph --path=github.com/donnol/tools/parser --func=Func.PrintCallGraph --ignore='std;github.com/donnol/tools/importpath' --depth=3
