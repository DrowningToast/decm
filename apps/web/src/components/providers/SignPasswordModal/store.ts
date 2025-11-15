import { create } from "zustand";

interface SignPasswordModalStore {
    isOpen: boolean;
    open: (
        title: string,
        description: string,
        showModeToggle: boolean,
        showSigningDetails: boolean,
        signingDetails: { contractAddress?: string; transactionType: string; details: string },
    ) => void;
    onClose: () => void;
    onSuccess?: (result: { type: "pin" | "password"; value: string }) => void;
    onError?: (error: unknown) => void;
    setOnSuccess: (callback: (result: { type: "pin" | "password"; value: string }) => void) => void;
    setOnError: (callback: (error: unknown) => void) => void;
    setOnClose: (callback: () => void) => void;
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

export const useSignPasswordModalStore = create<SignPasswordModalStore>((set) => ({
    isOpen: false,
    open: (
        title: string,
        description: string,
        showModeToggle: boolean,
        showSigningDetails: boolean,
        signingDetails: { contractAddress?: string; transactionType: string; details: string },
    ) =>
        set({
            isOpen: true,
            title,
            description,
            showModeToggle,
            showSigningDetails,
            signingDetails,
        }),
    onClose: () => set({ isOpen: false }),
    onSuccess: undefined,
    onError: undefined,
    setOnSuccess: (callback: (result: { type: "pin" | "password"; value: string }) => void) =>
        set({ onSuccess: callback }),
    setOnError: (callback: (error: unknown) => void) => set({ onError: callback }),
    setOnClose: (callback: () => void) => set({ onClose: callback }),
    title: "",
    description: "",
    showModeToggle: true,
    showSigningDetails: false,
    signingDetails: {
        transactionType: "",
        details: "",
    },
}));
