# Authentication Setup Guide

This application now includes Google OAuth authentication with role-based access control.

## Features

- **Google OAuth Sign-in**: Users sign in using their Google accounts
- **Role-based Access Control**: Two user roles - `standard` and `admin`
- **Admin User**: Only `tyler.osborn@impact.com` has admin access
- **Protected Routes**: 
  - Standard users can access: Home, Give Meal, Receive Meal pages
  - Admin users can access: All standard pages + Admin panel

## Setup Instructions

### 1. Google OAuth Configuration

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google+ API
4. Go to "Credentials" → "Create Credentials" → "OAuth 2.0 Client IDs"
5. Configure the OAuth consent screen
6. Add authorized redirect URIs:
   - Development: `http://localhost:8080/Api/Auth/Callback`
   - Production: `https://yourdomain.com/Api/Auth/Callback`

### 2. Environment Variables

Copy `.env.example` to `.env` and fill in the required values:

```bash
cp .env.example .env
```

Required environment variables:
- `GOOGLE_CLIENT_ID`: Your Google OAuth client ID
- `GOOGLE_CLIENT_SECRET`: Your Google OAuth client secret
- `GOOGLE_REDIRECT_URL`: OAuth callback URL (default: http://localhost:8080/Api/Auth/Callback)
- `SESSION_SECRET`: Random string for session encryption
- Database configuration variables (MYSQL_*)

### 3. Database Migration

The User model has been updated with new fields:
- `email`: User's Google email (unique)
- `google_id`: Google OAuth user ID (unique)
- `role`: User role ('standard' or 'admin')

The existing migration will automatically update the database schema.

### 4. Running the Application

1. Start the backend:
   ```bash
   go run main.go
   ```

2. Start the frontend:
   ```bash
   cd frontend
   npm run dev
   ```

### 5. Authentication Flow

1. Users visit the application
2. Unauthenticated users see a login prompt
3. Click "Sign in with Google" redirects to Google OAuth
4. After successful authentication, users are redirected back
5. Admin users (tyler.osborn@impact.com) get admin role automatically
6. All other users get standard role

### 6. API Endpoints

**Public:**
- `GET /Api/Auth/Login` - Get Google OAuth URL
- `GET /Api/Auth/Callback` - OAuth callback handler
- `GET /Api/Meal` - Get meals
- `GET /Api/Meal/Today` - Get today's meals

**Authenticated:**
- `POST /Api/Auth/Logout` - Logout user
- `GET /Api/Auth/Profile` - Get user profile
- `POST /Api/Donation` - Create donation
- `POST /Api/Donation/Claim` - Claim donation

**Admin Only:**
- `POST /Api/Meal/Upload` - Upload meals
- `GET /Api/Stats/Claims/Summary` - Get donation statistics

### 7. Frontend Routes

- `/` - Home page (public)
- `/login` - Login page (public)
- `/give-meal` - Give meal page (authenticated)
- `/receive-meal` - Receive meal page (authenticated)
- `/admin` - Admin panel (admin only)

The frontend includes route guards that automatically redirect unauthenticated users to the login page and prevent non-admin users from accessing admin routes.