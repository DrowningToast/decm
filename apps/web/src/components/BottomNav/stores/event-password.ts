import { create } from "zustand";

interface EventPasswordNavStore {
    password: string;
    setPassword: (password: string) => void;
    resetPassword: () => void;
    onSubmitCallback?: (password: string) => void;
    setOnSubmitCallback: (callback: (password: string) => void) => void;
}

export const useEventPasswordNavStore = create<EventPasswordNavStore>((set) => ({
    password: "",
    setPassword: (password: string) => set({ password }),
    resetPassword: () => set({ password: "" }),
    onSubmitCallback: undefined,
    setOnSubmitCallback: (callback: (password: string) => void) =>
        set({ onSubmitCallback: callback }),
}));
