import type { PropsWithChildren } from "react";

export function Table({ children }: PropsWithChildren) {
  return (
    <div className="panel-scroll overflow-x-auto rounded-[22px] border border-black/10 bg-white dark:border-white/10 dark:bg-neutral-950">
      <table className="min-w-full text-left text-sm">{children}</table>
    </div>
  );
}

export function THead({ children }: PropsWithChildren) {
  return (
    <thead className="border-b border-black/10 bg-black text-[11px] uppercase tracking-[0.2em] text-white/62 dark:border-white/10 dark:bg-white dark:text-black/58">
      {children}
    </thead>
  );
}

export function TBody({ children }: PropsWithChildren) {
  return <tbody className="divide-y divide-black/8 dark:divide-white/8">{children}</tbody>;
}

export function TR({ children }: PropsWithChildren) {
  return <tr className="bg-white transition odd:bg-white even:bg-neutral-50 hover:bg-neutral-100 dark:bg-neutral-950 dark:odd:bg-neutral-950 dark:even:bg-neutral-900 dark:hover:bg-neutral-900">{children}</tr>;
}

export function TH({ children }: PropsWithChildren) {
  return <th className="px-4 py-3.5 font-semibold">{children}</th>;
}

export function TD({ children, className = "", colSpan }: PropsWithChildren<{ className?: string; colSpan?: number }>) {
  return (
    <td className={`px-4 py-3.5 align-top text-slate-700 dark:text-slate-200 ${className}`} colSpan={colSpan}>
      {children}
    </td>
  );
}
