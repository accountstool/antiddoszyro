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
      ? "from-neutral-200/85 via-neutral-100/78 to-white dark:from-neutral-800/85 dark:via-neutral-900/78 dark:to-black"
      : tone === "slate"
        ? "from-zinc-100/95 via-zinc-50/70 to-white dark:from-zinc-900/88 dark:via-black dark:to-black"
        : "from-white via-neutral-100/75 to-neutral-50 dark:from-neutral-900/94 dark:via-neutral-950 dark:to-black";

  const orbClass =
    tone === "orange"
      ? "bg-black/10 text-black dark:bg-white/10 dark:text-white"
      : tone === "slate"
        ? "bg-neutral-700/10 text-neutral-800 dark:bg-neutral-300/10 dark:text-neutral-200"
        : "bg-black text-white dark:bg-white dark:text-black";

  return (
    <Card className={`bg-gradient-to-br ${toneClass}`}>
      <div className="flex items-start justify-between gap-4">
        <div>
          <p className="text-[11px] uppercase tracking-[0.18em] text-slate-500 dark:text-slate-400">{label}</p>
          <p className="mt-4 font-display text-[2rem] font-bold leading-none tracking-[-0.05em] text-slate-950 dark:text-slate-50">{value}</p>
        </div>
        <div className={`grid h-11 w-11 place-items-center rounded-2xl ${orbClass}`}>
          <div className="h-3.5 w-3.5 rounded-full bg-current opacity-80" />
        </div>
      </div>
    </Card>
  );
}
