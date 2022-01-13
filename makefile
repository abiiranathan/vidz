linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-s -w" -o vidz main.go 

windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags "-s -w" -o vidz.exe main.go 

darwin:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-s -w" -o vidz-mac main.go 
