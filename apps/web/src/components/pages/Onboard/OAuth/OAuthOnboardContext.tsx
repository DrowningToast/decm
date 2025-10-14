import { createContext, useEffect, useMemo, } from "react"
import { useForm, type UseFormReturn } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { ProfileSchema } from "../Profile"
import { z } from "zod"
import { useTranslation } from "react-i18next"
import { Form } from "@/components/ui/form"
import { useSignup } from "../useSignup"
import { OnboardRegistrationMethod, type OnboardCheckOnboardStatusResponse } from "@decm/api"
import { Error } from "../../Error"
import { useSearchParams } from "react-router-dom"
import { coreApiClient } from "@/lib/api/api"
import { LOCAL_STORAGE_KEYS, setLocalStorageItem } from "@/lib/constants/localStorage"
import { toast } from "sonner"
import { AxiosError } from "axios"

type OAuthOnboardContextType = {
    form: UseFormReturn<OAuthOnboardForm>
    handleSubmit: () => Promise<void>
}

const OAuthOnboardContext = createContext<OAuthOnboardContextType>(undefined as unknown as OAuthOnboardContextType)

const createOAuthOnboardFormSchema = (t: (key: string) => string) => {
    return z.object({
        password: z.string().min(6, { message: t("validation.passwordMin6") }).max(32, { message: t("validation.passwordMax32") }),
        accessToken: z.string(),
        expiresIn: z.number(),
    }).extend(ProfileSchema(t).shape)
}

type OAuthOnboardForm = z.infer<ReturnType<typeof createOAuthOnboardFormSchema>>

type OAuthOnboardErrorType = "missingAccessToken" | "missingExpiresIn"

const OAuthOnboardProvider: React.FC<React.PropsWithChildren> = ({ children }) => {
    const { t } = useTranslation()
    const { createAccount, upsertProfile, isLoading } = useSignup()
    const [searchParams] = useSearchParams()
    const accessToken = searchParams.get("access_token")
    const expiresIn = searchParams.get("expires_in")

    const OAuthOnboardFormSchema = createOAuthOnboardFormSchema(t)

    const form = useForm<OAuthOnboardForm>({
        resolver: zodResolver(OAuthOnboardFormSchema),
        mode: 'onChange'
    })

    const onSubmit = async (data: OAuthOnboardForm) => {
        if (isLoading) {
            toast.error(t("errors.generic"))
            return;
        }

        let status: OnboardCheckOnboardStatusResponse | undefined = undefined
        try {
            status = await coreApiClient.v1.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                access_token: data.accessToken,
                expires_in: data.expiresIn,
            })
        } catch (error) {
            if (error instanceof AxiosError) {
                switch (error.response?.status) {
                    case 401:
                        toast.error(t("flow.check_onboard_status.unauthenticated_response"))
                        return
                }
            }
        }

        let credential_id: string | undefined = undefined
        if (!status?.authentication_credential_id) {
            try {
                const { credential_id: _credential_id, jwt } = await createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: data.accessToken ?? "",
                    expiresIn: data.expiresIn ?? undefined,
                    password: data.password,
                })
                setLocalStorageItem(LOCAL_STORAGE_KEYS.JWT, jwt)
                credential_id = _credential_id
            } catch (error) {
                if (error instanceof AxiosError) {
                    switch (error.response?.status) {
                        case 409:
                            toast.error(t("flow.oauth_google.create_account_error_duplicate"))
                            return
                        case 401:
                            toast.error(t("flow.oauth_google.create_account_error_expired_token"))
                            return
                        case 500:
                            toast.error(t("flow.oauth_google.create_account_error_generic"))
                            return
                        default:
                            toast.error(t("flow.oauth_google.create_account_error_generic"))
                            return
                    }
                }
                toast.error(t("flow.oauth_google.create_account_error_generic"))
                return
            }
        }
        try {
            await upsertProfile({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: data.accessToken ?? "",
                expiresIn: data.expiresIn ?? undefined,
                password: data.password,
                authenticationCredentialId: credential_id ?? "",
                profile: {
                    authentication_credential_id: credential_id ?? "",
                    email: data.email,
                    phone_number: data.phoneNumber,
                    first_name: data.firstName,
                    last_name: data.lastName,
                    is_email_public: data.isEmailPublic,
                    is_phone_number_public: data.isPhoneNumberPublic,
                    is_first_name_public: data.isFirstNamePublic,
                    is_last_name_public: data.isLastNamePublic,
                }
            })
        } catch (error) {
            if (error instanceof AxiosError) {
                switch (error.response?.status) {
                    case 409:
                        toast.error(t("flow.oauth_google.create_account_error_duplicate"))
                        return
                    case 401:
                        toast.error(t("flow.oauth_google.create_account_error_expired_token"))
                        return
                    case 500:
                        toast.error(t("flow.oauth_google.create_account_error_generic"))
                        return
                    default:
                        toast.error(t("flow.oauth_google.create_account_error_generic"))
                        return
                }
            }
            toast.error(t("flow.oauth_google.create_account_error_generic"))
            return
        }

    }

    const handleSubmit = async () => {
        const data = form.getValues()
        const isValid = await form.trigger()
        if (!isValid) {
            return
        }
        await onSubmit(data)
    }

    const errorType: OAuthOnboardErrorType | undefined = useMemo(() => {
        if (form.formState.errors.accessToken) {
            return "missingAccessToken"
        }
        if (form.formState.errors.expiresIn) {
            return "missingExpiresIn"
        }
        return undefined
    }, [form.formState.errors])

    useEffect(() => {
        const init = async () => {
            if (!accessToken || !expiresIn) {
                return
            }
            form.setValue("accessToken", accessToken)
            form.setValue("expiresIn", parseInt(expiresIn ?? "0"))
        }
        init()
    }, [accessToken, form, expiresIn])

    if (errorType) {
        return <Error />
    }

    return (
        <OAuthOnboardContext.Provider value={{ form, handleSubmit }}>
            <Form {...form}>
                <form>
                    {children}
                </form>
            </Form>
        </OAuthOnboardContext.Provider>
    )
}

export {
    OAuthOnboardProvider,
    OAuthOnboardContext,
}

export type {
    OAuthOnboardForm,
}