# 🚀 Mini-Project: GameGear E-commerce (Microservice Architecture + Kong API Gateway)

![Go](https://img.shields.io/badge/Go-1.25.3-00ADD8?style=for-the-badge&logo=go)
![Gin](https://img.shields.io/badge/Gin-Framework-008ECF?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql)
![GORM](https://img.shields.io/badge/GORM-B93527?style=for-the-badge)
![JWT](https://img.shields.io/badge/Auth-JWT-FF6F00?style=for-the-badge)
![Kong](https://img.shields.io/badge/API%20Gateway-Kong-003459?style=for-the-badge&logo=kong)
![Konga](https://img.shields.io/badge/Admin%20UI-Konga-00C7B7?style=for-the-badge)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker)

โปรเจกต์นี้คือระบบ **E-commerce** สำหรับอุปกรณ์เกมมิ่ง บนสถาปัตยกรรม **Microservices** โดยมี **Kong API Gateway (DB Mode)** เป็นจุดทางเข้าเดียว (single entry point) พร้อมกับ **Konga UI** สำหรับจัดการ Gateway แบบ GUI เพื่อควบคุม Routing, Authentication (JWT), Rate Limiting, CORS, และ Logging ให้กับบริการภายใน (`users-service`, `shop-service`, `admin-service`).

---

## 📋 Table of Contents

- [🏛️ ภาพรวมระบบ (System Overview)](#%EF%B8%8F-ภาพรวมระบบ-  - [องค์ประกอบหลัก](#องค์ประกอบหลัก)
  - [🔑 รายละเอียดเทคนิค](#-รายละเอียดเทคนิค)
  - [Request Flow Example](#request-flow-example)
- [🚀 เริ่มต้นใช้งาน (Getting Started)](#-เริ่มต้นใช้งาน-getting-started)
  - [ทำไมเลือก DB Mode + Konga](#ทำไมเลือก-db-mode--konga)
- [⚡ ติดตั้งและรันระบบ (Installation & Setup)](#-ติดตั้งและรันระบบ-installation--setup)
  - [🎯 ภาพรวมของส่วนนี้](#-ภาพรวมของส่วนนี้)
  - [Prerequisites](#prerequisites)
  - [เริ่ม Kong Gateway + Konga](#1-เริ่ม-kong-gateway--konga)
  - [รัน Services บนเครื่อง](#2-รัน-services-สำหรับ-developers)
  - [ตรวจสอบ Services](#3-ตรวจสอบ-services)
  - [เปิด Konga UI และผูกกับ Kong](#4-เปิด-konga-ui-และผูกกับ-kong)
  - [ตั้งค่า Services และ Routes ใน Konga](#5-ตั้งค่า-services-และ-routes-ใน-konga)
  - [ทดสอบผ่าน Kong Proxy](#53-ทดสอบผ่าน-kong-proxy)
  - [เปิด Kong Proxy เป็น Public URL](#54-เปิด-kong-proxy-เป็น-public-url-ตัวเลือกด้วย-ngrok)
- [🔌 ตั้งค่า Plugins](#-ตั้งค่า-plugins)
  - [CORS Plugin สำหรับ Frontend](#1-cors-plugin-สำหรับ-frontend)
  - [JWT Plugin สำหรับ Authentication](#2-jwt-plugin-สำหรับ-authentication)
  - [Rate Limiting Plugin ป้องกัน DDoS](#3-rate-limiting-plugin-ป้องกัน-ddos)
- [📡 API Endpoints Overview](#-api-endpoints-overview-ผ่าน-kong-proxy)
  - [Users Service](#users-service-users)
  - [Shop Service](#shop-service-shop)
  - [Admin Service](#admin-service-admin)
  - [Quick Endpoint Checks](#quick-endpoint-checks)
- [📁 ไฟล์ Docker Compose](#-ไฟล์-docker-compose)
- [🌐 รายการ Ports](#-รายการ-ports)
- [📚 เอกสารและแหล่งข้อมูล](#-เอกสารและแหล่งข้อมูล)
  - [Project Documentation (Google Docs)](#project-documentation-google-docs)
  - [Kong Documentation](#kong-documentation)
  - [Konga Documentation](#konga-documentation)
  - [Go + Gin Framework](#go--gin-framework)
- [📦 โครงสร้างโปรเจค (Project Structure)](#-โครงสร้างโปรเจค-project-structure)
  - [Local Development Setup](#local-development-setup)
  - [Dependencies Management](#dependencies-management)
- [📂 Service Repositories](#-service-repositories)
- [👥 ทีมพัฒนา (Development Team)](#-ทีมพัฒนา-development-team)

-included)
- [📞 ติดต่อและสนับสนุน](#-ติดต่อและสนับสนุน)

---

## 🏛️ ภาพรวมระบบ (System Overview)

### องค์ประกอบหลัก

ระบบ **GameGear E-commerce** ใช้สถาปัตยกรรม **Microservices Architecture** ที่ประกอบด้วย 3 ชั้นหลัก (3-Tier Architecture):

#### 🌐 **Client Layer (ชั้นลูกค้า)**

- **Web Application** - เว็บไซต์หลักสำหรับลูกค้า (Customer-facing Web App)
- **Mobile App** - แอปพลิเคชันมือถือ (Mobile Application)
- **Admin Dashboard** - หน้าจอจัดการสำหรับผู้ดูแลระบบ (Administrative Interface)

#### 🦍 **API Gateway Layer (ชั้น Gateway)**

- **Kong Gateway** - จุดทางเข้าเดียว (Single Entry Point) ที่จัดการ:
  - **Routing** - นำทางคำขอไปยัง Service ที่ถูกต้อง (Request Routing)
  - **Authentication** - ตรวจสอบ JWT Token (Token Validation)
  - **Rate Limiting** - จำกัดจำนวนคำขอต่อวินาที (Request Throttling)
  - **CORS** - อนุญาตการเรียกข้าม Domain (Cross-Origin Resource Sharing)
  - **Logging** - บันทึกการเข้าถึงและข้อผิดพลาด (Request/Response Logging)
- **Konga UI** - หน้าจอจัดการ Kong แบบ GUI (Web-based Admin Dashboard)
- **PostgreSQL Databases** - เก็บการตั้งค่า Kong และ Konga (Configuration Storage)

#### 🏢 **Microservices Layer (ชั้นบริการ)**

- **👤 Users Service** - จัดการผู้ใช้และการยืนยันตัวตนของทั้งสมาชิกทั่วไปและแอดมิน (Role-based User & Admin Authentication)
- **🛍️ Shop Service** - จัดการสินค้า, ตะกร้า, คำสั่งซื้อ รวมถึงช่องทางสำหรับแอดมินในการดูแลแคตตาล็อก (Product Catalog, Shopping Cart, Orders & Admin Catalog APIs)
c
- **🗄️ PostgreSQL Databases** - ฐานข้อมูลแยกตาม Service (Service-specific Databases)

#### 🔄 **การทำงานของระบบ (System Flow)**

1. **Client** ส่งคำขอ HTTP มาที่ **Kong Gateway** (Port 8000) - _Client sends HTTP request to Kong Gateway_
2. **Kong** ตรวจสอบ JWT Token, Rate Limiting, CORS - _Kong validates JWT token, applies rate limiting and CORS_
3. **Kong** นำทางคำขอไปยัง **Service ที่ถูกต้อง** ตาม Route - _Kong routes request to appropriate service based on routing rules_
4. **Service** ประมวลผลและส่งคำตอบกลับผ่าน **Kong** - _Service processes request and returns response through Kong_
5. **Kong** ส่งคำตอบกลับไปยัง **Client** - _Kong forwards response back to client_

#### 🔐 **การยืนยันตัวตนและสิทธิ์ของแอดมิน**

- **แหล่งเก็บข้อมูลบัญชีกลาง**: การสมัคร, ล็อกอิน, ลืมรหัสผ่าน และรีเซ็ตรหัสผ่านของแอดมินถูกจัดการโดย `users-service` เช่นเดียวกับผู้ใช้ทั่วไป แต่จะระบุ `role = "admin"` ภายใน JWT
- **Admin Service ไม่เก็บรหัสผ่าน**: `admin-service` ทำหน้าที่เรียกใช้ API อื่นผ่าน Kong เท่านั้น และตรวจสอบความถูกต้องของคำขอด้วย JWT ที่ออกจาก `users-service`
- **สิทธิ์จัดการสินค้า**: `shop-service` เปิดเผยชุด API `/api/admin/products` (POST/PUT/DELETE) สำหรับแอดมินโดยเฉพาะ พร้อม Middleware ตรวจสอบว่า JWT มีค่า `role=admin` ก่อนอนุญาตให้แก้ไขข้อมูลสินค้า
- **การสื่อสารผ่าน Kong**: Kong Forward `Authorization` header ไปยัง backend services เพื่อให้ middleware ตรวจสอบต่อได้ และสามารถเสริม Plugin เช่น JWT/ACL ได้ตามต้องการ
- **Shared Secret**: `users-service` และ `shop-service` ต้องตั้งค่า `JWT_SECRET_KEY` เหมือนกัน (ผ่าน `.env`) เพื่อให้สามารถตรวจสอบโทเคนร่วมกันได้

#### 🎯 **ข้อดีของสถาปัตยกรรมนี้ (Architecture Benefits)**

- **Scalability** - สามารถขยายแต่ละ Service แยกกันได้ (Independent scaling of services)
- **Maintainability** - แก้ไข Service เดียวไม่กระทบ Service อื่น (Isolated maintenance and updates)
- **Security** - มีจุดควบคุมความปลอดภัยที่เดียว (Centralized security control via Kong)
- **Flexibility** - ใช้เทคโนโลยีต่างกันได้ในแต่ละ Service (Technology-agnostic service development)
- **Team Development** - ทีมสามารถพัฒนา Service แยกกันได้ (Independent team development)

```mermaid
graph TB
    subgraph "Client Layer"
        WEB[🌐 Web Application]
        MOBILE[📱 Mobile App]
        ADMIN_UI[🛡️ Admin Dashboard]
    end

    subgraph "API Gateway Layer"
        KONGA[🖥️ Konga Admin UI<br/>:1337]
        KONG[🦍 Kong Gateway<br/>Proxy :8000<br/>Admin API :8001]
        KONG_DB[(🗄️ Kong PostgreSQL<br/>Configuration DB)]
        KONGA_DB[(🗄️ Konga PostgreSQL<br/>UI Database)]
    end

    subgraph "Microservices Layer"
        USERS[👤 Users Service<br/>:8081<br/>Auth & User Management]
        SHOP[🛍️ Shop Service<br/>:8082<br/>Products, Cart, Orders]
        ADMIN[🛡️ Admin Service<br/>:8083<br/>Backend Admin Panel]

        USERS_DB[(🗄️ Users DB<br/>PostgreSQL)]
        SHOP_DB[(🗄️ Shop DB<br/>PostgreSQL)]
    end

    %% Client connections
    WEB -->|HTTP Request| KONG
    MOBILE -->|HTTP Request| KONG
    ADMIN_UI -->|HTTP Request| KONG

    %% Konga manages Kong
    KONGA -.->|Manage Config| KONG
    KONGA -->|Store Settings| KONGA_DB
    KONG -->|Store Routes/Services| KONG_DB

    %% Kong routes to services
    KONG -->|/users/*| USERS
    KONG -->|/shop/*| SHOP
    KONG -->|/admin/*| ADMIN

    %% Admin service calls other services
    ADMIN -.->|API Calls via Kong| USERS
    ADMIN -.->|API Calls via Kong| SHOP

    %% Services connect to databases
    USERS -->|GORM| USERS_DB
    SHOP -->|GORM| SHOP_DB

    style KONG fill:#003459,stroke:#00ffff,stroke-width:3px,color:#fff
    style KONGA fill:#00C7B7,stroke:#00ffff,stroke-width:2px,color:#fff
    style USERS fill:#4A90E2,stroke:#333,stroke-width:2px,color:#fff
    style SHOP fill:#7B68EE,stroke:#333,stroke-width:2px,color:#fff
    style ADMIN fill:#FF6B6B,stroke:#333,stroke-width:2px,color:#fff
```

### 🔑 รายละเอียดเทคนิค

| Component            | Description                                                              | Port                              | Technology         |
| -------------------- | ------------------------------------------------------------------------ | --------------------------------- | ------------------ |
| **🦍 Kong Gateway**  | API Gateway หลัก - จัดการ Routing, JWT Validation, Rate Limiting, CORS   | 8000 (Proxy)<br/>8001 (Admin API) | Kong 3.4 (DB Mode) |
| **🖥️ Konga UI**      | Web-based Admin Dashboard สำหรับจัดการ Kong แบบ GUI                      | 1337                              | Konga (Node.js)    |
| **👤 Users Service** | จัดการผู้ใช้ทุกประเภท รวมถึงออก JWT พร้อม `role` สำหรับสมาชิกและแอดมิน | 8081                              | Go + Gin + GORM    |
| **🛍️ Shop Service**  | จัดการสินค้า, ตะกร้า, คำสั่งซื้อ และเปิด API สำหรับแอดมินดูแลแคตตาล็อก | 8082                              | Go + Gin + GORM    |
| **🛡️ Admin Service** | ระบบหลังบ้านที่ประสานงานกับ Users/Shop ผ่าน Kong โดยไม่เก็บข้อมูลรับรอง | 8083                              | Go + Gin           |
| **🗄️ PostgreSQL**    | ฐานข้อมูลสำหรับแต่ละ Service + Kong + Konga                              | 5432 (varies)                     | PostgreSQL 15/17   |

### Request Flow Example

```mermaid
sequenceDiagram
    participant AdminClient
    participant Kong
    participant UsersService
    participant ShopService

    AdminClient->>Kong: POST /users/api/auth/login<br/>{email, password}
    Kong->>UsersService: /api/auth/login
    UsersService-->>Kong: {token: "jwt", user: {role: "admin"}}
    Kong-->>AdminClient: JWT + Admin Profile

    Note over AdminClient: เก็บ JWT (role = admin)

    AdminClient->>Kong: POST /shop/api/admin/products<br/>Authorization: Bearer jwt
    Kong->>Kong: Validate JWT & ตรวจสอบ Plugin (ถ้ามี)
    Kong->>ShopService: /api/admin/products
    ShopService-->>Kong: {data: {...}}
    Kong-->>AdminClient: Product Created
```

---

---

## 🚀 เริ่มต้นใช้งาน (Getting Started)

โปรเจกต์นี้ใช้ **Kong Gateway (DB Mode)** พร้อม **Konga Admin UI** เพื่อให้สามารถจัดการ Services, Routes, และ Plugins ผ่าน Web Interface ได้แบบ Real-time

### ทำไมเลือก DB Mode + Konga?

| Feature                  | DB Mode + Konga ✅          | DB-less Mode ❌              |
| ------------------------ | --------------------------- | ---------------------------- |
| **GUI Management**       | มี Konga UI ใช้งานง่าย      | ไม่มี (ต้องแก้ YAML เอง)     |
| **Real-time Changes**    | แก้ได้ทันที ไม่ต้อง restart | ต้อง reload config ทุกครั้ง  |
| **Team Collaboration**   | หลายคนจัดการได้พร้อมกัน     | ต้อง commit YAML แล้ว deploy |
| **Plugin Configuration** | ตั้งค่าผ่าน UI สะดวก        | ต้องเขียน config เอง         |
| **Suitable For**         | Development & Production    | CI/CD, Static Config         |

---

## ⚡ ติดตั้งและรันระบบ (Installation & Setup)

> 👤 **ผู้รับผิดชอบ Kong Gateway Setup**: **วรรธนโรจน์ บุตรดี** (Project Manager)  
> 📝 หากมีปัญหาเกี่ยวกับ Kong Gateway, Konga UI หรือการตั้งค่า Routes/Plugins กรุณาติดต่อ

### 🎯 **ภาพรวมของส่วนนี้**

ส่วนนี้จะอธิบายขั้นตอนการตั้งค่า **Kong API Gateway** และ **Konga Admin UI** เพื่อให้ระบบ Microservices ทำงานร่วมกันได้ โดยมีขั้นตอนหลักดังนี้:

1. **🦍 Kong Gateway + Konga** - ติดตั้งและรัน Kong Gateway พร้อม Konga UI
2. **▶️ รัน Services บนเครื่อง** - เปิด users-service, shop-service, admin-service แบบ local
3. **🔧 ตั้งค่า Services และ Routes** - กำหนดค่า Services และ Routes ใน Konga
4. **🔌 ตั้งค่า Plugins** - กำหนดค่า CORS, JWT, Rate Limiting
5. **🧪 ทดสอบการเชื่อมต่อ** - ตรวจสอบว่าระบบทำงานถูกต้อง

> 💡 **หมายเหตุ**: ส่วนนี้เป็นขั้นตอนการตั้งค่า **Gateway และ Microservices** เพื่อให้ระบบทำงานร่วมกันได้

### Prerequisites

- Docker และ Docker Compose ติดตั้งแล้ว
- Port ว่าง: `8000`, `8001`, `1337`, `8081`, `8082`, `8083`

### 1. เริ่ม Kong Gateway + Konga

> 🎯 **จุดประสงค์**: ติดตั้งและรัน Kong Gateway พร้อม Konga UI เพื่อเป็นจุดกลางในการจัดการ Microservices

ใช้ไฟล์ `docker-compose.kong.yml` ที่มีอยู่ใน admin-service directory:

```bash
# จาก root directory (GameGear-Ecommerce/)
cd admin-service
docker compose -f docker-compose.kong.yml up -d
```

**หมายเหตุ:** Kong Gateway จัดการโดย PM (วรรธนโรจน์) ใน admin-service เท่านั้น

### 2. รัน Services (สำหรับ Developers)

**สำหรับ Users Service (ณิชพน):**

```bash
cd users-service
go run cmd/api/main.go
```

**สำหรับ Shop Service (ณัฐพงษ์ + วายุ):**

```bash
cd shop-service
go run cmd/api/main.go
```

**สำหรับ Admin Service (วรรธนโรจน์):**

```bash
cd admin-service
go run cmd/api/main.go
```

> 💡 ถ้า Kong รันใน Docker บนเครื่องเดียวกัน ให้ใช้ hostname `host.docker.internal` เพื่อให้ Kong เข้าถึง service ที่เปิดบนเครื่อง host ได้ (Linux บางระบบต้องเพิ่ม DNS เองหรือใช้ IP ของเครื่องแทน)

### 3. ตรวจสอบ Services

```bash
# ตรวจสอบว่า Kong ทำงาน
curl -i http://localhost:8001/

# ตรวจสอบ Konga UI
curl -i http://localhost:1337/
```

### 4. เปิด Konga UI และผูกกับ Kong

1. เปิด **http://localhost:1337** แล้วสมัคร/ล็อกอินครั้งแรก  
2. กด **+ Connection** → กรอก
   - `Name`: `GameGear Kong` (หรือชื่อที่ต้องการ)
   - `Kong Admin URL`: `http://kong:8001` *(ถ้า Kong รันใน Docker เดียวกัน)*  
     หรือ `http://host.docker.internal:8001` *(ถ้า Kong กับบริการอยู่คนละคอนเทนเนอร์/เครื่อง)*  
   - `Default Port`: 8001  
   กด **Create Connection** และ **Activate**

---

### 5. ตั้งค่า Services และ Routes ใน Konga

ทำซ้ำขั้นตอนต่อไปนี้กับทั้งสามบริการ (`users-service`, `shop-service`, `admin-service`)

### 5.1 สร้าง Service

- ไปที่ **Services → + ADD NEW SERVICE**
- กรอก:
  ```
  Name: <service-name>       // users-service, shop-service, admin-service
  Protocol: http
  Host: host.docker.internal
  Port: <service-port>       // 8081, 8082, 8083
  Path / Retries / Timeouts: ค่าเริ่มต้น
  ```
- กด **Create**

### 5.2 เพิ่ม Route ให้ Service

- ใน service ที่สร้าง → แท็บ **Routes → + ADD ROUTE**
- กรอก:
  ```
  Name: <service-name>-route
  Paths: <route-path>        // /users, /shop, /admin
  Strip Path: เปิด (true)
  Methods / Hosts / Headers: เว้นว่าง
  ```
- กด **Create**

> ต้องการแก้ไขหรือลบ ใช้ปุ่ม **Edit / Delete** ในหน้าเดียวกัน หรือจะเรียก Kong Admin API (`/services`, `/routes`) ก็ได้

### 5.3 ทดสอบผ่าน Kong Proxy

```bash
curl -sS http://localhost:8000/users/healthz
curl -sS http://localhost:8000/shop/healthz
curl -sS http://localhost:8000/admin/healthz
```

เมื่อต้องการทดสอบฟีเจอร์จริง ให้เรียก path ย่อยผ่าน Kong เช่น `/users/api/auth/login`, `/shop/api/products`, `/admin/api/admin/products`

### 5.4 เปิด Kong Proxy เป็น Public URL (ตัวเลือกด้วย ngrok)

1. ติดตั้งและล็อกอิน ngrok จาก [ngrok.com/download](https://ngrok.com/download) (ทำครั้งแรกเท่านั้น)
2. เปิดเทอร์มินัลใหม่ ขณะที่ Kong และ services ยังรันอยู่
3. สั่ง
   ```bash
   ngrok http 8000
   ```
4. คัดลอก **Forwarding URL** ที่ได้ เช่น `https://xxxx.ngrok-free.app`
5. ใช้ URL นี้แทน `http://localhost:8000` เมื่อแชร์ API ภายนอก
   ```bash
   curl -sS https://xxxx.ngrok-free.app/users/healthz
   ```

> โปรดเปิดเทอร์มินัล ngrok ค้างไว้ระหว่างใช้งาน และ URL แบบฟรีจะเปลี่ยนทุกครั้งที่เริ่ม ngrok ใหม่

---

## 🔌 ตั้งค่า Plugins

> 🎯 **จุดประสงค์**: กำหนดค่า Plugins (CORS, JWT, Rate Limiting) เพื่อเพิ่มความปลอดภัยและความสามารถให้กับ Kong Gateway

### 1. CORS Plugin สำหรับ Frontend

ใน Konga UI:

1. ไปที่ **Plugins** → **Add Global Plugin**
2. เลือก **CORS**
3. ตั้งค่า:
   ```
   origins: *
   methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
   headers: Authorization, Content-Type
   credentials: true
   max_age: 3600
   ```

### 2. JWT Plugin สำหรับ Authentication

สำหรับ Routes ที่ต้องการ Authentication:

1. ไปที่ Service/Route ที่ต้องการป้องกัน
2. เพิ่ม **JWT Plugin**
3. ตั้งค่า:
   ```
   key_claim_name: iss
   secret_is_base64: false
   ```

### 3. Rate Limiting Plugin ป้องกัน DDoS

1. เพิ่ม **Rate Limiting Plugin** แบบ Global หรือต่อ Service
2. ตั้งค่า:

| Field | Value  |
| ----- | ------ |
| minute | 100   |
| hour   | 10000 |
| policy | local |

---

## 📡 API Endpoints Overview (ผ่าน Kong Proxy)

> Base URL: `http://localhost:8000` (ถ้าใช้ ngrok ให้แทนด้วย URL จากหัวข้อ 5.4)

### Users Service (`/users`)

**Member Authentication**

- `POST /api/auth/register` — ลงทะเบียนสมาชิกใหม่
  ```bash
  curl -X POST http://localhost:8000/users/api/auth/register \
    -H "Content-Type: application/json" \
    -d '{
      "username": "tawan123",
      "display_name": "Tawan Gamer",
      "email": "tawan@example.com",
      "password": "password123",
      "confirm_password": "password123"
    }'
  ```
- `POST /api/auth/login` — เข้าสู่ระบบด้วยอีเมลหรือชื่อผู้ใช้
  ```bash
  curl -X POST http://localhost:8000/users/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"identifier":"tawan@example.com","password":"password123"}'
  ```
- `POST /api/auth/forgot-password` — ขออีเมลรีเซ็ตรหัสผ่าน
  ```bash
  curl -X POST http://localhost:8000/users/api/auth/forgot-password \
    -H "Content-Type: application/json" \
    -d '{"email":"tawan@example.com"}'
  ```
- `POST /api/auth/reset-password` — ตั้งรหัสผ่านใหม่ด้วยโทเคน
  ```bash
  curl -X POST http://localhost:8000/users/api/auth/reset-password \
    -H "Content-Type: application/json" \
    -d '{
      "token":"RESET_TOKEN_HERE",
      "new_password":"newpassword123",
      "confirm_password":"newpassword123"
    }'
  ```
- `POST /api/auth/logout` — ออกจากระบบและเพิกถอนโทเคน
  ```bash
  curl -X POST http://localhost:8000/users/api/auth/logout \
    -H "Authorization: Bearer MEMBER_JWT"
  ```

**Member Profile**

- `GET /api/user/profile` — ดูข้อมูลโปรไฟล์ผู้ใช้ปัจจุบัน
  ```bash
  curl http://localhost:8000/users/api/user/profile \
    -H "Authorization: Bearer MEMBER_JWT"
  ```
- `PUT /api/user/profile` — แก้ไขข้อมูลโปรไฟล์ / รูปภาพ / ชื่อแสดง / เปลี่ยนรหัสผ่าน / ลบบัญชี
  ```bash
  curl -X PUT http://localhost:8000/users/api/user/profile \
    -H "Authorization: Bearer MEMBER_JWT" \
    -H "Content-Type: application/json" \
    -d '{
      "display_name": "Tawan Updated",
      "profile_image": "https://example.com/avatar.png"
    }'
  ```
  - หากต้องการ **ลบบัญชี** ให้ส่ง `delete_account_flag=true` พร้อมรหัสผ่านปัจจุบัน เช่น
    ```bash
    curl -X PUT http://localhost:8000/users/api/user/profile \
      -H "Authorization: Bearer MEMBER_JWT" \
      -H "Content-Type: application/json" \
      -d '{
        "delete_account_flag": true,
        "password": "currentPassword123"
      }'
    ```

**Admin Authentication (ใช้โดย admin-service)**

- `POST /api/admin/login` — ล็อกอินแอดมินเพื่อรับ JWT
  ```bash
  curl -X POST http://localhost:8000/users/api/admin/login \
    -H "Content-Type: application/json" \
    -d '{"email":"admin@gamegear.com","password":"admin123"}'
  ```
- `POST /api/admin/register` — สร้างบัญชีแอดมินใหม่
  ```bash
  curl -X POST http://localhost:8000/users/api/admin/register \
    -H "Content-Type: application/json" \
    -d '{
      "email": "new.admin@gamegear.com",
      "password": "securePassword123",
      "confirm_password": "securePassword123",
      "display_name": "Operations Lead"
    }'
  ```
- `POST /api/admin/forgot-password` — ขออีเมลรีเซ็ตรหัสผ่านแอดมิน
  ```bash
  curl -X POST http://localhost:8000/users/api/admin/forgot-password \
    -H "Content-Type: application/json" \
    -d '{"email":"admin@gamegear.com"}'
  ```
- `POST /api/admin/reset-password` — ตั้งรหัสผ่านใหม่สำหรับแอดมิน
  ```bash
  curl -X POST http://localhost:8000/users/api/admin/reset-password \
    -H "Content-Type: application/json" \
    -d '{
      "token":"RESET_TOKEN_FROM_EMAIL",
      "new_password":"NewStrongPassword123",
      "confirm_password":"NewStrongPassword123"
    }'
  ```
- `POST /api/admin/logout` — ออกจากระบบแอดมิน
  ```bash
  curl -X POST http://localhost:8000/users/api/admin/logout \
    -H "Authorization: Bearer ADMIN_JWT"
  ```

> รายละเอียดเพิ่มเติมและตัวอย่างอื่น ๆ ดูได้ที่ [users-service README – JSON Request Examples](https://github.com/Wattanaroj2567/users-service/tree/develop?tab=readme-ov-file#23-json-request-examples)

### Shop Service (`/shop`)

**สินค้า (public)**

- `GET /api/products` — ดึงรายการสินค้า (รองรับ paging/filter)
  ```bash
  curl "http://localhost:8000/shop/api/products?page=1&limit=12"
  ```
- `GET /api/products/:id` — ดูรายละเอียดสินค้าแต่ละชิ้น
  ```bash
  curl http://localhost:8000/shop/api/products/1
  ```

**ตะกร้า (ต้องมี Member JWT)**

- `GET /api/cart` — ดูรายการสินค้าในตะกร้าปัจจุบัน
  ```bash
  curl http://localhost:8000/shop/api/cart \
    -H "Authorization: Bearer MEMBER_JWT"
  ```
- `POST /api/cart/add` — เพิ่มสินค้าเข้าในตะกร้า
  ```bash
  curl -X POST http://localhost:8000/shop/api/cart/add \
    -H "Authorization: Bearer MEMBER_JWT" \
    -H "Content-Type: application/json" \
    -d '{"product_id":42,"quantity":2}'
  ```
- `PUT /api/cart/update` — ปรับจำนวนสินค้าในตะกร้า
  ```bash
  curl -X PUT http://localhost:8000/shop/api/cart/update \
    -H "Authorization: Bearer MEMBER_JWT" \
    -H "Content-Type: application/json" \
    -d '{"cart_item_id":101,"quantity":3}'
  ```
- `DELETE /api/cart/remove` — ลบสินค้าออกจากตะกร้า
  ```bash
  curl -X DELETE http://localhost:8000/shop/api/cart/remove \
    -H "Authorization: Bearer MEMBER_JWT" \
    -H "Content-Type: application/json" \
    -d '{"cart_item_id":101}'
  ```

**คำสั่งซื้อ**

- `POST /api/orders`
  ```bash
  curl -X POST http://localhost:8000/shop/api/orders \
    -H "Authorization: Bearer MEMBER_JWT" \
    -H "Content-Type: application/json" \
    -d '{
      "cart_id": 12,
      "shipping_address": "123 ถนนสุขุมวิท กรุงเทพฯ 10110",
      "payment_method": "credit_card",
      "notes": "กรุณาส่งช่วง 9:00-17:00"
    }'
  ```
- `GET /api/orders/history`
  ```bash
  curl http://localhost:8000/shop/api/orders/history \
    -H "Authorization: Bearer MEMBER_JWT"
  ```

> ดูตัวอย่างอื่น ๆ ได้ที่ [shop-service README – JSON Request Examples](https://github.com/Wattanaroj2567/shop-service/tree/develop?tab=readme-ov-file#24-json-request-examples)

### Admin Service (`/admin`)

**Auth แอดมิน**

- `POST /api/admin/login` — ล็อกอินแอดมินผ่าน Kong
  ```bash
  curl -X POST http://localhost:8000/admin/api/admin/login \
    -H "Content-Type: application/json" \
    -d '{"email":"admin@gamegear.com","password":"admin123"}'
  ```
- `POST /api/admin/register` — สมัครแอดมินใหม่ (proxy ไป users-service)
  ```bash
  curl -X POST http://localhost:8000/admin/api/admin/register \
    -H "Content-Type: application/json" \
    -d '{
      "email": "new.admin@gamegear.com",
      "password": "securePassword123",
      "confirm_password": "securePassword123",
      "display_name": "Operations Lead"
    }'
  ```

**จัดการสินค้า (ต้องมี Admin JWT)**

- `GET /api/admin/products` — ดึงทั้งหมดมาจัดการในแดชบอร์ด
  ```bash
  curl http://localhost:8000/admin/api/admin/products \
    -H "Authorization: Bearer ADMIN_JWT"
  ```
- `POST /api/admin/products` — เพิ่มสินค้าใหม่ (ยิงไป shop-service)
  ```bash
  curl -X POST http://localhost:8000/admin/api/admin/products \
    -H "Authorization: Bearer ADMIN_JWT" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Gaming Keyboard RGB",
      "description": "คีย์บอร์ดเกมมิ่งพร้อมไฟ RGB",
      "price": 2599,
      "stock": 30,
      "category_id": 2,
      "image_url": "https://example.com/keyboard.jpg"
    }'
  ```
- `PUT /api/admin/products/:id` — ปรับปรุงข้อมูลสินค้าเดิม
  ```bash
  curl -X PUT http://localhost:8000/admin/api/admin/products/1 \
    -H "Authorization: Bearer ADMIN_JWT" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Gaming Keyboard RGB Pro",
      "description": "คีย์บอร์ดเกมมิ่งรุ่น Pro",
      "price": 2999,
      "stock": 20,
      "category_id": 2,
      "image_url": "https://example.com/keyboard-pro.jpg"
    }'
  ```
- `DELETE /api/admin/products/:id` — ลบสินค้าออกจากระบบ
  ```bash
  curl -X DELETE http://localhost:8000/admin/api/admin/products/1 \
    -H "Authorization: Bearer ADMIN_JWT"
  ```

**จัดการคำสั่งซื้อ**

- `GET /api/admin/orders`
  ```bash
  curl http://localhost:8000/admin/api/admin/orders \
    -H "Authorization: Bearer ADMIN_JWT"
  ```
- `PUT /api/admin/orders/:id/status`
  ```bash
  curl -X PUT http://localhost:8000/admin/api/admin/orders/1/status \
    -H "Authorization: Bearer ADMIN_JWT" \
    -H "Content-Type: application/json" \
    -d '{"status":"shipped"}'
  ```

> JSON ตัวอย่างทั้งหมดอยู่ใน [admin-service README – JSON Request Examples](https://github.com/Wattanaroj2567/admin-service/tree/develop?tab=readme-ov-file#24-json-request-examples)

### Quick Endpoint Checks

```bash
# Local ผ่าน Kong
curl -sS http://localhost:8000/users/healthz
curl -sS http://localhost:8000/shop/healthz
curl -sS http://localhost:8000/admin/healthz

# ผ่าน ngrok (แทน <NGROK_URL> ด้วย URL จากหัวข้อ 5.4)
curl -sS https://<NGROK_URL>/users/healthz
curl -sS https://<NGROK_URL>/shop/healthz
curl -sS https://<NGROK_URL>/admin/healthz
```

---

## 📁 ไฟล์ Docker Compose

```yaml
# =====================================================
# 🌉 Docker Compose – Kong Gateway (DB mode) + Konga UI
# =====================================================

networks:
  kong-net:
    driver: bridge

volumes:
  kong_data: # เก็บข้อมูลฐาน Kong (PostgreSQL)
  konga_data: # เก็บข้อมูลฐาน Konga (PostgreSQL)

services:
  # 🦍 Kong Gateway (DB mode)
  kong:
    image: kong:3.4
    container_name: kong-gateway
    networks: [kong-net]
    environment:
      KONG_DATABASE: "postgres"
      KONG_PG_HOST: kong-database
      KONG_PG_USER: kong
      KONG_PG_PASSWORD: kong
      KONG_PG_DATABASE: kong
      KONG_ADMIN_LISTEN: 0.0.0.0:8001
      KONG_PROXY_LISTEN: 0.0.0.0:8000
    ports:
      - "8000:8000" # Proxy Port
      - "8001:8001" # Admin API Port
    extra_hosts:
      - "host.docker.internal:host-gateway"
    depends_on:
      kong-database:
        condition: service_healthy
      kong-migrations:
        condition: service_completed_successfully
    healthcheck:
      test: ["CMD", "kong", "health"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # ⚙️ Kong Migrations
  kong-migrations:
    image: kong:3.4
    container_name: kong-migrations
    networks: [kong-net]
    environment:
      KONG_DATABASE: "postgres"
      KONG_PG_HOST: kong-database
      KONG_PG_USER: kong
      KONG_PG_PASSWORD: kong
      KONG_PG_DATABASE: kong
    command: >
      /bin/sh -lc "
        kong migrations bootstrap &&
        kong migrations up &&
        kong migrations finish
      "
    depends_on:
      kong-database:
        condition: service_healthy
    restart: "no"

  # 🗄️ PostgreSQL (Kong Database)
  kong-database:
    image: postgres:17
    container_name: kong-database
    networks: [kong-net]
    environment:
      POSTGRES_USER: kong
      POSTGRES_DB: kong
      POSTGRES_PASSWORD: kong
    volumes:
      - kong_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U kong"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # 🗄️ PostgreSQL (Konga Database)
  konga-database:
    image: postgres:9.5
    container_name: konga-database
    networks: [kong-net]
    environment:
      POSTGRES_USER: konga
      POSTGRES_DB: konga
      POSTGRES_PASSWORD: konga
    volumes:
      - konga_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U konga"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # ⚙️ Konga Prepare
  konga-prepare:
    image: pantsel/konga:latest
    container_name: konga-prepare
    networks: [kong-net]
    command: >
      -c prepare
      -a postgres
      -u postgresql://konga:konga@konga-database:5432/konga
    depends_on:
      konga-database:
        condition: service_healthy
    restart: "no"

  # 🖥️ Konga Admin UI
  konga:
    image: pantsel/konga:latest
    container_name: konga
    networks: [kong-net]
    environment:
      NODE_ENV: production
      DB_ADAPTER: postgres
      DB_HOST: konga-database
      DB_PORT: 5432
      DB_USER: konga
      DB_PASSWORD: konga
      DB_DATABASE: konga
      TOKEN_SECRET: some-secret-token-change-this
    ports:
      - "1337:1337"
    depends_on:
      konga-database:
        condition: service_healthy
      konga-prepare:
        condition: service_completed_successfully
    healthcheck:
      test: ["CMD-SHELL", "curl -fsS http://localhost:1337 || exit 1"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
```

---

## 🌐 รายการ Ports

| Service               | Port | Access URL            | Description                        |
| --------------------- | ---- | --------------------- | ---------------------------------- |
| 🦍 **Kong Proxy**     | 8000 | http://localhost:8000 | API Gateway - ใช้เรียก API ทั้งหมด |
| 🔧 **Kong Admin API** | 8001 | http://localhost:8001 | Admin API (Dev only)               |
| 🖥️ **Konga UI**       | 1337 | http://localhost:1337 | Web Dashboard จัดการ Kong          |
| 👤 **Users Service**  | 8081 | http://localhost:8081 | Direct access (Dev only)           |
| 🛍️ **Shop Service**   | 8082 | http://localhost:8082 | Direct access (Dev only)           |
| 🛡️ **Admin Service**  | 8083 | http://localhost:8083 | Direct access (Dev only)           |

> ⚠️ **Production**: ควรปิด Direct Access และให้เรียกผ่าน Kong Proxy (port 8000) เท่านั้น

---
## 📚 เอกสารและแหล่งข้อมูล

### Project Documentation (Google Docs)

เอกสารสำหรับการพัฒนาและส่งงาน:

| Document                       | Description                                          | Link                                                                                                                         |
| ------------------------------ | ---------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| 📝 **Step 1: Topic**           | หัวข้อและแนวคิดโปรเจกต์                              | [ส่งงานกลุ่ม Step 1 Topic](https://docs.google.com/document/d/1riGCYJim9EKQ7PVJ8JHrOqRWG2BxqlOrxZOhKO79I_k/edit?usp=sharing) |
| 📋 **Step 2: Planning**        | แผนการพัฒนาและ Architecture                          | [ส่งงานกลุ่ม Step 2](https://docs.google.com/document/d/1ieKM8_QVyB3G1ZUcbO9glSLP2tnpLTFH9ZcyCC3pK9c/edit?usp=drive_link)    |
| 🗄️ **Step 3: Database Design** | โครงสร้างฐานข้อมูลสำหรับระบบ E-commerce (SQL Schema) | [ส่งงานกลุ่ม Step 3](https://docs.google.com/document/d/1__wV-m49uxHTo2lNu0A6gCIJaP7aPcUNByorOFLFAiw/edit?usp=sharing)       |

> 💡 **หมายเหตุ**: เอกสารเหล่านี้เป็นส่วนหนึ่งของการส่งงานและใช้อ้างอิงในการพัฒนา กรุณาอ่านก่อนเริ่มทำงาน

---

### Kong Documentation

- [Kong Gateway Docs](https://docs.konghq.com/gateway/latest/)
- [Kong Admin API Reference](https://docs.konghq.com/gateway/latest/admin-api/)
- [Kong Plugins Hub](https://docs.konghq.com/hub/)

### Konga Documentation

- [Konga GitHub](https://github.com/pantsel/konga)
- [Konga Setup Guide](https://github.com/pantsel/konga#konga)

### Go + Gin Framework

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)

---

## 📦 โครงสร้างโปรเจค (Project Structure)

แต่ละ Service ใช้ Go Module สำหรับจัดการ dependencies:

| Service           | Module Name                         | Description                      |
| ----------------- | ----------------------------------- | -------------------------------- |
| **Users Service** | `github.com/gamegear/users-service` | Authentication & User Management |
| **Shop Service**  | `github.com/gamegear/shop-service`  | Products, Cart & Orders          |
| **Admin Service** | `github.com/gamegear/admin-service` | Backend Admin Panel              |

### Local Development Setup

สำหรับการพัฒนาในเครื่อง local แต่ละ service ใช้ `replace` directive ใน `go.mod`:

```go
// ใน admin-service/go.mod
replace github.com/gamegear/users-service => ../users-service
replace github.com/gamegear/shop-service => ../shop-service

// ใน shop-service/go.mod
replace github.com/gamegear/users-service => ../users-service
```

### Dependencies Management

- **Go Version**: 1.25.1 (ทุก service)
- **Database**: GORM + PostgreSQL (users-service, shop-service)
- **Web Framework**: Gin (ทุก service)
- **Authentication**: JWT (ทุก service)
- **API Testing**: Postman (ทุก service)

---

## 📂 Service Repositories

| Service Repository                                                          | Description                                             | Team Member                |
| --------------------------------------------------------------------------- | ------------------------------------------------------- | -------------------------- |
| 👤 **[Users Service](https://github.com/Wattanaroj2567/users-service.git)** | จัดการผู้ใช้และการยืนยันตัวตน (สมัคร, ล็อกอิน, โปรไฟล์) | ณิชพน มานิตย์              |
| 🛍️ **[Shop Service](https://github.com/Wattanaroj2567/shop-service.git)**   | จัดการสินค้า, ตะกร้า, คำสั่งซื้อ                        | ณัฐพงษ์ ดีบุตร, วายุ กอคูณ |
| 🛡️ **[Admin Service](https://github.com/Wattanaroj2567/admin-service.git)** | ระบบหลังบ้าน (Admin Panel)                              | วรรธนโรจน์ บุตรดี          |

---

## 👥 ทีมพัฒนา (Development Team)

| Profile                                                                                                                       | Name                  | Responsibility            |
| ----------------------------------------------------------------------------------------------------------------------------- | --------------------- | ------------------------- |
| [<img src="https://github.com/Wattanaroj2567.png" width="60" height="60"/>](https://github.com/Wattanaroj2567)                | **วรรธนโรจน์ บุตรดี** | Project Manager & Backend |
| [<img src="https://avatars.githubusercontent.com/u/159878532?v=4" width="60" height="60"/>](https://github.com/Natthaphong66) | **ณัฐพงษ์ ดีบุตร**    | Backend Developer         |
| [<img src="https://avatars.githubusercontent.com/u/159880199?v=4" width="60" height="60"/>](https://github.com/nitchapon66)   | **ณิชพน มานิตย์**     | Backend Developer         |
| [<img src="https://avatars.githubusercontent.com/u/160033040?v=4" width="60" height="60"/>](https://github.com/FUJIKOTH)      | **วายุ กอคูณ**        | Backend Developer         |

---
