# Frontend Architecture

## ภาพรวม

Frontend Architecture ของ DECM Platform ใช้ React 19 ด้วย Vite เป็น build tool ประกอบด้วยหลายชั้นที่ทำงานร่วมกัน เพื่อให้เกิดการแยกความสัมพันธ์และทำให้โครงสร้างเป็นระเบียบ

## Diagram สถาปัตยกรรม

(อ้างอิงจาก FigJam: Frontend Architecture diagram)

## องค์ประกอบหลัก

### 1. End User (ผู้ใช้ปลายทาง)

ผู้ใช้เข้าถึงแอปพลิเคชันผ่านเว็บเบราว์เซอร์ และ visit website ของ DECM Platform

### 2. React Vite Web Server

**Web Component** รับ visits จากผู้ใช้ปลายทาง ทำหน้าที่เป็น entry point ของ web application และเชื่อมต่อกับ React Router

**Package (Generouté)** จัดเก็บไฟล์การตั้งค่า routing Configure routing rules โดยส่งต่อไปยัง Web component และใช้ directory-based routing pattern

**Pages** ประกอบไปด้วย route pages ทั้งหมด แต่ละ page เป็น minimal wrapper สำหรับ route ใช้ directory-based routing โดยได้รับการ configure จาก Generouté และส่งต่อไปยัง API Client เพื่อดึงข้อมูล

**API Client** เป็น generated client จากไฟล์ OpenAPI specifications ทำหน้าที่เป็น bridge ระหว่าง frontend กับ Backend API และส่งไปยัง Server

### 3. React Router

**Routing System** จัดการการสลับ routes/pages รับ configuration จาก Generouté และส่งต่อไปยัง Web component

**Directory-based Routing** ใช้ไฟล์ directory structure เพื่อ define routes และได้รับการ configure จาก Generouté

### 4. API Client

**Generation** ทำการ generate จากไฟล์ OpenAPI specifications ของ Backend และใช้ไฟล์ generated TypeScript client แทนการเขียนด้วยมือ

### 5. Backend

Backend API (Go Fiber) ให้บริการ endpoints และส่ง OpenAPI specifications เพื่อให้ Frontend สามารถ generate API client

## Workflows

### 1. Initialization Workflow

เมื่อ End user visit website จะเกิดขึ้นตามลำดับดังนี้:

1. End user เข้าถึง Web component จาก React Vite
2. Web component เชื่อมต่อกับ React Router เพื่อ routing
3. React Router ได้รับ configuration จาก Generouté
4. Generouté define routing rules
5. Web component ทำการ routing ไปยัง Pages ที่เหมาะสม
6. Pages import และ render page components
7. Page components ใช้ API Client เพื่อดึงข้อมูล
8. API Client ส่ง requests ไปยัง Backend Server

### 2. API Call Workflow

เมื่อ Page component ต้องการข้อมูล จะเกิดขึ้นตามลำดับดังนี้:

1. Page component เรียก API Client (generated TypeScript client)
2. API Client ส่ง request ไปยัง Backend API server
3. Backend API ส่งข้อมูลกลับ
4. API Client ส่ง response กลับไปยัง Page component

## Directory Structure

Frontend source code จัดระเบียบในโครงสร้างดังนี้:

**Pages Directory** มี route pages (Vite routes) ที่กำหนด URL paths รวมทั้ง index route signup route และ oauth routes เป็นต้น

**Components Directory** ที่สำคัญมี:

- Pages subdirectory มี page UI components ที่มี UI logic ทั้งหมด เช่น LandingPage SignUp GoogleOAuth
- UI subdirectory มี Radix UI base components เช่น button card dialog เป็นต้น
- Providers subdirectory มี context providers และ helmet components สำหรับ SEO
- Typography subdirectory มี typography component สำหรับ consistent text rendering

**Hooks Directory** มี custom React hooks สำหรับ reusable logic เช่น useWallet useAuth เป็นต้น

**Config Directory** มี configuration files สำหรับ environment variables และ API settings

**Lib Directory** มี utility functions และ helper functions

**Context Directory** มี React Context สำหรับ global state management เช่น AuthContext

**Router Configuration** อยู่ใน router.ts file

**App Entry Point** อยู่ใน index.tsx

