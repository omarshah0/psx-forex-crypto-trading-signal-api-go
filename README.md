# Go REST API with Social Authentication

A production-ready Go REST API boilerplate with PostgreSQL, MongoDB, Redis, social authentication (Google & Facebook), JWT tokens, request logging, and role-based access control.

## Features

- 🔐 **Social Authentication**: Google and Facebook OAuth2 integration with account linking and profile picture sync
- 🎫 **JWT Tokens**: Secure access (15min) and refresh (7-15 days) tokens with Redis storage
- 🗄️ **Multiple Databases**: PostgreSQL for data, MongoDB for logs, Redis for caching
- 📝 **Request Logging**: Comprehensive request/response logging with sensitive data masking
- 🛡️ **Security**: Rate limiting, CORS, blocked users, admin-only endpoints
- 🍪 **Dual Client Support**: Cookie-based auth for web, Bearer token for mobile apps
- 📊 **Trading Signals**: CRUD operations with admin-only write access
- 👤 **User Profiles**: Automatic profile picture import from OAuth providers
- ⚙️ **Configurable**: All settings via environment variables
- 🚀 **Production Ready**: Graceful shutdown, health checks, structured logging

## Tech Stack

- **Language**: Go 1.25+
- **Router**: Gorilla Mux
- **Databases**: PostgreSQL (latest), MongoDB (latest), Redis (latest)
- **Authentication**: OAuth2, JWT
- **Migrations**: golang-migrate

## Project Structure

```
rest-api-with-social-auth/
├── cmd/api/                    # Application entry point
│   └── main.go
├── internal/
│   ├── config/                 # Configuration management
│   ├── database/               # Database connections
│   ├── models/                 # Data models
│   ├── repositories/           # Data access layer
│   ├── services/               # Business logic
│   ├── handlers/               # HTTP handlers
│   ├── middleware/             # HTTP middleware
│   └── utils/                  # Utility functions
├── migrations/                 # Database migrations
├── docker-compose.yml          # Local development databases
├── env.example                 # Environment variables template
└── go.mod                      # Go dependencies
```

## Getting Started

### Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose (for local databases)
- Google OAuth credentials (optional)
- Facebook OAuth credentials (optional)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd rest-api-with-social-auth
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Start databases with Docker Compose**
   ```bash
   docker-compose up -d
   ```

4. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env and add your configuration
   ```

5. **Run database migrations**
   ```bash
   migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/rest_api_db?sslmode=disable" up
   ```

6. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```

The server will start on `http://localhost:8080`

## Configuration

All configuration is done through environment variables. See `env.example` for all available options.

### Required Environment Variables

```bash
# JWT Secrets (REQUIRED - Generate secure random strings)
JWT_ACCESS_SECRET=your-super-secret-access-key
JWT_REFRESH_SECRET=your-super-secret-refresh-key

# OAuth Credentials (if enabled)
OAUTH_GOOGLE_CLIENT_ID=your-google-client-id
OAUTH_GOOGLE_CLIENT_SECRET=your-google-client-secret
OAUTH_FACEBOOK_CLIENT_ID=your-facebook-app-id
OAUTH_FACEBOOK_CLIENT_SECRET=your-facebook-app-secret
```

### OAuth Setup

#### Google OAuth
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URI: `http://localhost:8080/auth/google/callback`
6. Copy Client ID and Client Secret to `.env`

