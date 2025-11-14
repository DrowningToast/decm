/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_CORE_BACKEND_API: string;
    readonly VITE_WALLETCONNECT_PROJECT_ID: string;
    readonly VITE_GOOGLE_MAPS_API_KEY: string;
    readonly DEV: boolean;
    readonly VITE_USE_MOCK_API: boolean;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
