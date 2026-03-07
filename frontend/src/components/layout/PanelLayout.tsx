import type { ReactNode } from "react";
import { Link, NavLink } from "react-router-dom";
import { BarChart3, Globe, History, LayoutDashboard, Moon, Settings, ShieldAlert, Sun, Users } from "lucide-react";
import { useTranslation } from "react-i18next";

import { useAuth } from "../../app/auth";
import { useUI } from "../../app/ui";
import { Button } from "../ui/Button";

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

  return (
    <div className="min-h-screen bg-panel-grid bg-[size:28px_28px]">
      <div className="mx-auto flex min-h-screen max-w-[1600px] gap-6 px-4 py-4 lg:px-6">
        <aside className="hidden w-72 shrink-0 rounded-[28px] border border-slate-200/70 bg-white/85 p-5 shadow-panel backdrop-blur dark:border-slate-800 dark:bg-slate-950/75 lg:block">
          <Link to="/" className="block rounded-3xl bg-gradient-to-br from-teal-500 via-teal-600 to-cyan-700 p-5 text-white">
            <div className="text-xs uppercase tracking-[0.2em] opacity-80">ShieldPanel</div>
            <div className="mt-3 text-2xl font-extrabold">{t("brand.tagline")}</div>
            <div className="mt-2 text-sm opacity-80">{t("brand.description")}</div>
          </Link>

          <nav className="mt-8 space-y-2">
            {navigation.map((item) => {
              const Icon = item.icon;
              return (
                <NavLink
                  key={item.to}
                  to={item.to}
                  className={({ isActive }) =>
                    `flex items-center gap-3 rounded-2xl px-4 py-3 text-sm font-semibold transition ${
                      isActive
                        ? "bg-teal-50 text-teal-700 dark:bg-teal-950/40 dark:text-teal-300"
                        : "text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-900"
                    }`
                  }
                >
                  <Icon size={18} />
                  {t(item.key)}
                </NavLink>
              );
            })}
          </nav>
        </aside>

        <div className="flex min-w-0 flex-1 flex-col gap-4">
          <header className="flex flex-col gap-3 rounded-[28px] border border-slate-200/70 bg-white/85 p-4 shadow-panel backdrop-blur dark:border-slate-800 dark:bg-slate-950/75 lg:flex-row lg:items-center lg:justify-between">
            <div>
              <p className="text-xs uppercase tracking-[0.2em] text-slate-500">{t("header.welcome")}</p>
              <h1 className="text-2xl font-extrabold text-slate-900 dark:text-slate-100">{user?.displayName || user?.username}</h1>
            </div>
            <div className="flex flex-wrap items-center gap-3">
              <select
                className="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-900"
                value={language}
                onChange={(event) => setLanguage(event.target.value)}
              >
                <option value="en">English</option>
                <option value="vi">Tieng Viet</option>
              </select>
              <Button variant="secondary" onClick={toggleTheme}>
                {theme === "dark" ? <Sun size={16} /> : <Moon size={16} />}
              </Button>
              <Button variant="secondary" onClick={() => void logout()}>
                {t("actions.logout")}
              </Button>
            </div>
          </header>

          <main className="min-w-0">{children}</main>
        </div>
      </div>
    </div>
  );
}
