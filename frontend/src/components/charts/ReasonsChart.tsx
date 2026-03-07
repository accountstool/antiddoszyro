import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";

import type { RankedMetric } from "../../types/api";

const colors = ["#0f766e", "#c2410c", "#0369a1", "#7c3aed", "#475569", "#e11d48"];

export function ReasonsChart({ data }: { data: RankedMetric[] }) {
  return (
    <div className="h-72">
      <ResponsiveContainer width="100%" height="100%">
        <PieChart>
          <Pie data={data} innerRadius={60} outerRadius={90} paddingAngle={3} dataKey="value" nameKey="name">
            {data.map((entry, index) => (
              <Cell key={entry.name} fill={colors[index % colors.length]} />
            ))}
          </Pie>
          <Tooltip />
        </PieChart>
      </ResponsiveContainer>
    </div>
  );
}
