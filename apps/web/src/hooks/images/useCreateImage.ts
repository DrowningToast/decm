import { useQuery } from "@tanstack/react-query";
import { QUERY_KEY } from "@/lib/queryKeys";

export function useCreateImage(url: string) {
    const {
        data: image,
        isLoading: isLoadingImage,
        isError: isErrorImage,
    } = useQuery({
        queryKey: QUERY_KEY.image.byUrl(url),
        queryFn: async () => {
            const response = await fetch(url);
            const blob = await response.blob();
            return new File([blob], url, { type: blob.type });
        },
    });

    return {
        image,
        isLoadingImage,
        isErrorImage,
    };
}
