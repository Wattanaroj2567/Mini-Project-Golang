# Admin Service

บริการสำหรับจัดการระบบหลังบ้าน (Admin Panel) ของโปรเจกต์ **GameGear E-commerce** โดยทำหน้าที่เป็น **Coordinator (Admin Gateway)** เพื่อประสานงานกับบริการอื่นๆ ผ่าน Kong API Gateway

> สำหรับภาพรวมสถาปัตยกรรมทั้งระบบและการติดตั้ง Kong Gateway โปรดดูที่ [Root README](../README.md) และ [docs/architecture.md](../docs/architecture.md)

---

## 1. Architectural Role

`admin-service` ถูกออกแบบตามรูปแบบ Coordinator Pattern โดยมีลักษณะสำคัญดังนี้:

- **ไม่มีฐานข้อมูลเป็นของตนเอง (No Database)**: ไม่มีการเชื่อมต่อฐานข้อมูลโดยตรง เพื่อรักษาหลักการ Data Ownership ของแต่ละ Microservice
- **API Coordinator**: ทำหน้าที่รับคำขอจาก Admin Dashboard และแปลงคำขอเพื่อเรียกต่อไปยัง `users-service` และ `shop-service` ผ่าน Kong Gateway
- **Stateless Operation**: ไม่จัดเก็บ Session หรือ Token โดยตรง แต่จะใช้ JWT ที่ออกโดย `users-service` ในการยืนยันตัวตนข้ามบริการ

---

## 2. API Endpoints

ทุก Endpoint ภายในบริการนี้ออกแบบมาสำหรับบทบาทผู้ดูแลระบบ:

| Method | Path | Auth Required | คำอธิบาย |
| :--- | :--- | :---: | :--- |
| `POST` | `/api/admin/register` | No | สมัครบัญชีแอดมินใหม่ (Proxy ไปยัง `users-service`) |
| `POST` | `/api/admin/login` | No | เข้าสู่ระบบแอดมิน (Proxy ไปยัง `users-service`) |
| `POST` | `/api/admin/forgot-password` | No | ร้องขอการรีเซ็ตรหัสผ่านแอดมิน |
| `POST` | `/api/admin/reset-password` | No | ตั้งรหัสผ่านแอดมินใหม่ด้วย Reset Token |
| `POST` | `/api/admin/logout` | Yes (Admin JWT) | ออกจากระบบแอดมิน |
| `GET` | `/api/admin/products` | Yes (Admin JWT) | ดึงรายการสินค้าทั้งหมดผ่าน `shop-service` |
| `POST` | `/api/admin/products` | Yes (Admin JWT) | เพิ่มสินค้าใหม่ผ่าน `shop-service` |
| `PUT` | `/api/admin/products/:id` | Yes (Admin JWT) | ปรับปรุงข้อมูลสินค้าผ่าน `shop-service` |
| `DELETE` | `/api/admin/products/:id` | Yes (Admin JWT) | ลบสินค้าออกจากระบบผ่าน `shop-service` |
| `GET` | `/api/admin/orders` | Yes (Admin JWT) | ดึงรายการคำสั่งซื้อทั้งหมดของระบบ |
| `PUT` | `/api/admin/orders/:id/status` | Yes (Admin JWT) | อัปเดตสถานะของคำสั่งซื้อ |
| `GET` | `/healthz` | No | ตรวจสอบสถานะการทำงานของ Service (Healthcheck) |

> ตัวอย่าง Request และ Response Payload ฉบับเต็ม ดูได้ที่ [docs/api-reference.md](../docs/api-reference.md#3-admin-service-coordinator-endpoints-admin)

---

## 3. Project Structure

```
admin-service/
├── cmd/
│   └── api/
│       └── main.go                 # จุดเริ่มต้นและตั้งค่า Routing
├── internal/
│   ├── clients/                    # HTTP Clients สำหรับเรียก Upstream Services
│   │   ├── auth_client.go
│   │   └── shop_service_client.go
│   ├── handlers/                   # Controller จัดการ HTTP Requests
│   │   ├── auth_handler.go
│   │   ├── order_handler.go
│   │   ├── product_handler.go
│   │   └── routes.go
│   ├── models/                     # DTO & Contract Models
│   │   ├── auth.go
│   │   ├── order.go
│   │   └── product.go
│   └── services/                   # Business Logic & Orchestration
│       ├── auth_service.go
│       ├── order_service.go
│       ├── product_service.go
│       └── errors.go
├── docker-compose.kong.yml         # โครงสร้างพื้นฐานของ Kong & Konga
├── .env.example
├── go.mod
└── go.sum
```

---

## 4. Environment Configuration

คัดลอกไฟล์คอนฟิกตัวอย่างและกำหนดค่าตามสภาพแวดล้อมของคุณ:

```bash
cp .env.example .env
```

| Variable | Default Value | คำอธิบาย |
| :--- | :--- | :--- |
| `APPLICATION_PORT` | `8083` | พอร์ตสำหรับรัน Admin Service |
| `SHOP_SERVICE_URL` | `http://localhost:8000/shop` | URL สำหรับเข้าถึง Shop Service ผ่าน Kong |
| `ADMIN_AUTH_SERVICE_URL` | `http://localhost:8000/users` | URL สำหรับเข้าถึง Users Service ผ่าน Kong |
| `FRONTEND_URL` | `http://localhost:3000` | URL สำหรับระบบหน้าบ้านของผู้ดูแล |

---

## 5. Getting Started

1. **ติดตั้ง Dependencies**:
   ```bash
   go mod tidy
   ```
2. **รัน Service**:
   ```bash
   go run cmd/api/main.go
   ```
3. **ตรวจสอบความพร้อมการทำงาน**:
   ```bash
   curl -i http://localhost:8083/healthz
   ```
