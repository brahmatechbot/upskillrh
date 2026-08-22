.PHONY: test build run db-setup

test:
	cd backend && go test ./...

build:
	cd backend && go build ./...

run:
	cd backend && go run .

db-setup:
	cd backend && ./scripts/setup_dev_db.sh
