# 🚀 Mini-Project: GameGear E-commerce (Microservice Architecture)

![Go](https://img.shields.io/badge/Go-1.24.6-00ADD8?style=for-the-badge\&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge\&logo=postgresql)
![Microservices](https://img.shields.io/badge/Architecture-Microservices-2ea44f?style=for-the-badge)

โปรเจกต์นี้คือระบบ **E-commerce** สำหรับขายอุปกรณ์เกมมิ่ง ซึ่งถูกพัฒนาขึ้นโดยใช้สถาปัตยกรรมแบบ **Microservice** โดยแบ่งการทำงานออกเป็น Service ย่อยๆ ที่ทำงานร่วมกัน

---

## 🏛️ System Architecture Overview

สถาปัตยกรรมของโปรเจกต์ประกอบด้วย Service หลัก 3 ส่วนที่ทำงานแยกจากกัน แต่สื่อสารกันผ่าน **API** โดยมี API Gateway เป็นจุดรับ Request จากผู้ใช้งาน

---

## 📂 Service Repositories

โปรเจกต์นี้ประกอบด้วย Repository ย่อย 3 ส่วน โดยแต่ละส่วนมีหน้าที่และความรับผิดชอบต่างกัน

| Service Repository                       | Description                                                    | Team Member         |
| ---------------------------------------- | -------------------------------------------------------------- | ------------------- |
| 👤 **[users-service](https://github.com/Wattanaroj2567/users-service.git)**  | จัดการผู้ใช้และการยืนยันตัวตนทั้งหมด (สมัคร, ล็อกอิน, โปรไฟล์) | ณิชพน มานิตย์       |
| 🛍️ **[shop-service](https://github.com/Wattanaroj2567/shop-service.git)**   | จัดการการซื้อขายทั้งหมด (สินค้า, ตะกร้า, คำสั่งซื้อ)           | ณัฐพงษ์ ดีบุตร, วายุ กอคูณ   |
| 🛡️ **[admin-service](https://github.com/Wattanaroj2567/admin-service.git)** | จัดการระบบหลังบ้านสำหรับ Admin โดยจะเรียกใช้ Service อื่นๆ     | วรรธนโรจน์ บุตรดี         |

---

## 🤝 Development Team

| Name                  | Responsibility        |
| --------------------- | --------------------- |
| **วรรธนโรจน์ บุตรดี** | Project Manager & Dev |
| **ณัฐพงษ์ ดีบุตร**    | Backend Developer     |
| **ณิชพน มานิตย์**     | Backend Developer     |
| **วายุ กอคูณ**        | Backend Developer     |

---

## 🚀 How to Run the Entire Project

### 1. Clone all repositories

```bash
git clone https://github.com/Wattanaroj2567/users-service.git
git clone https://github.com/Wattanaroj2567/shop-service.git
git clone https://github.com/Wattanaroj2567/admin-service.git
```

### 2. Follow `README.md` in each service

เข้าไปในโฟลเดอร์ของแต่ละ Service แล้วทำตามขั้นตอนใน `README.md` ของ Service นั้นๆ เพื่อตั้งค่าและรันโปรแกรม
