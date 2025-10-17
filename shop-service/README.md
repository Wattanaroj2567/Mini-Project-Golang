# 🛍️ Shop Service — GameGear E-commerce

บริการสำหรับจัดการประสบการณ์การซื้อขายทั้งหมดของโปรเจกต์ **GameGear E-commerce** ตั้งแต่การแสดงสินค้า, การจัดการตะกร้า, ไปจนถึงการสร้างคำสั่งซื้อ เป็นหนึ่งใน microservices ที่ทำงานร่วมกับ **Kong API Gateway**

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

- เป็น **เจ้าของข้อมูล (Data Owner)** สำหรับตารางที่เกี่ยวข้องกับการซื้อขาย เช่น `products`, `carts`, `orders` เป็นต้น
- ข้อมูลของสินค้าและคำสั่งซื้อ **จะไม่ถูกเข้าถึงโดยตรงจาก service อื่น** — ต้องเรียกผ่าน API ของ Shop Service เท่านั้น
- บริการอื่น (Users, Admin) ที่ต้องการใช้ข้อมูลเหล่านี้ ต้องเรียกผ่าน Gateway เพื่อรักษาความถูกต้องและความปลอดภัยของข้อมูล

---

## ✨ คุณสมบัติและ Endpoints

### ลูกค้าทั่วไป (Member Endpoints)

| Feature              | Method | Path                      | Auth Required        | Description                                   |
| -------------------- | ------ | ------------------------- | -------------------- | --------------------------------------------- |
| Browse Products      | GET    | `/api/products`          | No                   | ดูรายการสินค้า พร้อมตัวกรองและแบ่งหน้า     |
| View Product Details | GET    | `/api/products/:id`      | No                   | ดูรายละเอียดสินค้าแต่ละชิ้น                  |
| View Cart            | GET    | `/api/cart`              | Yes (Member JWT)     | ดูสินค้าในตะกร้าของฉัน                       |
| Add To Cart          | POST   | `/api/cart/add`          | Yes (Member JWT)     | เพิ่มสินค้าลงในตะกร้า                       |
| Update Cart Item     | PUT    | `/api/cart/update`       | Yes (Member JWT)     | แก้ไขจำนวนสินค้าภายในตะกร้า                |
| Remove From Cart     | DELETE | `/api/cart/remove`       | Yes (Member JWT)     | ลบสินค้าจากตะกร้า                            |
| Checkout Order       | POST   | `/api/orders`            | Yes (Member JWT)     | สร้างคำสั่งซื้อจากตะกร้า                     |
| Order History        | GET    | `/api/orders/history`    | Yes (Member JWT)     | ดูประวัติคำสั่งซื้อของตนเอง                 |

### ผู้ดูแลระบบ (Admin Endpoints) — ใช้ผ่าน Kong คู่กับ JWT ที่มี `role=admin`

| Feature                 | Method | Path                       | Auth Required    | Description                                                          |
| ----------------------- | ------ | -------------------------- | ---------------- | -------------------------------------------------------------------- |
| Create Product          | POST   | `/api/products`           | Yes (Admin JWT)  | เพิ่มสินค้าใหม่ลงฐานข้อมูล                                          |
| Update Product          | PUT    | `/api/products/:id`       | Yes (Admin JWT)  | แก้ไขข้อมูลสินค้าเดิม                                                |
| Delete Product          | DELETE | `/api/products/:id`       | Yes (Admin JWT)  | ลบสินค้าออกจากระบบ                                                   |
| List All Orders         | GET    | `/api/orders`             | Yes (Admin JWT)  | ดูคำสั่งซื้อทั้งหมดในระบบ (สำหรับแดชบอร์ดผู้ดูแล)                  |
| Update Order Status     | PUT    | `/api/orders/:id/status`  | Yes (Admin JWT)  | เปลี่ยนสถานะคำสั่งซื้อ (`pending`, `processing`, `shipped`, ฯลฯ)    |

