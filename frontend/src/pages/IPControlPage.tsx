import type { ReactNode } from "react";
import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { api, unwrap } from "../api/client";
import { Button } from "../components/ui/Button";
import { Card } from "../components/ui/Card";
import { Input } from "../components/ui/Input";
import { Modal } from "../components/ui/Modal";
import { HeaderMetric, PageHeader } from "../components/ui/PageHeader";
import { Spinner } from "../components/ui/Spinner";
import { Table, TBody, TD, TH, THead, TR } from "../components/ui/Table";
import type { IPEntry, TemporaryBan } from "../types/api";
import { formatDate } from "../utils/format";

type ListType = "blacklist" | "whitelist";

export function IPControlPage() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const [listType, setListType] = useState<ListType>("blacklist");
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState({ ip: "", cidr: "", reason: "" });

  const entriesQuery = useQuery({
    queryKey: ["ip-control", listType],
    queryFn: () => unwrap<IPEntry[]>(api.get(`/ip-control/${listType}?pageSize=100`))
  });
  const bansQuery = useQuery({
    queryKey: ["temporary-bans"],
    queryFn: () => unwrap<TemporaryBan[]>(api.get("/ip-control/bans?pageSize=100"))
  });

  const saveMutation = useMutation({
    mutationFn: () => unwrap(api.post(`/ip-control/${listType}`, form)),
    onSuccess: () => {
      toast.success(t("messages.saved"));
      setOpen(false);
      setForm({ ip: "", cidr: "", reason: "" });
      void queryClient.invalidateQueries({ queryKey: ["ip-control", listType] });
    }
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => unwrap(api.delete(`/ip-control/${id}`)),
    onSuccess: () => {
      toast.success(t("messages.deleted"));
      void queryClient.invalidateQueries({ queryKey: ["ip-control", listType] });
    }
  });

  if (entriesQuery.isLoading || bansQuery.isLoading) {
    return (
      <div className="grid min-h-[50vh] place-items-center">
        <Spinner />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title={t("ipControl.title")}
        subtitle={t("ipControl.subtitle")}
        actions={
          <>
            <Button variant={listType === "blacklist" ? "primary" : "secondary"} onClick={() => setListType("blacklist")}>
              {t("ipControl.blacklist")}
            </Button>
            <Button variant={listType === "whitelist" ? "primary" : "secondary"} onClick={() => setListType("whitelist")}>
              {t("ipControl.whitelist")}
            </Button>
            <Button onClick={() => setOpen(true)}>{t("actions.add")}</Button>
          </>
        }
      >
        <HeaderMetric label={t("ipControl.temporaryBans")} value={String((bansQuery.data ?? []).length)} tone="warm" />
        <HeaderMetric label={t("ipControl.activeEntries")} value={String((entriesQuery.data ?? []).length)} tone="accent" />
      </PageHeader>

      <Card>
        <h3 className="font-display text-2xl font-bold tracking-[-0.04em]">{t(`ipControl.${listType}`)}</h3>
        <Table>
          <THead>
            <TR>
              <TH>{t("ipControl.ip")}</TH>
              <TH>{t("ipControl.cidr")}</TH>
              <TH>{t("ipControl.reason")}</TH>
              <TH>{t("statistics.time")}</TH>
              <TH>{t("actions.actions")}</TH>
            </TR>
          </THead>
          <TBody>
            {(entriesQuery.data ?? []).map((item) => (
              <TR key={item.id}>
                <TD>{item.ip || "-"}</TD>
                <TD>{item.cidr || "-"}</TD>
                <TD>{item.reason || "-"}</TD>
                <TD>{formatDate(item.createdAt)}</TD>
                <TD>
                  <button className="inline-flex rounded-full bg-neutral-800 px-3 py-1.5 text-sm font-semibold text-white dark:bg-neutral-200 dark:text-black" onClick={() => deleteMutation.mutate(item.id)}>
                    {t("actions.delete")}
                  </button>
                </TD>
              </TR>
            ))}
          </TBody>
        </Table>
      </Card>

      <Card>
        <h3 className="font-display text-2xl font-bold tracking-[-0.04em]">{t("ipControl.temporaryBans")}</h3>
        <Table>
          <THead>
            <TR>
              <TH>{t("ipControl.ip")}</TH>
              <TH>{t("ipControl.reason")}</TH>
              <TH>{t("statistics.status")}</TH>
              <TH>{t("statistics.time")}</TH>
            </TR>
          </THead>
          <TBody>
            {(bansQuery.data ?? []).map((item) => (
              <TR key={item.id}>
                <TD>{item.ip}</TD>
                <TD>{item.reason}</TD>
                <TD>{item.source}</TD>
                <TD>{formatDate(item.expiresAt)}</TD>
              </TR>
            ))}
          </TBody>
        </Table>
      </Card>

      <Modal open={open} title={t("actions.add")} onClose={() => setOpen(false)}>
        <div className="grid gap-4">
          <Field label={t("ipControl.ip")}>
            <Input value={form.ip} onChange={(event) => setForm((current) => ({ ...current, ip: event.target.value }))} />
          </Field>
          <Field label={t("ipControl.cidr")}>
            <Input value={form.cidr} onChange={(event) => setForm((current) => ({ ...current, cidr: event.target.value }))} />
          </Field>
          <Field label={t("ipControl.reason")}>
            <Input value={form.reason} onChange={(event) => setForm((current) => ({ ...current, reason: event.target.value }))} />
          </Field>
          <div className="flex justify-end">
            <Button onClick={() => saveMutation.mutate()}>{t("actions.save")}</Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div>
      <label className="mb-2 block text-sm font-semibold text-slate-700 dark:text-slate-200">{label}</label>
      {children}
    </div>
  );
}
