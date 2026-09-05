# Getting Started Guide

คู่มือนี้สรุปขั้นตอนการเตรียมสภาพแวดล้อม การติดตั้ง และการรันระบบ **GameGear E-commerce** ทั้ง Kong Gateway, Konga UI และ Microservices ทั้งสามตัว

---

## 1. Prerequisites

ก่อนเริ่มติดตั้งระบบ ตรวจสอบให้แน่ใจว่าเครื่องของคุณติดตั้งเครื่องมือต่อไปนี้เรียบร้อยแล้ว:

- **Git**: เวอร์ชัน 2.30+
- **Go**: เวอร์ชัน 1.25+
- **Docker & Docker Compose**: Docker Desktop หรือ Docker Engine 24.0+
- **Ports ที่ต้องไม่ถูกใช้งานล่วงหน้า**: `8000`, `8001`, `1337`, `8081`, `8082`, `8083`

---

## 2. Infrastructure Setup (Kong Gateway & Konga UI)

Kong Gateway และ Konga UI ถูกกำหนดค่าผ่านไฟล์ Docker Compose ซึ่งอยู่ที่โฟลเดอร์ `admin-service`:

```bash
# นำทางไปยังโฟลเดอร์ admin-service
cd admin-service

# เริ่มต้นการทำงานของ Kong Database, Migrations, Kong Gateway และ Konga UI
docker compose -f docker-compose.kong.yml up -d
```

### การตรวจสอบสถานะ Containers
ตรวจสอบว่าทุก Container ทำงานปกติและ Healthy:

```bash
docker compose -f docker-compose.kong.yml ps
```

ตรวจสอบการตอบสนองของ Gateway:
```bash
# ตรวจสอบ Kong Admin API
curl -i http://localhost:8001/

# ตรวจสอบ Kong Proxy
curl -i http://localhost:8000/

# ตรวจสอบ Konga UI
curl -i http://localhost:1337/
```

> รายละเอียดขั้นตอนการผูก Service และ Plugin ใน Konga ดูได้ที่ [docs/gateway-configuration.md](gateway-configuration.md)

---

## 3. Running Microservices Locally

โปรเจกต์ประกอบด้วย 3 บริการหลัก ซึ่งสามารถรันแบบ Local แยกตาม Terminal ได้ดังนี้:

### 3.1 Users Service (Port 8081)
รับผิดชอบระบบยืนยันตัวตน สมาชิก และการจัดการ JWT

```bash
cd users-service

# ติดตั้ง Dependencies
go mod tidy

# คัดลอกและตั้งค่า Environment
cp .env.example .env

# รัน Service
go run cmd/api/main.go
```

### 3.2 Shop Service (Port 8082)
รับผิดชอบระบบสินค้า ตะกร้าสินค้า และคำสั่งซื้อ

```bash
cd shop-service

# ติดตั้ง Dependencies
go mod tidy

# คัดลอกและตั้งค่า Environment
cp .env.example .env

# รัน Service
go run cmd/api/main.go
```

### 3.3 Admin Service (Port 8083)
รับผิดชอบระบบแดชบอร์ดหลังบ้านและการประสานงานข้ามบริการ

```bash
cd admin-service

# ติดตั้ง Dependencies
go mod tidy

# คัดลอกและตั้งค่า Environment
cp .env.example .env

# รัน Service
go run cmd/api/main.go
```

---

## 4. Healthcheck Verification

เมื่อเปิดทั้ง Gateway และ Microservices ครบแล้ว สามารถตรวจสอบความพร้อมของระบบผ่าน Kong Proxy (Port 8000):

```bash
# ตรวจสอบ Users Service ผ่าน Gateway
curl -sS http://localhost:8000/users/healthz

# ตรวจสอบ Shop Service ผ่าน Gateway
curl -sS http://localhost:8000/shop/healthz

# ตรวจสอบ Admin Service ผ่าน Gateway
curl -sS http://localhost:8000/admin/healthz
```

ทุก Endpoint ควรตอบกลับสถานะ `200 OK`

---

## 5. Team Collaboration with ngrok (Optional)

กรณีที่สมาชิกในทีมแยกกันพัฒนาคนละ Service บนเครื่องของตนเอง สามารถเปิดอุโมงค์เครือข่ายด้วย ngrok เพื่อให้ผู้ดูแล Gateway เชื่อมต่อเข้ามาได้:

1. **ติดตั้งและ Authenticate ngrok**:
   ```bash
   ngrok config add-authtoken <YOUR_AUTHTOKEN>
   ```
2. **เปิดอุโมงค์ตาม Port ของ Service**:
   ```bash
   # ตัวอย่าง: เปิดอุโมงค์สำหรับ Users Service (Port 8081)
   ngrok http 8081
   ```
3. **ส่ง Forwarding URL**:
   คัดลอก Public URL ที่ได้ (เช่น `https://xxxx.ngrok-free.app`) ส่งให้ผู้ดูแล Kong Gateway เพื่อนำไปอัปเดต Service Upstream Host ใน Konga
4. **ทดสอบ URL ผ่าน ngrok**:
   ```bash
   curl -sS https://xxxx.ngrok-free.app/healthz
   ```
