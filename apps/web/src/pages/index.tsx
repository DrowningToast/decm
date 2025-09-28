import { FaviconHelmet } from '@/components/providers/helmets/FaviconHelmet';
import { Typography } from '@/components/typography/typography';

const IndexPage = () => {
    return (
        <>
            <FaviconHelmet
                title="Home | DECM - Decentralized Event Management"
                description="Welcome to DECM - Web 3.0 platform for NFT ticketing, digital credentials, and academic identity verification"
            />
            <>
                <section className='flex flex-col gap-y-2 items-center p-8 from-primary-alt to-primary bg-gradient-to-l'>
                    <Typography variant="header" tag='h1' size="header" color="primary" className='text-shadow-lg/20 text-shadow-primary'>Themis Certification</Typography>
                    <Typography variant="text" tag='p' size="base" color="muted">Welcome to the dummy index page</Typography>
                </section>
            </>
        </>
    );
};

export default IndexPage;
