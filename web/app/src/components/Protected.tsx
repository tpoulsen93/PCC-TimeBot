import { Navigate } from "react-router-dom";
import type { ReactNode } from "react";
import { useAuth } from "../auth";

export function Protected({
  children,
  adminOnly = false,
}: {
  children: ReactNode;
  adminOnly?: boolean;
}) {
  const { employee, loading } = useAuth();

  if (loading) {
    return (
      <div className="app">
        <div className="loading">Loading…</div>
      </div>
    );
  }

  if (!employee) {
    return <Navigate to="/login" replace />;
  }

  if (adminOnly && !employee.isAdmin) {
    return <Navigate to="/submit" replace />;
  }

  return <>{children}</>;
}
