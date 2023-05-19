.PHONY: run test updatedep d-up d-down

updatedep: 
	go get -u ./...
	go generate ./...
run:
	go run .
test:
	go test -v ./...

d-up:
	docker-compose -f docker-compose.yaml up -d
d-down:
	docker-compose -f docker-compose.yaml down -v
