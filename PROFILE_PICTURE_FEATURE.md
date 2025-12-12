# Profile Picture Feature from OAuth Providers

## Overview
This feature adds support for retrieving and storing user profile pictures from OAuth providers (Google and Facebook) during authentication. The profile picture field is optional - if the OAuth provider doesn't provide a picture, the field will be null.

## Changes Made

### 1. Database Migration
**Files**: 
- `migrations/000006_add_profile_picture.up.sql`
- `migrations/000006_add_profile_picture.down.sql`

Added a new `profile_picture` column to the `users` table to store the profile picture URL (up to 1000 characters).

```sql
ALTER TABLE users ADD COLUMN profile_picture VARCHAR(1000);
```

### 2. User Model Updates
**File**: `internal/models/user.go`

- Added `ProfilePicture *string` field to the `User` struct
- Added `ProfilePicture *string` field to the `UserCreate` struct
- Both fields use pointer types to allow null values (optional field)
- JSON tags include `omitempty` to exclude null values from API responses

### 3. OAuth Service Updates
**File**: `internal/services/oauth_service.go`

#### OAuthUserInfo Struct
- Added `Picture *string` field to capture profile picture URLs from OAuth providers

#### Google OAuth Integration
- Modified `ExchangeGoogleCode()` to extract the `Picture` field from Google's userinfo response
- Google provides the picture URL directly in the `userInfo.Picture` field

#### Facebook OAuth Integration
- Modified `ExchangeFacebookCode()` to request and extract profile picture
- Updated Facebook Graph API call to include `picture` field in the request
- Facebook returns picture in a nested structure: `picture.data.url`
- Added custom struct to parse Facebook's nested response format

### 4. User Repository Updates
**File**: `internal/repositories/user_repository.go`

Updated all user-related database queries to include the `profile_picture` field:

- `Create()` - Save profile picture when creating OAuth users
- `GetByID()` - Retrieve profile picture
- `GetByEmail()` - Retrieve profile picture
- `Update()` - Update profile picture
- `CreateWithPassword()` - Handle profile picture for password-based users
- `GetByEmailWithPassword()` - Retrieve profile picture
- `GetByVerificationToken()` - Retrieve profile picture
- `GetByResetToken()` - Retrieve profile picture

### 5. Auth Service Updates
**File**: `internal/services/auth_service.go`

#### AuthenticateWithOAuth()
- Pass profile picture to `UserCreate` when creating new users
- Update existing users' profile picture if:
  - OAuth provider provides a picture
  - User doesn't already have a profile picture (preserves user's existing picture)

#### AuthenticateWithOAuthUserInfo()
- Similar logic applied for mobile ID token verification flows
- Pass profile picture when creating new users
- Update existing users' profile picture only if not already set

### 6. Auth Handler Updates
**File**: `internal/handlers/auth_handler.go`

#### VerifyGoogleIDToken()
- Modified to extract `picture` field from Google's token info response
- Creates `OAuthUserInfo` with profile picture for authentication

#### VerifyFacebookAccessToken()
- Updated Facebook Graph API call to include `picture` field
- Parses Facebook's nested picture structure
- Creates `OAuthUserInfo` with profile picture for authentication

## API Response Example

When a user authenticates via OAuth with a profile picture, the response will include:

```json
{
  "status": "success",
  "type": "auth",
  "message": "Authentication successful",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "profile_picture": "https://lh3.googleusercontent.com/...",
      "email_verified": true,
      "blocked": false,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "is_admin": false
  }
}
```

If no profile picture is available, the field will be omitted from the JSON response due to the `omitempty` tag:

```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "email_verified": true,
    ...
  }
}
```

## Behavior Details

### Profile Picture Priority
- **New users**: Profile picture from OAuth provider is saved on first authentication
- **Existing users without profile picture**: Profile picture from OAuth provider is added
- **Existing users with profile picture**: Existing picture is preserved (not overwritten)

### Optional Nature
- The feature gracefully handles cases where OAuth providers don't provide pictures
- Profile picture updates that fail are logged but don't block authentication
- The field uses nullable types (`*string`) throughout the codebase

### Supported OAuth Providers
1. **Google OAuth**: Picture URL from `userInfo.Picture`
2. **Facebook OAuth**: Picture URL from `picture.data.url`

## Running the Migration

To apply the database migration:

```bash
# Using migrate CLI
migrate -path ./migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" up

# Or if using golang-migrate in code, it will auto-apply on startup
```

To rollback:

```bash
migrate -path ./migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" down 1
```

## Testing

After applying the migration and deploying the code:

1. Test Google OAuth authentication and verify profile picture is stored
2. Test Facebook OAuth authentication and verify profile picture is stored
3. Test that existing users without pictures get updated on next OAuth login
4. Test that existing users with pictures keep their current picture
5. Verify API responses include profile picture when available
6. Test that authentication still works if OAuth provider doesn't provide a picture

## Notes

- Profile pictures are stored as URLs, not downloaded/uploaded to your server
- The URLs point to the OAuth provider's CDN (Google/Facebook servers)
- URLs may expire based on OAuth provider policies - consider implementing a refresh mechanism in the future
- Maximum URL length is 1000 characters (adjust if needed)
- The feature is fully backward compatible with existing users
