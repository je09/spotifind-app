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

# Build GUI for all platforms
gui-darwin64:
	echo "Building GUI for Darwin"
	wails build -platform darwin/amd64 -o spotifind-gui-macos
	mv ./build/bin/spotifind-gui.app ./bin/spotifind-gui-macos.app

# Build GUI for all platforms
gui-darwinarm64:
	echo "Building GUI for Darwin"
	wails build -platform darwin/arm64 -o spotifind-gui-macos-arm64
	mv ./build/bin/spotifind-gui.app ./bin/spotifind-gui-macos-arm64.app

gui-windows64:
	echo "Building GUI for Windows"
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CGO_CXXFLAGS="-IC:\msys64\mingw64\include" wails build -ldflags '-extldflags "-static"' -skipbindings

# unstable for now
gui-windowsarm64:
	echo "Building GUI for Windows"
	wails build -platform windows/arm64 -o spotifind-gui-windows-arm64

gui-linux64:
	echo "Building GUI for Linux"
	wails build -platform linux/amd64 -o spotifind-gui-linux

gui-linuxarm64:
	echo "Building GUI for Linux"
	wails build -platform linux/arm64 -o spotifind-gui-linux-arm64

gui-darwin: gui-darwin64 gui-darwinarm64
gui-win: gui-windows64
gui-linux: gui-linux64 gui-linuxarm64

gui: gui-linux gui-win # mac is not supported for cross-compilation

build: cli gui