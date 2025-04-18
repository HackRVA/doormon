
all: doormon arm

.PHONY:doormon
doormon:
	go build -o dist/doormon

.PHONY:arm
arm:
	GOOS=linux GOARCH=arm GOARM=6 go build -o dist/arm/doormon

run-broker:
	docker compose -f ./deployments/docker-compose.yml up -d

stop-broker:
	docker compose -f ./deployments/docker-compose.yml down

clean:
	rm -rf dist