> 🔐 **หมายเหตุ**: admin-service จะเรียก endpoint เหล่านี้ผ่าน Kong เช่น `http://localhost:8000/shop/api/products` โดยต้องแนบ JWT จาก users-service ที่มี `role=admin` เสมอ

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
│   │   ├── cart_handler.go
│   │   ├── order_handler.go
│   │   ├── product_handler.go
│   │   └── routes.go
│   ├── middleware/
│   │   └── auth.go
│   ├── services/
│   │   ├── cart_service.go
│   │   ├── order_service.go
│   │   └── product_service.go
│   ├── repositories/
│   │   ├── cart_item_repository.go
│   │   ├── cart_repository.go
│   │   ├── order_repository.go
│   │   └── product_repository.go
│   ├── security/
│   │   └── token_validator.go
│   └── models/
│       ├── cart_model.go
│       ├── cart_payloads.go
│       ├── order_model.go
│       ├── order_payloads.go
│       ├── order_admin_payloads.go
│       ├── product_model.go
│       └── product_payloads.go
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

คำอธิบายไฟล์:

- `cmd/api/main.go` — บูต service, เชื่อม DB, migrate model, อ่าน `JWT_SECRET_KEY` และประกาศ route พร้อม middleware
- `internal/handlers/product_handler.go` — รับคำขอ `/api/products` ทั้ง read (สาธารณะ) และ admin CRUD
- `internal/handlers/cart_handler.go` — จัดการ `/api/cart`, `/api/cart/add`, `/api/cart/update`, `/api/cart/remove`
- `internal/handlers/order_handler.go` — จัดการ `/api/orders` สำหรับสมาชิก และ `/api/orders`, `/api/orders/:id/status` สำหรับแอดมิน
- `internal/handlers/routes.go` — รวมการประกาศ route ทั้งหมด พร้อม inject middleware แยก member/admin
- `internal/middleware/auth.go` — Gin middleware สำหรับตรวจสอบ JWT และ role (`member`, `admin`)
- `internal/services/product_service.go` — โครง logic สำหรับค้นหา/จัดการสินค้า (list + admin create/update/delete)
- `internal/services/cart_service.go` — โครง logic สำหรับจัดการตะกร้าสินค้า
- `internal/services/order_service.go` — โครง logic สำหรับสร้างคำสั่งซื้อ, ดึงประวัติ, และ admin list/update status
- `internal/repositories/product_repository.go` — ทำงานกับตาราง `products`
- `internal/repositories/cart_repository.go` — จัดการตาราง `carts`
- `internal/repositories/cart_item_repository.go` — จัดการตาราง `cart_items`
- `internal/repositories/order_repository.go` — ทำงานกับตาราง `orders` และ `order_items` รวมถึง update status
- `internal/security/token_validator.go` — Utility สำหรับตรวจสอบความถูกต้องของ JWT ด้วย secret
- `internal/models/product_model.go` — GORM model ของสินค้า (ฟิลด์สอดคล้อง README)
- `internal/models/cart_model.go` — GORM model สำหรับตะกร้าและรายการสินค้าในตะกร้า
- `internal/models/order_model.go` — GORM model สำหรับคำสั่งซื้อและรายการภายในคำสั่งซื้อ
- `internal/models/product_payloads.go` — DTO สำหรับ filter, create, update, response ของสินค้า
- `internal/models/cart_payloads.go` — DTO สำหรับการเพิ่ม/แก้ไข/ลบสินค้าและการตอบกลับตะกร้า
- `internal/models/order_payloads.go` — DTO สำหรับสร้างคำสั่งซื้อ, รายการประวัติ, และรายละเอียดคำสั่งซื้อ
- `internal/models/order_admin_payloads.go` — DTO สำหรับ admin list orders และคำขออัปเดตสถานะ
- **.env.example** — ตัวอย่างไฟล์สำหรับตั้งค่าคอนฟิก เช่น Database URL, JWT, Email
- **README.md** — เอกสารอธิบายรายละเอียดการติดตั้งและใช้งาน Service

