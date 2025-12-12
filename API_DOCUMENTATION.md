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

## WebSocket Support

Currently, this API does not support WebSocket connections. All communication is done via HTTP REST endpoints. If real-time updates are needed, implement polling on the client side or consider adding WebSocket support in a future version.

