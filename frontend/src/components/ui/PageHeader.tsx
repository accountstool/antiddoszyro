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
      <div className="relative overflow-hidden rounded-[inherit] bg-gradient-to-br from-white via-neutral-100/78 to-neutral-200/58 p-6 dark:from-black dark:via-neutral-950 dark:to-neutral-900">
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
          "border-black bg-black text-white dark:border-white dark:bg-white dark:text-black",
        tone === "warm" &&
          "border-neutral-400/80 bg-neutral-100/90 text-black dark:border-neutral-600 dark:bg-neutral-900 dark:text-white",
        tone === "default" &&
          "border-black/12 bg-white/84 text-slate-900 dark:border-white/10 dark:bg-neutral-950/82 dark:text-slate-100"
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
