import { createContext } from "react"
import { useForm, type UseFormReturn } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { ProfileSchema } from "../Profile"
import { z } from "zod"

type OAuthOnboardContextType = {
    form: UseFormReturn<OAuthOnboardForm>
}

const OAuthOnboardContext = createContext<OAuthOnboardContextType>(undefined as unknown as OAuthOnboardContextType)

const OAuthOnboardFormSchema = z.object({
    password: z.string().min(6),
}).extend(ProfileSchema.shape)

type OAuthOnboardForm = z.infer<typeof OAuthOnboardFormSchema>

const OAuthOnboardProvider: React.FC<React.PropsWithChildren> = ({ children }) => {

    const form = useForm<OAuthOnboardForm>({
        resolver: zodResolver(OAuthOnboardFormSchema),
    })

    return (
        <OAuthOnboardContext.Provider value={{ form }}>
            {children}
        </OAuthOnboardContext.Provider>
    )
}

export {
    OAuthOnboardProvider,
    OAuthOnboardContext,
}

export type {
    OAuthOnboardForm,
    OAuthOnboardFormSchema,
}