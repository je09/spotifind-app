# Build the CLI for all platforms
cli-darwin64:
	echo "Building CLI for Darwin"
	GOOS=darwin GOARCH=amd64 go build -o bin/spotifind-cli-darwin ./cli

cli-darwinarm64:
	echo "Building CLI for Darwin"
	GOOS=darwin GOARCH=arm64 go build -o bin/spotifind-cli-darwin-arm64 ./cli

cli-linux64:
	echo "Building CLI for Linux"
	GOOS=linux GOARCH=amd64 go build -o bin/spotifind-cli-linux ./cli

cli-linuxarm64:
	echo "Building CLI for Linux"
	GOOS=linux GOARCH=arm64 go build -o bin/spotifind-cli-linux-arm64 ./cli

cli-windows64:
	echo "Building CLI for Windows"
	GOOS=windows GOARCH=amd64 go build -o bin/spotifind-cli-windows.exe ./cli

cli-windowsarm64:
	echo "Building CLI for Windows"
	GOOS=windows GOARCH=arm64 go build -o bin/spotifind-cli-windows-arm64.exe ./cli

cli-all: cli-darwin64 cli-darwinarm64 cli-linux64 cli-linuxarm64 cli-windows64 cli-windowsarm64
cli: cli-all

build: cli