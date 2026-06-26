import { NavLink, useNavigate } from "react-router-dom";
import type { ReactNode } from "react";
import { useAuth } from "../auth";

interface NavTab {
  to: string;
  label: string;
  icon: string;
  adminOnly?: boolean;
}

const TABS: NavTab[] = [
  { to: "/submit", label: "Submit", icon: "🕐" },
  { to: "/history", label: "History", icon: "📋" },
  { to: "/timecard", label: "Timecard", icon: "📅" },
  { to: "/admin", label: "Admin", icon: "⚙️", adminOnly: true },
];

export function Layout({
  title,
  subtitle,
  children,
}: {
  title: string;
  subtitle?: string;
  children: ReactNode;
}) {
  const { employee, logout } = useAuth();
  const navigate = useNavigate();

  async function handleLogout() {
    await logout();
    navigate("/login", { replace: true });
  }

  const tabs = TABS.filter((t) => !t.adminOnly || employee?.isAdmin);

  return (
    <div className="app">
      <header className="header">
        <div className="header-row">
          <div>
            <h1>{title}</h1>
            {subtitle && <div className="subtitle">{subtitle}</div>}
          </div>
          <button className="btn btn-secondary btn-sm" onClick={handleLogout}>
            Sign out
          </button>
        </div>
      </header>

      <main className="content">{children}</main>

      <nav className="bottom-nav">
        {tabs.map((tab) => (
          <NavLink
            key={tab.to}
            to={tab.to}
            className={({ isActive }) =>
              `nav-item${isActive ? " active" : ""}`
            }
          >
            <span className="icon">{tab.icon}</span>
            <span>{tab.label}</span>
          </NavLink>
        ))}
      </nav>
    </div>
  );
}
