import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { api, unwrap } from "../api/client";
import { Button } from "../components/ui/Button";
import { Card } from "../components/ui/Card";
import { Input } from "../components/ui/Input";
import { HeaderMetric, PageHeader } from "../components/ui/PageHeader";
import { Spinner } from "../components/ui/Spinner";
import type { SystemSetting } from "../types/api";

export function SettingsPage() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const [draft, setDraft] = useState<Record<string, string>>({});

  const settingsQuery = useQuery({
    queryKey: ["settings"],
    queryFn: () => unwrap<SystemSetting[]>(api.get("/settings"))
  });

  const saveMutation = useMutation({
    mutationFn: () => unwrap(api.put("/settings", { values: draft })),
    onSuccess: () => {
      toast.success(t("messages.saved"));
      void queryClient.invalidateQueries({ queryKey: ["settings"] });
    }
  });

  const settings = useMemo(() => settingsQuery.data ?? [], [settingsQuery.data]);

  if (settingsQuery.isLoading) {
    return (
      <div className="grid min-h-[50vh] place-items-center">
        <Spinner />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title={t("settings.title")}
        subtitle={t("settings.subtitle")}
        actions={<Button onClick={() => saveMutation.mutate()}>{t("actions.save")}</Button>}
      >
        <HeaderMetric label={t("settings.keysCount")} value={String(settings.length)} tone="accent" />
      </PageHeader>

      <Card>
        <div className="grid gap-4">
          {settings.map((item) => (
            <div key={item.key} className="grid gap-2 md:grid-cols-[220px_1fr] md:items-center">
              <div>
                <div className="font-semibold">{item.key}</div>
                <div className="text-xs uppercase tracking-[0.16em] text-slate-500">{item.type}</div>
              </div>
              <Input
                value={draft[item.key] ?? item.value}
                onChange={(event) => setDraft((current) => ({ ...current, [item.key]: event.target.value }))}
              />
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
}
