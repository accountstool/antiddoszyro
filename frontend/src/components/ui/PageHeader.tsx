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
      <div className="relative overflow-hidden rounded-[inherit] bg-gradient-to-br from-white/75 via-white/55 to-teal-50/80 p-6 dark:from-slate-950/50 dark:via-slate-950/35 dark:to-teal-950/25">
        <div className="absolute inset-0 panel-grid-overlay opacity-50" />
        <div className="relative flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
          <div className="max-w-3xl">
            <h1 className="font-display text-3xl font-bold tracking-[-0.04em] text-slate-950 dark:text-slate-50 md:text-[2.3rem]">
              {title}
            </h1>
            {subtitle ? (
              <p className="mt-2 max-w-2xl text-sm leading-6 text-slate-600 dark:text-slate-300">{subtitle}</p>
            ) : null}
            {children ? <div className="mt-5 flex flex-wrap gap-3">{children}</div> : null}
          </div>
          {actions ? <div className="relative flex shrink-0 flex-wrap items-center gap-3">{actions}</div> : null}
        </div>
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
        "rounded-2xl border px-4 py-3 shadow-sm backdrop-blur",
        tone === "accent" &&
          "border-teal-200/80 bg-teal-50/80 text-teal-900 dark:border-teal-900/70 dark:bg-teal-950/30 dark:text-teal-100",
        tone === "warm" &&
          "border-orange-200/80 bg-orange-50/80 text-orange-900 dark:border-orange-900/70 dark:bg-orange-950/25 dark:text-orange-100",
        tone === "default" &&
          "border-white/70 bg-white/72 text-slate-900 dark:border-slate-800/80 dark:bg-slate-900/55 dark:text-slate-100"
      )}
    >
      <div className="text-[11px] uppercase tracking-[0.18em] text-slate-500 dark:text-slate-400">{label}</div>
      <div className="mt-2 font-display text-lg font-bold tracking-[-0.03em]">{value}</div>
    </div>
  );
}
