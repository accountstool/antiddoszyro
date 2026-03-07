import clsx from "clsx";
import type { ButtonHTMLAttributes, PropsWithChildren } from "react";

type Props = PropsWithChildren<
  ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: "primary" | "secondary" | "danger" | "ghost";
    block?: boolean;
  }
>;

export function Button({ children, className, variant = "primary", block, ...props }: Props) {
  return (
    <button
      className={clsx(
        "inline-flex items-center justify-center gap-2 rounded-2xl px-4 py-2.5 text-sm font-semibold tracking-[-0.01em] transition duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-transparent disabled:cursor-not-allowed disabled:opacity-60",
        block && "w-full",
        variant === "primary" &&
          "bg-black text-white shadow-lg shadow-black/18 hover:translate-y-[-1px] hover:bg-neutral-800 hover:shadow-xl hover:shadow-black/20 focus:ring-black dark:bg-white dark:text-black dark:hover:bg-neutral-200 dark:focus:ring-white",
        variant === "secondary" &&
          "border border-black/12 bg-white/78 text-neutral-900 shadow-sm shadow-black/5 hover:border-black/30 hover:bg-white dark:border-white/12 dark:bg-neutral-950/82 dark:text-neutral-100 dark:hover:border-white/22 dark:hover:bg-neutral-900",
        variant === "danger" &&
          "bg-neutral-800 text-white shadow-lg shadow-black/14 hover:translate-y-[-1px] hover:bg-black hover:shadow-xl hover:shadow-black/18 focus:ring-neutral-700 dark:bg-neutral-200 dark:text-black dark:hover:bg-white",
        variant === "ghost" &&
          "bg-transparent text-slate-600 hover:bg-black/5 hover:text-black dark:text-slate-300 dark:hover:bg-white/6 dark:hover:text-white",
        className
      )}
      {...props}
    >
      {children}
    </button>
  );
}
