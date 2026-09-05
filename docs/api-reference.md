# API Reference Guide

คู่มืออ้างอิง API Endpoints ทั้งหมดของระบบ **GameGear E-commerce** โดยทุก Endpoint ในเอกสารนี้อ้างอิงการเรียกผ่าน **Kong API Gateway Proxy (`http://localhost:8000`)**

---

## Standards & Headers

- **Base URL**: `http://localhost:8000`
- **Content-Type**: `application/json` (สำหรับ Request ที่มี Body)
- **Authorization**: `Bearer <JWT_TOKEN>` (สำหรับ Protected Endpoints)

---

## 1. Users Service Endpoints (`/users`)

### 1.1 Member Authentication

#### Register Member
- **Method**: `POST`
- **Path**: `/users/api/auth/register`
- **Auth Required**: No

```bash
curl -X POST http://localhost:8000/users/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "gamer123",
    "display_name": "Pro Gamer",
    "email": "gamer@example.com",
    "password": "password123",
    "confirm_password": "password123"
  }'
```

#### Login Member
- **Method**: `POST`
- **Path**: `/users/api/auth/login`
- **Auth Required**: No

```bash
curl -X POST http://localhost:8000/users/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "gamer@example.com",
    "password": "password123"
  }'
```
*Response คืน JWT Token ที่มี claim `role: "member"`*

#### Forgot Password
- **Method**: `POST`
- **Path**: `/users/api/auth/forgot-password`
- **Auth Required**: No

```bash
curl -X POST http://localhost:8000/users/api/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "gamer@example.com"}'
```

#### Reset Password
- **Method**: `POST`
- **Path**: `/users/api/auth/reset-password`
- **Auth Required**: No

```bash
curl -X POST http://localhost:8000/users/api/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "RESET_TOKEN_STRING",
    "new_password": "newpassword123",
    "confirm_password": "newpassword123"
  }'
```

#### Logout Member
- **Method**: `POST`
- **Path**: `/users/api/auth/logout`
- **Auth Required**: Yes (Member JWT)

```bash
curl -X POST http://localhost:8000/users/api/auth/logout \
  -H "Authorization: Bearer <MEMBER_JWT>"
```

---

### 1.2 Member Profile

#### Get Profile
- **Method**: `GET`
- **Path**: `/users/api/user/profile`
- **Auth Required**: Yes (Member JWT)

```bash
curl -sS http://localhost:8000/users/api/user/profile \
  -H "Authorization: Bearer <MEMBER_JWT>"
```

#### Update Profile
- **Method**: `PUT`
- **Path**: `/users/api/user/profile`
- **Auth Required**: Yes (Member JWT)

```bash
curl -X PUT http://localhost:8000/users/api/user/profile \
  -H "Authorization: Bearer <MEMBER_JWT>" \
  -H "Content-Type: application/json" \
  -d '{
    "display_name": "Updated Gamer Name",
    "profile_image": "https://example.com/avatar.png"
  }'
```

#### Delete Account
- **Method**: `DELETE`
- **Path**: `/users/api/user/profile`
- **Auth Required**: Yes (Member JWT)

```bash
curl -X DELETE http://localhost:8000/users/api/user/profile \
  -H "Authorization: Bearer <MEMBER_JWT>" \
  -H "Content-Type: application/json" \
  -d '{"password": "currentPassword123"}'
```

---

### 1.3 Admin Authentication (Called by admin-service or Admin Client)

#### Admin Login
- **Method**: `POST`
- **Path**: `/users/api/admin/login`
- **Auth Required**: No

```bash
curl -X POST http://localhost:8000/users/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@gamegear.com",
    "password": "adminPassword123"
  }'
```
*Response คืน JWT Token ที่มี claim `role: "admin"`*

#### Admin Registration
- **Method**: `POST`
- **Path**: `/users/api/admin/register`
- **Auth Required**: No

```bash
curl -X POST http://localhost:8000/users/api/admin/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "operations@gamegear.com",
    "password": "securePassword123",
    "confirm_password": "securePassword123",
    "display_name": "Operations Lead"
  }'
```

---

## 2. Shop Service Endpoints (`/shop`)

### 2.1 Public Product Catalog

#### List Products (with Pagination)
- **Method**: `GET`
- **Path**: `/shop/api/products?page=1&limit=12`
- **Auth Required**: No

```bash
curl -sS "http://localhost:8000/shop/api/products?page=1&limit=12"
```

#### Get Product By ID
- **Method**: `GET`
- **Path**: `/shop/api/products/:id`
- **Auth Required**: No

```bash
curl -sS http://localhost:8000/shop/api/products/1
```

---

### 2.2 Cart Operations (Member)

#### Get Current Cart
- **Method**: `GET`
- **Path**: `/shop/api/cart`
- **Auth Required**: Yes (Member JWT)

```bash
curl -sS http://localhost:8000/shop/api/cart \
  -H "Authorization: Bearer <MEMBER_JWT>"
```

#### Add Item to Cart
- **Method**: `POST`
- **Path**: `/shop/api/cart/add`
- **Auth Required**: Yes (Member JWT)

```bash
curl -X POST http://localhost:8000/shop/api/cart/add \
  -H "Authorization: Bearer <MEMBER_JWT>" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'
```

