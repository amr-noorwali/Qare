import * as React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';
const buttonVariants = cva('interactive-lift inline-flex items-center justify-center gap-2 rounded-xl border px-6 py-3 text-sm font-bold transition duration-300 disabled:opacity-50', { variants: { variant: { primary: 'border-ink bg-ink text-white hover:bg-white hover:text-ink', outline: 'border-ink/25 bg-white/60 text-ink hover:border-ink hover:bg-white', ghost: 'border-transparent text-ink hover:bg-ink/5' } }, defaultVariants: { variant: 'primary' } });
export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement>, VariantProps<typeof buttonVariants> {}
export function Button({ className, variant, ...props }: ButtonProps) { return <button className={cn(buttonVariants({ variant }), className)} {...props} />; }
