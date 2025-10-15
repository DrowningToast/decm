import { createContext, useContext, useEffect, useMemo, } from "react"
import { useForm, type UseFormReturn } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { ProfileSchema } from "../Profile"
import { z } from "zod"
import { useTranslation } from "react-i18next"
import { Form } from "@/components/ui/form"
import { useSignup } from "../useSignup"
import { OnboardRegistrationMethod, type OnboardCheckOnboardStatusResponse } from "@decm/api"
import { ErrorPage } from "../../Error"
import { coreApiClient } from "@/lib/api/api"
import { toast } from "sonner"
import { handleUniversalError } from "@/common/Err"
import { USECASE_IDS } from "@/constants/usecase"
import { OnboardPageContext } from "../../../../pages/onboard/[method]"

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
    const { accessToken, expiresIn, onboardStatus } = useContext(OnboardPageContext)
    const { createAccount, upsertProfile, isLoading } = useSignup()

    const OAuthOnboardFormSchema = createOAuthOnboardFormSchema(t)

    const form = useForm<OAuthOnboardForm>({
        resolver: zodResolver(OAuthOnboardFormSchema),
        mode: 'onChange'
    })

    const onSubmit = async (data: OAuthOnboardForm) => {
        console.log(isLoading);
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
            if (error instanceof Error) {
                handleUniversalError(t, error)
            }
        }

        let credential_id: string | undefined = undefined
        if (!status?.authentication_credential_id) {
            try {
                const { credential_id: _credential_id, } = await createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: data.accessToken ?? "",
                    expiresIn: data.expiresIn ?? undefined,
                    password: data.password,
                })
                credential_id = _credential_id
            } catch (error) {
                if (error instanceof Error) {
                    handleUniversalError(t, error, {
                        onDuplicateEntry: () => {
                            toast.error(t("flow.oauth_google.create_account_error_duplicate"))
                            return
                        },
                        onUnauthorized: () => {
                            toast.error(t("flow.oauth_google.create_account_error_expired_token"))
                            return
                        },
                        onInternalServerError: () => {
                            toast.error(t("flow.oauth_google.create_account_error_generic"))
                            return
                        }
                    }, USECASE_IDS.OAUTH_GOOGLE_CREATE_ACCOUNT);
                    return
                }
            }
        }

        try {
            console.log(credential_id);
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
            if (error instanceof Error) {
                handleUniversalError(t, error, {
                    onInternalServerError: () => {
                        toast.error(t("flow.oauth_google.create_account_error_generic"))
                        return
                    },
                    onDuplicateEntry: () => {
                        toast.error(t("flow.oauth_google.create_account_error_duplicate"))
                        return
                    },
                }, USECASE_IDS.OAUTH_GOOGLE_CREATE_PROFILE);
            }
        }

    }

    // TODOL Determine is valid
    const handleSubmit = async () => {
        const data = form.getValues()
        if (!data.accessToken || !data.expiresIn) {
            return
        }
        if (!onboardStatus?.authentication_credential_id) {
            if (!data.password) {
                toast.error(t("flow.oauth_google.create_account_error_password"))
                return
            }
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
    console.log(form.formState.errors)

    useEffect(() => {
        const init = async () => {
            if (!accessToken || !expiresIn) {
                return
            }
            form.setValue("accessToken", accessToken)
            form.setValue("expiresIn", expiresIn)
        }
        init()
    }, [accessToken, form, expiresIn])

    if (errorType) {
        return <ErrorPage />
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