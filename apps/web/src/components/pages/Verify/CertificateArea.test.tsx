import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { CertificateArea } from "./CertificateArea";

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
        i18n: { language: "en" },
    }),
}));

const mockShareData = {
    payload: {
        header: { "@context": [], id: "vc-1", issuanceDate: "", issuer: "", type: [] },
        data: {
            certificateId: "cert-1",
            certificateTitle: "My Cert",
            certificateSubtitle: "",
            certificateTokenId: "1",
            eventName: "Event",
            eventDescription: "",
            issuedAt: "",
            issuerAddresses: "",
            issuerId: "",
            receiverAddress: "",
            status: "VALID",
            userId: "",
            backendEncryptedUserData: "",
            encryptedUserData: "",
        },
        proof: {
            hash: "",
            encryptedByBackendRawData: "",
            encryptedByUserRawData: "",
            signMessage: "{}",
            host: { publicKey: "", signature: "" },
            issuers: [],
        },
    },
    contract: {
        eventCertificateContractAddress: "0xabc",
        certificateTokenId: "1",
    },
};

describe("CertificateArea", () => {
    it("shows empty frame when no handle", () => {
        render(
            <CertificateArea
                handle={null}
                shareStatus={null}
                shareData={null}
                isLoading={false}
                isError={false}
                errorStatus={undefined}
                imageUrl={undefined}
                isImageLoading={false}
            />,
        );
        expect(screen.getByText("certificateVerify.emptyFrame")).toBeInTheDocument();
    });

    it("shows spinner while loading", () => {
        render(
            <CertificateArea
                handle="abc"
                shareStatus={null}
                shareData={null}
                isLoading={true}
                isError={false}
                errorStatus={undefined}
                imageUrl={undefined}
                isImageLoading={false}
            />,
        );
        expect(document.querySelector(".animate-spin")).toBeInTheDocument();
    });

    it("shows not found error on non-rate-limited error", () => {
        render(
            <CertificateArea
                handle="abc"
                shareStatus={null}
                shareData={null}
                isLoading={false}
                isError={true}
                errorStatus={404}
                imageUrl={undefined}
                isImageLoading={false}
            />,
        );
        expect(screen.getByText("certificateVerify.notFound")).toBeInTheDocument();
    });

    it("shows rate limited message on 429 error", () => {
        render(
            <CertificateArea
                handle="abc"
                shareStatus={null}
                shareData={null}
                isLoading={false}
                isError={true}
                errorStatus={429}
                imageUrl={undefined}
                isImageLoading={false}
            />,
        );
        expect(screen.getByText("certificateVerify.rateLimited")).toBeInTheDocument();
    });

    it("shows empty frame when password locked", () => {
        render(
            <CertificateArea
                handle="abc"
                shareStatus="PASSWORD_LOCKED"
                shareData={null}
                isLoading={false}
                isError={false}
                errorStatus={undefined}
                imageUrl={undefined}
                isImageLoading={false}
            />,
        );
        expect(screen.getByText("certificateVerify.emptyFrame")).toBeInTheDocument();
    });

    it("shows certificate title when ready and no image", () => {
        render(
            <CertificateArea
                handle="abc"
                shareStatus="READY"
                shareData={mockShareData}
                isLoading={false}
                isError={false}
                errorStatus={undefined}
                imageUrl={undefined}
                isImageLoading={false}
            />,
        );
        expect(screen.getByText("My Cert")).toBeInTheDocument();
    });

    it("shows image when ready and imageUrl provided", () => {
        render(
            <CertificateArea
                handle="abc"
                shareStatus="READY"
                shareData={mockShareData}
                isLoading={false}
                isError={false}
                errorStatus={undefined}
                imageUrl="http://example.com/cert.png"
                isImageLoading={false}
            />,
        );
        const img = screen.getByRole("img");
        expect(img).toHaveAttribute("src", "http://example.com/cert.png");
        expect(img).toHaveAttribute("alt", "My Cert");
    });
});
