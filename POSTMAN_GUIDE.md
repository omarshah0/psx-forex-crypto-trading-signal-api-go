# Postman Collection Guide

This guide will help you import and use the Postman collection for testing the API.

## 📥 Importing the Collection

### Step 1: Import Collection
1. Open Postman
2. Click **Import** button (top left)
3. Select **File** tab
4. Choose `postman_collection.json` from this directory
5. Click **Import**

### Step 2: Import Environment (Optional but Recommended)
1. Click **Import** button again
2. Select `postman_environment.json`
3. Click **Import**
4. Select the environment from dropdown (top right) - "Go REST API - Local Development"

## 🔧 Setup

### Configure Environment Variables

The collection uses these environment variables:
- `base_url`: API base URL (default: `http://localhost:8080`)
- `access_token`: JWT access token (auto-populated after OAuth login)
- `refresh_token`: JWT refresh token (auto-populated after OAuth login)

You can view/edit these by:
1. Click the eye icon (👁️) next to environment dropdown
2. Click **Edit** to modify values

## 🚀 Testing the API

### 1. Health Check (No Auth Required)

Test if the API is running:
- Request: **Health → Health Check**
- Expected Response: 200 OK with all services healthy

### 2. Authentication Flow

#### Option A: OAuth Login (Browser Required)

Since OAuth requires browser interaction:

1. **Start OAuth Flow:**
   - Copy the URL from **Authentication → Google OAuth - Initiate** or **Facebook OAuth - Initiate**
   - Paste in your browser
   - Complete the OAuth login
   - You'll be redirected back with tokens in cookies

2. **Extract Tokens:**
   After OAuth callback, you'll get a JSON response with tokens. Manually copy:
   - `access_token` → Set in environment variable
   - `refresh_token` → Set in environment variable

#### Option B: Use Existing Tokens

If you already have tokens from a previous session:
1. Click the eye icon (👁️) next to environment
2. Edit `access_token` and `refresh_token`
3. Paste your tokens

### 3. Trading Signals (Authenticated)

Once you have an access token:

#### List All Signals
- Request: **Trading Signals → List All Trading Signals**
- Auth: Automatically uses `{{access_token}}`
- Query params: `limit` and `offset` for pagination

#### Get Single Signal
- Request: **Trading Signals → Get Single Trading Signal**
- Change the ID in URL path as needed

#### Create Signal (Admin Only)
- Request: **Trading Signals → Create Trading Signal (Admin)**
- Multiple examples provided:
  - Basic signal with result
  - LONG signal without result
  - SHORT signal with LOSS
- Requires admin privileges

#### Update Signal (Admin Only)
- Request: **Trading Signals → Update Trading Signal (Admin)**
- Two examples:
  - Partial update (result and return only)
  - Full update (all fields)
- All fields are optional

#### Delete Signal (Admin Only)
- Request: **Trading Signals → Delete Trading Signal (Admin)**
- Change the ID in URL path

### 4. Token Refresh

When your access token expires (15 minutes):

- Request: **Authentication → Refresh Token (with body)**
- Or: **Authentication → Refresh Token (with cookie)**
- Automatically updates `access_token` and `refresh_token` in environment

### 5. Logout

- Request: **Authentication → Logout**
- Revokes refresh token
- You'll need to login again via OAuth

## 🔐 Authentication Methods

### Bearer Token (Mobile/API)

All authenticated requests use Bearer token:
```
Authorization: Bearer {{access_token}}
```

This is configured automatically in the collection.

### Cookie-based (Web)

If testing with cookies:
1. Use browser to complete OAuth flow
2. Cookies are automatically set
3. For cookie-based refresh, use: **Refresh Token (with cookie)**

## 📋 Request Examples

### Create Trading Signal - LONG
```json
{
    "symbol": "BTCUSDT",
    "stop_loss_price": 42000.50,
    "entry_price": 43000.00,
    "take_profit_price": 45000.00,
    "type": "LONG",
    "result": "WIN",
    "return": 4.65
}
```

