routes -> handler -> service -> repository -> model

handler: request/response
service: logic
repository: db interaction

## Cài đặt Golang-Migrate

```curl
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

hoặc cài đặt theo tài liệu

```curl
https://github.com/golang-migrate/migrate/tree/v4.18.3/cmd/migrate
```

kiểm tra version
```curl
migrate --version
```