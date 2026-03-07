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
import { Select } from "../components/ui/Select";
import { Spinner } from "../components/ui/Spinner";
import { Table, TBody, TD, TH, THead, TR } from "../components/ui/Table";
import type { User } from "../types/api";
import { formatDate } from "../utils/format";

type UserForm = {
  username: string;
  email: string;
  displayName: string;
  role: string;
  language: string;
  password: string;
};

const initialUser: UserForm = {
  username: "",
  email: "",
  displayName: "",
  role: "admin",
  language: "en",
  password: ""
};

export function UsersPage() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<UserForm>(initialUser);

  const usersQuery = useQuery({
    queryKey: ["users"],
    queryFn: () => unwrap<User[]>(api.get("/users?pageSize=100"))
  });

  const saveMutation = useMutation({
    mutationFn: async () => {
      if (editingId) {
        return unwrap<User>(api.put(`/users/${editingId}`, form));
      }
      return unwrap<User>(api.post("/users", form));
    },
    onSuccess: () => {
      toast.success(t("messages.saved"));
      setOpen(false);
      setEditingId(null);
      setForm(initialUser);
      void queryClient.invalidateQueries({ queryKey: ["users"] });
    }
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => unwrap(api.delete(`/users/${id}`)),
    onSuccess: () => {
      toast.success(t("messages.deleted"));
      void queryClient.invalidateQueries({ queryKey: ["users"] });
    }
  });

  if (usersQuery.isLoading) {
    return (
      <div className="grid min-h-[50vh] place-items-center">
        <Spinner />
      </div>
    );
  }

  const users = usersQuery.data ?? [];

  return (
    <div className="space-y-6">
      <PageHeader
        title={t("users.title")}
        subtitle={t("users.subtitle")}
        actions={
          <Button
            onClick={() => {
              setEditingId(null);
              setForm(initialUser);
              setOpen(true);
            }}
          >
            {t("actions.add")}
          </Button>
        }
      >
        <HeaderMetric label={t("users.totalUsers")} value={String(users.length)} tone="accent" />
      </PageHeader>

      <Card>
        <Table>
          <THead>
            <TR>
              <TH>{t("users.username")}</TH>
              <TH>{t("users.role")}</TH>
              <TH>{t("users.language")}</TH>
              <TH>{t("users.lastLogin")}</TH>
              <TH>{t("actions.actions")}</TH>
            </TR>
          </THead>
          <TBody>
            {users.map((item) => (
              <TR key={item.id}>
                <TD>
                  <div className="font-semibold">{item.displayName}</div>
                  <div className="text-xs text-slate-500">{item.email}</div>
                </TD>
                <TD>{item.role}</TD>
                <TD>{item.language}</TD>
                <TD>{formatDate(item.lastLoginAt)}</TD>
                <TD className="space-x-2 whitespace-nowrap">
                  <button
                    className="inline-flex rounded-full bg-slate-100 px-3 py-1.5 text-sm font-semibold text-slate-700 dark:bg-slate-900 dark:text-slate-200"
                    onClick={() => {
                      setEditingId(item.id);
                      setForm({
                        username: item.username,
                        email: item.email,
                        displayName: item.displayName,
                        role: item.role,
                        language: item.language,
                        password: ""
                      });
                      setOpen(true);
                    }}
                  >
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
      </Card>

      <Modal open={open} title={editingId ? t("actions.edit") : t("actions.add")} onClose={() => setOpen(false)}>
        <div className="grid gap-4 md:grid-cols-2">
          <Field label={t("users.username")}>
            <Input value={form.username} onChange={(event) => setForm((current) => ({ ...current, username: event.target.value }))} />
          </Field>
          <Field label={t("users.email")}>
            <Input value={form.email} onChange={(event) => setForm((current) => ({ ...current, email: event.target.value }))} />
          </Field>
          <Field label={t("users.displayName")}>
            <Input value={form.displayName} onChange={(event) => setForm((current) => ({ ...current, displayName: event.target.value }))} />
          </Field>
          <Field label={t("users.role")}>
            <Select value={form.role} onChange={(event) => setForm((current) => ({ ...current, role: event.target.value }))}>
              <option value="owner">Owner</option>
              <option value="admin">Admin</option>
              <option value="viewer">Viewer</option>
            </Select>
          </Field>
          <Field label={t("users.language")}>
            <Select value={form.language} onChange={(event) => setForm((current) => ({ ...current, language: event.target.value }))}>
              <option value="en">English</option>
              <option value="vi">Tiếng Việt</option>
            </Select>
          </Field>
          <Field label={t("users.password")}>
            <Input type="password" value={form.password} onChange={(event) => setForm((current) => ({ ...current, password: event.target.value }))} />
          </Field>
          <div className="md:col-span-2 flex justify-end">
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
