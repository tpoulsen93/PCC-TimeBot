import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { api, ApiError } from "../api";
import { useAuth } from "../auth";

export function LoginPage() {
  const [email, setEmail] = useState("");
  const [sent, setSent] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { employee, loading, refresh } = useAuth();
  const navigate = useNavigate();
  const [params] = useSearchParams();

  useEffect(() => {
    if (!loading && employee) {
      navigate("/submit", { replace: true });
    }
  }, [loading, employee, navigate]);

  useEffect(() => {
    if (params.get("error")) {
      setError("That sign-in link is invalid or has expired. Request a new one.");
    }
  }, [params]);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      const res = await api.post<{ dev?: boolean }>("/auth/request-link", {
        email: email.trim(),
      });
      if (res.dev) {
        await refresh();
        navigate("/submit", { replace: true });
        return;
      }
      setSent(true);
    } catch (err) {
      setError(
        err instanceof ApiError
          ? err.message
          : "Something went wrong. Please try again."
      );
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="login">
      <div className="brand">
        <div className="logo">🕐</div>
        <h1>PCC TimeBot</h1>
        <p>Sign in to submit your hours</p>
      </div>

      {sent ? (
        <div className="card">
          <h2>Check your email</h2>
          <p className="muted">
            If <strong>{email.trim()}</strong> is registered, we've sent a
            secure sign-in link. It expires in 15 minutes.
          </p>
          <button
            className="btn btn-secondary"
            onClick={() => {
              setSent(false);
              setError(null);
            }}
          >
            Use a different email
          </button>
        </div>
      ) : (
        <form className="card" onSubmit={handleSubmit}>
          {error && <div className="message error">{error}</div>}
          <div className="field">
            <label htmlFor="email">Email address</label>
            <input
              id="email"
              className="input"
              type="email"
              autoComplete="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          <button
            className="btn"
            type="submit"
            disabled={submitting || !email.trim()}
          >
            {submitting ? "Sending…" : "Send sign-in link"}
          </button>
        </form>
      )}
    </div>
  );
}
