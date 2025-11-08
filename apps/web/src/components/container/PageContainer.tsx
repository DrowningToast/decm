import type { ReactNode } from "react";
import { cva } from "class-variance-authority";
import { cn } from "@/lib/utils";

const pageContainerVariants = cva("py-10 min-h-dvh lg:py-16", {
    variants: {
        bgColor: {
            default: "bg-black",
            primary: "from-[#EB5331] to-[#362927] bg-gradient-to-b",
        },
    },
    defaultVariants: {
        bgColor: "default",
    },
});

interface PageContainerProps {
    children: ReactNode;
    title: string;
    description?: string;
    bgColor?: "default" | "primary";
    className?: string;
}

export default function PageContainer({
    children,
    // title,
    // description,
    className,
    bgColor = "default",
}: PageContainerProps) {
    const _className = cn(
        pageContainerVariants({
            bgColor,
        }),
        "relative",
        className,
    );

    const navigate = useNavigate();

    return (
        <div className={_className}>
            <div className="max-w-[1440px] mx-auto space-y-8">{children}</div>
        </div>
    );
}
