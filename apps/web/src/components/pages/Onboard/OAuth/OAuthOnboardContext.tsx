import { createContext, useEffect, useMemo, } from "react"
import { useForm, type UseFormReturn } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { ProfileSchema } from "../Profile"
import { z } from "zod"
import { useTranslation } from "react-i18next"
import { Form } from "@/components/ui/form"
import { useSignup } from "../useSignup"
import { OnboardRegistrationMethod } from "@decm/api"
import { Error } from "../../Error"
import { useSearchParams } from "react-router-dom"

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
    const { signup } = useSignup()
    const [searchParams] = useSearchParams()
    const accessToken = searchParams.get("access_token")
    const expiresIn = searchParams.get("expires_in")

    const OAuthOnboardFormSchema = createOAuthOnboardFormSchema(t)

    const form = useForm<OAuthOnboardForm>({
        resolver: zodResolver(OAuthOnboardFormSchema),
        mode: 'onChange'
    })

    const onSubmit = async (data: OAuthOnboardForm) => {
        console.log(data)
        await signup({
            method: OnboardRegistrationMethod.RegistrationMethodGoogle,
            accessToken: data.accessToken ?? "",
            expiresIn: data.expiresIn ?? undefined,
            profile: data,
            password: data.password,
        })
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