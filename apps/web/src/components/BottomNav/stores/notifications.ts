import { create } from "zustand";

interface SearchNotificationNavStore {
    searchQuery: string;
    setSearchQuery: (query: string) => void;
}

export const useSearchNotificationNavStore = create<SearchNotificationNavStore>((set) => ({
    searchQuery: "",
    setSearchQuery: (query: string) => set({ searchQuery: query }),
}));
