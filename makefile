migrate-install:
	curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash

migrate-up:
	migrate -path ./schema -database 'postgres://$(db_user):$(db_password)@$(db_host):$(db_port)/$(db_name)?sslmode=disable' -verbose up

migrate-down:
	migrate -path ./schema -database 'postgres://$(db_user):$(db_password)@$(db_host):$(db_port)/$(db_name)?sslmode=disable' down -all
