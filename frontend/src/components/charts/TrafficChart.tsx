import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

import type { TimePoint } from "../../types/api";

export function TrafficChart({ data }: { data: TimePoint[] }) {
  return (
    <div className="h-80">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data}>
          <defs>
            <linearGradient id="allowedGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#0f766e" stopOpacity={0.45} />
              <stop offset="95%" stopColor="#0f766e" stopOpacity={0.05} />
            </linearGradient>
            <linearGradient id="blockedGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#c2410c" stopOpacity={0.45} />
              <stop offset="95%" stopColor="#c2410c" stopOpacity={0.05} />
            </linearGradient>
          </defs>
          <CartesianGrid strokeDasharray="3 3" stroke="#cbd5e1" opacity={0.2} />
          <XAxis dataKey="label" hide />
          <YAxis allowDecimals={false} />
          <Tooltip />
          <Area type="monotone" dataKey="allowed" stroke="#0f766e" fill="url(#allowedGradient)" />
          <Area type="monotone" dataKey="blocked" stroke="#c2410c" fill="url(#blockedGradient)" />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
