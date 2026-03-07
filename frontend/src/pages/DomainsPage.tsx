import type { ReactNode } from "react";
import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { api, unwrap } from "../api/client";
import { Button } from "../components/ui/Button";
import { Card } from "../components/ui/Card";
import { Input } from "../components/ui/Input";
import { Modal } from "../components/ui/Modal";
import { HeaderMetric, PageHeader } from "../components/ui/PageHeader";
import { Select } from "../components/ui/Select";
import { Spinner } from "../components/ui/Spinner";
import { Table, TBody, TD, TH, THead, TR } from "../components/ui/Table";
import type { Domain, DomainDetail } from "../types/api";
import { formatDate } from "../utils/format";

type DomainForm = {
  id?: string;
  name: string;
  originHost: string;
  originPort: number;
  originProtocol: "http" | "https";
  originServerName: string;
  enabled: boolean;
  protectionEnabled: boolean;
  protectionMode: string;
  challengeMode: string;
  cloudflareMode: boolean;
  sslAutoIssue: boolean;
  sslEnabled: boolean;
  forceHttps: boolean;
  rateLimitRps: number;
  rateLimitBurst: number;
  badBotMode: boolean;
  headerValidation: boolean;
  jsChallengeEnabled: boolean;
  allowedMethods: string;
  notes: string;
  rulesText: string;
};

const emptyForm: DomainForm = {
  name: "",
  originHost: "",
  originPort: 80,
  originProtocol: "http",
  originServerName: "",
  enabled: true,
  protectionEnabled: true,
  protectionMode: "basic",
  challengeMode: "cookie",
  cloudflareMode: false,
  sslAutoIssue: false,
  sslEnabled: false,
  forceHttps: false,
  rateLimitRps: 20,
  rateLimitBurst: 40,
  badBotMode: true,
  headerValidation: true,
  jsChallengeEnabled: false,
  allowedMethods: "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
  notes: "",
  rulesText: JSON.stringify(
    [
      { name: "Block traversal", type: "path", pattern: "../", action: "block", enabled: true, priority: 10 }
    ],
    null,
    2
  )
};

