import { useSearchCertificateNavStore } from "@/components/Botto/stores";
import FuzzySearch from "fuzzy-search";

// TODO: Replace with real ceritifcate interfaces
export interface Certificate {
    id: string;
    name: string;
    issuer: string;
    issuedDate: string;
    status: "completed" | "pending";
}

// TODO: Replace with actual data from API
const mockCertificates: Certificate[] = [
    {
        id: "1",
        name: "Participation award",
        issuer: "ToBeIT69",
        issuedDate: "2025-09-24",
        status: "completed",
    },
    {
        id: "2",
        name: "Winning team award",
        issuer: "ToBeIT69",
        issuedDate: "2025-09-24",
        status: "pending",
    },
    {
        id: "3",
        name: "Participation award",
        issuer: "ToBeIT69",
        issuedDate: "2025-09-24",
        status: "pending",
    },
    {
        id: "4",
        name: "Participation award",
        issuer: "ToBeIT69",
        issuedDate: "2025-09-24",
        status: "pending",
    },
    {
        id: "5",
        name: "Participation award",
        issuer: "ToBeIT69",
        issuedDate: "2025-09-24",
        status: "pending",
    },
];

export const useCertificatesListUsecase = () => {
    const { searchQuery } = useSearchCertificateNavStore();

    const searcher = new FuzzySearch(mockCertificates, ["name", "status"], {
        caseSensitive: false,
    });

    const filteredCertificaes = searcher.search(searchQuery);

    return { certificates: filteredCertificaes };
};
