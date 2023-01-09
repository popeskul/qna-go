# qna-go

TODO:
- [x] linter
- [x] unit tests
- [x] authorization and authentication
- [x] docker
- [ ] docker-compose
- [x] query cache
- [x] swagger
- [ ] postman script for testing]
- [x] grpc logger

## Installation docker db
```bash
docker run --name=qna_db      -e POSTGRES_PASSWORD=12345 -p 5432:5432 -d --rm postgres
docker run --name=qna_test_db -e POSTGRES_PASSWORD=12345 -p 5436:5432 -d --rm postgres

docker run --rm -d --name audit-logo-mongo -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=123 -p 27017:27017 mongo:latest
docker run --rm -d --name audit-logo-mongo-test -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=123 -p 27011:27017 mongo:latest
```

```bash
make migrate-up db_user=postgres db_password=12345 db_host=localhost db_port=5432 db_name=postgres
make migrate-up db_user=postgres db_password=12345 db_host=localhost db_port=5436 db_name=postgres
```
