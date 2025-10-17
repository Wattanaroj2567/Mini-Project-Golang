# 👤 Users Service — GameGear E-commerce

บริการสำหรับจัดการผู้ใช้และการยืนยันตัวตนทั้งหมดภายในระบบ **GameGear E-commerce**
เป็นหนึ่งใน microservices ที่ทำงานร่วมกับ **Kong API Gateway**

> 📖 **ดูเอกสารหลักของระบบ**: สำหรับ Kong Gateway setup, Architecture overview และการ integrate ทั้งระบบ → [Main README](https://github.com/Wattanaroj2567/Mini-Project-Golang)

---

## 📋 Table of Contents

- [🏛️ ภาพรวมระบบ (System Overview)](#%EF%B8%8F-ภาพรวมระบบ-system-overview)
- [✨ คุณสมบัติและ Endpoints](#-คุณสมบัติและ-endpoints)
- [📂 โครงสร้างโปรเจค (Project Structure)](#-โครงสร้างโปรเจค-project-structure)
- [📦 Module Structure](#-module-structure)
- [🚀 เริ่มต้นใช้งาน (Getting Started)](#-เริ่มต้นใช้งาน-getting-started)
- [📝 API Documentation](#-api-documentation)
- [📞 ติดต่อและสนับสนุน](#-ติดต่อและสนับสนุน)

---

## 🏛️ ภาพรวมระบบ (System Overview)

- เป็น **เจ้าของข้อมูล (Data Owner)** สำหรับตารางผู้ใช้: `users`, `password_reset_tokens`
- ข้อมูลผู้ใช้ **ไม่ถูกเข้าถึงโดยตรงจาก service อื่น** — ต้องเรียกผ่าน API ของ Users Service เท่านั้น
- บริการอื่น (Shop, Admin) ถ้าต้องการข้อมูลผู้ใช้ จะต้องเรียก API นี้ผ่าน Gateway

---

## ✨ คุณสมบัติและ Endpoints

### สมาชิกทั่วไป (Member Endpoints)

| Feature                 | Method | Path                         | Auth Required        | Description                                             |
| ----------------------- | ------ | ---------------------------- | -------------------- | ------------------------------------------------------- |
| User Registration       | POST   | `/api/auth/register`        | No                   | ลงทะเบียนผู้ใช้ใหม่                                    |
| Login                   | POST   | `/api/auth/login`           | No                   | เข้าสู่ระบบและรับ JWT token                            |
| Logout                  | POST   | `/api/auth/logout`          | Yes (Member JWT)     | ออกจากระบบและเพิกถอนโทเคนปัจจุบัน                    |
| Forgot Password Request | POST   | `/api/auth/forgot-password` | No                   | ขออีเมลสำหรับขั้นตอนรีเซ็ตรหัสผ่าน                    |
| Reset Password          | POST   | `/api/auth/reset-password`  | No                   | ตั้งรหัสผ่านใหม่ด้วยโทเคนที่ได้รับ                     |
| View Profile            | GET    | `/api/user/profile`         | Yes (Member JWT)     | ดูรายละเอียดโปรไฟล์ของตนเอง                           |
| Update Profile          | PUT    | `/api/user/profile`         | Yes (Member JWT)     | แก้ไขข้อมูลส่วนตัว, เปลี่ยนรหัสผ่าน หรือปิดบัญชี      |

### ผู้ดูแลระบบ (Admin Endpoints) — ใช้สำหรับ admin-service เท่านั้น

| Feature                       | Method | Path                              | Auth Required     | Description                                                               |
| ----------------------------- | ------ | --------------------------------- | ----------------- | ------------------------------------------------------------------------- |
| Admin Registration            | POST   | `/api/admin/register`             | No                | สร้างบัญชีผู้ดูแลระบบใหม่ (กำหนด role = `admin`)                         |
| Admin Login                   | POST   | `/api/admin/login`                | No                | เข้าสู่ระบบสำหรับแอดมินและออก JWT ที่ฝังข้อมูล role                     |
| Admin Logout                  | POST   | `/api/admin/logout`               | Yes (Admin JWT)   | ออกจากระบบและเพิกถอนโทเคนแอดมิน                                       |
| Admin Forgot Password Request | POST   | `/api/admin/forgot-password`      | No                | เริ่มกระบวนการรีเซ็ตรหัสผ่านสำหรับแอดมิน                               |
| Admin Reset Password          | POST   | `/api/admin/reset-password`       | No                | ตั้งรหัสผ่านใหม่ด้วยโทเคนที่ได้รับ                                       |

> 🔐 **หมายเหตุ**: endpoints ชุด `/api/admin/*` ถูกเรียกใช้งานโดย `admin-service` ผ่าน Kong Gateway (เช่น `http://localhost:8000/users/api/admin/login`) และจะคืน JWT ที่ต้องมี `role=admin` เพื่อใช้ควบคุมสิทธิ์ใน services อื่น

นอกจากนี้ ควรมี endpoint สำหรับ **Healthcheck**:

```
GET /healthz → 200 OK
```

---

## 📂 โครงสร้างโปรเจค (Project Structure)

```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── profile_handler.go
│   │   └── routes.go
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── profile_service.go
│   │   └── token_service.go
│   ├── repositories/
│   │   ├── password_reset_repository.go
│   │   └── user_repository.go
│   └── models/
│       ├── auth_payloads.go
│       ├── password_reset_token.go
│       ├── profile_payloads.go
│       └── user.go
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

คำอธิบายไฟล์:

- `cmd/api/main.go` — บูตระบบทั้งหมด, เชื่อมต่อฐานข้อมูล, รัน migration และประกาศ route หลัก
- `internal/handlers/auth_handler.go` — รับคำขอจาก `/api/auth/*`, แปลง input และเรียก `AuthService`
- `internal/handlers/profile_handler.go` — จัดการ `/api/user/profile` สำหรับดูและอัปเดตข้อมูลสมาชิก
- `internal/handlers/routes.go` — รวมการประกาศ Gin route group สำหรับ auth, profile และ admin
- `internal/services/auth_service.go` — โครง business logic สำหรับ Register/Login/Forgot/Reset/Logout ทั้งสมาชิกและแอดมิน (role-based)
- `internal/services/profile_service.go` — โครง service สำหรับ get/update/delete โปรไฟล์
- `internal/services/token_service.go` — interface สำหรับการออก/เพิกถอน JWT (ยังไม่กำหนด implementation)
- `internal/repositories/user_repository.go` — ฟังก์ชัน CRUD ต่อฐานข้อมูลตาราง `users`
- `internal/repositories/password_reset_repository.go` — จัดการตาราง `password_reset_tokens` (สร้าง/ค้นหา/ลบ token)
- `internal/models/user.go` — GORM model ของผู้ใช้ในระบบ
- `internal/models/password_reset_token.go` — GORM model สำหรับ token รีเซ็ตรหัสผ่าน
- `internal/models/auth_payloads.go` — DTO ที่ใช้ bind/response สำหรับ endpoints auth
- `internal/models/profile_payloads.go` — DTO สำหรับการดูและอัปเดตโปรไฟล์ รวมถึงคำขอปิดบัญชี
- **.env.example** — ตัวอย่างไฟล์สำหรับตั้งค่าคอนฟิก เช่น Database URL, JWT, Email
- **README.md** — เอกสารอธิบายรายละเอียดการติดตั้งและใช้งาน Service

---

## 📦 Module Structure

Service นี้ใช้ Go Module สำหรับจัดการ dependencies:

| Property              | Value                               |
| --------------------- | ----------------------------------- |
| **Module Name**       | `github.com/gamegear/users-service` |
| **Go Version**        | 1.25.1                              |
| **Main Dependencies** | Gin, GORM, PostgreSQL, JWT          |

### Local Development Setup

สำหรับการพัฒนาในเครื่อง local service นี้ใช้ `replace` directive ใน `go.mod`:

```go
// ใน admin-service/go.mod และ shop-service/go.mod
replace github.com/gamegear/users-service => ../users-service
```

### Dependencies Management

- **Web Framework**: Gin (HTTP router)
- **Database**: GORM + PostgreSQL driver
- **Authentication**: JWT tokens
- **API Documentation**: Postman
- **Environment**: godotenv for .env files
- **Password Hashing**: bcrypt

### Import Statements

เมื่อต้องการ import จาก services อื่น ให้ใช้ module name แทน local path:

```go
// ✅ ถูกต้อง - ใช้ module name
import "github.com/gamegear/shop-service/internal/models"

// ❌ ผิด - อย่าใช้ local path
import "../shop-service/internal/models"
```

**หมายเหตุ**: users-service เป็น base service ที่ไม่ต้อง import จาก services อื่น แต่หากต้องการ import ให้ใช้ module name

---

## 🚀 เริ่มต้นใช้งาน (Getting Started)

### 1. Clone & เข้าโฟลเดอร์

```bash
git clone https://github.com/Wattanaroj2567/users-service.git
cd users-service
```

### 2. ติดตั้ง dependencies

```bash
go mod tidy
```

### 3. ตั้งค่า Database

ตรวจสอบว่า PostgreSQL ทำงานอยู่ แล้วสร้างฐานข้อมูล:

```sql
CREATE DATABASE gamegear_users_db;
```

### 4. สร้างไฟล์ `.env`

**สร้างไฟล์ `.env` ใหม่** โดยคัดลอกเนื้อหาจาก `.env.example` และแก้ไขค่าต่างๆ ตามต้องการ:

```env
# Application Configuration
APPLICATION_PORT=8081
APP_NAME="GameGear Users Service"
APP_VERSION="1.0.0"

# Database Configuration
DATABASE_URL="host=localhost user=postgres password=your_password dbname=gamegear_users_db port=5432 sslmode=disable"
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=gamegear_users_db

# JWT Configuration
JWT_SECRET_KEY="your_super_secret_jwt_key_here_make_it_long_and_secure"
JWT_EXPIRATION_HOURS=24

# CORS Configuration
FRONTEND_URL="http://localhost:3000"
ALLOWED_ORIGINS="http://localhost:3000,http://localhost:8080"

# Email Configuration (สำหรับ Forgot Password)
EMAIL_HOST="smtp.gmail.com"
EMAIL_PORT=587
EMAIL_USERNAME="your_email@gmail.com"
EMAIL_PASSWORD="your_app_password"
FROM_EMAIL="noreply@gamegear.com"

# Logging Configuration
LOG_LEVEL=info
LOG_FILE="logs/users-service.log"
```

**💡 หมายเหตุ:**

- **JWT_SECRET_KEY**: ต้องใช้ค่าเดียวกันกับ `shop-service` และ `admin-service` เพื่อให้สามารถตรวจสอบ JWT ข้าม service ได้ และสามารถสร้างได้จาก [jwt.io](https://jwt.io) หรือ random string generator
- **Database Password**: ใช้รหัสผ่าน PostgreSQL ของคุณ
- **Email Settings**: ใช้ Gmail App Password สำหรับการส่งอีเมล
- **ไฟล์ `.env.example`**: เก็บไว้เป็น template สำหรับครั้งต่อไป

### 5. ไฟล์ที่ต้องแก้ไข/เขียนโค้ด

> 👨‍💻 **Developer**: ณิชพน มานิตย์

คุณต้องเขียนโค้ดเฉพาะใน **2 จุดหลัก**:

#### 5.1 **`cmd/api/main.go`**

- เริ่มต้น Gin server
- Setup routes
- Connect database
- Middleware configuration

#### 5.2 **โฟลเดอร์ `internal/`**

| Folder               | ต้องทำ       | หมายเหตุ                          |
| -------------------- | ------------ | --------------------------------- |
| ✅ **handlers/**     | ✅ ต้องทำ    | เขียน HTTP handlers (controllers) |
| ❌ **models/**       | ❌ ไม่ต้องทำ | **PM ทำให้แล้ว** (วรรธนโรจน์)     |
| ✅ **repositories/** | ✅ ต้องทำ    | เขียน database operations (CRUD)  |
| ✅ **services/**     | ✅ ต้องทำ    | เขียน business logic              |

> 💡 **หมายเหตุ**: ในโค้ดจะมี **TODO comments** บอกว่าต้องทำอะไรบ้าง ให้ทำตามที่ระบุไว้ในโค้ด

#### 5.3 **ไฟล์อื่นๆ**

- ✅ **`.env`** - ตั้งค่า environment variables
- ✅ **`go.mod`** - เพิ่ม dependencies ตามต้องการ

#### 5.4 **ไฟล์ที่ไม่ต้องแก้**

- ❌ `internal/models/` - PM ทำให้แล้ว
- ❌ `README.md` - มีให้แล้ว (แต่สามารถเพิ่มเติมได้)

### 6. รันโปรเจกต์

```bash
go run cmd/api/main.go
```

ตอนเริ่มต้นจะทำการ migrate ตารางอัตโนมัติ และรันที่ `http://localhost:8081`

### 6.1 แชร์ service ให้ทีม Kong ผ่าน ngrok

> ใช้เมื่อต้องการให้ทีมอื่นเรียก Users Service ผ่าน Kong ขณะพัฒนา

1. ปรับไฟล์ `.env` ให้ `APPLICATION_PORT=8080` (หรือรันด้วย `PORT=8080 go run ...`) เพื่อให้ตรงกับช่องทางมาตรฐาน
2. รัน service ตามปกติให้ฟังอยู่ที่พอร์ต 8080
3. เปิดเทอร์มินัลใหม่แล้วสั่ง
   ```bash
   ngrok http 8080
   ```
4. คัดลอก URL จาก `Forwarding` (เช่น `https://<hash>.ngrok-free.app`) และส่งให้เพื่อนที่ดูแล Kong Gateway เพื่อตั้งค่า Service/Route ให้ชี้มาที่ปลายทางของคุณ
5. ทดสอบ URL โดยตรงก่อนส่ง (เช่น `curl https://<hash>.ngrok-free.app/healthz`) แล้วค่อยตรวจสอบอีกครั้งผ่าน Kong หลังทีม Gateway ผูกเสร็จ

### 7. 📋 Checklist

- [ ] เขียน `cmd/api/main.go`
- [ ] เขียน handlers ใน `internal/handlers/` (ทำตาม TODO comments)
- [ ] เขียน repositories ใน `internal/repositories/` (ทำตาม TODO comments)
- [ ] เขียน services ใน `internal/services/` (ทำตาม TODO comments)
- [ ] ตั้งค่า `.env`
- [ ] ทดสอบ API ด้วย Postman
- [ ] ทดสอบผ่าน Kong Gateway

### 8. 🚀 การเอาขึ้น Github (Git Workflow)

#### 8.1 Clone Repository

**ขั้นตอนที่ 1: Clone Repository**

```bash
git clone https://github.com/Wattanaroj2567/users-service.git
```

**ผลลัพธ์ที่คาดหวัง:** Repository จะถูกดาวน์โหลดมาในโฟลเดอร์ `users-service`

**ขั้นตอนที่ 2: เข้าไปในโฟลเดอร์**

```bash
cd users-service
```

**ผลลัพธ์ที่คาดหวัง:** เปลี่ยน directory ไปยัง `users-service`

**ขั้นตอนที่ 3: ตรวจสอบ branch ปัจจุบัน**

```bash
git branch
```

**ผลลัพธ์ที่คาดหวัง:** ควรเห็น `* develop` (develop branch เป็นค่าเริ่มต้น)

#### 8.2 Development & Testing

**ขั้นตอนที่ 4: ตรวจสอบสถานะไฟล์**

```bash
git status
```

**ผลลัพธ์ที่คาดหวัง:** แสดงไฟล์ที่แก้ไข (modified files) และไฟล์ใหม่ (untracked files)

**ขั้นตอนที่ 5: เพิ่มไฟล์ที่แก้ไข**

```bash
git add .
```

**ผลลัพธ์ที่คาดหวัง:** ไฟล์ทั้งหมดถูกเพิ่มเข้า staging area

**ขั้นตอนที่ 6: Commit การเปลี่ยนแปลง**

```bash
git commit -m "feat: implement user authentication and profile management"
```

**ผลลัพธ์ที่คาดหวัง:** แสดงจำนวนไฟล์ที่เปลี่ยนแปลงและ commit hash

#### 8.3 Push to Develop Branch

**ขั้นตอนที่ 7: Push ไปยัง develop branch**

```bash
git push origin develop
```

**ผลลัพธ์ที่คาดหวัง:** การ push สำเร็จและแสดง URL ของ repository

#### 8.4 Final Merge to Main (PM ทำ)

```bash
# PM จะ merge develop ไป main เมื่อทุกอย่างสมบูรณ์
git checkout main
git merge develop
git push origin main
```

#### 8.5 Branch Strategy

| Branch    | Purpose                | Who       |
| --------- | ---------------------- | --------- |
| `develop` | การพัฒนาหลัก (Default) | Developer |
| `main`    | Production Ready       | PM        |

> 💡 **คำแนะนำ**: เปิดไฟล์โค้ดและดู **TODO comments** ที่มีอยู่แล้ว จะบอกว่าต้องทำอะไรบ้างในแต่ละส่วน

---
## 📝 API Documentation

### 2.1 ระบบล็อกอินและลงทะเบียน (Authentication)

| Endpoint            | Method | Auth Required    | Body / Parameters                                                                                                             | Format    |
| ------------------- | ------ | ---------------- | ----------------------------------------------------------------------------------------------------------------------------- | --------- |
| `/api/auth/register`| POST   | No               | - username: ชื่อผู้ใช้ (ต้องไม่ซ้ำ)<br>- display_name: ชื่อเล่นที่จะแสดง<br>- email: ต้องเป็นรูปแบบอีเมล<br>- password, confirm_password | JSON Body |
| `/api/auth/login`   | POST   | No               | - identifier: อีเมลหรือชื่อผู้ใช้สำหรับเข้าสู่ระบบ<br>- password: รหัสผ่าน                                                  | JSON Body |
| `/api/auth/forgot-password` | POST | No        | - email: ที่ลงทะเบียนไว้ เพื่อรับลิงก์รีเซ็ตรหัสผ่าน                                                                           | JSON Body |
| `/api/auth/reset-password`  | POST | No        | - token: โทเคนจากอีเมล<br>- new_password: รหัสผ่านใหม่<br>- confirm_password: ยืนยันรหัสผ่านใหม่                              | JSON Body |
| `/api/auth/logout`  | POST   | Yes (Member JWT) | - Authorization: Bearer {JWT}                                                                                                  | HTTP Header |

### 2.2 การจัดการข้อมูลบัญชีผู้ใช้ (User Profile)

| Endpoint            | Method | Auth Required    | Body / Parameters                                                                                                                                                                                       | Format                |
| ------------------- | ------ | ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------- |
| `/api/user/profile` | GET    | Yes (Member JWT) | - Authorization: Bearer {JWT}                                                                                                                                                                           | HTTP Header           |
| `/api/user/profile` | PUT    | Yes (Member JWT) | - username, display_name, email, profile_image (เมื่ออัปเดตข้อมูลทั่วไป)<br>- old_password, new_password, confirm_password (เมื่อเปลี่ยนรหัสผ่าน)<br>- delete_account_flag (true/false) + password (เมื่อปิดบัญชี) | JSON Body + HTTP Header |

### 2.3 JSON Request Examples

> 🧪 **สำหรับการทดสอบ API**: ใช้ JSON examples ด้านล่างเพื่อทดสอบ API ตาม use case หลัก  
> 📱 **แนะนำใช้ Postman**: จะช่วยจัดการ Headers และ Body ได้สะดวก

#### 2.3.1 Authentication Requests

**POST http://localhost:8081/api/auth/register**

ใช้สำหรับสร้างบัญชีผู้ใช้ใหม่ พร้อมกำหนดข้อมูลเริ่มต้น

> วางใน Body => Raw (JSON)

```json
{
  "username": "tawan123",
  "display_name": "Tawan Gamer",
  "email": "tawan@example.com",
  "password": "password123",
  "confirm_password": "password123"
}
```

**POST http://localhost:8081/api/auth/login**

ใช้สำหรับเข้าสู่ระบบ โดยรับได้ทั้งอีเมลหรือชื่อผู้ใช้

> วางใน Body => Raw (JSON)

```json
{
  "identifier": "tawan@example.com",
  "password": "password123"
}
```

> ระบบรองรับทั้งอีเมลหรือชื่อผู้ใช้ในฟิลด์ `identifier`

**POST http://localhost:8081/api/auth/forgot-password**

ใช้สำหรับส่งคำขอรับอีเมลรีเซ็ตรหัสผ่าน

> วางใน Body => Raw (JSON)

```json
{
  "email": "tawan@example.com"
}
```

**POST http://localhost:8081/api/auth/reset-password**

ใช้สำหรับตั้งรหัสผ่านใหม่ด้วยโทเคนที่ได้รับจากอีเมล

> วางใน Body => Raw (JSON)

```json
{
  "token": "RESET_TOKEN_HERE",
  "new_password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**POST http://localhost:8081/api/auth/logout**

ใช้สำหรับออกจากระบบและยกเลิกโทเคนปัจจุบัน

> วางใน Headers => Key: Authorization, Value: Bearer YOUR_JWT_TOKEN

```text
Authorization: Bearer YOUR_JWT_TOKEN
```

#### 2.3.2 User Profile Requests

> ทุกคำขอจำเป็นต้องแนบ `Authorization: Bearer YOUR_JWT_TOKEN`

**GET http://localhost:8081/api/user/profile**

ใช้สำหรับดึงรายละเอียดโปรไฟล์ของผู้ใช้ที่ล็อกอิน

> วางใน Headers => Key: Authorization, Value: Bearer YOUR_JWT_TOKEN

```text
Authorization: Bearer YOUR_JWT_TOKEN
```

**PUT http://localhost:8081/api/user/profile**

ใช้สำหรับแก้ไขข้อมูลโปรไฟล์ เช่น อีเมล ชื่อเล่น และรูปภาพ

> วางใน Headers => Key: Authorization, Value: Bearer YOUR_JWT_TOKEN  
> วางใน Body => Raw (JSON)

```json
{
  "username": "new_tawan123",
  "display_name": "Tawan Pro Gamer",
  "email": "new_tawan@example.com",
  "profile_image": "https://cdn.example.com/u/newpic.png"
}
```

**PUT http://localhost:8081/api/user/profile** 

ใช้สำหรับเปลี่ยนรหัสผ่านเดิมเป็นรหัสผ่านใหม่

> วางใน Headers => Key: Authorization, Value: Bearer YOUR_JWT_TOKEN  
> วางใน Body => Raw (JSON)

```json
{
  "old_password": "password123",
  "new_password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**PUT http://localhost:8081/api/user/profile**

ใช้สำหรับปิดบัญชีผู้ใช้ โดยต้องยืนยันด้วยรหัสผ่านปัจจุบัน

> วางใน Headers => Key: Authorization, Value: Bearer YOUR_JWT_TOKEN  
> วางใน Body => Raw (JSON)

```json
{
  "delete_account_flag": true,
  "password": "password123"
}
```

### 2.4 Error Responses

**ตัวอย่าง Error Response:**


```json
{
  "status": "error",
  "message": "Invalid credentials",
  "error": "Email or password is incorrect"
}
```

**HTTP Status Codes:**

- `200` - Success
- `400` - Bad Request (ข้อมูลไม่ถูกต้อง)
- `401` - Unauthorized (ไม่มีสิทธิ์เข้าถึง)
- `404` - Not Found (ไม่พบข้อมูล)
- `500` - Internal Server Error (ข้อผิดพลาดของเซิร์ฟเวอร์)

> 🧪 **วิธีทดสอบ**: ใช้ Postman เพื่อส่ง JSON requests ตาม examples ด้านบน และตรวจสอบ response ที่ได้รับ

---

## 📞 ติดต่อและสนับสนุน

### 👥 ทีมพัฒนา

| Role                | Name              | Contact                                     |
| ------------------- | ----------------- | ------------------------------------------- |
| **Developer**       | ณิชพน มานิตย์     | [GitHub](https://github.com/nitchapon66)    |
| **Project Manager** | วรรธนโรจน์ บุตรดี | [GitHub](https://github.com/Wattanaroj2567) |

### 🔗 ลิงก์ที่เกี่ยวข้อง

- **Main Project**: [Mini-Project-Golang](https://github.com/Wattanaroj2567/Mini-Project-Golang)
- **Kong Gateway Setup**: [Main README - Kong Setup](https://github.com/Wattanaroj2567/Mini-Project-Golang#-quick-start-ติดตั้งและรัน-kong--konga)
- **System Architecture**: [Main README - Architecture](https://github.com/Wattanaroj2567/Mini-Project-Golang#%EF%B8%8F-ภาพรวมระบบ-system-overview)

### 📚 เอกสารเพิ่มเติม

- **Troubleshooting**: [Main README - Troubleshooting](https://github.com/Wattanaroj2567/Mini-Project-Golang#-แก้ไขปัญหา-troubleshooting)
- **Ports Summary**: [Main README - Ports](https://github.com/Wattanaroj2567/Mini-Project-Golang#-รายการ-ports)

---

## ✅ สรุป

- README นี้อัปเดตให้สอดคล้องกับ **แนวทางหลักของโปรเจกต์ Mini-Project-Golang**
- รองรับทั้งการพัฒนา, ทดสอบ และเชื่อมต่อกับ Kong Gateway
- มีวิธีรันแบบ local, Docker, Kong Integration และ Remote Dev พร้อมใช้งานจริง
