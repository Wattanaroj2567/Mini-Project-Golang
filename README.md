# 🚀 Mini-Project: GameGear E-commerce (Microservice Architecture)

![Go](https://img.shields.io/badge/Go-1.24.6-00ADD8?style=for-the-badge\&logo=go)
![Gin](https://img.shields.io/badge/Gin-Framework-008ECF?style=for-the-badge\&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge\&logo=postgresql)
![GORM](https://img.shields.io/badge/GORM-B93527?style=for-the-badge)
![JWT](https://img.shields.io/badge/Auth-JWT-FF6F00?style=for-the-badge)
![Swagger](https://img.shields.io/badge/API-Swagger-85EA2D?style=for-the-badge\&logo=swagger)

โปรเจกต์นี้คือระบบ **E-commerce** สำหรับขายอุปกรณ์เกมมิ่ง ซึ่งถูกพัฒนาขึ้นโดยใช้สถาปัตยกรรมแบบ **Microservice** และใช้ **Gin Framework** เป็น Backend หลักในการพัฒนา API

---

## 🏛️ System Architecture Overview

สถาปัตยกรรมของโปรเจกต์ประกอบด้วย Service หลัก 3 ส่วนที่ทำงานแยกจากกัน ได้แก่ `users-service`, `shop-service`, และ `admin-service` โดยมี **API Gateway กลาง** เป็นจุดรับ Request หลักจาก Client แล้วส่งต่อไปยัง Service ที่เกี่ยวข้อง

### องค์ประกอบหลัก

1. **Client Apps**: ฝั่งผู้ใช้งาน (Web และ Mobile) ที่จะส่งคำขอไปยังระบบ Backend
2. **API Gateway**: ทำหน้าที่เป็นประตูทางเข้าหลัก รับคำขอจากผู้ใช้ทั้งหมด แล้วกระจายไปยัง Service ที่เหมาะสม (Users / Shop / Admin)
3. **Users Service**: จัดการข้อมูลผู้ใช้ทั้งหมด (สมัคร, ล็อกอิน, โปรไฟล์, การรีเซ็ตรหัสผ่าน)
4. **Shop Service**: จัดการข้อมูลสินค้า, ตะกร้า และคำสั่งซื้อ
5. **Admin Service**: ใช้สำหรับ Admin Panel โดยเฉพาะ เพื่อเรียกใช้ Users Service หรือ Shop Service ผ่าน API Gateway

### ทำไมต้องใช้ API Gateway?

* **Separation of Concerns**: แยกหน้าที่ Gateway ออกจาก Logic ของ Admin Service
* **Scalability**: ป้องกัน bottleneck และรองรับการเพิ่ม Service ใหม่ ๆ ได้ง่าย
* **Flexibility**: เพิ่ม/เปลี่ยน Service ได้โดยแก้ไขเฉพาะ config ของ Gateway
* **Security**: Gateway เป็นด่านแรกที่ตรวจสอบ Auth และป้องกันการเข้าถึงโดยตรงถึง Service ภายใน

---

## 🛠️ Tech Stack

* **Language & Framework:** Go + Gin Web Framework
* **Database:** PostgreSQL
* **ORM:** GORM
* **Authentication:** JWT
* **API Documentation:** Swagger (OpenAPI)
* **Architecture:** Microservices + API Gateway + Standard Go Layout

---

## 📂 Service Repositories

| Service Repository                                                           | Description                                             | Team Member                |
| ---------------------------------------------------------------------------- | ------------------------------------------------------- | -------------------------- |
| 👤 **[users-service](https://github.com/Wattanaroj2567/users-service.git)**  | จัดการผู้ใช้และการยืนยันตัวตน (สมัคร, ล็อกอิน, โปรไฟล์) | ณิชพน มานิตย์              |
| 🛍️ **[shop-service](https://github.com/Wattanaroj2567/shop-service.git)**   | จัดการสินค้า, ตะกร้า, คำสั่งซื้อ                        | ณัฐพงษ์ ดีบุตร, วายุ กอคูณ |
| 🛡️ **[admin-service](https://github.com/Wattanaroj2567/admin-service.git)** | ระบบหลังบ้าน (Admin Panel)                              | วรรธนโรจน์ บุตรดี          |

---

## 🤝 Development Team

| Name                  | Responsibility            |
| --------------------- | ------------------------- |
| **วรรธนโรจน์ บุตรดี** | Project Manager & Backend |
| **ณัฐพงษ์ ดีบุตร**    | Backend Developer         |
| **ณิชพน มานิตย์**     | Backend Developer         |
| **วายุ กอคูณ**        | Backend Developer         |

---

## 🚀 How to Run the Entire Project

### 1. Clone all repositories

```bash
git clone https://github.com/Wattanaroj2567/users-service.git
git clone https://github.com/Wattanaroj2567/shop-service.git
git clone https://github.com/Wattanaroj2567/admin-service.git
```

### 2. Setup & Run Each Service (Manual)

เข้าไปในแต่ละโฟลเดอร์ของ Service แล้วทำตาม `README.md` ของ Service นั้น ๆ เพื่อตั้งค่าและรันโปรเจกต์

* **Users Service (Port: 8080)**

```bash
cd users-service
# ทำตามขั้นตอนใน README.md
```

* **Shop Service (Port: 8081)**

```bash
cd shop-service
# ทำตามขั้นตอนใน README.md
```

* **Admin Service (Port: 8082)**

```bash
cd admin-service
# ทำตามขั้นตอนใน README.md
```

---

## 🐋 Run with Docker (Recommended)

> ต้องติดตั้ง Docker และ Docker Compose ก่อน

### A) Quick Start

```bash
docker compose -f docker-compose.dev.yml up --build
```

เปิดใช้งานที่:

* Users Service → [http://localhost:8080](http://localhost:8080)
* Shop Service → [http://localhost:8081](http://localhost:8081)
* Admin Service → [http://localhost:8082](http://localhost:8082)

หยุดบริการทั้งหมด:

```bash
docker compose -f docker-compose.dev.yml down
```

### B) Example: `docker-compose.dev.yml` (ย่อ)

```yaml
version: "3.9"
services:
  users-db:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: gamegear_users_db
      POSTGRES_USER: dev
      POSTGRES_PASSWORD: dev
  shop-db:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: gamegear_shop_db
      POSTGRES_USER: dev
      POSTGRES_PASSWORD: dev

  users-service:
    build: { context: ./users-service, dockerfile: Dockerfile.dev }
    environment:
      DATABASE_URL: postgres://dev:dev@users-db:5432/gamegear_users_db?sslmode=disable
      APPLICATION_PORT: 8080
    ports: ["8080:8080"]
    depends_on: [users-db]
    volumes: ["./users-service:/app"]
    command: air

  shop-service:
    build: { context: ./shop-service, dockerfile: Dockerfile.dev }
    environment:
      DATABASE_URL: postgres://dev:dev@shop-db:5432/gamegear_shop_db?sslmode=disable
      APPLICATION_PORT: 8081
    ports: ["8081:8081"]
    depends_on: [shop-db]
    volumes: ["./shop-service:/app"]
    command: air

  admin-service:
    build: { context: ./admin-service, dockerfile: Dockerfile.dev }
    environment:
      USER_SERVICE_URL: http://users-service:8080
      SHOP_SERVICE_URL: http://shop-service:8081
      APPLICATION_PORT: 8082
    ports: ["8082:8082"]
    depends_on: [users-service, shop-service]
    volumes: ["./admin-service:/app"]
    command: air
```

### C) Example: `Dockerfile.dev` (ทุก service ใช้โครงเดียวกัน)

```dockerfile
FROM golang:1.22-alpine
RUN apk add --no-cache git bash build-base \
 && go install github.com/cosmtrek/air@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
EXPOSE 8080
CMD ["air"]
```

> โหมด Dev ใช้ `air` เพื่อ hot-reload โค้ดทันทีเมื่อแก้ไฟล์

---

## ✅ Notes

* แต่ละ Service มี `.env` ของตัวเอง แยกค่าคอนฟิกชัดเจน
* ใน Docker network สามารถอ้างอิงกันด้วยชื่อ service ได้เลย (เช่น `http://users-service:8080`)
* ถ้าจะทำ production ให้ใช้ Dockerfile แบบ multi-stage build และ compose อีกรูปแบบหนึ่ง 
