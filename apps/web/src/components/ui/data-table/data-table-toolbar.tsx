import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { X } from "lucide-react";

interface DataTableToolbarProps {
    searchValue: string;
    onSearchChange: (value: string) => void;
    searchPlaceholder?: string;
}

export function DataTableToolbar({
    searchValue,
    onSearchChange,
    searchPlaceholder = "Search...",
}: DataTableToolbarProps) {
    const isFiltered = searchValue.length > 0;

    return (
        <div className="flex items-center justify-between">
            <div className="flex flex-1 items-center space-x-2">
                <Input
                    placeholder={searchPlaceholder}
                    value={searchValue}
                    onChange={(event) => onSearchChange(event.target.value)}
                    className="h-8 w-[150px] lg:w-[250px]"
                />
                {isFiltered && (
                    <Button
                        variant="secondary-light"
                        onClick={() => onSearchChange("")}
                        className="h-8 px-2 lg:px-3"
                    >
                        Reset
                        <X className="ml-2 h-4 w-4" />
                    </Button>
                )}
            </div>
        </div>
    );
}
