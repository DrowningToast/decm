import { cn } from '@/lib/utils';
import { ThemisBlack, ThemisWhite } from '@/components/icons';

interface PublicNavbarProps {
    className?: string;
    variant?: "light" | "dark";
}

export function PublicNavbar({ className, variant = "light" }: PublicNavbarProps) {
    return (
        <nav className={cn(
            'fixed top-0 left-0 right-0 z-50 transition-all bg-foreground-alt',
            cn({
                'bg-white/10': variant === 'light',
            }),
            className
        )}>
            <div className="flex items-center justify-center px-4 py-1.5 md:px-8 md:py-4 max-w-[1512px] mx-auto">
                {/* Logo */}
                <a href="/" className="flex items-center transition-transform hover:scale-105">
                    <div className='size-9 md:size-12 grid place-items-center relative'>
                        {variant === "light" ? <ThemisWhite className='w-full h-full' /> : <ThemisBlack className='w-full h-full' />}
                    </div>
                </a>
            </div>
        </nav>
    );
}
