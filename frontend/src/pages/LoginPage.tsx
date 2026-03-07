import { useState } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { useAuth } from "../app/auth";
import { Button } from "../components/ui/Button";
import { Card } from "../components/ui/Card";
import { Input } from "../components/ui/Input";

export function LoginPage() {
  const { t } = useTranslation();
  const { login } = useAuth();
  const [identifier, setIdentifier] = useState("admin@shieldpanel.local");
  const [password, setPassword] = useState("ChangeMe123!");
  const [rememberMe, setRememberMe] = useState(true);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setLoading(true);
    try {
      await login({ identifier, password, rememberMe });
    } catch (error) {
      toast.error(t("login.failed"));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="grid min-h-screen place-items-center px-4 py-12">
      <div className="grid w-full max-w-5xl gap-6 lg:grid-cols-[1.2fr_0.9fr]">
        <div className="rounded-[32px] bg-gradient-to-br from-teal-600 via-teal-700 to-cyan-900 p-8 text-white shadow-panel">
          <p className="text-xs uppercase tracking-[0.2em] text-teal-100">ShieldPanel</p>
          <h1 className="mt-6 max-w-xl text-5xl font-extrabold leading-tight">{t("brand.tagline")}</h1>
          <p className="mt-4 max-w-lg text-lg text-teal-50/90">{t("brand.description")}</p>
          <div className="mt-10 grid gap-4 md:grid-cols-3">
            {["feature.multiDomain", "feature.challenge", "feature.stats"].map((key) => (
              <div key={key} className="rounded-2xl border border-white/15 bg-white/10 p-4 text-sm">
                {t(key)}
              </div>
            ))}
          </div>
        </div>

        <Card className="p-8">
          <h2 className="text-2xl font-extrabold">{t("login.title")}</h2>
          <p className="mt-2 text-sm text-slate-500">{t("login.subtitle")}</p>

          <form className="mt-8 space-y-4" onSubmit={handleSubmit}>
            <div>
              <label className="mb-2 block text-sm font-semibold">{t("login.identifier")}</label>
              <Input value={identifier} onChange={(event) => setIdentifier(event.target.value)} />
            </div>
            <div>
              <label className="mb-2 block text-sm font-semibold">{t("login.password")}</label>
              <Input type="password" value={password} onChange={(event) => setPassword(event.target.value)} />
            </div>
            <label className="flex items-center gap-2 text-sm text-slate-600 dark:text-slate-300">
              <input checked={rememberMe} onChange={(event) => setRememberMe(event.target.checked)} type="checkbox" />
              {t("login.rememberMe")}
            </label>
            <Button block disabled={loading} type="submit">
              {loading ? t("actions.loading") : t("login.submit")}
            </Button>
          </form>
        </Card>
      </div>
    </div>
  );
}
