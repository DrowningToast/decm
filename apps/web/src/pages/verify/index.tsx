import { useState } from "react";
import { useSearchParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import { PublicNavbar } from "@/components/layouts/navigations/PublicNavbar";
import { Input } from "@/components/ui/input";

export default function VerifyPage() {
    const [searchParams] = useSearchParams();
    const handle = searchParams.get("handle");

    const [password, setPassword] = useState("");
    const [passwordError, setPasswordError] = useState("");
    const [unlockedData, setUnlockedData] = useState<{
        header: object;
        data: object;
        proof: object;
    } | null>(null);
    const [isUnlocking, setIsUnlocking] = useState(false);

    const {
        data: statusData,
        isLoading,
        isError,
    } = useQuery({
        queryKey: ["certificateShareStatus", handle],
        queryFn: () =>
            coreApiClient.v1.getCertificateShareStatus({ handle: handle! }).then((res) => res.data),
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
                const res = await coreApiClient.v1.getCertificateShareDataWithPassword(
                    { handle },
                    { password },
                );
                setUnlockedData(res.data);
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
                    <h1>{statusData.certificate.certificate_title}</h1>
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
