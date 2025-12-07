import { useSignPasswordModalStore } from "@/components/providers/SignPasswordModal/store";
import { useAuth } from "@/context/AuthContext";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

export const usePasswordPrompt = () => {
    const { user } = useAuth();

    const { open, setOnSuccess, setOnError } = useSignPasswordModalStore();

    const mutation = useMutation({
        mutationKey: ["sign-tx"],
        mutationFn: async ({
            eventContractAddress,
            transactionType,
            title,
            description,
            details,
            onClose: _onClose,
        }: {
            eventContractAddress: string;
            transactionType: string;
            title: string;
            description: string;
            details: string;
            onClose?: () => void;
        }) => {
            // eslint-disable-next-line no-async-promise-executor
            return new Promise<string>(async (resolve, reject) => {
                setOnSuccess(({ value }: { value: string }) => {
                    resolve(value);
                });
                setOnError(() => {
                    reject();
                });

                if (user?.solutionStatus === "SYSTEM_MANAGED") {
                    open(
                        title,
                        description,
                        true,
                        true,
                        {
                            contractAddress: eventContractAddress,
                            transactionType,
                            details,
                        },
                        () => {
                            console.log("helo");
                            reject();
                            _onClose?.();
                        },
                    );
                } else if (user?.solutionStatus === "BYOK") {
                    toast.error("The system doesn't support signing with self custody yet.");
                }
            });
        },
    });

    return mutation;
};
