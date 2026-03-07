import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { api, unwrap } from "../api/client";
import { ReasonsChart } from "../components/charts/ReasonsChart";
import { TrafficChart } from "../components/charts/TrafficChart";
import { Card } from "../components/ui/Card";
import { HeaderMetric, PageHeader } from "../components/ui/PageHeader";
import { Spinner } from "../components/ui/Spinner";
import { StatCard } from "../components/ui/StatCard";
import type { DashboardSummary, RankedMetric, TimePoint } from "../types/api";
import { formatNumber, formatPercent } from "../utils/format";

const emptySummary: DashboardSummary = {
  healthy: true,
  totalDomains: 0,
  blockedToday: 0,
  currentRps: 0,
  currentBlockedPerSecond: 0,
  topAttackedDomain: "",
  topAttackingIp: "",
  totalRequests24h: 0,
  allowed24h: 0,
  blocked24h: 0,
  challenged24h: 0,
  challengePassRate: 0
};

const emptyCharts = {
  last24h: [] as TimePoint[],
  last7d: [] as TimePoint[],
  topIps: [] as RankedMetric[],
  topDomains: [] as RankedMetric[],
  topReasons: [] as RankedMetric[]
};

export function DashboardPage() {
  const { t } = useTranslation();
  const summaryQuery = useQuery({
    queryKey: ["dashboard", "summary"],
    queryFn: () => unwrap<DashboardSummary>(api.get("/dashboard/summary"))
  });
  const chartsQuery = useQuery({
    queryKey: ["dashboard", "charts"],
    queryFn: () =>
      unwrap<{
        last24h: TimePoint[];
        last7d: TimePoint[];
        topIps: RankedMetric[];
        topDomains: RankedMetric[];
        topReasons: RankedMetric[];
      }>(api.get("/dashboard/charts"))
  });

  if (summaryQuery.isLoading || chartsQuery.isLoading) {
    return (
      <div className="grid min-h-[50vh] place-items-center">
        <Spinner />
      </div>
    );
  }

  const summary = summaryQuery.data ?? emptySummary;
  const charts = chartsQuery.data ?? emptyCharts;
  const hasError = summaryQuery.isError || chartsQuery.isError;

  return (
    <div className="space-y-6">
      <PageHeader title={t("nav.dashboard")} subtitle={t("brand.description")}>
        <HeaderMetric label={t("statistics.incoming")} value={formatNumber(summary.totalRequests24h)} tone="accent" />
        <HeaderMetric label={t("statistics.challenged")} value={formatNumber(summary.challenged24h)} tone="warm" />
        <HeaderMetric label={t("dashboard.topAttackedDomain")} value={summary.topAttackedDomain || "-"} tone="warm" />
      </PageHeader>

      {hasError ? (
        <Card className="border-black/12 bg-neutral-100/80 text-neutral-900 dark:border-white/10 dark:bg-neutral-900/80 dark:text-neutral-100">
          <h2 className="text-lg font-bold">{t("messages.requestFailed")}</h2>
          <p className="mt-1 text-sm opacity-80">{t("messages.partialData")}</p>
        </Card>
      ) : null}

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <StatCard label={t("dashboard.totalDomains")} value={formatNumber(summary.totalDomains)} />
        <StatCard label={t("statistics.allowed")} tone="slate" value={formatNumber(summary.allowed24h)} />
        <StatCard label={t("statistics.challenged")} tone="orange" value={formatNumber(summary.challenged24h)} />
        <StatCard label={t("dashboard.challengePassRate")} value={formatPercent(summary.challengePassRate)} />
      </div>

      <div className="grid gap-6 xl:grid-cols-[1.6fr_1fr]">
        <Card>
          <div className="mb-5 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
            <div>
              <h2 className="font-display text-2xl font-bold tracking-[-0.04em] text-slate-950 dark:text-slate-50">{t("dashboard.requests24h")}</h2>
              <p className="mt-1 text-sm text-slate-600 dark:text-slate-300">
                {t("dashboard.topAttackedDomain")}: <span className="font-semibold text-slate-950 dark:text-slate-50">{summary.topAttackedDomain || "-"}</span>
              </p>
            </div>
            <div className="flex flex-wrap gap-2">
              <MetaPill label={t("dashboard.currentRps")} value={formatNumber(summary.currentRps)} />
              <MetaPill label={t("dashboard.topIps")} value={summary.topAttackingIp || "-"} />
              <MetaPill label={t("statistics.blocked")} value={formatNumber(summary.blocked24h)} />
            </div>
          </div>
          <TrafficChart data={charts.last24h} />
        </Card>

        <Card>
          <div className="mb-4">
            <h2 className="font-display text-2xl font-bold tracking-[-0.04em] text-slate-950 dark:text-slate-50">{t("dashboard.topReasons")}</h2>
            <p className="mt-1 text-sm text-slate-600 dark:text-slate-300">{t("dashboard.reasonBreakdown")}</p>
          </div>
          <ReasonsChart data={charts.topReasons} />
        </Card>
      </div>

      <div className="grid gap-6 xl:grid-cols-3">
        <MetricList title={t("dashboard.topIps")} items={charts.topIps} />
        <MetricList title={t("dashboard.topDomains")} items={charts.topDomains} />
        <MetricList title={t("dashboard.topReasons")} items={charts.topReasons} />
      </div>
    </div>
  );
}

function MetricList({ title, items }: { title: string; items: RankedMetric[] }) {
  const { t } = useTranslation();
  const peak = items[0]?.value || 1;

  return (
    <Card>
      <h2 className="font-display text-xl font-bold tracking-[-0.04em] text-slate-950 dark:text-slate-50">{title}</h2>
      <div className="mt-4 space-y-3">
        {items.length === 0 ? (
          <div className="rounded-2xl border border-black/10 bg-neutral-50 px-4 py-3 text-sm text-slate-500 dark:border-white/10 dark:bg-neutral-900/70 dark:text-slate-300">
            {t("messages.noDataYet")}
          </div>
        ) : null}
        {items.map((item) => {
          const width = `${Math.max(12, Math.round((item.value / peak) * 100))}%`;
          return (
            <div key={item.name} className="rounded-2xl border border-black/10 bg-neutral-50 px-4 py-3 dark:border-white/10 dark:bg-neutral-900/70">
              <div className="flex items-center justify-between gap-3">
                <span className="truncate text-sm font-semibold text-slate-950 dark:text-slate-50">{item.name}</span>
                <span className="font-mono text-sm text-slate-700 dark:text-slate-200">{formatNumber(item.value)}</span>
              </div>
              <div className="mt-3 h-2 rounded-full bg-neutral-200 dark:bg-neutral-800">
                <div className="h-full rounded-full bg-black dark:bg-white" style={{ width }} />
              </div>
            </div>
          );
        })}
      </div>
    </Card>
  );
}

function MetaPill({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-full border border-black/10 bg-neutral-50 px-3.5 py-2 dark:border-white/10 dark:bg-neutral-900">
      <div className="text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500 dark:text-slate-400">{label}</div>
      <div className="mt-1 text-sm font-semibold text-slate-950 dark:text-slate-50">{value}</div>
    </div>
  );
}
