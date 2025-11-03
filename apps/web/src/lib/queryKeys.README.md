# React Query Keys - Centralized Management

This document explains how to use the centralized query keys system for React Query in the DECM application.

## Overview

All React Query keys are centralized in `/lib/queryKeys.ts` to:

- **Prevent typos** when referencing query keys
- **Enable autocomplete** for better developer experience
- **Make refactoring easier** by having a single source of truth
- **Improve cache invalidation** by providing consistent key structures

## File Location

```
apps/web/src/lib/queryKeys.ts
```

## Basic Usage

### Importing

```tsx
import { queryKeys } from "@/lib/queryKeys";
```

### Using in `useQuery`

```tsx
import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/queryKeys";
import { coreApiClient } from "@/lib/api/api";

export function useEvent(eventId: string) {
    const { data, isLoading } = useQuery({
        queryKey: queryKeys.event.byId(eventId),
        queryFn: () => coreApiClient.v1.getEventById({ eventId }),
    });

    return { data, isLoading };
}
```

### Using in Cache Invalidation

```tsx
import { useMutation } from "@tanstack/react-query";
import { queryClient } from "@/lib/api/queryClient";
import { queryKeys } from "@/lib/queryKeys";

export function useUpdateEvent(eventId: string) {
    const { mutateAsync } = useMutation({
        mutationFn: (data) => updateEventApi(eventId, data),
        onSuccess: () => {
            // Invalidate specific event
            queryClient.invalidateQueries({
                queryKey: queryKeys.event.byId(eventId),
            });

            // Or invalidate all events
            queryClient.invalidateQueries({
                queryKey: queryKeys.event.all,
            });
        },
    });

    return { mutateAsync };
}
```

## Available Query Keys

### User & Authentication

```tsx
// User profile
queryKeys.user.profile;
// Result: ["user", "profile"]
```

### Onboarding

```tsx
// Sign message
queryKeys.onboard.signMessage;
// Result: ["getSignMessage"]

// Onboard status - all
queryKeys.onboard.status.all;
// Result: ["onboardStatus"]

// Onboard status - wallet
queryKeys.onboard.status.wallet(signSignature);
// Result: ["onboardStatus", signSignature]

// Onboard status - Google
queryKeys.onboard.status.google(accessToken, expiresIn);
// Result: ["onboardStatus", accessToken, expiresIn]
```

### Images

```tsx
// Image by URL
queryKeys.image.byUrl(url);
// Result: ["image", url]
```

### Events

```tsx
// All events (for invalidation)
queryKeys.event.all;
// Result: ["event"]

// Event by ID
queryKeys.event.byId(eventId);
// Result: ["event", eventId]

// Event registration config
queryKeys.event.registrationConfig(eventId);
// Result: ["event-registration-config", eventId]

// Event issuers
queryKeys.event.issuers.byEventId(eventId);
// Result: ["event", eventId, "issuers"]

// Event certificate config
queryKeys.event.certificate.config(eventId);
// Result: ["event", eventId, "certificate", "config"]
```

### Host Events

```tsx
// All host events (for invalidation)
queryKeys.hostEvents.all;
// Result: ["host-events"]

// Host events list with pagination
queryKeys.hostEvents.list(userId, rowsPerPage, offset);
// Result: ["host-events", userId, rowsPerPage, offset]
```

### Issuers

```tsx
// Verified issuers
queryKeys.issuers.verified;
// Result: ["issuers"]
```

## Patterns and Best Practices

### 1. Query Key Structure

Query keys follow a hierarchical structure:

```tsx
queryKeys.{domain}.{operation}(params)
```

**Examples:**

- `queryKeys.event.byId(eventId)` - Single event
- `queryKeys.event.issuers.byEventId(eventId)` - Event's issuers
- `queryKeys.event.certificate.config(eventId)` - Event's certificate config

### 2. Static vs Dynamic Keys

**Static keys** (constant arrays):

```tsx
queryKeys.user.profile;
queryKeys.onboard.signMessage;
queryKeys.issuers.verified;
```

**Dynamic keys** (factory functions):

```tsx
queryKeys.event.byId(eventId);
queryKeys.image.byUrl(url);
queryKeys.hostEvents.list(userId, rowsPerPage, offset);
```

### 3. Invalidation Patterns

**Invalidate specific item:**

```tsx
queryClient.invalidateQueries({
    queryKey: queryKeys.event.byId(eventId),
});
```

**Invalidate all items in a domain:**

