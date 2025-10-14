# Google Maps API - Fixing 403 Error

## 🚨 The 403 Error You're Seeing

The error `Failed to load resource: the server responded with a status of 403` means Google Maps API is rejecting your request. This is almost always due to **API key restrictions**.

---

## ✅ Quick Fix Steps

### Step 1: Check Your API Key Restrictions

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Select your project
3. Navigate to **"APIs & Services" → "Credentials"**
4. Find your API key and click the **pencil icon** to edit it

### Step 2: Configure HTTP Referrers (Most Important!)

Under **"Application restrictions"**:

1. Select **"HTTP referrers (web sites)"**
2. Click **"Add an item"**
3. Add these referrers:
    ```
    localhost:3000/*
    127.0.0.1:3000/*
    localhost:*/*
    127.0.0.1:*/*
    ```
4. **Important**: Add each one separately!

### Step 3: Configure API Restrictions

Under **"API restrictions"**:

1. Select **"Restrict key"**
2. From the dropdown, select **ONLY**:
    - ✅ Maps Embed API
3. Uncheck everything else
4. Click **"Save"**

### Step 4: Wait for Changes to Propagate

⏳ **Important:** Changes can take 1-5 minutes to take effect.

After saving, wait 2-3 minutes before testing again.

---

## 🔍 Verify Your Current Settings

Your current API key (visible in the error): `AIzaSyCXYP0bmOQONa_HedUayJ-B8QtYd6zIUwA`

### Check These Settings in Google Cloud Console:

**Application restrictions:**

- ✅ Should be: "HTTP referrers"
- ❌ NOT: "None" or "IP addresses"

**HTTP referrers should include:**

```
localhost:3000/*
127.0.0.1:3000/*
localhost:*/*
127.0.0.1:*/*
```

**API restrictions:**

- ✅ "Restrict key" is selected
- ✅ "Maps Embed API" is checked
- ❌ No other APIs should be checked

---

## 🧪 Testing After Changes

### 1. Clear Browser Cache

```bash
# In Chrome/Edge: Ctrl+Shift+Delete (Windows) or Cmd+Shift+Delete (Mac)
# Or use Incognito/Private mode
```

### 2. Wait 2-3 Minutes

Google's changes need time to propagate across their servers.

### 3. Refresh Your App

```bash
# Hard refresh in browser
Ctrl+Shift+R (Windows) or Cmd+Shift+R (Mac)
```

### 4. Check if Map Loads

- Navigate to your event form
- Go to Step 2 (Venue Information)
- Enter a location (e.g., "Central World Bangkok")
- Map should now appear!

---

## 🔧 Common Issues & Solutions

### Issue 1: Still Getting 403 After Restrictions

**Possible causes:**

- Changes haven't propagated yet (wait 5 minutes)
- Wrong API is restricted (make sure it's "Maps Embed API")
- Typo in referrer URL (check for extra spaces)

**Solution:**

1. Double-check the HTTP referrers list
2. Make sure you clicked "Save"
3. Wait 5 full minutes
4. Try in Incognito mode

### Issue 2: "This page can't load Google Maps correctly"

**Cause:** Billing not set up or API not enabled

**Solution:**

1. Go to "APIs & Services" → "Library"
2. Search for "Maps Embed API"
3. Make sure it shows "API Enabled" (green)
4. Go to "Billing" and verify billing account is linked

### Issue 3: "For development purposes only" watermark

**Cause:** API key issue or billing issue

**Solution:**

1. Verify billing is set up
2. Check that your API key has no quotas exceeded
3. Try creating a new API key

---

## 🎯 Correct Configuration Example

Here's what your API key settings should look like:

```
API Key: AIzaSyCXYP0bmOQONa_HedUayJ-B8QtYd6zIUwA

Application restrictions:
  Type: HTTP referrers (web sites)
  Referrers:
    - localhost:3000/*
    - 127.0.0.1:3000/*
    - localhost:*/*
    - 127.0.0.1:*/*
    - your-domain.com/*  (add when deploying)

API restrictions:
  Type: Restrict key
  APIs:
    ✅ Maps Embed API

Key restrictions:
  ⚠️ None (leave blank)
```

---

## 📋 Quick Checklist

Before you continue, verify:

- [ ] API key is added to `apps/web/.env.local`
- [ ] "Maps Embed API" is enabled in Google Cloud Console
- [ ] Billing is set up (even for free tier)
- [ ] HTTP referrers include `localhost:3000/*`
- [ ] API restrictions include ONLY "Maps Embed API"
- [ ] You've waited 2-3 minutes after saving changes
- [ ] Browser cache is cleared
- [ ] Dev server is restarted (`pnpm dev`)

---

## 🆘 Still Not Working?

### Option 1: Create a New API Key

Sometimes it's easier to start fresh:

1. Go to "APIs & Services" → "Credentials"
2. Click "Create Credentials" → "API Key"
3. Configure restrictions immediately (see Step 2 above)
4. Copy the new key to your `.env.local` file
5. Restart dev server

### Option 2: Use a Less Restrictive Configuration (Temporarily)

**⚠️ WARNING: Only for testing! Not for production!**

Temporarily set:

- Application restrictions: **"None"**
- API restrictions: **"Don't restrict key"**

If this works, then you know it's a restriction issue. Re-add restrictions one by one to find the problem.

**Remember to re-enable restrictions before deploying!**

---

## 📚 Additional Resources

- [Google Maps Embed API Documentation](https://developers.google.com/maps/documentation/embed/get-started)
- [API Key Best Practices](https://developers.google.com/maps/api-security-best-practices)
- [Troubleshooting Guide](https://developers.google.com/maps/documentation/embed/get-started#troubleshooting)

---

## 💡 Pro Tip

If you're still having issues, check the browser's **Console** tab in Developer Tools for more detailed error messages. They often provide specific information about why the request was rejected.

**Your current error showed:**

```
Failed to load resource: the server responded with a status of 403
```

This specifically means your API key restrictions are blocking the request. Follow Step 2 above carefully!

---

**Need More Help?**

If you've followed all steps and it's still not working:

1. Take a screenshot of your API key restrictions page
2. Check the browser console for detailed error messages
3. Verify the API key in `.env.local` matches the one in Google Cloud Console
