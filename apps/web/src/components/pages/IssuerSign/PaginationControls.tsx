import React from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";

interface PaginationControlsProps {
    currentPage: number;
    totalPages: number;
    rowsPerPage: number;
    onPageChange: (page: number) => void;
    onRowsPerPageChange: (rowsPerPage: number) => void;
}

export const PaginationControls: React.FC<PaginationControlsProps> = ({
    currentPage,
    totalPages,
    rowsPerPage,
    onPageChange,
    onRowsPerPageChange,
}) => {
    const { t } = useTranslation();

    return (
        <div className="flex justify-between items-center text-sm text-secondary mt-6">
            <div className="flex items-center space-x-2">
                <Typography variant="text" tag="span" className="text-sm">
                    {t("issuer.sign.pagination.rowsPerPage")}
                </Typography>
                <select
                    value={rowsPerPage}
                    onChange={(e) => onRowsPerPageChange(Number(e.target.value))}
                    className="border border-gray-300 rounded px-2 py-1 text-sm"
                >
                    <option value={10}>10</option>
                    <option value={25}>25</option>
                    <option value={50}>50</option>
                </select>
            </div>
            <div>
                <Typography variant="text" tag="span" className="text-sm">
                    {Math.min(totalPages * rowsPerPage, currentPage * rowsPerPage)}-
                    {Math.min(totalPages * rowsPerPage, (currentPage + 1) * rowsPerPage)} of{" "}
                    {totalPages * rowsPerPage}
                </Typography>
                <Typography variant="text" tag="span" className="ml-4">
                    Page {currentPage} of {totalPages}
                </Typography>
            </div>
            <div className="flex space-x-2">
                <button
                    className="btn-secondary p-2 rounded-lg disabled:opacity-50"
                    disabled={currentPage <= 1}
                    onClick={() => onPageChange(currentPage - 1)}
                >
                    <svg
                        className="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth="2"
                            d="M15 19l-7-7 7-7"
                        ></path>
                    </svg>
                </button>
                <button
                    className="btn-secondary p-2 rounded-lg disabled:opacity-50"
                    disabled={currentPage >= totalPages}
                    onClick={() => onPageChange(currentPage + 1)}
                >
                    <svg
                        className="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth="2"
                            d="M9 5l7 7-7 7"
                        ></path>
                    </svg>
                </button>
            </div>
        </div>
    );
};
