# QA Environment สถาปัตยกรรม

## ภาพรวม

QA Environment เป็นสภาพแวดล้อมทดสอบที่ใช้ในการทดสอบและตรวจสอบระบบ DECM Platform ก่อนที่จะนำไปใช้งานจริง สภาพแวดล้อมนี้ใช้หลาย cloud services เพื่อรองรับการทำงานของแอปพลิเคชัน

## Diagram สถาปัตยกรรม

(อ้างอิงจาก FigJam: QA Environment diagram)

## คอมโพเนนต์

### 1. Frontend - Vercel

**Server**

- บริการโฮสต์ frontend application (React)
- รับการขอ (requests) จากผู้ใช้ Desktop
- สร้าง build จาก Vercel Compute

**Compute**

- รับการทำให้ triggered จาก GitHub Package
- ทำการ build frontend application
- ส่งใช้งาน (Deploy) ไปยัง Vercel Server

### 2. การควบคุมเวอร์ชัน - GitHub

**Package**

- จัดเก็บโค้ดต้นฉบับ (source code) และ registry ของ package
- ทำให้เกิด pipeline เมื่อมีการรวมการเปลี่ยนแปลง (merge changes)
- ส่งข้อมูลไปยัง Vercel Compute เพื่อสร้าง (build) และส่งใช้งาน (deploy)

### 3. โครงสร้างพื้นฐาน Backend - Huawei Cloud

**Elastic Computing Service**

- ให้บริการทรัพยากรคำนวณ (compute resources) สำหรับบริการ backend
- รองรับการส่งใช้งาน Docker container

**Docker Container**

ภายใน container มี services หลัก 3 ตัว:

1. **Service**
    - เลเยอร์ API service ที่รับการขอ (requests) จาก Desktop
    - ส่งต่อการขอ (requests) ไปยัง Backend API

2. **Backend API**
    - เลเยอร์ที่จัดการตรรกะทางธุรกิจ (business logic)
    - รับการขอ (requests) จาก Service
    - ทำการสอบถาม (query) และการดำเนินการเปลี่ยนแปลง (mutation) ไปยัง Database
    - สื่อสารกับ AWS Storage สำหรับสินทรัพย์ (assets)

3. **Database**
    - ฐานข้อมูล PostgreSQL
    - เก็บข้อมูลส่วนบุคคล (PII) ที่เข้ารหัส
    - รับการสอบถาม (query) และการดำเนินการเปลี่ยนแปลง (mutation) จาก Backend API

### 4. ที่เก็บข้อมูล - บริการ AWS Cloud

**ที่เก็บข้อมูล (Storage)**

- เก็บสินทรัพย์แบบคงที่ (static assets) และไฟล์
- ให้บริการ URLs ที่ลงนาม (presigned URLs) สำหรับผู้ใช้ desktop
- รองรับการโหลดสินทรัพย์ (assets) จาก frontend

## ขั้นตอนการทำงาน (Workflows)

### 1. ขั้นตอนการพัฒนา (Development Workflow)

ผู้พัฒนาแก้ไขโค้ดในเครื่องเอง จากนั้น commit และ push ไปยัง GitHub แล้วรวมการเปลี่ยนแปลง (merge changes) ไปยัง QA branch

เมื่อมีเหตุการณ์การรวม (merge event) GitHub Package จะทำให้เกิด pipeline ไปยัง Vercel Compute เพื่อทำการสร้าง (build) frontend application แล้วส่งใช้งาน (deploy) ไปยัง Vercel Server

สำหรับ backend มีตัวเลือกการส่งใช้งาน (deploy) แบบด้วยตนเอง (manual) โดยผู้ใช้สามารถรัน docker compose command เพื่อส่งใช้งาน (deploy) ไปยัง Huawei Cloud Docker Container ซึ่งจะทำการส่งใช้งาน (deploy) Service, Backend API และเชื่อมต่อไปยัง Database

### 2. ขั้นตอนการเข้าถึงของผู้ใช้ (User Access Workflow)

ผู้ใช้ Desktop เข้าถึง Vercel Server (Frontend) เพื่อใช้งาน web application และจะโหลดสินทรัพย์ (assets) จาก AWS Storage ผ่าน URLs ที่ลงนาม (presigned URLs)

เมื่อผู้ใช้กดปุ่มหรือส่งการกระทำ (action) ต่างๆ จะส่งการเรียก API (API call) ไปยัง Service ใน Docker Container เพื่อส่งต่อการขอ (request) ไปยัง เลเยอร์ Backend API

Backend API จะทำการสอบถาม (query) และการดำเนินการเปลี่ยนแปลง (mutation) ไปยัง Database และเข้าถึง AWS Storage (ถ้าจำเป็น) แล้วส่งข้อมูลกลับมาให้ Service ซึ่งจะส่งกลับไปยังผู้ใช้ Desktop

## การไหลของข้อมูล (Data Flow)

### การไหลของการขอ (Request Flow)

