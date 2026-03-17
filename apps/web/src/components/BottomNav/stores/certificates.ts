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
    onClickShareable: () => void;
    setOnClickShareable: (onClickShareable: () => void) => void;
    isShareableLoading: boolean;
}

export const useCertificateDetailNavStore = create<CertificateDetailNavStore>((set) => ({
    certificateId: null,
    isClaimed: false,
    imageUrl: null,
    setCertificateId: (certificateId: string | null) => set({ certificateId }),
    setIsClaimed: (isClaimed: boolean) => set({ isClaimed }),
    setImageUrl: (imageUrl: string | null) => set({ imageUrl }),
    isShareableLoading: false,
    // callback
    onClickShareable: () => {},
    setOnClickShareable: (onClickShareable: () => void) =>
        set({
            onClickShareable: () => {
                set({ isShareableLoading: true });
                onClickShareable();
                set({ isShareableLoading: false });
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
    setShareableHandle: (handle: string | null) => void;
    setIsPasswordDialogOpen: (open: boolean) => void;
    // callback
    onClickPassword: (isPasswordProtected: boolean) => void;
    setOnClickPassword: (onPasswordProtected: () => void) => void;
    isPasswordLoading: boolean;
    onChangePublish: (isPublished: boolean) => void;
    setOnChangePublish: (onPublishStatus: (isPublished: boolean) => void) => void;
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
        setShareableHandle: (handle: string | null) => set({ shareableHandle: handle }),
        setIsPasswordDialogOpen: (open: boolean) => set({ isPasswordDialogOpen: open }),
        onClickPassword: () => {},
        isPasswordLoading: false,
        onChangePublish: () => {},
        isPublishLoading: false,
        setOnClickPassword: (cn: () => void) =>
            set({
                onClickPassword: () => {
                    set({ isPasswordLoading: true });
                    cn();
                    set({ isPasswordLoading: false });
                },
            }),
        setOnChangePublish: (cn: (isPublished: boolean) => void) =>
            set({
                onChangePublish: (isPublished: boolean) => {
                    set({ isPublishLoading: true });
                    cn(isPublished);
                    set({ isPublishLoading: false });
                },
            }),
    }),
);
