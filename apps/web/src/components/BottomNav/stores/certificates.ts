import { create } from "zustand";

interface SearchCertificateNavStore {
    searchQuery: string;
    setSearchQuery: (searchQuery: string) => void;
}

export const useSearchCertificateNavStore = create<SearchCertificateNavStore>((set) => ({
    searchQuery: "",
    setSearchQuery: (searchQuery: string) => set({ searchQuery }),
}));

interface CertificateDetailNavStore {
    certificateId: string | null;
    isClaimed: boolean;
    imageUrl: string | null;
    setCertificateId: (certificateId: string | null) => void;
    setIsClaimed: (isClaimed: boolean) => void;
    setImageUrl: (imageUrl: string | null) => void;
}

export const useCertificateDetailNavStore = create<CertificateDetailNavStore>((set) => ({
    certificateId: null,
    isClaimed: false,
    imageUrl: null,
    setCertificateId: (certificateId: string | null) => set({ certificateId }),
    setIsClaimed: (isClaimed: boolean) => set({ isClaimed }),
    setImageUrl: (imageUrl: string | null) => set({ imageUrl }),
}));
