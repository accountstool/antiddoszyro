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
      <Card className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h2 className="text-2xl font-extrabold">{t("domains.title")}</h2>
          <p className="text-sm text-slate-500">{t("domains.subtitle")}</p>
        </div>
        <div className="flex gap-3">
          <Input placeholder={t("actions.search")} value={search} onChange={(event) => setSearch(event.target.value)} />
          <Button onClick={openCreate}>{t("domains.add")}</Button>
        </div>
      </Card>

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
                  <TD className="space-x-2">
                    <Link className="text-sm font-semibold text-teal-700" to={`/domains/${item.id}`}>
                      {t("actions.view")}
                    </Link>
                    <button className="text-sm font-semibold text-slate-600" onClick={() => void openEdit(item)}>
                      {t("actions.edit")}
                    </button>
                    <button className="text-sm font-semibold text-orange-700" onClick={() => deleteMutation.mutate(item.id)}>
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
              className="min-h-48 w-full rounded-xl border border-slate-200 bg-white/80 px-3 py-2.5 text-sm outline-none focus:border-teal-500 focus:ring-2 focus:ring-teal-500/20 dark:border-slate-700 dark:bg-slate-900"
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
      <label className="mb-2 block text-sm font-semibold">{label}</label>
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
