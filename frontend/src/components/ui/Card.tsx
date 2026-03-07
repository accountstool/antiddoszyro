import clsx from "clsx";
import type { PropsWithChildren } from "react";

export function Card({ children, className }: PropsWithChildren<{ className?: string }>) {
  return (
    <section
      className={clsx(
        "relative overflow-hidden rounded-[28px] border border-white/75 bg-white/82 p-5 shadow-panel shadow-slate-900/10 backdrop-blur-xl before:absolute before:inset-x-6 before:top-0 before:h-px before:bg-gradient-to-r before:from-transparent before:via-teal-500/40 before:to-transparent dark:border-slate-800/90 dark:bg-slate-950/72 dark:shadow-black/20",
        className
      )}
    >
      {children}
    </section>
  );
}
