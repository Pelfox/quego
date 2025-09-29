import type { Execution } from '@/types/execution';
import { DashoardLayout } from '@components/dashboard-layout';
import { getBadgeIcon, getBadgeType } from '@lib/badge';
import { formatDuration } from '@lib/format';
import { Badge } from '@ui/badge';
import { TableBody, TableCell, TableHeader, TableHeaderCell, TableRoot, TableRow } from '@ui/table';
import { Loader2Icon, ServerCrashIcon } from 'lucide-react';
import useSWR from 'swr';

export function IndexPage() {
  const { data, isLoading, error, mutate } = useSWR('/executions', {
    refreshInterval: 5000,
    revalidateOnFocus: true,
  });

  return (
    <DashoardLayout title="Workflows" mutate={mutate}>
      <div className="w-full overflow-hidden border border-neutral-800 rounded-2xl">
        <div className="flex items-center w-full justify-between border-b border-neutral-800 px-6 py-5">
          <h2 className="text-lg font-semibold">
            Recent Executions
          </h2>
          <p className="text-sm text-neutral-400">
            Showing
            {' '}
            {data?.length ?? 0}
            {' '}
            results
          </p>
        </div>
        {isLoading && (
          <div className="p-12 flex items-center justify-center w-full">
            <Loader2Icon className="animate-spin" />
          </div>
        )}
        {error && (
          <div className="p-12 flex flex-col items-center justify-center w-full">
            <ServerCrashIcon />
            <div className="mt-3 text-center">
              <h3 className="mt-2 text-lg font-semibold">Uh-oh! Loading failed</h3>
              <p>{error.message ?? 'Unknown error.'}</p>
            </div>
          </div>
        )}
        {(!isLoading && !error) && (
          <TableRoot>
            <TableHeader>
              <TableHeaderCell>Execution ID</TableHeaderCell>
              <TableHeaderCell>Function name</TableHeaderCell>
              <TableHeaderCell>Status</TableHeaderCell>
              <TableHeaderCell>Duration</TableHeaderCell>
              <TableHeaderCell>Started</TableHeaderCell>
              <TableHeaderCell>Trigger</TableHeaderCell>
            </TableHeader>
            <TableBody>
              {data && data.map((execution: Execution) => (
                <TableRow key={execution.id}>
                  <TableCell>{execution.id}</TableCell>
                  <TableCell>{execution.trigger.function_name}</TableCell>
                  <TableCell>
                    <Badge icon={getBadgeIcon(execution.status)} type={getBadgeType(execution.status)}>
                      {execution.status}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    {formatDuration(execution.started_at, execution.finished_at)}
                  </TableCell>
                  <TableCell>
                    {execution.started_at && new Date(execution.started_at).toLocaleString()}
                  </TableCell>
                  <TableCell>{execution.trigger.trigger_type}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </TableRoot>
        )}
      </div>
    </DashoardLayout>
  );
}
