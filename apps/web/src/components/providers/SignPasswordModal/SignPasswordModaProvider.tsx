import { HeadlessPasswordPinModal } from "@/components/ui/headless-password-pin-modal";

type SignPasswordModalProviderProps = React.PropsWithChildren;

export const SignPasswordModalProvider: React.FC<SignPasswordModalProviderProps> = ({
    children,
}) => {
    return (
        <>
            <HeadlessPasswordPinModal />
            {children}
        </>
    );
};
