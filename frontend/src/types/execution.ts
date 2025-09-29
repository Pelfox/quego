import type { Trigger } from '@/types/trigger';

export interface Execution {
  id: string;
  trigger_id: string;
  trigger: Trigger;
  status: 'PENDING' | 'RUNNING' | 'COMPLETED' | 'FAILED';
  started_at?: string;
  finished_at?: string;
}
