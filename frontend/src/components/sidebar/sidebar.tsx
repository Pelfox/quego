import { SidebarItem } from '@components/sidebar/sidebar-item';
import { WorkflowIcon } from 'lucide-react';

export function Sidebar() {
  return (
    <div className="w-64 min-h-screen border-r border-neutral-800">
      <div className="p-6 border-b border-neutral-800 flex items-center gap-2 h-18">
        <div className="w-7 h-7 rounded-md bg-orange-500 flex items-center justify-center font-semibold">Q</div>
        <span className="font-semibold text-lg">quego</span>
      </div>
      <div className="p-4 space-y-2">
        <SidebarItem icon={WorkflowIcon} title="Workflows" href="/" />
      </div>
    </div>
  );
}
