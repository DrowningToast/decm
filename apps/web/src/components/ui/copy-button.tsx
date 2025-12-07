import { Copy, Check } from "lucide-react";
import { useState } from "react";
import { cn } from "@/lib/utils";

interface CopyButtonProps {
    text: string;
    className?: string;
    iconSize?: number;
}

export function CopyButton({ text, className, iconSize = 14 }: CopyButtonProps) {
    const [copied, setCopied] = useState(false);

    const handleCopy = async () => {
        try {
            await navigator.clipboard.writeText(text);
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
        } catch (err) {
            console.error("Failed to copy to clipboard:", err);
        }
    };

    return (
        <button
            onClick={handleCopy}
            className={cn(
                "p-1 hover:bg-foreground/10 rounded transition-colors inline-flex items-center justify-center",
                className,
            )}
            title={copied ? "Copied!" : "Copy to clipboard"}
        >
            {copied ? (
                <Check className="text-green-500" style={{ width: iconSize, height: iconSize }} />
            ) : (
                <Copy
                    className="text-muted-foreground"
                    style={{ width: iconSize, height: iconSize }}
                />
            )}
        </button>
    );
}
