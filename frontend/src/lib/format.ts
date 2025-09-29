export function formatDuration(startedAt?: string | null, finishedAt?: string | null): string {
  if (!startedAt || !finishedAt)
    return '';

  const start = new Date(startedAt);
  const end = new Date(finishedAt);

  if (Number.isNaN(start.getTime()) || Number.isNaN(end.getTime()))
    return '';

  const ms = Math.abs(end.getTime() - start.getTime());
  const seconds = Math.floor(ms / 1000) % 60;
  const minutes = Math.floor(ms / (1000 * 60)) % 60;
  const hours = Math.floor(ms / (1000 * 60 * 60));

  if (hours > 0)
    return `${hours}h ${minutes}m ${seconds}s`;
  else if (minutes > 0)
    return `${minutes}m ${seconds}s`;
  else
    return `${seconds}s`;
}
