# Admin Authentication System

## Overview
Basic authentication system added to protect the AdminPage (`/admin` route) of the lunch-order application.

## Implementation Details

### Password
- **Default Password:** `admin123`
- Location: `frontend/src/composables/useAuth.ts` (ADMIN_PASSWORD constant)

### Files Added/Modified

1. **`frontend/src/composables/useAuth.ts`** (NEW)
   - Authentication composable managing login/logout state
   - LocalStorage persistence for authentication state
   - Simple password validation

2. **`frontend/src/components/AdminLogin.vue`** (NEW)
   - Login form with password field
   - Error handling and toast notifications
   - Clean PrimeVue-styled interface

3. **`frontend/src/components/AdminScreen.vue`** (MODIFIED)
   - Conditional rendering based on authentication state
   - Added logout functionality with confirmation
   - Admin dashboard header with logout button

4. **`frontend/src/main.ts`** (MODIFIED)
   - Added route metadata for authentication requirements
   - Route guard setup (currently minimal, handled by component)

## Features

### Authentication Flow
1. User navigates to `/admin`
2. If not authenticated, login form is displayed
3. User enters password and clicks Login
4. On successful authentication:
   - User sees admin dashboard
   - Authentication state is saved to localStorage
5. User can logout using the Logout button
6. Authentication state persists across browser sessions

### Security Features
- Password-based authentication
- State persistence with localStorage
- Logout functionality with state cleanup
- Error handling for invalid passwords
- Toast notifications for user feedback

### UI/UX Features
- Clean, responsive login form
- Consistent styling with existing PrimeVue components
- Loading states during authentication
- Success/error toast messages
- Enter key support for login form

## Usage

### Login
1. Navigate to `/admin`
2. Enter password: `admin123`
3. Click "Login" or press Enter

### Logout
1. Click the "Logout" button in the admin dashboard header
2. User is returned to login form

### Changing Password
To change the admin password:
1. Edit `ADMIN_PASSWORD` constant in `frontend/src/composables/useAuth.ts`
2. Rebuild the application

## Technical Notes

- **Authentication Type:** Frontend-only, password-based
- **State Management:** Vue 3 Composition API with reactive refs
- **Persistence:** localStorage
- **Framework Integration:** Vue Router, PrimeVue components
- **Error Handling:** Toast notifications via PrimeVue ToastService

## Future Enhancements

For production environments, consider:
- Backend authentication with JWT tokens
- Password hashing and secure storage
- Session management
- Role-based access control
- Multi-factor authentication
- Audit logging

## Security Considerations

This is a basic authentication system suitable for simple admin access control. The password is stored in plain text in the source code and should be changed before deployment. For production use, implement proper backend authentication.