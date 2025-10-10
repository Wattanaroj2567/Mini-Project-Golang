# 🚀 Mini-Project: GameGear E-commerce (Microservice Architecture + Kong API Gateway)

![Go](https://img.shields.io/badge/Go-1.24.6-00ADD8?style=for-the-badge\&logo=go)
![Gin](https://img.shields.io/badge/Gin-Framework-008ECF?style=for-the-badge\&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge\&logo=postgresql)
![GORM](https://img.shields.io/badge/GORM-B93527?style=for-the-badge)
![JWT](https://img.shields.io/badge/Auth-JWT-FF6F00?style=for-the-badge)
![Swagger](https://img.shields.io/badge/API-Swagger-85EA2D?style=for-the-badge\&logo=swagger)
![Kong](https://img.shields.io/badge/API%20Gateway-Kong-003459?style=for-the-badge\&logo=kong)

โปรเจกต์นี้คือระบบ **E-commerce** สำหรับอุปกรณ์เกมมิ่ง บนสถาปัตยกรรม **Microservices** โดยมี **Kong API Gateway** เป็นจุดทางเข้าเดียว (single entry point) เพื่อจัดการ Routing, Auth (JWT), และ Observability ให้กับบริการภายใน (`users-service`, `shop-service`, `admin-service`).

---

## 🏛️ System Architecture Overview

![System Architecture](https://drive.google.com/uc?export=view\&id=1PaLRNsrbhVisgQUEvg1LMezzYZ7zcBhw)

### องค์ประกอบหลัก

* **Client Apps (Web/Mobile)** → เรียก API ผ่าน Gateway เดียว
* **Kong API Gateway** → ตรวจสอบ JWT, จัดเส้นทาง, บันทึก Metrics/Logs
* **Users Service** → สมัคร/ล็อกอิน/โปรไฟล์/สิทธิ์
* **Shop Service** → สินค้า/ตะกร้า/คำสั่งซื้อ
* **Admin Service** → แผงผู้ดูแล เรียกใช้ Users/Shop ผ่าน Gateway

---

## 📘 API Documentation (OpenAPI / Swagger)

ทุก Service มี Swagger UI สำหรับอ้างอิงและทดสอบ API:

| Service       | Swagger Endpoint                                                                     | Description                |
| ------------- | ------------------------------------------------------------------------------------ | -------------------------- |
| Users Service | [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) | สมัคร, ล็อกอิน, โปรไฟล์    |
| Shop Service  | [http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html) | สินค้า, ตะกร้า, คำสั่งซื้อ |
| Admin Service | [http://localhost:8082/swagger/index.html](http://localhost:8082/swagger/index.html) | ฟังก์ชันผู้ดูแลระบบ        |

> Swagger ถูกสร้างจาก OpenAPI Spec ในโฟลเดอร์ `/docs` ของแต่ละ service

---

## 🧩 Kong Gateway Integration (DB-less, Declarative)

เพื่อให้ dev ทุกเครื่องรันได้ทันทีโดย **ไม่ต้อง Bootstrap DB ของ Kong**, เราใช้โหมด **DB-less** พร้อมไฟล์ประกาศ `kong.yml` ซึ่งกำหนด **services + routes** ไว้ล่วงหน้า

### 1) `kong.yml` (ประกาศ Services/Routes ที่ต้องใช้)

```yaml
_format_version: "3.0"

services:
  - name: users-service
    url: http://users-service:8080
    routes:
      - name: users-route
        paths: [ "/users" ]

  - name: shop-service
    url: http://shop-service:8081
    routes:
      - name: shop-route
        paths: [ "/shop" ]

  - name: admin-service
    url: http://admin-service:8082
    routes:
      - name: admin-route
        paths: [ "/admin" ]
```

### 2) `docker-compose.yml` (พร้อมใช้งานจริงสำหรับ Dev)

```yaml
version: "3.9"

services:
  # 🚪 Kong API Gateway (DB-less)
  kong:
    image: kong:latest
    container_name: kong
    environment:
      KONG_DATABASE: off                 # ใช้โหมด DB-less
      KONG_DECLARATIVE_CONFIG: /kong.yml # โหลด routes/services จากไฟล์ประกาศ
      KONG_PROXY_LISTEN: 0.0.0.0:8000
      KONG_ADMIN_LISTEN: 0.0.0.0:8001
    ports:
      - "8000:8000"   # Public Proxy
      - "8001:8001"   # Admin API (Dev เท่านั้น)
    volumes:
      - ./kong.yml:/kong.yml:ro
    depends_on:
      - users-service
      - shop-service
      - admin-service

  # 👤 Users Service
  users-service:
    build: ./users-service
    container_name: users-service
    environment:
      APPLICATION_PORT: 8080
    ports:
      - "8080:8080"
    # ✅ ให้แอปฟังที่ 0.0.0.0:8080 ด้านในคอนเทนเนอร์

  # 🛍️ Shop Service
  shop-service:
    build: ./shop-service
    container_name: shop-service
    environment:
      APPLICATION_PORT: 8081
    ports:
      - "8081:8081"

  # 🛡️ Admin Service
  admin-service:
    build: ./admin-service
    container_name: admin-service
    environment:
      APPLICATION_PORT: 8082
    ports:
      - "8082:8082"
```

### 3) วิธีใช้งาน

```bash
# เริ่มระบบทั้งหมด
docker compose up -d --build

# ทดสอบเรียกผ่าน Gateway
# Users Login → จะถูก proxy ไปยัง users-service:8080/login
curl http://localhost:8000/users/login -i

# ดูสถานะ Kong (Admin API – ใช้เฉพาะ Dev)
curl http://localhost:8001/ -s | jq
```

> ถ้าต้องการเปลี่ยน prefix ของ route เช่น `/api/users` ให้แก้ `paths` ใน `kong.yml` แล้ว restart คอนเทนเนอร์ `kong`

---

## 🧠 ทำไมเลือก DB-less สำหรับ Dev?

* ⚡ **เร็วและง่าย**: ไม่ต้องตั้งค่า Postgres + migrations ของ Kong
* 🧩 **ประกาศครั้งเดียว**: เพื่อนร่วมทีม `docker compose up` ก็ได้ routes พร้อมใช้ทันที
* 🔁 **ทำซ้ำได้**: ลดความต่างของสภาพแวดล้อมแต่ละเครื่อง

> หากจะไป Production ค่อยเปลี่ยนไปใช้ **Postgres-backed** + `decK`/GitOps ได้โดยไม่กระทบ services ภายใน

---

## 🐋 Healthcheck (แนะนำให้เพิ่มในแต่ละ Service)

ตัวอย่าง (Gin): เปิด endpoint `/healthz` คืน 200 OK

```go
r := gin.Default()
r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
```

เรียกผ่าน Gateway ได้ที่:

* `http://localhost:8000/users/healthz`
* `http://localhost:8000/shop/healthz`
* `http://localhost:8000/admin/healthz`

---

## 📂 Service Repositories

| Service Repository                                                                                                 | Description                                             | Team Member                |
| ------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------- | -------------------------- |
| 👤 **[Users Service](https://github.com/Wattanaroj2567/users-service.git)**  | จัดการผู้ใช้และการยืนยันตัวตน (สมัคร, ล็อกอิน, โปรไฟล์) | ณิชพน มานิตย์              |
| 🛍️ **[Shop Service](https://github.com/Wattanaroj2567/shop-service.git)**   | จัดการสินค้า, ตะกร้า, คำสั่งซื้อ                        | ณัฐพงษ์ ดีบุตร, วายุ กอคูณ |
| 🛡️ **[Admin Service](https://github.com/Wattanaroj2567/admin-service.git)** | ระบบหลังบ้าน (Admin Panel)                              | วรรธนโรจน์ บุตรดี          |

---

## 🤝 Development Team

| Profile                                                                                   | Name                  | Responsibility            |
| ----------------------------------------------------------------------------------------- | --------------------- | ------------------------- |
| [<img src="https://github.com/Wattanaroj2567.png" width="60" height="60"/>](https://github.com/Wattanaroj2567)               | **วรรธนโรจน์ บุตรดี** | Project Manager & Backend |
| [<img src="https://avatars.githubusercontent.com/u/159878532?v=4" width="60" height="60"/>](https://github.com/Natthaphong66) | **ณัฐพงษ์ ดีบุตร**    | Backend Developer         |
| [<img src="https://avatars.githubusercontent.com/u/159880199?v=4" width="60" height="60"/>](https://github.com/nitchapon66) | **ณิชพน มานิตย์**     | Backend Developer         |
| [<img src="https://avatars.githubusercontent.com/u/160033040?v=4" width="60" height="60"/>](https://github.com/FUJIKOTH) | **วายุ กอคูณ**        | Backend Developer         |

---

## ✅ Summary

* ใช้ **Kong (DB-less)** เป็น API Gateway กลาง กำหนด routes/services ด้วย `kong.yml`
* ทุก service มี Swagger UI และควรมี `/healthz` สำหรับ readiness
* โครงสร้างนี้พร้อมใช้งานจริงสำหรับ dev: **คัดลอก `docker-compose.yml` + `kong.yml` ไปวาง แล้ว `docker compose up -d` ได้ทันที**
