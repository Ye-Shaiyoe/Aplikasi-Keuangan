import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { GoogleOAuthProvider } from '@react-oauth/google';
import { AuthProvider, useAuth } from './context/AuthContext';
import { ThemeProvider } from './context/ThemeContext';
import Layout from './components/Layout';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Transactions from './pages/Transactions';
import Categories from './pages/Categories';
import Reports from './pages/Reports';
import SavingsGoals from './pages/SavingsGoals';
import Budgets from './pages/Budgets';
import Recurring from './pages/Recurring';
import InsightCharts from './pages/InsightCharts';
import InsightTables from './pages/InsightTables';
import AdvancedAnalytics from './pages/AdvancedAnalytics';
import Docs from './pages/Docs';

function ProtectedRoute({ children }) {
  const { isAuthenticated } = useAuth();
  if (!isAuthenticated) return <Navigate to="/login" replace />;
  return <Layout>{children}</Layout>;
}

function PublicRoute({ children }) {
  const { isAuthenticated } = useAuth();
  if (isAuthenticated) return <Navigate to="/" replace />;
  return children;
}

function AppRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<PublicRoute><Login /></PublicRoute>} />
      <Route path="/register" element={<PublicRoute><Register /></PublicRoute>} />
      <Route path="/" element={<ProtectedRoute><Dashboard /></ProtectedRoute>} />
      <Route path="/transactions" element={<ProtectedRoute><Transactions /></ProtectedRoute>} />
      <Route path="/categories" element={<ProtectedRoute><Categories /></ProtectedRoute>} />
      <Route path="/reports" element={<ProtectedRoute><Reports /></ProtectedRoute>} />
      <Route path="/savings" element={<ProtectedRoute><SavingsGoals /></ProtectedRoute>} />
      <Route path="/budgets" element={<ProtectedRoute><Budgets /></ProtectedRoute>} />
      <Route path="/recurring" element={<ProtectedRoute><Recurring /></ProtectedRoute>} />
      <Route path="/insights" element={<Navigate to="/insights/charts" replace />} />
      <Route path="/insights/charts" element={<ProtectedRoute><InsightCharts /></ProtectedRoute>} />
      <Route path="/insights/analytics" element={<ProtectedRoute><AdvancedAnalytics /></ProtectedRoute>} />
      <Route path="/insights/tables" element={<ProtectedRoute><InsightTables /></ProtectedRoute>} />
      <Route path="/docs" element={<ProtectedRoute><Docs /></ProtectedRoute>} />
    </Routes>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <ThemeProvider>
        <AuthProvider>
          <GoogleOAuthProvider clientId={import.meta.env.VITE_GOOGLE_CLIENT_ID || ''}>
            <AppRoutes />
          </GoogleOAuthProvider>
        </AuthProvider>
      </ThemeProvider>
    </BrowserRouter>
  );
}
