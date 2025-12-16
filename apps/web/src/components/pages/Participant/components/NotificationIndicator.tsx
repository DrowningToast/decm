import { useMyProfile } from "@/hooks/useMyProfile";

export const NotificationIndicator: React.FC = () => {
    const { data: profile, isLoading, isError } = useMyProfile();

    if (isLoading || isError) {
        return null;
    }
    if (!profile || !profile.unreadInboxMessageCount) {
        return null;
    }

    return (
        <div className="flex items-center gap-1.5 flex-shrink-0 animate-pulse">
            <div className="w-3 h-3 bg-accent rounded-full" />
            <span className="text-sm font-semibold text-accent">
                {profile.unreadInboxMessageCount > 99 ? "99+" : profile.unreadInboxMessageCount}
            </span>
        </div>
    );
};
