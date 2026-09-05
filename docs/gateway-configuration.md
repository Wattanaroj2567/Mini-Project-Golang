# Kong Gateway & Konga Configuration Guide

เอกสารนี้อธิบายการตั้งค่า **Kong API Gateway** และการควบคุมผ่าน **Konga Admin UI** สำหรับการเชื่อมต่อ Microservices, การกำหนดเส้นทางคำขอ (Routes) และการเปิดใช้งาน Plugins ด้านความปลอดภัย

---

## 1. Connecting Konga to Kong Admin API

1. เปิดเบราว์เซอร์ไปยัง `http://localhost:1337`
2. สมัครบัญชีผู้ดูแลระบบของ Konga ในการเข้าใช้งานครั้งแรก
3. ไปที่เมนู **Connections** -> คลิก **Create Connection**
4. กำหนดค่าการเชื่อมต่อ:
   - **Name**: `GameGear Kong`
   - **Kong Admin URL**:
     - `http://kong:8001` (หาก Konga และ Kong อยู่ใน Docker Network เดียวกัน)
     - หรือ `http://host.docker.internal:8001` (หากรันอยู่นอก Container)
   - **Default Port**: `8001`
5. คลิก **Create Connection** จากนั้นกด **Activate**

---

## 2. Registering Services & Routes

ในการให้ Kong Gateway ส่งต่อ Request ไปยัง Microservices ทั้งสามตัว ต้องลงทะเบียน Service และ Route ดังต่อไปนี้:

### 2.1 Configuration Mapping Table

| Service Name | Host | Port | Route Name | Paths | Strip Path |
| :--- | :--- | :---: | :--- | :--- | :---: |
| `users-service` | `host.docker.internal` | `8081` | `users-route` | `/users` | `true` |
| `shop-service` | `host.docker.internal` | `8082` | `shop-route` | `/shop` | `true` |
| `admin-service` | `host.docker.internal` | `8083` | `admin-route` | `/admin` | `true` |

> **หมายเหตุสำหรับ Host**:
> - หาก Microservice รันอยู่บนเครื่องโฮสต์ (Local Machine) ให้ใช้ `host.docker.internal` เพื่อให้ Kong ในคอนเทนเนอร์เชื่อมต่อออกมาได้
> - หากเพื่อนร่วมทีมใช้ ngrok ให้ระบุ Host เป็นโดเมนที่ได้จาก ngrok เช่น `xxxx.ngrok-free.app` และตั้ง Protocol เป็น `https` พอร์ต `443`

### 2.2 ขั้นตอนการเพิ่ม Service ใน Konga
1. ไปที่เมนู **Services** -> คลิก **+ ADD NEW SERVICE**
2. ใส่ค่า **Name**, **Protocol** (`http`), **Host**, และ **Port** ตามตารางด้านบน
3. คลิก **Save Service**

### 2.3 ขั้นตอนการเพิ่ม Route ในแต่ละ Service
1. คลิกเข้าไปที่ Service ที่สร้างขึ้น -> ไปที่แท็บ **Routes** -> คลิก **+ ADD ROUTE**
2. ใส่ค่า **Name** (เช่น `users-route`)
3. ในช่อง **Paths** พิมพ์ Path นำหน้า (เช่น `/users`) แล้วกด Enter
4. ตรวจสอบให้แน่ใจว่า **Strip Path** ถูกเลือกเป็น **`true`** เพื่อให้ Kong ตัด Prefix `/users` ออกก่อนส่งต่อไปยัง Service หลังบ้าน
5. คลิก **Save Route**

---

## 3. Configuring Gateway Plugins

Kong Gateway มีระบบ Plugin เพื่อจัดการ Cross-cutting Concerns ในจุดศูนย์กลาง:

### 3.1 CORS Plugin (Global)
เปิดใช้งานเพื่อให้ Client Frontend สามารถเรียก API ข้ามโดเมนได้:

1. ไปที่เมนู **Plugins** ใน Konga -> คลิก **Add Global Plugin**
2. เลือก **CORS** และระบุค่า:
   - **origins**: `*` (หรือระบุ URL Frontend เช่น `http://localhost:3000`)
   - **methods**: `GET, POST, PUT, PATCH, DELETE, OPTIONS`
   - **headers**: `Authorization, Content-Type, Accept`
   - **credentials**: `true`
   - **max_age**: `3600`
3. คลิก **Save Plugin**

### 3.2 JWT Authentication Plugin
ใช้สำหรับตรวจสอบความถูกต้องของ Token ก่อนอนุญาตให้คำขอผ่านไปยัง Protected Routes:

1. ไปที่ Route หรือ Service ที่ต้องการป้องกัน
2. คลิกแท็บ **Plugins** -> คลิก **+ ADD PLUGIN** -> เลือก **JWT**
3. ระบุค่า:
   - **key_claim_name**: `iss`
   - **secret_is_base64**: `false`
   - **run_on_preflight**: `true`
4. คลิก **Save Plugin**

### 3.3 Rate Limiting Plugin
ป้องกันการเรียกใช้งานที่ถี่เกินไปหรือการโจมตีแบบ Denial of Service:

1. ไปที่เมนู **Plugins** (สามารถเลือกผูกแบบ Global หรือเจาะจงเฉพาะ Service)
2. เลือก **Rate Limiting** และระบุค่า:
   - **minute**: `100` (จำกัด 100 requests ต่อนาที)
   - **hour**: `10000` (จำกัด 10,000 requests ต่อชั่วโมง)
   - **policy**: `local` (ใช้ Memory ของ Kong สำหรับ Single Instance)
3. คลิก **Save Plugin**

---

## 4. Verification via Kong Proxy

เมื่อตั้งค่าเสร็จสิ้น ให้ทดสอบผ่าน Proxy พอร์ต `8000`:

```bash
# ทดสอบการเรียกเส้นทาง /users
curl -i http://localhost:8000/users/healthz

# ทดสอบ CORS Preflight Request
curl -I -X OPTIONS http://localhost:8000/shop/healthz \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: GET"
```
