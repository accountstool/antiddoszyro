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
      ? "from-orange-100/90 via-orange-50/75 to-white dark:from-orange-950/35 dark:via-orange-950/15 dark:to-slate-950"
      : tone === "slate"
        ? "from-slate-100/95 via-slate-50/70 to-white dark:from-slate-900 dark:via-slate-950/70 dark:to-slate-950"
        : "from-teal-100/90 via-cyan-50/70 to-white dark:from-teal-950/30 dark:via-slate-950/80 dark:to-slate-950";

  const orbClass =
    tone === "orange"
      ? "bg-orange-500/15 text-orange-700 dark:bg-orange-400/15 dark:text-orange-200"
      : tone === "slate"
        ? "bg-slate-500/15 text-slate-700 dark:bg-slate-400/15 dark:text-slate-200"
        : "bg-teal-500/15 text-teal-700 dark:bg-teal-400/15 dark:text-teal-200";

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
