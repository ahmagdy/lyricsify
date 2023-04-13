.PHONY: run test updatedep

updatedep: 
	go get -u ./...
run:
	go run .
test:
	go test -v ./...