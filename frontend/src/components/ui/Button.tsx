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
        "rounded-xl px-4 py-2.5 text-sm font-semibold transition focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-transparent disabled:cursor-not-allowed disabled:opacity-60",
        block && "w-full",
        variant === "primary" &&
          "bg-tide text-white shadow-lg shadow-teal-900/10 hover:bg-teal-700 focus:ring-teal-500",
        variant === "secondary" &&
          "border border-slate-300/60 bg-white/80 text-slate-800 hover:bg-slate-50 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100 dark:hover:bg-slate-800",
        variant === "danger" &&
          "bg-ember text-white hover:bg-orange-700 focus:ring-orange-500",
        variant === "ghost" &&
          "bg-transparent text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800",
        className
      )}
      {...props}
    >
      {children}
    </button>
  );
}
