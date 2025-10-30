import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";

interface StyledTabsProps {
    defaultValue: string;
    children: React.ReactNode;
    className?: string;
}

interface StyledTabsTriggerProps {
    value: string;
    children: React.ReactNode;
}

interface StyledTabsContentProps {
    value: string;
    children: React.ReactNode;
    className?: string;
}

/**
 * Styled Tabs component with consistent DECM styling
 * Uses Cormorant Garamond font and primary color scheme
 */
export const StyledTabs = ({ defaultValue, children, className }: StyledTabsProps) => {
    return (
        <Tabs defaultValue={defaultValue} className={className}>
            {children}
        </Tabs>
    );
};

/**
 * Styled TabsList with consistent background and height
 */
export const StyledTabsList = ({ children }: { children: React.ReactNode }) => {
    return <TabsList className="w-full h-10 bg-[#E9DEDE]">{children}</TabsList>;
};

/**
 * Styled TabsTrigger with Cormorant Garamond font and primary color active state
 */
export const StyledTabsTrigger = ({ value, children }: StyledTabsTriggerProps) => {
    return (
        <TabsTrigger
            value={value}
            className="data-[state=active]:bg-primary text-gray-900 data-[state=active]:text-white"
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