1. **Desktop** → **Vercel Server**: ผู้ใช้เข้าถึง frontend application
2. **Desktop** → **AWS Storage**: โหลดสินทรัพย์แบบคงที่ (static assets) ผ่าน URLs ที่ลงนาม (presigned URLs)
3. **Desktop** → **Service**: เรียก API จาก frontend
4. **Service** → **Backend API**: ส่งต่อการขอ (requests) ไปยัง เลเยอร์ตรรกะทางธุรกิจ (business logic layer)
5. **Backend API** → **Database**: ประมวลผลการสอบถาม (queries) และการดำเนินการเปลี่ยนแปลง (mutations)
6. **Backend API** → **AWS Storage**: เข้าถึงไฟล์/สินทรัพย์ (files/assets) (ถ้าจำเป็น)

### การไหลของการตอบกลับ (Response Flow)

1. **Database** → **Backend API**: ส่งคืนผลการสอบถาม (query results)
2. **AWS Storage** → **Backend API**: ส่งคืนข้อมูลไฟล์ (file data) (ถ้าจำเป็น)
3. **Backend API** → **Service**: ส่งคืนข้อมูลที่ประมวลผล
4. **Service** → **Desktop**: ส่งคืนการตอบกลับ API (API response)

## กลยุทธ์การส่งใช้งาน (Deployment Strategy)

### Frontend (อัตโนมัติ)

- **Trigger**: เหตุการณ์ GitHub merge/push
- **Process**: GitHub Package → Vercel Compute → สร้าง (Build) → ส่งใช้งาน (Deploy) ไปยัง Vercel Server
- **Environment**: Pipeline CI/CD อัตโนมัติ

### Backend (ด้วยตนเอง - Manual)

- **Trigger**: การกระทำด้วยตนเอง
- **Process**: ผู้ใช้รัน docker compose command
- **Target**: บริการ Huawei Cloud Elastic Computing
- **Environment**: Docker containers ใน Huawei Cloud

## การพิจารณาด้านความปลอดภัย (Security Considerations)

1. **ข้อมูลส่วนบุคคลที่เข้ารหัส (Encrypted PII)**: Database เก็บข้อมูลข้อมูลส่วนบุคคล (PII) ที่เข้ารหัสโดยใช้ AES-GCM encryption ที่ application layer
2. **URLs ที่ลงนาม (Presigned URLs)**: AWS Storage ให้บริการการเข้าถึงสินทรัพย์ (assets) ที่มีความปลอดภัยและจำกัดเวลา
3. **การแยก Container**: Services ทำงานในแยก Docker containers
4. **การตรวจสอบสิทธิ์ API (API Authentication)**: เลเยอร์ Service ยืนยันและตรวจสอบสิทธิ์ของการขอ (requests)

## ลักษณะของสภาพแวดล้อม

- **วัตถุประสงค์**: ทดสอบและ Quality Assurance
- **ความเสถียร**: อาจรวมถึงฟีเจอร์ทดลองและการตั้งค่าต่างๆ
- **การเข้าถึง**: จำกัดไว้สำหรับทีม QA และผู้มีส่วนได้ส่วนเสีย (stakeholders)
- **ข้อมูล**: อาจใช้ข้อมูลทดสอบสำหรับการทดสอบ

## Tech Stack

- **Frontend**: React 19 (โฮสต์บน Vercel)
- **Backend**: Go Fiber API (Docker บน Huawei Cloud) พร้อม AES-GCM encryption
- **Database**: PostgreSQL เก็บข้อมูลที่เข้ารหัสแล้ว
- **Storage**: AWS S3 (หรือบริการที่เก็บข้อมูลที่เข้ากันได้)
- **Containerization**: Docker
- **CI/CD**: GitHub Actions + Vercel

## การเชื่อมต่อระบบเครือข่าย

**จากผู้ใช้ Desktop**

- เข้าถึง Vercel Server ผ่าน HTTPS สำหรับเข้าใช้งาน web application
- เรียก API ไปยัง Service ผ่าน HTTPS
- โหลดสินทรัพย์ (assets) จาก AWS Storage ผ่าน HTTPS โดยใช้ URLs ที่ลงนาม (presigned URLs)

**จาก Service ไปยัง Backend API**

- ส่งต่อการขอ (requests) ผ่านการเชื่อมต่อภายใน (Internal connection)

**จาก Backend API**

- ทำการสอบถาม (query) และการดำเนินการเปลี่ยนแปลง (mutation) ไปยัง Database ผ่านการเชื่อมต่อ PostgreSQL
- เข้าถึง AWS Storage ผ่าน HTTPS/S3 API

**จาก GitHub Package**

- ส่ง Webhook ไปยัง Vercel Compute เพื่อทำให้เกิด pipeline CI/CD

**จากผู้ใช้**

- ทำการควบคุมเวอร์ชัน (version control) ผ่าน Git ไปยัง GitHub Package
- ทำการส่งใช้งาน (deployment) ด้วยตนเอง (manual) ไปยัง Huawei Cloud ผ่าน Docker

## หมายเหตุ

- QA Environment แยกออกมาจาก Production Environment เพื่อความปลอดภัย
- การส่งใช้งาน (deployment) ด้วยตนเอง (manual) สำหรับ backend ช่วยให้ทีมสามารถควบคุมเวลาและการตั้งค่า (timing และ configuration) ได้
- AWS Storage และบริการ Huawei Cloud สามารถทำงานร่วมกันผ่าน public internet
- Frontend และ Backend สามารถขยายขนาด (scale) แยกกันได้ตามความต้องการ
