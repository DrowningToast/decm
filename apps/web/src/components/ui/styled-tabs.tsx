import * as React from "react";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";
import type { ClassValue } from "clsx";

interface StyledTabsProps {
    defaultValue: string;
    children: React.ReactNode;
    className?: string;
    onValueChange?: (value: string) => void;
}

interface StyledTabsListProps {
    children: React.ReactNode;
    className?: string;
}

interface StyledTabsTriggerProps {
    value: string;
    children: React.ReactNode;
    className?: ClassValue;
}

interface StyledTabsContentProps {
    value: string;
    children?: React.ReactNode;
    className?: string;
}

/**
 * Styled Tabs component with consistent DECM styling
 * Uses Cormorant Garamond font and primary color scheme
 */
export const StyledTabs = ({
    defaultValue,
    children,
    className,
    onValueChange,
}: StyledTabsProps) => {
    return (
        <Tabs defaultValue={defaultValue} className={className} onValueChange={onValueChange}>
            {children}
        </Tabs>
    );
};

/**
 * Styled TabsList with consistent background and height
 */
export const StyledTabsList = ({ children, className }: StyledTabsListProps) => {
    return (
        <TabsList className={cn("w-full h-10 bg-[#E9DEDE] py-1.5", className)}>{children}</TabsList>
    );
};

/**
 * Styled TabsTrigger with Cormorant Garamond font and primary color active state
 */
export const StyledTabsTrigger = ({ value, children, className }: StyledTabsTriggerProps) => {
    return (
        <TabsTrigger
            value={value}
            className={cn(
                "data-[state=active]:bg-primary text-gray-900 data-[state=active]:text-white",
                className,
            )}
            style={{
                fontFamily: "Cormorant Garamond",
                fontSize: "16px",
            }}
        >
            {children}
        </TabsTrigger>
    );
};

/**
 * Styled TabsContent with consistent margin top
 */
export const StyledTabsContent = ({ value, children, className }: StyledTabsContentProps) => {
    return (
        <TabsContent value={value} className={className || "mt-6"}>
            {children}
        </TabsContent>
    );
};
