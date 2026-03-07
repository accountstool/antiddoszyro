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
  const currentPage =
    navigation.find((item) => location.pathname === item.to || (item.to !== "/" && location.pathname.startsWith(`${item.to}/`)))?.key ??
    "nav.dashboard";

  return (
    <div className="relative min-h-screen">
      <div className="pointer-events-none fixed inset-0 panel-grid-overlay opacity-30" />

      <div className="mx-auto max-w-[1680px] px-4 py-4 lg:px-6">
        <div className="grid gap-4 xl:grid-cols-[320px_minmax(0,1fr)]">
          <aside className="hidden xl:flex xl:sticky xl:top-4 xl:h-[calc(100vh-2rem)] xl:flex-col xl:rounded-[32px] xl:border xl:border-black/90 xl:bg-black xl:p-5 xl:text-white xl:shadow-lift xl:backdrop-blur-xl dark:xl:border-white/12 dark:xl:bg-neutral-950">
            <Link to="/" className="relative block overflow-hidden rounded-[28px] bg-mesh-teal p-6 text-white shadow-lift">
              <div className="absolute inset-0 panel-grid-overlay opacity-15" />
              <div className="relative">
                <div className="inline-flex items-center rounded-full border border-white/20 bg-white/8 px-3 py-1 text-[11px] uppercase tracking-[0.22em] text-white/78">
                  ShieldPanel
                </div>
                <div className="mt-5 font-display text-[2rem] font-bold leading-[1.05] tracking-[-0.05em] text-balance">
                  {t("brand.tagline")}
                </div>
                <div className="mt-3 text-sm leading-6 text-white/64">{t("brand.description")}</div>
              </div>
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
                          : "border-white/8 text-white/72 hover:border-white/16 hover:bg-white/8 hover:text-white"
                      }`
                    }
                  >
                    {({ isActive }) => (
                      <>
                        <span
                          className={`grid h-10 w-10 place-items-center rounded-2xl border transition ${
                            isActive
                              ? "border-black/10 bg-black text-white"
                              : "border-white/12 bg-white/6 text-white/82 group-hover:bg-white group-hover:text-black"
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

            <div className="grid gap-3">
              <div className="rounded-2xl border border-white/10 bg-white/6 p-4">
                <div className="text-[11px] uppercase tracking-[0.2em] text-white/45">{t("nav.statistics")}</div>
                <div className="mt-2 text-sm font-semibold text-white">{t("feature.stats")}</div>
              </div>
              <div className="rounded-2xl border border-white/10 bg-white/6 p-4">
                <div className="text-[11px] uppercase tracking-[0.2em] text-white/45">{t("nav.ipControl")}</div>
                <div className="mt-2 text-sm font-semibold text-white">{t("feature.challenge")}</div>
              </div>
            </div>
          </aside>

          <div className="flex min-w-0 flex-1 flex-col gap-4">
            <header className="rounded-[32px] border border-black/10 bg-white/76 p-4 shadow-panel backdrop-blur-xl dark:border-white/10 dark:bg-neutral-950/82">
              <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
                <div>
                  <div className="text-[11px] uppercase tracking-[0.22em] text-slate-500 dark:text-slate-400">{t(currentPage)}</div>
                  <div className="mt-2 flex flex-wrap items-center gap-3">
                    <h1 className="font-display text-3xl font-bold tracking-[-0.05em] text-slate-950 dark:text-slate-50">
                      {user?.displayName || user?.username}
                    </h1>
                    <span className="rounded-full border border-black bg-black px-3 py-1 text-[11px] uppercase tracking-[0.18em] text-white dark:border-white dark:bg-white dark:text-black">
                      {(user?.role || "admin").replace("_", " ")}
                    </span>
                  </div>
                  <p className="mt-2 max-w-2xl text-sm leading-6 text-slate-600 dark:text-slate-300">{t("brand.description")}</p>
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
                    const active = location.pathname === item.to;
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
