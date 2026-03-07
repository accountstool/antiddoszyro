import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import clsx from "clsx";

import { api, unwrap } from "../api/client";
import { Button } from "../components/ui/Button";
import { Card } from "../components/ui/Card";
import { HeaderMetric, PageHeader } from "../components/ui/PageHeader";
import { Spinner } from "../components/ui/Spinner";
import { Table, TBody, TD, TH, THead, TR } from "../components/ui/Table";
import type { DomainDetail, RequestLog, StatsOverview } from "../types/api";
import { formatDate, formatNumber, formatPercent } from "../utils/format";

const emptyDomainDetail: DomainDetail = {
  domain: {
    id: "",
    name: "",
    originHost: "",
    originPort: 0,
    originProtocol: "http",
    originServerName: "",
    enabled: false,
    protectionEnabled: false,
    protectionMode: "off",
    challengeMode: "off",
    cloudflareMode: false,
    sslAutoIssue: false,
    sslEnabled: false,
    forceHttps: false,
    rateLimitRps: 0,
    rateLimitBurst: 0,
    badBotMode: false,
    headerValidation: false,
    jsChallengeEnabled: false,
    allowedMethods: [],
    notes: "",
    createdAt: "",
    updatedAt: ""
  },
  rules: [],
  nginxStatus: "unknown"
};

