# Test Subscription Access Control

## Critical Security Fix Applied ✅

The `GET /api/trading-signals/{id}` endpoint now properly enforces subscription-based access control.

## How It Works

### Access Rules
A user can view a signal **ONLY IF**:
1. The signal has `free_for_all: true`, **OR**
2. The user has an **active subscription** that matches both:
   - Signal's `asset_class` (FOREX, CRYPTO, or PSX)
   - Signal's `duration_type` (SHORT_TERM or LONG_TERM)

### Active Subscription Criteria
- `is_active = true`
- `expires_at > CURRENT_TIMESTAMP`

## Testing Steps

### Step 1: Create a Test Signal (Admin)

```bash
# Login as admin first to get admin token

# Create a FOREX SHORT_TERM signal (not free-for-all)
curl -X POST http://localhost:8080/api/admin/trading-signals \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "EURUSD",
    "asset_class": "FOREX",
    "duration_type": "SHORT_TERM",
    "stop_loss_price": 1.0850,
    "entry_price": 1.0900,
    "take_profit_price": 1.1000,
    "type": "LONG",
    "free_for_all": false,
    "comments": "Test signal for subscription check"
  }'

# Note the signal ID from response (e.g., id: 1)
```

### Step 2: Try to Access WITHOUT Subscription (Regular User)

```bash
# Login as regular user (not admin)
# Get user access token

# Try to access the signal by ID
curl http://localhost:8080/api/trading-signals/1 \
  -H "Authorization: Bearer <user_token>"

# Expected Result: 403 Forbidden
{
  "status": "error",
  "type": "forbidden",
  "error": {
    "code": 403,
    "message": "You don't have access to this signal. Subscribe to the appropriate package to view this signal."
  }
}
```

### Step 3: Subscribe to Matching Package

```bash
# Get available packages first
curl http://localhost:8080/api/packages \
  -H "Authorization: Bearer <user_token>"

# Find FOREX + SHORT_TERM package (e.g., id: 1 - "Forex Short Term - Monthly")

# Subscribe to the package
curl -X POST http://localhost:8080/api/subscriptions \
  -H "Authorization: Bearer <user_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "package_ids": [1]
  }'

# Verify subscription is active
curl http://localhost:8080/api/subscriptions/active \
  -H "Authorization: Bearer <user_token>"
```

### Step 4: Try to Access WITH Subscription

```bash
# Now try to access the same signal
curl http://localhost:8080/api/trading-signals/1 \
  -H "Authorization: Bearer <user_token>"

# Expected Result: 200 OK with signal data
{
  "status": "success",
  "type": "resource",
  "data": {
    "id": 1,
    "symbol": "EURUSD",
    "asset_class": "FOREX",
    "duration_type": "SHORT_TERM",
    "stop_loss_price": 1.0850,
    "entry_price": 1.0900,
    "take_profit_price": 1.1000,
    "type": "LONG",
    "free_for_all": false,
    "comments": "Test signal for subscription check",
    ...
  },
  "message": "Trading signal retrieved successfully"
}
```

### Step 5: Test Free-For-All Signal

```bash
# Create a free-for-all signal (Admin)
curl -X POST http://localhost:8080/api/admin/trading-signals \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "GBPUSD",
    "asset_class": "FOREX",
    "duration_type": "LONG_TERM",
    "stop_loss_price": 1.2650,
    "entry_price": 1.2700,
    "take_profit_price": 1.2850,
    "type": "LONG",
    "free_for_all": true,
    "comments": "Free signal for everyone"
  }'

# Try to access as user WITHOUT subscription (different user or unsubscribe first)
curl http://localhost:8080/api/trading-signals/2 \
  -H "Authorization: Bearer <user_token>"

# Expected Result: 200 OK (free-for-all signals are accessible to all)
```

### Step 6: Test Wrong Subscription Type

