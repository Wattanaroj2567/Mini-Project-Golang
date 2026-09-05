# Shop Service

บริการสำหรับจัดการแคตตาล็อกสินค้า ระบบตะกร้าสินค้า และคำสั่งซื้อของโปรเจกต์ **GameGear E-commerce** โดยทำหน้าที่เป็น **Data Owner** สำหรับข้อมูลการค้าทั้งหมด

> สำหรับภาพรวมสถาปัตยกรรมทั้งระบบและการติดตั้ง Kong Gateway โปรดดูที่ [Root README](../README.md) และ [docs/architecture.md](../docs/architecture.md)

---

## 1. Architectural Role

`shop-service` มีบทบาทและความรับผิดชอบหลักในระบบดังนี้:

- **Data Owner**: ถือครองฐานข้อมูลและตารางที่เกี่ยวข้องกับการซื้อขาย ได้แก่ `products`, `carts`, `cart_items`, `orders`, และ `order_items`
- **Role-based Access Control**: มี Middleware ภายในเพื่อตรวจสอบ JWT Token และ Role ของผู้ใช้งาน (อนุญาตเฉพาะคำขอที่มี `role: "admin"` ในการจัดการสินค้าและปรับสถานะคำสั่งซื้อ)
- **Data Isolation**: ไม่อนุญาตให้ Service อื่นเข้าถึงฐานข้อมูลโดยตรง การดึงหรือแก้ไขข้อมูลต้องทำผ่าน REST API เท่านั้น

---

## 2. API Endpoints

### 2.1 Member Endpoints
| Method | Path | Auth Required | คำอธิบาย |
| :--- | :--- | :---: | :--- |
| `GET` | `/api/products` | No | ดูรายการสินค้า (รองรับ Pagination & Filter) |
| `GET` | `/api/products/:id` | No | ดูรายละเอียดสินค้าเฉพาะรายการ |
| `GET` | `/api/cart` | Yes (Member JWT) | ดูรายการสินค้าในตะกร้าของผู้ใช้ |
| `POST` | `/api/cart/add` | Yes (Member JWT) | เพิ่มสินค้าลงในตะกร้า |
| `PUT` | `/api/cart/update` | Yes (Member JWT) | แก้ไขจำนวนสินค้าในตะกร้า |
| `DELETE` | `/api/cart/remove` | Yes (Member JWT) | ลบรายการสินค้าออกจากตะกร้า |
| `POST` | `/api/orders` | Yes (Member JWT) | สั่งซื้อสินค้าจากตะกร้าปัจจุบัน |
| `GET` | `/api/orders/history` | Yes (Member JWT) | ดูประวัติคำสั่งซื้อของตนเอง |

### 2.2 Admin Endpoints
| Method | Path | Auth Required | คำอธิบาย |
| :--- | :--- | :---: | :--- |
| `POST` | `/api/products` | Yes (Admin JWT) | เพิ่มรายการสินค้าใหม่เข้าสู่ระบบ |
| `PUT` | `/api/products/:id` | Yes (Admin JWT) | แก้ไขรายละเอียดสินค้าเดิม |
| `DELETE` | `/api/products/:id` | Yes (Admin JWT) | ลบรายการสินค้าออกจากระบบ |
| `GET` | `/api/orders` | Yes (Admin JWT) | ดูรายการคำสั่งซื้อทั้งหมดในระบบ |
| `PUT` | `/api/orders/:id/status` | Yes (Admin JWT) | ปรับปรุงสถานะคำสั่งซื้อ |
| `GET` | `/healthz` | No | ตรวจสอบสถานะการทำงานของ Service (Healthcheck) |

> ตัวอย่าง Request และ Response Payload ฉบับเต็ม ดูได้ที่ [docs/api-reference.md](../docs/api-reference.md#2-shop-service-endpoints-shop)

---

## 3. Project Structure

```
shop-service/
├── cmd/
│   └── api/
│       └── main.go                 # จุดเริ่มต้น เชื่อมต่อ DB และตั้งค่า Routes
├── internal/
│   ├── handlers/                   # HTTP Request Handlers (Cart, Order, Product)
│   ├── middleware/                 # JWT Authentication & Role Authorization Middleware
│   ├── models/                     # GORM Database Models และ Request/Response DTOs
│   ├── repositories/               # Database Access Layer (GORM)
│   ├── security/                   # Token Validator Utility
│   └── services/                   # Business Logic Layer
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
| `PORT` | `8082` | พอร์ตสำหรับรัน Shop Service |
| `DB_HOST` | `localhost` | ที่อยู่โฮสต์ของฐานข้อมูล PostgreSQL |
| `DB_PORT` | `5432` | พอร์ตของฐานข้อมูล PostgreSQL |
| `DB_USER` | `postgres` | ชื่อผู้ใช้ฐานข้อมูล |
| `DB_PASSWORD` | `postgres` | รหัสผ่านฐานข้อมูล |
| `DB_NAME` | `shop_db` | ชื่อฐานข้อมูลสำหรับ Shop Service |
| `JWT_SECRET_KEY` | - | กุญแจลับสำหรับตรวจสอบ JWT (ต้องตรงกับ `users-service`) |

---

## 5. Getting Started

1. **ติดตั้ง Dependencies**:
   ```bash
   go mod tidy
   ```
2. **รัน Service (พร้อม Auto-migration ของ GORM)**:
   ```bash
   go run cmd/api/main.go
   ```
3. **ตรวจสอบความพร้อมการทำงาน**:
   ```bash
   curl -i http://localhost:8082/healthz
   ```
