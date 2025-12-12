# Quick Start Guide

Get your Go REST API up and running in 5 minutes!

## Prerequisites Check

Before you begin, ensure you have:
- ✅ Go 1.25+ installed (`go version`)
- ✅ Docker and Docker Compose installed (`docker --version`)
- ✅ Git installed (`git --version`)

## Step-by-Step Setup

### 1. Clone and Navigate
```bash
git clone <your-repo-url>
cd rest-api-with-social-auth
```

### 2. Install Go Dependencies
```bash
go mod download
```

### 3. Start Databases
```bash
docker-compose up -d
```

Wait about 5 seconds for databases to initialize.

### 4. Generate JWT Secrets
```bash
go run scripts/generate_secrets.go
```

Copy the generated secrets for the next step.

### 5. Configure Environment
```bash
cp env.example .env
```

Edit `.env` and update at minimum:
```bash
JWT_ACCESS_SECRET=<paste-generated-access-secret>
JWT_REFRESH_SECRET=<paste-generated-refresh-secret>
```

### 6. Run Database Migrations
```bash
# Install golang-migrate if you haven't already
# macOS: brew install golang-migrate
# Linux: See https://github.com/golang-migrate/migrate

# Run migrations
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/rest_api_db?sslmode=disable" up
```

Or use the Makefile:
```bash
make migrate-up
```

### 7. Start the Server
```bash
go run cmd/api/main.go
```

Or use the Makefile:
```bash
make run
```

### 8. Test the API
```bash
curl http://localhost:8080/health
```

You should see a healthy response! 🎉

## What's Next?

### Create an Admin User
Connect to PostgreSQL and run:
```bash
docker exec -it rest_api_postgres psql -U postgres -d rest_api_db
```

Then execute:
```sql
-- Create a user
INSERT INTO users (email, name, blocked) 
VALUES ('admin@example.com', 'Admin User', false);

-- Make them an admin
INSERT INTO admins (user_id) 
VALUES ((SELECT id FROM users WHERE email = 'admin@example.com'));

-- Exit
\q
```

### Configure OAuth (Optional)

#### Google OAuth:
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create OAuth credentials
3. Add redirect URI: `http://localhost:8080/auth/google/callback`
4. Update `.env`:
   ```bash
   OAUTH_GOOGLE_ENABLED=true
   OAUTH_GOOGLE_CLIENT_ID=your-client-id
   OAUTH_GOOGLE_CLIENT_SECRET=your-client-secret
   ```

#### Facebook OAuth:
1. Go to [Facebook Developers](https://developers.facebook.com/)
2. Create an app and add Facebook Login
3. Add redirect URI: `http://localhost:8080/auth/facebook/callback`
4. Update `.env`:
   ```bash
   OAUTH_FACEBOOK_ENABLED=true
   OAUTH_FACEBOOK_CLIENT_ID=your-app-id
   OAUTH_FACEBOOK_CLIENT_SECRET=your-app-secret
   ```

### Test the Endpoints

Visit in browser:
- Health: http://localhost:8080/health
- Google OAuth: http://localhost:8080/auth/google (if configured)
- Facebook OAuth: http://localhost:8080/auth/facebook (if configured)

## Using Makefile Commands

We've included a Makefile for common tasks:

```bash
make help           # Show all available commands
make build          # Build the application
make run            # Run the application
make test           # Run tests
make docker-up      # Start Docker containers
make docker-down    # Stop Docker containers
make migrate-up     # Run database migrations
make migrate-down   # Rollback migrations
make setup          # Complete setup (docker + migrations)
make dev            # Start development (docker + migrations + run)
```

## Quick Development Workflow

For daily development, just run:
```bash
make dev
```

This will:
1. Start Docker containers
2. Run migrations
3. Start the server

## Troubleshooting

### Port Already in Use
If you get "address already in use" errors:
```bash
# Find and kill process on port 8080
lsof -ti:8080 | xargs kill -9

# Or change port in .env
SERVER_PORT=3000
```

### Database Connection Failed
```bash
# Check if containers are running
docker-compose ps

# Restart containers
docker-compose restart

# Check logs
docker-compose logs postgres
docker-compose logs mongodb
docker-compose logs redis
```

### Migration Errors
```bash
# Reset database (⚠️ deletes all data)
docker-compose down -v
docker-compose up -d
sleep 5
make migrate-up
```

## Project Structure Overview

```
rest-api-with-social-auth/
├── cmd/api/main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration
│   ├── database/                # DB connections
│   ├── handlers/                # HTTP handlers
│   ├── middleware/              # Middleware
│   ├── models/                  # Data models
│   ├── repositories/            # Data access
│   ├── services/                # Business logic
│   └── utils/                   # Utilities
├── migrations/                  # SQL migrations
├── docker-compose.yml           # Docker setup
├── .env                         # Configuration (create from env.example)
└── README.md                    # Full documentation
```

## API Testing

### With cURL
```bash
# Health check
curl http://localhost:8080/health

# List trading signals (after OAuth login)
curl -H "Authorization: Bearer <token>" \
     http://localhost:8080/api/trading-signals
```

### With Postman
1. Import the endpoints from `API_DOCUMENTATION.md`
2. Set base URL: `http://localhost:8080`
3. Authenticate via OAuth
4. Copy access_token to Postman environment

## Next Steps

1. **Read the full README.md** for detailed documentation
2. **Check API_DOCUMENTATION.md** for complete API reference
3. **Customize the code** for your specific needs
4. **Deploy to production** following the deployment guide in README.md

## Need Help?

- 📖 Read the [README.md](README.md)
- 📚 Check [API_DOCUMENTATION.md](API_DOCUMENTATION.md)
- 🐛 Check logs: `docker-compose logs -f`
- 💬 Review the code comments

## Development Tips

1. **Use the Makefile** - It makes life easier!
2. **Check logs** - MongoDB logs all requests for debugging
3. **Test endpoints** - Use the health endpoint to verify setup
4. **Environment variables** - Never commit your .env file
5. **Database migrations** - Always create migrations for schema changes

Happy coding! 🚀

