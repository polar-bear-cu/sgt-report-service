# Subglutee Project - Report Service

สรุป/วิเคราะห์การใช้จ่าย subscription ของ user

- REST `:8086` - frontend/gateway (`/api/v1/reports/...`)
- gRPC client -> subscription-service (`ListCharges`)
- ไม่มี DB - compute จาก charges ทุก request (stateless), charges ถูกบันทึกโดย subscription-service ตอน billing rollover

### Structure

```
routes/         map path -> controller
middlewares/    JWT auth (ตั้ง user_id ใน context)
controllers/    HTTP <-> DTO, CSV export
usecases/       compute report จาก charges (แยกเดือนตามเวลา Asia/Bangkok)
clients/        gRPC client -> subscription-service (proto adapter)
dtos/           request/response struct
config/         env loader
```

### API

ทุก endpoint ใต้ `/api/v1` ต้องแนบ `Authorization: Bearer <access token>` (JWT จาก sgt-auth-service, HS256, ใช้ `JWT_SECRET` เดียวกัน) ไม่มี / ไม่ถูกต้อง / หมดอายุ -> `401`
ทุก endpoint เห็นเฉพาะ charges ของ user เจ้าของ token (`sub`) เท่านั้น

| Method | Path | ใช้ทำอะไร |
| --- | --- | --- |
| `GET` | `/api/v1/reports/monthly` | ยอดรายเดือน + ยอดรวม + ค่าเฉลี่ย |
| `GET` | `/api/v1/reports/subscriptions` | ยอดรวมของแต่ละ subscription ในช่วงเดือน |
| `GET` | `/api/v1/reports/charges` | charges ของ 1 เดือน + ยอดแยก category |
| `GET` | `/api/v1/reports/export` | ดาวน์โหลด charges ในช่วงเดือนเป็น CSV |
| `GET` | `/health` | health check (ไม่ต้องใช้ token) |

- เดือนเป็น `YYYY-MM` เวลา Asia/Bangkok ผิดรูปแบบ -> `400`
- charge = การตัดเงิน 1 ครั้ง, `type` / `amount` เป็นค่า ณ วันตัดเงิน, `name` / `category` / `status` เป็นค่าปัจจุบันของ subscription
- subscription-service ใช้งานไม่ได้ -> `502`

#### Monthly

| Param | ค่า | Default |
| --- | --- | --- |
| `months` | 1 - 24 | `12` |

item สุดท้ายคือเดือนปัจจุบันเสมอ (ยังไม่จบ), เดือนที่ไม่มี charge -> `0`, `averagePerMonth = total / months`

```json
{ "months": [{ "month": "2026-10", "total": 139 }], "total": 17691, "averagePerMonth": 1474.25 }
```

#### Subscriptions / Export

| Param | ค่า | Default |
| --- | --- | --- |
| `from` | `YYYY-MM` (รวม) | 11 เดือนก่อน `to` |
| `to` | `YYYY-MM` (รวม) | เดือนปัจจุบัน |

ช่วงยาวสุด 24 เดือน, subscriptions เรียงจากยอดมากไปน้อย, `type` เป็นของ charge ล่าสุด

```json
{ "from": "2025-11", "to": "2026-10", "total": 17691, "items": [{ "subscriptionId": "...", "name": "Netflix", "category": "streaming", "type": "monthly", "status": "active", "count": 11, "total": 4609 }] }
```

export ได้ไฟล์ `sgt-charges-<from>_<to>.csv` (UTF-8 + BOM สำหรับ Excel)

```csv
date,subscription,category,type,amount
2026-09-05,Spotify,music,monthly,139.00
```

#### Charges

| Param | ค่า | Default |
| --- | --- | --- |
| `month` | `YYYY-MM` | เดือนปัจจุบัน |

items เรียงจากเก่าไปใหม่, categories เรียงจากยอดมากไปน้อย

```json
{ "month": "2026-09", "total": 1462, "count": 5, "categories": [{ "category": "technology", "total": 725 }], "items": [{ "subscriptionId": "...", "name": "Spotify", "category": "music", "type": "monthly", "amount": 139, "chargedAt": "2026-09-05T00:00:00+07:00" }] }
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

### Run alternatively (container)

```terminal
make image
make container
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
