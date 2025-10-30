# Production Environment สถาปัตยกรรม

## ภาพรวม

Production Environment เป็นสภาพแวดล้อมการใช้งานจริงสำหรับระบบ DECM Platform ที่เสิร์ฟบริการให้แก่ผู้ใช้ปลายทาง สภาพแวดล้อมนี้ออกแบบมาเพื่อความเสถียร ความปลอดภัย และความสามารถในการขยายตัว โดยรวมทั้ง Version Control System และ VPS Infrastructure

## Diagram สถาปัตยกรรม

(อ้างอิงจาก FigJam: Production Environment diagram)

## คอมโพเนนต์

### 1. ผู้ใช้ (End User)

**Desktop**

- ผู้ใช้ปลายทางเข้าถึงแอปพลิเคชันผ่านเว็บเบราว์เซอร์
- ทำการ visit Frontend และเรียก API

### 2. VPS Infrastructure

VPS (Virtual Private Server) เป็นหัวใจของ Production Environment ที่รวมรวมคอมโพเนนต์ทั้งหมด

#### Docker Container

ภายใน Docker container มี services หลักดังนี้:

**Service**

- เลเยอร์ API service layer ที่รับ การขอ (requests) จาก End user
- ส่งต่อ การขอ (requests) ไปยัง เลเยอร์ Frontend service และ Backend API layer

**Frontend Service (React Vite)**

- บริการให้ web interface จาก VPS โดยตรง
- ใช้ React กับ Vite เพื่อประสิทธิภาพการโหลดที่เร็ว
- สร้างเสร็จแล้วและติดตั้งบน VPS

**Backend API (Core API - Go Fiber)**

- Business logic layer ที่จัดการการประมวลผลข้อมูล
- รับ การขอ (requests) จาก Service
- ทำการ สอบถาม (query) และการเปลี่ยนแปลง (mutation) ไปยัง Database
- เข้าถึง Storage เพื่อเก็บและดึง UGC (User Generated Content)

**Database**

- PostgreSQL database ที่จัดเก็บข้อมูลทั้งหมด
- เก็บข้อมูล PII ที่เข้ารหัส
- รับ query และ mutation จาก Backend API

**Storage**

- Local storage ใน VPS สำหรับเก็บ UGC (User Generated Content) เช่น รูปภาพ วิดีโอ เอกสาร
- ให้บริการการอ่านและเขียนข้อมูล UGC

### 3. Version Control - GitHub

**Package Repository**

- จัดเก็บ source code ของทั้ง frontend และ backend
- ผู้พัฒนาทำการ merge changes ไปยังสาขา production

### 4. Developer

**Manual ส่งใช้งานment**

- ผู้พัฒนาทำการ manual deployment ของ application ไปยัง VPS
- ใช้ docker compose command เพื่อ deploy Docker container บน VPS

## Workflows

### 1. Development and ส่งใช้งานment Workflow

ผู้พัฒนาแก้ไขโค้ดและ commit ไปยัง GitHub จากนั้น merge changes ไปยังสาขา production

เมื่อมี merge event ผู้พัฒนาจะทำการ manual deployment ไปยัง VPS โดยรัน docker compose command เพื่อ update และ redeploy Docker container ที่มี Service, เลเยอร์ Frontend service, Backend API และ Database

### 2. User Access Workflow

End user เข้าถึง เลเยอร์ Frontend service ผ่าน Desktop web browser

เมื่อผู้ใช้โต้ตอบกับ Frontend จะส่ง API call ไปยัง Service ใน Docker container

Service ส่งต่อ การขอ (requests) ไปยัง เลเยอร์ Frontend service และ Backend API layer สำหรับการประมวลผล

Backend API ทำการ สอบถาม (query) และการเปลี่ยนแปลง (mutation) ไปยัง Database และเข้าถึง Storage เพื่อดึงหรือเก็บ UGC

ข้อมูลจะถูกส่งกลับมาผ่าน Service ไปยัง เลเยอร์ Frontend service เพื่อแสดงผลให้แก่ผู้ใช้

## Data Flow

### Request Flow

1. **End user** → **Service**: ผู้ใช้ visit Frontend และเรียก API
2. **Service** → **เลเยอร์ Frontend service**: ส่งต่อ การขอ (requests) เพื่อแสดงผลหน้าเว็บ
3. **Service** → **Backend API**: ส่งต่อ การขอ (requests) สำหรับการประมวลผล business logic
4. **Backend API** → **Database**: ประมวลผล queries และ mutations
5. **Backend API** → **Storage**: เข้าถึง UGC สำหรับการอ่านหรือเขียนข้อมูล

