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
    setCertificateId: (certificateId: string | null) => void;
    setIsClaimed: (isClaimed: boolean) => void;
}

export const useCertificateDetailNavStore = create<CertificateDetailNavStore>((set) => ({
    certificateId: null,
    isClaimed: false,
    setCertificateId: (certificateId: string | null) => set({ certificateId }),
    setIsClaimed: (isClaimed: boolean) => set({ isClaimed }),
}));
