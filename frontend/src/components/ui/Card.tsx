import clsx from "clsx";
import type { PropsWithChildren } from "react";

export function Card({ children, className }: PropsWithChildren<{ className?: string }>) {
  return (
    <section
      className={clsx(
        "rounded-3xl border border-slate-200/70 bg-white/85 p-5 shadow-panel backdrop-blur dark:border-slate-800 dark:bg-slate-950/70",
        className
      )}
    >
      {children}
    </section>
  );
}
