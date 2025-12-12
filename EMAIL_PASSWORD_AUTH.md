# Email/Password Authentication Guide

Complete guide for implementing and using email/password authentication in addition to social OAuth.

## Overview

This API supports three authentication methods:
1. **Google OAuth** (enabled/disabled via config)
2. **Facebook OAuth** (enabled/disabled via config)
3. **Email/Password** (enabled/disabled via config) ← **YOU ARE HERE**

All three can be enabled simultaneously, and accounts are automatically linked by email address.

---

## ⚙️ Configuration

### Enable Email/Password Authentication

In your `.env` file:

```bash
# Enable email/password authentication
EMAIL_PASSWORD_AUTH_ENABLED=true

# Require email verification before login (optional)
REQUIRE_EMAIL_VERIFICATION=false

# Token expiry times
VERIFICATION_TOKEN_EXPIRY=24h
RESET_TOKEN_EXPIRY=1h
```

### Email Service Configuration

For sending verification and reset emails:

```bash
# Enable email service (logs to console if false)
EMAIL_SERVICE_ENABLED=false

# Email settings
EMAIL_FROM_ADDRESS=noreply@yourapp.com
EMAIL_FROM_NAME=Your App Name

# Frontend URL for email links
FRONTEND_URL=http://localhost:3000
```

**Note:** When `EMAIL_SERVICE_ENABLED=false`, emails are logged to console instead of sent. This is useful for development.

---

## 📊 Database Schema

The migration `000005_add_password_fields.up.sql` adds:

```sql
-- Password field (nullable for OAuth-only users)
hashed_password VARCHAR(255)

-- Email verification
email_verified BOOLEAN DEFAULT FALSE
verification_token VARCHAR(255)
verification_token_expires TIMESTAMP

-- Password reset
reset_token VARCHAR(255)
reset_token_expires TIMESTAMP
```

### Run Migration

```bash
make migrate-up
# OR
migrate -path migrations -database "your-postgres-url" up
```

---

## 🔐 Security Features

### Password Requirements
- Minimum 8 characters
- Validated on registration and change

### Password Hashing
- Uses bcrypt with cost factor 12
- Industry-standard secure hashing

### Token Security
- Verification tokens: 24-hour expiry (configurable)
- Reset tokens: 1-hour expiry (configurable)
- Cryptographically secure random generation
- One-time use (cleared after use)

### Rate Limiting
- All endpoints protected by IP-based rate limiting
- 100 requests per minute default

### Email Enumeration Protection
- Forgot password doesn't reveal if email exists
- Consistent responses regardless of user existence

---

## 🚀 API Endpoints

### 1. Register

**POST** `/auth/register`

Creates a new user account with email and password.

**Request:**
```json
{
    "email": "user@example.com",
    "name": "John Doe",
    "password": "SecurePass123!"
}
```

**Response (201 Created):**
```json
{
    "status": "success",
    "type": "auth",
    "data": {
        "user": {
            "id": 1,
            "email": "user@example.com",
            "name": "John Doe",
            "email_verified": false,
            "blocked": false,
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z"
        },
        "message": "Registration successful. Please check your email to verify your account."
    },
    "message": "User registered successfully"
}
```

**Errors:**
- `409 Conflict`: Email already exists
- `400 Bad Request`: Validation errors

**What happens:**
1. User account created
2. Password hashed with bcrypt
3. Verification token generated
4. Verification email sent (or logged if disabled)

---

### 2. Login

**POST** `/auth/login`

Authenticates user with email and password.

**Request:**
```json
{
    "email": "user@example.com",
    "password": "SecurePass123!"
}
```

