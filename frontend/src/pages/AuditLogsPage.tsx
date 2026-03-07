import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { api, unwrap } from "../api/client";
import { Card } from "../components/ui/Card";
import { Spinner } from "../components/ui/Spinner";
import { Table, TBody, TD, TH, THead, TR } from "../components/ui/Table";
import type { AuditLog } from "../types/api";
import { formatDate } from "../utils/format";

export function AuditLogsPage() {
  const { t } = useTranslation();
  const auditQuery = useQuery({
    queryKey: ["audit-logs"],
    queryFn: () => unwrap<AuditLog[]>(api.get("/audit-logs?pageSize=100"))
  });

  if (auditQuery.isLoading) {
    return (
      <div className="grid min-h-[50vh] place-items-center">
        <Spinner />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <Card>
        <h2 className="text-2xl font-extrabold">{t("audit.title")}</h2>
        <p className="text-sm text-slate-500">{t("audit.subtitle")}</p>
      </Card>

      <Card>
        <Table>
          <THead>
            <TR>
              <TH>{t("statistics.time")}</TH>
              <TH>{t("audit.user")}</TH>
              <TH>{t("audit.action")}</TH>
              <TH>{t("audit.target")}</TH>
              <TH>{t("statistics.ip")}</TH>
              <TH>{t("audit.details")}</TH>
            </TR>
          </THead>
          <TBody>
            {(auditQuery.data ?? []).map((item) => (
              <TR key={item.id}>
                <TD>{formatDate(item.createdAt)}</TD>
                <TD>{item.username}</TD>
                <TD>{item.action}</TD>
                <TD>{item.entityType}:{item.entityId}</TD>
                <TD>{item.ipAddress}</TD>
                <TD className="max-w-lg truncate">{item.details || "-"}</TD>
              </TR>
            ))}
          </TBody>
        </Table>
      </Card>
    </div>
  );
}
