import { create } from "zustand";

interface EventInvitationNavStore {
    onAcceptCallback?: () => void;
    setOnAcceptCallback: (callback: () => void) => void;
}

export const useEventInvitationNavStore = create<EventInvitationNavStore>((set) => ({
    onAcceptCallback: undefined,
    setOnAcceptCallback: (callback: () => void) => set({ onAcceptCallback: callback }),
}));
