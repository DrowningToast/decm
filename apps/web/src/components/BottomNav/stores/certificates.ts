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
    setCertificateId: (certificateId: string | null) => void;
}

export const useCertificateDetailNavStore = create<CertificateDetailNavStore>((set) => ({
    certificateId: null,
    setCertificateId: (certificateId: string | null) => set({ certificateId }),
}));
