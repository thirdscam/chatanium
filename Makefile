sqlc:
	sqlc generate -f database/sqlc.yaml

run:
	go1.23.2 run main.go

build:
	GOOS=freebsd GOARCH=386 go1.23.2 build -o bin/chatanium-b$(date +%s)-freebsd-i386 main.go
	GOOS=linux GOARCH=386 go1.23.2 build -o bin/chatanium-b$(date +%s)-linux-i386 main.go
	GOOS=windows GOARCH=386 go1.23.2 build -o bin/chatanium-b$(date +%s)-windows-i386 main.go

build_modules:
	rm -rf ./modules/*.so > /dev/null 2>&1
	for dir in $(shell find ./modules -type d -mindepth 1 -maxdepth 1); do \
		( \
			cd $$dir && \
			go1.23.2 build -buildmode=plugin -o $$(basename $$dir).so . && \
			mv $$(basename $$dir).so ../$$(basename $$dir).so \
		); \
	done

get_lines:
	find . -name '*.go' -not -path "./src/Database/Internal/*" | xargs wc -l | sort -nr