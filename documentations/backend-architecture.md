# Backend Architecture

## ภาพรวม

Backend Architecture ของ DECM Platform ใช้ Go Fiber framework โดยปฏิบัติตาม clean architecture pattern ที่มีการแยกความสัมพันธ์อย่างชัดเจน ระหว่าง 3 ชั้นหลัก ได้แก่ Handler Layer, UseCase Layer และ Repository Layer

## Diagram สถาปัตยกรรม

(อ้างอิงจาก FigJam: Backend Architecture diagram)

## ส่วนประกอบหลัก

### 1. End User (ผู้ใช้ปลายทาง)

ผู้ใช้ส่ง HTTP request ไปยัง Backend API เพื่อเรียกใช้ services ต่างๆ

### 2. Core API System

**Handler Layer** รับ HTTP requests จาก clients และจัดการการ validate input ส่งต่อไปยัง UseCase layer และจัดการการแปลง errors เป็น HTTP responses

**Middleware** ใช้สำหรับ cross-cutting concerns เช่น authentication authorization logging request validation เป็นต้น

**Swagger/OpenAPI Documentation** ถูก define ใน Handler layer เพื่อให้ auto-generate API documentation

**UseCase Layer** ประกอบด้วย business logic ที่ orchestrates การทำงานระหว่าง repositories และ services ทั้ง usecase และ services ทำงานร่วมกัน

**Repository Layer** จัดการการเข้าถึง database ทำการ encrypt/decrypt PII fields และการทำ query ทั้งหมด

**Data Gateway** define interface ที่เป็นสัญญาระหว่าง Repository layer กับ external services

### 3. Database Tools & Migrations

**Database Migration Files** เก็บไฟล์ migration scripts ที่ใช้ manage database schema changes

**Database Queries Files** เก็บไฟล์ SQL queries ที่ใช้โดย sqlc เพื่อ generate Go code

**Terminal** ใช้สำหรับรัน migration commands และ generate queries

### 4. Generated Queries Package

Generated Go code จาก SQL files ผ่าน sqlc ประกอบด้วย type-safe database operations

### 5. API Client Package

Generated TypeScript client จาก OpenAPI specifications เพื่อให้ Frontend สามารถเรียก Backend API ได้

## Workflows

### 1. Request Handling Workflow

เมื่อ End user ส่ง HTTP request ไปยัง Backend API จะเกิดขึ้นตามลำดับดังนี้:

1. End user ส่ง request ไปยัง Backend API
2. Middleware ประมวลผล request (authentication authorization logging)
3. Handler parse request body validate input
4. Handler เรียก UseCase ที่เกี่ยวข้อง
5. UseCase ประมวลผล business logic และ coordinate กับ repositories
6. Repository ทำ database queries encrypt/decrypt PII data
7. Repository ส่งคืนข้อมูลไปยัง UseCase
8. UseCase ส่งคืน result ไปยัง Handler
9. Handler format response และส่งกลับไปยัง client

### 2. Database & Query Generation Workflow

Database schema definitions และ queries ถูก manage ดังนี้:

1. User define database migrations ใน migration files
2. Migration CLI tool ใช้สำหรับรัน migrations up/down
3. User write SQL queries ใน queries files
4. sqlc tool ใช้สำหรับ generate Go code จาก SQL
5. Generated queries package ถูก import ใน Repository layer
6. Generated API Client ถูก generate จาก OpenAPI specifications

## Architecture Layers

### Handler Layer (Internal/Handler)

Handler layer จัดการ HTTP request/response และการ validate input

**Responsibility**: HTTP handling input validation error mapping OpenAPI documentation

**File Structure**: แต่ละ feature มี handler files ที่ define methods สำหรับ actions ต่างๆ รวมทั้ง routes.go ไฟล์ที่ register routes

**Pattern**: Handler struct มี UseCase dependency ที่ inject เข้ามา แต่ละ handler method มี Swagger annotations สำหรับ auto-generate API documentation

**Input Validation**: Request body ต้อง validate ทั้ง parsing และ validation logic

**Error Mapping**: Errors ต้องถูก wrap และส่งกลับเป็น customerror instances ให้ middleware handle

### UseCase Layer (Internal/UseCase)

UseCase layer ประกอบด้วย business logic ที่ independent จาก HTTP specifics

**Responsibility**: Business logic orchestration transaction management validation rules

**File Structure**: แต่ละ feature มี usecase files ที่ implement specific use cases

**Pattern**: UseCase struct มี Repository dependencies ที่ inject เข้ามา แต่ละ use case method รับ context และ input parameters จากนั้นส่งคืน output หรือ error

**Business Rules**: Business validation logic ควร implement ใน use case layer ก่อนจะ access database

**Transaction Management**: Complex operations ที่ต้อง maintain data consistency ควร ใช้ transactions

### Repository Layer (Internal/Repository)

Repository layer จัดการ data access และ database operations

**Responsibility**: Database operations data access PII encryption/decryption query execution

**File Structure**: แต่ละ feature มี repository files ที่ implement CRUD operations

**Pattern**: Repository struct มี database queries reference ที่ generate จาก sqlc แต่ละ repository method handle database operations พร้อม error handling

**PII Encryption**: ทุก PII fields ต้อง encrypt ที่ application layer (Go code) ก่อน write ไปยัง database และ decrypt หลัง read จาก database โดยใช้ AES-GCM encryption

