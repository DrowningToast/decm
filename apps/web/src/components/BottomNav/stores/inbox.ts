import { create } from "zustand";

interface InboxNavStore {
    onViewEventCallback?: () => void;
    onViewCertificateCallback?: () => void;
    setOnViewEventCallback: (callback: () => void) => void;
    setOnViewCertificateCallback: (callback: () => void) => void;
    resetCallbacks: () => void;
}

export const useInboxNavStore = create<InboxNavStore>((set) => ({
    onViewEventCallback: undefined,
    onViewCertificateCallback: undefined,
    setOnViewEventCallback: (callback) => set({ onViewEventCallback: callback }),
    setOnViewCertificateCallback: (callback) => set({ onViewCertificateCallback: callback }),
    resetCallbacks: () =>
        set({
            onViewEventCallback: undefined,
            onViewCertificateCallback: undefined,
        }),
}));
