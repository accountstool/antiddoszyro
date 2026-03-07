import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "./Button";

export function Modal({
  open,
  title,
  onClose,
  children
}: PropsWithChildren<{ open: boolean; title: string; onClose: () => void }>) {
  const { t } = useTranslation();
  if (!open) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-50 grid place-items-center bg-slate-950/50 p-4">
      <div className="w-full max-w-3xl rounded-3xl border border-slate-200 bg-white p-6 shadow-panel dark:border-slate-800 dark:bg-slate-950">
        <div className="mb-5 flex items-center justify-between">
          <h2 className="text-xl font-bold">{title}</h2>
          <Button variant="ghost" onClick={onClose}>
            {t("actions.close")}
          </Button>
        </div>
        {children}
      </div>
    </div>
  );
}
