import { cn } from '@/lib/utils';
import { ThemisWhite } from '@/components/icons';

interface PublicNavbarProps {
    variant?: "primary" | "dark" | "light";
    className?: string;
}

export function PublicNavbar({ variant = "primary", className }: PublicNavbarProps) {

    // Define navbar styles based on variant
    const variantStyles = {
        primary: 'bg-[#eb5331]',
        'dark': 'bg-secondary-foreground',
        'light': 'bg-foreground',
    };

    return (
        <nav className={cn(
            'fixed top-0 left-0 right-0 z-50 transition-all',
            variantStyles[variant],
            className
        )}>
            <div className="flex items-center justify-center px-4 py-1.5 md:px-8 md:py-4 max-w-[1512px] mx-auto">
                {/* Logo */}
                <a href="/" className="flex items-center transition-transform hover:scale-105">
                    <div className='size-9 md:size-12 grid place-items-center relative'>
                        <ThemisWhite className='w-full h-full' />
                    </div>
                </a>
            </div>
        </nav>
    );
}
