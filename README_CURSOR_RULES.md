# Cursor Rules & Documentation - README

**Status**: ✅ **COMPLETE**  
**Generated**: October 30, 2024  
**Total Coverage**: 1,682 lines across 4 files

## 🎯 What Was Generated?

Comprehensive Cursor IDE rules and documentation for the DECM platform, with special focus on **PII encryption**.

### 📁 4 New/Updated Files

| File                                                   | Size      | Purpose                                  |
| ------------------------------------------------------ | --------- | ---------------------------------------- |
| **`.cursorules`**                                      | 435 lines | Main Cursor rules for the entire project |
| **`documentations/pii-encryption.md`**                 | 700 lines | Comprehensive PII encryption guide       |
| **`documentations/pii-encryption-quick-reference.md`** | 172 lines | Quick 5-minute reference guide           |
| **`documentations/CURSOR_RULES_SUMMARY.md`**           | 375 lines | Overview and navigation guide            |

## 🚀 Quick Start (3 Steps)

### Step 1: Understand the Rules

```bash
# Open and read the main Cursor Rules
cat .cursorules  # 435 lines - all project conventions
```

### Step 2: Understand PII Encryption (Pick One)

```bash
# Option A: Quick version (5 min read)
cat documentations/pii-encryption-quick-reference.md

# Option B: Comprehensive version (20 min read)
cat documentations/pii-encryption.md
```

### Step 3: Use as Reference

```bash
# When implementing new features
# → Reference .cursorules for architecture
# → Copy patterns from pii-encryption.md
# → Use environment variable PII_ENCRYPTION_KEY
```

## 📚 Documentation Structure

### Layer 1: Main Cursor Rules (`.cursorules`)

**What**: Comprehensive project conventions  
**Who**: All developers  
**When**: Reference while developing  
**15 Sections**:

1. Project Identity
2. PII Encryption ⭐
3. Backend Architecture
4. Swagger/OpenAPI
5. Error Handling
6. Validation Pattern
7. SQL & sqlc Conventions
8. Frontend Conventions
9. Environment Configuration
10. Development Endpoints
11. Common Commands
12. Code Quality
13. Type Safety
14. Testing
15. Prohibited Patterns

### Layer 2: PII Encryption Documentation

#### 2a. Comprehensive Guide (`documentations/pii-encryption.md`)

**What**: Deep-dive on encryption  
**Who**: Developers implementing PII fields  
**When**: Implementing new repository methods  
**Key Sections**:

- Critical principles
- PII fields definition (11 fields)
- Architecture explanation
- 6 Implementation patterns (CREATE, READ, UPDATE, SEARCH, DELETE, BATCH)
- Available functions
- Database schema
- Error handling
- Testing strategies
- Security best practices
- Production considerations
- Migration guide
- Troubleshooting

#### 2b. Quick Reference (`documentations/pii-encryption-quick-reference.md`)

**What**: Fast implementation guide  
**Who**: Developers who need to code NOW  
**When**: Implementing encryption quickly  
**Key Sections**:

- 10-second summary
- Copy-paste patterns
- Function reference
- Common mistakes
- Setup instructions
- Troubleshooting table

### Layer 3: Navigation Guide (`documentations/CURSOR_RULES_SUMMARY.md`)

**What**: Overview of all documentation  
**Who**: First-time users, project leads  
**When**: Getting oriented or planning  
**Key Sections**:

- File summary
- Architecture layers
- Quick navigation
- Reading order
- Cross-references
- Training resources

### Layer 4: Updated (`CLAUDE.md`)

**What**: AI Assistant context  
**Who**: Claude in Cursor IDE  
**When**: Always loaded in IDE  
**Updated**: Section "Key Cursor Rules"

## 🔐 PII Encryption Focus

### 11 PII Fields Covered

**Authentication (2)**:

- `google_connector_ref`
- `github_connector_ref`

**Profile (9)**:

- `email`
- `first_name` / `last_name`
- `phone_number`
- `address`
- `bio`
- `profile_picture_url`
- `academic_institution`
- `academic_email`

