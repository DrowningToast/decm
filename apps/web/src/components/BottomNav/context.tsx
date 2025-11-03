import clsx, { type ClassValue } from "clsx";
import { createContext, useContext } from "react";

interface BottomContainerContextType {
    children?: React.ReactNode;
    className?: ClassValue;
    onBack?: () => void;
}

// eslint-disable-next-line react-refresh/only-export-components
export const BottomContainerContext = createContext<BottomContainerContextType>(
    {} as BottomContainerContextType,
);

export const BottomContainerProvider = ({
    children,
    className,
    onBack,
}: BottomContainerContextType) => {
    const cn = clsx(className);

    return (
        <BottomContainerContext.Provider value={{ onBack, className: cn }}>
            {children}
        </BottomContainerContext.Provider>
    );
};

// eslint-disable-next-line react-refresh/only-export-components
export const useBottomContainerContext = () => {
    const context = useContext(BottomContainerContext);
    if (!context) {
        throw new Error("useBottomContainerContext must be used within a BottomContainerProvider");
    }
    return context;
};
