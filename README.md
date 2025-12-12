# Trading Signal Subscription API

A production-ready Go REST API for trading signal subscriptions with multi-tier pricing, social authentication, email notifications, and real-time alerts via Telegram/Discord.

## 🎯 Overview

This API provides a complete subscription-based trading signal platform supporting:
- **Forex, Crypto, and PSX (Pakistan Stock Exchange)** signals
- **Short-term (day trading)** and **Long-term (swing trading)** signals
- **Multi-tier pricing** with monthly, 6-month, and yearly billing cycles
- **Subscription-based access control** with automatic expiry management
- **Social authentication** (Google & Facebook OAuth2)
- **Email/Password authentication** with verification
- **Push notifications** via Telegram, Discord, and Expo (mobile)

## ✨ Key Features

### 📊 Trading Signals
- Multi-asset class support (Forex, Crypto, PSX)
- Short-term and long-term signal types
- Entry, Stop-Loss, Take-Profit prices
- Signal results tracking (Win/Loss/Breakeven)
- Free-for-all promotional signals
- Admin-only signal creation with auto-notifications

### 💳 Subscription System
- 18 package combinations (3 assets × 2 durations × 3 billing cycles)
- Multi-package subscriptions (users can subscribe to multiple packages)
- Price protection (existing subscriptions unaffected by price changes)
- Flexible expiry (30, 180, or 365 days based on billing cycle)
- Automatic subscription expiry management
- Payment history tracking

### 🔐 Authentication & Security
- Social OAuth (Google & Facebook) with profile picture sync
- Email/Password authentication with email verification
- JWT tokens (access & refresh) with Redis storage
- Role-based access control (User & Admin)
- Subscription-based signal visibility
- Rate limiting and CORS protection

### 📧 Notifications & Emails
- **Email providers**: Resend API or SMTP (Gmail, SendGrid, etc.)
- **Telegram notifications**: Bot sends signal alerts to channel/group
- **Discord notifications**: Webhook integration with formatted embeds
- **Expo push notifications**: Placeholder for mobile apps
- Subscription confirmation emails
- Password reset and verification emails

### 🗄️ Database Architecture
- **PostgreSQL**: Users, subscriptions, packages, payments, signals
- **MongoDB**: Request/response logging with sensitive data masking
- **Redis**: JWT token storage, rate limiting, caching

## 🏗️ System Architecture

```
┌─────────────┐
│   Client    │ (Web/Mobile)
└──────┬──────┘
       │
       ↓
┌─────────────────────────────────────┐
│         API Gateway                 │
│  (Authentication + Rate Limiting)   │
└──────────────┬──────────────────────┘
               │
       ┌───────┴───────┐
       │               │
       ↓               ↓
┌─────────────┐  ┌─────────────┐
│  User API   │  │  Admin API  │
│             │  │             │
│ • Get Pkgs  │  │ • Create    │
│ • Subscribe │  │ • Update    │
│ • View Sigs │  │ • Delete    │
└──────┬──────┘  └──────┬──────┘
       │                │
       └────────┬───────┘
                │
       ┌────────┴────────┐
       │                 │
       ↓                 ↓
┌──────────────┐  ┌──────────────┐
│  PostgreSQL  │  │   MongoDB    │
│  (Data)      │  │   (Logs)     │
└──────────────┘  └──────────────┘
       │
       ↓
┌──────────────┐
│    Redis     │
│  (Cache)     │
└──────────────┘
```

## 🚀 Quick Start

### Prerequisites

- **Go** 1.21 or higher
- **Docker & Docker Compose** (for local databases)
- **Google OAuth credentials** (optional, for social login)
- **Facebook OAuth credentials** (optional, for social login)
- **Email provider** (Resend API key or SMTP credentials)
- **Telegram Bot** (optional, for notifications)
- **Discord Webhook** (optional, for notifications)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd forex-crypto-psx-stocks-signal-app/go
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Start databases with Docker Compose**
   ```bash
   docker-compose up -d
   ```
   
   This starts:
   - PostgreSQL on `localhost:5432`
   - MongoDB on `localhost:27017`
   - Redis on `localhost:6379`

