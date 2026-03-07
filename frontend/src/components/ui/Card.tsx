import clsx from "clsx";
import type { PropsWithChildren } from "react";

export function Card({ children, className }: PropsWithChildren<{ className?: string }>) {
  return (
    <section
      className={clsx(
        "rounded-[26px] border border-black/10 bg-white p-6 shadow-[0_18px_44px_rgba(15,15,15,0.07)] dark:border-white/10 dark:bg-neutral-950 dark:shadow-[0_18px_44px_rgba(0,0,0,0.28)]",
        className
      )}
    >
      {children}
    </section>
  );
}
