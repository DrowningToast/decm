import { useQuery } from "@tanstack/react-query";

export function useCreateImage(url: string) {
    const {
        data: image,
        isLoading: isLoadingImage,
        isError: isErrorImage,
    } = useQuery({
        queryKey: ["image", url],
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
