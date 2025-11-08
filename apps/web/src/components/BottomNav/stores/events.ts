import { create } from "zustand";

export const useSearchEventNavStore = create<{
    searchQuery: string;
    setSearchQuery: (query: string) => void;
}>((set) => ({
    searchQuery: "",
    setSearchQuery: (query: string) => set({ searchQuery: query }),
}));
