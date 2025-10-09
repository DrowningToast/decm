import { createContext } from "react"
import { useForm, type UseFormReturn } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { ProfileSchema } from "../Profile"
import { z } from "zod"
import { useTranslation } from "react-i18next"
import { Form } from "@/components/ui/form"

type OAuthOnboardContextType = {
    form: UseFormReturn<OAuthOnboardForm>
    onSubmit: (data: OAuthOnboardForm) => void
}

const OAuthOnboardContext = createContext<OAuthOnboardContextType>(undefined as unknown as OAuthOnboardContextType)

const createOAuthOnboardFormSchema = (t: (key: string) => string) => {
    return z.object({
        password: z.string().min(6, { message: t("validation.passwordMin6") }).max(32, { message: t("validation.passwordMax32") }),

    }).extend(ProfileSchema(t).shape)
}

type OAuthOnboardForm = z.infer<ReturnType<typeof createOAuthOnboardFormSchema>>

const OAuthOnboardProvider: React.FC<React.PropsWithChildren> = ({ children }) => {
    const { t } = useTranslation()

    const OAuthOnboardFormSchema = createOAuthOnboardFormSchema(t)

    const form = useForm<OAuthOnboardForm>({
        resolver: zodResolver(OAuthOnboardFormSchema),
        mode: 'onChange'
    })

    const onSubmit = (data: OAuthOnboardForm) => {
        // TODO: Implement submit logic
        alert("submit")
        console.log(data)
    }

    return (
        <OAuthOnboardContext.Provider value={{ form, onSubmit }}>
            <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)}>
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