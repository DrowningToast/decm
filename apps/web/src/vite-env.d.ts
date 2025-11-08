/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_CORE_BACKEND_API: string;
    readonly VITE_WALLETCONNECT_PROJECT_ID: string;
    readonly DEV: boolean;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
