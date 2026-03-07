import clsx from "clsx";
import type { SelectHTMLAttributes } from "react";

export function Select(props: SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      {...props}
      className={clsx(
        "w-full rounded-2xl border border-white/80 bg-white/72 px-4 py-3 text-sm text-slate-900 shadow-sm shadow-slate-900/5 outline-none transition focus:border-teal-300 focus:bg-white focus:ring-4 focus:ring-teal-500/12 dark:border-slate-800 dark:bg-slate-950/78 dark:text-slate-100 dark:focus:border-teal-700 dark:focus:bg-slate-950",
        props.className
      )}
    />
  );
}
