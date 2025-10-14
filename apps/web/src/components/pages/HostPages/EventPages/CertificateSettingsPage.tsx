import { useState, useRef } from "react";
import { useTranslation } from "react-i18next";
import { useParams } from "react-router-dom";
import DOMPurify from "dompurify";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
} from "@/components/ui/dialog";
import { Checkbox } from "@/components/ui/checkbox";
import PageContainer from "@/components/container/PageContainer";
import SectionContainer from "@/components/container/SectionContainer";
import { Search, Trash2, Upload, AlertCircle, Info, Image as ImageIcon } from "lucide-react";

// Type definitions
interface Issuer {
    id: string;
    name: string;
    email: string;
    organization?: string;
}

interface DetectedKeyword {
    keyword: string;
    x: number;
    y: number;
    count: number;
}

export const CertificateSettingsPage = () => {
    const { t } = useTranslation();
    const { eventId } = useParams<{ eventId: string }>();
    const fileInputRef = useRef<HTMLInputElement>(null);

    // State for Step 1: Issuer Settings
    const [searchQuery, setSearchQuery] = useState("");
    const [selectedIssuers, setSelectedIssuers] = useState<Issuer[]>([]);
    const [isSearching, setIsSearching] = useState(false);
    const [searchResults, setSearchResults] = useState<Issuer[]>([]);
    const [isSearchModalOpen, setIsSearchModalOpen] = useState(false);
    const [tempSelectedIds, setTempSelectedIds] = useState<Set<string>>(new Set());

    // State for Step 2: Certificate Template Settings
    const [svgFile, setSvgFile] = useState<File | null>(null);
    const [svgPreview, setSvgPreview] = useState<string>("");
    const [detectedKeywords, setDetectedKeywords] = useState<DetectedKeyword[]>([]);
    const [isLoading, setIsLoading] = useState(false);

    // Available keywords for certificate templates
    const availableKeywords = [
        { keyword: "{{ eventName }}", mandatory: true },
        { keyword: "{{ name }}", mandatory: true },
        { keyword: "{{ academicInstitutionName }}", mandatory: false },
        { keyword: "{{ startDate }}", mandatory: false },
        { keyword: "{{ endDate }}", mandatory: false },
    ];

    // Check if mandatory keywords are detected
    const mandatoryKeywords = availableKeywords.filter((kw) => kw.mandatory);
    const detectedKeywordStrings = detectedKeywords.map((k) => k.keyword);
    const missingMandatoryKeywords = mandatoryKeywords.filter(
        (kw) => !detectedKeywordStrings.includes(kw.keyword),
    );
    const hasMissingMandatory = missingMandatoryKeywords.length > 0;

    // Mock issuer search function
    const handleSearchIssuers = async () => {
        setIsSearching(true);
        try {
            // TODO: Implement actual API call to search issuers
            console.log("Searching for:", searchQuery);
            // Simulate API delay
            await new Promise((resolve) => setTimeout(resolve, 1000));

            // Mock results - replace with actual API data
            const mockResults: Issuer[] = [
                {
                    id: "1",
                    name: "John Doe",
                    email: "john.doe@university.edu",
                    organization: "Computer Science Department",
                },
                {
                    id: "2",
                    name: "Jane Smith",
                    email: "jane.smith@university.edu",
                    organization: "Engineering Department",
                },
                {
                    id: "3",
                    name: "Dr. Michael Brown",
                    email: "m.brown@university.edu",
                    organization: "Mathematics Department",
                },
                {
                    id: "4",
                    name: "Prof. Sarah Wilson",
                    email: "s.wilson@university.edu",
                    organization: "Physics Department",
                },
            ];

            // Set search results and open modal
            setSearchResults(mockResults);
            setIsSearchModalOpen(true);

            // Pre-select already selected issuers
            const alreadySelectedIds = new Set(selectedIssuers.map((issuer) => issuer.id));
            setTempSelectedIds(alreadySelectedIds);
        } catch (error) {
            console.error("Error searching issuers:", error);
        } finally {
            setIsSearching(false);
        }
    };

    const handleToggleIssuerSelection = (issuerId: string) => {
        const newSelection = new Set(tempSelectedIds);
        if (newSelection.has(issuerId)) {
            newSelection.delete(issuerId);
        } else {
            newSelection.add(issuerId);
        }
        setTempSelectedIds(newSelection);
    };

    const handleConfirmIssuerSelection = () => {
        // Get all selected issuers from search results
        const newSelectedIssuers = searchResults.filter((issuer) => tempSelectedIds.has(issuer.id));

        // Merge with existing selections (avoid duplicates)
        const merged = [...selectedIssuers];
        newSelectedIssuers.forEach((issuer) => {
            if (!merged.some((existing) => existing.id === issuer.id)) {
                merged.push(issuer);
            }
        });

        // Update selected issuers
        setSelectedIssuers(merged);

        // Close modal and reset
        setIsSearchModalOpen(false);
        setSearchResults([]);
        setTempSelectedIds(new Set());
    };

    const handleCancelIssuerSelection = () => {
        setIsSearchModalOpen(false);
        setSearchResults([]);
        setTempSelectedIds(new Set());
    };

    const handleRemoveIssuer = (issuerId: string) => {
        setSelectedIssuers(selectedIssuers.filter((issuer) => issuer.id !== issuerId));
    };

    // SVG file handling
    const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (file && file.type === "image/svg+xml") {
            setSvgFile(file);

            // Read file for preview
            const reader = new FileReader();
            reader.onload = (e) => {
                const content = e.target?.result as string;
                setSvgPreview(content);
                parseKeywordsFromSVG(content);
            };
            reader.readAsText(file);
        }
    };

    // Parse keywords from SVG content
    const parseKeywordsFromSVG = (svgContent: string) => {
        const keywords: DetectedKeyword[] = [];
        const keywordRegex = /\{\{\s*(\w+)\s*\}\}/g;

        // Parse SVG to extract text elements with keywords and their positions
        const parser = new DOMParser();
        const svgDoc = parser.parseFromString(svgContent, "image/svg+xml");
        const textElements = svgDoc.querySelectorAll("text, tspan");

        textElements.forEach((element) => {
            const text = element.textContent || "";
            const matches = text.matchAll(keywordRegex);

            for (const match of matches) {
                const keyword = `{{ ${match[1]} }}`;
                const x = parseFloat(element.getAttribute("x") || "0");
                const y = parseFloat(element.getAttribute("y") || "0");

                // Check if keyword already exists
                const existingKeyword = keywords.find((k) => k.keyword === keyword);
                if (existingKeyword) {
                    existingKeyword.count++;
                } else {
                    keywords.push({ keyword, x, y, count: 1 });
                }
            }
        });

        setDetectedKeywords(keywords);
    };

    const handleSubmit = async () => {
        setIsLoading(true);
        try {
            // TODO: Implement API call to save certificate settings
            console.log("Event ID:", eventId);
            console.log("Selected Issuers:", selectedIssuers);
            console.log("SVG File:", svgFile);
            console.log("Detected Keywords:", detectedKeywords);

            // Simulate API delay
            await new Promise((resolve) => setTimeout(resolve, 1500));

            alert(t("certificateSettings.saveSuccess"));
        } catch (error) {
            console.error("Error saving certificate settings:", error);
            alert(t("certificateSettings.saveError"));
        } finally {
            setIsLoading(false);
        }
    };

    const handleCancel = () => {
        // Reset form or navigate back
        setSelectedIssuers([]);
        setSvgFile(null);
        setSvgPreview("");
        setDetectedKeywords([]);
    };

    return (
        <PageContainer
            title={t("certificateSettings.pageTitle")}
            description={t("certificateSettings.pageDescription")}
        >
            <SectionContainer>
                <div className="space-y-8">
                    {/* Step 1: Issuer Settings */}
                    <div className="space-y-6">
                        <div>
                            <Typography
                                variant="header"
                                tag="h2"
                                className="text-xl font-bold mb-2"
                            >
                                {t("certificateSettings.step1.title")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-sm text-muted-foreground"
                            >
                                {t("certificateSettings.step1.description")}
                            </Typography>
                        </div>

                        {/* Search Issuers */}
                        <div className="space-y-4 rounded-lg border p-6">
                            <div>
                                <Label htmlFor="issuer-search">
                                    <Typography
                                        variant="text"
                                        tag="span"
                                        className="text-sm font-medium"
                                    >
                                        {t("certificateSettings.step1.searchLabel")}
                                    </Typography>
                                </Label>
                                <div className="flex gap-2 mt-2">
                                    <Input
                                        id="issuer-search"
                                        type="text"
                                        placeholder={t(
                                            "certificateSettings.step1.searchPlaceholder",
                                        )}
                                        value={searchQuery}
                                        onChange={(e) => setSearchQuery(e.target.value)}
                                        onKeyDown={(e) => {
                                            if (e.key === "Enter") {
                                                handleSearchIssuers();
                                            }
                                        }}
                                        disabled={isSearching}
                                    />
                                    <Button
                                        type="button"
                                        onClick={handleSearchIssuers}
                                        disabled={isSearching || !searchQuery.trim()}
                                        className="min-w-[100px]"
                                    >
                                        <Search className="h-4 w-4 mr-2" />
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="font-medium"
                                        >
                                            {isSearching
                                                ? t("common.searching")
                                                : t("certificateSettings.step1.searchButton")}
                                        </Typography>
                                    </Button>
                                </div>
                            </div>

                            {/* Selected Issuers Table */}
                            {selectedIssuers.length > 0 && (
                                <div className="mt-6">
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        className="text-sm font-medium mb-3"
                                    >
                                        {t("certificateSettings.step1.selectedIssuers")} (
                                        {selectedIssuers.length})
                                    </Typography>
                                    <div className="rounded-md border">
                                        <Table>
                                            <TableHeader>
                                                <TableRow>
                                                    <TableHead>
                                                        {t("certificateSettings.step1.table.name")}
                                                    </TableHead>
                                                    <TableHead>
                                                        {t("certificateSettings.step1.table.email")}
                                                    </TableHead>
                                                    <TableHead>
                                                        {t(
                                                            "certificateSettings.step1.table.organization",
                                                        )}
                                                    </TableHead>
                                                    <TableHead className="w-[100px]">
                                                        {t(
                                                            "certificateSettings.step1.table.actions",
                                                        )}
                                                    </TableHead>
                                                </TableRow>
                                            </TableHeader>
                                            <TableBody>
                                                {selectedIssuers.map((issuer) => (
                                                    <TableRow key={issuer.id}>
                                                        <TableCell className="font-medium">
                                                            {issuer.name}
                                                        </TableCell>
                                                        <TableCell>{issuer.email}</TableCell>
                                                        <TableCell>
                                                            {issuer.organization || "-"}
                                                        </TableCell>
                                                        <TableCell>
                                                            <Button
                                                                type="button"
                                                                variant="secondary-light"
                                                                size="sm"
                                                                onClick={() =>
                                                                    handleRemoveIssuer(issuer.id)
                                                                }
                                                            >
                                                                <Trash2 className="h-4 w-4 text-destructive" />
                                                            </Button>
                                                        </TableCell>
                                                    </TableRow>
                                                ))}
                                            </TableBody>
                                        </Table>
                                    </div>
                                </div>
                            )}

                            {/* Alert about issuer settings */}
                            <Alert variant="warning">
                                <AlertCircle className="h-4 w-4" />
                                <AlertTitle>
                                    {t("certificateSettings.step1.alert.title")}
                                </AlertTitle>
                                <AlertDescription>
                                    {t("certificateSettings.step1.alert.description")}
                                </AlertDescription>
                            </Alert>
                        </div>
                    </div>

                    {/* Step 2: Certificate Template Settings */}
                    <div className="space-y-6">
                        <div>
                            <Typography
                                variant="header"
                                tag="h2"
                                className="text-xl font-bold mb-2"
                            >
                                {t("certificateSettings.step2.title")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-sm text-muted-foreground"
                            >
                                {t("certificateSettings.step2.description")}
                            </Typography>
                        </div>

                        <div className="space-y-6 rounded-lg border p-6">
                            {/* Instructions */}
                            <div className="space-y-4">
                                <Typography
                                    variant="header"
                                    tag="h3"
                                    className="text-base font-semibold"
                                >
                                    {t("certificateSettings.step2.instructions.title")}
                                </Typography>

                                {/* Instruction steps with mockup images */}
                                <div className="space-y-4">
                                    {[1, 2, 3].map((step) => (
                                        <div key={step} className="flex gap-4">
                                            <div className="flex-shrink-0 w-32 h-24 bg-muted rounded-md flex items-center justify-center border-2 border-dashed">
                                                <ImageIcon className="h-8 w-8 text-muted-foreground" />
                                            </div>
                                            <div className="flex-1">
                                                <Typography
                                                    variant="text"
                                                    tag="p"
                                                    className="text-sm font-medium mb-1"
                                                >
                                                    {t(
                                                        `certificateSettings.step2.instructions.step${step}.title`,
                                                    )}
                                                </Typography>
                                                <Typography
                                                    variant="text"
                                                    tag="p"
                                                    className="text-xs text-muted-foreground"
                                                >
                                                    {t(
                                                        `certificateSettings.step2.instructions.step${step}.description`,
                                                    )}
                                                </Typography>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            {/* Available Keywords Info */}
                            <Alert variant="info">
                                <Info className="h-4 w-4" />
                                <AlertTitle>
                                    {t("certificateSettings.step2.keywords.title")}
                                </AlertTitle>
                                <AlertDescription>
                                    <div className="mt-2">
                                        <Typography
                                            variant="text"
                                            tag="p"
                                            className="text-xs mb-2 text-blue-700"
                                        >
                                            {t("certificateSettings.step2.keywords.description")}
                                        </Typography>
                                        <div className="flex flex-wrap gap-2">
                                            {availableKeywords.map((kw) => (
                                                <div
                                                    key={kw.keyword}
                                                    className="inline-flex items-center gap-1"
                                                >
                                                    <code className="px-2 py-1 bg-blue-100 text-blue-700 dark:text-blue-400 rounded text-xs font-mono">
                                                        {kw.keyword}
                                                    </code>
                                                    {kw.mandatory && (
                                                        <span className="px-1.5 py-0.5 bg-red-100 text-red-700 rounded text-[10px] font-semibold">
                                                            Required
                                                        </span>
                                                    )}
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                </AlertDescription>
                            </Alert>

                            {/* SVG Upload */}
                            <div className="space-y-2">
                                <Label htmlFor="svg-upload">
                                    <Typography
                                        variant="text"
                                        tag="span"
                                        className="text-sm font-medium"
                                    >
                                        {t("certificateSettings.step2.upload.label")}
                                    </Typography>
                                </Label>
                                <div className="flex gap-2">
                                    <Input
                                        id="svg-upload"
                                        type="file"
                                        ref={fileInputRef}
                                        accept=".svg,image/svg+xml"
                                        onChange={handleFileSelect}
                                        className="hidden"
                                    />
                                    <Button
                                        type="button"
                                        variant="secondary-light"
                                        onClick={() => fileInputRef.current?.click()}
                                        className="w-full h-30"
                                    >
                                        <Upload className="h-4 w-4 mr-2" />
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="font-medium text-black"
                                        >
                                            {svgFile
                                                ? svgFile.name
                                                : t("certificateSettings.step2.upload.button")}
                                        </Typography>
                                    </Button>
                                </div>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-xs text-muted-foreground"
                                >
                                    {t("certificateSettings.step2.upload.hint")}
                                </Typography>
                            </div>

                            {/* SVG Preview */}
                            {svgPreview && (
                                <div className="space-y-2">
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        className="text-sm font-medium"
                                    >
                                        {t("certificateSettings.step2.preview.title")}
                                    </Typography>
                                    <div className="rounded-md border bg-muted/30 p-4">
                                        <div className="w-full flex items-center justify-center">
                                            <div
                                                className="w-full max-w-4xl [&>svg]:w-full [&>svg]:h-auto [&>svg]:max-h-[500px] [&>svg]:object-contain"
                                                dangerouslySetInnerHTML={{
                                                    __html: DOMPurify.sanitize(svgPreview),
                                                }}
                                            />
                                        </div>
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>

                    {/* Action Buttons */}
                    <div className="flex justify-end gap-4 pt-4">
                        <Button
                            type="button"
                            variant="secondary-dark"
                            size="lg"
                            onClick={handleCancel}
                            disabled={isLoading}
                            className="min-w-[150px]"
                        >
                            <Typography variant="text" tag="span" className="font-medium">
                                {t("common.cancel")}
                            </Typography>
                        </Button>
                        <Button
                            type="button"
                            variant="primary"
                            size="lg"
                            onClick={handleSubmit}
                            disabled={
                                isLoading ||
                                selectedIssuers.length === 0 ||
                                !svgFile ||
                                detectedKeywords.length === 0 ||
                                hasMissingMandatory
                            }
                            className="min-w-[150px]"
                        >
                            <Typography variant="text" tag="span" className="font-medium">
                                {isLoading
                                    ? t("common.loading")
                                    : t("certificateSettings.confirmButton")}
                            </Typography>
                        </Button>
                    </div>
                </div>
            </SectionContainer>

            {/* Issuer Selection Modal */}
            <Dialog open={isSearchModalOpen} onOpenChange={setIsSearchModalOpen}>
                <DialogContent className="max-w-4xl max-h-[80vh] overflow-y-auto">
                    <DialogHeader>
                        <DialogTitle>
                            <Typography variant="header" tag="h2" className="text-xl font-bold">
                                {t("certificateSettings.step1.modal.title")}
                            </Typography>
                        </DialogTitle>
                    </DialogHeader>

                    <div className="mt-4">
                        {searchResults.length === 0 ? (
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-center text-muted-foreground py-8"
                            >
                                {t("certificateSettings.step1.modal.noResults")}
                            </Typography>
                        ) : (
                            <>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-sm text-muted-foreground mb-4"
                                >
                                    {t("certificateSettings.step1.modal.description")} (
                                    {tempSelectedIds.size}{" "}
                                    {t("certificateSettings.step1.modal.selected")})
                                </Typography>

                                <div className="rounded-md border">
                                    <Table>
                                        <TableHeader>
                                            <TableRow>
                                                <TableHead className="w-[50px]"></TableHead>
                                                <TableHead>
                                                    {t("certificateSettings.step1.table.name")}
                                                </TableHead>
                                                <TableHead>
                                                    {t("certificateSettings.step1.table.email")}
                                                </TableHead>
                                                <TableHead>
                                                    {t(
                                                        "certificateSettings.step1.table.organization",
                                                    )}
                                                </TableHead>
                                            </TableRow>
                                        </TableHeader>
                                        <TableBody>
                                            {searchResults.map((issuer) => (
                                                <TableRow
                                                    key={issuer.id}
                                                    className="cursor-pointer"
                                                    onClick={() =>
                                                        handleToggleIssuerSelection(issuer.id)
                                                    }
                                                >
                                                    <TableCell>
                                                        <Checkbox
                                                            checked={tempSelectedIds.has(issuer.id)}
                                                            onCheckedChange={() =>
                                                                handleToggleIssuerSelection(
                                                                    issuer.id,
                                                                )
                                                            }
                                                            onClick={(e) => e.stopPropagation()}
                                                        />
                                                    </TableCell>
                                                    <TableCell className="font-medium">
                                                        {issuer.name}
                                                    </TableCell>
                                                    <TableCell>{issuer.email}</TableCell>
                                                    <TableCell>
                                                        {issuer.organization || "-"}
                                                    </TableCell>
                                                </TableRow>
                                            ))}
                                        </TableBody>
                                    </Table>
                                </div>
                            </>
                        )}
                    </div>

                    <DialogFooter className="mt-6">
                        <Button
                            type="button"
                            variant="secondary-dark"
                            onClick={handleCancelIssuerSelection}
                        >
                            <Typography variant="text" tag="span" className="font-medium">
                                {t("common.cancel")}
                            </Typography>
                        </Button>
                        <Button
                            type="button"
                            variant="primary"
                            onClick={handleConfirmIssuerSelection}
                            disabled={tempSelectedIds.size === 0}
                        >
                            <Typography variant="text" tag="span" className="font-medium">
                                {t("certificateSettings.step1.modal.chooseButton")} (
                                {tempSelectedIds.size})
                            </Typography>
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </PageContainer>
    );
};
