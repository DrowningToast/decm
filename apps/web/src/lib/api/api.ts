import { env } from "@/config/env";
import { CoreApi, HttpClient } from "@decm/api";

export const coreApiClient = new CoreApi(
    new HttpClient({
        baseURL: env.VITE_CORE_BACKEND_API,
        withCredentials: true, // Required for sending/receiving cookies in cross-origin requests
    }),
);

type CoreApiType = typeof coreApiClient;

export { CoreApi, type CoreApiType };
