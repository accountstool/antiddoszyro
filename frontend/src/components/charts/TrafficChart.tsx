import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

import type { TimePoint } from "../../types/api";

export function TrafficChart({ data }: { data: TimePoint[] }) {
  return (
    <div className="h-80">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data}>
          <defs>
            <linearGradient id="allowedGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#0f766e" stopOpacity={0.5} />
              <stop offset="95%" stopColor="#0f766e" stopOpacity={0.04} />
            </linearGradient>
            <linearGradient id="blockedGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#ea580c" stopOpacity={0.45} />
              <stop offset="95%" stopColor="#ea580c" stopOpacity={0.04} />
            </linearGradient>
          </defs>
          <CartesianGrid strokeDasharray="3 3" stroke="#94a3b8" opacity={0.16} vertical={false} />
          <XAxis axisLine={false} dataKey="label" hide tickLine={false} />
          <YAxis allowDecimals={false} axisLine={false} tickLine={false} width={34} />
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
          <Area activeDot={{ r: 4 }} strokeWidth={2.5} type="monotone" dataKey="allowed" stroke="#0f766e" fill="url(#allowedGradient)" />
          <Area activeDot={{ r: 4 }} strokeWidth={2.5} type="monotone" dataKey="blocked" stroke="#ea580c" fill="url(#blockedGradient)" />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
