import { CoreApi, HttpClient } from "@decm/api";

const httpClient = new HttpClient({});
const coreApi = new CoreApi(httpClient);

export { coreApi };
