import { useState } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { Globe, Moon, Shield, Sun, Zap } from "lucide-react";

import { useAuth } from "../app/auth";
import { useUI } from "../app/ui";
import { Button } from "../components/ui/Button";
import { Card } from "../components/ui/Card";
import { Input } from "../components/ui/Input";
import { Select } from "../components/ui/Select";

export function LoginPage() {
  const { t } = useTranslation();
  const { login } = useAuth();
  const { language, setLanguage, theme, toggleTheme } = useUI();
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
    <div className="relative min-h-screen overflow-hidden px-4 py-8">
      <div className="pointer-events-none absolute inset-0 panel-grid-overlay opacity-35" />

      <div className="mx-auto flex w-full max-w-7xl justify-end gap-3 pb-5">
        <Select className="w-[154px]" value={language} onChange={(event) => setLanguage(event.target.value)}>
          <option value="en">English</option>
          <option value="vi">Tieng Viet</option>
        </Select>
        <Button variant="secondary" onClick={toggleTheme}>
          {theme === "dark" ? <Sun size={16} /> : <Moon size={16} />}
        </Button>
      </div>

      <div className="mx-auto grid w-full max-w-7xl gap-6 lg:grid-cols-[1.15fr_0.85fr]">
        <div className="relative overflow-hidden rounded-[36px] bg-mesh-teal p-8 text-white shadow-panel lg:p-10">
          <div className="absolute inset-0 panel-grid-overlay opacity-15" />
          <div className="relative">
            <div className="inline-flex items-center gap-2 rounded-full border border-white/20 bg-white/6 px-4 py-2 text-[11px] uppercase tracking-[0.22em] text-white/70">
              <Shield size={14} />
              ShieldPanel
            </div>

            <h1 className="mt-8 max-w-2xl font-display text-5xl font-bold leading-[0.95] tracking-[-0.06em] text-balance lg:text-7xl">
              {t("brand.tagline")}
            </h1>
            <p className="mt-5 max-w-xl text-base leading-7 text-white/82 lg:text-lg">{t("brand.description")}</p>

            <div className="mt-10 grid gap-4 md:grid-cols-3">
              {[
                { key: "feature.multiDomain", icon: Globe },
                { key: "feature.challenge", icon: Zap },
                { key: "feature.stats", icon: Shield }
              ].map((item) => {
                const Icon = item.icon;
                return (
                    <div key={item.key} className="rounded-[24px] border border-white/12 bg-white/6 p-5 backdrop-blur">
                    <div className="grid h-10 w-10 place-items-center rounded-2xl bg-white/10">
                      <Icon size={18} />
                    </div>
                    <p className="mt-4 text-sm leading-6 text-white/85">{t(item.key)}</p>
                  </div>
                );
              })}
            </div>

            <div className="mt-8 grid gap-4 md:grid-cols-2">
              <div className="rounded-[24px] border border-white/10 bg-white/6 p-5 backdrop-blur">
                <div className="text-[11px] uppercase tracking-[0.22em] text-white/65">{t("nav.dashboard")}</div>
                <div className="mt-3 font-display text-3xl font-bold tracking-[-0.05em]">24/7</div>
                <div className="mt-2 text-sm text-white/75">{t("feature.challenge")}</div>
              </div>
              <div className="rounded-[24px] border border-white/10 bg-white/6 p-5 backdrop-blur">
                <div className="text-[11px] uppercase tracking-[0.22em] text-white/65">{t("nav.domains")}</div>
                <div className="mt-3 font-display text-3xl font-bold tracking-[-0.05em]">Nginx</div>
                <div className="mt-2 text-sm text-white/75">{t("feature.multiDomain")}</div>
              </div>
            </div>
          </div>
        </div>

        <Card className="self-center p-8 lg:p-9">
          <div className="inline-flex items-center rounded-full border border-black bg-black px-3 py-1 text-[11px] uppercase tracking-[0.22em] text-white dark:border-white dark:bg-white dark:text-black">
            {t("login.title")}
          </div>
          <h2 className="mt-5 font-display text-4xl font-bold tracking-[-0.05em] text-slate-950 dark:text-slate-50">
            {t("login.submit")}
          </h2>
          <p className="mt-3 text-sm leading-6 text-slate-500 dark:text-slate-400">{t("login.subtitle")}</p>

          <form className="mt-8 space-y-4" onSubmit={handleSubmit}>
            <div>
              <label className="mb-2 block text-sm font-semibold text-slate-700 dark:text-slate-200">{t("login.identifier")}</label>
              <Input value={identifier} onChange={(event) => setIdentifier(event.target.value)} />
            </div>
            <div>
              <label className="mb-2 block text-sm font-semibold text-slate-700 dark:text-slate-200">{t("login.password")}</label>
              <Input type="password" value={password} onChange={(event) => setPassword(event.target.value)} />
            </div>
            <label className="flex items-center gap-2 text-sm text-slate-600 dark:text-slate-300">
              <input checked={rememberMe} className="h-4 w-4 accent-black dark:accent-white" onChange={(event) => setRememberMe(event.target.checked)} type="checkbox" />
              {t("login.rememberMe")}
            </label>
            <Button block disabled={loading} type="submit">
              {loading ? t("actions.loading") : t("login.submit")}
            </Button>
          </form>

          <div className="mt-6 rounded-[24px] border border-slate-200/80 bg-slate-50/85 p-4 text-sm text-slate-600 dark:border-slate-800 dark:bg-slate-900/70 dark:text-slate-300">
            <div className="font-semibold text-slate-900 dark:text-slate-100">admin@shieldpanel.local</div>
            <div className="mt-1 font-mono text-[13px]">ChangeMe123!</div>
          </div>
        </Card>
      </div>
    </div>
  );
}