**Global Styles** อยู่ใน index.css

## Page Component Pattern

### Pages Directory (Route Files)

Pages ใน /pages/ ควรเป็น minimal wrapper เท่านั้น โดยจัดการเฉพาะ SEO metadata ผ่าน helmet components และ import page component จาก /components/pages/

### Components/Pages Directory (UI Logic)

Page components ใน /components/pages/ ประกอบด้วย UI logic ทั้งหมด รวมทั้ง data fetching state management และ event handlers

## UI Components

UI components ใช้ Radix UI primitives ร่วมกับ Tailwind CSS สำหรับ styling ทุก component มี variant support ผ่าน Class Variance Authority (CVA) เพื่อให้สามารถ customize styling ได้อย่างยืดหยุ่น

## Routing

Frontend ใช้ React Router v6 สำหรับการจัดการ routes โดย configuration อยู่ใน router.ts file และใช้ lazy loading pattern เพื่อลด bundle size

Vite สนับสนุน directory-based routing ผ่าน glob imports ซึ่งช่วยให้สามารถ auto-generate routes ได้

## API Integration

API Client ถูก generate จาก OpenAPI specifications ของ Backend API โดยรัน command pnpm gen-api:core

Generated file อยู่ที่ packages/api/src/apis/core/api.ts

Page components ใช้ API client เพื่อ fetch data โดยเรียก methods จาก DefaultApi class พร้อม error handling

## Styling

Frontend ใช้ Tailwind CSS สำหรับ utility-based styling ที่ให้ความยืดหยุ่นสูง

Class Variance Authority (CVA) ใช้สำหรับ variant-based styling เพื่อให้สามารถ manage component variants ได้อย่างเป็นระเบียบ

## State Management

React Context ใช้สำหรับ global state management เช่น authentication state user preferences เป็นต้น

Custom hooks ใช้สำหรับ reusable logic ที่ซ่อนรายละเอียด implementation detail ไว้ และให้ interface ที่สะอาดให้กับ components

## Typography Component

Typography component ใช้สำหรับ consistent text rendering โดยรองรับ multiple variants เช่น header body caption เป็นต้น ทำให้มั่นใจว่า typography ทั่วทั้ง application มี consistency

## Form Handling

Form handling ใช้ React Hook Form สำหรับ efficient form management

Zod ใช้สำหรับ schema validation ทั้ง frontend และ backend ซึ่งช่วยให้ data validation มี consistency

Form components รองรับ real-time validation error messages และ accessibility attributes

## Internationalization (i18n)

react-i18next ใช้สำหรับ translations รองรับ multiple languages และ language switching

ทุก text strings ใน application ควรใช้ translation keys แทนการ hard-code text

## Build & Development

**Development Server** ถูก start โดย pnpm dev command และสามารถเข้าถึงได้ที่ http://localhost:3000

**Production Build** ถูก create โดย pnpm build command

**Type Checking** ถูกทำโดย pnpm check-types command

**Linting** ถูกทำโดย pnpm lint command

## Technology Stack

- **Framework**: React 19
- **Build Tool**: Vite
- **Routing**: React Router v6
- **UI Components**: Radix UI
- **Styling**: Tailwind CSS + CVA
- **Form Management**: React Hook Form + Zod
- **Internationalization**: react-i18next
- **API Client**: Generated from OpenAPI
- **State Management**: React Context + Custom Hooks
- **Package Manager**: pnpm

## ข้อกำหนดในการพัฒนา

- Frontend ทั้งหมดใช้ React 19 ด้วย TypeScript
- ใช้ directory-based routing แบบ Vite
- Page components แยกออกจาก route files เพื่อ separation of concerns
- UI components ใช้ Radix UI + Tailwind CSS สำหรับ styling
- Form validation ใช้ Zod + React Hook Form เพื่อ end-to-end type safety
- API client ถูก generate จาก OpenAPI specifications เพื่อลด boilerplate code
- ทุกข้อความใช้ i18n translations เพื่อรองรับ multiple languages
- Global state ใช้ React Context + Custom Hooks
- ไม่ใช้ dangerouslySetInnerHTML เพื่อความปลอดภัย
- Lazy loading ใช้สำหรับ routes เพื่อให้ initial bundle size เล็กลง
