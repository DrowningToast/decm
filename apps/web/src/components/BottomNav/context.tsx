import clsx, { type ClassValue } from "clsx";
import { createContext, useContext } from "react";

interface BottomContainerContextType {
    className?: string;
    onBack?: () => void;
}

interface BottomContainerProviderProps {
    children: React.ReactNode;
    className?: ClassValue;
    onBack?: () => void;
}

export const BottomContainerContext = createContext<BottomContainerContextType | undefined>(
    undefined,
);

export const BottomContainerProvider = ({
    children,
    className,
    onBack,
}: BottomContainerProviderProps) => {
    const cn = clsx(className);

    return (
        <BottomContainerContext.Provider value={{ onBack, className: cn }}>
            {children}
        </BottomContainerContext.Provider>
    );
};

export const useBottomContainerContext = () => {
    const context = useContext(BottomContainerContext);
    if (!context) {
        throw new Error("useBottomContainerContext must be used within a BottomContainerProvider");
    }
    return context;
};
