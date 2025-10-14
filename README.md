# 🚀 Mini-Project: GameGear E-commerce (Microservice Architecture + Kong API Gateway)

![Go](https://img.shields.io/badge/Go-1.25.1-00ADD8?style=for-the-badge&logo=go)
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

- [🏛️ ภาพรวมระบบ (System Overview)](#%EF%B8%8F-ภาพรวมระบบ-system-overview)
  - [องค์ประกอบหลัก](#-องค์ประกอบหลัก)
  - [Request Flow Example](#request-flow-example)
- [🚀 เริ่มต้นใช้งาน (Getting Started)](#-เริ่มต้นใช้งาน-getting-started)
  - [ทำไมเลือก DB Mode + Konga](#ทำไมเลือก-db-mode--konga)
- [⚡ ติดตั้งและรันระบบ (Installation & Setup)](#-ติดตั้งและรันระบบ-installation--setup)
  - [Prerequisites](#-prerequisites)
  - [Remote Development (ngrok หรือ Tunnel)](#-remote-development-ngrok-หรือ-tunnel)
  - [1. เริ่ม Kong Gateway + Konga](#1-เริ่ม-kong-gateway--konga)
  - [2️. รัน Services (สำหรับ Developers)](#2-รัน-services-สำหรับ-developers)
  - [3️. ตรวจสอบ Services](#3-ตรวจสอบ-services)
  - [4️. เปิด Konga UI และตั้งค่าครั้งแรก](#4-เปิด-konga-ui-และตั้งค่าครั้งแรก)
- [🔧 ตั้งค่า Services และ Routes](#-ตั้งค่า-services-และ-routes)
  - [1. เพิ่ม Users Service](#1-เพิ่ม-users-service)
  - [2. เพิ่ม Route สำหรับ Users Service](#2-เพิ่ม-route-สำหรับ-users-service)
  - [3. ทำซ้ำสำหรับ Shop และ Admin Service](#3-ทำซ้ำสำหรับ-shop-และ-admin-service)
  - [ทดสอบการเชื่อมต่อ](#ทดสอบการเชื่อมต่อ)
- [🔌 ตั้งค่า Plugins](#-ตั้งค่า-plugins)
  - [1. CORS Plugin สำหรับ Frontend](#1-cors-plugin-สำหรับ-frontend)
  - [2. JWT Plugin สำหรับ Authentication](#2-jwt-plugin-สำหรับ-authentication)
  - [3. Rate Limiting Plugin ป้องกัน DDoS](#3-rate-limiting-plugin-ป้องกัน-ddos)
- [📁 ไฟล์ Docker Compose](#-ไฟล์-docker-compose)
- [🌐 รายการ Ports](#-รายการ-ports)
- [🐋 Health Check Endpoints](#-health-check-endpoints)
  - [ทดสอบผ่าน Kong Gateway](#ทดสอบผ่าน-kong-gateway)
- [🔍 แก้ไขปัญหา (Troubleshooting)](#-แก้ไขปัญหา-troubleshooting)
  - [Problem: Kong ไม่สามารถเชื่อมต่อกับ Services](#problem-kong-ไม่สามารถเชื่อมต่อกับ-services)
  - [Problem: Konga ไม่สามารถเชื่อมต่อกับ Kong](#problem-konga-ไม่สามารถเชื่อมต่อกับ-kong)
  - [Problem: JWT Authentication ไม่ทำงาน](#problem-jwt-authentication-ไม่ทำงาน)
  - [Problem: Database Migration Errors](#problem-database-migration-errors)
  - [Problem: Port Already in Use](#problem-port-already-in-use)
  - [Useful Commands](#useful-commands)
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
- [✅ สรุปคุณสมบัติ (Features Summary)](#-สรุปคุณสมบัติ-features-summary)
  - [What You Get](#what-you-get)
  - [Quick Commands Recap](#quick-commands-recap)
  - [Documentation Structure](#documentation-structure)
  - [Learning Resources Included](#learning-resources-included)
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

- **👤 Users Service** - จัดการผู้ใช้และการยืนยันตัวตน (User Management & Authentication)
- **🛍️ Shop Service** - จัดการสินค้า, ตะกร้า, คำสั่งซื้อ (Product Catalog, Shopping Cart, Orders)
- **🛡️ Admin Service** - ระบบหลังบ้านที่ประสานงานกับ Services อื่น (Backend Admin Panel)
- **🗄️ PostgreSQL Databases** - ฐานข้อมูลแยกตาม Service (Service-specific Databases)

#### 🔄 **การทำงานของระบบ (System Flow)**

1. **Client** ส่งคำขอ HTTP มาที่ **Kong Gateway** (Port 8000) - _Client sends HTTP request to Kong Gateway_
2. **Kong** ตรวจสอบ JWT Token, Rate Limiting, CORS - _Kong validates JWT token, applies rate limiting and CORS_
3. **Kong** นำทางคำขอไปยัง **Service ที่ถูกต้อง** ตาม Route - _Kong routes request to appropriate service based on routing rules_
4. **Service** ประมวลผลและส่งคำตอบกลับผ่าน **Kong** - _Service processes request and returns response through Kong_
5. **Kong** ส่งคำตอบกลับไปยัง **Client** - _Kong forwards response back to client_

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
| **👤 Users Service** | จัดการผู้ใช้และการยืนยันตัวตน (Register, Login, Profile, Password Reset) | 8081                              | Go + Gin + GORM    |
| **🛍️ Shop Service**  | จัดการสินค้า, ตะกร้า, คำสั่งซื้อ                                         | 8082                              | Go + Gin + GORM    |
| **🛡️ Admin Service** | ระบบหลังบ้าน (Admin Panel) - ประสานงานกับ Users/Shop ผ่าน Kong           | 8083                              | Go + Gin           |
| **🗄️ PostgreSQL**    | ฐานข้อมูลสำหรับแต่ละ Service + Kong + Konga                              | 5432 (varies)                     | PostgreSQL 15/17   |

### Request Flow Example

```mermaid
sequenceDiagram
    participant Client
    participant Kong
    participant UsersService
    participant ShopService

    Client->>Kong: POST /users/login<br/>{email, password}
    Kong->>UsersService: Forward to :8081/login
    UsersService-->>Kong: {token: "jwt-token", user: {...}}
    Kong-->>Client: Response + JWT

    Note over Client: Client stores JWT

    Client->>Kong: GET /shop/products<br/>Authorization: Bearer jwt-token
    Kong->>Kong: Validate JWT
    Kong->>ShopService: Forward to :8081/products
    ShopService-->>Kong: {products: [...]}
    Kong-->>Client: Product List
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

1. **🌐 Remote Development** - อธิบายการใช้ ngrok สำหรับการพัฒนาแบบทีม
2. **🦍 Kong Gateway + Konga** - ติดตั้งและรัน Kong Gateway พร้อม Konga UI
3. **🔧 ตั้งค่า Services และ Routes** - กำหนดค่า Services และ Routes ใน Konga
4. **🔌 ตั้งค่า Plugins** - กำหนดค่า CORS, JWT, Rate Limiting
5. **🧪 ทดสอบการเชื่อมต่อ** - ตรวจสอบว่าระบบทำงานถูกต้อง

> 💡 **หมายเหตุ**: ส่วนนี้เป็นขั้นตอนการตั้งค่า **Gateway และ Microservices** เพื่อให้ระบบทำงานร่วมกันได้

### Prerequisites

- Docker และ Docker Compose ติดตั้งแล้ว
- Port ว่าง: `8000`, `8001`, `1337`, `8081`, `8082`, `8083`

### 🌐 Remote Development (ngrok หรือ Tunnel)

> 🎯 **จุดประสงค์**: อธิบายวิธีการใช้ ngrok สำหรับการพัฒนาแบบทีม เพื่อให้ services สื่อสารกันได้ผ่าน Kong Gateway

ในการพัฒนาแบบทีม แต่ละ service จะรัน **local** บนเครื่องของแต่ละคน (ไม่ใช้ Docker Compose) ดังนั้นจึงต้องใช้ **ngrok** เพื่อสร้าง tunnel ให้ services สื่อสารกันได้ผ่านอินเทอร์เน็ต

#### 🎯 **เหตุผลที่ต้องใช้ ngrok:**

| เหตุผล                             | รายละเอียด                                                                                                                                            |
| ---------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| 🌍 **การทำงานคนละเครื่อง**         | แต่ละคนรัน service บนเครื่องของตัวเอง (users-service, shop-service, admin-service) ต้องใช้ ngrok แชร์ URL ให้กันเพื่อให้เชื่อมต่อได้ผ่าน Kong Gateway |
| ⚙️ **Integration กับ Kong**        | Kong สามารถใช้ URL จาก ngrok เป็นปลายทาง service ได้โดยตรง เช่น `USER_SERVICE_URL=https://abc1234.ngrok.io` เพื่อ proxy ไปยังเครื่องเพื่อน            |
| 🚀 **สะดวกต่อการทดสอบ**            | ใช้ทดสอบระบบรวม (Integration Test) ระหว่าง service จริง โดยไม่ต้อง deploy ขึ้นเซิร์ฟเวอร์กลาง                                                         |
| 🎓 **ใช้สาธิต/ส่งอาจารย์ได้ทันที** | สามารถเปิดระบบในเครื่องแล้วแชร์ให้ผู้สอนหรือผู้ทดสอบเข้ามาเรียก API ได้ผ่านอินเทอร์เน็ต                                                               |

#### 📋 **วิธีการทำงาน:**

1. **แต่ละคนรัน service บนเครื่องของตัวเอง** (local development)
2. **ใช้ ngrok สร้าง tunnel** เพื่อแชร์ URL ให้เพื่อนร่วมทีม
3. **Kong Gateway ใช้ ngrok URL** เป็นปลายทางของแต่ละ service
4. **ทดสอบระบบรวม** ผ่าน Kong Gateway ที่ port 8000

> 🔸 **สรุป**: แต่ละ service รัน local + ใช้ ngrok → Kong เชื่อมต่อได้ครบทุก service

### 1️. เริ่ม Kong Gateway + Konga

> 🎯 **จุดประสงค์**: ติดตั้งและรัน Kong Gateway พร้อม Konga UI เพื่อเป็นจุดกลางในการจัดการ Microservices

ใช้ไฟล์ `docker-compose.kong.yml` ที่มีอยู่ใน admin-service directory:

```bash
# จาก root directory (GameGear-Ecommerce/)
cd admin-service
docker compose -f docker-compose.kong.yml up -d
```

**หมายเหตุ:** Kong Gateway จัดการโดย PM (วรรธนโรจน์) ใน admin-service เท่านั้น

### 2️. รัน Services (สำหรับ Developers)

**สำหรับ Users Service (ณิชพน):**

```bash
cd users-service
go run cmd/api/main.go

# Terminal ใหม่
ngrok http 8081
# แชร์ URL: https://abc123.ngrok.io
```

**สำหรับ Shop Service (ณัฐพงษ์ + วายุ):**

```bash
cd shop-service
go run cmd/api/main.go

# Terminal ใหม่
ngrok http 8082
# แชร์ URL: https://def456.ngrok.io
```

**สำหรับ Admin Service (วรรธนโรจน์):**

```bash
cd admin-service
go run cmd/api/main.go

# Terminal ใหม่
ngrok http 8083
# แชร์ URL: https://ghi789.ngrok.io
```

### 3️. ตรวจสอบ Services

```bash
# ตรวจสอบว่า Kong ทำงาน
curl -i http://localhost:8001/

# ตรวจสอบ Konga UI
curl -i http://localhost:1337/
```

### 4️. เปิด Konga UI และตั้งค่าครั้งแรก

1. เปิด browser ไปที่ **http://localhost:1337**
2. สร้างบัญชี Admin ครั้งแรก:

   - **Username**: `admin`
   - **Email**: `admin@gamegear.local`
   - **Password**: `admin123` (หรือตามต้องการ)

3. เพิ่ม Kong Connection:

   - **Name**: `GameGear Kong`
   - **Kong Admin URL**: `http://kong:8001` (ถ้ารันใน Docker network)
     หรือ `http://host.docker.internal:8001` (ถ้ารัน Kong แยก)

4. คลิก **Activate** เพื่อเชื่อมต่อ

---

## 🔧 ตั้งค่า Services และ Routes

> 🎯 **จุดประสงค์**: กำหนดค่า Services และ Routes ใน Konga UI เพื่อให้ Kong Gateway รู้จักและนำทางคำขอไปยัง Microservices ที่ถูกต้อง

### ขั้นตอนการเพิ่ม Service ใหม่

### 1. เพิ่ม Users Service

1. ใน Konga UI ไปที่ **Services** → **Add New Service**
2. กรอกข้อมูล:
   ```
   Name: users-service
   Protocol: https
   Host: abc123.ngrok.io  (URL จาก ngrok ของ ณิชพน)
   Port: 443
   Path: /
   ```
3. คลิก **Submit**

### 2. เพิ่ม Route สำหรับ Users Service

1. ไปที่ **Routes** ของ `users-service` → **Add Route**
2. กรอกข้อมูล:
   ```
   Name: users-route
   Paths: /users
   Strip Path: false
   ```
3. คลิก **Submit**

### 3. ทำซ้ำสำหรับ Shop และ Admin Service

**Shop Service:**

```
Service Name: shop-service
Host: def456.ngrok.io  (URL จาก ngrok ของ ณัฐพงษ์/วายุ)
Protocol: https
Port: 443
Route Path: /shop
```

**Admin Service:**

```
Service Name: admin-service
Host: ghi789.ngrok.io  (URL จาก ngrok ของ วรรธนโรจน์)
Protocol: https
Port: 443
Route Path: /admin
```

### ทดสอบการเชื่อมต่อ

```bash
# ทดสอบผ่าน Kong Gateway
curl http://localhost:8000/users/healthz
curl http://localhost:8000/shop/healthz
curl http://localhost:8000/admin/healthz
```

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
   ```
   minute: 100
   hour: 10000
   policy: local
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

## 🐋 Health Check Endpoints

แต่ละ Service ควรมี `/healthz` endpoint สำหรับตรวจสอบสถานะ:

### ตัวอย่าง Implementation (Gin)

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()

    // Healthcheck endpoint
    r.GET("/healthz", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "service": "users-service",
            "version": "1.0.0",
        })
    })

    // ... other routes

    r.Run(":8081")
}
```

### ทดสอบผ่าน Kong Gateway

```bash
# ตรวจสอบ Service ทั้งหมด
curl http://localhost:8000/users/healthz
curl http://localhost:8000/shop/healthz
curl http://localhost:8000/admin/healthz

# Expected Response
{
  "status": "ok",
  "service": "users-service",
  "version": "1.0.0"
}
```

---

## 🔍 แก้ไขปัญหา (Troubleshooting)

### Problem: Kong ไม่สามารถเชื่อมต่อกับ Services

**Symptoms:**

- Error: `502 Bad Gateway`
- Kong logs: `upstream connect error or disconnect/reset before headers`

**Solutions:**

1. **ตรวจสอบว่า Services ทำงานอยู่:**

   ```bash
   curl http://localhost:8081/healthz  # Users
   curl http://localhost:8082/healthz  # Shop
   curl http://localhost:8083/healthz  # Admin
   ```

2. **ตรวจสอบ Docker Network:**

   ```bash
   docker network ls
   docker network inspect kong-net
   ```

3. **แก้ไข Host Configuration ใน Konga:**
   - ถ้ารัน Services นอก Docker: ใช้ `host.docker.internal`
   - ถ้ารัน Services ใน Docker network เดียวกัน: ใช้ `users-service`, `shop-service`, etc.

### Problem: Konga ไม่สามารถเชื่อมต่อกับ Kong

**Solutions:**

1. **ตรวจสอบว่า Kong Admin API ทำงาน:**

   ```bash
   curl http://localhost:8001/
   ```

2. **ใช้ Kong Admin URL ที่ถูกต้อง:**
   - ถ้ารันใน Docker network: `http://kong:8001`
   - ถ้ารันจากเครื่องโฮสต์: `http://localhost:8001`

### Problem: CORS Errors ใน Browser

**Solutions:**

1. **เพิ่ม CORS Plugin แบบ Global:**

   - ไปที่ Konga → Plugins → Add Global Plugin → CORS
   - ตั้งค่า `origins: *` (Dev) หรือระบุ domain (Production)

2. **ตรวจสอบ Headers:**
   ```bash
   curl -X OPTIONS http://localhost:8000/users/login \
     -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: POST" \
     -v
   ```

### Problem: JWT Authentication ไม่ทำงาน

**Solutions:**

1. **ตรวจสอบว่า JWT Plugin ติดตั้งแล้ว:**

   - ไปที่ Route/Service ใน Konga → Plugins → ดู JWT

2. **ตรวจสอบ JWT Token Format:**

   ```bash
   # Token ต้องส่งใน Authorization header
   curl http://localhost:8000/users/profile \
     -H "Authorization: Bearer <your-jwt-token>"
   ```

3. **ตรวจสอบ Secret Key ตรงกันระหว่าง Service และ Kong**

### Problem: Database Migration Errors

**Solutions:**

```bash
# รัน migrations ใหม่
docker compose -f admin-service/docker-compose.kong.yml down -v
docker compose -f admin-service/docker-compose.kong.yml up -d

# ตรวจสอบ logs
docker logs kong-migrations
docker logs kong-gateway
```

### Problem: Port Already in Use

**Solutions:**

```bash
# ตรวจสอบ port ที่ใช้อยู่
netstat -ano | findstr :8000  # Windows
lsof -i :8000                  # Mac/Linux

# หยุด process ที่ใช้ port หรือเปลี่ยน port ใน docker-compose
```

### Useful Commands

```bash
# ดู logs ทั้งหมด
docker compose -f admin-service/docker-compose.kong.yml logs -f

# ดู logs เฉพาะ service
docker logs kong-gateway -f
docker logs konga -f

# Restart services
docker compose -f admin-service/docker-compose.kong.yml restart

# ลบทั้งหมดและเริ่มใหม่ (ระวัง: จะลบข้อมูลทั้งหมด)
docker compose -f admin-service/docker-compose.kong.yml down -v
docker compose -f admin-service/docker-compose.kong.yml up -d

# ตรวจสอบสถานะ containers
docker ps

# ตรวจสอบ Kong configuration
curl http://localhost:8001/services
curl http://localhost:8001/routes
curl http://localhost:8001/plugins
```

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

## ✅ สรุปคุณสมบัติ (Features Summary)

### What You Get

- ✅ **Kong Gateway (DB Mode)** - API Gateway หลักพร้อม PostgreSQL backing
- ✅ **Konga Admin UI** - จัดการ Kong แบบ GUI สะดวก ไม่ต้องใช้ Command Line
- ✅ **3 Microservices** - Users, Shop, Admin แยกกันชัดเจน
- ✅ **Mermaid Diagrams** - Architecture และ Flow diagrams ในเอกสาร
- ✅ **Docker Compose Ready** - รันทั้งระบบด้วยคำสั่งเดียว
- ✅ **Healthcheck Endpoints** - ตรวจสอบสถานะ service ได้ทันที
- ✅ **Plugin Support** - CORS, JWT, Rate Limiting พร้อมใช้งาน
- ✅ **Troubleshooting Guide** - มีคำแนะนำแก้ปัญหาที่พบบ่อย

### Quick Commands Recap

```bash
# 1. เริ่ม Kong Gateway + Konga
# จาก root directory (GameGear-Ecommerce/)
docker compose -f admin-service/docker-compose.kong.yml up -d

# 2. เปิด Konga UI
open http://localhost:1337

# 3. ตั้งค่า Services และ Routes (ผ่าน Konga UI)

# 4. ทดสอบผ่าน Kong Gateway
curl http://localhost:8000/users/healthz
curl http://localhost:8000/shop/healthz
curl http://localhost:8000/admin/healthz
```

### Documentation Structure

```
GameGear-Ecommerce/
├── Mini-Project-Golang/ (Main Repository)
│   └── README.md (Main - คุณอยู่ที่นี่)
├── users-service/ (Users Service Repository)
│   ├── README.md (Detailed setup for Users Service)
│   ├── go.mod
│   ├── go.sum
│   ├── .env.example
│   └── internal/ (handlers, models, repositories, services)
├── shop-service/ (Shop Service Repository)
│   ├── README.md (Detailed setup for Shop Service)
│   ├── go.mod
│   ├── go.sum
│   ├── .env.example
│   └── internal/ (handlers, models, repositories, services)
└── admin-service/ (Admin Service Repository)
    ├── README.md (Detailed setup for Admin Service)
    ├── docker-compose.kong.yml (Kong Gateway + Konga UI)
    ├── go.mod
    ├── go.sum
    ├── .env.example
    └── internal/ (handlers, models, services, clients)
```

> 💡 **หมายเหตุ**: แต่ละ service เป็น repository แยกกันบน GitHub

### Learning Resources Included

- **Architecture Diagrams**: Mermaid diagrams แสดงโครงสร้างระบบและ Request Flow
- **Setup Guides**: คำแนะนำทีละขั้นตอนสำหรับ Kong และ Konga
- **Plugin Configuration**: ตัวอย่างการตั้งค่า CORS, JWT, Rate Limiting
- **API Testing**: คำแนะนำการทดสอบ API ด้วย Postman
- **Troubleshooting**: วิธีแก้ปัญหาที่พบบ่อยพร้อมคำสั่ง
- **Best Practices**: แนวทางที่ดีสำหรับ Development และ Production

---

## 📞 ติดต่อและสนับสนุน

หากมีปัญหาหรือคำถาม:

1. ตรวจสอบ **Troubleshooting Section** ด้านบน
2. ดู **README.md** ของแต่ละ Service สำหรับรายละเอียดเพิ่มเติม
3. ตรวจสอบ logs: `docker logs kong-gateway -f`
4. ติดต่อทีมพัฒนา (ดูรายชื่อด้านบน)

---

**Made with ❤️ by GameGear Development Team**
