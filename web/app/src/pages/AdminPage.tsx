import { useEffect, useState } from "react";
import { Layout } from "../components/Layout";
import { api, ApiError } from "../api";
import type { Employee, AdminTimecardsResponse } from "../types";

interface NewEmployee {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  isAdmin: boolean;
}

const EMPTY_NEW: NewEmployee = {
  firstName: "",
  lastName: "",
  email: "",
  phone: "",
  isAdmin: false,
};

function EmployeesSection() {
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);
  const [form, setForm] = useState<NewEmployee>(EMPTY_NEW);
  const [creating, setCreating] = useState(false);
  const [showForm, setShowForm] = useState(false);

  async function load() {
    try {
      const res = await api.get<{ employees: Employee[] }>("/admin/employees");
      setEmployees(res.employees ?? []);
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Failed to load employees."
      );
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void load();
  }, []);

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setMessage(null);
    setCreating(true);
    try {
      await api.post<Employee>("/admin/employees", {
        firstName: form.firstName.trim(),
        lastName: form.lastName.trim(),
        email: form.email.trim(),
        phone: form.phone.trim(),
        supervisorId: null,
        isAdmin: form.isAdmin,
      });
      setMessage(`Added ${form.firstName} ${form.lastName}.`);
      setForm(EMPTY_NEW);
      setShowForm(false);
      setLoading(true);
      await load();
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Failed to create employee."
      );
    } finally {
      setCreating(false);
    }
  }

  async function toggleAdmin(emp: Employee) {
    setError(null);
    try {
      await api.put<Employee>(`/admin/employees/${emp.id}`, {
        isAdmin: !emp.isAdmin,
      });
      setEmployees((prev) =>
        prev.map((e) =>
          e.id === emp.id ? { ...e, isAdmin: !emp.isAdmin } : e
        )
      );
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Failed to update employee."
      );
    }
  }

  return (
    <>
      {message && <div className="message success">{message}</div>}
      {error && <div className="message error">{error}</div>}

      <div className="card">
        <h2>Team</h2>
        {loading ? (
          <p className="muted">Loading…</p>
        ) : employees.length === 0 ? (
          <p className="muted">No employees yet.</p>
        ) : (
          employees.map((e) => (
            <div className="list-item" key={e.id}>
              <div>
                <div className="primary">
                  {e.firstName} {e.lastName}
                  {e.isAdmin && <span className="badge">Admin</span>}
                </div>
                <div className="secondary">{e.email}</div>
              </div>
              <button
                className="btn btn-secondary btn-sm"
                onClick={() => toggleAdmin(e)}
              >
                {e.isAdmin ? "Revoke admin" : "Make admin"}
              </button>
            </div>
          ))
        )}
      </div>

      {showForm ? (
        <form className="card" onSubmit={handleCreate}>
          <h2>Add employee</h2>
          <div className="row">
            <div className="field">
              <label htmlFor="fn">First name</label>
              <input
                id="fn"
                className="input"
                value={form.firstName}
                onChange={(e) =>
                  setForm({ ...form, firstName: e.target.value })
                }
                required
              />
            </div>
            <div className="field">
              <label htmlFor="ln">Last name</label>
              <input
                id="ln"
                className="input"
                value={form.lastName}
                onChange={(e) =>
                  setForm({ ...form, lastName: e.target.value })
                }
                required
              />
            </div>
          </div>
          <div className="field">
            <label htmlFor="em">Email</label>
            <input
              id="em"
              className="input"
              type="email"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              required
            />
          </div>
          <div className="field">
            <label htmlFor="ph">Phone</label>
            <input
              id="ph"
              className="input"
              value={form.phone}
              onChange={(e) => setForm({ ...form, phone: e.target.value })}
            />
          </div>
          <div className="field">
            <label>
              <input
                type="checkbox"
                checked={form.isAdmin}
                onChange={(e) =>
                  setForm({ ...form, isAdmin: e.target.checked })
                }
              />{" "}
              Grant admin access
            </label>
          </div>
          <div className="row">
            <button
              type="button"
              className="btn btn-secondary"
              onClick={() => {
                setShowForm(false);
                setForm(EMPTY_NEW);
              }}
            >
              Cancel
            </button>
            <button className="btn" type="submit" disabled={creating}>
              {creating ? "Adding…" : "Add"}
            </button>
          </div>
        </form>
      ) : (
        <button className="btn" onClick={() => setShowForm(true)}>
          Add employee
        </button>
      )}
    </>
  );
}

