# Subscription System Guide

## Overview

This trading signal platform uses a subscription-based model where users pay to access specific types of trading signals. The system supports three asset classes (Forex, Crypto, PSX) with two duration types (Short-Term, Long-Term) and three billing cycles (Monthly, 6 Months, Yearly).

## Package Structure

### Asset Classes
- **FOREX**: Foreign exchange currency pairs
- **CRYPTO**: Cryptocurrency trading signals
- **PSX**: Pakistan Stock Exchange signals

### Duration Types
- **SHORT_TERM**: Day trading signals (intraday positions)
- **LONG_TERM**: Swing trading signals (1-7 days holding period)

### Billing Cycles
- **MONTHLY**: 30 days access
- **SIX_MONTHS**: 180 days access
- **YEARLY**: 365 days access

### Example Packages

| Package | Asset Class | Duration Type | Billing Cycle | Price | Duration |
|---------|-------------|---------------|---------------|-------|----------|
| Forex Short Term - Monthly | FOREX | SHORT_TERM | MONTHLY | $10 | 30 days |
| Forex Short Term - 6 Months | FOREX | SHORT_TERM | SIX_MONTHS | $50 | 180 days |
| Forex Short Term - Yearly | FOREX | SHORT_TERM | YEARLY | $80 | 365 days |
| Forex Long Term - Monthly | FOREX | LONG_TERM | MONTHLY | $15 | 30 days |
| Crypto Short Term - Monthly | CRYPTO | SHORT_TERM | MONTHLY | $8 | 30 days |
| PSX Long Term - Yearly | PSX | LONG_TERM | YEARLY | $80 | 365 days |

## How Subscriptions Work

### 1. Package Selection
Users can subscribe to multiple packages simultaneously. For example:
- Forex Short Term - Monthly ($10)
- Crypto Long Term - Yearly ($95)
- PSX Long Term - Monthly ($10)
- **Total: $115**

### 2. Payment Processing
Currently, the system uses a dummy payment processor for development. In production, this will be replaced with:
- Stripe integration
- Binance Pay
- Manual payment verification

### 3. Subscription Activation
Upon successful payment:
- Subscription records are created with expiry dates
- User receives confirmation email
- Access is granted immediately
- Expiry date = Subscription date + duration_days

**Example:**
- Subscribe on: January 16, 2024
- Package: Monthly (30 days)
- Expires on: February 15, 2024

### 4. Price Protection
When a user subscribes to a package, the current price is locked in for that subscription period:
- Price stored in `price_paid` field
- If admin updates package price from $10 to $13
- Existing subscribers continue with $10 until expiry
- New subscribers pay $13
- When renewing, users pay the NEW price

## Signal Visibility Rules

Users can see trading signals based on:

### 1. Free-for-All Signals
Signals marked as `free_for_all: true` are visible to ALL authenticated users, regardless of subscriptions.

### 2. Subscription-Based Signals
Users see signals that match their active subscriptions:
- Signal must match BOTH `asset_class` AND `duration_type`
- Subscription must be active (`is_active: true`)
- Subscription must not be expired (`expires_at > current time`)

### Example Scenarios

**Scenario 1: User with Forex Short Term subscription**
- ✅ Can see: Forex Short Term signals
- ✅ Can see: All free-for-all signals
- ❌ Cannot see: Forex Long Term signals
- ❌ Cannot see: Crypto Short Term signals

**Scenario 2: User with multiple subscriptions**
- Subscriptions: Forex Short Term + Crypto Long Term
- ✅ Can see: Forex Short Term signals
- ✅ Can see: Crypto Long Term signals
- ✅ Can see: All free-for-all signals
- ❌ Cannot see: PSX signals (not subscribed)

**Scenario 3: User with expired subscription**
- Had: Forex Short Term (expired yesterday)
- ✅ Can see: Only free-for-all signals
- ❌ Cannot see: Any subscription-based signals

## Subscription Lifecycle

### Active Subscription
- Status: `is_active: true`
- Current date < `expires_at`
- User has full access to signals

### Expired Subscription
- Status: `is_active: false` (auto-updated)
- Current date >= `expires_at`
- User loses access to signals
- History preserved for reference

### Renewal Process
1. User's subscription expires
2. User browses packages (may see new prices)
3. User subscribes again
4. New subscription created with new expiry date
5. User pays current package price (not old price)

## Admin Features

### Package Management
Admins can:
- Create new packages with custom pricing
- Update existing packages (including prices)
- Deactivate packages (`is_active: false`)
- Delete packages (if no active subscriptions)

### Price Updates
- Admin updates package price via `PUT /api/admin/packages/{id}`
- Existing active subscriptions are NOT affected
- New subscriptions use updated price
- Transparent to users (they see price when subscribing)

