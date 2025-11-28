# Zitadel Multi-Tenant Vue 3 Application

A Vue 3 TypeScript application with Zitadel multi-tenant authentication using PKCE flow.

## Features

- **Multi-Tenant Authentication**: Organization-scoped login with Zitadel
- **PKCE Flow**: Secure OAuth 2.0 authorization code flow with PKCE
- **Protected Routes**: All routes require authentication except login and callback
- **Auto Token Refresh**: Silent token refresh in the background
- **Modern UI**: Material Design with Vuetify 3
- **Type Safety**: Full TypeScript support

## Tech Stack

- **Vue 3** - Progressive JavaScript framework
- **TypeScript** - Type-safe development
- **Vite** - Fast build tool
- **Pinia** - State management
- **Vuetify 3** - Material Design component framework
- **Vue Router 4** - Official router
- **oidc-client-ts** - OpenID Connect client library

## Prerequisites

Before running this application, you need:

1. **Zitadel Backend Running**: Ensure the Go backend is running with Zitadel
   ```bash
   cd /Users/pablogarciavivo/code/github.com/github.com/zitadel
   make up
   go run main.go
   ```

2. **Client ID**: Copy the Client ID from the backend output

## Setup

### 1. Install Dependencies

```bash
npm install
```

### 2. Configure Environment Variables

Create a `.env.local` file in the root of the `/web` directory:

```bash
cp .env.example .env.local
```

Edit `.env.local` and add your Client ID from the backend output:

```bash
VITE_ZITADEL_DOMAIN=http://localhost:8080
VITE_ZITADEL_CLIENT_ID=<YOUR_CLIENT_ID_HERE>
VITE_REDIRECT_URI=http://localhost:3000/auth/callback
VITE_POST_LOGOUT_URI=http://localhost:3000
```

**To get the Client ID:**
1. Run the backend: `go run main.go`
2. Look for the output: `📋 Shared Client ID: XXXXXXXXX`
3. Copy that value to your `.env.local`

### 3. Run Development Server

```bash
npm run dev
```

