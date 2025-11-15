import { useSignPasswordModalStore } from "@/components/providers/SignPasswordModal/store";
import { useAuth } from "@/context/AuthContext";
import { useMutation } from "@tanstack/react-query";
import { title } from "process";
import { toast } from "sonner";

export const usePasswordPrompt = () => {
    const { user } = useAuth();

    const { open, setOnSuccess, setOnClose, setOnError } = useSignPasswordModalStore();

    const mutation = useMutation({
        mutationKey: ["sign-tx"],
        mutationFn: async ({
            eventContractAddress,
            transactionType,
            description,
            details,
        }: {
            eventContractAddress: string;
            transactionType: string;
            description: string;
            details: string;
        }) => {
            // eslint-disable-next-line no-async-promise-executor
            return new Promise(async (resolve, reject) => {
                setOnSuccess(({ value }: { value: string }) => {
                    resolve(value);
                });
                setOnError(() => {
                    reject();
                });
                setOnClose(() => {
                    // clearTimeout(timer);
                });

                if (user?.solutionStatus === "SYSTEM_MANAGED") {
                    open(title, description, true, true, {
                        contractAddress: eventContractAddress,
                        transactionType,
                        details,
                    });
                } else if (user?.solutionStatus === "BYOK") {
                    toast.error("The system doesn't support signing with self custody yet.");
                }
            });
        },
    });

    return mutation;
};
