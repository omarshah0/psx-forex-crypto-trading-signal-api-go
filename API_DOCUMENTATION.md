# API Documentation

Complete API documentation for the Go REST API with Social Authentication.

## Base URL

```
http://localhost:8080
```

## Authentication

The API supports two authentication methods:

### 1. Cookie-based (Web Applications)
Tokens are automatically stored in HTTP-only cookies after OAuth login.

### 2. Bearer Token (Mobile Applications)
Include the access token in the Authorization header:
```
Authorization: Bearer <access_token>
```

## Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "status": "success",
  "type": "resource" | "collection" | "action" | "auth",
  "data": { ... },
  "message": "Operation completed successfully"
}
```

### Error Response
```json
{
  "status": "error",
  "type": "bad_request" | "unauthorized" | "forbidden" | "not_found" | "validation_error" | "internal_server_error",
  "error": {
    "code": 400,
    "message": "Error description"
  }
}
```

## Endpoints

### Health Check

#### GET /health
Check the health status of all services (PostgreSQL, MongoDB, Redis).

**Authentication:** Not required

**Response:**
```json
{
  "status": "success",
  "type": "action",
  "data": {
    "postgres": "healthy",
    "mongodb": "healthy",
    "redis": "healthy"
  },
  "message": "All services are healthy"
}
```

---

## Authentication Endpoints

### Google OAuth

#### GET /auth/google
Initiates the Google OAuth flow. Redirects to Google's authorization page.

**Authentication:** Not required

**Response:** 302 Redirect to Google OAuth

---

#### GET /auth/google/callback
Google OAuth callback endpoint. Handles the OAuth response from Google.

**Authentication:** Not required

**Query Parameters:**
- `code` (string, required): Authorization code from Google
- `state` (string, required): CSRF protection token

**Response:**
```json
{
  "status": "success",
  "type": "auth",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "blocked": false,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "is_admin": false
  },
  "message": "Authentication successful"
}
```

**Cookies Set:**
- `access_token`: JWT access token (15 minutes)
- `refresh_token`: JWT refresh token (7 days)

---

### Facebook OAuth

#### GET /auth/facebook
Initiates the Facebook OAuth flow. Redirects to Facebook's authorization page.

**Authentication:** Not required

**Response:** 302 Redirect to Facebook OAuth

---

#### GET /auth/facebook/callback
Facebook OAuth callback endpoint. Handles the OAuth response from Facebook.

**Authentication:** Not required

**Query Parameters:**
- `code` (string, required): Authorization code from Facebook
- `state` (string, required): CSRF protection token

**Response:** Same as Google callback

**Cookies Set:** Same as Google callback

---

### Code Exchange (Frontend-Initiated OAuth)

These endpoints are designed for modern web and mobile applications where the frontend handles the OAuth flow and exchanges the authorization code for JWT tokens.

#### POST /auth/google/exchange
Exchange Google authorization code for JWT tokens.

**Authentication:** Not required

**Request Body:**
```json
{
  "code": "4/0AX4XfWh..."
}
```

**Response:**
```json
{
  "status": "success",
  "type": "auth",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "email_verified": true,
      "blocked": false,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "is_admin": false,
    "expires_in": 900
  },
  "message": "Authentication successful"
}
```

**Error Responses:**
- `403 Forbidden`: Google OAuth is not enabled
- `400 Bad Request`: Missing or invalid code
- `401 Unauthorized`: Invalid authorization code or authentication failed

---

#### POST /auth/facebook/exchange
Exchange Facebook authorization code for JWT tokens.

**Authentication:** Not required

**Request Body:**
```json
{
  "code": "AQB..."
}
```

**Response:** Same format as Google exchange

**Error Responses:**
- `403 Forbidden`: Facebook OAuth is not enabled
- `400 Bad Request`: Missing or invalid code
- `401 Unauthorized`: Invalid authorization code or authentication failed

---

### ID Token Verification (Mobile SDK Flow)

These endpoints verify ID tokens obtained from native mobile SDKs (e.g., Google Sign-In SDK, Facebook SDK).

#### POST /auth/google/verify
Verify Google ID token from mobile SDK.

**Authentication:** Not required

**Request Body:**
```json
{
  "id_token": "eyJhbGciOiJSUzI1NiIs..."
}
```

**Response:** Same format as code exchange endpoints

**Error Responses:**
- `400 Bad Request`: Missing or invalid id_token
- `401 Unauthorized`: Invalid ID token or verification failed

---

#### POST /auth/facebook/verify
Verify Facebook access token from mobile SDK.

**Authentication:** Not required

**Request Body:**
```json
{
  "access_token": "EAAG..."
}
```

**Response:** Same format as code exchange endpoints

**Error Responses:**
- `400 Bad Request`: Missing or invalid access_token
- `401 Unauthorized`: Invalid access token or verification failed

---

### Refresh Token

#### POST /auth/refresh
Refresh the access token using a valid refresh token.

**Authentication:** Refresh token required (cookie or body)

**Request Body (optional if using cookie):**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:**
```json
{
  "status": "success",
  "type": "auth",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "blocked": false,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "is_admin": false
  },
  "message": "Tokens refreshed successfully"
}
```

**Cookies Updated:**
- `access_token`: New JWT access token
- `refresh_token`: New JWT refresh token

---

### Logout

#### POST /auth/logout
Logout the current user by revoking their refresh token.

**Authentication:** Required

**Response:**
```json
{
  "status": "success",
  "type": "action",
  "data": null,
  "message": "Logged out successfully"
}
```

**Cookies Cleared:**
- `access_token`
- `refresh_token`

---

## Trading Signals Endpoints

### List All Trading Signals

#### GET /api/trading-signals
Retrieve a list of all trading signals with pagination.

**Authentication:** Required

**Query Parameters:**
- `limit` (integer, optional, default: 50, max: 100): Number of results to return
- `offset` (integer, optional, default: 0): Number of results to skip

**Example Request:**
```
GET /api/trading-signals?limit=20&offset=0
```

**Response:**
```json
{
  "status": "success",
  "type": "collection",
  "data": {
    "signals": [
      {
        "id": 1,
        "symbol": "BTCUSDT",
        "stop_loss_price": 42000.50,
        "entry_price": 43000.00,
        "take_profit_price": 45000.00,
        "type": "LONG",
        "result": "WIN",
        "return": 4.65,
        "created_by": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "limit": 20,
    "offset": 0
  },
  "message": "Trading signals retrieved successfully"
}
```

---

### Get Single Trading Signal

#### GET /api/trading-signals/{id}
Retrieve a single trading signal by ID.

**Authentication:** Required

**Path Parameters:**
- `id` (integer, required): Trading signal ID

**Example Request:**
```
GET /api/trading-signals/1
```

**Response:**
```json
{
  "status": "success",
  "type": "resource",
  "data": {
    "id": 1,
    "symbol": "BTCUSDT",
    "stop_loss_price": 42000.50,
    "entry_price": 43000.00,
    "take_profit_price": 45000.00,
    "type": "LONG",
    "result": "WIN",
    "return": 4.65,
    "created_by": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "message": "Trading signal retrieved successfully"
}
```

---

### Create Trading Signal

#### POST /api/trading-signals
Create a new trading signal.

**Authentication:** Required (Admin only)

**Request Body:**
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

**Field Descriptions:**
- `symbol` (string, required): Trading pair symbol
- `stop_loss_price` (number, required): Stop loss price (must be > 0)
- `entry_price` (number, required): Entry price (must be > 0)
- `take_profit_price` (number, required): Take profit price (must be > 0)
- `type` (string, required): Signal type ("LONG" or "SHORT")
- `result` (string, optional): Signal result ("WIN", "LOSS", or "BREAKEVEN")
- `return` (number, optional): Return percentage

**Response:**
```json
{
  "status": "success",
  "type": "resource",
  "data": {
    "id": 1,
    "symbol": "BTCUSDT",
    "stop_loss_price": 42000.50,
    "entry_price": 43000.00,
    "take_profit_price": 45000.00,
    "type": "LONG",
    "result": "WIN",
    "return": 4.65,
    "created_by": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "message": "Trading signal created successfully"
}
```

**Status Code:** 201 Created

---

### Update Trading Signal

#### PUT /api/trading-signals/{id}
Update an existing trading signal.

**Authentication:** Required (Admin only)

**Path Parameters:**
- `id` (integer, required): Trading signal ID

**Request Body (all fields optional):**
```json
{
  "symbol": "ETHUSDT",
  "stop_loss_price": 2500.00,
  "entry_price": 2600.00,
  "take_profit_price": 2800.00,
  "type": "SHORT",
  "result": "LOSS",
  "return": -3.85
}
```

**Response:**
```json
{
  "status": "success",
  "type": "resource",
  "data": {
    "id": 1,
    "symbol": "ETHUSDT",
    "stop_loss_price": 2500.00,
    "entry_price": 2600.00,
    "take_profit_price": 2800.00,
    "type": "SHORT",
    "result": "LOSS",
    "return": -3.85,
    "created_by": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  },
  "message": "Trading signal updated successfully"
}
```

---

### Delete Trading Signal

#### DELETE /api/trading-signals/{id}
Delete a trading signal.

**Authentication:** Required (Admin only)

**Path Parameters:**
- `id` (integer, required): Trading signal ID

**Response:**
```json
{
  "status": "success",
  "type": "action",
  "data": null,
  "message": "Trading signal deleted successfully"
}
```

---

## Error Codes

| Status Code | Error Type | Description |
|-------------|------------|-------------|
| 400 | bad_request | Invalid request parameters |
| 400 | validation_error | Request validation failed |
| 401 | unauthorized | Authentication required or token invalid |
| 403 | forbidden | Insufficient permissions (e.g., admin required) |
| 404 | not_found | Resource not found |
| 409 | conflict | Resource conflict (e.g., duplicate entry) |
| 429 | rate_limit_exceeded | Too many requests |
| 500 | internal_server_error | Internal server error |
| 503 | service_unavailable | Service temporarily unavailable |

---

## Rate Limiting

The API implements rate limiting to prevent abuse:
- **Default Limit:** 100 requests per minute per IP address
- **Response:** HTTP 429 Too Many Requests

When rate limit is exceeded:
```json
{
  "status": "error",
  "type": "rate_limit_exceeded",
  "error": {
    "code": 429,
    "message": "Too many requests. Please try again later."
  }
}
```

---

## Account Linking

If the same email address is used across different OAuth providers (Google and Facebook), the accounts will be automatically linked to the same user account. This allows users to sign in with either provider.

**Example:**
1. User signs in with Google (email: user@example.com)
2. User later signs in with Facebook (same email: user@example.com)
3. The Facebook OAuth is linked to the existing user account
4. User can now authenticate with either Google or Facebook

---

## Testing with cURL

### Health Check
```bash
curl http://localhost:8080/health
```

### Get Trading Signals (with Bearer token)
```bash
curl -H "Authorization: Bearer <access_token>" \
     http://localhost:8080/api/trading-signals
```

### Create Trading Signal (Admin only)
```bash
curl -X POST \
     -H "Authorization: Bearer <access_token>" \
     -H "Content-Type: application/json" \
     -d '{
       "symbol": "BTCUSDT",
       "stop_loss_price": 42000.50,
       "entry_price": 43000.00,
       "take_profit_price": 45000.00,
       "type": "LONG"
     }' \
     http://localhost:8080/api/trading-signals
```

### Refresh Token
```bash
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"refresh_token": "<refresh_token>"}' \
     http://localhost:8080/auth/refresh
```

### Logout
```bash
curl -X POST \
     -H "Authorization: Bearer <access_token>" \
     http://localhost:8080/auth/logout
```

---

## Postman Collection

For easier testing, you can import this API into Postman:

1. Create a new collection in Postman
2. Add environment variables:
   - `base_url`: http://localhost:8080
   - `access_token`: (will be set automatically after login)
3. Add the endpoints listed above
4. Use `{{base_url}}` and `{{access_token}}` in your requests

---

---

## Package Endpoints

### GET /api/packages
Get all active packages available for subscription.

**Authentication:** Required

**Query Parameters:**
- `limit` (integer, optional): Number of results to return (default: 100, max: 100)
- `offset` (integer, optional): Number of results to skip (default: 0)

**Response:**
```json
{
  "status": "success",
  "type": "collection",
  "data": {
    "packages": [
      {
        "id": 1,
        "name": "Forex Short Term - Monthly",
        "asset_class": "FOREX",
        "duration_type": "SHORT_TERM",
        "billing_cycle": "MONTHLY",
        "duration_days": 30,
        "price": 10.00,
        "description": "Access to Forex day trading signals for 1 month",
        "is_active": true,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 18,
    "limit": 100,
    "offset": 0
  },
  "message": "Packages retrieved successfully"
}
```

### GET /api/packages/{id}
Get a specific package by ID.

**Authentication:** Required

**Response:** Same as single package object above.

---

## Subscription Endpoints

### POST /api/subscriptions
Subscribe to one or more packages.

**Authentication:** Required

**Request Body:**
```json
{
  "package_ids": [1, 5, 10]
}
```

**Response:**
```json
{
  "status": "success",
  "type": "resource",
  "data": {
    "subscriptions": [
      {
        "id": 1,
        "user_id": 123,
        "package_id": 1,
        "price_paid": 10.00,
        "subscribed_at": "2024-01-16T00:00:00Z",
        "expires_at": "2024-02-15T00:00:00Z",
        "is_active": true,
        "created_at": "2024-01-16T00:00:00Z",
        "updated_at": "2024-01-16T00:00:00Z",
        "package": {
          "id": 1,
          "name": "Forex Short Term - Monthly",
          "asset_class": "FOREX",
          "duration_type": "SHORT_TERM",
          "billing_cycle": "MONTHLY",
          "duration_days": 30,
          "price": 10.00
        }
      }
    ],
    "total_amount": 32.00,
    "message": "Successfully subscribed to 3 package(s)"
  },
  "message": "Successfully subscribed to 3 package(s)"
}
```

### GET /api/subscriptions/active
Get all active subscriptions for the authenticated user.

**Authentication:** Required

**Response:**
```json
{
  "status": "success",
  "type": "collection",
  "data": {
    "subscriptions": [
      {
        "id": 1,
        "user_id": 123,
        "package_id": 1,
        "price_paid": 10.00,
        "subscribed_at": "2024-01-16T00:00:00Z",
        "expires_at": "2024-02-15T00:00:00Z",
        "is_active": true,
        "package": { ... }
      }
    ],
    "total": 3
  },
  "message": "Active subscriptions retrieved successfully"
}
```

### GET /api/subscriptions/history
Get all subscriptions (active and expired) for the authenticated user.

**Authentication:** Required

**Query Parameters:**
- `limit` (integer, optional): Number of results to return (default: 50, max: 100)
- `offset` (integer, optional): Number of results to skip (default: 0)

**Response:** Same structure as active subscriptions with pagination.

### POST /api/subscriptions/check-access
Check if user has access to specific asset class and duration type.

**Authentication:** Required

**Request Body:**
```json
{
  "asset_class": "FOREX",
  "duration_type": "SHORT_TERM"
}
```

**Response:**
```json
{
  "status": "success",
  "type": "resource",
  "data": {
    "has_access": true,
    "expires_at": "2024-02-15T00:00:00Z",
    "package_id": 1,
    "price_paid": 10.00
  },
  "message": "Access check completed successfully"
}
```

---

## Payment Endpoints

### GET /api/payments/history
Get payment history for the authenticated user.

**Authentication:** Required

**Query Parameters:**
- `limit` (integer, optional): Number of results to return (default: 50, max: 100)
- `offset` (integer, optional): Number of results to skip (default: 0)

**Response:**
```json
{
  "status": "success",
  "type": "collection",
  "data": {
    "payments": [
      {
        "id": 1,
        "user_id": 123,
        "package_id": 1,
        "amount": 10.00,
        "payment_method": "dummy",
        "payment_status": "COMPLETED",
        "transaction_id": "dummy-123-1705449600",
        "metadata": "{\"note\": \"Dummy payment for development\"}",
        "created_at": "2024-01-16T00:00:00Z",
        "package": {
          "id": 1,
          "name": "Forex Short Term - Monthly",
          "asset_class": "FOREX",
          "duration_type": "SHORT_TERM"
        }
      }
    ],
    "total": 5,
    "limit": 50,
    "offset": 0
  },
  "message": "Payment history retrieved successfully"
}
```

---

## Updated Trading Signal Endpoints

### GET /api/trading-signals
Get trading signals visible to the authenticated user based on their active subscriptions and free-for-all signals.

**Authentication:** Required

**Query Parameters:**
- `limit` (integer, optional): Number of results to return (default: 50, max: 100)
- `offset` (integer, optional): Number of results to skip (default: 0)

**Response:**
```json
{
  "status": "success",
  "type": "collection",
  "data": {
    "signals": [
      {
        "id": 1,
        "symbol": "EURUSD",
        "asset_class": "FOREX",
        "duration_type": "SHORT_TERM",
        "stop_loss_price": 1.0850,
        "entry_price": 1.0900,
        "take_profit_price": 1.1000,
        "type": "LONG",
        "result": null,
        "return": null,
        "free_for_all": false,
        "comments": "Strong bullish momentum on 4H chart",
        "created_by": 1,
        "created_at": "2024-01-16T00:00:00Z",
        "updated_at": "2024-01-16T00:00:00Z"
      }
    ],
    "total": 25,
    "limit": 50,
    "offset": 0
  },
  "message": "Trading signals retrieved successfully"
}
```

**Note:** Users will only see:
1. Signals marked as `free_for_all: true` (visible to all users)
2. Signals matching their active subscriptions (asset_class and duration_type)

---

## Admin Endpoints

All admin endpoints require admin privileges. Add `/admin` prefix and use admin authentication.

### GET /api/admin/trading-signals
Get all trading signals (no filtering).

**Authentication:** Admin Required

**Response:** Same as regular signals endpoint but returns ALL signals.

### POST /api/admin/trading-signals
Create a new trading signal. Automatically sends notifications to configured channels.

**Authentication:** Admin Required

**Request Body:**
```json
{
  "symbol": "BTCUSDT",
  "asset_class": "CRYPTO",
  "duration_type": "LONG_TERM",
  "stop_loss_price": 42000.50,
  "entry_price": 43000.00,
  "take_profit_price": 45000.00,
  "type": "LONG",
  "free_for_all": false,
  "comments": "Bitcoin showing strong support at 42k"
}
```

### PUT /api/admin/trading-signals/{id}
Update a trading signal.

**Authentication:** Admin Required

**Request Body:** Partial update of signal fields.

### DELETE /api/admin/trading-signals/{id}
Delete a trading signal.

**Authentication:** Admin Required

### POST /api/admin/packages
Create a new package.

**Authentication:** Admin Required

**Request Body:**
```json
{
  "name": "Forex Short Term - Monthly",
  "asset_class": "FOREX",
  "duration_type": "SHORT_TERM",
  "billing_cycle": "MONTHLY",
  "duration_days": 30,
  "price": 10.00,
  "description": "Access to Forex day trading signals for 1 month"
}
```

### PUT /api/admin/packages/{id}
Update a package (including price changes).

**Authentication:** Admin Required

**Request Body:** Partial update of package fields.

**Note:** Price changes do NOT affect existing active subscriptions.

### DELETE /api/admin/packages/{id}
Delete a package.

**Authentication:** Admin Required

### POST /api/admin/payments
Manually record a payment (dummy implementation for development).

**Authentication:** Admin Required

**Request Body:**
```json
{
  "user_id": 123,
  "package_id": 1,
  "amount": 10.00,
  "payment_method": "manual",
  "payment_status": "COMPLETED",
  "transaction_id": "manual-txn-123"
}
```

---

## WebSocket Support

Currently, this API does not support WebSocket connections. All communication is done via HTTP REST endpoints. If real-time updates are needed, implement polling on the client side or consider adding WebSocket support in a future version.

