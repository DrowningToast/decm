import { Dialog, DialogContent } from "@/components/ui/dialog";
import { useCertificateDetailNavStore } from "@/components/BottomNav/stores/certificates";

export const CertificateShareModal = () => {
    const { isShareModalOpen, setIsShareModalOpen } = useCertificateDetailNavStore();

    return (
        <Dialog open={isShareModalOpen} onOpenChange={setIsShareModalOpen}>
            <DialogContent />
        </Dialog>
    );
};
