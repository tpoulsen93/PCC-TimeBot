import { useMemo, useState } from "react";
import { Layout } from "../components/Layout";
import { api, ApiError } from "../api";
import { calculateTime } from "../timeCalc";
import type { SubmitResponse } from "../types";

function todayISO(): string {
  const now = new Date();
  const tzOffset = now.getTimezoneOffset() * 60000;
  return new Date(now.getTime() - tzOffset).toISOString().slice(0, 10);
}

export function SubmitPage() {
  const [date, setDate] = useState(todayISO());
  const [start, setStart] = useState("");
  const [end, setEnd] = useState("");
  const [lunch, setLunch] = useState("0");
  const [extra, setExtra] = useState("0");
  const [location, setLocation] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const preview = useMemo(
    () => calculateTime(start, end, lunch, extra),
    [start, end, lunch, extra]
  );

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setMessage(null);
    if (preview.error) {
      setError(preview.error);
      return;
    }
    setSubmitting(true);
    try {
      const res = await api.post<SubmitResponse>("/timecards", {
        date,
        start,
        end,
        lunch: lunch.trim() === "" ? "0" : lunch.trim(),
        extra: extra.trim() === "" ? "0" : extra.trim(),
        location: location.trim(),
      });
      setMessage(`Submitted ${res.hours} hours for ${res.date}.`);
      setStart("");
      setEnd("");
      setLunch("0");
      setExtra("0");
      setLocation("");
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Failed to submit hours."
      );
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <Layout title="Submit Time" subtitle="Log your hours for the day">
      <form className="card" onSubmit={handleSubmit}>
        {message && <div className="message success">{message}</div>}
        {error && <div className="message error">{error}</div>}

        <div className="field">
          <label htmlFor="date">Date</label>
          <input
            id="date"
            className="input"
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
            required
          />
        </div>

        <div className="row">
          <div className="field">
            <label htmlFor="start">Start time</label>
            <input
              id="start"
              className="input"
              placeholder="9:00am"
              value={start}
              onChange={(e) => setStart(e.target.value)}
              required
            />
          </div>
          <div className="field">
            <label htmlFor="end">End time</label>
            <input
              id="end"
              className="input"
              placeholder="5:00pm"
              value={end}
              onChange={(e) => setEnd(e.target.value)}
              required
            />
          </div>
        </div>

        <div className="row">
          <div className="field">
            <label htmlFor="lunch">Lunch (hrs)</label>
            <input
              id="lunch"
              className="input"
              inputMode="decimal"
              value={lunch}
              onChange={(e) => setLunch(e.target.value)}
            />
          </div>
          <div className="field">
            <label htmlFor="extra">Extra (hrs)</label>
            <input
              id="extra"
              className="input"
              inputMode="decimal"
              value={extra}
              onChange={(e) => setExtra(e.target.value)}
            />
          </div>
        </div>

        <div className="field">
          <label htmlFor="location">Location</label>
          <input
            id="location"
            className="input"
            placeholder="Where did you work?"
            value={location}
            onChange={(e) => setLocation(e.target.value)}
          />
        </div>

        <div className="preview">
          <div className="label">Total hours</div>
          <div className="total">
            {preview.hours !== null ? preview.hours.toFixed(2) : "—"}
          </div>
        </div>

        <button className="btn" type="submit" disabled={submitting}>
          {submitting ? "Submitting…" : "Submit hours"}
        </button>
      </form>
    </Layout>
  );
}
