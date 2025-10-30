import { Button } from "@/components/ui/button";
import { PublicNavbar } from "@/components/layouts/navigations/PublicNavbar";
import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";

/**
 * Component Examples - Usage Guide
 *
 * This file demonstrates how to use the Button and Navbar components
 * with the Figma design system variants.
 */

export function ButtonExamples() {
    const { t } = useTranslation();
    return (
        <div className="space-y-8 p-8">
            <div className="space-y-4">
                <h2 className="text-2xl font-bold">Figma Button Variants</h2>

                {/* Primary Button - Web3 Wallet */}
                <div className="space-y-2">
                    <p className="text-sm text-muted-foreground">Primary Variant (Web3/Wallet)</p>
                    <Button variant="primary" size="xl">
                        Sign in with Web3 Wallet Provider
                    </Button>
                </div>

                {/* Secondary Dark Button - Google */}
                <div className="space-y-2">
                    <p className="text-sm text-muted-foreground">Secondary Dark Variant (Google)</p>
                    <Button variant="secondary-dark" size="xl">
                        Sign up with Google Account
                    </Button>
                </div>

                {/* Secondary Light Button - Portfolio */}
                <div className="space-y-2">
                    <p className="text-sm text-muted-foreground">
                        Secondary Light Variant (Portfolio)
                    </p>
                    <Button variant="secondary-light" size="xl">
                        Start building your portfolio
                    </Button>
                </div>
            </div>

            <div className="space-y-4">
                <Typography variant="header" tag="h2">
                    {t("examples.buttonSizes")}
                </Typography>

                <div className="flex flex-wrap items-center gap-4">
                    <Button variant="primary" size="sm">
                        {t("examples.buttons.small")}
                    </Button>
                    <Button variant="primary" size="default">
                        {t("examples.buttons.default")}
                    </Button>
                    <Button variant="primary" size="lg">
                        {t("examples.buttons.large")}
                    </Button>
                    <Button variant="primary" size="xl">
                        {t("examples.buttons.extraLarge")}
                    </Button>
                </div>
            </div>
        </div>
    );
}

export function NavbarExamples() {
    return (
        <div className="space-y-0">
            <h2 className="text-2xl font-bold p-8">Navbar Variants</h2>

            {/* Primary Navbar */}
            <div className="relative h-20 mb-4">
                <p className="text-sm text-muted-foreground px-8 pb-2">Primary Light Variant</p>
                <PublicNavbar variant="light" className="relative" />
            </div>

            {/* Secondary Dark Navbar */}
            <div className="relative h-20 mb-4">
                <p className="text-sm text-muted-foreground px-8 pb-2">Secondary Dark Variant</p>
                <PublicNavbar variant="dark" className="relative" />
            </div>

            {/* Secondary Light Navbar */}
            <div className="relative h-20 mb-4">
                <p className="text-sm text-muted-foreground px-8 pb-2">Secondary Light Variant</p>
                <PublicNavbar variant="light" className="relative" />
            </div>
        </div>
    );
}

/**
 * Usage in Your Pages:
 *
 * // Import the components
 * import { Button } from '@/components/ui/button';
 * import { PublicNavbar } from '@/components/layouts/navigations/PublicNavbar';
 *
 * // Use in your component
 * function MyPage() {
 *   return (
 *     <>
 *       <PublicNavbar variant="primary" />
 *       <main className="pt-20">
 *         <Button variant="primary" size="xl">
 *           Sign in with Web3 Wallet Provider
 *         </Button>
 *       </main>
 *     </>
 *   );
 * }
 */
