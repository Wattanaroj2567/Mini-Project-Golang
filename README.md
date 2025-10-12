# 🚀 Mini-Project: GameGear E-commerce (Microservice Architecture + Kong API Gateway)

![Go](https://img.shields.io/badge/Go-1.25.1-00ADD8?style=for-the-badge&logo=go)
![Gin](https://img.shields.io/badge/Gin-Framework-008ECF?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql)
![GORM](https://img.shields.io/badge/GORM-B93527?style=for-the-badge)
![JWT](https://img.shields.io/badge/Auth-JWT-FF6F00?style=for-the-badge)
![Swagger](https://img.shields.io/badge/API-Swagger-85EA2D?style=for-the-badge&logo=swagger)
![Kong](https://img.shields.io/badge/API%20Gateway-Kong-003459?style=for-the-badge&logo=kong)
![Konga](https://img.shields.io/badge/Admin%20UI-Konga-00C7B7?style=for-the-badge)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker)

โปรเจกต์นี้คือระบบ **E-commerce** สำหรับอุปกรณ์เกมมิ่ง บนสถาปัตยกรรม **Microservices** โดยมี **Kong API Gateway (DB Mode)** เป็นจุดทางเข้าเดียว (single entry point) พร้อมกับ **Konga UI** สำหรับจัดการ Gateway แบบ GUI เพื่อควบคุม Routing, Authentication (JWT), Rate Limiting, CORS, และ Logging ให้กับบริการภายใน (`users-service`, `shop-service`, `admin-service`).

---

## 📋 Table of Contents

- [System Architecture Overview](#system-architecture-overview)
- [Request Flow](#request-flow)
- [Kong API Gateway](#kong-api-gateway)
  - [Why Kong Gateway?](#why-kong-gateway)
  - [DB Mode vs DB-less Mode](#db-mode-vs-db-less-mode)
- [Quick Start: ติดตั้งและรัน Kong + Konga](#quick-start-ติดตั้งและรัน-kong--konga)
  - [Prerequisites](#prerequisites)
  - [1️⃣ เริ่ม Kong Gateway + Konga](#1-เริ่ม-kong-gateway--konga)
  - [2️⃣ รัน Services (สำหรับ Developers)](#2-รัน-services-สำหรับ-developers)
  - [3️⃣ ตรวจสอบ Services](#3-ตรวจสอบ-services)
  - [4️⃣ เปิด Konga UI และตั้งค่าครั้งแรก](#4-เปิด-konga-ui-และตั้งค่าครั้งแรก)
  - [5️⃣ ตั้งค่า Services และ Routes ใน Konga](#5-ตั้งค่า-services-และ-routes-ใน-konga)
- [การตั้งค่า Services และ Routes ใน Konga](#การตั้งค่า-services-และ-routes-ใน-konga)
  - [ขั้นตอนการเพิ่ม Service ใหม่](#ขั้นตอนการเพิ่ม-service-ใหม่)
  - [ทดสอบการเชื่อมต่อ](#ทดสอบการเชื่อมต่อ)
- [Adding Essential Plugins](#adding-essential-plugins)
  - [CORS Plugin](#cors-plugin)
  - [JWT Plugin](#jwt-plugin)
  - [Rate Limiting Plugin](#rate-limiting-plugin)
- [โครงสร้าง Docker Compose (admin-service/docker-compose.kong.yml)](#โครงสร้าง-docker-compose-admin-servicedocker-composekongyml)
- [Ports Summary](#ports-summary)
- [Healthcheck Endpoints](#healthcheck-endpoints)
- [Troubleshooting](#troubleshooting)
  - [❌ Problem: Kong Not Starting](#-problem-kong-not-starting)
  - [❌ Problem: Konga Connection Failed](#-problem-konga-connection-failed)
  - [❌ Problem: Database Migration Errors](#-problem-database-migration-errors)
  - [❌ Problem: Port Already in Use](#-problem-port-already-in-use)
  - [📋 Useful Commands](#-useful-commands)
- [Additional Resources](#additional-resources)
  - [Project Documentation (Google Docs)](#project-documentation-google-docs)
- [Development Team](#development-team)
- [Service Repositories](#service-repositories)
- [Summary & Key Features](#summary--key-features)
  - [Quick Commands Recap](#quick-commands-recap)
  - [Documentation Structure](#documentation-structure)
  - [Learning Resources Included](#learning-resources-included)
- [Support & Contact](#support--contact)

---

## 🏛️ System Architecture Overview

```mermaid
graph TB
    subgraph "Client Layer"
        WEB[🌐 Web Application]
        MOBILE[📱 Mobile App]
        ADMIN_UI[🛡️ Admin Dashboard]
    end

    subgraph "API Gateway Layer"
        KONGA[🖥️ Konga Admin UI<br/>:1337]
        KONG[🦍 Kong Gateway<br/>Proxy :8888<br/>Admin API :8001]
        KONG_DB[(🗄️ Kong PostgreSQL<br/>Configuration DB)]
        KONGA_DB[(🗄️ Konga PostgreSQL<br/>UI Database)]
    end

    subgraph "Microservices Layer"
        USERS[👤 Users Service<br/>:8080<br/>Auth & User Management]
        SHOP[🛍️ Shop Service<br/>:8081<br/>Products, Cart, Orders]
        ADMIN[🛡️ Admin Service<br/>:8082<br/>Backend Admin Panel]

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

### 🔑 องค์ประกอบหลัก

| Component            | Description                                                              | Port                              | Technology         |
| -------------------- | ------------------------------------------------------------------------ | --------------------------------- | ------------------ |
| **🦍 Kong Gateway**  | API Gateway หลัก - จัดการ Routing, JWT Validation, Rate Limiting, CORS   | 8000 (Proxy)<br/>8001 (Admin API) | Kong 3.4 (DB Mode) |
| **🖥️ Konga UI**      | Web-based Admin Dashboard สำหรับจัดการ Kong แบบ GUI                      | 1337                              | Konga (Node.js)    |
| **👤 Users Service** | จัดการผู้ใช้และการยืนยันตัวตน (Register, Login, Profile, Password Reset) | 8080                              | Go + Gin + GORM    |
| **🛍️ Shop Service**  | จัดการสินค้า, ตะกร้า, คำสั่งซื้อ                                         | 8081                              | Go + Gin + GORM    |
| **🛡️ Admin Service** | ระบบหลังบ้าน (Admin Panel) - ประสานงานกับ Users/Shop ผ่าน Kong           | 8082                              | Go + Gin           |
| **🗄️ PostgreSQL**    | ฐานข้อมูลสำหรับแต่ละ Service + Kong + Konga                              | 5432 (varies)                     | PostgreSQL 15/17   |

### 📊 Request Flow Example

```mermaid
sequenceDiagram
    participant Client
    participant Kong
    participant UsersService
    participant ShopService

    Client->>Kong: POST /users/login<br/>{email, password}
    Kong->>UsersService: Forward to :8080/login
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

## 📘 API Documentation (OpenAPI / Swagger)

ทุก Service มี Swagger UI สำหรับอ้างอิงและทดสอบ API:

| Service       | Direct Access (Dev Only)                                                      | Via Kong Gateway (Production)                    | Description                |
| ------------- | ----------------------------------------------------------------------------- | ------------------------------------------------ | -------------------------- |
| Users Service | [localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) | `http://localhost:8000/users/swagger/index.html` | สมัคร, ล็อกอิน, โปรไฟล์    |
| Shop Service  | [localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html) | `http://localhost:8000/shop/swagger/index.html`  | สินค้า, ตะกร้า, คำสั่งซื้อ |
| Admin Service | [localhost:8082/swagger/index.html](http://localhost:8082/swagger/index.html) | `http://localhost:8000/admin/swagger/index.html` | ฟังก์ชันผู้ดูแลระบบ        |

> 💡 **Best Practice**: ใน Production ควรปิด Direct Access และเรียกผ่าน Kong Gateway เท่านั้น

---

## 🦍 Kong Gateway + Konga Setup (DB Mode)

โปรเจกต์นี้ใช้ **Kong Gateway (DB Mode)** พร้อม **Konga Admin UI** เพื่อให้สามารถจัดการ Services, Routes, และ Plugins ผ่าน Web Interface ได้แบบ Real-time

### 🎯 ทำไมเลือก DB Mode + Konga?

| Feature                  | DB Mode + Konga ✅          | DB-less Mode ❌              |
| ------------------------ | --------------------------- | ---------------------------- |
| **GUI Management**       | มี Konga UI ใช้งานง่าย      | ไม่มี (ต้องแก้ YAML เอง)     |
| **Real-time Changes**    | แก้ได้ทันที ไม่ต้อง restart | ต้อง reload config ทุกครั้ง  |
| **Team Collaboration**   | หลายคนจัดการได้พร้อมกัน     | ต้อง commit YAML แล้ว deploy |
| **Plugin Configuration** | ตั้งค่าผ่าน UI สะดวก        | ต้องเขียน config เอง         |
| **Suitable For**         | Development & Production    | CI/CD, Static Config         |

---

## 🚀 Quick Start: ติดตั้งและรัน Kong + Konga

> 👤 **ผู้รับผิดชอบ Kong Gateway Setup**: **วรรธนโรจน์ บุตรดี** (Project Manager)  
> 📝 หากมีปัญหาเกี่ยวกับ Kong Gateway, Konga UI หรือการตั้งค่า Routes/Plugins กรุณาติดต่อ

### 📋 Prerequisites

- Docker และ Docker Compose ติดตั้งแล้ว
- Port ว่าง: `8000`, `8001`, `1337`, `8080`, `8081`, `8082`

### 1️⃣ เริ่ม Kong Gateway + Konga

ใช้ไฟล์ `docker-compose.kong.yml` ที่มีอยู่ใน admin-service directory:

```bash
# จาก root directory (GameGear-Ecommerce/)
cd admin-service
docker compose -f docker-compose.kong.yml up -d
```

**หมายเหตุ:** Kong Gateway จัดการโดย PM (วรรธนโรจน์) ใน admin-service เท่านั้น

### 2️⃣ รัน Services (สำหรับ Developers)

**สำหรับ Users Service (ณิชพน):**

```bash
cd users-service
go run cmd/api/main.go

# Terminal ใหม่
ngrok http 8080
# แชร์ URL: https://abc123.ngrok.io
```

**สำหรับ Shop Service (ณัฐพงษ์ + วายุ):**

```bash
cd shop-service
go run cmd/api/main.go

# Terminal ใหม่
ngrok http 8081
# แชร์ URL: https://def456.ngrok.io
```

**สำหรับ Admin Service (วรรธนโรจน์):**

```bash
cd admin-service
go run cmd/api/main.go

# Terminal ใหม่
ngrok http 8082
# แชร์ URL: https://ghi789.ngrok.io
```

### 3️⃣ ตรวจสอบ Services

```bash
# ตรวจสอบว่า Kong ทำงาน
curl -i http://localhost:8001/

# ตรวจสอบ Konga UI
curl -i http://localhost:1337/
```

### 4️⃣ เปิด Konga UI และตั้งค่าครั้งแรก

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

## 🔧 การตั้งค่า Services และ Routes ใน Konga

### ขั้นตอนการเพิ่ม Service ใหม่

#### 1. เพิ่ม Users Service

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

#### 2. เพิ่ม Route สำหรับ Users Service

1. ไปที่ **Routes** ของ `users-service` → **Add Route**
2. กรอกข้อมูล:
   ```
   Name: users-route
   Paths: /users
   Strip Path: false
   ```
3. คลิก **Submit**

#### 3. ทำซ้ำสำหรับ Shop และ Admin Service

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

## 🔌 เพิ่ม Plugins ที่จำเป็น

### 1. CORS Plugin (สำหรับ Frontend)

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

### 2. JWT Plugin (สำหรับ Authentication)

สำหรับ Routes ที่ต้องการ Authentication:

1. ไปที่ Service/Route ที่ต้องการป้องกัน
2. เพิ่ม **JWT Plugin**
3. ตั้งค่า:
   ```
   key_claim_name: iss
   secret_is_base64: false
   ```

### 3. Rate Limiting Plugin (ป้องกัน DDoS)

1. เพิ่ม **Rate Limiting Plugin** แบบ Global หรือต่อ Service
2. ตั้งค่า:
   ```
   minute: 100
   hour: 10000
   policy: local
   ```

---

## 📁 โครงสร้าง Docker Compose (admin-service/docker-compose.kong.yml)

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

## 🌐 Ports Summary

| Service               | Port | Access URL            | Description                        |
| --------------------- | ---- | --------------------- | ---------------------------------- |
| 🦍 **Kong Proxy**     | 8000 | http://localhost:8000 | API Gateway - ใช้เรียก API ทั้งหมด |
| 🔧 **Kong Admin API** | 8001 | http://localhost:8001 | Admin API (Dev only)               |
| 🖥️ **Konga UI**       | 1337 | http://localhost:1337 | Web Dashboard จัดการ Kong          |
| 👤 **Users Service**  | 8080 | http://localhost:8080 | Direct access (Dev only)           |
| 🛍️ **Shop Service**   | 8081 | http://localhost:8081 | Direct access (Dev only)           |
| 🛡️ **Admin Service**  | 8082 | http://localhost:8082 | Direct access (Dev only)           |

> ⚠️ **Production**: ควรปิด Direct Access และให้เรียกผ่าน Kong Proxy (port 8000) เท่านั้น

---

## 🐋 Healthcheck Endpoints

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

    r.Run(":8080")
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

## 🔍 Troubleshooting

### ❌ Problem: Kong ไม่สามารถเชื่อมต่อกับ Services

**Symptoms:**

- Error: `502 Bad Gateway`
- Kong logs: `upstream connect error or disconnect/reset before headers`

**Solutions:**

1. **ตรวจสอบว่า Services ทำงานอยู่:**

   ```bash
   curl http://localhost:8080/healthz  # Users
   curl http://localhost:8081/healthz  # Shop
   curl http://localhost:8082/healthz  # Admin
   ```

2. **ตรวจสอบ Docker Network:**

   ```bash
   docker network ls
   docker network inspect kong-net
   ```

3. **แก้ไข Host Configuration ใน Konga:**
   - ถ้ารัน Services นอก Docker: ใช้ `host.docker.internal`
   - ถ้ารัน Services ใน Docker network เดียวกัน: ใช้ `users-service`, `shop-service`, etc.

### ❌ Problem: Konga ไม่สามารถเชื่อมต่อกับ Kong

**Solutions:**

1. **ตรวจสอบว่า Kong Admin API ทำงาน:**

   ```bash
   curl http://localhost:8001/
   ```

2. **ใช้ Kong Admin URL ที่ถูกต้อง:**
   - ถ้ารันใน Docker network: `http://kong:8001`
   - ถ้ารันจากเครื่องโฮสต์: `http://localhost:8001`

### ❌ Problem: CORS Errors ใน Browser

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

### ❌ Problem: JWT Authentication ไม่ทำงาน

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

### ❌ Problem: Database Migration Errors

**Solutions:**

```bash
# รัน migrations ใหม่
docker compose -f admin-service/docker-compose.kong.yml down -v
docker compose -f admin-service/docker-compose.kong.yml up -d

# ตรวจสอบ logs
docker logs kong-migrations
docker logs kong-gateway
```

### ❌ Problem: Port Already in Use

**Solutions:**

```bash
# ตรวจสอบ port ที่ใช้อยู่
netstat -ano | findstr :8000  # Windows
lsof -i :8000                  # Mac/Linux

# หยุด process ที่ใช้ port หรือเปลี่ยน port ใน docker-compose
```

### 📋 Useful Commands

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

## 📚 Additional Resources

### 📖 Project Documentation (Google Docs)

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
- [Swag (Swagger for Go)](https://github.com/swaggo/swag)

---

## 📂 Service Repositories

| Service Repository                                                          | Description                                             | Team Member                |
| --------------------------------------------------------------------------- | ------------------------------------------------------- | -------------------------- |
| 👤 **[Users Service](https://github.com/Wattanaroj2567/users-service.git)** | จัดการผู้ใช้และการยืนยันตัวตน (สมัคร, ล็อกอิน, โปรไฟล์) | ณิชพน มานิตย์              |
| 🛍️ **[Shop Service](https://github.com/Wattanaroj2567/shop-service.git)**   | จัดการสินค้า, ตะกร้า, คำสั่งซื้อ                        | ณัฐพงษ์ ดีบุตร, วายุ กอคูณ |
| 🛡️ **[Admin Service](https://github.com/Wattanaroj2567/admin-service.git)** | ระบบหลังบ้าน (Admin Panel)                              | วรรธนโรจน์ บุตรดี          |

---

## 🤝 Development Team

| Profile                                                                                                                       | Name                  | Responsibility            |
| ----------------------------------------------------------------------------------------------------------------------------- | --------------------- | ------------------------- |
| [<img src="https://github.com/Wattanaroj2567.png" width="60" height="60"/>](https://github.com/Wattanaroj2567)                | **วรรธนโรจน์ บุตรดี** | Project Manager & Backend |
| [<img src="https://avatars.githubusercontent.com/u/159878532?v=4" width="60" height="60"/>](https://github.com/Natthaphong66) | **ณัฐพงษ์ ดีบุตร**    | Backend Developer         |
| [<img src="https://avatars.githubusercontent.com/u/159880199?v=4" width="60" height="60"/>](https://github.com/nitchapon66)   | **ณิชพน มานิตย์**     | Backend Developer         |
| [<img src="https://avatars.githubusercontent.com/u/160033040?v=4" width="60" height="60"/>](https://github.com/FUJIKOTH)      | **วายุ กอคูณ**        | Backend Developer         |

---

## ✅ Summary & Key Features

### 🎯 What You Get

- ✅ **Kong Gateway (DB Mode)** - API Gateway หลักพร้อม PostgreSQL backing
- ✅ **Konga Admin UI** - จัดการ Kong แบบ GUI สะดวก ไม่ต้องใช้ Command Line
- ✅ **3 Microservices** - Users, Shop, Admin แยกกันชัดเจน
- ✅ **Swagger Documentation** - ทุก service มีเอกสาร API แบบ Interactive
- ✅ **Mermaid Diagrams** - Architecture และ Flow diagrams ในเอกสาร
- ✅ **Docker Compose Ready** - รันทั้งระบบด้วยคำสั่งเดียว
- ✅ **Healthcheck Endpoints** - ตรวจสอบสถานะ service ได้ทันที
- ✅ **Plugin Support** - CORS, JWT, Rate Limiting พร้อมใช้งาน
- ✅ **Troubleshooting Guide** - มีคำแนะนำแก้ปัญหาที่พบบ่อย

### 🚀 Quick Commands Recap

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

# 5. ดู Swagger Documentation
open http://localhost:8080/swagger/index.html  # Users
open http://localhost:8081/swagger/index.html  # Shop
open http://localhost:8082/swagger/index.html  # Admin
```

### 📖 Documentation Structure

```
GameGear-Ecommerce/
├── Mini-Project-Golang/ (Main Repository)
│   └── README.md (Main - คุณอยู่ที่นี่)
├── users-service/ (Users Service Repository)
│   ├── README.md (Detailed setup for Users Service)
│   ├── go.mod
│   ├── go.sum
│   └── internal/ (handlers, models, repositories, services)
├── shop-service/ (Shop Service Repository)
│   ├── README.md (Detailed setup for Shop Service)
│   ├── go.mod
│   ├── go.sum
│   └── internal/ (handlers, models, repositories, services)
└── admin-service/ (Admin Service Repository)
    ├── README.md (Detailed setup for Admin Service)
    ├── docker-compose.kong.yml (Kong Gateway + Konga UI)
    ├── go.mod
    ├── go.sum
    └── internal/ (handlers, models, repositories, services)
```

> 💡 **หมายเหตุ**: แต่ละ service เป็น repository แยกกันบน GitHub

### 🎓 Learning Resources Included

- **Architecture Diagrams**: Mermaid diagrams แสดงโครงสร้างระบบและ Request Flow
- **Setup Guides**: คำแนะนำทีละขั้นตอนสำหรับ Kong และ Konga
- **Plugin Configuration**: ตัวอย่างการตั้งค่า CORS, JWT, Rate Limiting
- **Troubleshooting**: วิธีแก้ปัญหาที่พบบ่อยพร้อมคำสั่ง
- **Best Practices**: แนวทางที่ดีสำหรับ Development และ Production

---

## 👥 Development Team

### Team Members

| Role                  | Name              | Service                      | GitHub                                               |
| --------------------- | ----------------- | ---------------------------- | ---------------------------------------------------- |
| **Project Manager**   | วรรธนโรจน์ บุตรดี | admin-service                | [@Wattanaroj2567](https://github.com/Wattanaroj2567) |
| **Backend Developer** | ณิชพน มานิตย์     | users-service                | [@Wattanaroj2567](https://github.com/Wattanaroj2567) |
| **Backend Developer** | ณัฐพงษ์ ดีบุตร    | shop-service (Products)      | [@Wattanaroj2567](https://github.com/Wattanaroj2567) |
| **Backend Developer** | วายุ กอคูณ        | shop-service (Cart & Orders) | [@Wattanaroj2567](https://github.com/Wattanaroj2567) |

### Responsibilities

- **วรรธนโรจน์ บุตรดี**: Kong Gateway management, Admin service, Project coordination
- **ณิชพน มานิตย์**: Users service, Authentication & Authorization
- **ณัฐพงษ์ ดีบุตร**: Product catalog, Category management
- **วายุ กอคูณ**: Cart management, Order processing

---

## 📚 Service Repositories

### Individual Service Repositories

| Service                | Repository                                                                   | Description                                       | Developer                  |
| ---------------------- | ---------------------------------------------------------------------------- | ------------------------------------------------- | -------------------------- |
| **Main Documentation** | [Mini-Project-Golang](https://github.com/Wattanaroj2567/Mini-Project-Golang) | Main project documentation and Kong Gateway setup | วรรธนโรจน์ บุตรดี          |
| **Users Service**      | [users-service](https://github.com/Wattanaroj2567/users-service)             | Authentication, user management, JWT tokens       | ณิชพน มานิตย์              |
| **Shop Service**       | [shop-service](https://github.com/Wattanaroj2567/shop-service)               | Product catalog, cart, orders                     | ณัฐพงษ์ ดีบุตร, วายุ กอคูณ |
| **Admin Service**      | [admin-service](https://github.com/Wattanaroj2567/admin-service)             | Admin dashboard, Kong Gateway management          | วรรธนโรจน์ บุตรดี          |

### Repository Structure

```
GameGear-Ecommerce/
├── Mini-Project-Golang/ (Main Repository)
│   └── README.md (Main documentation)
├── users-service/ (Users Service Repository)
│   ├── README.md (Service documentation)
│   ├── go.mod
│   ├── go.sum
│   └── internal/ (handlers, models, repositories, services)
├── shop-service/ (Shop Service Repository)
│   ├── README.md (Service documentation)
│   ├── go.mod
│   ├── go.sum
│   └── internal/ (handlers, models, repositories, services)
└── admin-service/ (Admin Service Repository)
    ├── README.md (Service documentation)
    ├── docker-compose.kong.yml (Kong Gateway + Konga UI)
    ├── go.mod
    ├── go.sum
    └── internal/ (handlers, models, repositories, services)
```

> 💡 **หมายเหตุ**: แต่ละ service เป็น repository แยกกันบน GitHub

---

## 📞 Support & Contact

หากมีปัญหาหรือคำถาม:

1. ตรวจสอบ **Troubleshooting Section** ด้านบน
2. ดู **README.md** ของแต่ละ Service สำหรับรายละเอียดเพิ่มเติม
3. ตรวจสอบ logs: `docker logs kong-gateway -f`
4. ติดต่อทีมพัฒนา (ดูรายชื่อด้านบน)

---

**Made with ❤️ by GameGear Development Team**
