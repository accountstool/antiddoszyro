import { Card } from "./Card";

export function StatCard({
  label,
  value,
  tone = "teal"
}: {
  label: string;
  value: string;
  tone?: "teal" | "orange" | "slate";
}) {
  const toneClass =
    tone === "orange"
      ? "border-neutral-300 bg-neutral-100 text-black dark:border-neutral-700 dark:bg-neutral-900 dark:text-white"
      : tone === "slate"
        ? "border-black bg-black text-white dark:border-white dark:bg-white dark:text-black"
        : "border-black/10 bg-white text-black dark:border-white/10 dark:bg-neutral-950 dark:text-white";

  const orbClass =
    tone === "orange"
      ? "border-neutral-400 bg-white text-black dark:border-neutral-600 dark:bg-neutral-950 dark:text-white"
      : tone === "slate"
        ? "border-white/18 bg-white/12 text-white dark:border-black/10 dark:bg-black/8 dark:text-black"
        : "border-black/10 bg-black text-white dark:border-white/12 dark:bg-white dark:text-black";

  return (
    <Card className={toneClass}>
      <div className="flex items-start justify-between gap-4">
        <div>
          <p
            className={`text-[11px] uppercase tracking-[0.18em] ${
              tone === "slate" ? "text-white/62 dark:text-black/58" : "text-slate-500 dark:text-slate-400"
            }`}
          >
            {label}
          </p>
          <p
            className={`mt-4 font-display text-[2rem] font-bold leading-none tracking-[-0.05em] ${
              tone === "slate" ? "text-white dark:text-black" : "text-slate-950 dark:text-slate-50"
            }`}
          >
            {value}
          </p>
        </div>
        <div className={`grid h-11 w-11 place-items-center rounded-2xl border ${orbClass}`}>
          <div className="h-3.5 w-3.5 rounded-full bg-current opacity-80" />
        </div>
      </div>
    </Card>
  );
}
