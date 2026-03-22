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
    onClickShareable: () => Promise<void>;
    setOnClickShareable: (onClickShareable: () => Promise<void>) => void;
    isShareableLoading: boolean;
    isShareModalOpen: boolean;
    setIsShareModalOpen: (open: boolean) => void;
}

export const useCertificateDetailNavStore = create<CertificateDetailNavStore>((set) => ({
    certificateId: null,
    isClaimed: false,
    imageUrl: null,
    setCertificateId: (certificateId: string | null) => set({ certificateId }),
    setIsClaimed: (isClaimed: boolean) => set({ isClaimed }),
    setImageUrl: (imageUrl: string | null) => set({ imageUrl }),
    isShareableLoading: false,
    isShareModalOpen: false,
    setIsShareModalOpen: (open: boolean) => set({ isShareModalOpen: open }),
    // callback
    onClickShareable: async () => {},
    setOnClickShareable: (onClickShareable: () => Promise<void>) =>
        set({
            onClickShareable: async () => {
                set({ isShareableLoading: true });
                try {
                    await onClickShareable();
                } finally {
                    set({ isShareableLoading: false });
                }
            },
        }),
}));

interface CertificateDetailsSharedNavStore {
    shareableUrl: string | null;
    shareableHandle: string | null;
    isPasswordProtected: boolean;
    isPublished: boolean;
    isPasswordDialogOpen: boolean;
    setIsPublished: (isPublished: boolean) => void;
    setIsPasswordProtected: (isPasswordProtected: boolean) => void;
    setShareableHandle: (handle: string | null) => void;
    setIsPasswordDialogOpen: (open: boolean) => void;
    // callback
    onClickPassword: (isPasswordProtected: boolean) => Promise<void>;
    setOnClickPassword: (onPasswordProtected: () => Promise<void>) => void;
    isPasswordLoading: boolean;
    onChangePublish: (isPublished: boolean) => Promise<void>;
    setOnChangePublish: (onPublishStatus: (isPublished: boolean) => Promise<void>) => void;
    isPublishLoading: boolean;
}

export const useCertificateDetailsSharedNavStore = create<CertificateDetailsSharedNavStore>(
    (set) => ({
        isPasswordProtected: false,
        isPublished: false,
        isPasswordDialogOpen: false,
        shareableUrl: null,
        shareableHandle: null,
        setIsPublished: (isPublished: boolean) => set({ isPublished }),
        setIsPasswordProtected: (isPasswordProtected: boolean) => set({ isPasswordProtected }),
        setShareableHandle: (handle: string | null) => set({ shareableHandle: handle }),
        setIsPasswordDialogOpen: (open: boolean) => set({ isPasswordDialogOpen: open }),
        onClickPassword: async () => {},
        isPasswordLoading: false,
        onChangePublish: async () => {},
        isPublishLoading: false,
        setOnClickPassword: (cn: () => Promise<void>) =>
            set({
                onClickPassword: async () => {
                    set({ isPasswordLoading: true });
                    try {
                        await cn();
                    } finally {
                        set({ isPasswordLoading: false });
                    }
                },
            }),
        setOnChangePublish: (cn: (isPublished: boolean) => Promise<void>) =>
            set({
                onChangePublish: async (isPublished: boolean): Promise<void> => {
                    set({ isPublishLoading: true });
                    try {
                        await cn(isPublished);
                    } finally {
                        set({ isPublishLoading: false });
                    }
                },
            }),
    }),
);
