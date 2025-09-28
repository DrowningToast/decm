import { CustomHelmet } from "@/components/providers/helmets/CustomHelmet"
import { useHelmet } from "@/hooks/useHelmet"

export const LandingPage = () => {

    const helmetData = useHelmet({
        pageType: 'home',
        title: 'Landing | DECM - Decentralized Event Management',
        description: 'Landing page for DECM - Decentralized Event Management'
    })

    return (
        <>
            <CustomHelmet
                title={helmetData.title}
                description={helmetData.description}
                themeColor={helmetData.themeColor}
            />
            <div className="flex flex-col items-center">

            </div>
        </>
    )
}