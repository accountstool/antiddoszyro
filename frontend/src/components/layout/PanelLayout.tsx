import type { ReactNode } from "react";
import { Link, NavLink, useLocation } from "react-router-dom";
import { BarChart3, Globe, History, LayoutDashboard, Moon, Settings, ShieldAlert, Sun, Users } from "lucide-react";
import { useTranslation } from "react-i18next";

import { useAuth } from "../../app/auth";
import { useUI } from "../../app/ui";
import { Button } from "../ui/Button";
import { Select } from "../ui/Select";

const navigation = [
  { to: "/", icon: LayoutDashboard, key: "nav.dashboard" },
  { to: "/domains", icon: Globe, key: "nav.domains" },
  { to: "/statistics", icon: BarChart3, key: "nav.statistics" },
  { to: "/ip-control", icon: ShieldAlert, key: "nav.ipControl" },
  { to: "/settings", icon: Settings, key: "nav.settings" },
  { to: "/users", icon: Users, key: "nav.users" },
  { to: "/audit-logs", icon: History, key: "nav.auditLogs" }
];

export function PanelLayout({ children }: { children: ReactNode }) {
  const { t } = useTranslation();
  const { logout, user } = useAuth();
  const { language, setLanguage, theme, toggleTheme } = useUI();
  const location = useLocation();
  const isActive = (path: string) => location.pathname === path || (path !== "/" && location.pathname.startsWith(`${path}/`));
  const currentPage =
    navigation.find((item) => isActive(item.to))?.key ?? "nav.dashboard";

  return (
    <div className="relative min-h-screen">
      <div className="pointer-events-none fixed inset-0 panel-grid-overlay opacity-40" />

      <div className="mx-auto max-w-[1720px] px-4 py-5 lg:px-6">
        <div className="grid gap-5 xl:grid-cols-[296px_minmax(0,1fr)]">
          <aside className="hidden xl:flex xl:sticky xl:top-5 xl:h-[calc(100vh-2.5rem)] xl:flex-col xl:rounded-[30px] xl:border xl:border-black xl:bg-black xl:p-5 xl:text-white xl:shadow-lift dark:xl:border-white/12 dark:xl:bg-neutral-950">
            <Link to="/" className="block rounded-[26px] border border-white/10 bg-gradient-to-br from-neutral-900 to-black p-6 text-white">
              <div className="inline-flex items-center rounded-full border border-white/14 px-3 py-1 text-[11px] uppercase tracking-[0.24em] text-white/72">
                ShieldPanel
              </div>
              <div className="mt-5 font-display text-[1.9rem] font-bold leading-[1.02] tracking-[-0.06em] text-balance">
                {t("brand.tagline")}
              </div>
              <div className="mt-4 text-sm leading-6 text-white/60">{t("brand.description")}</div>
            </Link>

            <nav className="mt-6 flex-1 space-y-2">
              {navigation.map((item) => {
                const Icon = item.icon;
                return (
                  <NavLink
                    key={item.to}
                    to={item.to}
                    className={({ isActive }) =>
                      `group flex items-center gap-3 rounded-2xl border px-4 py-3 text-sm font-semibold transition ${
                        isActive
                          ? "border-white bg-white text-black shadow-sm"
                          : "border-transparent bg-transparent text-white/68 hover:border-white/12 hover:bg-white/6 hover:text-white"
                      }`
                    }
                  >
                    {({ isActive }) => (
                      <>
                        <span
                          className={`grid h-10 w-10 place-items-center rounded-2xl border transition ${
                            isActive
                              ? "border-black/10 bg-black text-white"
                              : "border-white/12 bg-white/6 text-white/82 group-hover:border-white/18 group-hover:bg-white group-hover:text-black"
                          }`}
                        >
                          <Icon size={18} />
                        </span>
                        <span>{t(item.key)}</span>
                      </>
                    )}
                  </NavLink>
                );
              })}
            </nav>

            <div className="space-y-3">
              <div className="rounded-[24px] border border-white/10 bg-white/[0.04] p-4">
                <div className="text-[11px] uppercase tracking-[0.2em] text-white/42">{t("header.welcome")}</div>
                <div className="mt-2 text-base font-semibold text-white">{user?.displayName || user?.username}</div>
                <div className="mt-1 text-sm text-white/56">{user?.email}</div>
              </div>
              <div className="rounded-[24px] border border-white/10 bg-white/[0.04] p-4">
                <div className="text-[11px] uppercase tracking-[0.2em] text-white/42">{t("nav.statistics")}</div>
                <div className="mt-2 text-sm font-semibold leading-6 text-white">{t("feature.stats")}</div>
              </div>
            </div>
          </aside>

          <div className="flex min-w-0 flex-1 flex-col gap-5">
            <header className="rounded-[26px] border border-black/10 bg-white p-4 shadow-panel dark:border-white/10 dark:bg-neutral-950">
              <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
                <div className="min-w-0">
                  <div className="text-[11px] uppercase tracking-[0.24em] text-slate-500 dark:text-slate-400">ShieldPanel</div>
                  <div className="mt-2 flex flex-wrap items-center gap-3">
                    <h1 className="font-display text-[1.95rem] font-bold tracking-[-0.05em] text-slate-950 dark:text-slate-50">
                      {t(currentPage)}
                    </h1>
                    <span className="rounded-full border border-black/10 bg-neutral-100 px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.18em] text-slate-700 dark:border-white/10 dark:bg-neutral-900 dark:text-slate-200">
                      {(user?.role || "admin").replace("_", " ")}
                    </span>
                  </div>
                  <p className="mt-1 truncate text-sm text-slate-500 dark:text-slate-400">
                    {user?.displayName || user?.username} | {user?.email}
                  </p>
                </div>

                <div className="flex flex-wrap items-center gap-3">
                  <Select className="min-w-[154px]" value={language} onChange={(event) => setLanguage(event.target.value)}>
                    <option value="en">English</option>
                    <option value="vi">Tieng Viet</option>
                  </Select>
                  <Button variant="secondary" onClick={toggleTheme}>
                    {theme === "dark" ? <Sun size={16} /> : <Moon size={16} />}
                  </Button>
                  <Button variant="secondary" onClick={() => void logout()}>
                    {t("actions.logout")}
                  </Button>
                </div>
              </div>

              <div className="mt-4 xl:hidden">
                <div className="panel-scroll flex gap-2 overflow-x-auto pb-1">
                  {navigation.map((item) => {
                    const Icon = item.icon;
                    const active = isActive(item.to);
                    return (
                      <NavLink
                        key={item.to}
                        to={item.to}
                        className={`inline-flex shrink-0 items-center gap-2 rounded-2xl border px-4 py-2.5 text-sm font-semibold transition ${
                          active
                            ? "border-black bg-black text-white dark:border-white dark:bg-white dark:text-black"
                            : "border-black/10 bg-white/70 text-slate-700 dark:border-white/10 dark:bg-neutral-950/70 dark:text-slate-200"
                        }`}
                      >
                        <Icon size={16} />
                        {t(item.key)}
                      </NavLink>
                    );
                  })}
                </div>
              </div>
            </header>

            <main className="min-w-0 pb-6">{children}</main>
          </div>
        </div>
      </div>
    </div>
  );
}
