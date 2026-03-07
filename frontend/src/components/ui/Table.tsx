import type { PropsWithChildren } from "react";

export function Table({ children }: PropsWithChildren) {
  return (
    <div className="panel-scroll overflow-x-auto">
      <table className="min-w-full text-left text-sm">{children}</table>
    </div>
  );
}

export function THead({ children }: PropsWithChildren) {
  return <thead className="border-b border-slate-200 text-xs uppercase tracking-[0.18em] text-slate-500 dark:border-slate-800">{children}</thead>;
}

export function TBody({ children }: PropsWithChildren) {
  return <tbody className="divide-y divide-slate-100 dark:divide-slate-900">{children}</tbody>;
}

export function TR({ children }: PropsWithChildren) {
  return <tr>{children}</tr>;
}

export function TH({ children }: PropsWithChildren) {
  return <th className="px-4 py-3 font-semibold">{children}</th>;
}

export function TD({ children, className = "" }: PropsWithChildren<{ className?: string }>) {
  return <td className={`px-4 py-3 align-top ${className}`}>{children}</td>;
}
