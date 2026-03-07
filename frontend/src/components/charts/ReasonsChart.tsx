import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";
import { useTranslation } from "react-i18next";

import type { RankedMetric } from "../../types/api";
import { formatNumber, formatPercent } from "../../utils/format";

const colors = ["#e5e7eb", "#cbd5e1", "#94a3b8", "#6b7280", "#a1a1aa", "#525252"];

export function ReasonsChart({ data }: { data: RankedMetric[] }) {
  const { t } = useTranslation();
  const total = data.reduce((sum, item) => sum + item.value, 0);

  if (data.length === 0) {
    return (
      <div className="grid h-72 place-items-center rounded-[22px] border border-dashed border-black/12 bg-neutral-50 text-sm text-slate-500 dark:border-white/12 dark:bg-neutral-900/60 dark:text-slate-300">
        {t("messages.noDataYet")}
      </div>
    );
  }

  return (
    <div className="grid gap-5 lg:grid-cols-[220px_1fr] lg:items-center">
      <div className="h-72">
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie
              data={data}
              innerRadius={62}
              outerRadius={92}
              paddingAngle={3}
              cornerRadius={8}
              dataKey="value"
              nameKey="name"
              stroke="rgba(255,255,255,0.28)"
              strokeWidth={1.2}
            >
              {data.map((entry, index) => (
                <Cell key={entry.name} fill={colors[index % colors.length]} />
              ))}
            </Pie>
            <Tooltip
              contentStyle={{
                borderRadius: "18px",
                border: "1px solid rgba(148,163,184,0.22)",
                background: "rgba(9,17,31,0.94)",
                color: "#f8fafc",
                boxShadow: "0 24px 60px rgba(2,8,23,0.24)"
              }}
              itemStyle={{ color: "#e2e8f0" }}
              labelStyle={{ color: "#94a3b8" }}
            />
          </PieChart>
        </ResponsiveContainer>
      </div>
      <div className="space-y-3">
        {data.map((entry, index) => (
          <div key={entry.name} className="rounded-2xl border border-black/10 bg-neutral-50 px-4 py-3 dark:border-white/10 dark:bg-neutral-900/70">
            <div className="flex items-start justify-between gap-4">
              <div className="flex min-w-0 items-center gap-3">
                <span className="mt-0.5 h-3 w-3 shrink-0 rounded-full" style={{ backgroundColor: colors[index % colors.length] }} />
                <span className="truncate text-sm font-semibold text-slate-950 dark:text-slate-50">{entry.name}</span>
              </div>
              <span className="font-mono text-sm text-slate-700 dark:text-slate-200">{formatNumber(entry.value)}</span>
            </div>
            <div className="mt-2 text-xs text-slate-500 dark:text-slate-400">
              {formatPercent(total > 0 ? entry.value / total : 0)}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
