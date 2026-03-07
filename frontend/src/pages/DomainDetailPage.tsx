import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { api, unwrap } from "../api/client";
import { Button } from "../components/ui/Button";
import { Card } from "../components/ui/Card";
import { HeaderMetric, PageHeader } from "../components/ui/PageHeader";
import { Spinner } from "../components/ui/Spinner";
import { Table, TBody, TD, TH, THead, TR } from "../components/ui/Table";
import type { DomainDetail, RequestLog, StatsOverview } from "../types/api";
import { formatDate } from "../utils/format";

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
            <Button variant="secondary" onClick={() => sslMutation.mutate("issue")}>
              {t("actions.issueSsl")}
            </Button>
            <Button variant="secondary" onClick={() => sslMutation.mutate("renew")}>
              {t("actions.renewSsl")}
            </Button>
          </>
        }
      >
        <HeaderMetric label={t("statistics.incoming")} value={String(stats.incomingRequests)} tone="accent" />
        <HeaderMetric label={t("statistics.blocked")} value={String(stats.blockRequests)} tone="warm" />
        <HeaderMetric label={t("statistics.uniqueIps")} value={String(stats.uniqueIps)} />
      </PageHeader>

      {hasError ? (
        <Card className="border-amber-300/70 bg-amber-50/80 text-amber-900 dark:border-amber-700 dark:bg-amber-950/40 dark:text-amber-100">
          <h2 className="text-lg font-bold">{t("messages.requestFailed")}</h2>
          <p className="mt-1 text-sm opacity-80">{t("messages.domainPartial")}</p>
        </Card>
      ) : null}

      <div className="grid gap-6 xl:grid-cols-[1.2fr_1fr]">
        <Card>
          <h3 className="font-display text-2xl font-bold tracking-[-0.04em]">{t("domains.detailStats")}</h3>
          <div className="mt-4 grid gap-4 md:grid-cols-2">
            {[
              [t("statistics.incoming"), stats.incomingRequests],
              [t("statistics.allowed"), stats.allowedRequests],
              [t("statistics.blocked"), stats.blockRequests],
              [t("statistics.challenged"), stats.challengedRequests],
              [t("statistics.uniqueIps"), stats.uniqueIps],
              [t("statistics.peakRps"), stats.peakRps]
            ].map(([label, value]) => (
              <div
                key={String(label)}
                className="rounded-2xl border border-white/80 bg-white/70 px-4 py-4 dark:border-slate-800 dark:bg-slate-900/70"
              >
                <div className="text-xs uppercase tracking-[0.16em] text-slate-500">{label}</div>
                <div className="mt-2 font-display text-2xl font-bold tracking-[-0.04em]">{value}</div>
              </div>
            ))}
          </div>
        </Card>

        <Card>
          <h3 className="font-display text-2xl font-bold tracking-[-0.04em]">{t("domains.rules")}</h3>
          <div className="mt-4 space-y-3">
            {detail.rules.map((rule) => (
              <div key={rule.id || rule.name} className="rounded-2xl border border-slate-200 px-4 py-3 dark:border-slate-800">
                <div className="font-semibold">{rule.name}</div>
                <div className="mt-1 text-sm text-slate-500">
                  {rule.type} | {rule.action} | {rule.pattern}
                </div>
              </div>
            ))}
          </div>
        </Card>
      </div>

      <Card>
        <h3 className="font-display text-2xl font-bold tracking-[-0.04em]">{t("domains.recentLogs")}</h3>
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
                <TD colSpan={5} className="text-center text-slate-500">
                  {t("messages.noLogsYet")}
                </TD>
              </TR>
            ) : null}
            {logs.map((item) => (
              <TR key={item.id}>
                <TD>{formatDate(item.createdAt)}</TD>
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