---

## 🚀 เริ่มต้นใช้งาน (Getting Started)

### 1. Clone & เข้าโฟลเดอร์

```bash
git clone https://github.com/Wattanaroj2567/shop-service.git
cd shop-service
```

### 2. ติดตั้ง dependencies

```bash
go mod tidy
```

### 3. ตั้งค่า Database

ตรวจสอบว่า PostgreSQL ทำงานอยู่ แล้วสร้างฐานข้อมูล:

```sql
CREATE DATABASE gamegear_shop_db;
```

### 4. สร้างไฟล์ `.env`

**สร้างไฟล์ `.env` ใหม่** โดยคัดลอกเนื้อหาจาก `.env.example` และแก้ไขค่าต่างๆ ตามต้องการ:

```env
# Application Configuration
APPLICATION_PORT=8082
APP_NAME="GameGear Shop Service"
APP_VERSION="1.0.0"

# Database Configuration
DATABASE_URL="host=localhost user=postgres password=your_password dbname=gamegear_shop_db port=5432 sslmode=disable"
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=gamegear_shop_db

# JWT Configuration
JWT_SECRET_KEY="your_super_secret_jwt_key_here_make_it_long_and_secure"
JWT_EXPIRATION_HOURS=24

# CORS Configuration
FRONTEND_URL="http://localhost:3000"
ALLOWED_ORIGINS="http://localhost:3000,http://localhost:8080"

# External Service URLs
USER_SERVICE_URL="http://localhost:8081"
ADMIN_SERVICE_URL="http://localhost:8083"


# Inventory Configuration
LOW_STOCK_THRESHOLD=10
MAX_CART_ITEMS=50

# Logging Configuration
LOG_LEVEL=info
LOG_FILE="logs/shop-service.log"
```

**💡 หมายเหตุ:**

- **JWT_SECRET_KEY**: ต้องใช้ค่าเดียวกับ users-service เพื่อให้ middleware ตรวจสอบ token ข้าม service ได้
- **Database Password**: ใช้รหัสผ่าน PostgreSQL ของคุณ
- **Payment Method**: เป็นการจำลอง (Mock) - ไม่ต้องเชื่อมต่อ Payment Gateway จริง
- **ไฟล์ `.env.example`**: เก็บไว้เป็น template สำหรับครั้งต่อไป

### 5. ไฟล์ที่ต้องแก้ไข/เขียนโค้ด

> 👨‍💻 **Developers**:
>
> - **ณัฐพงษ์ ดีบุตร** - Product Catalog (branch: `feature/product-catalog`)
> - **วายุ กอคูณ** - Cart & Order (branch: `feature/cart-and-order`)

คุณต้องเขียนโค้ดเฉพาะใน **2 จุดหลัก**:

#### 5.1 **`cmd/api/main.go`**

- เริ่มต้น Gin server
- Setup routes
- Connect database
- Middleware configuration

#### 5.2 **โฟลเดอร์ `internal/`**

