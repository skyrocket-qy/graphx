clear-db:
	if [ -e "./gorm.db" ]; then \
		rm "./gorm.db"; \
		echo "Database file ./gorm.db deleted."; \
	else \
		echo "Database file ./gorm.db does not exist. Nothing to delete."; \
	fi

run:
	go run .

test-run: clear-db run

gen-apidoc:
	swag init -g internal/delivery/*

build-image:
	docker build -t hrbac .

run-container:
	docker run -p 3000:3000 hrbac