```bash
# User has FOREX + SHORT_TERM subscription
# Try to access CRYPTO + SHORT_TERM signal

# Create CRYPTO signal (Admin)
curl -X POST http://localhost:8080/api/admin/trading-signals \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "asset_class": "CRYPTO",
    "duration_type": "SHORT_TERM",
    "stop_loss_price": 42000,
    "entry_price": 43000,
    "take_profit_price": 45000,
    "type": "LONG",
    "free_for_all": false,
    "comments": "Crypto signal"
  }'

# Try to access with FOREX subscription
curl http://localhost:8080/api/trading-signals/3 \
  -H "Authorization: Bearer <user_token>"

# Expected Result: 403 Forbidden (wrong asset class)
```

## What Changed in the Code

### Handler (`internal/handlers/trading_signal_handler.go`)
```go
func (h *TradingSignalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    // 1. Get authenticated user ID
    userID, ok := middleware.GetUserIDFromContext(r.Context())
    
    // 2. Parse signal ID from URL
    id, err := strconv.ParseInt(idStr, 10, 64)
    
    // 3. CHECK ACCESS (NEW!)
    hasAccess, err := h.service.CheckUserAccessToSignal(userID, id)
    
    // 4. Return 403 if no access
    if !hasAccess {
        utils.SendError(w, http.StatusForbidden, ...)
        return
    }
    
    // 5. Only retrieve signal if access verified
    signal, err := h.service.GetByID(id)
}
```

### Service (`internal/services/trading_signal_service.go`)
```go
func (s *TradingSignalService) CheckUserAccessToSignal(userID, signalID int64) (bool, error) {
    return s.repo.CheckUserAccess(userID, signalID)
}
```

### Repository (`internal/repositories/trading_signal_repository.go`)
```go
func (r *TradingSignalRepository) CheckUserAccess(userID, signalID int64) (bool, error) {
    query := `
        SELECT EXISTS (
            SELECT 1 FROM trading_signals ts
            WHERE ts.id = $2
            AND (
                ts.free_for_all = true
                OR EXISTS (
                    SELECT 1 FROM user_subscriptions us
                    JOIN packages p ON us.package_id = p.id
                    WHERE us.user_id = $1
                    AND us.is_active = true
                    AND us.expires_at > CURRENT_TIMESTAMP
                    AND p.asset_class = ts.asset_class
                    AND p.duration_type = ts.duration_type
                )
            )
        )
    `
    // Returns true if user has access, false otherwise
}
```

## Security Benefits

✅ **Subscription model enforced** - No bypass via direct ID access
✅ **Consistent access control** - Both list and detail endpoints use same logic
✅ **Efficient query** - Single database query using EXISTS
✅ **Clear error messages** - Users know they need to subscribe
✅ **Free-for-all support** - Promotional signals still accessible
✅ **Time-based expiry** - Expired subscriptions automatically denied

## Common Issues

### Issue: Still seeing signals without subscription
**Solution:** 
1. Restart the server after applying the fix
2. Clear any cached tokens
3. Verify the signal's `free_for_all` is false
4. Check user's active subscriptions

### Issue: Can't access signal even with subscription
**Solution:**
1. Verify subscription is active: `GET /api/subscriptions/active`
2. Check subscription hasn't expired
3. Ensure asset_class and duration_type match exactly
4. Check package ID matches signal type

### Issue: Admin can't access signals
**Solution:**
Admins should use the admin endpoint:
```bash
GET /api/admin/trading-signals/{id}
```
Regular user endpoint enforces subscription even for admins.

## Database Query Explanation

The access check query does:
1. Checks if signal exists with given ID
2. Returns true if EITHER:
   - Signal is `free_for_all = true`
   - OR user has active subscription matching signal type
3. Uses EXISTS for optimal performance
4. Single query, no N+1 issues

## Postman Testing

Update your Postman collection with these scenarios:
1. **Folder: Access Control Tests**
   - No Subscription → Get Signal (expect 403)
   - Subscribe → Get Signal (expect 200)
   - Wrong Subscription → Get Signal (expect 403)
   - Free Signal → Get Signal (expect 200)

---

**Status:** ✅ CRITICAL SECURITY FIX APPLIED

The subscription bypass vulnerability has been fixed. Users can only access signals they have active subscriptions for (or free-for-all signals).

