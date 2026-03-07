import clsx from "clsx";
import type { PropsWithChildren, ReactNode } from "react";

import { Card } from "./Card";

export function PageHeader({
  title,
  subtitle,
  actions,
  children,
  className
}: PropsWithChildren<{
  title: string;
  subtitle?: string;
  actions?: ReactNode;
  className?: string;
}>) {
  return (
    <Card className={clsx("p-0", className)}>
      <div className="rounded-[inherit] bg-white px-6 py-6 dark:bg-neutral-950">
        <div className="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
          <div className="min-w-0 max-w-4xl">
            <div className="text-[11px] font-semibold uppercase tracking-[0.24em] text-slate-500 dark:text-slate-400">ShieldPanel</div>
            <h1 className="mt-3 font-display text-3xl font-bold tracking-[-0.05em] text-slate-950 dark:text-slate-50 md:text-[2.35rem]">
              {title}
            </h1>
            {subtitle ? (
              <p className="mt-2 max-w-2xl text-sm leading-6 text-slate-600 dark:text-slate-300">{subtitle}</p>
            ) : null}
          </div>
          {actions ? <div className="flex shrink-0 flex-wrap items-center gap-3">{actions}</div> : null}
        </div>

        {children ? <div className="mt-6 grid gap-3 sm:grid-cols-2 xl:grid-cols-3">{children}</div> : null}
      </div>
    </Card>
  );
}

export function HeaderMetric({
  label,
  value,
  tone = "default"
}: {
  label: string;
  value: string;
  tone?: "default" | "accent" | "warm";
}) {
  return (
    <div
      className={clsx(
        "min-w-[150px] rounded-2xl border px-4 py-3.5 shadow-sm",
        tone === "accent" &&
          "border-black bg-black text-white dark:border-white dark:bg-white dark:text-black",
        tone === "warm" &&
          "border-neutral-300 bg-neutral-100 text-black dark:border-neutral-700 dark:bg-neutral-900 dark:text-white",
        tone === "default" &&
          "border-black/10 bg-white text-slate-900 dark:border-white/10 dark:bg-neutral-950 dark:text-slate-100"
      )}
    >
      <div className={clsx(
        "text-[11px] uppercase tracking-[0.18em]",
        tone === "accent" ? "text-white/70 dark:text-black/70" : "text-slate-500 dark:text-slate-400"
      )}>{label}</div>
      <div className="mt-2 font-display text-lg font-bold tracking-[-0.03em]">{value}</div>
    </div>
  );
}