#### Facebook OAuth
1. Go to [Facebook Developers](https://developers.facebook.com/)
2. Create a new app
3. Add Facebook Login product
4. Configure Valid OAuth Redirect URIs: `http://localhost:8080/auth/facebook/callback`
5. Copy App ID and App Secret to `.env`

## API Endpoints

### Health Check
- `GET /health` - Check service health

### Authentication
- `GET /auth/google` - Initiate Google OAuth flow
- `GET /auth/google/callback` - Google OAuth callback
- `POST /auth/google/exchange` - Exchange Google auth code for JWT tokens (includes profile picture)
- `POST /auth/google/verify` - Verify Google ID token (includes profile picture)
- `GET /auth/facebook` - Initiate Facebook OAuth flow
- `GET /auth/facebook/callback` - Facebook OAuth callback
- `POST /auth/facebook/exchange` - Exchange Facebook auth code for JWT tokens (includes profile picture)
- `POST /auth/facebook/verify` - Verify Facebook access token (includes profile picture)
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - Logout (requires authentication)

### Trading Signals (Authenticated)
- `GET /api/trading-signals` - List all signals (authenticated users)
- `GET /api/trading-signals/{id}` - Get single signal (authenticated users)
- `POST /api/trading-signals` - Create signal (admin only)
- `PUT /api/trading-signals/{id}` - Update signal (admin only)
- `DELETE /api/trading-signals/{id}` - Delete signal (admin only)

## Response Format

### Success Response
```json
{
  "status": "success",
  "type": "resource",
  "data": {
    "id": "123",
    "name": "John Doe"
  },
  "message": "User fetched successfully"
}
```

### Error Response
```json
{
  "status": "error",
  "type": "bad_request",
  "error": {
    "code": 400,
    "message": "Invalid request parameters"
  }
}
```

## Authentication Methods

### For Web Applications (Cookie-based)
Tokens are automatically set as HTTP-only cookies during OAuth flow.

### For Mobile Applications (Bearer Token)
Include the access token in the Authorization header:
```
Authorization: Bearer <access_token>
```

## Database Migrations

The project includes migrations for all database schema changes, including the profile picture feature added in migration `000006_add_profile_picture`.

### Create a new migration
```bash
migrate create -ext sql -dir migrations -seq create_table_name
```

### Run migrations
```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/rest_api_db?sslmode=disable" up
```

### Rollback migrations
```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/rest_api_db?sslmode=disable" down
```

### Current Migrations
1. `000001` - Create users table
2. `000002` - Create admins table
3. `000003` - Create OAuth providers table
4. `000004` - Create trading signals table
5. `000005` - Add password fields to users
6. `000006` - Add profile picture to users (for OAuth integration)

## Testing with Postman

A complete Postman collection is included for easy API testing:

1. **Import Collection**: Import `postman_collection.json` into Postman
2. **Import Environment**: Import `postman_environment.json` (optional)
3. **Set Tokens**: After OAuth login, set `access_token` and `refresh_token` in environment
4. **Start Testing**: All endpoints are pre-configured with Bearer token authentication

See [POSTMAN_GUIDE.md](POSTMAN_GUIDE.md) for detailed instructions.

## Creating Admin Users

Admins are created manually in the database:

```sql
-- First, create or identify a user
INSERT INTO users (email, name, blocked) VALUES ('admin@example.com', 'Admin User', false);

-- Then, make them an admin
INSERT INTO admins (user_id) VALUES ((SELECT id FROM users WHERE email = 'admin@example.com'));
```

## Blocking Users

To block a user from accessing the application:

```sql
UPDATE users SET blocked = true WHERE email = 'user@example.com';
```

## Request Logging

All requests are logged to MongoDB with:
- Request method, path, headers, body, query params
- Response status, headers, body
- User ID (if authenticated)
- IP address and duration

### Sensitive Data Masking

Configure sensitive keys in `.env`:
```bash
LOG_SENSITIVE_KEYS=password,hashed_password,token,secret,access_token,refresh_token
```

Values for these keys will be replaced with `****` in logs.

### Ignored Keys

Keys can be completely excluded from logs:
```bash
LOG_IGNORED_KEYS=internal_field,debug_info
```

## Rate Limiting

Default rate limit is 100 requests per minute per IP address. Configurable in `cmd/api/main.go`.

## Security Best Practices

1. **Always use HTTPS in production** - Set `COOKIE_SECURE=true`
2. **Keep secrets secure** - Never commit `.env` file
3. **Rotate JWT secrets regularly**
4. **Use strong passwords for databases**
5. **Configure CORS properly** - Update allowed origins in production
6. **Enable rate limiting**
7. **Monitor logs for suspicious activity**

## Docker Compose Services

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Stop and remove volumes (⚠️ deletes all data)
docker-compose down -v
```

## Documentation

- [API Documentation](./API_DOCUMENTATION.md) - Complete API reference
- [Quick Start Guide](./QUICK_START.md) - Get started in minutes
- [Email/Password Auth Guide](./EMAIL_PASSWORD_AUTH.md) - Email/password authentication setup
- [Profile Picture Feature](./PROFILE_PICTURE_FEATURE.md) - OAuth profile picture implementation details
- [Postman Collection Guide](./POSTMAN_GUIDE.md) - How to use the included Postman collection

## OAuth Flow

**Code Exchange (for React Web & React Native):**
1. Frontend opens Google/Facebook OAuth
2. OAuth redirects back with code
3. Frontend sends code to: `POST /auth/google/exchange` or `POST /auth/facebook/exchange`
4. Backend fetches user info including profile picture from OAuth provider
5. Backend returns JWT tokens and user data (including profile picture if available)

**SDK Flow (for React Native with native SDKs):**
1. Use Google Sign-In SDK or Facebook SDK
2. Get ID token or access token
3. Send to: `POST /auth/google/verify` or `POST /auth/facebook/verify`
4. Backend verifies token and fetches profile picture from OAuth provider
5. Backend returns JWT tokens and user data (including profile picture if available)

### Profile Picture Handling

When users authenticate via OAuth providers:
- **Google**: Profile picture URL is automatically fetched from Google's user info
- **Facebook**: Profile picture URL is automatically fetched from Facebook Graph API
- **New Users**: Profile picture is saved during account creation
- **Existing Users**: Profile picture is updated only if not already set (preserves user's existing picture)
- **Optional Field**: If OAuth provider doesn't provide a picture, authentication continues normally

The profile picture URL is included in the user object returned after successful authentication:

```json
{
  "status": "success",
  "type": "auth",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "profile_picture": "https://lh3.googleusercontent.com/...",
      "email_verified": true,
      "blocked": false
    },
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "is_admin": false
  }
}
```

## Development

### Building
```bash
go build -o bin/api cmd/api/main.go
```

### Running
```bash
./bin/api
```

### Testing
```bash
go test ./...
```

## Production Deployment

1. Set all environment variables securely
2. Use managed database services (not Docker containers)
3. Enable HTTPS with valid SSL certificates
4. Set `ENVIRONMENT=production`
5. Configure proper CORS origins
6. Enable `COOKIE_SECURE=true`
7. Use a process manager (systemd, supervisord) or container orchestration
8. Set up monitoring and alerting
9. Configure backup strategies for databases
10. Use a reverse proxy (nginx, caddy) for SSL termination

## Troubleshooting

### Database connection issues
- Ensure Docker containers are running: `docker-compose ps`
- Check database URLs in `.env`
- Verify ports are not in use: `lsof -i :5432`, `lsof -i :27017`, `lsof -i :6379`

### OAuth errors
- Verify client IDs and secrets
- Check redirect URLs match exactly
- Ensure OAuth providers are enabled in `.env`

### Token validation errors
- Check JWT secrets are set correctly
- Verify Redis is running
- Ensure tokens haven't expired

## License

MIT License - feel free to use this boilerplate for your projects.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

