import * as React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';
const buttonVariants = cva('inline-flex items-center justify-center gap-2 rounded-xl px-5 py-2.5 text-sm font-bold transition disabled:opacity-50', { variants: { variant: { primary: 'bg-moss text-white hover:bg-ink', outline: 'border border-moss/30 bg-white text-moss hover:bg-cream', ghost: 'text-moss hover:bg-moss/10' } }, defaultVariants: { variant: 'primary' } });
export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement>, VariantProps<typeof buttonVariants> {}
export function Button({ className, variant, ...props }: ButtonProps) { return <button className={cn(buttonVariants({ variant }), className)} {...props} />; }
