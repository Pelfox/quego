import type { LucideIcon } from 'lucide-react';
import type { PropsWithChildren } from 'react';
import { tv } from 'tailwind-variants';

const badge = tv({
  base: 'w-fit px-2.5 py-1 lowercase rounded-full text-xs border flex items-center justify-center gap-1 select-none',
  variants: {
    type: {
      success: 'bg-green-500/15 border-green-400/30',
      neutral: 'bg-neutral-200/5 border-neutral-400/30',
      danger: 'bg-red-500/20 border-red-400/20',
      medium: 'bg-yellow-300/20 border-yellow-400/30',
    },
  },
});

interface BadgeProps extends PropsWithChildren {
  icon?: LucideIcon;
  type?: keyof typeof badge.variants.type;
}

export function Badge({ type = 'neutral', icon: Icon, children }: BadgeProps) {
  return (
    <div className={badge({ type })}>
      {Icon && <Icon size="13" />}
      {children}
    </div>
  );
}
