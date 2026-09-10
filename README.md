# Subglutee Project - Report Service

สรุป/วิเคราะห์การใช้จ่าย subscription ของ user

- REST `:8086` - frontend/gateway (`/api/v1/reports/me`)
- gRPC client -> subscription-service (`GetSubscriptionsForReport`)
- ไม่มี DB - compute จาก subscription ทุก request (stateless)
