# Users Service

บริการสำหรับจัดการผู้ใช้งาน การยืนยันตัวตน และการบริหารจัดการสิทธิ์ของโปรเจกต์ **GameGear E-commerce** โดยทำหน้าที่เป็น **Data Owner** สำหรับข้อมูลบัญชีผู้ใช้และเป็นผู้ออก **JWT Token** ประจำระบบ

> สำหรับภาพรวมสถาปัตยกรรมทั้งระบบและการติดตั้ง Kong Gateway โปรดดูที่ [Root README](../README.md) และ [docs/architecture.md](../docs/architecture.md)

---

## 1. Architectural Role

`users-service` มีบทบาทและความรับผิดชอบหลักในระบบดังนี้:

- **Centralized Identity & Access Management**: จัดการการสมัครสมาชิก, การล็อกอิน, การขอและรีเซ็ตรหัสผ่าน สำหรับทั้งสมาชิกทั่วไปและผู้ดูแลระบบ
- **Data Owner**: ถือครองฐานข้อมูลและตารางที่เกี่ยวข้องกับความปลอดภัย ได้แก่ `users` และ `password_reset_tokens`
- **JWT Issuer**: เป็นผู้สร้างและลงลายเซ็นดิจิทัลให้กับ JSON Web Token (JWT) โดยฝังข้อมูล Claim ที่สำคัญ เช่น `user_id`, `email`, และ `role` (`member` หรือ `admin`)

---

## 2. API Endpoints

### 2.1 Member Endpoints
| Method | Path | Auth Required | คำอธิบาย |
| :--- | :--- | :---: | :--- |
| `POST` | `/api/auth/register` | No | ลงทะเบียนผู้ใช้ใหม่ |
| `POST` | `/api/auth/login` | No | เข้าสู่ระบบและรับ JWT Token |
| `POST` | `/api/auth/forgot-password` | No | ขอรับคำร้องรีเซ็ตรหัสผ่าน |
| `POST` | `/api/auth/reset-password` | No | กำหนดรหัสผ่านใหม่ด้วย Reset Token |
| `POST` | `/api/auth/logout` | Yes (Member JWT) | ออกจากระบบและเพิกถอน Token |
| `GET` | `/api/user/profile` | Yes (Member JWT) | ดูข้อมูลโปรไฟล์ของผู้ใช้ปัจจุบัน |
| `PUT` | `/api/user/profile` | Yes (Member JWT) | แก้ไขข้อมูลส่วนตัวหรือเปลี่ยนรหัสผ่าน |
| `DELETE` | `/api/user/profile` | Yes (Member JWT) | ขอลบบัญชีผู้ใช้ |

### 2.2 Admin Endpoints (สำหรับ admin-service)
| Method | Path | Auth Required | คำอธิบาย |
| :--- | :--- | :---: | :--- |
| `POST` | `/api/admin/register` | No | สร้างบัญชีแอดมินใหม่ (กำหนด `role: "admin"`) |
| `POST` | `/api/admin/login` | No | เข้าสู่ระบบแอดมินและรับ JWT ที่มีสิทธิ์ Admin |
| `POST` | `/api/admin/forgot-password` | No | เริ่มกระบวนการรีเซ็ตรหัสผ่านสำหรับแอดมิน |
| `POST` | `/api/admin/reset-password` | No | ตั้งรหัสผ่านแอดมินใหม่ด้วย Reset Token |
| `POST` | `/api/admin/logout` | Yes (Admin JWT) | ออกจากระบบแอดมิน |
| `GET` | `/healthz` | No | ตรวจสอบสถานะการทำงานของ Service (Healthcheck) |

> ตัวอย่าง Request และ Response Payload ฉบับเต็ม ดูได้ที่ [docs/api-reference.md](../docs/api-reference.md#1-users-service-endpoints-users)

---

## 3. Project Structure

```
users-service/
├── cmd/
│   └── api/
│       └── main.go                 # จุดเริ่มต้น เชื่อมต่อ DB และตั้งค่า Routes
├── internal/
│   ├── handlers/                   # HTTP Request Handlers (Auth, Profile)
│   ├── models/                     # GORM Database Models และ Request/Response DTOs
│   ├── repositories/               # Database Access Layer (GORM)
│   └── services/                   # Business Logic & Token Issuance Services
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
| `PORT` | `8081` | พอร์ตสำหรับรัน Users Service |
| `DB_HOST` | `localhost` | ที่อยู่โฮสต์ของฐานข้อมูล PostgreSQL |
| `DB_PORT` | `5432` | พอร์ตของฐานข้อมูล PostgreSQL |
| `DB_USER` | `postgres` | ชื่อผู้ใช้ฐานข้อมูล |
| `DB_PASSWORD` | `postgres` | รหัสผ่านฐานข้อมูล |
| `DB_NAME` | `users_db` | ชื่อฐานข้อมูลสำหรับ Users Service |
| `JWT_SECRET_KEY` | - | กุญแจลับสำหรับเซ็นและตรวจสอบ JWT |
| `JWT_EXPIRY_HOURS` | `24` | อายุการใช้งานของ Access Token (ชั่วโมง) |

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
   curl -i http://localhost:8081/healthz
   ```
