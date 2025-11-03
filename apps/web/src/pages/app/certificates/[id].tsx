import { useParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { CertificateDetail } from "@/components/pages/Participant/Certificates/CertificateDetail";
import { Typography } from "@/components/typography/typography";

const CertificateDetailPage = () => {
    const { id } = useParams<{ id: string }>();
    const { t } = useTranslation();

    if (!id) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <Typography variant="text" tag="p" color="muted">
                    {t("common.error")}
                </Typography>
            </div>
        );
    }

    return (
        <section className="relative z-10 w-full">
            <CertificateDetail certificateId={id} />
        </section>
    );
};

export default CertificateDetailPage;
