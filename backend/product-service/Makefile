
migrate-up:
	go run migration/cmd/main.go migrate up

migrate-down:
	go run migration/cmd/main.go migrate down

migrate-status:
	go run migration/cmd/main.go migrate status


migrate-create:
	@read -p "Migration açıklaması: " desc; \
	./scripts/create_migration.sh "$$desc"

seed:
	go run migration/cmd/main.go seed

.PHONY: devdb migrate-up migrate-down migrate-status migrate-create
