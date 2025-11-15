import { PasswordPinModal } from "@/components/ui/password-pin-modal";
import { useSignPasswordModalStore } from "./store";

type SignPasswordModalProviderProps = React.PropsWithChildren;

export const SignPasswordModalProvider: React.FC<SignPasswordModalProviderProps> = ({
    children,
}) => {
    const { isOpen, ...store } = useSignPasswordModalStore();

    return (
        <>
            {isOpen ? <PasswordPinModal isOpen={isOpen} {...store} /> : null}
            {children}
        </>
    );
};
