# Cursor Rules & Documentation Summary

**Generated**: October 30, 2024

## 📋 Overview

This document summarizes the Cursor Rules and comprehensive PII encryption documentation for the DECM platform. These files provide complete guidance for developing secure, maintainable code.

## 📁 Generated Files

### 1. `.cursorules` (Main Cursor Rules File)

**Location**: `.cursorules`  
**Size**: ~13 KB  
**Purpose**: Comprehensive Cursor IDE rules for the DECM project

**Contents**:

- Project identity and tech stack
- Critical rules (pnpm, PII encryption, architecture)
- Backend 3-layer architecture (Handler → UseCase → Repository)
- Swagger/OpenAPI requirements
- Error handling patterns
- Validation patterns
- SQL & sqlc conventions
- Frontend conventions (React 19)
- Environment configuration
- Development endpoints
- Common commands
- Code quality standards
- Type safety guidelines
- Testing patterns
- Prohibited patterns
- File structure reference
- API-first development workflow
- Troubleshooting guide

### 2. `documentations/pii-encryption.md` (Comprehensive PII Guide)

**Location**: `documentations/pii-encryption.md`  
**Size**: ~20 KB  
**Purpose**: Deep-dive guide on PII encryption implementation

**Contents**:

- Overview and critical principles
- PII fields definition (auth and profile fields)
- Encryption architecture (algorithm, key management, implementation location)
- Encryption flow diagram
- **6 Implementation Patterns**:
    - CREATE operation (insert with encryption)
    - READ operation (select and decrypt)
    - UPDATE operation (encrypt then update)
    - SEARCH operation (deterministic encryption)
    - DELETE operation (no decryption needed)
    - LIST/BATCH operations (decrypt all)
- Available functions (pgmapper & encryptutils)
- Database schema examples
- Environment configuration
- Error handling
- Testing strategies (unit & integration tests)
- Security best practices (DO's and DON'Ts)
- Production considerations (key management, rotation, audit)
- Migration guide for existing data
- Troubleshooting guide
- References

### 3. `documentations/pii-encryption-quick-reference.md` (Quick Reference)

**Location**: `documentations/pii-encryption-quick-reference.md`  
**Size**: ~4.6 KB  
**Purpose**: Quick reference for developers who need to implement encryption fast

**Contents**:

- 10-second summary
- PII fields checklist
- Implementation checklist
- Copy-paste patterns (CREATE, READ, SEARCH, DELETE)
- Function reference table
- Common mistakes with examples
- Setup instructions
- Encryption flow diagram
- Test template
- Troubleshooting table
- Links to full documentation
- Key rules summary

### 4. Updated `CLAUDE.md` (AI Assistant Guide)

**Location**: `CLAUDE.md`  
**Purpose**: Updated to reference new Cursor Rules and documentation

**Changes**:

- Updated Key Cursor Rules section
- References `.cursorules` file
- Direct link to PII encryption guide
- Emphasis on encryption importance

## 🔐 PII Encryption Documentation Architecture

### Documentation Layers

```
└── PII Encryption Documentation
    ├── 1. Quick Reference (for fast implementation)
    │   └── documentations/pii-encryption-quick-reference.md
    │
    ├── 2. Comprehensive Guide (for understanding)
    │   └── documentations/pii-encryption.md
    │
    ├── 3. Cursor Rules (for development)
    │   └── .cursorules (Section 2)
    │
    └── 4. This Summary (orientation)
        └── documentations/CURSOR_RULES_SUMMARY.md
```

### Key Sections in Comprehensive Guide

| Section                   | Purpose                                 |
| ------------------------- | --------------------------------------- |
| Encryption Architecture   | Understand how encryption works         |
| Implementation Patterns   | Copy-paste ready code examples          |
| Available Functions       | Reference for pgmapper functions        |
| Database Schema           | SQL examples with encrypted TEXT fields |
| Error Handling            | How to handle encryption errors         |
| Testing                   | Unit and integration test examples      |
| Security Best Practices   | DO's and DON'Ts with explanations       |
| Production Considerations | Key rotation, audit, monitoring         |
| Migration Guide           | How to encrypt existing data            |
| Troubleshooting           | Common errors and solutions             |

## 🎯 Quick Navigation Guide

### I want to...