### Create Trading Signal - SHORT
```json
{
    "symbol": "ETHUSDT",
    "stop_loss_price": 2700.00,
    "entry_price": 2600.00,
    "take_profit_price": 2400.00,
    "type": "SHORT",
    "result": "LOSS",
    "return": -7.69
}
```

### Update Trading Signal (Partial)
```json
{
    "result": "BREAKEVEN",
    "return": 0.00
}
```

## 🎯 Common Workflows

### First Time Setup
1. Health Check → Verify API is running
2. OAuth Login (in browser) → Get tokens
3. Set tokens in environment
4. Test List All Signals → Verify authentication
5. Create/Update/Delete signals (if admin)

### Daily Development
1. Refresh Token (if expired)
2. Test your endpoints
3. Check request logs in MongoDB

### Testing Admin Features
You need to be an admin user:
1. Create user via OAuth
2. Manually add to admins table in database:
   ```sql
   INSERT INTO admins (user_id) VALUES (<your-user-id>);
   ```
3. Login again to get new token with admin status
4. Test admin-only endpoints

## 🔍 Testing Tips

### View Request/Response
- Click on any request in the collection
- Send request
- View response in bottom panel
- Check **Body**, **Headers**, **Cookies** tabs

### Test Scripts
The "Refresh Token" requests include test scripts that automatically:
- Extract tokens from response
- Update environment variables
- Log success messages

### Request Organization
Requests are organized in folders:
- **Health**: No auth required
- **Authentication**: OAuth and token management
- **Trading Signals**: CRUD operations

### Quick Testing
1. **Ctrl/Cmd + Click** on a request to open in new tab
2. **Ctrl/Cmd + Enter** to send request
3. Use **Collections Runner** for automated testing

## 🐛 Troubleshooting

### 401 Unauthorized
- Token expired → Use Refresh Token request
- Token invalid → Login again via OAuth
- Token missing → Set `access_token` in environment

### 403 Forbidden
- Admin required → Check if your user is in admins table
- Login again to get updated token with admin status

### 404 Not Found
- Check URL path (especially IDs)
- Verify API is running on correct port

### 500 Internal Server Error
- Check API logs
- Verify database connections
- Check request body format

## 📊 Response Format

### Success Response
```json
{
  "status": "success",
  "type": "resource",
  "data": { ... },
  "message": "Operation successful"
}
```

### Error Response
```json
{
  "status": "error",
  "type": "bad_request",
  "error": {
    "code": 400,
    "message": "Error description"
  }
}
```

## 🔄 Auto-Update Tokens

The collection includes scripts that automatically update tokens after refresh:

```javascript
// Automatically runs after "Refresh Token" requests
if (pm.response.code === 200) {
    var jsonData = pm.response.json();
    pm.environment.set("access_token", jsonData.data.access_token);
    pm.environment.set("refresh_token", jsonData.data.refresh_token);
}
```

## 📝 Variables Reference

| Variable | Description | Set By |
|----------|-------------|--------|
| `base_url` | API base URL | Manual/Environment |
| `access_token` | JWT access token (15 min) | Auto after OAuth/Refresh |
| `refresh_token` | JWT refresh token (7 days) | Auto after OAuth/Refresh |

## 🚀 Advanced Usage

### Collection Runner
1. Click **Collections** → Three dots → **Run collection**
2. Select requests to run
3. Set iterations and delay
4. View results and test results

### Pre-request Scripts
Add to individual requests or collection level:
```javascript
// Example: Add timestamp to request
pm.environment.set("timestamp", new Date().toISOString());
```

### Tests
Add assertions to validate responses:
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has data", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.exist;
});
```

## 💡 Pro Tips

1. **Organize your environments**: Create separate environments for dev, staging, prod
2. **Use variables**: Use `{{variable}}` syntax for dynamic values
3. **Save examples**: Save successful responses as examples for documentation
4. **Share collections**: Export and share with your team
5. **Monitor APIs**: Set up Postman Monitors for automated testing

## 📚 Additional Resources

- [Postman Documentation](https://learning.postman.com/docs/)
- [API Documentation](API_DOCUMENTATION.md)
- [Quick Start Guide](QUICK_START.md)
- [README](README.md)

---

Happy Testing! 🎉

