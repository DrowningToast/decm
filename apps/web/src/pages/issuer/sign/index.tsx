import React, { useState } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { usePendingEvents, useSignedEvents } from "@/hooks/useIssuerEvents";
import { useAuth } from "@/context/AuthContext";
import {
    StyledTabs,
    StyledTabsContent,
    StyledTabsList,
    StyledTabsTrigger,
} from "@/components/ui/styled-tabs";
import { EventTable } from "@/components/pages/IssuerSign/EventTable";
import { PaginationControls } from "@/components/pages/IssuerSign/PaginationControls";

import { useNavigate } from "@/router";

const IssuerSignPage = () => {
    const { t } = useTranslation();
    const navigate = useNavigate();
    // Pagination state
    const [currentPage, setCurrentPage] = useState(1);
    const [rowsPerPage, setRowsPerPage] = useState(10);

    const { user } = useAuth();

    // Fetch data using React Query
    const {
        data: pendingEvents = [],
        isLoading: isPendingLoading,
        error: pendingError,
    } = usePendingEvents({
        limit: rowsPerPage,
        offset: (currentPage - 1) * rowsPerPage,
        issuer_credential_id: user?.authenticationCredentialId || "",
    });

    const {
        data: signedEvents = [],
        isLoading: isSignedLoading,
        error: signedError,
    } = useSignedEvents({
        limit: rowsPerPage,
        offset: (currentPage - 1) * rowsPerPage,
        issuer_credential_id: user?.authenticationCredentialId || "",
    });

    // Calculate pagination
    const pendingTotalPages = Math.ceil(pendingEvents.length / rowsPerPage);
    const signedTotalPages = Math.ceil(signedEvents.length / rowsPerPage);

    // Handle loading states
    if (isPendingLoading || isSignedLoading) {
        return (
            <div className="container mx-auto px-4 py-8">
                <Typography variant="text" tag="p" color="foreground" className="text-center">
                    {t("common.loading")}
                </Typography>
            </div>
        );
    }

    // Handle error states
    if (pendingError || signedError) {
        return (
            <div className="container mx-auto px-4 py-8">
                <Typography variant="text" tag="p" color="destructive" className="text-center">
                    {t("common.error")}
                </Typography>
            </div>
        );
    }

    return (
        <div title={t("issuer.sign.pageTitle")}>
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                {/* Page Header */}
                <div className="mb-8">
                    <Typography
                        variant="header"
                        tag="h1"
                        color="foreground"
                        className="text-4xl font-bold mb-2"
                    >
                        {t("issuer.sign.pageTitle")}
                    </Typography>
                    <Typography variant="text" tag="p" color="foreground-alt" className="text-lg">
                        {t("issuer.sign.pageDescription")}
                    </Typography>
                </div>

                {/* Tabs Section */}
                <StyledTabs defaultValue="pending">
                    <StyledTabsList>
                        <StyledTabsTrigger value="pending">
                            {t("issuer.sign.pendingTab")}
                        </StyledTabsTrigger>
                        <StyledTabsTrigger value="signed">
                            {t("issuer.sign.signedTab")}
                        </StyledTabsTrigger>
                    </StyledTabsList>

                    <StyledTabsContent value="pending">
                        <EventTable
                            events={pendingEvents}
                            type="pending"
                            onActionClick={(eventId) => {
                                navigate("/issuer/sign/:eventId", {
                                    params: {
                                        eventId,
                                    },
                                });
                            }}
                        />
                        <PaginationControls
                            currentPage={currentPage}
                            totalPages={pendingTotalPages}
                            rowsPerPage={rowsPerPage}
                            onPageChange={setCurrentPage}
                            onRowsPerPageChange={setRowsPerPage}
                        />
                    </StyledTabsContent>
                    <StyledTabsContent value="signed">
                        <EventTable
                            events={signedEvents}
                            type="signed"
                            onActionClick={(eventId) => {
                                navigate("/issuer/sign/:eventId", {
                                    params: {
                                        eventId,
                                    },
                                });
                            }}
                        />
                        <PaginationControls
                            currentPage={currentPage}
                            totalPages={signedTotalPages}
                            rowsPerPage={rowsPerPage}
                            onPageChange={setCurrentPage}
                            onRowsPerPageChange={setRowsPerPage}
                        />
                    </StyledTabsContent>
                </StyledTabs>
            </div>
        </div>
    );
};

export default IssuerSignPage;