4. **Set up environment variables**
   ```bash
   cp env.example .env
   ```
   
   Edit `.env` and configure:
   - JWT secrets (required)
   - Database URLs (default values work with Docker Compose)
   - OAuth credentials (if using social login)
   - Email provider settings
   - Notification service credentials

5. **Run database migrations**
   ```bash
   # Install golang-migrate if not already installed
   # macOS: brew install golang-migrate
   # Linux: See https://github.com/golang-migrate/migrate
   
   migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/rest_api_db?sslmode=disable" up
   ```
   
   Or use the Makefile:
   ```bash
   make migrate-up
   ```

6. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```
   
   Or use the Makefile:
   ```bash
   make run
   ```

The server will start on `http://localhost:8080`

## ⚙️ Configuration Guide

### 1. Database Setup

#### Using Docker Compose (Development)

```bash
# Start all databases
docker-compose up -d

# Check if running
docker-compose ps

# View logs
docker-compose logs -f postgres
docker-compose logs -f mongodb
docker-compose logs -f redis
```

#### Using Managed Services (Production)

Update `.env` with your managed database URLs:

```bash
# PostgreSQL (example: AWS RDS, DigitalOcean, Render)
POSTGRES_URL=postgres://username:password@host:5432/database?sslmode=require

# MongoDB (example: MongoDB Atlas)
MONGODB_URL=mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority

# Redis (example: Redis Cloud, AWS ElastiCache)
REDIS_URL=redis://:password@host:port
# OR for URL format:
REDIS_URL=redis://user:password@host:port/db
```

### 2. JWT Configuration (Required)

Generate secure random secrets:

```bash
# Generate secrets using the helper script
go run scripts/generate_secrets.go

# Or manually generate with openssl
openssl rand -base64 32
```

Add to `.env`:
```bash
JWT_ACCESS_SECRET=your-generated-secret-here
JWT_REFRESH_SECRET=your-different-generated-secret-here
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h  # 7 days
```

### 3. Email Provider Setup

Choose **ONE** email provider:

#### Option A: Resend (Recommended)

