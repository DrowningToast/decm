import { FaviconHelmet } from '@/components/providers/helmets/FaviconHelmet';
import { LandingPage } from '@/components/pages/LandingPage/LandingPage';

const IndexPage = () => {
    return (
        <>
            <FaviconHelmet
                title="Home | DECM - Decentralized Event Management"
                description="Welcome to DECM - Web 3.0 platform for NFT ticketing, digital credentials, and academic identity verification"
            />
            <LandingPage />
        </>
    );
};

export default IndexPage;
