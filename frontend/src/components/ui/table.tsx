import type { PropsWithChildren } from 'react';

export function TableRoot({ children }: PropsWithChildren) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full">
        {children}
      </table>
    </div>
  );
}

export function TableHeader({ children }: PropsWithChildren) {
  return (
    <thead className="text-xs border-b border-neutral-800">
      <tr>{children}</tr>
    </thead>
  );
}

export function TableHeaderCell({ children }: PropsWithChildren) {
  return (
    <th className="first:px-6 py-3 text-left font-medium uppercase tracking-wider text-neutral-400">
      {children}
    </th>
  );
}

export function TableBody({ children }: PropsWithChildren) {
  return (
    <tbody className="divide-y divide-neutral-800">
      {children}
    </tbody>
  );
}

export function TableRow({ children }: PropsWithChildren) {
  return (
    <tr className="hover:bg-neutral-900/50 transition-colors">
      {children}
    </tr>
  );
}

export function TableCell({ children }: PropsWithChildren) {
  return (
    <td className="first:px-6 py-4 text-sm">
      {children}
    </td>
  );
}
