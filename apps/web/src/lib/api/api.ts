import { env } from "@/config/env";
import { CoreApi, HttpClient } from "@decm/api";

export const coreApi = new CoreApi(
	new HttpClient({
		baseURL: env.VITE_CORE_BACKEND_API,
	})
);
