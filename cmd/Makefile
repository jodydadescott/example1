default:
	echo "make what? linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64, or all?"
	exit 2

linux-amd64:
	mkdir -p build
	env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o build/linux-amd64 cmd.go 

linux-arm64:
	mkdir -p build
	env GOOS=linux CGO_ENABLED=0 GOARCH=arm64 go build -o build/linux-arm64 cmd.go

darwin-amd64:
	mkdir -p build
	env GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64 cmd.go 

darwin-arm64:
	mkdir -p build
	env GOOS=darwin GOARCH=arm64 go build -o build/darwin-arm64 cmd.go  

windows-amd64:
	mkdir -p build
	env GOOS=windows GOARCH=amd64 go build -o build/windows-amd64 cmd.go 

all:
	$(MAKE) linux-amd64
	$(MAKE) linux-arm64
	$(MAKE) darwin-amd64
	$(MAKE) darwin-arm64
	$(MAKE) windows-amd64

clean:
	$(RM) -rf build