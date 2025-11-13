import { useState } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

interface ImageLoaderProps extends React.ImgHTMLAttributes<HTMLImageElement> {
    src?: string;
    alt: string;
    skeletonClassName?: string;
}

export function ImageLoader({
    src,
    alt,
    className,
    skeletonClassName,
    ...props
}: ImageLoaderProps): React.ReactElement {
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [hasError, setHasError] = useState<boolean>(false);

    const handleLoad = (): void => {
        setIsLoading(false);
    };

    const handleError = (): void => {
        setIsLoading(false);
        setHasError(true);
    };

    return (
        <div className={cn("relative", className)}>
            {isLoading && (
                <Skeleton className={cn("absolute inset-0 w-full h-full", skeletonClassName)} />
            )}
            {!hasError && (
                <img
                    src={src}
                    alt={alt}
                    onLoad={handleLoad}
                    onError={handleError}
                    className={cn("w-full h-full object-cover", className, {
                        "opacity-0": isLoading,
                        "opacity-100 transition-opacity duration-300": !isLoading,
                    })}
                    {...props}
                />
            )}
            {hasError && (
                <div className={cn("flex items-center justify-center bg-muted", className)}>
                    <span className="text-muted-foreground text-sm">Failed to load image</span>
                </div>
            )}
        </div>
    );
}
