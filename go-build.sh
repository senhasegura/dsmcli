# Build for Linux Alpine
env GOOS=linux GARCH=amd64 go build -v -a -tags cgo -installsuffix cgo -ldflags "-linkmode external -extldflags -static" -o build/linux/dsm .

# Build for Windows
env GOOS=windows CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -v -o build/windows/dsm.exe .