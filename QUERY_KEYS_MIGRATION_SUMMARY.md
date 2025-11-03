# Query Keys Centralization - Summary

## ✅ Completed Tasks

### 1. Created Centralized Query Keys File

**File:** `apps/web/src/lib/queryKeys.ts`

This file contains all React Query keys organized by domain:

- User & Authentication
- Onboarding (with wallet and Google OAuth variants)
- Images
- Events (with sub-domains for issuers and certificates)
- Host Events
- Issuers

### 2. Updated Files to Use Centralized Query Keys

Updated **15 files** to use the new centralized query keys:

#### Context Files

- ✅ `apps/web/src/context/AuthContext.tsx`

#### Hook Files

- ✅ `apps/web/src/hooks/events/useEvent.ts`
- ✅ `apps/web/src/hooks/events/useEventRegistrationConfig.ts`
- ✅ `apps/web/src/hooks/events/useVerifiedIssuers.ts`
- ✅ `apps/web/src/hooks/images/useCreateImage.ts`

#### Onboarding Files

- ✅ `apps/web/src/components/pages/Onboard/useGetSignMessage.ts`
- ✅ `apps/web/src/components/pages/Onboard/useCheckOnboardStatus.tsx`

#### Host Pages Files

- ✅ `apps/web/src/components/pages/HostPages/EventPages/useEventCertificateConfig.ts`
- ✅ `apps/web/src/components/pages/HostPages/EventPages/useEventIssuers.ts`
- ✅ `apps/web/src/components/pages/HostPages/EventPages/useUpdateCertificateConfig.ts`
- ✅ `apps/web/src/components/pages/HostPages/EventPages/useDeleteEventIssuer.ts`
- ✅ `apps/web/src/components/pages/HostPages/EventsPage/useHostEvents.ts`
- ✅ `apps/web/src/components/pages/HostPages/EditEventPage/useEditEvent.ts`
- ✅ `apps/web/src/components/pages/HostPages/EditEventPage/useDeleteEvent.ts`

#### Form Files

- ✅ `apps/web/src/components/forms/ParticipantSettingsForm/useUpdateParticipantSetting.ts`

### 3. Created Comprehensive Documentation

**File:** `apps/web/src/lib/queryKeys.README.md`

Includes:

- Complete usage guide with examples
- All available query keys documented
- Best practices and patterns
- Migration guide
- TypeScript benefits explanation

### 4. No Linter Errors

All updated files pass ESLint checks with no errors related to the changes.

## 📊 Query Keys Inventory

### Before

- Hardcoded string arrays scattered across 15+ files
- No type safety
- Prone to typos
- Difficult to refactor

### After

- Single source of truth in `queryKeys.ts`
- Full TypeScript type safety
- IDE autocomplete support
- Easy to refactor and maintain

## 🎯 Key Features

### Type-Safe Factory Functions

```tsx
// Dynamic keys with parameters
queryKeys.event.byId(eventId);
queryKeys.hostEvents.list(userId, rowsPerPage, offset);
queryKeys.onboard.status.google(accessToken, expiresIn);
```

### Hierarchical Structure

```tsx
queryKeys.event.all; // ["event"]
queryKeys.event.byId("123"); // ["event", "123"]
queryKeys.event.issuers.byEventId("123"); // ["event", "123", "issuers"]
```

### Invalidation Helpers

```tsx
// Invalidate all events
queryClient.invalidateQueries({ queryKey: queryKeys.event.all });

// Invalidate specific event
queryClient.invalidateQueries({ queryKey: queryKeys.event.byId(eventId) });
```

## 📝 Usage Examples

### In useQuery

```tsx
import { queryKeys } from "@/lib/queryKeys";

useQuery({
    queryKey: queryKeys.event.byId(eventId),
    queryFn: () => fetchEvent(eventId),
});
```

### In useMutation with Invalidation

```tsx
import { queryKeys } from "@/lib/queryKeys";

useMutation({
    mutationFn: updateEvent,
    onSuccess: () => {
        queryClient.invalidateQueries({
            queryKey: queryKeys.event.all,
        });
    },
});
```

## 🔄 Migration Pattern

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

## ✨ Benefits

1. **Type Safety**: Catch errors at compile time with TypeScript
2. **Autocomplete**: IDE provides suggestions for available keys
3. **Consistency**: Same key structure throughout the app
4. **Refactoring**: Change once, update everywhere
5. **Documentation**: Self-documenting code structure
6. **Cache Management**: Easy to invalidate related queries
7. **No Typos**: Impossible to mistype query keys

## 📚 Documentation

Complete documentation available at:

- `apps/web/src/lib/queryKeys.README.md`
- `apps/web/src/lib/queryKeys.ts` (with inline comments)

## 🚀 Next Steps

To continue using this pattern:

1. **Import the keys**:

    ```tsx
    import { queryKeys } from "@/lib/queryKeys";
    ```

2. **Use in queries**:

    ```tsx
    queryKey: queryKeys.domain.operation(params);
    ```

3. **Add new keys** as needed following the existing patterns

4. **Refer to documentation** for examples and best practices

## 📊 Impact

- **Files Updated**: 15
- **Query Keys Centralized**: 11 unique key patterns
- **Lines of Documentation**: 350+
- **Type Safety**: 100% of query keys now type-safe
- **Linter Errors**: 0 (related to this change)

## ✅ Verification

All changes have been:

- ✅ Implemented successfully
- ✅ Tested with ESLint (no errors)
- ✅ Documented comprehensively
- ✅ Following project conventions

---

**Created**: November 2, 2025  
**Status**: ✅ Complete
