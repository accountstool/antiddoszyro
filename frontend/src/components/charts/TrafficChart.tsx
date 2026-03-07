import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { useTranslation } from "react-i18next";

import type { TimePoint } from "../../types/api";

export function TrafficChart({ data }: { data: TimePoint[] }) {
  const { t } = useTranslation();

  if (data.length === 0) {
    return (
      <div className="grid h-80 place-items-center rounded-[22px] border border-dashed border-black/12 bg-neutral-50 text-sm text-slate-500 dark:border-white/12 dark:bg-neutral-900/60 dark:text-slate-300">
        {t("messages.noDataYet")}
      </div>
    );
  }

  return (
    <div className="h-80">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data}>
          <defs>
            <linearGradient id="allowedGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#111111" stopOpacity={0.34} />
              <stop offset="95%" stopColor="#111111" stopOpacity={0.04} />
            </linearGradient>
            <linearGradient id="blockedGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#6b7280" stopOpacity={0.32} />
              <stop offset="95%" stopColor="#6b7280" stopOpacity={0.04} />
            </linearGradient>
            <linearGradient id="challengeGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#cbd5e1" stopOpacity={0.28} />
              <stop offset="95%" stopColor="#cbd5e1" stopOpacity={0.03} />
            </linearGradient>
          </defs>
          <CartesianGrid strokeDasharray="3 3" stroke="#94a3b8" opacity={0.14} vertical={false} />
          <XAxis axisLine={false} dataKey="label" tick={{ fill: "#94a3b8", fontSize: 11 }} tickLine={false} />
          <YAxis allowDecimals={false} axisLine={false} tick={{ fill: "#94a3b8", fontSize: 11 }} tickLine={false} width={34} />
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
          <Area activeDot={{ r: 4 }} strokeWidth={2.5} type="monotone" dataKey="allowed" stroke="#111111" fill="url(#allowedGradient)" />
          <Area activeDot={{ r: 4 }} strokeWidth={2.5} type="monotone" dataKey="challenge" stroke="#cbd5e1" fill="url(#challengeGradient)" />
          <Area activeDot={{ r: 4 }} strokeWidth={2.5} type="monotone" dataKey="blocked" stroke="#6b7280" fill="url(#blockedGradient)" />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
