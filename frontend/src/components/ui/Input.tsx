import clsx from "clsx";
import type { InputHTMLAttributes } from "react";

export function Input(props: InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      {...props}
      className={clsx(
        "w-full rounded-2xl border border-white/80 bg-white/72 px-4 py-3 text-sm text-slate-900 shadow-sm shadow-slate-900/5 outline-none ring-0 transition placeholder:text-slate-400 focus:border-teal-300 focus:bg-white focus:ring-4 focus:ring-teal-500/12 dark:border-slate-800 dark:bg-slate-950/78 dark:text-slate-100 dark:focus:border-teal-700 dark:focus:bg-slate-950",
        props.className
      )}
    />
  );
}