### 3 Implementation Approaches

| Approach           | Time   | Use Case                             |
| ------------------ | ------ | ------------------------------------ |
| 📖 Read Full Guide | 20 min | Want to understand encryption deeply |
| ⚡ Quick Reference | 5 min  | Want to implement fast               |
| 🔍 Copy Pattern    | 2 min  | Know what pattern you need           |

### 6 Implementation Patterns

1. **CREATE** - Insert with encryption
2. **READ** - Select and decrypt
3. **UPDATE** - Encrypt then update
4. **SEARCH** - Deterministic encryption
5. **DELETE** - No decryption needed
6. **BATCH** - Decrypt multiple rows

## 🎯 Use Cases & Navigation

### Use Case 1: "I'm new to the project"

→ Start here:

1. `.cursorules` Section 1 (Project Identity) - 5 min
2. `.cursorules` Section 2 (PII Encryption) - 10 min
3. `pii-encryption-quick-reference.md` - 5 min
4. Read one pattern from `pii-encryption.md` - 10 min

**Time**: 30 minutes

### Use Case 2: "I need to implement profile encryption"

→ Start here:

1. `pii-encryption-quick-reference.md` - 5 min
2. Copy pattern from `pii-encryption.md` - 2 min
3. Adapt to your entity - 10 min
4. Write tests - 5 min

**Time**: 22 minutes

### Use Case 3: "I'm debugging encryption errors"

→ Go to:

1. `pii-encryption.md` → Troubleshooting - 5 min
2. Check environment variable - 1 min
3. Verify implementation - 5 min

**Time**: 11 minutes

### Use Case 4: "I need to set up for production"

→ Read:

1. `pii-encryption.md` → Production Considerations - 10 min
2. `pii-encryption.md` → Migration Guide - 10 min
3. Plan key rotation - 5 min

**Time**: 25 minutes

## 📊 Encryption Implementation at a Glance

```
Repository Layer (ONLY place for encryption)
├─ Before INSERT/UPDATE
│  ├─ Encrypt PII with: pgmapper.EncryptStringPtrToPgText()
│  ├─ Use key from: r.piiEncryptionKey (environment variable)
│  └─ Algorithm: AES-256-GCM
│
├─ Execute SQL query
│  └─ Data stored as encrypted TEXT
│
└─ After SELECT
   ├─ Decrypt PII with: pgmapper.DecryptPgTextToStringPtr()
   ├─ Use same key: r.piiEncryptionKey
   └─ Return decrypted to UseCase
```

## 🔧 One-Time Setup

```bash
# 1. Generate 32-byte encryption key
openssl rand -base64 32
# Output: REDACTED_PII_ENCRYPTION_KEY

# 2. Add to .env file
echo 'PII_ENCRYPTION_KEY=REDACTED_PII_ENCRYPTION_KEY' >> .env

# 3. Verify configuration
cat .env | grep PII_ENCRYPTION_KEY

# 4. Initialize repository with key
# (Usually done in main.go or config loading)
```

## 📋 Developer Checklist

Before committing repository code with PII:

- [ ] Read `.cursorules` Section 2 ✅
- [ ] Identified all PII fields ✅
- [ ] Using pgmapper functions (not custom) ✅
- [ ] Encrypting before INSERT/UPDATE ✅
- [ ] Decrypting after SELECT ✅
- [ ] Handling encryption errors ✅
- [ ] Testing encrypt/decrypt roundtrip ✅
- [ ] Environment variable set ✅
- [ ] No hardcoded keys ✅
- [ ] Swagger docs complete ✅

## 🏗️ Architecture Principles

### 3-Layer Backend Pattern ✅

```
Handler Layer (HTTP)
    ↓ parse & validate
UseCase Layer (Business Logic)
    ↓ coordinate repositories
Repository Layer (Data Access)
    ├─ Encrypt/Decrypt PII ⭐
    └─ Database operations
```

### PII Encryption Location ⭐