1. Sign up at [resend.com](https://resend.com)
2. Get your API key
3. Verify your domain

```bash
EMAIL_SERVICE_ENABLED=true
EMAIL_PROVIDER=resend
RESEND_API_KEY=re_xxxxxxxxxxxxx
EMAIL_FROM_ADDRESS=noreply@yourdomain.com
EMAIL_FROM_NAME=Trading Signals
```

#### Option B: SMTP (Gmail, SendGrid, etc.)

For Gmail:
1. Enable 2FA on your Google account
2. Generate an [App Password](https://support.google.com/accounts/answer/185833)

```bash
EMAIL_SERVICE_ENABLED=true
EMAIL_PROVIDER=smtp
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
EMAIL_FROM_ADDRESS=your-email@gmail.com
EMAIL_FROM_NAME=Trading Signals
```

For other SMTP providers:
- **SendGrid**: smtp.sendgrid.net:587
- **Mailgun**: smtp.mailgun.org:587
- **AWS SES**: email-smtp.region.amazonaws.com:587

#### Option C: Mock (Development Only)

```bash
EMAIL_SERVICE_ENABLED=false
```
Emails will be logged to console instead.

### 4. OAuth Setup (Optional)

#### Google OAuth

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project
3. Enable Google+ API
4. Create OAuth 2.0 credentials:
   - Application type: Web application
   - Authorized redirect URIs: `http://localhost:8080/auth/google/callback`
5. Copy credentials to `.env`:

```bash
OAUTH_GOOGLE_ENABLED=true
OAUTH_GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
OAUTH_GOOGLE_CLIENT_SECRET=your-client-secret
OAUTH_GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
```

#### Facebook OAuth

1. Go to [Facebook Developers](https://developers.facebook.com/)
2. Create a new app (Consumer app type)
3. Add "Facebook Login" product
4. Configure settings:
   - Valid OAuth Redirect URIs: `http://localhost:8080/auth/facebook/callback`
5. Copy credentials to `.env`:

```bash
OAUTH_FACEBOOK_ENABLED=true
OAUTH_FACEBOOK_CLIENT_ID=your-app-id
OAUTH_FACEBOOK_CLIENT_SECRET=your-app-secret
OAUTH_FACEBOOK_REDIRECT_URL=http://localhost:8080/auth/facebook/callback
```

### 5. Notification Services (Optional)

#### Telegram Notifications

1. Create a bot with [@BotFather](https://t.me/botfather)
2. Get your bot token
3. Add bot to your channel/group
4. Get chat ID:
   ```bash
   # Send a message to your bot, then:
   curl https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates
   # Look for "chat":{"id": -1001234567890}
   ```

```bash
TELEGRAM_NOTIFICATIONS_ENABLED=true
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
TELEGRAM_CHAT_ID=-1001234567890
```

#### Discord Notifications

1. Open your Discord server settings
2. Go to Integrations → Webhooks
3. Create a new webhook
4. Copy the webhook URL

```bash
DISCORD_NOTIFICATIONS_ENABLED=true
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/123456789/abcdefghijklmnop
```

#### Expo Push Notifications (Coming Soon)

```bash
EXPO_NOTIFICATIONS_ENABLED=false
```
Implementation placeholder for mobile app integration.

### 6. Subscription Settings

```bash
SUBSCRIPTION_DEFAULT_EXPIRY_DAYS=30
```

### 7. Email/Password Authentication

```bash
EMAIL_PASSWORD_AUTH_ENABLED=true
REQUIRE_EMAIL_VERIFICATION=true
VERIFICATION_TOKEN_EXPIRY=24h
RESET_TOKEN_EXPIRY=1h
FRONTEND_URL=http://localhost:3000
```

## 📦 Package System

### Available Packages (Seeded by Default)

| Asset Class | Duration | Monthly | 6 Months | Yearly |
|-------------|----------|---------|----------|--------|
| **Forex** | Short Term | $10 (30d) | $50 (180d) | $80 (365d) |
| **Forex** | Long Term | $15 (30d) | $75 (180d) | $120 (365d) |
| **Crypto** | Short Term | $8 (30d) | $40 (180d) | $65 (365d) |
| **Crypto** | Long Term | $12 (30d) | $60 (180d) | $95 (365d) |
| **PSX** | Short Term | $5 (30d) | $25 (180d) | $40 (365d) |
| **PSX** | Long Term | $10 (30d) | $50 (180d) | $80 (365d) |

**Total: 18 packages** automatically seeded during migrations.

### Managing Packages

Admins can create, update, or delete packages via API:

```bash
# Create new package
POST /api/admin/packages

# Update package price (doesn't affect existing subscriptions)
PUT /api/admin/packages/{id}

# Delete package
DELETE /api/admin/packages/{id}
```

## 🎫 Subscription Flow

1. **Browse Packages** → `GET /api/packages`
2. **Subscribe** → `POST /api/subscriptions` with package IDs
3. **Payment Processed** (currently dummy, integrate Stripe/Binance Pay)
4. **Subscription Activated** with expiry date
5. **Confirmation Email Sent**
6. **Access Granted** to matching signals

### Subscription Example

```bash
# User subscribes to Forex Short Term + Crypto Long Term
POST /api/subscriptions
{
  "package_ids": [1, 10]  # IDs from GET /api/packages
}

# User can now see:
# ✅ Forex + Short Term signals
# ✅ Crypto + Long Term signals
# ✅ All free-for-all signals
# ❌ PSX signals (not subscribed)
# ❌ Forex Long Term signals (not subscribed)
```

## 📡 API Endpoints

### Public Endpoints

- `GET /health` - Health check

### Authentication Endpoints

- `POST /auth/register` - Email/password registration
- `POST /auth/login` - Email/password login
- `GET /auth/verify-email` - Email verification
- `POST /auth/forgot-password` - Request password reset
- `POST /auth/reset-password` - Reset password with token
- `POST /auth/google/exchange` - Google OAuth code exchange
- `POST /auth/facebook/exchange` - Facebook OAuth code exchange
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - Logout (requires auth)

### User Endpoints (Authenticated)

**Packages:**
- `GET /api/packages` - Browse available packages
- `GET /api/packages/{id}` - Get package details

**Subscriptions:**
- `POST /api/subscriptions` - Subscribe to packages
- `GET /api/subscriptions/active` - View active subscriptions
- `GET /api/subscriptions/history` - View subscription history
- `POST /api/subscriptions/check-access` - Check access to signal type

**Payments:**
- `GET /api/payments/history` - View payment history

**Trading Signals:**
- `GET /api/trading-signals` - List visible signals (filtered by subscription)
- `GET /api/trading-signals/{id}` - Get signal details (requires access)

**Profile:**
- `GET /api/profile` - Get user profile

### Admin Endpoints (Admin Role Required)

**Trading Signals:**
- `GET /api/admin/trading-signals` - List all signals (no filtering)
- `POST /api/admin/trading-signals` - Create signal (triggers notifications)
- `PUT /api/admin/trading-signals/{id}` - Update signal
- `DELETE /api/admin/trading-signals/{id}` - Delete signal

**Packages:**
- `POST /api/admin/packages` - Create package
- `PUT /api/admin/packages/{id}` - Update package (price changes don't affect existing subscriptions)
- `DELETE /api/admin/packages/{id}` - Delete package

**Payments:**
- `POST /api/admin/payments` - Manually record payment

See [API_DOCUMENTATION.md](API_DOCUMENTATION.md) for complete API reference.

## 🔒 Security Features

### Subscription-Based Access Control
- Users can only view signals matching their active subscriptions
- Direct signal ID access is protected (no subscription bypass)
- Free-for-all signals visible to all authenticated users
- Expired subscriptions automatically lose access

### Authentication & Authorization
- JWT tokens with Redis blacklisting
- Role-based access control (User/Admin)
- Email verification for password auth
- OAuth profile sync with picture import

### Data Protection
- Sensitive data masking in logs
- HTTP-only cookies for web clients
- CORS configuration
- Rate limiting (100 req/min per IP)
- Blocked user management

## 👨‍💼 Admin Management

### Creating an Admin User

```sql
-- Connect to PostgreSQL
psql -U postgres -d rest_api_db

-- First, create or identify a user
INSERT INTO users (email, name, email_verified, blocked) 
VALUES ('admin@example.com', 'Admin User', true, false)
RETURNING id;

-- Then, make them an admin (use the ID from above)
INSERT INTO admins (user_id) VALUES (1);
```

Or create via registration then promote:
```sql
-- After user registers via API
INSERT INTO admins (user_id) 
SELECT id FROM users WHERE email = 'admin@example.com';
```

### Blocking Users

```sql
UPDATE users SET blocked = true WHERE email = 'user@example.com';
```

## 🧪 Testing

### Using Postman

1. **Import Collection**: `postman_collection.json`
2. **Import Environment**: `postman_environment.json`
3. **Login**: Use any auth endpoint to get tokens
4. **Test Endpoints**: All endpoints pre-configured

See [POSTMAN_GUIDE.md](POSTMAN_GUIDE.md) and [TEST_SUBSCRIPTION_ACCESS.md](TEST_SUBSCRIPTION_ACCESS.md) for testing guides.

### Manual Testing Flow

```bash
# 1. Register user
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","name":"Test User","password":"password123"}'

# 2. Browse packages
curl http://localhost:8080/api/packages \
  -H "Authorization: Bearer <token>"

# 3. Subscribe
curl -X POST http://localhost:8080/api/subscriptions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"package_ids":[1,5]}'

# 4. View signals
curl http://localhost:8080/api/trading-signals \
  -H "Authorization: Bearer <token>"
```

## 📚 Documentation

- **[API_DOCUMENTATION.md](API_DOCUMENTATION.md)** - Complete API reference
- **[SUBSCRIPTION_GUIDE.md](SUBSCRIPTION_GUIDE.md)** - Subscription system explained
- **[QUICK_START.md](QUICK_START.md)** - Quick start guide
- **[EMAIL_PASSWORD_AUTH.md](EMAIL_PASSWORD_AUTH.md)** - Email auth setup
- **[POSTMAN_GUIDE.md](POSTMAN_GUIDE.md)** - Postman collection guide
- **[BUG_FIXES.md](BUG_FIXES.md)** - Recent bug fixes
- **[IMPLEMENTATION_COMPLETE.md](IMPLEMENTATION_COMPLETE.md)** - Implementation details

## 🐳 Docker Compose Services

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Stop and remove all data (⚠️ Warning: deletes databases)
docker-compose down -v

# Restart specific service
docker-compose restart postgres
```

## 🏭 Production Deployment

### Environment Configuration

1. **Set all required environment variables**
2. **Use managed database services** (not Docker)
3. **Enable HTTPS** with valid SSL certificates
4. **Set security flags:**
   ```bash
   ENVIRONMENT=production
   COOKIE_SECURE=true
   ```
5. **Configure CORS** for your frontend domains
6. **Use strong passwords** for all services
7. **Enable monitoring** and alerting
8. **Set up backups** for PostgreSQL
9. **Use reverse proxy** (nginx/caddy) for SSL termination
10. **Configure log rotation** for MongoDB logs

### Deployment Checklist

- [ ] Database migrations applied
- [ ] Environment variables set securely
- [ ] HTTPS enabled with valid certificate
- [ ] CORS configured for production domains
- [ ] Rate limiting enabled
- [ ] Admin users created
- [ ] Email service configured and tested
- [ ] OAuth providers configured (if used)
- [ ] Notification services tested (if enabled)
- [ ] Monitoring and alerting set up
- [ ] Backup strategy implemented
- [ ] Log retention policy configured

## 🛠️ Development Tools

### Makefile Commands

```bash
make run            # Run the application
make build          # Build binary
make migrate-up     # Run migrations
make migrate-down   # Rollback migrations
make test           # Run tests
make clean          # Clean build artifacts
```

### Database Migrations

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq add_new_feature

# Apply migrations
migrate -path migrations -database "$POSTGRES_URL" up

# Rollback one migration
migrate -path migrations -database "$POSTGRES_URL" down 1

# Check migration version
migrate -path migrations -database "$POSTGRES_URL" version
```

### Current Migrations

1. `000001` - Create users table
2. `000002` - Create admins table
3. `000003` - Create OAuth providers table
4. `000004` - Create trading signals table
5. `000005` - Add password fields to users
6. `000006` - Add profile picture to users
7. `000007` - Update trading signals (asset class, duration type, free-for-all)
8. `000008` - Create packages table
9. `000009` - Create user subscriptions table
10. `000010` - Create payment history table
11. `000011` - Seed initial 18 packages

## 🔍 Troubleshooting

### Database Connection Issues

```bash
# Check if containers are running
docker-compose ps

# Check PostgreSQL
docker-compose logs postgres

# Test connection
psql -U postgres -h localhost -p 5432 -d rest_api_db

# Check ports
lsof -i :5432    # PostgreSQL
lsof -i :27017   # MongoDB
lsof -i :6379    # Redis
```

### Migration Issues

```bash
# Check current version
migrate -path migrations -database "$POSTGRES_URL" version

# Force version (if stuck)
migrate -path migrations -database "$POSTGRES_URL" force 11

# Start fresh (⚠️ Warning: deletes data)
docker-compose down -v
docker-compose up -d
make migrate-up
```

### Email Not Sending

1. Check `EMAIL_SERVICE_ENABLED=true`
2. Verify provider credentials
3. Check logs for errors: `docker-compose logs -f api`
4. Test with mock mode first: `EMAIL_SERVICE_ENABLED=false`

### OAuth Issues

1. Verify client IDs and secrets
2. Check redirect URLs match exactly (including http/https)
3. Ensure APIs are enabled (Google+ API for Google)
4. Test with OAuth playground first

### Subscription Access Issues

1. Check active subscriptions: `GET /api/subscriptions/active`
2. Verify subscription hasn't expired
3. Ensure signal's asset_class and duration_type match subscription
4. Check if signal is free-for-all

## 📄 License

MIT License - feel free to use this for your projects.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📞 Support

For issues, questions, or feature requests, please open an issue on GitHub.

---

Built with ❤️ using Go, PostgreSQL, MongoDB, and Redis
