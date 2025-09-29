export interface Trigger {
  id: string;
  function_name: string;
  trigger_type: 'EVENT' | 'CRON';
}
