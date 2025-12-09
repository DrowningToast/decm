# Certificate Claim - Testing Guide

## 🧪 Quick Testing Guide

### Prerequisites

1. ✅ Backend server running (`pnpm dev:core`)
2. ✅ Frontend server running (`pnpm dev`)
3. ✅ Database with test data
4. ✅ User authenticated with valid account

---

## 📋 Test Scenarios

### Scenario 1: Successful Certificate Claim ✅

**Setup**:

- User has an unclaimed certificate
- Certificate is published (has `event_certificate_address`)
- User is authenticated
- User knows their account password

**Steps**:

1. Navigate to `/app/certificates`
2. Click on an unclaimed certificate (no token ID)
3. Verify certificate detail page loads:
    - Certificate image displays
    - Certificate name and event name shown
    - "Claim Certificate" button visible
4. Click "Claim Certificate" button
5. Password prompt modal opens
6. Enter valid account password
7. Click submit

**Expected Result**:

- ✅ Button shows "Claiming..." with spinner
- ✅ Success toast appears: "Certificate claimed successfully!"
- ✅ Certificate status updates to "completed"
- ✅ "Certificate Claimed" badge appears
- ✅ "Claim Certificate" button replaced with claimed badge
- ✅ Certificate list refreshes automatically

---

### Scenario 2: Wrong Password ❌

**Steps**:

1. Navigate to unclaimed certificate detail page
2. Click "Claim Certificate"
3. Enter **incorrect** password
4. Click submit

**Expected Result**:

- ❌ Error toast appears with backend error message
- ❌ Certificate remains unclaimed
- ❌ Can retry claiming

---

### Scenario 3: Already Claimed Certificate ℹ️

**Steps**:

1. Navigate to a certificate that's already claimed
2. Verify page state

**Expected Result**:

- ℹ️ No "Claim Certificate" button visible
- ℹ️ "Certificate Claimed" badge displayed
- ℹ️ Certificate shows "Certificate issued on [date]"

---

### Scenario 4: Unauthenticated User 🚫

**Steps**:

1. Log out
2. Try to access `/app/certificates/{id}`

**Expected Result**:

- 🚫 Redirected to login page
- 🚫 Cannot access certificate detail

---

### Scenario 5: Certificate Not Published ⚠️

**Setup**:

- Certificate config exists
- Certificate NOT published (no contract address)

**Steps**:

1. Try to claim certificate

**Expected Result**:

- ⚠️ Error from backend
- ⚠️ Toast shows appropriate error message

---

## 🔍 Things to Check

### Visual Verification

- [ ] Certificate image loads correctly
- [ ] Loading spinner shows during claim
- [ ] Button text changes during states
- [ ] Toast notifications appear
- [ ] Claimed badge has correct styling
- [ ] Layout is responsive on mobile

### Functional Verification

- [ ] API call is made to correct endpoint
- [ ] Password is sent securely
- [ ] Certificate list refreshes after claim
- [ ] Inbox messages update after claim
- [ ] Multiple rapid clicks don't cause double-claim
- [ ] Error handling works correctly

### Network Verification (DevTools)

```
Network Tab:
1. Click claim button
2. Look for POST request to:
   /api/v1/certificates/claim/{certificate_id}
3. Request payload should contain:
   { "account_password": "..." }
4. Response should be 200 OK with certificate object
```

### Console Verification

```
Check browser console:
- No errors during claim process
- Success logs after claim
- Proper error logs if claim fails
```

---

## 🐛 Common Issues & Solutions

### Issue 1: "Failed to claim certificate" error

**Possible Causes**:

1. Wrong password → Re-enter correct password
2. Certificate already claimed → Check certificate status
3. Certificate not published → Verify event certificate config
4. User not in event_attendees → Check event registration
5. Backend server not running → Start backend server

**Debug Steps**:

```bash
# Check backend logs
cd apps/backend
# Look for error messages in terminal

# Check network tab in browser DevTools
# Look at response error message
```

---

### Issue 2: Button stuck in "Claiming..." state

