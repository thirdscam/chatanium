sqlc:
	sqlc generate -f database/sqlc.yaml

run:
	go run main.go

build:
	GOOS=freebsd GOARCH=386 go build -o bin/chatanium-b$(date +%s)-freebsd-i386 main.go
	GOOS=linux GOARCH=386 go build -o bin/chatanium-b$(date +%s)-linux-i386 main.go
	GOOS=windows GOARCH=386 go build -o bin/chatanium-b$(date +%s)-windows-i386 main.go

build_modules:
	go build -o ./modules -buildmode=plugin ./modules/**

get_lines:
	find . -name '*.go' -not -path "./src/Database/Internal/*" | xargs wc -l | sort -nr