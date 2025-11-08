import { useParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { CertificateDetail } from "@/components/pages/Participant/Certificates/CertificateDetail";
import { Typography } from "@/components/typography/typography";
import { PageContainer } from "@/components/container/PageContainer";

const CertificateDetailPage = () => {
    const { id } = useParams<{ id: string }>();
    const { t } = useTranslation();

    if (!id) {
        return (
            <PageContainer className="flex items-center justify-center min-h-screen">
                <Typography variant="text" tag="p" color="muted">
                    {t("common.error")}
                </Typography>
            </PageContainer>
        );
    }

    return (
        <PageContainer className="relative z-10 w-full">
            <CertificateDetail certificateId={id} />
        </PageContainer>
    );
};

export default CertificateDetailPage;
