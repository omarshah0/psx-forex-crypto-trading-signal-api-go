# Postman Collection Update

## Summary

The Postman collection has been updated with all new subscription system endpoints.

## New Folders Added

### 1. Packages (5 endpoints)
- ✅ **Get All Packages** - Browse available subscription packages
- ✅ **Get Package by ID** - View specific package details
- ✅ **Create Package (Admin)** - Add new packages
- ✅ **Update Package (Admin)** - Modify package details and prices
- ✅ **Delete Package (Admin)** - Remove packages

### 2. Subscriptions (4 endpoints)
- ✅ **Subscribe to Packages** - Subscribe to one or multiple packages
- ✅ **Get Active Subscriptions** - View current active subscriptions
- ✅ **Get Subscription History** - View all subscriptions (active & expired)
- ✅ **Check Access** - Verify access to specific signal type

### 3. Payments (2 endpoints)
- ✅ **Get Payment History** - View payment transaction history
- ✅ **Record Payment (Admin)** - Manually record payments

### 4. Trading Signals (Updated) (5 endpoints)
- ✅ **Get User Signals (Filtered)** - View signals based on subscriptions
- ✅ **Get All Signals (Admin)** - View all signals without filtering
- ✅ **Create Signal (Admin)** - Create new signal with notifications
- ✅ **Update Signal (Admin)** - Update signal details and results
- ✅ **Delete Signal (Admin)** - Remove signals

## How to Use

### 1. Import Collection
```bash
# In Postman:
1. Click "Import" button
2. Select "postman_collection.json"
3. Collection will be imported with all folders
```

### 2. Environment Variables
The collection uses these variables:
- `{{base_url}}` - Default: http://localhost:8080
- `{{access_token}}` - Automatically set after login
- `{{refresh_token}}` - Automatically set after login

### 3. Testing Flow

#### A. User Flow (Complete Subscription Journey)
```
1. Register/Login → Get access_token
2. GET /api/packages → Browse available packages
3. POST /api/subscriptions → Subscribe to packages
   Body: {"package_ids": [1, 5, 10]}
4. GET /api/subscriptions/active → Verify active subscriptions
5. GET /api/trading-signals → View signals (filtered by subscription)
6. GET /api/payments/history → Check payment records
```

#### B. Admin Flow (Signal Management)
```
1. Login as Admin → Get admin access_token
2. GET /api/admin/trading-signals → View all signals
3. POST /api/admin/trading-signals → Create new signal
   Body: {
     "symbol": "EURUSD",
     "asset_class": "FOREX",
     "duration_type": "SHORT_TERM",
     "entry_price": 1.0900,
     "stop_loss_price": 1.0850,
     "take_profit_price": 1.1000,
     "type": "LONG",
     "free_for_all": false,
     "comments": "Strong bullish momentum"
   }
4. PUT /api/admin/trading-signals/1 → Update signal result
5. POST /api/admin/packages → Create new package
6. PUT /api/admin/packages/1 → Update package price
```

## Request Examples

### Subscribe to Multiple Packages
```json
POST /api/subscriptions
{
  "package_ids": [1, 5, 10]
}

Response:
{
  "status": "success",
  "data": {
    "subscriptions": [...],
    "total_amount": 32.00,
    "message": "Successfully subscribed to 3 package(s)"
  }
}
```

### Create Trading Signal
```json
POST /api/admin/trading-signals
{
  "symbol": "BTCUSDT",
  "asset_class": "CRYPTO",
  "duration_type": "LONG_TERM",
  "stop_loss_price": 42000.50,
  "entry_price": 43000.00,
  "take_profit_price": 45000.00,
  "type": "LONG",
  "free_for_all": false,
  "comments": "Bitcoin showing strong support"
}
```

### Check Access
```json
POST /api/subscriptions/check-access
{
  "asset_class": "FOREX",
  "duration_type": "SHORT_TERM"
}

Response:
{
  "status": "success",
  "data": {
    "has_access": true,
    "expires_at": "2024-02-15T00:00:00Z",
    "package_id": 1,
    "price_paid": 10.00
  }
}
```

### Update Package Price
```json
PUT /api/admin/packages/1
{
  "price": 13.00,
  "description": "Updated pricing"
}
```

## Testing Scenarios

### Scenario 1: New User Subscription
1. Login as regular user
2. Get all packages
3. Subscribe to Forex Short Term + Crypto Long Term
4. Check active subscriptions
5. Get trading signals (should see only subscribed types)

### Scenario 2: Price Update Protection
1. Login as admin
2. Create subscription for user at $10
3. Update package price to $13
4. Verify user's active subscription still shows $10
5. New user subscribing should see $13

### Scenario 3: Signal Visibility
1. Admin creates signal with `free_for_all: true`
2. All users can see it (even without subscription)
3. Admin creates signal with `free_for_all: false`
4. Only users with matching subscription can see it

### Scenario 4: Subscription Expiry
1. User subscribes to monthly package (30 days)
2. Check expires_at date
3. After expiry, GET /api/trading-signals returns only free signals
4. User must resubscribe to regain access

## Authentication

All endpoints (except health check) require authentication:
```
Authorization: Bearer {{access_token}}
```

The collection automatically:
1. Saves tokens after login
2. Uses saved tokens in subsequent requests
3. Refreshes tokens when expired

## Admin Endpoints

Admin-only endpoints require:
1. Valid access_token
2. User must have admin role in database
3. Admin verification middleware checks role

To test admin endpoints:
1. Create admin user in database
2. Login with admin credentials
3. Use returned access_token

## Notes

1. **Package IDs**: The seeded packages have IDs 1-18
   - Forex: IDs 1-6
   - Crypto: IDs 7-12
   - PSX: IDs 13-18

2. **Signal Types**:
   - `type`: LONG (buy) or SHORT (sell)
   - `asset_class`: FOREX, CRYPTO, PSX
   - `duration_type`: SHORT_TERM, LONG_TERM

3. **Result Values**: WIN, LOSS, BREAKEVEN (optional)

4. **Payment Status**: PENDING, COMPLETED, FAILED, REFUNDED

5. **Billing Cycles**: MONTHLY (30 days), SIX_MONTHS (180 days), YEARLY (365 days)

## Troubleshooting

### 401 Unauthorized
- Ensure access_token is set in environment
- Token may be expired, refresh or re-login
- Check if endpoint requires admin role

### 400 Bad Request
- Verify request body matches schema
- Check required fields are present
- Validate enum values (asset_class, duration_type, etc.)

### 404 Not Found
- Check if resource ID exists
- Package IDs should be 1-18 (from seed data)
- Signal IDs start from 1

### 500 Internal Server Error
- Check server logs
- Verify database migrations ran successfully
- Ensure all services (Postgres, MongoDB, Redis) are running

## Environment Setup

Before testing:
1. Run database migrations
2. Seed packages (migration 000011)
3. Create at least one admin user
4. Start the server
5. Update environment variables if using different ports

## Collection Structure

```
Go REST API with Social Auth
├── Health (1 endpoint)
├── OAuth - Code Exchange (2 endpoints)
├── OAuth - Token Verification (2 endpoints)
├── Email/Password Auth (7 endpoints)
├── Token Management (2 endpoints)
├── Profile (1 endpoint)
├── Packages (5 endpoints) ← NEW
├── Subscriptions (4 endpoints) ← NEW
├── Payments (2 endpoints) ← NEW
├── Trading Signals (Updated) (5 endpoints) ← UPDATED
└── Trading Signals - Legacy (4 endpoints) ← OLD
```

Total: 35+ endpoints covering complete API functionality

