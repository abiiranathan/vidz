linux:
	GOOS=linux GOARCH=amd64 go build -o vidz main.go

windows:
	GOOS=windows GOARCH=amd64 go build -o vidz.exe main.go

darwin:
	GOOS=darwin GOARCH=amd64 go build -o vidz-mac main.go