**Response (200 OK):**
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
            "blocked": false
        },
        "access_token": "eyJhbGciOiJIUzI1NiIs...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
        "is_admin": false
    },
    "message": "Login successful"
}
```

**Cookies Set** (for web apps):
- `access_token`: 15-minute expiry
- `refresh_token`: 7-day expiry

**Errors:**
- `401 Unauthorized`: Invalid credentials
- `403 Forbidden`: Account blocked

**Note:** If `REQUIRE_EMAIL_VERIFICATION=true`, unverified users cannot login.

---

### 3. Verify Email

**GET** `/auth/verify-email?token=TOKEN`

Verifies user's email address.

**Query Parameters:**
- `token`: Verification token from email

**Response (200 OK):**
```json
{
    "status": "success",
    "type": "action",
    "data": null,
    "message": "Email verified successfully"
}
```

**Errors:**
- `400 Bad Request`: Invalid or expired token

**What happens:**
1. Token validated and expiry checked
2. `email_verified` set to `true`
3. Token cleared from database

---

### 4. Forgot Password

**POST** `/auth/forgot-password`

Initiates password reset process.

**Request:**
```json
{
    "email": "user@example.com"
}
```

**Response (200 OK):**
```json
{
    "status": "success",
    "type": "action",
    "data": null,
    "message": "If the email exists, a password reset link has been sent"
}
```

**Note:** Always returns success to prevent email enumeration.

**What happens:**
1. If email exists:
   - Reset token generated
   - Reset email sent (or logged)
   - Token expires in 1 hour
2. If email doesn't exist:
   - No action taken
   - Same response returned

---

### 5. Reset Password

**POST** `/auth/reset-password`

Resets password using reset token.

**Request:**
```json
{
    "token": "reset-token-from-email",
    "new_password": "NewSecurePass123!"
}
```

**Response (200 OK):**
```json
{
    "status": "success",
    "type": "action",
    "data": null,
    "message": "Password reset successfully"
}
```

**Errors:**
- `400 Bad Request`: Invalid/expired token or validation error

**What happens:**
1. Token validated
2. New password hashed
3. Password updated
4. Reset token cleared
5. All refresh tokens revoked (logout from all devices)
6. Confirmation email sent

---

### 6. Change Password

**POST** `/auth/change-password`

**Authentication Required:** Bearer token

Changes password for authenticated user.

**Request:**
```json
{
    "old_password": "SecurePass123!",
    "new_password": "NewSecurePass123!"
}
```

**Response (200 OK):**
```json
{
    "status": "success",
    "type": "action",
    "data": null,
    "message": "Password changed successfully. Please login again."
}
```

**Errors:**
- `401 Unauthorized`: Not authenticated
- `400 Bad Request`: Current password incorrect

**What happens:**
1. Old password verified
2. New password hashed
3. Password updated
4. All refresh tokens revoked (logout from all devices)
5. Confirmation email sent

---

### 7. Resend Verification

**POST** `/auth/resend-verification`

Resends email verification link.

**Request:**
```json
{
    "email": "user@example.com"
}
```

**Response (200 OK):**
```json
{
    "status": "success",
    "type": "action",
    "data": null,
    "message": "Verification email sent successfully"
}
```

**Errors:**
- `400 Bad Request`: Email already verified or user not found

**What happens:**
1. New verification token generated
2. New verification email sent

---

## 🔄 Account Linking

Users with the same email across different authentication methods are automatically linked:

**Example:**
1. User registers with email/password: `user@example.com`
2. Same user logs in with Google OAuth using `user@example.com`
3. Accounts are linked automatically
4. User can now login with either method

**Benefits:**
- Single user record
- Consistent user experience
- No duplicate accounts

---

## 🧪 Testing

### Using Postman

Import `postman_collection_v2.json` which includes all email/password endpoints.

### Manual Testing Flow

```bash
# 1. Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "SecurePass123!"
  }'

# 2. Check console for verification token (if EMAIL_SERVICE_ENABLED=false)
# Or check your email

# 3. Verify email
curl "http://localhost:8080/auth/verify-email?token=YOUR_TOKEN"

# 4. Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'

