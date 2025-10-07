import { Typography } from '@/components/typography/typography';
import { PublicNavbar } from '@/components/layouts/navigations/PublicNavbar';
import { Link } from '@/router';
import { useTranslation } from 'react-i18next';

export const LandingPage = () => {
    const { t } = useTranslation();

    return (
        <div className="min-h-screen bg-background">
            <PublicNavbar />

            {/* Hero Section */}
            <section className="relative min-h-[812px] md:min-h-[982px] overflow-hidden">
                {/* Gradient Background */}
                <div className="absolute inset-0 bg-gradient-to-t from-[#362927] via-[#8b3a28] to-[#eb5331]" />

                {/* Lady Justice Image */}
                <div className="absolute left-1/2 -translate-x-1/2 top-12 md:top-[470px] w-[428px] md:w-[674px] h-[769px] md:h-[1211px] opacity-60">
                    <img
                        src="/lady-justice.png"
                        alt="Lady Justice"
                        className="w-full h-full object-cover"
                    />
                </div>

                {/* Hero Content */}
                <div className="relative z-10 flex flex-col items-end md:items-center justify-end md:justify-start px-6 md:px-0 pt-[360px] md:pt-[167px] pb-16 md:pb-0 min-h-[812px] md:min-h-[982px]">
                    <div className="flex flex-col items-end md:items-center gap-12 md:gap-3 max-w-[638px] w-full">
                        {/* Title and Subtitle */}
                        <div className="flex flex-col items-end md:items-center gap-3 md:gap-6">
                            <Typography
                                variant="header"
                                tag="h1"
                                className="text-[36px] md:text-[64px] leading-[40px] md:leading-[56px] text-foreground text-right md:text-center font-['Cormorant_Garamond'] font-bold tracking-[0.06px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                            >
                                {t('landing.hero.title')}
                            </Typography>

                            <Typography
                                variant="text"
                                tag="p"
                                className="text-base md:text-2xl text-[#e9dede] text-right md:text-center max-w-[257px] md:max-w-[638px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {t('landing.hero.subtitle')}
                            </Typography>
                        </div>

                        {/* CTA Buttons */}
                        <div className="flex flex-col items-center gap-1 w-full md:w-auto">
                            <Link to="/signup">
                                <button className="bg-[#e9dede] hover:bg-[#fcfcfc] transition-colors rounded-xl px-8 py-3 text-primary font-medium text-base [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] w-[325px]">
                                    {t('landing.hero.ctaButton')}
                                </button>
                            </Link>
                            <Link to="/signup">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    className="text-[#e9dede] text-xs italic underline [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] hover:text-foreground transition-colors"
                                >
                                    {t('landing.hero.hasAccount')}
                                </Typography>
                            </Link>
                        </div>
                    </div>

                    {/* Scroll Indicator - Mobile Only */}
                    <Typography
                        variant="text"
                        tag="p"
                        className="md:hidden absolute bottom-4 right-6 text-[#b8b8b8] text-base [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {t('landing.hero.scrollIndicator')}
                    </Typography>
                </div>
            </section>

            {/* About & Features Section */}
            <section id="features" className="relative bg-gradient-to-b from-[#362927] to-[#0f1012] px-6 py-12 md:px-16 md:py-16">
                {/* About Themis */}
                <div className="max-w-[1384px] mx-auto mb-12 md:mb-24">
                    <div className="flex flex-col gap-3 md:gap-6">
                        <Typography
                            variant="header"
                            tag="h2"
                            className="text-[36px] md:text-[64px] leading-[40px] text-foreground font-['Bodoni_Moda'] font-semibold italic tracking-[0.06px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                        >
                            {t('landing.about.title')}
                        </Typography>

                        <Typography
                            variant="text"
                            tag="p"
                            className="text-base md:text-2xl text-foreground max-w-full md:max-w-[1384px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                        >
                            {t('landing.about.description')}
                        </Typography>
                    </div>
                </div>

                {/* Feature Cards */}
                <div className="max-w-[1384px] mx-auto mb-12 md:mb-24">
                    <div className="flex flex-col md:flex-row gap-6 md:gap-6">
                        {/* Card 1: Wallet Management */}
                        <div className="flex md:flex-col gap-3 md:gap-2.5 md:bg-foreground md:rounded-xl md:p-6 flex-1">
                            <img
                                src="/wallet-icon.svg"
                                alt="Wallet"
                                className="size-6 md:size-8 shrink-0"
                            />
                            <div className="flex flex-col gap-2">
                                <Typography
                                    variant="text"
                                    tag="h3"
                                    className="text-base md:text-2xl text-foreground md:text-primary font-normal md:font-semibold [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] md:[text-shadow:none]"
                                >
                                    {t('landing.features.wallet.title')}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-base text-[#b8b8b8] md:text-[#0f1012] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] md:[text-shadow:none]"
                                >
                                    {t('landing.features.wallet.description')}
                                </Typography>
                            </div>
                        </div>

                        {/* Card 2: Data Verification */}
                        <div className="flex md:flex-col gap-3 md:gap-2.5 md:bg-foreground md:rounded-xl md:p-6 flex-1">
                            <img
                                src="/receipt-text-icon.svg"
                                alt="Receipt"
                                className="size-6 md:size-8 shrink-0"
                            />
                            <div className="flex flex-col gap-2">
                                <Typography
                                    variant="text"
                                    tag="h3"
                                    className="text-base md:text-2xl text-foreground md:text-primary font-normal md:font-semibold [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] md:[text-shadow:none]"
                                >
                                    {t('landing.features.verification.title')}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-base text-[#b8b8b8] md:text-[#0f1012] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] md:[text-shadow:none]"
                                >
                                    {t('landing.features.verification.description')}
                                </Typography>
                            </div>
                        </div>

                        {/* Card 3: Download Certificates */}
                        <div className="flex md:flex-col gap-3 md:gap-2.5 md:bg-foreground md:rounded-xl md:p-6 flex-1">
                            <img
                                src="/shield-check-icon.svg"
                                alt="Shield"
                                className="size-6 md:size-8 shrink-0"
                            />
                            <div className="flex flex-col gap-2">
                                <Typography
                                    variant="text"
                                    tag="h3"
                                    className="text-base md:text-2xl text-foreground md:text-primary font-normal md:font-semibold [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] md:[text-shadow:none]"
                                >
                                    {t('landing.features.download.title')}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-base text-[#b8b8b8] md:text-[#0f1012] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] md:[text-shadow:none]"
                                >
                                    {t('landing.features.download.description')}
                                </Typography>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Get Started Section */}
                <div id="get-started" className="max-w-[1384px] mx-auto">
                    <div className="flex flex-col gap-6 md:gap-12">
                        <Typography
                            variant="header"
                            tag="h2"
                            className="text-base md:text-base text-foreground font-['Cormorant_Garamond'] font-semibold tracking-[0.06px]"
                        >
                            {t('landing.getStarted.title')}
                        </Typography>

                        {/* Mobile Links */}
                        <div className="md:hidden flex flex-col gap-6">
                            <ul className="list-disc list-inside space-y-1">
                                <li>
                                    <Typography
                                        variant="text"
                                        tag="span"
                                        className="text-[#b8b8b8] text-base italic underline [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                    >
                                        {t('landing.getStarted.mobile.verificationPortal')}
                                    </Typography>
                                </li>
                                <li>
                                    <Typography
                                        variant="text"
                                        tag="span"
                                        className="text-[#b8b8b8] text-base italic underline [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                    >
                                        {t('landing.getStarted.mobile.myCertificate')}
                                    </Typography>
                                </li>
                                <li>
                                    <Typography
                                        variant="text"
                                        tag="span"
                                        className="text-[#b8b8b8] text-base italic underline [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                    >
                                        {t('landing.getStarted.mobile.viewServices')}
                                    </Typography>
                                </li>
                            </ul>
                        </div>

                        {/* Desktop Links Grid */}
                        <div className="hidden md:grid md:grid-cols-5 gap-8">
                            {/* Public */}
                            <div className="flex flex-col gap-3">
                                <Typography variant="text" tag="h3" className="text-foreground text-base">
                                    {t('landing.getStarted.public.title')}
                                </Typography>
                                <ul className="list-disc list-inside space-y-1">
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.public.event')}
                                        </Typography>
                                    </li>
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.public.verifyCertificates')}
                                        </Typography>
                                    </li>
                                </ul>
                            </div>

                            {/* Students */}
                            <div className="flex flex-col gap-3">
                                <Typography variant="text" tag="h3" className="text-foreground text-base">
                                    {t('landing.getStarted.students.title')}
                                </Typography>
                                <ul className="list-disc list-inside space-y-1">
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.students.signIn')}
                                        </Typography>
                                    </li>
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.students.signUp')}
                                        </Typography>
                                    </li>
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.students.certificateManagement')}
                                        </Typography>
                                    </li>
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.students.myEvents')}
                                        </Typography>
                                    </li>
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.students.inbox')}
                                        </Typography>
                                    </li>
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.students.identities')}
                                        </Typography>
                                    </li>
                                </ul>
                            </div>

                            {/* Organizer */}
                            <div className="flex flex-col gap-3">
                                <Typography variant="text" tag="h3" className="text-foreground text-base">
                                    {t('landing.getStarted.organizer.title')}
                                </Typography>
                                <ul className="list-disc list-inside space-y-1">
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.organizer.register')}
                                        </Typography>
                                    </li>
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.organizer.signIn')}
                                        </Typography>
                                    </li>
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.organizer.eventManagement')}
                                        </Typography>
                                    </li>
                                </ul>
                            </div>

                            {/* Issuer */}
                            <div className="flex flex-col gap-3">
                                <Typography variant="text" tag="h3" className="text-foreground text-base">
                                    {t('landing.getStarted.issuer.title')}
                                </Typography>
                                <ul className="list-disc list-inside space-y-1">
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.issuer.signInIAM')}
                                        </Typography>
                                    </li>
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.issuer.register')}
                                        </Typography>
                                    </li>
                                </ul>
                            </div>

                            {/* Developer */}
                            <div className="flex flex-col gap-3">
                                <Typography variant="text" tag="h3" className="text-foreground text-base">
                                    {t('landing.getStarted.developer.title')}
                                </Typography>
                                <ul className="list-disc list-inside space-y-1">
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.developer.manageIntegration')}
                                        </Typography>
                                    </li>
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.developer.apiDocuments')}
                                        </Typography>
                                    </li>
                                    <li>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-[#b8b8b8] text-base italic underline cursor-pointer hover:text-foreground transition-colors [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                        >
                                            {t('landing.getStarted.developer.readDocs')}
                                        </Typography>
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        </div>
    );
};