The application will be available at [http://localhost:3000](http://localhost:3000)

## Usage

### Login Flow

1. Visit [http://localhost:3000](http://localhost:3000)
2. You'll be redirected to the login page
3. Enter your organization domain (e.g., "Tesla-150405", "Pepsi-150405", "Nike-150405")
4. Click "Continue to Login"
5. You'll be redirected to Zitadel's login page with your organization's branding
6. Enter your credentials (username and password)
7. After successful authentication, you'll be redirected back to the dashboard

### Test Credentials

The backend creates test organizations with the following credentials:

| Organization | Domain Format | Username | Password |
|-------------|---------------|----------|----------|
| Tesla | Tesla-{suffix} | user@Tesla.com | Password123! |
| Pepsi | Pepsi-{suffix} | user@Pepsi.com | Password123! |
| Nike | Nike-{suffix} | user@Nike.com | Password123! |

**Note**: The `{suffix}` is a timestamp generated when you run `go run main.go`. Check the backend output to see the exact domain names.

### Pages

- **`/login`** - Organization selection page (public)
- **`/auth/callback`** - OAuth callback handler (public)
- **`/dashboard`** - Main dashboard showing user info (protected)
- **`/profile`** - Detailed user profile and token information (protected)

### Logout

Click the user icon in the top-right corner and select "Logout" to sign out.

## Project Structure

```
/web
├── src/
│   ├── components/
│   │   ├── auth/
│   │   │   └── AuthCallback.vue     # OAuth callback handler
│   │   └── layout/
│   │       ├── AppHeader.vue        # Navigation header
│   │       └── AppLayout.vue        # Main layout wrapper
│   ├── views/
│   │   ├── Login.vue                # Organization selection
│   │   ├── Dashboard.vue            # User dashboard
│   │   └── Profile.vue              # User profile
│   ├── router/
│   │   ├── index.ts                 # Route definitions
│   │   └── guards.ts                # Auth guards
│   ├── stores/
│   │   └── auth.ts                  # Pinia auth store
│   ├── config/
│   │   └── oidc.ts                  # OIDC configuration
│   ├── types/
│   │   ├── auth.ts                  # Auth type definitions
│   │   └── user.ts                  # User type definitions
│   ├── plugins/
│   │   └── vuetify.ts               # Vuetify configuration
│   ├── App.vue                      # Root component
│   └── main.ts                      # Application entry
├── .env.example                     # Environment template
├── .env.local                       # Local environment (gitignored)
├── package.json
├── tsconfig.json
├── vite.config.ts
└── README.md
```

## How It Works

### Authentication Flow

1. **Organization Selection**: User enters organization domain on login page
2. **UserManager Creation**: Frontend creates a UserManager instance with organization-scoped authority
3. **PKCE Challenge**: UserManager generates code challenge and verifier
4. **Redirect to Zitadel**: User is redirected to Zitadel with:
   - Client ID
   - Organization scope: `urn:zitadel:iam:org:domain:primary:{org_domain}`
   - PKCE code challenge
5. **Zitadel Login**: User authenticates on organization's branded login page
6. **Callback**: Zitadel redirects back with authorization code
7. **Token Exchange**: Frontend exchanges code for tokens using PKCE verifier
8. **Token Storage**: Tokens are stored in localStorage
9. **Access Granted**: User can access protected routes

### Route Guards

All routes except `/login` and `/auth/callback` require authentication. The auth guard:
1. Checks if user is authenticated
2. If not, tries to load user from localStorage
3. If still not authenticated, redirects to login
4. Saves intended destination for post-login redirect

### State Management

Pinia store (`stores/auth.ts`) manages:
- User authentication state
- Selected organization
- UserManager instance
- Loading and error states
- Login/logout operations
- Token refresh

### Token Refresh

The UserManager automatically handles token refresh:
- Silent refresh enabled
- Refresh starts 60 seconds before expiration
- No user interaction required

## Building for Production

```bash
npm run build
```

Built files will be in the `dist/` directory.

## Troubleshooting

### "Client ID is not set" Warning

Make sure you've created `.env.local` with the correct Client ID from the backend.

### "Organization not found" Error

This means the organization domain you entered doesn't exist or is misspelled. Check the backend output for the exact domain names.

### "Authentication callback failed" Error

- Ensure the backend is running
- Check that the redirect URI matches what's configured in Zitadel
- Verify the Client ID is correct

### Token Expired

Tokens expire after a certain period. The app automatically refreshes them, but if you're signed out:
1. Close the app
2. Reopen and login again

### CORS Errors

- Zitadel must be running on localhost:8080
- Frontend must be running on localhost:3000
- Both are required for the OAuth flow to work

## Development

### Type Checking

```bash
npm run type-check
```

### Build

```bash
npm run build
```

## Architecture Decisions

### Why oidc-client-ts Instead of vue-oidc-context?

We use `oidc-client-ts` directly because:
- Multi-tenant requires **dynamic organization scope** based on user input
- Organization must be selected **before** starting OAuth flow
- `vue-oidc-context` expects static config at initialization
- We need to construct: `urn:zitadel:iam:org:domain:primary:{org_domain}` dynamically

### Why Shared OIDC Application?

We use the backend's shared OIDC application because:
- Already configured with correct redirect URIs
- Aligns with multi-tenant architecture
- No additional backend setup required
- All tenants share the same Client ID

## Security

- **PKCE**: Protects against authorization code interception
- **localStorage**: Used for token storage (acceptable for SPA)
- **No Client Secret**: Public client with PKCE (no secret in frontend)
- **State Parameter**: CSRF protection handled by oidc-client-ts
- **HTTPS Required in Production**: Always use HTTPS in production

## License

MIT