function todayISO(): string {
  const now = new Date();
  const tzOffset = now.getTimezoneOffset() * 60000;
  return new Date(now.getTime() - tzOffset).toISOString().slice(0, 10);
}

function TimecardsSection() {
  const [data, setData] = useState<AdminTimecardsResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);
  const [sending, setSending] = useState(false);
  const [start, setStart] = useState("");
  const [end, setEnd] = useState(todayISO());

  async function load(params?: { start?: string; end?: string }) {
    setLoading(true);
    setError(null);
    try {
      const qs = new URLSearchParams();
      if (params?.start) qs.set("start", params.start);
      if (params?.end) qs.set("end", params.end);
      const suffix = qs.toString() ? `?${qs.toString()}` : "";
      const res = await api.get<AdminTimecardsResponse>(
        `/admin/timecards${suffix}`
      );
      setData(res);
      setStart(res.start);
      setEnd(res.end);
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Failed to load timecards."
      );
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function handleSend() {
    setSending(true);
    setError(null);
    setMessage(null);
    try {
      const qs = new URLSearchParams();
      if (start) qs.set("start", start);
      if (end) qs.set("end", end);
      const res = await api.post<{ sent: number }>(
        `/admin/timecards/send?${qs.toString()}`
      );
      setMessage(`Sent ${res.sent} timecard${res.sent === 1 ? "" : "s"}.`);
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Failed to send timecards."
      );
    } finally {
      setSending(false);
    }
  }

  return (
    <>
      {message && <div className="message success">{message}</div>}
      {error && <div className="message error">{error}</div>}

      <div className="card">
        <h2>Timecards</h2>
        <div className="row">
          <div className="field">
            <label htmlFor="ts">Start</label>
            <input
              id="ts"
              className="input"
              type="date"
              value={start}
              onChange={(e) => setStart(e.target.value)}
            />
          </div>
          <div className="field">
            <label htmlFor="te">End</label>
            <input
              id="te"
              className="input"
              type="date"
              value={end}
              onChange={(e) => setEnd(e.target.value)}
            />
          </div>
        </div>
        <button
          className="btn btn-secondary"
          onClick={() => load({ start, end })}
          disabled={loading}
        >
          {loading ? "Loading…" : "Update range"}
        </button>
      </div>

      {data && (
        <>
          <div className="card">
            {data.timecards.length === 0 ? (
              <p className="muted">No hours in this range.</p>
            ) : (
              data.timecards.map((t) => (
                <div className="list-item" key={t.employeeId}>
                  <div className="primary">{t.name}</div>
                  <div className="value">{t.totalHours.toFixed(2)}h</div>
                </div>
              ))
            )}
          </div>

          <div className="preview">
            <div className="label">Total hours · ${data.cost.toFixed(2)}</div>
            <div className="total">{data.totalHours.toFixed(2)}</div>
          </div>

          <button className="btn" onClick={handleSend} disabled={sending}>
            {sending ? "Sending…" : "Email timecards"}
          </button>
        </>
      )}
    </>
  );
}

export function AdminPage() {
  const [tab, setTab] = useState<"team" | "timecards">("team");

  return (
    <Layout title="Admin" subtitle="Manage your team">
      <div className="row" style={{ marginBottom: 16 }}>
        <button
          className={`btn ${tab === "team" ? "" : "btn-secondary"}`}
          onClick={() => setTab("team")}
        >
          Team
        </button>
        <button
          className={`btn ${tab === "timecards" ? "" : "btn-secondary"}`}
          onClick={() => setTab("timecards")}
        >
          Timecards
        </button>
      </div>

      {tab === "team" ? <EmployeesSection /> : <TimecardsSection />}
    </Layout>
  );
}
