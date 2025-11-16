import { create } from "zustand";

interface SignPasswordModalStore {
    isOpen: boolean;
    open: (
        title: string,
        description: string,
        showModeToggle: boolean,
        showSigningDetails: boolean,
        signingDetails: { contractAddress?: string; transactionType: string; details: string },
        onClose?: () => void,
    ) => void;
    close: () => void;
    onSuccess?: (result: { type: "pin" | "password"; value: string }) => void;
    onError?: (error: unknown) => void;
    onClose?: () => void;
    setOnSuccess: (callback: (result: { type: "pin" | "password"; value: string }) => void) => void;
    setOnError: (callback: (error: unknown) => void) => void;
    title: string;
    description: string;
    showModeToggle: boolean;
    showSigningDetails: boolean;
    signingDetails: {
        contractAddress?: string;
        transactionType: string;
        details: string;
    };
}

export const useSignPasswordModalStore = create<SignPasswordModalStore>((set, get) => ({
    isOpen: false,
    open: (
        title: string,
        description: string,
        showModeToggle: boolean,
        showSigningDetails: boolean,
        signingDetails: { contractAddress?: string; transactionType: string; details: string },
        onClose?: () => void,
    ) =>
        set({
            ...get(),
            isOpen: true,
            title,
            description,
            showModeToggle,
            showSigningDetails,
            signingDetails,
            onClose,
        }),
    close: () => {
        const state = get();
        state.onClose?.();
        set({ ...state, isOpen: false, onClose: undefined });
    },
    onSuccess: undefined,
    onError: undefined,
    setOnSuccess: (callback: (result: { type: "pin" | "password"; value: string }) => void) =>
        set({ ...get(), onSuccess: callback }),
    setOnError: (callback: (error: unknown) => void) => set({ ...get(), onError: callback }),
    title: "",
    description: "",
    showModeToggle: true,
    showSigningDetails: false,
    signingDetails: {
        transactionType: "",
        details: "",
    },
}));
