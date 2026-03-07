import clsx from "clsx";
import type { InputHTMLAttributes } from "react";

export function Input(props: InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      {...props}
      className={clsx(
        "w-full rounded-2xl border border-black/12 bg-white/78 px-4 py-3 text-sm text-neutral-900 shadow-sm shadow-black/5 outline-none ring-0 transition placeholder:text-slate-400 focus:border-black/40 focus:bg-white focus:ring-4 focus:ring-black/8 dark:border-white/10 dark:bg-neutral-950/82 dark:text-neutral-100 dark:focus:border-white/28 dark:focus:bg-neutral-950 dark:focus:ring-white/10",
        props.className
      )}
    />
  );
}