export function DomainsPage() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const [search, setSearch] = useState("");
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<DomainForm>(emptyForm);

  const domainsQuery = useQuery({
    queryKey: ["domains", search],
    queryFn: () => unwrap<Domain[]>(api.get(`/domains?search=${encodeURIComponent(search)}`))
  });

  const saveMutation = useMutation({
    mutationFn: async (form: DomainForm) => {
      const payload = toPayload(form);
      if (form.id) {
        return unwrap<Domain>(api.put(`/domains/${form.id}`, payload));
      }
      return unwrap<Domain>(api.post("/domains", payload));
    },
    onSuccess: () => {
      toast.success(t("messages.saved"));
      setModalOpen(false);
      setEditing(emptyForm);
      void queryClient.invalidateQueries({ queryKey: ["domains"] });
    },
    onError: () => toast.error(t("messages.requestFailed"))
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => unwrap(api.delete(`/domains/${id}`)),
    onSuccess: () => {
      toast.success(t("messages.deleted"));
      void queryClient.invalidateQueries({ queryKey: ["domains"] });
    }
  });

  const openCreate = () => {
    setEditing(emptyForm);
    setModalOpen(true);
  };

  const openEdit = async (domain: Domain) => {
    const detail = await unwrap<DomainDetail>(api.get(`/domains/${domain.id}`));
    setEditing({
      ...domain,
      allowedMethods: domain.allowedMethods.join(","),
      rulesText: JSON.stringify(
        detail.rules.map(({ name, type, pattern, action, enabled, priority }) => ({
          name,
          type,
          pattern,
          action,
          enabled,
          priority
        })),
        null,
        2
      )
    });
    setModalOpen(true);
  };

  const domains = useMemo(() => domainsQuery.data ?? [], [domainsQuery.data]);

  return (
    <div className="space-y-6">
      <PageHeader
        title={t("domains.title")}
        subtitle={t("domains.subtitle")}
        actions={
          <>
            <Input className="min-w-[220px]" placeholder={t("actions.search")} value={search} onChange={(event) => setSearch(event.target.value)} />
            <Button onClick={openCreate}>{t("domains.add")}</Button>
          </>
        }
      >
        <HeaderMetric label={t("dashboard.totalDomains")} value={String(domains.length)} tone="accent" />
        <HeaderMetric label={t("domains.sslDomains")} value={String(domains.filter((item) => item.sslEnabled).length)} />
      </PageHeader>

      <Card>
        {domainsQuery.isLoading ? (
          <div className="grid min-h-60 place-items-center">
            <Spinner />
          </div>
        ) : (
          <Table>
            <THead>
              <TR>
                <TH>{t("domains.domain")}</TH>
                <TH>{t("domains.origin")}</TH>
                <TH>{t("domains.mode")}</TH>
                <TH>{t("domains.ssl")}</TH>
                <TH>{t("domains.updatedAt")}</TH>
                <TH>{t("actions.actions")}</TH>
              </TR>
            </THead>
            <TBody>
              {domains.map((item) => (
                <TR key={item.id}>
                  <TD>
                    <div className="font-semibold">{item.name}</div>
                    <div className="text-xs text-slate-500">{item.enabled ? t("common.enabled") : t("common.disabled")}</div>
                  </TD>
                  <TD>{item.originProtocol}://{item.originHost}:{item.originPort}</TD>
                  <TD>{item.protectionMode}</TD>
                  <TD>{item.sslEnabled ? t("common.on") : t("common.off")}</TD>
                  <TD>{formatDate(item.updatedAt)}</TD>
                  <TD className="space-x-2 whitespace-nowrap">
                    <Link className="inline-flex rounded-full bg-black px-3 py-1.5 text-sm font-semibold text-white dark:bg-white dark:text-black" to={`/domains/${item.id}`}>
                      {t("actions.view")}
                    </Link>
                    <button className="inline-flex rounded-full bg-slate-100 px-3 py-1.5 text-sm font-semibold text-slate-700 dark:bg-slate-900 dark:text-slate-200" onClick={() => void openEdit(item)}>
                      {t("actions.edit")}
                    </button>
                    <button className="inline-flex rounded-full bg-neutral-800 px-3 py-1.5 text-sm font-semibold text-white dark:bg-neutral-200 dark:text-black" onClick={() => deleteMutation.mutate(item.id)}>
                      {t("actions.delete")}
                    </button>
                  </TD>
                </TR>
              ))}
            </TBody>
          </Table>
        )}
      </Card>

      <Modal open={modalOpen} title={editing.id ? t("domains.edit") : t("domains.add")} onClose={() => setModalOpen(false)}>
        <form
          className="grid gap-4 md:grid-cols-2"
          onSubmit={(event) => {
            event.preventDefault();
            saveMutation.mutate(editing);
          }}
        >
          <Field label={t("domains.domain")}>
            <Input value={editing.name} onChange={(event) => setEditing((current) => ({ ...current, name: event.target.value }))} />
          </Field>
          <Field label={t("domains.originHost")}>
            <Input value={editing.originHost} onChange={(event) => setEditing((current) => ({ ...current, originHost: event.target.value }))} />
          </Field>
          <Field label={t("domains.originPort")}>
            <Input type="number" value={editing.originPort} onChange={(event) => setEditing((current) => ({ ...current, originPort: Number(event.target.value) }))} />
          </Field>
          <Field label={t("domains.protocol")}>
            <Select value={editing.originProtocol} onChange={(event) => setEditing((current) => ({ ...current, originProtocol: event.target.value as "http" | "https" }))}>
              <option value="http">HTTP</option>
              <option value="https">HTTPS</option>
            </Select>
          </Field>
          <Field label={t("domains.protectionMode")}>
            <Select value={editing.protectionMode} onChange={(event) => setEditing((current) => ({ ...current, protectionMode: event.target.value }))}>
              <option value="off">Off</option>
              <option value="basic">Basic</option>
              <option value="aggressive">Aggressive</option>
              <option value="under_attack">Under Attack</option>
            </Select>
          </Field>
          <Field label={t("domains.challengeMode")}>
            <Select value={editing.challengeMode} onChange={(event) => setEditing((current) => ({ ...current, challengeMode: event.target.value }))}>
              <option value="off">Off</option>
              <option value="cookie">Cookie</option>
              <option value="js">JavaScript</option>
            </Select>
          </Field>
          <Field label={t("domains.rateLimitRps")}>
            <Input type="number" value={editing.rateLimitRps} onChange={(event) => setEditing((current) => ({ ...current, rateLimitRps: Number(event.target.value) }))} />
          </Field>
          <Field label={t("domains.rateLimitBurst")}>
            <Input type="number" value={editing.rateLimitBurst} onChange={(event) => setEditing((current) => ({ ...current, rateLimitBurst: Number(event.target.value) }))} />
          </Field>
          <Field className="md:col-span-2" label={t("domains.allowedMethods")}>
            <Input value={editing.allowedMethods} onChange={(event) => setEditing((current) => ({ ...current, allowedMethods: event.target.value }))} />
          </Field>
          <Field className="md:col-span-2" label={t("domains.rulesJson")}>
            <textarea
              className="min-h-48 w-full rounded-2xl border border-black/12 bg-white/78 px-4 py-3 text-sm text-slate-900 shadow-sm shadow-black/5 outline-none transition focus:border-black/40 focus:bg-white focus:ring-4 focus:ring-black/8 dark:border-white/10 dark:bg-neutral-950/82 dark:text-slate-100 dark:focus:border-white/28 dark:focus:bg-neutral-950"
              value={editing.rulesText}
              onChange={(event) => setEditing((current) => ({ ...current, rulesText: event.target.value }))}
            />
          </Field>
          <Field className="md:col-span-2" label={t("domains.notes")}>
            <Input value={editing.notes} onChange={(event) => setEditing((current) => ({ ...current, notes: event.target.value }))} />
          </Field>
          <div className="md:col-span-2 flex flex-wrap gap-4 text-sm">
            {[
              ["enabled", t("domains.enabled")],
              ["protectionEnabled", t("domains.protectionEnabled")],
              ["cloudflareMode", t("domains.cloudflareMode")],
              ["sslEnabled", t("domains.sslEnabled")],
              ["sslAutoIssue", t("domains.sslAutoIssue")],
              ["forceHttps", t("domains.forceHttps")],
              ["badBotMode", t("domains.badBotMode")],
              ["headerValidation", t("domains.headerValidation")],
              ["jsChallengeEnabled", t("domains.jsChallengeEnabled")]
            ].map(([field, label]) => (
              <label key={field} className="flex items-center gap-2">
                <input
                  className="h-4 w-4 accent-black dark:accent-white"
                  checked={Boolean(editing[field as keyof DomainForm])}
                  onChange={(event) => setEditing((current) => ({ ...current, [field]: event.target.checked }))}
                  type="checkbox"
                />
                {label}
              </label>
            ))}
          </div>
          <div className="md:col-span-2 flex justify-end gap-3">
            <Button type="button" variant="secondary" onClick={() => setModalOpen(false)}>
              {t("actions.cancel")}
            </Button>
            <Button disabled={saveMutation.isPending} type="submit">
              {t("actions.save")}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}

function Field({ label, children, className = "" }: { label: string; children: ReactNode; className?: string }) {
  return (
    <div className={className}>
      <label className="mb-2 block text-sm font-semibold text-slate-700 dark:text-slate-200">{label}</label>
      {children}
    </div>
  );
}

function toPayload(form: DomainForm) {
  return {
    ...form,
    allowedMethods: form.allowedMethods.split(",").map((value) => value.trim()).filter(Boolean),
    rules: JSON.parse(form.rulesText)
  };
}
