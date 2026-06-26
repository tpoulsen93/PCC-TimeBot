import { useEffect, useState } from "react";
import { Layout } from "../components/Layout";
import { api, ApiError } from "../api";
import type { SummaryResponse } from "../types";

function formatDate(iso: string): string {
  const d = new Date(iso + "T00:00:00");
  return d.toLocaleDateString(undefined, {
    weekday: "short",
    month: "short",
    day: "numeric",
  });
}

function formatRange(start: string, end: string): string {
  const s = new Date(start + "T00:00:00");
  const e = new Date(end + "T00:00:00");
  const opts: Intl.DateTimeFormatOptions = { month: "short", day: "numeric" };
  return `${s.toLocaleDateString(undefined, opts)} – ${e.toLocaleDateString(
    undefined,
    opts
  )}`;
}

export function TimecardPage() {
  const [summary, setSummary] = useState<SummaryResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const res = await api.get<SummaryResponse>("/timecards/summary");
        if (active) setSummary(res);
      } catch (err) {
        if (active) {
          setError(
            err instanceof ApiError ? err.message : "Failed to load timecard."
          );
        }
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, []);

  return (
    <Layout title="Timecard" subtitle="This week at a glance">
      {loading ? (
        <div className="loading">Loading…</div>
      ) : error ? (
        <div className="message error">{error}</div>
      ) : summary ? (
        <>
          <div className="card">
            <h2>{formatRange(summary.weekStart, summary.weekEnd)}</h2>
            {summary.days.length === 0 ? (
              <p className="muted">No hours logged this week yet.</p>
            ) : (
              summary.days.map((d, i) => (
                <div className="list-item" key={`${d.date}-${i}`}>
                  <div>
                    <div className="primary">{formatDate(d.date)}</div>
                    <div className="secondary">
                      {d.location || "No location"}
                    </div>
                  </div>
                  <div className="value">{d.hours.toFixed(2)}h</div>
                </div>
              ))
            )}
          </div>

          <div className="preview">
            <div className="label">Week total</div>
            <div className="total">{summary.totalHours.toFixed(2)}</div>
          </div>

          <div className="card">
            <div className="list-item">
              <div className="primary">Payday</div>
              <div className="value">{formatDate(summary.payday)}</div>
            </div>
          </div>
        </>
      ) : null}
    </Layout>
  );
}
