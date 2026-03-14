import { useState } from "react";
import { useSearchParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { certificateService } from "@/services/services";
import { PublicNavbar } from "@/components/layouts/navigations/PublicNavbar";
import { Input } from "@/components/ui/input";
import type { GetCertificateShareDataResult } from "@/services/CertificateService/mapper";

export default function VerifyPage() {
    const [searchParams] = useSearchParams();
    const handle = searchParams.get("handle");

    const [password, setPassword] = useState("");
    const [passwordError, setPasswordError] = useState("");
    const [unlockedData, setUnlockedData] = useState<GetCertificateShareDataResult | null>(null);
    const [isUnlocking, setIsUnlocking] = useState(false);

    const {
        data: statusData,
        isLoading,
        isError,
    } = useQuery({
        queryKey: ["certificateShareStatus", handle],
        queryFn: () => certificateService.getCertificateShareStatus(handle!),
        enabled: !!handle,
    });

    if (!handle) {
        return (
            <div>
                <PublicNavbar />
                <main className="pt-16 flex items-center justify-center min-h-screen">
                    <p>Invalid Link</p>
                </main>
            </div>
        );
    }

    if (isLoading) {
        return (
            <div>
                <PublicNavbar />
                <main className="pt-16 flex items-center justify-center min-h-screen">
                    <p>Loading certificate…</p>
                </main>
            </div>
        );
    }

    if (isError) {
        return (
            <div>
                <PublicNavbar />
                <main className="pt-16 flex items-center justify-center min-h-screen">
                    <p>Not Found</p>
                </main>
            </div>
        );
    }

    if (statusData?.status === "PASSWORD_LOCKED") {
        if (unlockedData) {
            return (
                <div>
                    <PublicNavbar />
                    <main className="pt-16 flex items-center justify-center min-h-screen">
                        <h1>Certificate Verified</h1>
                    </main>
                </div>
            );
        }

        const handleUnlock = async () => {
            if (!password) {
                setPasswordError("Password is required.");
                return;
            }
            setIsUnlocking(true);
            setPasswordError("");
            try {
                const result = await certificateService.getCertificateShareDataWithPassword(
                    handle,
                    password,
                );
                setUnlockedData(result);
            } catch {
                setPasswordError("Incorrect password. Please try again.");
            } finally {
                setIsUnlocking(false);
            }
        };

        return (
            <div>
                <PublicNavbar />
                <main className="pt-16 flex items-center justify-center min-h-screen">
                    <div>
                        <p>This certificate is password protected.</p>
                        <Input
                            type="password"
                            placeholder="Enter password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            aria-label="Password"
                        />
                        {passwordError && <p role="alert">{passwordError}</p>}
                        <button onClick={handleUnlock} disabled={isUnlocking}>
                            Unlock
                        </button>
                    </div>
                </main>
            </div>
        );
    }

    if (statusData?.status === "VALID_BUT_PENDING") {
        return (
            <div>
                <PublicNavbar />
                <main className="pt-16 flex items-center justify-center min-h-screen">
                    <p>Certificate Pending</p>
                </main>
            </div>
        );
    }

    if (statusData?.status === "READY" && statusData?.certificate) {
        return (
            <div>
                <PublicNavbar />
                <main className="pt-16 flex items-center justify-center min-h-screen">
                    <h1>{statusData.certificate.certificateTitle}</h1>
                </main>
            </div>
        );
    }

    return (
        <div>
            <PublicNavbar />
            <main className="pt-16 flex items-center justify-center min-h-screen">
                <p>Not Found</p>
            </main>
        </div>
    );
}
