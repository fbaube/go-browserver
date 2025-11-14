cd ./cmd/wasm/
GOOS=js GOARCH=wasm go build -o main.wasm . ; cp main.wasm ../../public
cd ../server
go run main.go
