import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { api, unwrap } from "../api/client";
import { Button } from "../components/ui/Button";
import { Card } from "../components/ui/Card";
import { Spinner } from "../components/ui/Spinner";
import { Table, TBody, TD, TH, THead, TR } from "../components/ui/Table";
import type { DomainDetail, RequestLog, StatsOverview } from "../types/api";
import { formatDate } from "../utils/format";

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

  const detail = detailQuery.data!;
  const stats = statsQuery.data!;
  const logs = logsQuery.data!;

  return (
    <div className="space-y-6">
      <Card>
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <h2 className="text-2xl font-extrabold">{detail.domain.name}</h2>
            <p className="mt-2 text-sm text-slate-500">
              {detail.domain.originProtocol}://{detail.domain.originHost}:{detail.domain.originPort} · {detail.domain.protectionMode}
            </p>
          </div>
          <div className="flex gap-3">
            <Button variant="secondary" onClick={() => sslMutation.mutate("issue")}>
              {t("actions.issueSsl")}
            </Button>
            <Button variant="secondary" onClick={() => sslMutation.mutate("renew")}>
              {t("actions.renewSsl")}
            </Button>
          </div>
        </div>
      </Card>

      <div className="grid gap-6 xl:grid-cols-[1.2fr_1fr]">
        <Card>
          <h3 className="text-lg font-bold">{t("domains.detailStats")}</h3>
          <div className="mt-4 grid gap-4 md:grid-cols-2">
            {[
              [t("statistics.incoming"), stats.incomingRequests],
              [t("statistics.allowed"), stats.allowedRequests],
              [t("statistics.blocked"), stats.blockRequests],
              [t("statistics.challenged"), stats.challengedRequests],
              [t("statistics.uniqueIps"), stats.uniqueIps],
              [t("statistics.peakRps"), stats.peakRps]
            ].map(([label, value]) => (
              <div key={String(label)} className="rounded-2xl bg-slate-50 px-4 py-4 dark:bg-slate-900">
                <div className="text-xs uppercase tracking-[0.16em] text-slate-500">{label}</div>
                <div className="mt-2 text-2xl font-bold">{value}</div>
              </div>
            ))}
          </div>
        </Card>

        <Card>
          <h3 className="text-lg font-bold">{t("domains.rules")}</h3>
          <div className="mt-4 space-y-3">
            {detail.rules.map((rule) => (
              <div key={rule.id || rule.name} className="rounded-2xl border border-slate-200 px-4 py-3 dark:border-slate-800">
                <div className="font-semibold">{rule.name}</div>
                <div className="mt-1 text-sm text-slate-500">
                  {rule.type} · {rule.action} · {rule.pattern}
                </div>
              </div>
            ))}
          </div>
        </Card>
      </div>

      <Card>
        <h3 className="text-lg font-bold">{t("domains.recentLogs")}</h3>
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