# 5. Use access token for authenticated requests
curl http://localhost:8080/api/trading-signals \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## 📧 Email Service Integration

### Current Implementation

The `EmailService` is a stub that logs to console. To implement actual email sending:

### Option 1: SendGrid

```go
// Install: go get github.com/sendgrid/sendgrid-go

func (s *EmailService) sendEmail(to, subject, body string) error {
    from := mail.NewEmail(s.fromName, s.fromAddress)
    toEmail := mail.NewEmail("", to)
    message := mail.NewSingleEmail(from, subject, toEmail, body, body)
    
    client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
    response, err := client.Send(message)
    
    if err != nil {
        return err
    }
    
    if response.StatusCode >= 400 {
        return fmt.Errorf("email failed with status: %d", response.StatusCode)
    }
    
    return nil
}
```

### Option 2: SMTP

```go
// Use net/smtp or gomail

import "gopkg.in/gomail.v2"

func (s *EmailService) sendEmail(to, subject, body string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", s.fromAddress)
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)
    
    d := gomail.NewDialer(
        os.Getenv("SMTP_HOST"),
        587,
        os.Getenv("SMTP_USERNAME"),
        os.Getenv("SMTP_PASSWORD"),
    )
    
    return d.DialAndSend(m)
}
```

### Option 3: AWS SES, Mailgun, etc.

Similar integration - install SDK and implement `sendEmail` method.

---

## 🔒 Production Checklist

Before deploying to production:

- [ ] Set `EMAIL_SERVICE_ENABLED=true`
- [ ] Configure real email service (SendGrid, etc.)
- [ ] Set strong JWT secrets
- [ ] Enable HTTPS (`COOKIE_SECURE=true`)
- [ ] Configure `FRONTEND_URL` to production domain
- [ ] Set `REQUIRE_EMAIL_VERIFICATION=true` (recommended)
- [ ] Configure CORS for production domain
- [ ] Set up monitoring for failed email sends
- [ ] Test password reset flow end-to-end
- [ ] Test verification email flow end-to-end
- [ ] Implement email templates (HTML)
- [ ] Add logging for security events
- [ ] Configure backup email provider

---

## 🐛 Troubleshooting

### Emails not being sent

**Check:**
1. `EMAIL_SERVICE_ENABLED=true` in `.env`
2. Email service credentials configured
3. Check console logs for errors
4. Verify email service API key/credentials

### Cannot login after registration

**Check:**
1. If `REQUIRE_EMAIL_VERIFICATION=true`, verify email first
2. Check if password meets requirements (8+ characters)
3. Verify user is not blocked in database
4. Check console logs for errors

### Token expired errors

**Check:**
1. Tokens have short expiry (verification: 24h, reset: 1h)
2. Request new token using resend-verification or forgot-password
3. Check server time is correct

### Password reset not working

**Check:**
1. Token from email matches exactly (no spaces)
2. Token not expired (1 hour default)
3. New password meets requirements
4. Check console logs for errors

---

## 📚 Additional Resources

- [API Documentation](API_DOCUMENTATION.md) - Complete API reference
- [README](README.md) - Project overview and setup
- [Postman Collection](postman_collection_v2.json) - API testing
- [Quick Start](QUICK_START.md) - 5-minute setup guide

---

## 🎯 Summary

**Email/Password Authentication is now:**
- ✅ Fully implemented
- ✅ Disabled by default (`EMAIL_PASSWORD_AUTH_ENABLED=false`)
- ✅ Easy to enable via environment variable
- ✅ Secure (bcrypt, tokens, rate limiting)
- ✅ Production-ready (just needs email service integration)
- ✅ Fully documented
- ✅ Includes all flows (register, login, verify, reset, change)
- ✅ Compatible with OAuth (account linking)

Enable it when ready by setting `EMAIL_PASSWORD_AUTH_ENABLED=true` in your `.env` file!