const emptyDomainStats: StatsOverview = {
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

export function DomainDetailPage() {
  const { t } = useTranslation();
  const { id } = useParams();
  const queryClient = useQueryClient();

  const detailQuery = useQuery({
    queryKey: ["domain", id],
    queryFn: () => unwrap<DomainDetail>(api.get(`/domains/${id}`)),
    enabled: Boolean(id)
  });
  const statsQuery = useQuery({
    queryKey: ["domain", id, "stats"],
    queryFn: () => unwrap<StatsOverview>(api.get(`/domains/${id}/stats`)),
    enabled: Boolean(id)
  });
  const logsQuery = useQuery({
    queryKey: ["domain", id, "logs"],
    queryFn: () => unwrap<RequestLog[]>(api.get(`/domains/${id}/logs?pageSize=10`)),
    enabled: Boolean(id)
  });
  const sslMutation = useMutation({
    mutationFn: (mode: "issue" | "renew") => unwrap(api.post(`/ssl/${mode}`, { domainId: id })),
    onSuccess: () => {
      toast.success(t("messages.saved"));
      void queryClient.invalidateQueries({ queryKey: ["domain", id] });
    }
  });

  if (detailQuery.isLoading || statsQuery.isLoading || logsQuery.isLoading) {
    return (
      <div className="grid min-h-[50vh] place-items-center">
        <Spinner />
      </div>
    );
  }

  const detail = detailQuery.data ?? emptyDomainDetail;
  const stats = statsQuery.data ?? emptyDomainStats;
  const logs = logsQuery.data ?? [];
  const hasError = detailQuery.isError || statsQuery.isError || logsQuery.isError;
  const overviewItems = [
    { label: t("domains.origin"), value: `${detail.domain.originHost || "-"}:${detail.domain.originPort || "-"}` },
    { label: t("domains.protocol"), value: detail.domain.originProtocol ? detail.domain.originProtocol.toUpperCase() : "-" },
    { label: "Nginx", value: detail.nginxStatus || "-" },
    { label: t("domains.updatedAt"), value: formatDate(detail.domain.updatedAt) },
    { label: t("domains.allowedMethods"), value: detail.domain.allowedMethods.length > 0 ? detail.domain.allowedMethods.join(", ") : "-" },
    { label: t("statistics.challengePassRate"), value: formatPercent(stats.challengePassRate) }
  ];
  const statusPills = [
    `${t("domains.mode")}: ${detail.domain.protectionMode || "-"}`,
    `${t("domains.challengeMode")}: ${detail.domain.challengeMode || "-"}`,
    `${t("domains.cloudflareMode")}: ${detail.domain.cloudflareMode ? t("common.on") : t("common.off")}`,
    `${t("domains.ssl")}: ${detail.domain.sslEnabled ? t("common.on") : t("common.off")}`,
    `${t("domains.forceHttps")}: ${detail.domain.forceHttps ? t("common.on") : t("common.off")}`,
    `${t("domains.badBotMode")}: ${detail.domain.badBotMode ? t("common.on") : t("common.off")}`
  ];

  return (
    <div className="space-y-6">
      <PageHeader
        title={detail.domain.name || t("domains.domain")}
        subtitle={
          detail.domain.name
            ? `${detail.domain.originProtocol}://${detail.domain.originHost}:${detail.domain.originPort} | ${detail.domain.protectionMode}`
            : t("domains.subtitle")
        }
        actions={
          <>
            <Button onClick={() => sslMutation.mutate("issue")}>
              {t("actions.issueSsl")}
            </Button>
            <Button variant="secondary" onClick={() => sslMutation.mutate("renew")}>
              {t("actions.renewSsl")}
            </Button>
          </>
        }
      >
        <HeaderMetric label={t("statistics.incoming")} value={formatNumber(stats.incomingRequests)} tone="accent" />
        <HeaderMetric label={t("statistics.blocked")} value={formatNumber(stats.blockRequests)} tone="warm" />
        <HeaderMetric label={t("statistics.uniqueIps")} value={formatNumber(stats.uniqueIps)} />
      </PageHeader>

      {hasError ? (
        <Card className="border-neutral-300 bg-neutral-100 text-neutral-900 dark:border-white/10 dark:bg-neutral-900 dark:text-neutral-100">
          <h2 className="text-lg font-bold">{t("messages.requestFailed")}</h2>
          <p className="mt-1 text-sm opacity-80">{t("messages.domainPartial")}</p>
        </Card>
      ) : null}

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1.35fr)_360px]">
        <Card>
          <div className="flex flex-col gap-6">
            <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
              <div>
                <div className="text-[11px] uppercase tracking-[0.22em] text-slate-500 dark:text-slate-400">{t("domains.detailStats")}</div>
                <h3 className="mt-2 font-display text-2xl font-bold tracking-[-0.04em]">{detail.domain.name || t("domains.domain")}</h3>
                <p className="mt-2 text-sm text-slate-500 dark:text-slate-400">
                  {detail.domain.originProtocol}://{detail.domain.originHost}:{detail.domain.originPort}
                </p>
              </div>
              <span className="inline-flex items-center rounded-full border border-black/10 bg-neutral-100 px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.18em] text-slate-700 dark:border-white/10 dark:bg-neutral-900 dark:text-slate-200">
                {detail.nginxStatus}
              </span>
            </div>

            <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
              {[
                [t("statistics.incoming"), stats.incomingRequests],
                [t("statistics.allowed"), stats.allowedRequests],
                [t("statistics.blocked"), stats.blockRequests],
                [t("statistics.challenged"), stats.challengedRequests],
                [t("statistics.uniqueIps"), stats.uniqueIps],
                [t("statistics.peakRps"), stats.peakRps]
              ].map(([label, value], index) => (
                <div
                  key={String(label)}
                  className={clsx(
                    "rounded-2xl border px-4 py-4",
                    index === 0
                      ? "border-black bg-black text-white dark:border-white dark:bg-white dark:text-black"
                      : "border-black/10 bg-neutral-50 text-slate-950 dark:border-white/10 dark:bg-neutral-900 dark:text-slate-50"
                  )}
                >
                  <div className={clsx("text-[11px] uppercase tracking-[0.18em]", index === 0 ? "text-white/65 dark:text-black/60" : "text-slate-500 dark:text-slate-400")}>
                    {label}
                  </div>
                  <div className="mt-3 font-display text-[2rem] font-bold leading-none tracking-[-0.05em]">{formatNumber(Number(value) || 0)}</div>
                </div>
              ))}
            </div>

            <div className="grid gap-3 sm:grid-cols-2">
              {overviewItems.map((item) => (
                <div key={item.label} className="rounded-2xl border border-black/10 bg-white px-4 py-3 dark:border-white/10 dark:bg-neutral-950">
                  <div className="text-[11px] uppercase tracking-[0.18em] text-slate-500 dark:text-slate-400">{item.label}</div>
                  <div className="mt-2 break-words text-sm font-semibold text-slate-900 dark:text-slate-100">{item.value}</div>
                </div>
              ))}
            </div>
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between gap-4">
            <div>
              <div className="text-[11px] uppercase tracking-[0.22em] text-slate-500 dark:text-slate-400">{t("domains.mode")}</div>
              <h3 className="mt-2 font-display text-2xl font-bold tracking-[-0.04em]">{t("domains.rules")}</h3>
            </div>
            <span className="inline-flex items-center rounded-full border border-black/10 bg-neutral-100 px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.18em] text-slate-700 dark:border-white/10 dark:bg-neutral-900 dark:text-slate-200">
              {detail.rules.length}
            </span>
          </div>

          <div className="mt-5 flex flex-wrap gap-2">
            {statusPills.map((item) => (
              <span
                key={item}
                className="inline-flex items-center rounded-full border border-black/10 bg-neutral-100 px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-slate-700 dark:border-white/10 dark:bg-neutral-900 dark:text-slate-200"
              >
                {item}
              </span>
            ))}
          </div>

          <div className="mt-5 space-y-3">
            {detail.rules.length === 0 ? (
              <div className="rounded-2xl border border-dashed border-black/12 bg-neutral-50 px-4 py-5 text-sm text-slate-500 dark:border-white/10 dark:bg-neutral-900 dark:text-slate-400">
                {t("messages.noDataYet")}
              </div>
            ) : null}
            {[
              ...detail.rules.map((rule) => (
                <div key={rule.id || rule.name} className="rounded-2xl border border-black/10 bg-white px-4 py-4 dark:border-white/10 dark:bg-neutral-950">
                  <div className="flex items-start justify-between gap-3">
                    <div className="font-semibold text-slate-950 dark:text-slate-50">{rule.name}</div>
                    <span className="inline-flex items-center rounded-full border border-black/10 bg-black px-2.5 py-1 text-[10px] font-semibold uppercase tracking-[0.16em] text-white dark:border-white/10 dark:bg-white dark:text-black">
                      {rule.action}
                    </span>
                  </div>
                  <div className="mt-2 text-sm text-slate-500 dark:text-slate-400">
                    {rule.type} | {rule.pattern}
                  </div>
                </div>
              ))
            ]}
          </div>
        </Card>
      </div>

      <Card className="overflow-hidden p-0">
        <div className="flex flex-col gap-2 border-b border-black/10 px-6 py-5 dark:border-white/10">
          <h3 className="font-display text-2xl font-bold tracking-[-0.04em]">{t("domains.recentLogs")}</h3>
          <p className="text-sm text-slate-500 dark:text-slate-400">{formatNumber(logs.length)} events</p>
        </div>
        <div className="p-6 pt-5">
          <Table>
            <THead>
              <TR>
                <TH>{t("statistics.time")}</TH>
                <TH>{t("statistics.ip")}</TH>
                <TH>{t("statistics.path")}</TH>
                <TH>{t("statistics.status")}</TH>
                <TH>{t("statistics.reason")}</TH>
              </TR>
            </THead>
            <TBody>
              {logs.length === 0 ? (
                <TR>
                  <TD colSpan={5} className="py-8 text-center text-slate-500">
                    {t("messages.noLogsYet")}
                  </TD>
                </TR>
              ) : null}
              {logs.map((item) => (
                <TR key={item.id}>
                  <TD>{formatDate(item.createdAt)}</TD>
                  <TD className="font-mono text-[13px]">{item.clientIp}</TD>
                  <TD className="font-mono text-[13px]">{item.path}</TD>
                  <TD>
                    <span className={decisionBadgeClass(item.decision)}>{item.decision}</span>
                  </TD>
                  <TD className="font-mono text-[13px] text-slate-500 dark:text-slate-400">{item.blockReason || "-"}</TD>
                </TR>
              ))}
            </TBody>
          </Table>
        </div>
      </Card>
    </div>
  );
}

function decisionBadgeClass(decision: string) {
  if (decision === "allowed") {
    return "inline-flex items-center rounded-full border border-black/10 bg-neutral-100 px-2.5 py-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-slate-700 dark:border-white/10 dark:bg-neutral-900 dark:text-slate-200";
  }
  if (decision === "blocked") {
    return "inline-flex items-center rounded-full border border-black bg-black px-2.5 py-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-white dark:border-white dark:bg-white dark:text-black";
  }
  if (decision === "challenge_passed") {
    return "inline-flex items-center rounded-full border border-neutral-300 bg-white px-2.5 py-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-slate-700 dark:border-white/10 dark:bg-neutral-950 dark:text-slate-200";
  }
  return "inline-flex items-center rounded-full border border-black/10 bg-neutral-50 px-2.5 py-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-slate-700 dark:border-white/10 dark:bg-neutral-900 dark:text-slate-200";
}
