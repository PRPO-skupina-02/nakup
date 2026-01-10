docs:
	swag fmt
	swag init --parseDependency -g api/api.go -o api/docs

test:
	godotenv go test ./...

test-update:
	godotenv go test ./... -update

install-cli-tools:
	go install github.com/joho/godotenv/cmd/godotenv@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest

migrate:
	migrate -database postgres://postgres:postgres@localhost:5433/nakup?sslmode=disable -path db/migrations up

migrate-down: 
	migrate -database postgres://postgres:postgres@localhost:5433/nakup?sslmode=disable -path db/migrations down

fixtures:
	godotenv go run ../common/tools/loadfixture/loadfixture.go db/fixtures/

swagger-clients:
	swagger generate client -f ../spored/api/docs/swagger.json -A spored -t ./clients/spored