### Response Flow

1. **Database** → **Backend API**: ส่งคืน ผลการสอบถาม (query results)
2. **Storage** → **Backend API**: ส่งคืน UGC data
3. **Backend API** → **Service**: ส่งคืนข้อมูลที่ประมวลผล
4. **เลเยอร์ Frontend service** → **Service**: ส่งคืน HTML/CSS/JavaScript ที่สร้างแล้ว
5. **Service** → **End user**: แสดงผลหน้าเว็บให้แก่ผู้ใช้

## ส่งใช้งานment Strategy

### Automated Pipeline

- **Trigger**: GitHub merge/push events
- **Process**: GitHub Package → Manual trigger by developer
- **Target**: VPS Docker Container
- **Environment**: Docker containers บน VPS
- **Components Updated**: Service, เลเยอร์ Frontend service, Backend API, Database

### ส่งใช้งานment Process

ผู้พัฒนาควบคุมการ deploy ด้วยตนเองโดย:

1. ทำการ merge code ไปยัง production branch ใน GitHub
2. รัน docker compose command เพื่ออัปเดต Docker container บน VPS
3. Docker จะ pull image ล่าสุด และ restart services

## Security Considerations

1. **Encrypted PII**: Database เก็บข้อมูลส่วนตัว (PII) ที่เข้ารหัสโดยใช้ AES-GCM encryption ที่ application layer
2. **Local Storage**: UGC เก็บไว้ใน VPS local storage ที่ได้รับการ backup
3. **Container Isolation**: Services ทำงานในแยก Docker containers บน VPS
4. **API Authentication**: เลเยอร์ Service ยืนยันและตรวจสอบสิทธิ์ของ การขอ (requests)
5. **VPS Security**: VPS ต้องมีการตั้งค่า firewall และ network security policies

## ลักษณะของสภาพแวดล้อม

- **วัตถุประสงค์**: บริการลูกค้ารายจริง Production
- **ความเสถียร**: ต้องมีความเสถียรสูง หลีกเลี่ยง downtime
- **การเข้าถึง**: เปิดให้ผู้ใช้ปลายทางเข้าถึงได้
- **ข้อมูล**: ใช้ข้อมูลจริงของผู้ใช้

## Technology Stack

- **Frontend**: React 19 + Vite (compiled and hosted on VPS)
- **Backend**: Go Fiber API (Docker on VPS) พร้อม AES-GCM encryption
- **Database**: PostgreSQL เก็บข้อมูลที่เข้ารหัสแล้ว
- **Storage**: Local VPS storage สำหรับ UGC
- **Containerization**: Docker
- **ส่งใช้งานment**: Manual docker compose deployment
- **Version Control**: GitHub

## การเชื่อมต่อระบบเครือข่าย

**จาก End user**

- เข้าถึง VPS Service ผ่าน HTTPS สำหรับเข้าใช้งาน web application
- ทำการ visit และเรียก API จาก Desktop

**ภายใน VPS Docker Container**

- Service ส่งต่อ การขอ (requests) ไปยัง เลเยอร์ Frontend service และ Backend API layer
- Backend API ทำการ สอบถาม (query) และการเปลี่ยนแปลง (mutation) ไปยัง Database ผ่าน PostgreSQL connection
- Backend API เข้าถึง Storage ที่เก็บ UGC

**จาก GitHub**

- ผู้พัฒนา pull changes จาก GitHub repository

**จาก Developer**

- ทำการ manual deployment ไปยัง VPS ผ่าน docker compose command

## หมายเหตุ

- Production Environment ใช้ Single VPS ซึ่งรวม Frontend, Backend, Database และ Storage ทั้งหมดในที่เดียว
- ต้องมีการ backup ข้อมูล Database และ UGC Storage เป็นประจำ
- ผู้พัฒนามีสิทธิ์ควบคุม deployment โดยตรง
- VPS ต้องมีการ monitoring และ alerting เพื่อตรวจจับปัญหาได้อย่างรวดเร็ว
- ควรมีการ scale plan เพื่อรับมือกับการเพิ่มจำนวนผู้ใช้ในอนาคต
