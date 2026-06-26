import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { AuthProvider } from "./auth";
import { Protected } from "./components/Protected";
import { LoginPage } from "./pages/LoginPage";
import { SubmitPage } from "./pages/SubmitPage";
import { HistoryPage } from "./pages/HistoryPage";
import { TimecardPage } from "./pages/TimecardPage";
import { AdminPage } from "./pages/AdminPage";

export function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route
            path="/submit"
            element={
              <Protected>
                <SubmitPage />
              </Protected>
            }
          />
          <Route
            path="/history"
            element={
              <Protected>
                <HistoryPage />
              </Protected>
            }
          />
          <Route
            path="/timecard"
            element={
              <Protected>
                <TimecardPage />
              </Protected>
            }
          />
          <Route
            path="/admin"
            element={
              <Protected adminOnly>
                <AdminPage />
              </Protected>
            }
          />
          <Route path="/" element={<Navigate to="/submit" replace />} />
          <Route path="*" element={<Navigate to="/submit" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}
