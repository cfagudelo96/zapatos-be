build-windows:
	env GOOS=windows GOARCH=amd64 go build -o bin/main.exe main.go
