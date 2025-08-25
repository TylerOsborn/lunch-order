# Admin Authentication System

## Overview
Basic authentication system added to protect the AdminPage (`/admin` route) of the lunch-order application. Authentication is now handled server-side with password stored as an environment variable.

## Implementation Details

### Password
- **Default Password:** `admin123` (fallback if environment variable not set)
- **Environment Variable:** `ADMIN_PASSWORD` (recommended for production)
- **Validation:** Server-side via `/Api/Admin/Login` endpoint

### Files Added/Modified

1. **`main.go`** (MODIFIED)
   - Added `/Api/Admin/Login` endpoint for authentication
   - Reads `ADMIN_PASSWORD` environment variable
   - Falls back to default password if environment variable not set
   - Returns JSON response with authentication status

2. **`frontend/src/composables/useAuth.ts`** (MODIFIED)
   - Updated to use server-side authentication endpoint
   - Async login function that calls `/Api/Admin/Login`
   - LocalStorage persistence for authentication state
   - Removed hardcoded password

3. **`frontend/src/components/AdminLogin.vue`** (MODIFIED)
   - Updated to handle async authentication
   - Error handling for network requests
   - Clean PrimeVue-styled interface

4. **`frontend/src/components/AdminScreen.vue`** (EXISTING)
   - Conditional rendering based on authentication state
   - Logout functionality with confirmation
   - Admin dashboard header with logout button

5. **`frontend/src/main.ts`** (EXISTING)
   - Route metadata for authentication requirements

## Features

### Authentication Flow
1. User navigates to `/admin`
2. If not authenticated, login form is displayed
3. User enters password and clicks Login
4. Frontend sends password to `/Api/Admin/Login` endpoint
5. Server validates password against `ADMIN_PASSWORD` environment variable
6. On successful authentication:
   - Server returns success response
   - Frontend saves authentication state to localStorage
   - User sees admin dashboard
7. User can logout using the Logout button
8. Authentication state persists across browser sessions

### Security Features
- **Server-side password validation**: Password never exposed to frontend
- **Environment variable storage**: Password stored as server secret
- **Secure API endpoint**: Password validation happens on backend
- **State persistence**: Authentication state saved locally
- **Logout functionality**: Complete state cleanup
- **Error handling**: Proper error responses for invalid credentials
- **Toast notifications**: User feedback for all authentication actions

### UI/UX Features
- Clean, responsive login form
- Consistent styling with existing PrimeVue components
- Loading states during authentication
- Success/error toast messages
- Enter key support for login form
- Async request handling

## Environment Setup

### Setting Admin Password
For production deployment on fly.io:
```bash
flyctl secrets set ADMIN_PASSWORD=your_secure_password
```

For local development:
```bash
export ADMIN_PASSWORD=your_secure_password
# or add to .env file:
echo "ADMIN_PASSWORD=your_secure_password" >> .env
```

If no environment variable is set, the system falls back to the default password `admin123`.

## API Endpoints

### POST /Api/Admin/Login
Authenticates admin user with password.

**Request:**
```json
{
  "password": "your_password"
}
```

**Response (Success - 200):**
```json
{
  "StatusCode": 200,
  "Data": {
    "authenticated": true
  }
}
```

**Response (Invalid Password - 401):**
```json
{
  "StatusCode": 401,
  "Error": "Invalid password"
}
```

## Usage

### Login
1. Navigate to `/admin`
2. Enter the admin password
3. Click "Login" or press Enter
4. System validates password with server

### Logout
1. Click the "Logout" button in the admin dashboard header
2. User is returned to login form

### Changing Password
To change the admin password:
1. Update the `ADMIN_PASSWORD` environment variable on the server
2. For fly.io: `flyctl secrets set ADMIN_PASSWORD=new_password`
3. No application rebuild required

## Technical Notes

- **Authentication Type:** Server-side password validation
- **State Management:** Vue 3 Composition API with reactive refs
- **Persistence:** localStorage for frontend state
- **API Integration:** RESTful endpoint for authentication
- **Framework Integration:** Vue Router, PrimeVue components
- **Error Handling:** Toast notifications via PrimeVue ToastService
- **Security:** Password stored as environment variable, never exposed to frontend

## Security Improvements

This updated system provides:
- ✅ Server-side password validation
- ✅ Environment variable password storage
- ✅ No password exposure in frontend code
- ✅ Secure API communication
- ✅ Proper error handling
- ✅ Session persistence

## Future Enhancements

For production environments, consider:
- JWT token-based authentication
- Password hashing (bcrypt)
- Session management with expiration
- Rate limiting for login attempts
- Role-based access control
- Multi-factor authentication
- Audit logging

## Security Considerations

This system provides basic but secure authentication suitable for admin access control. The password is stored as a server environment variable and validation happens server-side, preventing exposure in frontend code. For enhanced security in production, consider implementing additional measures like JWT tokens, rate limiting, and password hashing.