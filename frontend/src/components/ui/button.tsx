import type { ButtonHTMLAttributes } from 'react';
import { tv } from 'tailwind-variants';

const button = tv({
  base: 'rounded-lg border border-neutral-800 hover:bg-neutral-900 transition-colors cursor-pointer font-normal flex items-center justify-center gap-1.5',
  variants: {
    size: {
      icon: 'p-2',
      sm: 'px-3 py-2 text-xs',
      md: 'px-4 py-2 text-sm',
    },
    disabled: {
      true: 'opacity-80 cursor-default hover:bg-transparent select-none',
    },
  },
  defaultVariants: {
    size: 'md',
  },
});

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  size?: keyof typeof button.variants.size;
}

export function Button({ size, className, children, ...props }: ButtonProps) {
  return (
    <button {...props} className={button({ size, disabled: props.disabled, className })}>
      {children}
    </button>
  );
}
