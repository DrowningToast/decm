# ProtectedRoute Component - Usage Examples

The `ProtectedRoute` component now supports role-based access control using the `/check-roles` endpoint.

## Props

| Prop                     | Type              | Default         | Description                                  |
| ------------------------ | ----------------- | --------------- | -------------------------------------------- |
| `children`               | `React.ReactNode` | Required        | Content to render when access is granted     |
| `redirectTo`             | `Path`            | `/signup`       | Redirect path when user is not authenticated |
| `unauthorizedRedirectTo` | `Path`            | `/unauthorized` | Redirect path when role check fails          |
| `requireAuthenticated`   | `boolean`         | `true`          | Require user to be authenticated             |
| `requireHost`            | `boolean`         | `false`         | Require user to be a verified host/organizer |
| `requireIssuer`          | `boolean`         | `false`         | Require user to be a verified issuer         |
| `fallback`               | `React.ReactNode` | Default spinner | Custom loading component                     |

## Basic Usage

### 1. Protected Route (Authentication Only)

```tsx
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

export default function DashboardPage() {
    return (
        <ProtectedRoute>
            <div>Dashboard Content</div>
        </ProtectedRoute>
    );
}
```

### 2. Host-Only Route

```tsx
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

export default function CreateEventPage() {
    return (
        <ProtectedRoute requireHost={true}>
            <div>Create Event Form</div>
        </ProtectedRoute>
    );
}
```

### 3. Issuer-Only Route

```tsx
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

export default function IssueCertificatePage() {
    return (
        <ProtectedRoute requireIssuer={true}>
            <div>Issue Certificate Form</div>
        </ProtectedRoute>
    );
}
```

### 4. Multiple Role Requirements

```tsx
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

// User must be both a host AND an issuer
export default function AdvancedEventPage() {
    return (
        <ProtectedRoute requireHost={true} requireIssuer={true}>
            <div>Advanced Event Management</div>
        </ProtectedRoute>
    );
}
```

## Advanced Usage

### Custom Redirect Paths

```tsx
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

export default function AdminPage() {
    return (
        <ProtectedRoute requireHost={true} redirectTo="/signin" unauthorizedRedirectTo="/403">
            <div>Admin Panel</div>
        </ProtectedRoute>
    );
}
```

### Custom Loading Fallback

```tsx
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

const CustomLoader = () => (
    <div className="flex items-center justify-center min-h-screen">
        <span>Verifying permissions...</span>
    </div>
);

export default function HostPage() {
    return (
        <ProtectedRoute requireHost={true} fallback={<CustomLoader />}>
            <div>Host Dashboard</div>
        </ProtectedRoute>
    );
}
```

### Public Route (No Authentication Required)

```tsx
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

// Useful for pages that check roles but don't require authentication
export default function OptionalAuthPage() {
    return (
        <ProtectedRoute requireAuthenticated={false}>
            <div>Public Content</div>
        </ProtectedRoute>
    );
}
```

## Layout Patterns

### Protected Layout

```tsx
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { Outlet } from "react-router-dom";

const HostLayout = () => {
    return (
        <ProtectedRoute requireHost={true}>
            <div className="host-layout">
                <HostNavbar />
                <Outlet />
            </div>
        </ProtectedRoute>
    );
};

export default HostLayout;
```

### Nested Protection

```tsx
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

const EventManagementPage = () => {
    return (
        <ProtectedRoute>
            {/* First layer: Authentication only */}
            <div className="container">
                <h1>Event Management</h1>

                {/* Second layer: Host-only section */}
                <ProtectedRoute requireHost={true}>
                    <div className="host-section">
                        <button>Create Event</button>
                    </div>
                </ProtectedRoute>
            </div>
        </ProtectedRoute>
    );
};
```

## Behavior

### Loading States

1. **Initial Auth Check**: Shows loading spinner while checking authentication status
2. **Role Check**: If roles are required, shows loading spinner during API call
3. **Custom Fallback**: Use the `fallback` prop to customize the loading UI

### Redirect Behavior

1. **Not Authenticated**: Redirects to `redirectTo` path (default: `/signup`)
2. **Role Check Failed**: Redirects to `unauthorizedRedirectTo` path (default: `/unauthorized`)
3. **Success**: Renders children content

### Error Handling

- Shows toast error message on authentication failure
- Shows toast error message on role check failure
- Logs API errors to console for debugging
- Always redirects on error (fail-safe behavior)

## Testing

The component includes comprehensive tests covering:

- Authentication checks
- Role verification (host, issuer, both)
- Loading states
- Error handling
- Custom props
- Edge cases

See `ProtectedRoute.test.tsx` for examples.

## API Integration

The component uses the `/api/v1/auth/check-role` endpoint with query parameters:

```typescript
// Example API call
const response = await coreApi.authCheckRole({
    isHost: true,
    isIssuer: true,
});

// Response format
{
    is_host?: boolean;
    is_issuer?: boolean;
}
```

Only requested fields are returned in the response.