**Possible Causes**:

1. Backend request timed out
2. Blockchain transaction pending
3. Network error

**Solution**:

```
1. Refresh the page
2. Check if certificate was actually claimed
3. Look at backend logs for transaction status
4. Check blockchain explorer for transaction
```

---

### Issue 3: Certificate list doesn't refresh

**Possible Causes**:

1. Query invalidation not working
2. Cache issue

**Solution**:

```
1. Manually refresh the page
2. Check if queryClient.invalidateQueries is called
3. Verify query keys match
```

---

## 📊 Backend Verification

### Check Database After Claim

```sql
-- Verify certificate was updated with token ID
SELECT
  id,
  name,
  certificate_token_id,
  claimed_at
FROM event_certificates
WHERE id = 'certificate-id';

-- Expected: certificate_token_id should be populated
-- Expected: claimed_at should have a timestamp
```

### Check Blockchain Transaction

```bash
# Get transaction hash from backend logs
# Or from API response

# Visit block explorer
https://sepolia.etherscan.io/tx/{transaction_hash}

# Verify:
- Transaction successful
- Contract address matches
- MintNft function called
- Correct token ID returned
```

---

## 🎯 Testing Checklist

### Basic Flow

- [ ] View certificate list
- [ ] Click unclaimed certificate
- [ ] Certificate detail loads correctly
- [ ] Click "Claim Certificate" button
- [ ] Password prompt appears
- [ ] Enter password
- [ ] Claim succeeds
- [ ] Success toast shows
- [ ] Certificate updates to claimed
- [ ] Badge appears

### Error Cases

- [ ] Wrong password shows error
- [ ] Already claimed certificate handled
- [ ] Unpublished certificate handled
- [ ] Unauthenticated user redirected

### UI/UX

- [ ] Loading states work
- [ ] Responsive on mobile
- [ ] Translations correct (EN + TH)
- [ ] Accessibility features work

### Performance

- [ ] API calls efficient
- [ ] No memory leaks
- [ ] Query caching works
- [ ] Image loading optimized

---

## 🚀 API Testing with cURL

### Get Sign Message

```bash
curl -X GET \
  "http://localhost:8080/api/v1/certificates/claim/{certificate_id}/sign-message" \
  -H "Cookie: jwt=YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json"
```

### Claim with Password

```bash
curl -X POST \
  "http://localhost:8080/api/v1/certificates/claim/{certificate_id}" \
  -H "Cookie: jwt=YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "account_password": "your-password"
  }'
```

### Expected Response

```json
{
    "id": "cert-123",
    "name": "Workshop Certificate",
    "certificate_token_id": "456",
    "event_name": "Web3 Workshop",
    "event_id": "event-789",
    "receiver_credential_id": "cred-123",
    "receiver_name": "John Doe",
    "receiver_email": "john@example.com",
    "created_at": "2025-12-08T10:30:00Z",
    "claimed_at": "2025-12-08T15:45:00Z"
    // ... other fields
}
```

---

## 📝 Test Data Setup

### Create Test Certificate

```sql
-- Ensure event exists
-- Ensure user is in event_attendees
-- Ensure certificate config is published
-- Create unclaimed certificate
INSERT INTO event_certificates (
  id,
  event_id,
  receiver_credential_id,
  name,
  certificate_token_id -- NULL for unclaimed
) VALUES (
  gen_random_uuid(),
  'event-id',
  'user-credential-id',
  'Test Certificate',
  NULL
);
```

---

## ✅ Success Criteria

A successful test run should demonstrate:

1. ✅ User can view their certificates
2. ✅ User can claim unclaimed certificates
3. ✅ Password verification works
4. ✅ NFT is minted on blockchain
5. ✅ Database is updated correctly
6. ✅ UI updates in real-time
7. ✅ Error handling works properly
8. ✅ No console errors
9. ✅ Responsive on all devices
10. ✅ Translations work in both languages

---

## 🎉 Happy Testing!

If all test scenarios pass, the certificate claiming feature is working correctly and ready for production deployment! 🚀
