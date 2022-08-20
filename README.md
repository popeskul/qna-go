# qna-go

TODO:
- [x] linter
- [x] unit tests
- [x] authorization and authentication
- [x] docker
- [ ] docker-compose
- [ ] query cache
- [ ] swagger
- [ ] postman script for testing

## Installation docker db
```bash
docker run --name=qna_db      -e POSTGRES_PASSWORD=12345 -p 5432:5432 -d --rm postgres
docker run --name=qna_test_db -e POSTGRES_PASSWORD=12345 -p 5436:5432 -d --rm postgres
```

```bash
make migrate-up db_user=postgres db_password=12345 db_host=localhost db_port=5432 db_name=postgres
make migrate-up db_user=postgres db_password=12345 db_host=localhost db_port=5432 db_name=postgres
```
