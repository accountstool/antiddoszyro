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
          "bg-gradient-to-r from-teal-600 via-teal-500 to-cyan-600 text-white shadow-lg shadow-teal-900/20 hover:translate-y-[-1px] hover:shadow-xl hover:shadow-teal-900/20 focus:ring-teal-500",
        variant === "secondary" &&
          "border border-white/80 bg-white/75 text-slate-800 shadow-sm shadow-slate-900/5 hover:border-teal-200 hover:bg-white dark:border-slate-800 dark:bg-slate-900/75 dark:text-slate-100 dark:hover:border-slate-700 dark:hover:bg-slate-900",
        variant === "danger" &&
          "bg-gradient-to-r from-orange-600 to-orange-500 text-white shadow-lg shadow-orange-900/15 hover:translate-y-[-1px] hover:shadow-xl hover:shadow-orange-900/20 focus:ring-orange-500",
        variant === "ghost" &&
          "bg-transparent text-slate-600 hover:bg-white/60 hover:text-slate-900 dark:text-slate-300 dark:hover:bg-slate-900/70 dark:hover:text-white",
        className
      )}
      {...props}
    >
      {children}
    </button>
  );
}
