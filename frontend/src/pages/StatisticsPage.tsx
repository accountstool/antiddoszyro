import { useEffect, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { api, unwrap } from "../api/client";
import { Button } from "../components/ui/Button";
import { ReasonsChart } from "../components/charts/ReasonsChart";
import { TrafficChart } from "../components/charts/TrafficChart";
import { Card } from "../components/ui/Card";
import { Input } from "../components/ui/Input";
import { HeaderMetric, PageHeader } from "../components/ui/PageHeader";
import { Select } from "../components/ui/Select";
import { Spinner } from "../components/ui/Spinner";
import { Table, TBody, TD, TH, THead, TR } from "../components/ui/Table";
import type { Domain, Envelope, Pagination, RequestLog, StatsOverview } from "../types/api";
import { formatDate, formatNumber, formatPercent } from "../utils/format";

const emptyStats: StatsOverview = {
  incomingRequests: 0,
  allowedRequests: 0,
  blockRequests: 0,
  challengedRequests: 0,
  challengePassRate: 0,
  uniqueIps: 0,
  peakRps: 0,
  peakTime: "",
  topIps: [],
  topUserAgents: [],
  topDomains: [],
  topReasons: [],
  topCountries: [],
  requestSeries: []
};

export function StatisticsPage() {
  const { t } = useTranslation();
  const [domainId, setDomainId] = useState("");
  const [decision, setDecision] = useState("");
  const [reason, setReason] = useState("");
  const [page, setPage] = useState(1);
  const pageSize = 10;

  useEffect(() => {
    setPage(1);
  }, [domainId, decision, reason]);

  const domainsQuery = useQuery({
    queryKey: ["domains", "options"],
    queryFn: () => unwrap<Domain[]>(api.get("/domains?pageSize=100"))
  });

  const statsQuery = useQuery({
    queryKey: ["statistics", domainId, decision, reason],
    queryFn: () =>
      unwrap<StatsOverview>(
        api.get(`/statistics/overview?domainId=${domainId}&decision=${decision}&reason=${encodeURIComponent(reason)}`)
      )
  });

  const logsQuery = useQuery({
    queryKey: ["logs", domainId, decision, reason, page],
    queryFn: async () => {
      const response = await api.get<Envelope<RequestLog[]>>(
        `/logs?page=${page}&pageSize=${pageSize}&domainId=${domainId}&decision=${decision}&reason=${encodeURIComponent(reason)}`
      );
      return response.data;
    }
  });

  if (statsQuery.isLoading || logsQuery.isLoading || domainsQuery.isLoading) {
    return (
      <div className="grid min-h-[50vh] place-items-center">
        <Spinner />
      </div>
    );
  }

  const stats = statsQuery.data ?? emptyStats;
  const logs = logsQuery.data?.data ?? [];
  const pagination = logsQuery.data?.pagination;
  const hasError = statsQuery.isError || logsQuery.isError || domainsQuery.isError;

  return (
    <div className="space-y-6">
      <PageHeader title={t("nav.statistics")} subtitle={t("statistics.requestsOverTime")}>
        <HeaderMetric label={t("statistics.incoming")} value={formatNumber(stats.incomingRequests)} tone="accent" />
        <HeaderMetric label={t("statistics.allowed")} value={formatNumber(stats.allowedRequests)} tone="warm" />
        <HeaderMetric label={t("statistics.uniqueIps")} value={formatNumber(stats.uniqueIps)} />
      </PageHeader>

      {hasError ? (
        <Card className="border-black/12 bg-neutral-100/80 text-neutral-900 dark:border-white/10 dark:bg-neutral-900/80 dark:text-neutral-100">
          <h2 className="text-lg font-bold">{t("messages.requestFailed")}</h2>
          <p className="mt-1 text-sm opacity-80">{t("messages.statsUnavailable")}</p>
        </Card>
      ) : null}

      <Card className="overflow-hidden p-0">
        <div className="border-b border-black/10 px-6 py-5 dark:border-white/10">
          <h2 className="font-display text-2xl font-bold tracking-[-0.04em] text-slate-950 dark:text-slate-50">{t("statistics.filters")}</h2>
          <p className="mt-1 text-sm text-slate-600 dark:text-slate-300">{t("statistics.filtersDescription")}</p>
        </div>
        <div className="grid gap-4 p-6 md:grid-cols-3">
          <Select value={domainId} onChange={(event) => setDomainId(event.target.value)}>
            <option value="">{t("statistics.allDomains")}</option>
            {(domainsQuery.data ?? []).map((domain) => (
              <option key={domain.id} value={domain.id}>
                {domain.name}
              </option>
            ))}
          </Select>
          <Select value={decision} onChange={(event) => setDecision(event.target.value)}>
            <option value="">{t("statistics.allDecisions")}</option>
            <option value="allowed">{t("statistics.allowed")}</option>
            <option value="blocked">{t("statistics.blocked")}</option>
            <option value="challenged">{t("statistics.challenged")}</option>
          </Select>
          <Input placeholder={t("statistics.reason")} value={reason} onChange={(event) => setReason(event.target.value)} />
        </div>
      </Card>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-5">
        {[
          [t("statistics.incoming"), formatNumber(stats.incomingRequests)],
          [t("statistics.allowed"), formatNumber(stats.allowedRequests)],
          [t("statistics.blocked"), formatNumber(stats.blockRequests)],
          [t("statistics.challenged"), formatNumber(stats.challengedRequests)],
          [t("statistics.challengePassRate"), formatPercent(stats.challengePassRate)]
        ].map(([label, value], index) => (
          <Card
            key={String(label)}
            className={
              index === 0
                ? "border-black bg-black text-white dark:border-white dark:bg-white dark:text-black"
                : index === 4
                  ? "border-neutral-300 bg-neutral-100 text-black dark:border-neutral-700 dark:bg-neutral-900 dark:text-white"
                  : ""
            }
          >
            <p
              className={`text-xs uppercase tracking-[0.16em] ${
                index === 0 ? "text-white/68 dark:text-black/60" : "text-slate-500 dark:text-slate-400"
              }`}
            >
              {label}
            </p>
            <p className={`mt-2 text-3xl font-extrabold ${index === 0 ? "text-white dark:text-black" : "text-slate-950 dark:text-slate-50"}`}>{value}</p>
          </Card>
        ))}
      </div>

      <div className="grid gap-6 xl:grid-cols-[1.6fr_1fr]">
        <Card>
          <div className="mb-4 flex flex-wrap items-end justify-between gap-3">
            <div>
              <h2 className="font-display text-2xl font-bold tracking-[-0.04em] text-slate-950 dark:text-slate-50">{t("statistics.requestsOverTime")}</h2>
              <p className="mt-1 text-sm text-slate-600 dark:text-slate-300">{t("statistics.chartDescription")}</p>
            </div>
            <div className="rounded-full border border-black/10 bg-neutral-50 px-3.5 py-2 dark:border-white/10 dark:bg-neutral-900">
              <div className="text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500 dark:text-slate-400">{t("statistics.peakRps")}</div>
              <div className="mt-1 text-sm font-semibold text-slate-950 dark:text-slate-50">{formatNumber(stats.peakRps)}</div>
            </div>
          </div>
          <TrafficChart data={stats.requestSeries} />
        </Card>
        <Card>
          <div className="mb-4">
            <h2 className="font-display text-2xl font-bold tracking-[-0.04em] text-slate-950 dark:text-slate-50">{t("statistics.topReasons")}</h2>
            <p className="mt-1 text-sm text-slate-600 dark:text-slate-300">{t("statistics.reasonBreakdown")}</p>
          </div>
          <ReasonsChart data={stats.topReasons} />
        </Card>
      </div>

      <Card className="overflow-hidden p-0">
        <div className="flex flex-col gap-3 border-b border-black/10 px-6 py-5 dark:border-white/10 md:flex-row md:items-center md:justify-between">
          <div>
            <h2 className="font-display text-2xl font-bold tracking-[-0.04em] text-slate-950 dark:text-slate-50">{t("statistics.logs")}</h2>
            <p className="mt-1 text-sm text-slate-600 dark:text-slate-300">{paginationSummary(pagination, logs.length, t)}</p>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="secondary" disabled={!pagination || pagination.page <= 1} onClick={() => setPage((current) => Math.max(1, current - 1))}>
              {t("statistics.previousPage")}
            </Button>
            <Button
              variant="secondary"
              disabled={!pagination || pagination.page >= pagination.totalPages}
              onClick={() => setPage((current) => current + 1)}
            >
              {t("statistics.nextPage")}
            </Button>
          </div>
        </div>
        <div className="p-6 pt-5">
          <Table>
            <THead>
              <TR>
                <TH>{t("statistics.time")}</TH>
                <TH>{t("statistics.domain")}</TH>
                <TH>{t("statistics.ip")}</TH>
                <TH>{t("statistics.path")}</TH>
                <TH>{t("statistics.status")}</TH>
                <TH>{t("statistics.reason")}</TH>
              </TR>
            </THead>
            <TBody>
              {logs.length === 0 ? (
                <TR>
                  <TD colSpan={6} className="py-8 text-center text-slate-500 dark:text-slate-400">
                    {t("messages.noLogsYet")}
                  </TD>
                </TR>
              ) : null}
              {logs.map((item) => (
                <TR key={item.id}>
                  <TD>{formatDate(item.createdAt)}</TD>
                  <TD className="font-medium text-slate-950 dark:text-slate-50">{item.domainName}</TD>
                  <TD className="font-mono text-[13px]">{item.clientIp}</TD>
                  <TD className="max-w-[320px] truncate font-mono text-[13px]">{item.path}</TD>
                  <TD>
                    <span className="inline-flex rounded-full border border-black/10 bg-neutral-100 px-2.5 py-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-slate-700 dark:border-white/10 dark:bg-neutral-900 dark:text-slate-200">
                      {item.decision}
                    </span>
                  </TD>
                  <TD>{item.blockReason || "-"}</TD>
                </TR>
              ))}
            </TBody>
          </Table>
        </div>
        <div className="flex flex-col gap-3 border-t border-black/10 px-6 py-4 dark:border-white/10 md:flex-row md:items-center md:justify-between">
          <div className="text-sm text-slate-600 dark:text-slate-300">
            {t("statistics.pageOf", {
              page: pagination?.page ?? 1,
              totalPages: Math.max(1, pagination?.totalPages ?? 1)
            })}
          </div>
          <div className="flex items-center gap-2">
            <Button variant="secondary" disabled={!pagination || pagination.page <= 1} onClick={() => setPage((current) => Math.max(1, current - 1))}>
              {t("statistics.previousPage")}
            </Button>
            <Button
              variant="secondary"
              disabled={!pagination || pagination.page >= pagination.totalPages}
              onClick={() => setPage((current) => current + 1)}
            >
              {t("statistics.nextPage")}
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
}

function paginationSummary(
  pagination: Pagination | undefined,
  currentItems: number,
  t: (key: string, options?: Record<string, unknown>) => string
) {
  if (!pagination) {
    return t("statistics.showingRange", { start: 0, end: currentItems, total: currentItems });
  }

  const start = pagination.totalItems === 0 ? 0 : (pagination.page - 1) * pagination.pageSize + 1;
  const end = pagination.totalItems === 0 ? 0 : Math.min(pagination.page * pagination.pageSize, pagination.totalItems);
  return t("statistics.showingRange", { start, end, total: pagination.totalItems });
}
