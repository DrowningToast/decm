import { useState, useEffect, useRef } from "react";
import { Check, Copy } from "lucide-react";

export function CopyButton({ value }: { value: string }) {
    const [copied, setCopied] = useState(false);
    const timeoutRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);

    useEffect(
        () => () => {
            clearTimeout(timeoutRef.current);
        },
        [],
    );

    const handleCopy = () => {
        void navigator.clipboard.writeText(value).then(() => {
            setCopied(true);
            timeoutRef.current = setTimeout(() => setCopied(false), 1500);
        });
    };

    return (
        <button
            onClick={handleCopy}
            className="shrink-0 text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
            aria-label="Copy to clipboard"
        >
            {copied ? <Check className="w-3 h-3 text-green-500" /> : <Copy className="w-3 h-3" />}
        </button>
    );
}
