import clsx from "clsx";
import type { SelectHTMLAttributes } from "react";

export function Select(props: SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      {...props}
      className={clsx(
        "w-full rounded-2xl border border-black/14 bg-neutral-50 px-4 py-3 text-sm font-medium text-neutral-950 shadow-sm shadow-black/5 outline-none transition focus:border-black/35 focus:bg-white focus:ring-4 focus:ring-black/8 dark:border-white/10 dark:bg-neutral-900 dark:text-neutral-100 dark:focus:border-white/24 dark:focus:bg-neutral-900 dark:focus:ring-white/10",
        props.className
      )}
    />
  );
}
