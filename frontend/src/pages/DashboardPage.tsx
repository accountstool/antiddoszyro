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
        <HeaderMetric label={t("dashboard.totalDomains")} value={formatNumber(summary.totalDomains)} tone="accent" />
        <HeaderMetric label={t("dashboard.currentRps")} value={formatNumber(summary.currentRps)} />
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
        <StatCard label={t("dashboard.blockedToday")} tone="orange" value={formatNumber(summary.blockedToday)} />
        <StatCard label={t("dashboard.currentRps")} tone="slate" value={formatNumber(summary.currentRps)} />
        <StatCard label={t("dashboard.challengePassRate")} value={formatPercent(summary.challengePassRate)} />
      </div>

      <div className="grid gap-6 xl:grid-cols-[1.6fr_1fr]">
        <Card>
          <div className="mb-4 flex items-center justify-between">
            <div>
              <h2 className="font-display text-2xl font-bold tracking-[-0.04em]">{t("dashboard.requests24h")}</h2>
              <p className="text-sm text-slate-500">{t("dashboard.topAttackedDomain")}: {summary.topAttackedDomain || "-"}</p>
            </div>
          </div>
          <TrafficChart data={charts.last24h} />
        </Card>

        <Card>
          <h2 className="font-display text-2xl font-bold tracking-[-0.04em]">{t("dashboard.topReasons")}</h2>
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

  return (
    <Card>
      <h2 className="font-display text-xl font-bold tracking-[-0.04em]">{title}</h2>
      <div className="mt-4 space-y-3">
        {items.length === 0 ? (
          <div className="rounded-2xl bg-slate-50/80 px-4 py-3 text-sm text-slate-500 dark:bg-slate-900">
            {t("messages.noDataYet")}
          </div>
        ) : null}
        {items.map((item) => (
          <div key={item.name} className="flex items-center justify-between rounded-2xl bg-slate-50/80 px-4 py-3 dark:bg-slate-900">
            <span className="truncate text-sm font-medium">{item.name}</span>
            <span className="font-mono text-sm">{formatNumber(item.value)}</span>
          </div>
        ))}
      </div>
    </Card>
  );
}