- ❌ NOT in Handler (validates only)
- ❌ NOT in UseCase (business logic only)
- ✅ YES in Repository (data access layer)

## 🔑 Critical Rules

### Absolute (Non-negotiable)

1. Use **pnpm only** (never npm/yarn/bun)
2. Encrypt **PII at repository layer** (never skip)
3. Use **pgmapper utilities** (never implement your own)
4. Get key from **environment variable** (never hardcode)
5. **No encryption in SQL** (all in Go code)

### Best Practices

6. Handle **encryption errors** properly
7. Add **Swagger annotations** on all endpoints
8. Test **encrypt/decrypt roundtrips**
9. Document **PII fields clearly**
10. Plan **key rotation** for production

## 📞 Quick Links

| Need             | File                                | Section                   |
| ---------------- | ----------------------------------- | ------------------------- |
| Project rules    | `.cursorules`                       | All sections              |
| Encryption guide | `pii-encryption.md`                 | All sections              |
| Quick start      | `pii-encryption-quick-reference.md` | Copy-Paste Patterns       |
| Navigation       | `CURSOR_RULES_SUMMARY.md`           | Quick Navigation          |
| Code examples    | `pii-encryption.md`                 | CREATE/READ/UPDATE/SEARCH |
| Troubleshooting  | `pii-encryption.md`                 | Troubleshooting section   |
| Setup            | `pii-encryption-quick-reference.md` | Setup section             |
| Production       | `pii-encryption.md`                 | Production Considerations |

## 📈 Documentation Statistics

| Metric                  | Count |
| ----------------------- | ----- |
| Total Lines             | 1,682 |
| Files                   | 4     |
| Sections                | 15+   |
| Implementation Patterns | 6     |
| PII Fields Documented   | 11    |
| Code Examples           | 30+   |
| Test Templates          | 4     |
| Best Practices          | 20+   |

## 🎓 For Different Roles

### For Developers

1. Start with `.cursorules`
2. Reference `pii-encryption.md` while coding
3. Use `pii-encryption-quick-reference.md` for quick lookups

### For Architects

1. Review `.cursorules` Section 3 (Architecture)
2. Study `pii-encryption.md` (Encryption Flow)
3. Check `CURSOR_RULES_SUMMARY.md` (Architecture Alignment)

### For Security Team

1. Audit `.cursorules` Section 2
2. Review `pii-encryption.md` (Security Best Practices)
3. Check Production Considerations

### For Team Leads

1. Use `CURSOR_RULES_SUMMARY.md` for overview
2. Reference `documentations/CURSOR_RULES_SUMMARY.md` for training plan
3. Share `.cursorules` with team

## ✅ Success Criteria

After using these documents, you should be able to:

✅ Identify PII fields that need encryption  
✅ Implement encryption using pgmapper  
✅ Write tests for encryption/decryption  
✅ Debug encryption issues  
✅ Plan key rotation for production  
✅ Follow project conventions  
✅ Structure code using Handler → UseCase → Repository  
✅ Handle errors with customerror package

## 🚀 Next Steps

1. **Read `.cursorules`** (start here - 30 min)
2. **Read `pii-encryption-quick-reference.md`** (5 min)
3. **Copy pattern from `pii-encryption.md`** (as needed)
4. **Implement with guidance** (test + document)
5. **Refer to `CURSOR_RULES_SUMMARY.md`** (for navigation)

## 📝 File Locations Summary

```
decm/
├── .cursorules                              ← START HERE
├── README_CURSOR_RULES.md                   ← This file
├── CLAUDE.md                                ← Updated
└── documentations/
    ├── pii-encryption.md                    ← Comprehensive
    ├── pii-encryption-quick-reference.md    ← Quick start
    ├── CURSOR_RULES_SUMMARY.md              ← Navigation
    ├── backend-architecture.md              ← Context
    └── ... other docs
```

---

**Created**: October 30, 2024  
**Status**: ✅ Complete  
**Ready to Use**: Yes  
**Maintenance**: See `CURSOR_RULES_SUMMARY.md` → Document Maintenance
