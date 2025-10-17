# 🛡️ Admin Service — GameGear E-commerce

บริการสำหรับจัดการระบบหลังบ้าน (Admin Panel) ของโปรเจกต์ **GameGear E-commerce** เป็นหนึ่งใน microservices ที่ทำงานร่วมกับ **Kong API Gateway**

> 📖 **ดูเอกสารหลักของระบบ**: สำหรับ Kong Gateway setup, Architecture overview และการ integrate ทั้งระบบ → [Main README](https://github.com/Wattanaroj2567/Mini-Project-Golang.git)

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

`admin-service` ถูกออกแบบให้เป็น **"ผู้ประสานงาน" (Coordinator)** หรือ **"Admin Gateway"** โดย **ไม่มีฐานข้อมูลเป็นของตัวเอง** หน้าที่หลักคือเชื่อมต่อและประสานงานกับ Service อื่น ๆ ผ่าน API ของพวกมัน โดยไม่ยุ่งกับข้อมูลโดยตรง

**ตัวอย่าง:**

- เพิ่มสินค้า → เรียกใช้ API ของ **shop-service**
- จัดการผู้ใช้ → เรียกใช้ API ของ **users-service**

**Models ใน Admin Service:**

- **ไม่ใช่ Database Models** - เพราะไม่มีฐานข้อมูล
- **เป็น API Contract Models** - กำหนดโครงสร้างข้อมูลที่รับ-ส่งกับ frontend
- **Data Transformation** - แปลงข้อมูลจาก services อื่นให้เหมาะกับ admin panel
- **Response Standardization** - ให้ API response มีรูปแบบสม่ำเสมอ

แนวทางนี้ช่วยให้ **การเป็นเจ้าของข้อมูลยังคงอยู่กับ Service ที่รับผิดชอบ** และแยก Logic ฝั่ง Admin ออกมาเฉพาะ

---

## ✨ คุณสมบัติและ Endpoints

| Feature                | Method | Path                           | Auth Required        | Description                                             |
| ---------------------- | ------ | ------------------------------ | -------------------- | ------------------------------------------------------- |
| Admin Registration     | POST   | `/api/admin/register`          | No                   | สร้างบัญชี Admin ใหม่                                  |
| Admin Login            | POST   | `/api/admin/login`             | No                   | เข้าสู่ระบบ Admin Panel                                |
| Admin Forgot Password  | POST   | `/api/admin/forgot-password`   | No                   | ขออีเมลสำหรับรีเซ็ตรหัสผ่าน                           |
| Admin Reset Password   | POST   | `/api/admin/reset-password`    | No                   | ตั้งรหัสผ่านใหม่ด้วยโทเคน                              |
| Admin Logout           | POST   | `/api/admin/logout`            | Yes (Admin JWT)      | ออกจากระบบและเพิกถอนโทเคนปัจจุบัน                    |
| List Products          | GET    | `/api/admin/products`          | Yes (Admin JWT)      | ดึงข้อมูลสินค้าทั้งหมดผ่าน shop-service               |
| Create Product         | POST   | `/api/admin/products`          | Yes (Admin JWT)      | เพิ่มสินค้าใหม่ด้วยข้อมูลที่จำเป็น                     |
| Update Product         | PUT    | `/api/admin/products/:id`      | Yes (Admin JWT)      | ปรับรายละเอียดสินค้าที่มีอยู่                          |
| Delete Product         | DELETE | `/api/admin/products/:id`      | Yes (Admin JWT)      | ลบสินค้าออกจากระบบ                                     |
| List Orders            | GET    | `/api/admin/orders`            | Yes (Admin JWT)      | ดูคำสั่งซื้อทั้งหมดของระบบ                             |
| Update Order Status    | PUT    | `/api/admin/orders/:id/status` | Yes (Admin JWT)      | อัปเดตสถานะคำสั่งซื้อ (pending → delivered ฯลฯ)        |

> 🔐 **Authentication ทั้งหมด proxy ผ่าน users-service**: `admin-service` เรียก `users-service` ผ่าน Kong (`/users/api/admin/...`) เพื่อจัดการ register/login/forgot/reset โดยไม่เก็บรหัสผ่านเอง  
> 🛍️ **การจัดการสินค้า/คำสั่งซื้อ**: ใช้ HTTP client เชื่อมต่อ `shop-service` ผ่าน Kong (`/shop/api/...`) และทุกคำขอต้องแนบ JWT ที่มี `role=admin`

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
│   ├── clients/
│   │   ├── auth_client.go
│   │   └── shop_service_client.go
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── order_handler.go
│   │   ├── product_handler.go
│   │   └── routes.go
│   ├── models/
│   │   ├── auth.go
│   │   ├── order.go
│   │   └── product.go
│   └── services/
│       ├── auth_service.go
│       ├── order_service.go
│       ├── product_service.go
│       └── errors.go
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

คำอธิบายไฟล์:

- `cmd/api/main.go` — บูต service, ประกาศ routes, สร้าง clients/services/handlers ตาม configuration
- `internal/clients/auth_client.go` — HTTP client สำหรับเรียก endpoint auth ของ upstream (register/login/forgot/reset)
- `internal/clients/shop_service_client.go` — HTTP client สำหรับบริหารสินค้าและคำสั่งซื้อผ่าน shop-service
- `internal/handlers/auth_handler.go` — Handler สำหรับ `/api/admin/*` ที่เกี่ยวกับการยืนยันตัวตน
- `internal/handlers/product_handler.go` — Handler สำหรับ `/api/admin/products` (list/create/update/delete)
- `internal/handlers/order_handler.go` — Handler สำหรับ `/api/admin/orders` และการอัปเดตสถานะ
- `internal/handlers/routes.go` — รวมการ mount route ทั้งหมดลงใน Gin engine
- `internal/models/auth.go` — DTO สำหรับ request/response ด้าน auth ของ admin
- `internal/models/product.go` — DTO สำหรับจัดการสินค้า (summary และ upsert payload)
- `internal/models/order.go` — DTO สำหรับสรุป/รายละเอียดคำสั่งซื้อ และโครง request อัปเดตสถานะ
- `internal/services/auth_service.go` — โครง logic สำหรับประสานงานกับ `AuthClient`
- `internal/services/product_service.go` — โครง logic สำหรับประสานงานกับ shop-service ในการจัดการสินค้า
- `internal/services/order_service.go` — โครง logic สำหรับดึง/อัปเดตคำสั่งซื้อผ่าน shop-service
- `internal/services/errors.go` — ฟังก์ชัน/struct กลางสำหรับจัดการ service error ที่แชร์ให้ทุก service layer
- **.env.example** — ตัวอย่างไฟล์สำหรับตั้งค่าคอนฟิก เช่น Service URLs, CORS
- **README.md** — เอกสารอธิบายรายละเอียดการติดตั้งและใช้งาน Service

---

## 🚀 เริ่มต้นใช้งาน (Getting Started)

### 1. Clone & เข้าโฟลเดอร์

```bash
git clone https://github.com/Wattanaroj2567/admin-service.git
cd admin-service
```

### 2. ติดตั้ง dependencies

```bash
go mod tidy
```

### 3. สร้างไฟล์ `.env`

**สร้างไฟล์ `.env` ใหม่** โดยคัดลอกเนื้อหาจาก `.env.example` และแก้ไขค่าต่างๆ ตามต้องการ:

```env
# Application Configuration
APPLICATION_PORT=8083

# External Service URLs
SHOP_SERVICE_URL="http://localhost:8000/shop"
ADMIN_AUTH_SERVICE_URL="http://localhost:8000/users"

# CORS Configuration (optional หากมี reverse proxy จัดการอยู่แล้ว)
FRONTEND_URL="http://localhost:3000"
ALLOWED_ORIGINS="http://localhost:3000,http://localhost:8080"
```

**💡 หมายเหตุ:**

- **เรียกผ่าน Kong เป็นหลัก**: ค่าเริ่มต้นใช้ `http://localhost:8000/...` เพื่อให้ admin-service ผ่านกฎ Gateway เสมอ (routing, auth, rate limit)
- **กรณีอยู่ใน Docker Network**: เปลี่ยน host เป็น `http://kong:8000/...` เช่น `SHOP_SERVICE_URL="http://kong:8000/shop"`
- **ไฟล์ `.env.example`**: เก็บไว้เป็น template สำหรับตั้งค่าครั้งถัดไป

### 4. ไฟล์ที่ต้องแก้ไข/เขียนโค้ด

> 👨‍💻 **Developer & PM**: วรรธนโรจน์ บุตรดี

คุณต้องเขียนโค้ดเฉพาะใน **2 จุดหลัก**:

#### 4.1 **`cmd/api/main.go`**

- เริ่มต้น Gin server
- Setup routes
- Middleware configuration
- **ไม่ต้องต่อ database** (service นี้ไม่มี DB)

#### 4.2 **โฟลเดอร์ `internal/`**

| Folder           | ต้องทำ    | หมายเหตุ                                     |
| ---------------- | --------- | -------------------------------------------- |
| ✅ **handlers/** | ✅ ต้องทำ | เขียน HTTP handlers สำหรับ admin panel       |
| ✅ **services/** | ✅ ต้องทำ | เขียน business logic และ orchestration       |
| ✅ **clients/**  | ✅ ต้องทำ | เขียน HTTP clients สำหรับเรียก services อื่น |
| ✅ **models/**   | ✅ ต้องทำ | เขียน API Contract models                    |

> 💡 **หมายเหตุ**: ในโค้ดจะมี **TODO comments** บอกว่าต้องทำอะไรบ้าง ให้ทำตามที่ระบุไว้ในโค้ด

#### 4.3 **ไฟล์อื่นๆ**

- ✅ **`.env`** - ตั้งค่า environment variables
- ✅ **`go.mod`** - เพิ่ม dependencies ตามต้องการ

#### 4.4 **ไฟล์ที่ไม่ต้องแก้**

- ❌ `README.md` - มีให้แล้ว (แต่สามารถเพิ่มเติมได้)

### 5. รันโปรเจกต์

```bash
go run cmd/api/main.go
```

ตอนเริ่มต้นจะรันที่ `http://localhost:8083`

#### 5.1 แชร์ service ให้ทีม Kong ผ่าน ngrok

1. ปรับ `.env` ให้ `APPLICATION_PORT=8080` (หรือรันด้วย `PORT=8080`) เพื่อให้ตรงกับมาตรฐานที่ใช้เปิดอุโมงค์
2. รัน Admin Service ที่พอร์ต 8080 จากนั้นเปิดเทอร์มินัลใหม่และสั่ง
   ```bash
   ngrok http 8080
   ```
3. คัดลอก Forwarding URL ที่ได้ แล้วส่งให้เพื่อนที่ดูแล Kong Gateway เพื่อผูกกับ Service `/admin`
4. ทดสอบ URL โดยตรง (เช่น `curl https://<hash>.ngrok-free.app/healthz`) ก่อนตรวจสอบผ่าน Kong หลังจากทีม Gateway ตั้งค่าเสร็จ

### 6. 📋 Checklist

- [ ] เขียน `cmd/api/main.go`
- [ ] เขียน handlers ใน `internal/handlers/` (ทำตาม TODO comments)
- [ ] เขียน services ใน `internal/services/` (ทำตาม TODO comments)
- [ ] เขียน clients ใน `internal/clients/` (ทำตาม TODO comments)
- [ ] เขียน models ใน `internal/models/` (ทำตาม TODO comments)
- [ ] ตั้งค่า `.env`
- [ ] ทดสอบ API ด้วย Postman
- [ ] ทดสอบผ่าน Kong Gateway

### 7. 🚀 การเอาขึ้น Github (Git Workflow)

#### 7.1 Clone Repository

**ขั้นตอนที่ 1: Clone Repository**

```bash
git clone https://github.com/Wattanaroj2567/admin-service.git
```

**ผลลัพธ์ที่คาดหวัง:** Repository จะถูกดาวน์โหลดมาในโฟลเดอร์ `admin-service`

**ขั้นตอนที่ 2: เข้าไปในโฟลเดอร์**

```bash
cd admin-service
```

**ผลลัพธ์ที่คาดหวัง:** เปลี่ยน directory ไปยัง `admin-service`

**ขั้นตอนที่ 3: ตรวจสอบ branch ปัจจุบัน**

```bash
git branch
```

**ผลลัพธ์ที่คาดหวัง:** ควรเห็น `* develop` (develop branch เป็นค่าเริ่มต้น)

#### 7.2 Development & Testing

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
git commit -m "feat: implement admin panel with product, order, and user management"
```

**ผลลัพธ์ที่คาดหวัง:** แสดงจำนวนไฟล์ที่เปลี่ยนแปลงและ commit hash

#### 7.3 Push to Develop Branch

**ขั้นตอนที่ 7: Push ไปยัง develop branch**

```bash
git push origin develop
```

**ผลลัพธ์ที่คาดหวัง:** การ push สำเร็จและแสดง URL ของ repository

#### 7.4 Final Merge to Main (PM ทำ)

```bash
# PM จะ merge develop ไป main เมื่อทุกอย่างสมบูรณ์
git checkout main
git merge develop
git push origin main
```

#### 7.5 Branch Strategy

| Branch    | Purpose                | Who |
| --------- | ---------------------- | --- |
| `develop` | การพัฒนาหลัก (Default) | PM  |
| `main`    | Production Ready       | PM  |

---

## 📦 Module Structure

Service นี้ใช้ Go Module สำหรับจัดการ dependencies:

| Property              | Value                               |
| --------------------- | ----------------------------------- |
| **Module Name**       | `github.com/gamegear/admin-service` |
| **Go Version**        | 1.25.1                              |
| **Main Dependencies** | Gin, HTTP Client, JWT               |

### Local Development Setup

สำหรับการพัฒนาในเครื่อง local service นี้ใช้ `replace` directive ใน `go.mod`:

```go
// ใน shop-service/go.mod
replace github.com/gamegear/admin-service => ../admin-service

// ใน users-service/go.mod
replace github.com/gamegear/admin-service => ../admin-service
```

### Dependencies Management

- **Web Framework**: Gin (HTTP router)
- **HTTP Client**: สำหรับเรียก services อื่น
- **Authentication**: JWT tokens
- **API Documentation**: Postman
- **Environment**: godotenv for .env files

### Import Statements

เมื่อต้องการ import จาก services อื่น ให้ใช้ module name แทน local path:

```go
// ✅ ถูกต้อง - ใช้ module name
import "github.com/gamegear/shop-service/internal/models"

// ❌ ผิด - อย่าใช้ local path
import "../shop-service/internal/models"
```

**หมายเหตุ**: admin-service ต้อง import จาก shop-service และ users-service เพื่อใช้ข้อมูล

---

## 📝 API Documentation

### 2.1 ระบบ Authentication (Admin Authentication)

| Endpoint                     | Method | Auth Required    | Body / Parameters                                                                                          | Format      |
| ---------------------------- | ------ | ---------------- | ----------------------------------------------------------------------------------------------------------- | ----------- |
| `/api/admin/register`        | POST   | No               | - email: อีเมลของผู้ดูแลระบบ<br>- password: รหัสผ่าน<br>- confirm_password: ยืนยันรหัสผ่าน<br>- display_name: ชื่อที่แสดง | JSON Body   |
| `/api/admin/login`           | POST   | No               | - email: อีเมลของผู้ดูแลระบบ<br>- password: รหัสผ่าน                                                      | JSON Body   |
| `/api/admin/forgot-password` | POST   | No               | - email: อีเมลที่ลงทะเบียนไว้เพื่อรับลิงก์รีเซ็ตรหัสผ่าน                                                 | JSON Body   |
| `/api/admin/reset-password`  | POST   | No               | - token: โทเคนที่ได้รับจากอีเมล<br>- new_password: รหัสผ่านใหม่<br>- confirm_password: ยืนยันรหัสผ่านใหม่ | JSON Body   |
| `/api/admin/logout`          | POST   | Yes (Admin JWT)  | - Authorization: Bearer {JWT}                                                                               | HTTP Header |

### 2.2 ระบบจัดการสินค้า (Products Management)

| Endpoint                  | Method | Auth Required    | Body / Parameters                                                                                                                                      | Format                  |
| ------------------------- | ------ | ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------- |
| `/api/admin/products`     | GET    | Yes (Admin JWT)  | - Authorization: Bearer {JWT}                                                                                                                           | HTTP Header             |
| `/api/admin/products`     | POST   | Yes (Admin JWT)  | - name: ชื่อสินค้า<br>- description: คำอธิบาย<br>- price: ราคา<br>- stock: จำนวนสต็อก<br>- category_id: หมวดหมู่<br>- image_url: URL รูปสินค้า       | JSON Body + HTTP Header |
| `/api/admin/products/:id` | PUT    | Yes (Admin JWT)  | - id: รหัสสินค้าใน path<br>- name, description, price, stock, category_id, image_url (ข้อมูลที่ต้องการอัปเดต)                                         | JSON Body + HTTP Header |
| `/api/admin/products/:id` | DELETE | Yes (Admin JWT)  | - id: รหัสสินค้าใน path                                                                                                                                 | HTTP Header             |

### 2.3 ระบบจัดการคำสั่งซื้อ (Orders Management)

| Endpoint                        | Method | Auth Required    | Body / Parameters                                                                                                   | Format                  |
| ------------------------------- | ------ | ---------------- | -------------------------------------------------------------------------------------------------------------------- | ----------------------- |
| `/api/admin/orders`             | GET    | Yes (Admin JWT)  | - Authorization: Bearer {JWT}                                                                                       | HTTP Header             |
| `/api/admin/orders/:id/status`  | PUT    | Yes (Admin JWT)  | - id: รหัสคำสั่งซื้อใน path<br>- status: ค่าสถานะใหม่ (`pending`, `processing`, `shipped`, `delivered`, `cancelled`) | JSON Body + HTTP Header |

### 2.4 JSON Request Examples

> 🧪 **สำหรับการทดสอบ API**: ใช้ JSON examples ด้านล่างเพื่อทดสอบ flow หลัก  
> 🔐 **แนบ `Authorization: Bearer ADMIN_JWT_TOKEN` สำหรับทุกคำขอที่ต้องยืนยันตัวตน**

#### 2.4.1 Admin Authentication Requests

**POST http://localhost:8083/api/admin/register**

> วางใน Body => Raw (JSON)

ใช้สำหรับสร้างบัญชีผู้ดูแลระบบใหม่

```json
{
  "email": "new.admin@gamegear.com",
  "password": "securePassword123",
  "confirm_password": "securePassword123",
  "display_name": "Operations Lead"
}
```

**POST http://localhost:8083/api/admin/forgot-password**

> วางใน Body => Raw (JSON)

ใช้สำหรับส่งคำขอรีเซ็ตรหัสผ่านไปยังอีเมลที่ลงทะเบียนไว้

```json
{
  "email": "admin@gamegear.com"
}
```

**POST http://localhost:8083/api/admin/reset-password**

> วางใน Body => Raw (JSON)

ใช้สำหรับตั้งรหัสผ่านใหม่ด้วยโทเคนที่ได้รับจากอีเมล

```json
{
  "token": "RESET_TOKEN_FROM_EMAIL",
  "new_password": "NewStrongPassword123",
  "confirm_password": "NewStrongPassword123"
}
```

**POST http://localhost:8083/api/admin/login**

>วางใน Body => Raw ในช่อง Json

ใช้สำหรับเข้าสู่ระบบแผงผู้ดูแลและรับ JWT

```json
{
  "email": "admin@gamegear.com",
  "password": "admin123"
}
```

**POST http://localhost:8083/api/admin/logout**

>วางใน Headers => ตาราง Key, Value

ใช้สำหรับออกจากระบบและทำให้โทเคนหมดอายุ

```bash
Key: Authorization  Value: Bearer ADMIN_JWT_TOKEN
```

#### 2.4.2 Products Management Requests

**GET http://localhost:8083/api/admin/products**

>วางใน Headers => ตาราง Key, Value

ใช้สำหรับดึงรายการสินค้าทั้งหมดเพื่อแสดงในแผงผู้ดูแล

```bash
Key: Authorization  Value: Bearer ADMIN_JWT_TOKEN
```

**POST http://localhost:8083/api/admin/products**

>วางใน Body => Raw ในช่อง Json

ใช้สำหรับเพิ่มสินค้าใหม่เข้าสู่ระบบผ่าน shop-service

```json
{
  "name": "Gaming Keyboard RGB",
  "description": "คีย์บอร์ดเกมมิ่งพร้อมไฟ RGB สวยงาม",
  "price": 2599.0,
  "stock": 30,
  "category_id": 2,
  "image_url": "https://example.com/keyboard.jpg"
}
```

**PUT http://localhost:8083/api/admin/products/1**

>วางใน Body => Raw ในช่อง Json

ใช้สำหรับแก้ไขรายละเอียดสินค้าที่มีอยู่เดิม

```json
{
  "name": "Gaming Keyboard RGB Pro",
  "description": "คีย์บอร์ดเกมมิ่งพร้อมไฟ RGB สวยงาม รุ่น Pro",
  "price": 2999.0,
  "stock": 20,
  "category_id": 2,
  "image_url": "https://example.com/keyboard-pro.jpg"
}
```

**DELETE http://localhost:8083/api/admin/products/1**

>วางในช่อง URL และ Path นั้นๆ

ใช้สำหรับลบสินค้าออกจากระบบโดยอ้างอิงรหัสสินค้า

```bash
ระบุ Path หมายเลข id ด้านหลัง / ที่ได้รับจาก Response "เพิ่มสินค้า" ตามที่เราต้องการลบในแต่ละชิ้น เช่น id:1 ให้ระบุ /1
```

#### 2.4.3 Orders Management Requests

**GET http://localhost:8083/api/admin/orders**

>วางใน Headers => ตาราง Key, Value

ใช้สำหรับดึงรายการคำสั่งซื้อทั้งหมดในระบบ

```bash
Key: Authorization  Value: Bearer ADMIN_JWT_TOKEN
```

**PUT http://localhost:8083/api/admin/orders/1/status**

>วางใน Body => Raw ในช่อง Json

ใช้สำหรับอัปเดตสถานะคำสั่งซื้อที่เลือก

```json
{
  "status": "shipped"
}
```

> เปลี่ยนค่า `status` เป็น `processing`, `delivered` หรือ `cancelled` ได้ตามสถานะที่ต้องการอัปเดต

### 2.5 Error Responses

**ตัวอย่าง Error Response:**


```json
{
  "status": "error",
  "message": "Unauthorized access",
  "error": "Invalid JWT token"
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

| Role               | Name              | Contact                                     |
| ------------------ | ----------------- | ------------------------------------------- |
| **Developer & PM** | วรรธนโรจน์ บุตรดี | [GitHub](https://github.com/Wattanaroj2567) |

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
