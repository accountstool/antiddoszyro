import clsx from "clsx";
import type { PropsWithChildren } from "react";

export function Card({ children, className }: PropsWithChildren<{ className?: string }>) {
  return (
    <section
      className={clsx(
        "relative overflow-hidden rounded-[28px] border border-black/8 bg-white/82 p-5 shadow-panel shadow-black/8 backdrop-blur-xl before:absolute before:inset-x-6 before:top-0 before:h-px before:bg-gradient-to-r before:from-transparent before:via-black/18 before:to-transparent dark:border-white/10 dark:bg-neutral-950/82 dark:shadow-black/28",
        className
      )}
    >
      {children}
    </section>
  );
}
