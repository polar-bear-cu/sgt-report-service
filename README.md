# Subglutee Project - Report Service

สรุป/วิเคราะห์การใช้จ่าย subscription ของ user

- REST `:8086` - frontend/gateway (`/api/v1/reports/me`)
- gRPC client -> subscription-service (`GetSubscriptionsForReport`)
- ไม่มี DB - compute จาก subscription ทุก request (stateless)

### Structure

```
routes/         map path -> controller
controllers/    HTTP <-> DTO
usecases/       compute report จาก subscription list
clients/        gRPC client -> subscription-service (proto adapter)
dtos/           request/response struct
config/         env loader
```

### Prerequisite

- Go 1.26
- Docker
- `make` - `winget install ezwinports.make`
- tools:

```terminal
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install github.com/evilmartians/lefthook@latest
go install github.com/swaggo/swag/cmd/swag@v1.16.6
```

### Setup

```terminal
git clone https://github.com/polar-bear-cu/sgt-report-service.git
cd sgt-report-service
lefthook install
cp .env.example .env
go mod download
make run
```

### Useful Commands

Check `Makefile`

### API Docs (Swagger)

```terminal
make docs
```

Document is at http://localhost:8086/swagger/index.html (ต้อง `ENABLE_SWAGGER=true`)

#### ถ้าแก้ proto พร้อม service นี้

```terminal
cd ..
go work init ./sgt-proto ./sgt-report-service
```

หรือถ้ามี go.work แล้ว...

```terminal
go work use ./sgt-proto ./sgt-report-service
```