| Folder               | ต้องทำ       | หมายเหตุ                                        |
| -------------------- | ------------ | ----------------------------------------------- |
| ✅ **handlers/**     | ✅ ต้องทำ    | เขียน HTTP handlers สำหรับ products/cart/orders |
| ✅ **middleware/**   | ✅ ต้องทำ    | เขียน/ปรับ Middleware สำหรับตรวจสอบ JWT & role |
| ❌ **models/**       | ❌ ไม่ต้องทำ | **PM ทำให้แล้ว** (วรรธนโรจน์)                   |
| ✅ **repositories/** | ✅ ต้องทำ    | เขียน database operations (CRUD)                |
| ✅ **services/**     | ✅ ต้องทำ    | เขียน business logic                            |
| ❌ **security/**     | ❌ ไม่ต้องทำ | PM เตรียม `token_validator.go` สำหรับ reuse     |

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

ตอนเริ่มต้นจะทำการ migrate ตารางอัตโนมัติ และรันที่ `http://localhost:8082`

### 6.1 แชร์ service ให้ทีม Kong ผ่าน ngrok

1. หากต้องการให้ทีมอื่นเข้าถึง Shop Service ผ่าน Kong ระหว่างพัฒนา ให้ตั้ง `APPLICATION_PORT=8080` ใน `.env` หรือรันด้วย `PORT=8080`
2. รัน service ให้ฟังบนพอร์ต 8080 แล้วเปิดเทอร์มินัลใหม่
3. สั่ง
   ```bash
   ngrok http 8080
   ```
4. คัดลอก Forwarding URL (เช่น `https://<hash>.ngrok-free.app`) ส่งให้ผู้ดูแล Kong Gateway เพื่ออัปเดต Service/Route `/shop`
5. ตรวจสอบ URL โดยตรง (ตัวอย่าง `curl https://<hash>.ngrok-free.app/healthz`) ก่อนทดสอบผ่าน Kong หลังมีการเชื่อมต่อ

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
git clone https://github.com/Wattanaroj2567/shop-service.git
```

**ผลลัพธ์ที่คาดหวัง:** Repository จะถูกดาวน์โหลดมาในโฟลเดอร์ `shop-service`

**ขั้นตอนที่ 2: เข้าไปในโฟลเดอร์**

```bash
cd shop-service
```

**ผลลัพธ์ที่คาดหวัง:** เปลี่ยน directory ไปยัง `shop-service`

**ขั้นตอนที่ 3: ตรวจสอบ branch ปัจจุบัน**

```bash
git branch
```

**ผลลัพธ์ที่คาดหวัง:** ควรเห็น `* develop` (develop branch เป็นค่าเริ่มต้น)

#### 8.2 Checkout Feature Branch

**สำหรับ ณัฐพงษ์ (Product Catalog):**

**ขั้นตอนที่ 4: Checkout feature branch สำหรับ Product Catalog**

```bash
git checkout -b feature/product-catalog
```

**ผลลัพธ์ที่คาดหวัง:** เปลี่ยนไปยัง `feature/product-catalog` branch และแสดง "Switched to a new branch"

**สำหรับ วายุ (Cart & Order):**

**ขั้นตอนที่ 4: Checkout feature branch สำหรับ Cart & Order**

```bash
git checkout -b feature/cart-and-order
```

**ผลลัพธ์ที่คาดหวัง:** เปลี่ยนไปยัง `feature/cart-and-order` branch และแสดง "Switched to a new branch"

#### 8.3 Development & Testing

**ขั้นตอนที่ 5: ตรวจสอบสถานะไฟล์**

```bash
git status
```

**ผลลัพธ์ที่คาดหวัง:** แสดงไฟล์ที่แก้ไข (modified files) และไฟล์ใหม่ (untracked files)

**ขั้นตอนที่ 6: เพิ่มไฟล์ที่แก้ไข**

```bash
git add .
```

**ผลลัพธ์ที่คาดหวัง:** ไฟล์ทั้งหมดถูกเพิ่มเข้า staging area

**ขั้นตอนที่ 7: Commit การเปลี่ยนแปลง**

**สำหรับ ณัฐพงษ์ (Product Catalog):**

```bash
git commit -m "feat: implement product catalog management"
```

**สำหรับ วายุ (Cart & Order):**

```bash
git commit -m "feat: implement cart and order management"
```

**ผลลัพธ์ที่คาดหวัง:** แสดงจำนวนไฟล์ที่เปลี่ยนแปลงและ commit hash

#### 8.4 Push to Feature Branch

**ขั้นตอนที่ 8: Push ไปยัง feature branch ของตัวเอง**

**สำหรับ ณัฐพงษ์ (Product Catalog):**

```bash
git push origin feature/product-catalog
```

**สำหรับ วายุ (Cart & Order):**

```bash
git push origin feature/cart-and-order
```

**ผลลัพธ์ที่คาดหวัง:** การ push สำเร็จและแสดง URL ของ repository

#### 8.5 Final Merge (PM ทำ)

```bash
# PM จะ merge feature branches ไป develop
git checkout develop
git merge feature/product-catalog
git merge feature/cart-and-order
git push origin develop

# สุดท้าย merge develop ไป main
git checkout main
git merge develop
git push origin main
```

#### 8.6 Branch Strategy

| Branch                    | Purpose                 | Who     |
| ------------------------- | ----------------------- | ------- |
| `develop`                 | การพัฒนาหลัก (Default)  | PM      |
| `feature/product-catalog` | Product Management      | ณัฐพงษ์ |
| `feature/cart-and-order`  | Cart & Order Management | วายุ    |
| `main`                    | Production Ready        | PM      |

---
## 📦 Module Structure

Service นี้ใช้ Go Module สำหรับจัดการ dependencies:

| Property              | Value                              |
| --------------------- | ---------------------------------- |
| **Module Name**       | `github.com/gamegear/shop-service` |
| **Go Version**        | 1.25.1                             |
| **Main Dependencies** | Gin, GORM, PostgreSQL, JWT         |

### Local Development Setup

สำหรับการพัฒนาในเครื่อง local service นี้ใช้ `replace` directive ใน `go.mod`:

```go
// ใน admin-service/go.mod
replace github.com/gamegear/shop-service => ../shop-service

// ใน users-service/go.mod
replace github.com/gamegear/shop-service => ../shop-service
```

### Dependencies Management

- **Web Framework**: Gin (HTTP router)
- **Database**: GORM + PostgreSQL driver
- **Authentication**: JWT tokens
- **API Documentation**: Postman
- **Environment**: godotenv for .env files

### Import Statements

เมื่อต้องการ import จาก services อื่น ให้ใช้ module name แทน local path:

```go
// ✅ ถูกต้อง - ใช้ module name
import "github.com/gamegear/users-service/internal/models"

// ❌ ผิด - อย่าใช้ local path
import "../users-service/internal/models"
```

**หมายเหตุ**: shop-service อาจต้อง import จาก users-service เพื่อใช้ข้อมูลผู้ใช้

---

## 📝 API Documentation

### 2.1 ระบบสินค้า (Products)

| Endpoint            | Method | Auth Required | Body / Parameters                                                                                  | Format           |
| ------------------- | ------ | ------------- | --------------------------------------------------------------------------------------------------- | ---------------- |
| `/api/products`     | GET    | No            | - page: หน้าที่ต้องการ (optional)<br>- limit: จำนวนต่อหน้า (optional)<br>- category, search (optional) | Query Parameters |
| `/api/products/:id` | GET    | No            | - id: รหัสสินค้าที่ต้องการดู                                                                       | Path Parameter   |

### 2.2 ระบบตะกร้าสินค้า (Cart)

| Endpoint           | Method | Auth Required    | Body / Parameters                                                | Format                  |
| ------------------ | ------ | ---------------- | ----------------------------------------------------------------- | ----------------------- |
| `/api/cart`        | GET    | Yes (Member JWT) | - Authorization: Bearer {JWT}                                    | HTTP Header             |
| `/api/cart/add`    | POST   | Yes (Member JWT) | - product_id: รหัสสินค้า<br>- quantity: จำนวนที่ต้องการ          | JSON Body + HTTP Header |
| `/api/cart/update` | PUT    | Yes (Member JWT) | - cart_item_id: รหัสรายการในตะกร้า<br>- quantity: จำนวนใหม่     | JSON Body + HTTP Header |
| `/api/cart/remove` | DELETE | Yes (Member JWT) | - cart_item_id: รหัสรายการในตะกร้า                              | JSON Body + HTTP Header |

### 2.3 ระบบคำสั่งซื้อ (Orders)

| Endpoint              | Method | Auth Required    | Body / Parameters                                                                                  | Format                  |
| --------------------- | ------ | ---------------- | --------------------------------------------------------------------------------------------------- | ----------------------- |
| `/api/orders`         | POST   | Yes (Member JWT) | - cart_id: รหัสตะกร้าที่ต้องการเช็คเอาท์<br>- shipping_address: ที่อยู่จัดส่ง<br>- payment_method: วิธีการชำระเงิน (mock)<br>- notes: หมายเหตุเพิ่มเติม (optional) | JSON Body + HTTP Header |
| `/api/orders/history` | GET    | Yes (Member JWT) | - Authorization: Bearer {JWT}                                                                      | HTTP Header             |

### 2.4 JSON Request Examples

> 🧪 **สำหรับการทดสอบ API**: ใช้ตัวอย่างด้านล่างร่วมกับ Postman หรือ curl  
> 🔐 **โปรดแนบ `Authorization: Bearer YOUR_JWT_TOKEN` สำหรับทุกคำขอที่ต้องยืนยันตัวตน**

#### 2.4.0 Pagination & Response Format Guidelines

> ℹ️ **แนวทางสำหรับทีมพัฒนาที่จะเติม TODO ภายหลัง**  
> ใช้รูปแบบการตอบกลับให้สอดคล้องกันทุก endpoint ที่คืนลิสต์ข้อมูล

- **โครงสร้างการตอบกลับ** (ตัวอย่างสำหรับ `/api/products`):
  ```json
  {
    "items": [
      {
        "id": 1,
        "name": "Gaming Mouse RGB",
        "description": "...",
        "price": 1299,
        "stock": 50,
        "category_id": 1,
        "image_url": "https://example.com/mouse.jpg"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 12,
      "total_items": 48,
      "total_pages": 4
    }
  }
  ```
- **Default**: หากไม่ส่ง `page` และ `limit` ให้ตั้ง `page=1`, `limit=12`
- **Validation**: ตรวจสอบว่า `page` ≥ 1, `limit` อยู่ในช่วงที่ยอมรับ (เช่น 1–100) ถ้าไม่ผ่านให้ตอบ `400`
- **Order History**: ใช้รูปแบบเดียวกัน (`items` + `pagination`) เพื่อให้ Frontend/Kong ใช้โครงสร้างร่วมกันได้
- **Cart/Order Detail (single resource)**: ให้ส่งเป็น object เดียว เช่น `{ "cart": {...} }` หรือ `{ "order": {...} }`

#### 2.4.1 Product Catalog

- **GET http://localhost:8082/api/products?page=1&limit=12&category=headset&search=rgb**  
  _ผู้รับผิดชอบ: ณัฐพงษ์ ดีบุตร_  
  ใช้สำหรับค้นหา/กรองสินค้าเพื่อแสดงรายการ ผลลัพธ์จะส่งกลับตาม Query Parameters ที่กำหนด  
  **การตั้งค่าใน Postman**
  1. เลือก Method เป็น `GET` แล้ววาง URL `http://localhost:8082/api/products`
  2. ไปที่แท็บ **Params** แล้วเพิ่มคีย์ต่าง ๆ ดังนี้ (เลือกเฉพาะที่ต้องการ):
     - `page`: ลำดับหน้าของผลลัพธ์ เช่น `1`
     - `limit`: จำนวนสินค้าต่อหน้า เช่น `12`
     - `category`: slug หรือรหัสหมวดหมู่ เช่น `headset`
     - `search`: คำค้นหาที่ต้องการ เช่น `rgb`
  3. Postman จะประกอบ URL ให้เป็น `http://localhost:8082/api/products?page=1&limit=12&category=headset&search=rgb`
  4. กด **Send** เพื่อรับรายการสินค้าที่กรองแล้ว
  > หากไม่ระบุพารามิเตอร์ ระบบจะส่งสินค้าทั้งหมดของหน้าที่ 1 พร้อมค่า limit เริ่มต้น

- **GET http://localhost:8082/api/products/1**  
  _ผู้รับผิดชอบ: ณัฐพงษ์ ดีบุตร_  
  ใช้สำหรับดูรายละเอียดสินค้าเฉพาะชิ้น
  ระบุ `:id` ตามสินค้าที่ต้องการดูรายละเอียด

#### 2.4.2 Cart Requests

**GET http://localhost:8082/api/cart**

ใช้สำหรับดูรายการสินค้าในตะกร้าของผู้ใช้ปัจจุบัน

_ผู้รับผิดชอบ: วายุ กอคูณ_

> วางใน Headers => Key: Authorization, Value: Bearer YOUR_JWT_TOKEN

```text
Authorization: Bearer YOUR_JWT_TOKEN
```

**POST http://localhost:8082/api/cart/add**

ใช้สำหรับเพิ่มสินค้าพร้อมระบุจำนวนลงในตะกร้า

_ผู้รับผิดชอบ: วายุ กอคูณ_

> วางใน Headers => Key: Authorization, Value: Bearer YOUR_JWT_TOKEN  
> วางใน Body => Raw (JSON)

```json
{
  "product_id": 42,
  "quantity": 2
}
```

**PUT http://localhost:8082/api/cart/update**

ใช้สำหรับแก้ไขจำนวนสินค้าที่อยู่ในตะกร้า

_ผู้รับผิดชอบ: วายุ กอคูณ_

> วางใน Headers => Key: Authorization, Value: Bearer YOUR_JWT_TOKEN  
> วางใน Body => Raw (JSON)

```json
{
  "cart_item_id": 101,
  "quantity": 3
}
```

**DELETE http://localhost:8082/api/cart/remove**

ใช้สำหรับลบสินค้าชิ้นหนึ่งออกจากตะกร้า

_ผู้รับผิดชอบ: วายุ กอคูณ_

> วางใน Headers => Key: Authorization, Value: Bearer YOUR_JWT_TOKEN  
> วางใน Body => Raw (JSON)

```json
{
  "cart_item_id": 101
}
```

#### 2.4.3 Orders Requests

**POST http://localhost:8082/api/orders**

ใช้สำหรับสร้างคำสั่งซื้อจากตะกร้าปัจจุบัน

_ผู้รับผิดชอบ: วายุ กอคูณ_

> วางใน Headers => Key: Authorization, Value: Bearer YOUR_JWT_TOKEN  
> วางใน Body => Raw (JSON)

```json
{
  "cart_id": 12,
  "shipping_address": "123 ถนนสุขุมวิท แขวงคลองตัน เขตวัฒนา กรุงเทพฯ 10110",
  "payment_method": "credit_card",
  "notes": "กรุณาจัดส่งในช่วงเวลา 9:00-17:00"
}
```

**GET http://localhost:8082/api/orders/history**

ใช้สำหรับดูประวัติคำสั่งซื้อทั้งหมดของผู้ใช้

_ผู้รับผิดชอบ: วายุ กอคูณ_

> วางใน Headers => Key: Authorization, Value: Bearer YOUR_JWT_TOKEN

```text
Authorization: Bearer YOUR_JWT_TOKEN
```

### 2.5 Error Responses

**ตัวอย่าง Error Response:**

```json
{
  "status": "error",
  "message": "Product not found",
  "error": "The requested product does not exist"
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
| **Developers**      | ณัฐพงษ์ ดีบุตร    | [GitHub](https://github.com/Natthaphong66)  |
| **Developers**      | วายุ กอคูณ        | [GitHub](https://github.com/FUJIKOTH)       |
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
