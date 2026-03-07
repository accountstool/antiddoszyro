import type { ReactNode } from "react";
import { Navigate, Route, Routes } from "react-router-dom";

import { useAuth } from "../app/auth";
import { PanelLayout } from "../components/layout/PanelLayout";
import { Spinner } from "../components/ui/Spinner";
import { AuditLogsPage } from "../pages/AuditLogsPage";
import { DashboardPage } from "../pages/DashboardPage";
import { DomainDetailPage } from "../pages/DomainDetailPage";
import { DomainsPage } from "../pages/DomainsPage";
import { IPControlPage } from "../pages/IPControlPage";
import { LoginPage } from "../pages/LoginPage";
import { SettingsPage } from "../pages/SettingsPage";
import { StatisticsPage } from "../pages/StatisticsPage";
import { UsersPage } from "../pages/UsersPage";

function Protected({ children }: { children: ReactNode }) {
  const { ready, user } = useAuth();

  if (!ready) {
    return (
      <div className="grid min-h-screen place-items-center">
        <Spinner />
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return <PanelLayout>{children}</PanelLayout>;
}

export function AppRouter() {
  const { user, ready } = useAuth();

  if (!ready) {
    return (
      <div className="grid min-h-screen place-items-center">
        <Spinner />
      </div>
    );
  }

  return (
    <Routes>
      <Route path="/login" element={user ? <Navigate to="/" replace /> : <LoginPage />} />
      <Route
        path="/"
        element={
          <Protected>
            <DashboardPage />
          </Protected>
        }
      />
      <Route
        path="/domains"
        element={
          <Protected>
            <DomainsPage />
          </Protected>
        }
      />
      <Route
        path="/domains/:id"
        element={
          <Protected>
            <DomainDetailPage />
          </Protected>
        }
      />
      <Route
        path="/statistics"
        element={
          <Protected>
            <StatisticsPage />
          </Protected>
        }
      />
      <Route
        path="/ip-control"
        element={
          <Protected>
            <IPControlPage />
          </Protected>
        }
      />
      <Route
        path="/settings"
        element={
          <Protected>
            <SettingsPage />
          </Protected>
        }
      />
      <Route
        path="/users"
        element={
          <Protected>
            <UsersPage />
          </Protected>
        }
      />
      <Route
        path="/audit-logs"
        element={
          <Protected>
            <AuditLogsPage />
          </Protected>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