```tsx
queryClient.invalidateQueries({
    queryKey: queryKeys.event.all,
});
```

This invalidates all queries starting with `["event"]`, including:

- `["event", "123"]`
- `["event", "456", "issuers"]`
- `["event", "789", "certificate", "config"]`

### 4. TypeScript Benefits

All query keys are typed as `const` arrays, providing:

- **Type safety**: Compile-time checks for key structure
- **Autocomplete**: IDE suggestions for available keys
- **Immutability**: Keys cannot be accidentally modified

```tsx
// ✅ Good - Type-safe and autocomplete works
queryKeys.event.byId(eventId)[
    // ❌ Bad - No type safety, prone to typos
    ("event", eventId)
];
```

## Examples from the Codebase

### Example 1: Basic Query Hook

```tsx
// useEvent.ts
import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/queryKeys";
import { coreApiClient } from "@/lib/api/api";

export function useEvent(eventId: string) {
    const { data: event, isLoading } = useQuery({
        queryKey: queryKeys.event.byId(eventId),
        queryFn: () => coreApiClient.v1.getEventById({ eventId }),
    });

    return { event, isLoading };
}
```

### Example 2: Mutation with Cache Invalidation

```tsx
// useUpdateCertificateConfig.ts
import { useMutation } from "@tanstack/react-query";
import { queryClient } from "@/lib/api/queryClient";
import { queryKeys } from "@/lib/queryKeys";

export function useUpdateCertificateConfig(eventId: string) {
    const { mutateAsync } = useMutation({
        mutationFn: (data) => coreApiClient.v1.updateEventCertificateConfig({ eventId }, data),
        onSuccess: () => {
            // Invalidate all event-related queries
            queryClient.invalidateQueries({
                queryKey: queryKeys.event.all,
            });
        },
    });

    return { mutateAsync };
}
```

### Example 3: Conditional Query

```tsx
// useHostEvents.ts
import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/queryKeys";
import { useAuth } from "@/context/AuthContext";

export function useHostEvents(rowsPerPage: number, offset: number) {
    const { user } = useAuth();

    const { data, isLoading } = useQuery({
        queryKey: queryKeys.hostEvents.list(
            user?.authentication_credential_id,
            rowsPerPage,
            offset,
        ),
        queryFn: () =>
            coreApiClient.v1.getEventsByOwnerCredentialsId({
                ownerCredentialId: user?.authentication_credential_id ?? "",
                limit: rowsPerPage,
                offset: offset,
            }),
        enabled: !!user?.authentication_credential_id,
    });

    return { data, isLoading };
}
```

## Adding New Query Keys

When adding a new query key, follow these steps:

### 1. Add to `queryKeys.ts`

```tsx
export const queryKeys = {
    // ... existing keys ...

    // New domain
    certificates: {
        all: ["certificates"] as const,
        byId: (certificateId: string) => ["certificates", certificateId] as const,
        byUserId: (userId: string) => ["certificates", "user", userId] as const,
    },
} as const;
```

### 2. Use in Your Hook

```tsx
import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/queryKeys";

export function useCertificate(certificateId: string) {
    const { data } = useQuery({
        queryKey: queryKeys.certificates.byId(certificateId),
        queryFn: () => fetchCertificate(certificateId),
    });

    return { data };
}
```

### 3. Update This Documentation

Add the new key to the "Available Query Keys" section above.

## Migration Guide

To migrate existing code from hardcoded query keys:

**Before:**

```tsx
useQuery({
    queryKey: ["event", eventId],
    queryFn: () => fetchEvent(eventId),
});
```

**After:**

```tsx
import { queryKeys } from "@/lib/queryKeys";

useQuery({
    queryKey: queryKeys.event.byId(eventId),
    queryFn: () => fetchEvent(eventId),
});
```

## Benefits Recap

✅ **Type Safety**: Catch errors at compile time  
✅ **Autocomplete**: IDE helps you find the right key  
✅ **Consistency**: Same key structure everywhere  
✅ **Refactoring**: Change once, update everywhere  
✅ **Documentation**: Self-documenting key structure  
✅ **Cache Management**: Easy to invalidate related queries

## Related Documentation

- [TanStack Query Documentation](https://tanstack.com/query/latest)
- [React Query Best Practices](https://tkdodo.eu/blog/effective-react-query-keys)
- DECM Code Standards: `.cursor/rules/code-standards.mdc`

## Questions?

If you have questions about query keys or need to add new ones, refer to this documentation or check existing patterns in the codebase.
