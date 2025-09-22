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

* พัฒนาเป็น **Web Application** ด้วยภาษา Go
* ใช้สถาปัตยกรรมแบบ **Microservices** แบ่ง Service ออกเป็น Users, Shop และ Admin
* การเชื่อมต่อฐานข้อมูลใช้ **PostgreSQL** ผ่าน ORM **GORM**
* การยืนยันตัวตนและการเข้าถึงระบบใช้ **JWT (JSON Web Tokens)**
* เอกสาร API ใช้ **OpenAPI/Swagger**

---

## 🛠️ Tech Stack

* **Language & Framework:** Go + Gin Web Framework
* **Database:** PostgreSQL
* **ORM:** GORM
* **Authentication:** JWT
* **API Documentation:** Swagger (OpenAPI)
* **Architecture:** Microservices + Standard Go Layout

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

### 2. Setup & Run

เข้าไปในแต่ละโฟลเดอร์ของ Service แล้วทำตาม `README.md` ของ Service นั้น ๆ
(เช่น config `.env`, migrate database, และ `go run cmd/api/main.go`)
