import type { LucideIcon } from 'lucide-react';
import clsx from 'clsx';

export function SidebarItem({ icon: Icon, title, href }: { icon: LucideIcon; title: string; href: string }) {
  const pathname = window.location.pathname;
  return (
    <a href={href}>
      <div className={clsx(
        'w-full py-2 px-3 cursor-pointer rounded-lg flex items-center gap-3 hover:bg-neutral-900 transition-colors',
        pathname === href && 'bg-neutral-900',
      )}
      >
        <Icon size={18} />
        <span>{title}</span>
      </div>
    </a>
  );
}
