import { PageContainer } from "@/components/container/PageContainer";
import { MainPage } from "@/components/pages/Participant/MainPage";

const AppPage = () => {
    return (
        <>
            {/* Background image - positioned at bottom right, visible on both mobile and desktop */}
            <div className="absolute bottom-0 right-0 w-[424px] h-[424px] md:w-[500px] md:h-[500px] opacity-20 md:opacity-20 pointer-events-none">
                <img
                    src="/assets/scale.webp"
                    alt=""
                    className="w-full h-full object-cover object-center"
                />
            </div>
            <PageContainer className="relative z-10">
                <MainPage />
            </PageContainer>
        </>
    );
};

export default AppPage;
