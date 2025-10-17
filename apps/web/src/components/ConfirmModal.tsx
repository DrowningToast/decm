import type { ReactNode } from "react";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from "./ui/alert-dialog";
import { TriangleAlert } from "lucide-react";

interface ConfirmModalProps {
    title: string;
    message: string;
    onConfirm: () => void;
    onCancel?: () => void;
    children?: ReactNode;
    cancelText: string;
    confirmText: string;
    icon?: ReactNode;
    open?: boolean;
    onOpenChange?: (open: boolean) => void;
    destructive?: boolean;
}

export default function ConfirmModal({
    title,
    message,
    onConfirm,
    onCancel,
    children,
    cancelText,
    confirmText,
    icon = <TriangleAlert />,
    open,
    onOpenChange,
    destructive = false,
}: ConfirmModalProps) {
    return (
        <AlertDialog open={open} onOpenChange={onOpenChange}>
            <AlertDialogTrigger asChild>{children}</AlertDialogTrigger>
            <AlertDialogContent>
                <AlertDialogHeader className="flex flex-col items-start">
                    {icon && icon}
                    <AlertDialogTitle>{title}</AlertDialogTitle>
                    <AlertDialogDescription className="leading-relaxed">
                        {message}
                    </AlertDialogDescription>
                </AlertDialogHeader>

                <AlertDialogFooter>
                    <AlertDialogCancel onClick={onCancel} autoFocus={false}>
                        {cancelText}
                    </AlertDialogCancel>
                    <AlertDialogAction
                        onClick={onConfirm}
                        autoFocus={false}
                        className={destructive ? "bg-red-500 text-white hover:bg-red-600" : ""}
                    >
                        {confirmText}
                    </AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    );
}
