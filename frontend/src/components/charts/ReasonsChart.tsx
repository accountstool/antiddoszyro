import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";

import type { RankedMetric } from "../../types/api";

const colors = ["#0f766e", "#c2410c", "#0369a1", "#7c3aed", "#475569", "#e11d48"];

export function ReasonsChart({ data }: { data: RankedMetric[] }) {
  return (
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
  );
}
