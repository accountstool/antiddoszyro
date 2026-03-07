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
      ? "from-orange-100 to-orange-50 dark:from-orange-950/40 dark:to-orange-900/10"
      : tone === "slate"
        ? "from-slate-100 to-slate-50 dark:from-slate-900 dark:to-slate-950"
        : "from-teal-100 to-cyan-50 dark:from-teal-950/30 dark:to-slate-950";

  return (
    <Card className={`bg-gradient-to-br ${toneClass}`}>
      <p className="text-xs uppercase tracking-[0.16em] text-slate-500">{label}</p>
      <p className="mt-3 text-3xl font-extrabold text-slate-900 dark:text-slate-100">{value}</p>
    </Card>
  );
}
