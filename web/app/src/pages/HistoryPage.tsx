import { useEffect, useState } from "react";
import { Layout } from "../components/Layout";
import { api, ApiError } from "../api";
import type { HistoryEntry, HistoryResponse } from "../types";

function formatDate(iso: string): string {
  const d = new Date(iso + "T00:00:00");
  return d.toLocaleDateString(undefined, {
    weekday: "short",
    month: "short",
    day: "numeric",
  });
}

export function HistoryPage() {
  const [entries, setEntries] = useState<HistoryEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const res = await api.get<HistoryResponse>("/timecards/history");
        if (active) setEntries(res.entries ?? []);
      } catch (err) {
        if (active) {
          setError(
            err instanceof ApiError ? err.message : "Failed to load history."
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

  const total = entries.reduce((sum, e) => sum + e.hours, 0);

  return (
    <Layout title="History" subtitle="Your recent submissions">
      {loading ? (
        <div className="loading">Loading…</div>
      ) : error ? (
        <div className="message error">{error}</div>
      ) : entries.length === 0 ? (
        <div className="card center">
          <p className="muted">No submissions in the last 60 days.</p>
        </div>
      ) : (
        <>
          <div className="preview">
            <div className="label">Total · last 60 days</div>
            <div className="total">{total.toFixed(2)}</div>
          </div>
          <div className="card">
            {entries.map((e, i) => (
              <div className="list-item" key={`${e.date}-${i}`}>
                <div>
                  <div className="primary">{formatDate(e.date)}</div>
                  <div className="secondary">
                    {e.location || "No location"}
                  </div>
                </div>
                <div className="value">{e.hours.toFixed(2)}h</div>
              </div>
            ))}
          </div>
        </>
      )}
    </Layout>
  );
}
