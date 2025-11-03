import { create } from "zustand";

interface SearchIdentitiesNavStore {
    searchQuery: string;
    setSearchQuery: (query: string) => void;
}

export const useSearchIdentitiesNavStore = create<SearchIdentitiesNavStore>((set) => ({
    searchQuery: "",
    setSearchQuery: (query: string) => set({ searchQuery: query }),
}));