**Error Handling**: Database errors ต้องถูก parse และ wrap เป็น customerror instances

## Dependency Injection

Dependencies ถูก inject ผ่าน constructor functions แล้ว pass ไปยัง dependent layers

**Setup Pattern**: ใน main.go ให้ create repositories ก่อน จากนั้น create use cases และ pass repositories เข้าไป จากนั้น create handlers และ pass use cases เข้าไป สุดท้าย register routes

**Benefit**: ช่วยให้ easy to test และ decouple layers จากกัน

## Entity & DTO Pattern

**Domain Entity**: ใน internal/entity/ ประกอบด้วย data models ที่ represent domain concepts มี business methods ที่ operate on entities

**Request DTO**: ใน handler layer represent incoming request data มี validation tags ที่ define constraints

**Response DTO**: ใน handler layer represent outgoing response data ให้เป็น clean interface ไปยัง clients

## PII Encryption Pattern

**Encryption Location**: ทั้งหมด encryption/decryption ต้อง happen ใน repository layer เท่านั้น

**pgmapper Utility**: ใช้ pgmapper package functions สำหรับ encrypt/decrypt operations

**Encryption Key**: ดึงมาจาก environment configuration

**Database Storage**: PII fields เก็บเป็น TEXT type ใน database (เมื่อ encrypted)

**Decryption on Read**: เมื่อ read จาก database repository ต้อง decrypt ทั้งหมด PII fields ก่อนส่งไปยัง use case

**Encryption on Write**: เมื่อ write ไปยัง database repository ต้อง encrypt ทั้งหมด PII fields

## Error Handling

**Custom Errors**: ใช้ customerror package สำหรับ define standard error responses

**Error Wrapping**: ทั้ง handler use case และ repository ต้องส่งคืน customerror instances

**PostgreSQL Errors**: ใช้ pgerrutils package สำหรับ handle PostgreSQL-specific errors

**Validation Errors**: ใช้ validatorutils package สำหรับ validate structs

**HTTP Status Codes**: Errors ต้องถูก map ไปยัง appropriate HTTP status codes

## Database Migrations

**Migration Files**: เก็บไฟล์ up/down migrations ใน packages/database/migrations/

**Up Migrations**: สร้าง tables indexes triggers ตามลำดับ dependencies

**Down Migrations**: ลบ objects ในลำดับ reverse ไปยัง respect foreign keys

**Migration Naming**: ใช้ numbered prefix และ descriptive names

**Running Migrations**: ใช้ pnpm db:migrate command สำหรับรัน migrations

## Database Queries & sqlc

**SQL Query Files**: เก็บไฟล์ queries ใน packages/database/queries/

**sqlc Configuration**: ใน packages/database/sqlc.yaml กำหนด how to generate Go code

**Query Annotations**: ใช้ special comments ใน SQL files เพื่อ define query names และ parameters

**Generated Code**: sqlc generates type-safe Go code จาก SQL queries

**Using Generated Code**: import generated queries package ใน repository layer

## API Documentation

**OpenAPI/Swagger**: ทุก handler methods มี Swagger annotations

**Annotations Include**: @Summary @Description @Tags @Accept @Produce @Param @Success @Failure @Router

**Auto-generation**: Swagger documentation ถูก auto-generate จาก annotations

**Access Docs**: API documentation accessible ที่ http://localhost:8080/swagger/

## Services Layer

Services layer ใช้สำหรับ external integrations หรือ utilities

**Pattern**: Services implement specific functionality เช่น OAuth authentication email sending file storage เป็นต้น

**Usage**: UseCase layers เรียก services เพื่อ perform external operations

**Examples**: OAuth service S3 service authentication service เป็นต้น

## Configuration Management

**Config Structure**: ใน internal/config/ มี configuration loading logic

**Environment Variables**: ทั้ง database credentials API keys encryption keys ดึงมาจาก .env file

**Config Access**: Pass config ไปยัง dependencies ที่ต้องการ

## Technology Stack

- **Framework**: Go Fiber (Web framework)
- **Database**: PostgreSQL (Relational database)
- **Query Generator**: sqlc (SQL to Go code generator)
- **Validation**: go-playground/validator (Struct validation)
- **Error Handling**: Custom customerror package
- **Encryption**: AES-256-GCM (Application-layer encryption in Go)
- **Logging**: Structured logging with logrus or zap
- **Package Manager**: Go modules

## ข้อประจำในการสร้าง Backend

- Backend ทั้งหมดใช้ Go Fiber framework
- ปฏิบัติตาม clean architecture pattern ด้วย 3 layers
- Handler layer น้อยและ thin (เฉพาะ HTTP handling)
- UseCase layer ประกอบด้วย business logic
- Repository layer จัดการ database operations และ PII encryption
- ทุก handlers มี complete OpenAPI/Swagger annotations
- ทั้ง input validation ต้อง implement ที่ handler level
- ทั้ง PII fields ต้อง encrypt/decrypt ที่ repository level ใน Go application layer (ไม่ใช่ database level)
- ใช้ dependency injection ไม่มี global state
- ใช้ sqlc generated queries ไม่มี raw SQL
- Error handling ต้อง comprehensive ที่ทุก layers
- Database migrations ต้อง manage ด้วย migration files
