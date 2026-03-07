import clsx from "clsx";
import type { InputHTMLAttributes } from "react";

export function Input(props: InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      {...props}
      className={clsx(
        "w-full rounded-2xl border border-black/14 bg-neutral-50 px-4 py-3 text-sm font-medium text-neutral-950 shadow-sm shadow-black/5 outline-none ring-0 transition placeholder:text-neutral-400 focus:border-black/35 focus:bg-white focus:ring-4 focus:ring-black/8 dark:border-white/10 dark:bg-neutral-900 dark:text-neutral-100 dark:placeholder:text-neutral-500 dark:focus:border-white/24 dark:focus:bg-neutral-900 dark:focus:ring-white/10",
        props.className
      )}
    />
  );
}
