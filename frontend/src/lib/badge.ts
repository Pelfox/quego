import type { Execution } from '@/types/execution';
import { CheckIcon, FileStackIcon, Loader2Icon, ServerCrashIcon } from 'lucide-react';

export function getBadgeType(status: Execution['status']) {
  switch (status) {
    case 'COMPLETED':
      return 'success';
    case 'FAILED':
      return 'danger';
    case 'PENDING':
      return 'medium';
    case 'RUNNING':
    default:
      return 'neutral';
  }
}

export function getBadgeIcon(status: Execution['status']) {
  switch (status) {
    case 'COMPLETED':
      return CheckIcon;
    case 'FAILED':
      return ServerCrashIcon;
    case 'PENDING':
      return FileStackIcon;
    case 'RUNNING':
    default:
      return Loader2Icon;
  }
}
