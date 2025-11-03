import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/queryKeys";

export function useCreateImage(url: string) {
    const {
        data: image,
        isLoading: isLoadingImage,
        isError: isErrorImage,
    } = useQuery({
        queryKey: queryKeys.image.byUrl(url),
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
