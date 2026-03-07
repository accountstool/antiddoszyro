import type { PropsWithChildren } from "react";

export function Table({ children }: PropsWithChildren) {
  return (
    <div className="panel-scroll overflow-x-auto rounded-[24px] border border-black/10 bg-white/50 dark:border-white/10 dark:bg-neutral-950/40">
      <table className="min-w-full text-left text-sm">{children}</table>
    </div>
  );
}

export function THead({ children }: PropsWithChildren) {
  return (
    <thead className="border-b border-slate-200/80 bg-slate-50/75 text-[11px] uppercase tracking-[0.22em] text-slate-500 dark:border-slate-800 dark:bg-slate-900/55">
      {children}
    </thead>
  );
}

export function TBody({ children }: PropsWithChildren) {
  return <tbody className="divide-y divide-slate-100/90 dark:divide-slate-900">{children}</tbody>;
}

export function TR({ children }: PropsWithChildren) {
  return <tr className="transition hover:bg-white/70 dark:hover:bg-slate-900/55">{children}</tr>;
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