#### Update Cart Item Quantity
- **Method**: `PUT`
- **Path**: `/shop/api/cart/update`
- **Auth Required**: Yes (Member JWT)

```bash
curl -X PUT http://localhost:8000/shop/api/cart/update \
  -H "Authorization: Bearer <MEMBER_JWT>" \
  -H "Content-Type: application/json" \
  -d '{
    "cart_item_id": 10,
    "quantity": 3
  }'
```

#### Remove Item from Cart
- **Method**: `DELETE`
- **Path**: `/shop/api/cart/remove`
- **Auth Required**: Yes (Member JWT)

```bash
curl -X DELETE http://localhost:8000/shop/api/cart/remove \
  -H "Authorization: Bearer <MEMBER_JWT>" \
  -H "Content-Type: application/json" \
  -d '{"cart_item_id": 10}'
```

---

### 2.3 Orders (Member)

#### Create Order from Cart
- **Method**: `POST`
- **Path**: `/shop/api/orders`
- **Auth Required**: Yes (Member JWT)

```bash
curl -X POST http://localhost:8000/shop/api/orders \
  -H "Authorization: Bearer <MEMBER_JWT>" \
  -H "Content-Type: application/json" \
  -d '{
    "cart_id": 12,
    "shipping_address": "99/1 Rama 9 Rd, Huai Khwang, Bangkok 10310",
    "payment_method": "credit_card",
    "notes": "Deliver during office hours"
  }'
```

#### Get Order History
- **Method**: `GET`
- **Path**: `/shop/api/orders/history`
- **Auth Required**: Yes (Member JWT)

```bash
curl -sS http://localhost:8000/shop/api/orders/history \
  -H "Authorization: Bearer <MEMBER_JWT>"
```

---

### 2.4 Admin Product & Order Management (Direct Shop Service)

#### Create Product (Admin)
- **Method**: `POST`
- **Path**: `/shop/api/products`
- **Auth Required**: Yes (Admin JWT)

```bash
curl -X POST http://localhost:8000/shop/api/products \
  -H "Authorization: Bearer <ADMIN_JWT>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mechanical Keyboard RGB Pro",
    "description": "Hot-swappable mechanical gaming keyboard",
    "price": 2890,
    "stock": 50,
    "category_id": 2,
    "image_url": "https://example.com/keyboard.jpg"
  }'
```

#### Update Product (Admin)
- **Method**: `PUT`
- **Path**: `/shop/api/products/:id`
- **Auth Required**: Yes (Admin JWT)

```bash
curl -X PUT http://localhost:8000/shop/api/products/1 \
  -H "Authorization: Bearer <ADMIN_JWT>" \
  -H "Content-Type: application/json" \
  -d '{
    "price": 2690,
    "stock": 45
  }'
```

#### Delete Product (Admin)
- **Method**: `DELETE`
- **Path**: `/shop/api/products/:id`
- **Auth Required**: Yes (Admin JWT)

```bash
curl -X DELETE http://localhost:8000/shop/api/products/1 \
  -H "Authorization: Bearer <ADMIN_JWT>"
```

#### List All Orders (Admin)
- **Method**: `GET`
- **Path**: `/shop/api/orders`
- **Auth Required**: Yes (Admin JWT)

```bash
curl -sS http://localhost:8000/shop/api/orders \
  -H "Authorization: Bearer <ADMIN_JWT>"
```

#### Update Order Status (Admin)
- **Method**: `PUT`
- **Path**: `/shop/api/orders/:id/status`
- **Auth Required**: Yes (Admin JWT)

```bash
curl -X PUT http://localhost:8000/shop/api/orders/1/status \
  -H "Authorization: Bearer <ADMIN_JWT>" \
  -H "Content-Type: application/json" \
  -d '{"status": "shipped"}'
```

---

## 3. Admin Service Coordinator Endpoints (`/admin`)

Endpoints ของ Admin Service จะทำหน้าที่เป็น Coordinator สำหรับ Admin Dashboard Web Interface โดยจะเรียกต่อเชื่อมไปยัง Users และ Shop Service:

| Method | Path | Auth Required | หน้าที่ |
| :--- | :--- | :---: | :--- |
| `POST` | `/admin/api/admin/login` | No | ล็อกอินแอดมิน (Proxy ไปยัง users-service) |
| `POST` | `/admin/api/admin/register` | No | ลงทะเบียนแอดมิน (Proxy ไปยัง users-service) |
| `GET` | `/admin/api/admin/products` | Yes (Admin JWT) | ดึงรายการสินค้าทั้งหมดเพื่อนำไปจัดการ |
| `POST` | `/admin/api/admin/products` | Yes (Admin JWT) | สร้างสินค้าใหม่ผ่าน Coordinator logic |
| `PUT` | `/admin/api/admin/products/:id` | Yes (Admin JWT) | ปรับปรุงข้อมูลสินค้า |
| `DELETE` | `/admin/api/admin/products/:id` | Yes (Admin JWT) | ลบสินค้าออกจากระบบ |
| `GET` | `/admin/api/admin/orders` | Yes (Admin JWT) | ดึงรายการคำสั่งซื้อทั้งหมด |
| `PUT` | `/admin/api/admin/orders/:id/status` | Yes (Admin JWT) | อัปเดตสถานะของคำสั่งซื้อ |