### Signal Creation
When admin creates a signal:
1. Signal saved to database
2. Notifications sent automatically:
   - Telegram (if enabled)
   - Discord (if enabled)
   - Expo push (placeholder for future)
3. Signal becomes visible based on visibility rules

### Signal Types
Admins set:
- `asset_class`: FOREX, CRYPTO, or PSX
- `duration_type`: SHORT_TERM or LONG_TERM
- `type`: LONG (buy) or SHORT (sell)
- `free_for_all`: true/false
- Entry, Stop Loss, Take Profit prices
- Optional comments

## Email Notifications

### Subscription Confirmation
Users receive email after successful subscription:
- List of subscribed packages
- Expiry dates for each
- Total amount paid
- Confirmation number

### Email Providers
System supports multiple email providers (configured via env):
- **Resend**: Modern email API (recommended)
- **SMTP**: Traditional email (Gmail, etc.)
- **Mock**: Development mode (logs to console)

## Push Notifications

### Telegram
- Bot sends message to configured channel/group
- Format: "🚨 New [AssetClass] [Duration] Signal! Check the app."
- Includes asset and type information
- Minimal details (not full signal data)

### Discord
- Webhook sends embed to configured channel
- Formatted with colors and fields
- Includes timestamp and footer
- Professional appearance

### Expo (Future)
- Placeholder for mobile push notifications
- Will target users with active subscriptions
- Personalized based on user's subscriptions

## API Integration Examples

### 1. Get Available Packages
```bash
curl -H "Authorization: Bearer <token>" \
     http://localhost:8080/api/packages
```

### 2. Subscribe to Packages
```bash
curl -X POST \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"package_ids": [1, 5, 10]}' \
     http://localhost:8080/api/subscriptions
```

### 3. Check Active Subscriptions
```bash
curl -H "Authorization: Bearer <token>" \
     http://localhost:8080/api/subscriptions/active
```

### 4. Get Visible Signals
```bash
curl -H "Authorization: Bearer <token>" \
     http://localhost:8080/api/trading-signals
```

### 5. Check Access to Specific Type
```bash
curl -X POST \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"asset_class": "FOREX", "duration_type": "SHORT_TERM"}' \
     http://localhost:8080/api/subscriptions/check-access
```

## Database Schema

### Packages Table
- Stores all available subscription packages
- Admin can add/edit/delete
- Unique constraint on (asset_class, duration_type, billing_cycle)

### User Subscriptions Table
- Links users to packages they subscribed to
- Stores price_paid (for price protection)
- Tracks expiry dates
- Auto-deactivates on expiry

### Payment History Table
- Records all payment transactions
- Supports multiple payment methods
- JSONB metadata for flexibility
- Audit trail for all transactions

### Trading Signals Table
- Stores all trading signals
- Includes asset_class and duration_type
- Free-for-all flag for public signals
- Comments field for additional context

## Best Practices

### For Users
1. Subscribe to packages that match your trading style
2. Check expiry dates regularly
3. Renew before expiration to avoid access loss
4. Monitor both free-for-all and subscribed signals

### For Admins
1. Set reasonable prices for packages
2. Mark promotional signals as free-for-all
3. Update prices carefully (affects new subscriptions only)
4. Include detailed comments in signals
5. Test notifications before sending to production

### For Developers
1. Always check subscription expiry before showing signals
2. Run subscription cleanup job periodically
3. Monitor email delivery success rates
4. Test notification integrations regularly
5. Handle payment failures gracefully

## Troubleshooting

### User can't see signals
1. Check if subscription is active
2. Verify subscription hasn't expired
3. Confirm signal matches user's subscriptions
4. Check if signal is marked as free-for-all

### Price update not reflected
1. Check if user has existing active subscription
2. Existing subscriptions use old price (by design)
3. New subscriptions will use new price

### Email not received
1. Check EMAIL_SERVICE_ENABLED in config
2. Verify email provider credentials
3. Check spam folder
4. Review application logs

### Notification not sent
1. Verify notification services are enabled
2. Check API keys/tokens/webhook URLs
3. Review application logs
4. Test notification services manually

## Future Enhancements

1. **Auto-renewal**: Automatic subscription renewal before expiry
2. **Payment Gateway**: Integration with Stripe, Binance Pay
3. **Discounts**: Coupon codes and promotional pricing
4. **Free Trials**: Limited time free access
5. **Bundle Deals**: Discounted multi-package bundles
6. **Referral System**: Reward users for referrals
7. **Analytics**: Track subscription metrics
8. **Mobile App**: Native iOS/Android apps with Expo
9. **WebSocket**: Real-time signal updates
10. **Performance Tracking**: Automated result tracking