**Implement PII encryption in a new repository**
→ Start with `pii-encryption-quick-reference.md` (5 min read)
→ Copy pattern from `pii-encryption.md` (CREATE/READ/SEARCH/etc.)
→ Reference functions in `.cursorules` Section 2

**Understand the full encryption system**
→ Read `pii-encryption.md` from top to bottom (20 min read)
→ Study the 6 implementation patterns
→ Review security best practices section

**Debug encryption issues**
→ Go to `pii-encryption.md` → Troubleshooting section
→ Common errors: "ciphertext too short", "message authentication failed"
→ Check environment variable `PII_ENCRYPTION_KEY`

**Set up project for PII encryption**
→ Generate key: `openssl rand -base64 32`
→ Add to `.env`: `PII_ENCRYPTION_KEY=<key>`
→ Initialize repository with encryption key

**Review project conventions**
→ Read `.cursorules` file (comprehensive)
→ Sections 1-15 cover all major topics
→ Prohibited patterns at the end

**Migrate existing unencrypted data**
→ `pii-encryption.md` → Migration Guide: Existing Data
→ Ready-to-use migration function code

## 🏗️ Architecture Alignment

### Backend Architecture (3-Layer Pattern)

```
Handler Layer (HTTP)
    ↓ (parse/validate)
UseCase Layer (Business Logic)
    ↓ (coordinate)
Repository Layer (Data Access)
    ├─→ Encrypt PII before INSERT/UPDATE
    ├─→ Execute SQL queries
    └─→ Decrypt PII after SELECT
```

### Encryption is Repository Layer Only ✅

**NOT in Handler**: Parse/validate input only  
**NOT in UseCase**: Business logic only  
**YES in Repository**: Encrypt before INSERT/UPDATE, Decrypt after SELECT

## 🔑 Key Rules Summary

### Absolute Rules ⚠️

1. **Use pnpm only** - Never npm/yarn/bun
2. **Encrypt PII at repository layer** - Never skip
3. **Use pgmapper utilities** - Never implement own encryption
4. **Get key from environment** - Never hardcode
5. **NO encryption in SQL** - All in Go code
6. **Handle errors properly** - Don't silently ignore
7. **Always add Swagger annotations** - Required on all endpoints
8. **Never use `dangerouslySetInnerHTML`** - Use Typography component

### AES-256-GCM Encryption ✅

- **Algorithm**: AES-256-GCM (authenticated encryption)
- **Key Size**: 32 bytes (256 bits)
- **Nonce**: 12 bytes, randomly generated per encryption
- **Output**: Base64-encoded for safe transmission
- **Authentication**: Detects any data tampering

## 📚 File Cross-References

| Need            | File                              | Section                    |
| --------------- | --------------------------------- | -------------------------- |
| Quick start     | pii-encryption-quick-reference.md | Copy-Paste Patterns        |
| Full guide      | pii-encryption.md                 | Implementation Patterns    |
| Rules           | .cursorules                       | Section 2 (PII Encryption) |
| Examples        | pii-encryption.md                 | CREATE/READ/UPDATE/SEARCH  |
| Testing         | pii-encryption.md                 | Testing Encryption         |
| Troubleshooting | pii-encryption.md                 | Troubleshooting            |
| Setup           | pii-encryption-quick-reference.md | Setup                      |
| Production      | pii-encryption.md                 | Production Considerations  |

## 📊 PII Fields Covered

### Authentication (2 fields)

- `google_connector_ref`
- `github_connector_ref`

### Profile (9 fields)

- `first_name`
- `last_name`
- `email`
- `phone_number`
- `address`
- `bio`
- `profile_picture_url`
- `academic_institution`
- `academic_email`

## 🛠️ Development Workflow

### When Starting New Feature

1. **Read `.cursorules`** (5 min) - Understand project conventions
2. **Plan Architecture** (5 min) - Handler → UseCase → Repository
3. **Identify PII Fields** (2 min) - Which fields need encryption?
4. **Reference `pii-encryption.md`** (5 min) - Find matching pattern
5. **Implement Repository** (10 min) - Copy pattern and adapt
6. **Write Tests** (5 min) - Use test template from docs
7. **Add Swagger** (5 min) - Document API endpoint

### Environment Setup

