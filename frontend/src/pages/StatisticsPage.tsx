import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { api, unwrap } from "../api/client";
import { ReasonsChart } from "../components/charts/ReasonsChart";
import { TrafficChart } from "../components/charts/TrafficChart";
import { Card } from "../components/ui/Card";
import { Input } from "../components/ui/Input";
import { HeaderMetric, PageHeader } from "../components/ui/PageHeader";
import { Select } from "../components/ui/Select";
import { Spinner } from "../components/ui/Spinner";
import { Table, TBody, TD, TH, THead, TR } from "../components/ui/Table";
import type { Domain, RequestLog, StatsOverview } from "../types/api";
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
    queryKey: ["logs", domainId, decision, reason],
    queryFn: () =>
      unwrap<RequestLog[]>(api.get(`/logs?pageSize=20&domainId=${domainId}&decision=${decision}&reason=${encodeURIComponent(reason)}`))
  });

  if (statsQuery.isLoading || logsQuery.isLoading || domainsQuery.isLoading) {
    return (
      <div className="grid min-h-[50vh] place-items-center">
        <Spinner />
      </div>
    );
  }

  const stats = statsQuery.data ?? emptyStats;
  const logs = logsQuery.data ?? [];
  const hasError = statsQuery.isError || logsQuery.isError || domainsQuery.isError;

  return (
    <div className="space-y-6">
      <PageHeader title={t("nav.statistics")} subtitle={t("statistics.requestsOverTime")}>
        <HeaderMetric label={t("statistics.incoming")} value={formatNumber(stats.incomingRequests)} tone="accent" />
        <HeaderMetric label={t("statistics.blocked")} value={formatNumber(stats.blockRequests)} tone="warm" />
        <HeaderMetric label={t("statistics.uniqueIps")} value={formatNumber(stats.uniqueIps)} />
      </PageHeader>

      {hasError ? (
        <Card className="border-amber-300/70 bg-amber-50/80 text-amber-900 dark:border-amber-700 dark:bg-amber-950/40 dark:text-amber-100">
          <h2 className="text-lg font-bold">{t("messages.requestFailed")}</h2>
          <p className="mt-1 text-sm opacity-80">{t("messages.statsUnavailable")}</p>
        </Card>
      ) : null}

      <Card className="grid gap-4 md:grid-cols-4">
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
          <option value="allowed">Allowed</option>
          <option value="blocked">Blocked</option>
          <option value="challenged">Challenged</option>
        </Select>
        <Input placeholder={t("statistics.reason")} value={reason} onChange={(event) => setReason(event.target.value)} />
      </Card>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-5">
        {[
          [t("statistics.incoming"), formatNumber(stats.incomingRequests)],
          [t("statistics.allowed"), formatNumber(stats.allowedRequests)],
          [t("statistics.blocked"), formatNumber(stats.blockRequests)],
          [t("statistics.challenged"), formatNumber(stats.challengedRequests)],
          [t("statistics.challengePassRate"), formatPercent(stats.challengePassRate)]
        ].map(([label, value]) => (
          <Card key={String(label)}>
            <p className="text-xs uppercase tracking-[0.16em] text-slate-500">{label}</p>
            <p className="mt-2 text-3xl font-extrabold">{value}</p>
          </Card>
        ))}
      </div>

      <div className="grid gap-6 xl:grid-cols-[1.6fr_1fr]">
        <Card>
          <h2 className="font-display text-2xl font-bold tracking-[-0.04em]">{t("statistics.requestsOverTime")}</h2>
          <TrafficChart data={stats.requestSeries} />
        </Card>
        <Card>
          <h2 className="font-display text-2xl font-bold tracking-[-0.04em]">{t("statistics.topReasons")}</h2>
          <ReasonsChart data={stats.topReasons} />
        </Card>
      </div>

      <Card>
        <h2 className="font-display text-2xl font-bold tracking-[-0.04em]">{t("statistics.logs")}</h2>
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
                <TD colSpan={6} className="text-center text-slate-500">
                  {t("messages.noLogsYet")}
                </TD>
              </TR>
            ) : null}
            {logs.map((item) => (
              <TR key={item.id}>
                <TD>{formatDate(item.createdAt)}</TD>
                <TD>{item.domainName}</TD>
                <TD>{item.clientIp}</TD>
                <TD>{item.path}</TD>
                <TD>{item.decision}</TD>
                <TD>{item.blockReason || "-"}</TD>
              </TR>
            ))}
          </TBody>
        </Table>
      </Card>
    </div>
  );
}
