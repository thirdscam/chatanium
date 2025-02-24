sqlc:
	sqlc generate -f database/sqlc.yaml

run:
	go run main.go start

start:
	go run main.go start

new_module:
	go run main.go new

build:
	GOOS=freebsd GOARCH=386 go build -o bin/chatanium-b$(date +%s)-freebsd-i386 main.go
	GOOS=linux GOARCH=386 go build -o bin/chatanium-b$(date +%s)-linux-i386 main.go
	GOOS=windows GOARCH=386 go build -o bin/chatanium-b$(date +%s)-windows-i386 main.go

build_modules:
	rm -rf ./modules/*.so
	for dir in $$(find ./modules -mindepth 1 -maxdepth 1 -type d -o -type l); do \
		cd "$$dir" && \
		go build -buildmode=plugin -o "$$(basename $$dir).so" . && \
		cd ../.. && \
		mv "$$dir/$$(basename $$dir).so" "./modules/$$(basename $$dir).so"; \
	done

get_lines:
	find . -name '*.go' -not -path "./src/Database/Internal/*" | xargs wc -l | sort -nr