```bash
# Generate encryption key
openssl rand -base64 32
# → 7kF3mQ9pL2xN8vR5wT1jB4cD6hE9sG7k

# Add to .env
echo 'PII_ENCRYPTION_KEY=7kF3mQ9pL2xN8vR5wT1jB4cD6hE9sG7k' >> .env

# Verify it works
cd apps/backend && go test ./...
```

## 📖 Reading Order Recommendation

### First Time Users

1. `.cursorules` - Section 1 (Project Identity)
2. `.cursorules` - Section 2 (PII Encryption)
3. `pii-encryption-quick-reference.md` - Full document
4. `pii-encryption.md` - One implementation pattern

### Regular Development

1. `.cursorules` - Section 2 (PII Encryption) - Refresh
2. `pii-encryption-quick-reference.md` - Find pattern
3. `pii-encryption.md` - Reference for details

### Complex Tasks

1. `.cursorules` - Relevant sections
2. `pii-encryption.md` - Full document
3. `CLAUDE.md` - Project overview

## ✅ Quality Checklist

Before committing repository code:

- [ ] All PII fields identified and documented
- [ ] Encryption at repository layer only
- [ ] Using `pgmapper` functions (not custom encryption)
- [ ] Encrypt before INSERT/UPDATE
- [ ] Decrypt after SELECT
- [ ] Error handling on encryption operations
- [ ] Tested encrypt/decrypt roundtrip
- [ ] Environment variable configured
- [ ] No hardcoded keys
- [ ] Swagger documentation complete

## 🔗 File Locations

```
decm/
├── .cursorules                                    ← Main Cursor Rules
├── documentations/
│   ├── pii-encryption.md                        ← Comprehensive Guide
│   ├── pii-encryption-quick-reference.md        ← Quick Start
│   └── CURSOR_RULES_SUMMARY.md                  ← This File
├── apps/backend/common/
│   ├── pgmapper/                                ← Encryption Utilities
│   └── encryptutils/                            ← Low-level Functions
└── .env.example                                 ← Config Template
```

## 📝 Document Maintenance

### How to Update Documentation

1. **New Encryption Pattern?**
    - Add to `pii-encryption.md` → Implementation Patterns section
    - Add quick reference to `pii-encryption-quick-reference.md`
    - Update `.cursorules` → Section 2 if needed

2. **New Project Rule?**
    - Add to `.cursorules` with clear section
    - Reference in `CLAUDE.md` if critical
    - Update project guidelines

3. **New PII Field?**
    - Add to all three docs in PII Fields Definition
    - Create sample implementation in `pii-encryption.md`
    - Add to quick reference

## 🎓 Training Resources

### For New Developers

1. Read `.cursorules` (main project guide)
2. Read `pii-encryption-quick-reference.md` (5 min)
3. Copy pattern from `pii-encryption.md` (pick one)
4. Implement with guidance from docs
5. Test with template from `pii-encryption.md`

### For Security Review

1. Audit `.cursorules` → Section 2
2. Review `pii-encryption.md` → Security Best Practices
3. Check `pii-encryption.md` → Production Considerations
4. Review all PII fields identified

### For Architecture Review

1. Review `.cursorules` → Section 3 (Backend Architecture)
2. Check `pii-encryption.md` → Encryption Flow Diagram
3. Verify layering (Handler → UseCase → Repository)
4. Confirm encryption at repository layer only

## 🎯 Success Metrics

After reading and understanding these documents, you should be able to:

✅ **Identify PII fields** in any entity
✅ **Implement encryption** using pgmapper functions
✅ **Write tests** for encryption/decryption
✅ **Debug encryption issues** using troubleshooting guide
✅ **Plan key rotation** for production
✅ **Follow project conventions** from `.cursorules`
✅ **Structure backend code** using Handler → UseCase → Repository pattern
✅ **Handle errors properly** with customerror package

## 📞 Quick Reference Links

- **Cursor Rules**: `.cursorules` (all project conventions)
- **Encryption Guide**: `documentations/pii-encryption.md` (comprehensive)
- **Quick Start**: `documentations/pii-encryption-quick-reference.md` (5-min version)
- **Implementation**: `apps/backend/common/pgmapper/` (actual code)
- **Configuration**: `.env.example` (environment setup)
- **AI Assistant**: `CLAUDE.md` (for Claude in Cursor)

---

**Last Updated**: October 30, 2024  
**Status**: Complete ✅  
**Coverage**: PII Encryption + Project Conventions
