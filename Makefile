help:
	@echo Weave Makefile Commands
	@echo 'help' - Displays the command usage
	@echo 'build' - Builds the application into a binary/executable
	@echo 'install' - Installs the application
	@echo 'build-windows' - Builds the application for Windows platforms
	@echo 'build-darwin' - Builds the application for MacOSX platforms
	@echo 'build-linux' - Builds the application for Linux platforms
	@echo 'build-all' - Builds the application for all platforms
	@echo 'genprotos' - Generate protocol buffer files for the protos pkg

build:
	@echo Compiling Weave
	@go build .
	@echo Compile Complete. Run './weave(.exe)'

install:
	@echo Installing Weave
	@go install .
	@echo install Complete. Run 'weave'.

build-windows:
	@echo Cross Compiling Weave for Windows x86
	@GOOS=windows GOARCH=386 go build -o ./bin/weave-windows-x32.exe

	@echo Cross Compiling Weave for Windows x64
	@GOOS=windows GOARCH=amd64 go build -o ./bin/weave-windows-x64.exe

build-darwin:
	@echo Cross Compiling Weave for MacOSX x64
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/weave-darwin-x64

build-linux:
	@echo Cross Compiling Weave for Linux x32
	@GOOS=linux GOARCH=386 go build -o ./bin/weave-linux-x32

	@echo Cross Compiling Weave for Linux x64
	@GOOS=linux GOARCH=amd64 go build -o ./bin/weave-linux-x64

	@echo Cross Compiling Weave for Linux Arm32
	@GOOS=linux GOARCH=arm go build -o ./bin/weave-linux-arm32

	@echo Cross Compiling Weave for Linux Arm64
	@GOOS=linux GOARCH=arm64 go build -o ./bin/weave-linux-arm64

genprotos:
	@echo Generating Code for the protos library
	@protoc --go_out ./protos protos/entity.proto
	@protoc --go_out ./protos protos/query.proto
	@protoc --go_out ./protos protos/response.proto
	@protoc --go_out ./protos protos/message.proto
	@echo Code Generation Complete for protos library