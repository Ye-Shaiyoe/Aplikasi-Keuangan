import { useState, useEffect } from 'react';
import {
  Wallet, ArrowDownCircle, ArrowUpCircle, TrendingUp, TrendingDown,
  ArrowRight, PiggyBank, Target, Lightbulb, ShieldCheck, AlertTriangle,
  BarChart3, CreditCard, Sparkles, Calendar, ChevronRight, Activity,
} from 'lucide-react';
import {
  PieChart, Pie, Cell, ResponsiveContainer, Tooltip,
  AreaChart, Area, XAxis, YAxis, CartesianGrid, Legend,
} from 'recharts';
import {
  getSummary, getTransactions, getAdvancedAnalytics,
  getBudgetSummary, getSavingsGoals,
} from '../api/client';
import { Link } from 'react-router-dom';

// ── formatters ───────────────────────────────────────────────────────
const fmt = (v) => new Intl.NumberFormat('id-ID').format(v || 0);
const fmtC = (v) => {
  const abs = Math.abs(v || 0);
  if (abs >= 1_000_000_000) return `${(v / 1_000_000_000).toFixed(1)}M`;
  if (abs >= 1_000_000) return `${(v / 1_000_000).toFixed(1)}jt`;
  if (abs >= 1_000) return `${(v / 1_000).toFixed(0)}rb`;
  return fmt(v);
};

// ── loading skeleton ─────────────────────────────────────────────────
function LoadingSkeleton() {
  return (
    <div className="space-y-5 sm:space-y-6 page-enter">
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 sm:gap-4">
        {[...Array(4)].map((_, i) => <div key={i} className="skeleton h-28 sm:h-32 rounded-2xl" />)}
      </div>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <div className="skeleton h-52 rounded-2xl" />
        <div className="lg:col-span-2 skeleton h-52 rounded-2xl" />
      </div>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div className="skeleton h-72 rounded-2xl" />
        <div className="skeleton h-72 rounded-2xl" />
      </div>
    </div>
  );
}

// ── health score ring (SVG) ──────────────────────────────────────────
function HealthRing({ score, rating }) {
  const radius = 40;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (score / 100) * circumference;

  const ringColor =
    score >= 80 ? '#10b981' :
    score >= 60 ? '#3b82f6' :
    score >= 40 ? '#f59e0b' : '#ef4444';

  return (
    <div className="health-ring-animate flex flex-col items-center">
      <div className="relative">
        <svg width="100" height="100" className="transform -rotate-90">
          <circle
            cx="50" cy="50" r={radius}
            stroke="currentColor"
            className="text-gray-100 dark:text-gray-700"
            strokeWidth="8"
            fill="none"
          />
          <circle
            cx="50" cy="50" r={radius}
            stroke={ringColor}
            strokeWidth="8"
            fill="none"
            strokeLinecap="round"
            strokeDasharray={circumference}
            strokeDashoffset={offset}
            style={{ animation: 'health-ring-fill 1s ease-out forwards' }}
          />
        </svg>
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <span className="text-2xl font-bold text-gray-800 dark:text-gray-100">{score}</span>
          <span className="text-[10px] text-gray-400 dark:text-gray-500">/ 100</span>
        </div>
      </div>
      <span
        className="mt-2 text-xs font-semibold px-2.5 py-1 rounded-full"
        style={{
          color: ringColor,
          backgroundColor: `${ringColor}15`,
        }}
      >
        {rating}
      </span>
    </div>
  );
}

// ── insight card ─────────────────────────────────────────────────────
function InsightCard({ text, index }) {
  const icons = [Lightbulb, ShieldCheck, AlertTriangle, Target, Sparkles];
  const colors = [
    'text-amber-500 bg-amber-50 dark:bg-amber-900/20',
    'text-emerald-500 bg-emerald-50 dark:bg-emerald-900/20',
    'text-red-500 bg-red-50 dark:bg-red-900/20',
    'text-blue-500 bg-blue-50 dark:bg-blue-900/20',
    'text-purple-500 bg-purple-50 dark:bg-purple-900/20',
  ];
  const Icon = icons[index % icons.length];
  const colorClass = colors[index % colors.length];

  return (
    <div className="insight-card flex gap-3 p-3 rounded-xl bg-gray-50 dark:bg-gray-700/40 border border-gray-100 dark:border-gray-600/40 transition-all duration-200 cursor-default">
      <div className={`w-8 h-8 rounded-lg flex items-center justify-center shrink-0 ${colorClass}`}>
        <Icon size={16} />
      </div>
      <p className="text-xs text-gray-600 dark:text-gray-300 leading-relaxed">{text}</p>
    </div>
  );
}

// ── budget progress item ─────────────────────────────────────────────
function BudgetItem({ name, color, spent, amount }) {
  const pct = amount > 0 ? Math.min(100, Math.round((spent / amount) * 100)) : 0;
  const isOver = pct >= 100;

  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between text-xs">
        <div className="flex items-center gap-2">
          <div className="w-2.5 h-2.5 rounded-full shrink-0" style={{ backgroundColor: color || '#6366f1' }} />
          <span className="text-gray-600 dark:text-gray-300 font-medium truncate max-w-[120px]">{name}</span>
        </div>
        <span className={`font-semibold ${isOver ? 'text-red-500' : 'text-gray-500 dark:text-gray-400'}`}>
          {pct}%
        </span>
      </div>
      <div className="w-full h-2 bg-gray-100 dark:bg-gray-700 rounded-full overflow-hidden">
        <div
          className="progress-bar-fill h-full rounded-full"
          style={{
            width: `${Math.min(pct, 100)}%`,
            backgroundColor: isOver ? '#ef4444' : (color || '#6366f1'),
          }}
        />
      </div>
      <div className="flex justify-between text-[10px] text-gray-400 dark:text-gray-500">
        <span>Rp {fmtC(spent)}</span>
        <span>Rp {fmtC(amount)}</span>
      </div>
    </div>
  );
}

// ── savings goal item ────────────────────────────────────────────────
function SavingsItem({ name, current, target, color, deadline }) {
  const pct = target > 0 ? Math.min(100, Math.round((current / target) * 100)) : 0;

  return (
    <div className="flex items-center gap-3 p-3 rounded-xl bg-gray-50 dark:bg-gray-700/40 border border-gray-100 dark:border-gray-600/40">
      <div className="relative w-11 h-11 shrink-0">
        <svg width="44" height="44" className="transform -rotate-90">
          <circle cx="22" cy="22" r="18" stroke="currentColor" className="text-gray-200 dark:text-gray-600" strokeWidth="4" fill="none" />
          <circle
            cx="22" cy="22" r="18"
            stroke={color || '#6366f1'}
            strokeWidth="4" fill="none"
            strokeLinecap="round"
            strokeDasharray={2 * Math.PI * 18}
            strokeDashoffset={2 * Math.PI * 18 * (1 - pct / 100)}
          />
        </svg>
        <div className="absolute inset-0 flex items-center justify-center">
          <span className="text-[10px] font-bold text-gray-600 dark:text-gray-300">{pct}%</span>
        </div>
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-xs font-semibold text-gray-700 dark:text-gray-200 truncate">{name}</p>
        <p className="text-[10px] text-gray-400 dark:text-gray-500 mt-0.5">
          Rp {fmtC(current)} / Rp {fmtC(target)}
        </p>
        {deadline && (
          <p className="text-[10px] text-gray-400 dark:text-gray-500 flex items-center gap-1 mt-0.5">
            <Calendar size={9} />
            {new Date(deadline + 'T00:00:00').toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })}
          </p>
        )}
      </div>
    </div>
  );
}

// ── custom recharts tooltip ──────────────────────────────────────────
function CustomTooltip({ active, payload, label }) {
  if (!active || !payload?.length) return null;
  return (
    <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-2.5 shadow-lg">
      <p className="text-xs text-gray-500 dark:text-gray-400 mb-1.5 font-medium">{label}</p>
      {payload.map((p, i) => (
        <div key={i} className="flex items-center gap-2 text-xs">
          <div className="w-2 h-2 rounded-full" style={{ backgroundColor: p.color }} />
          <span className="text-gray-600 dark:text-gray-300">{p.name}:</span>
          <span className="font-semibold text-gray-800 dark:text-gray-100">Rp {fmt(p.value)}</span>
        </div>
      ))}
    </div>
  );
}

// ── DONUT COLORS ─────────────────────────────────────────────────────
const DEFAULT_COLORS = [
  '#6366f1', '#8b5cf6', '#ec4899', '#f43f5e', '#f97316',
  '#eab308', '#22c55e', '#06b6d4', '#3b82f6', '#a855f7',
];

// ══════════════════════════════════════════════════════════════════════
// ██  MAIN DASHBOARD COMPONENT
// ══════════════════════════════════════════════════════════════════════
export default function Dashboard() {
  const [summary, setSummary] = useState(null);
  const [recent, setRecent] = useState([]);
  const [analytics, setAnalytics] = useState(null);
  const [budgetData, setBudgetData] = useState(null);
  const [savings, setSavings] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const now = new Date();
    const month = now.getMonth() + 1;
    const year = now.getFullYear();

    Promise.all([
      getSummary(month, year),
      getTransactions({ limit: 5, month, year }),
      getAdvancedAnalytics().catch(() => null),
      getBudgetSummary(month, year).catch(() => null),
      getSavingsGoals().catch(() => []),
    ]).then(([s, t, a, b, sg]) => {
      setSummary(s);
      setRecent(t.data || []);
      setAnalytics(a);
      setBudgetData(b);
      setSavings(Array.isArray(sg) ? sg : []);
    }).finally(() => setLoading(false));
  }, []);

  if (loading) return <LoadingSkeleton />;

  // ── derived data ─────────────────────────────────────────────────
  const balance = summary?.balance || 0;
  const income = summary?.total_income || 0;
  const expense = summary?.total_expense || 0;
  const savingsRate = analytics?.savings_rate ?? (income > 0 ? ((income - expense) / income * 100) : 0);
  const monthName = new Date().toLocaleString('id', { month: 'long', year: 'numeric' });

  // Pie data
  const pieData = (summary?.by_category || [])
    .filter((c) => c.total > 0)
    .map((c, i) => ({
      name: c.category_name,
      value: c.total,
      color: c.category_color || DEFAULT_COLORS[i % DEFAULT_COLORS.length],
    }));

  // Area chart data from analytics monthly_metrics
  const monthlyChartData = (analytics?.monthly_metrics || []).map((m) => ({
    name: m.month_name || `${m.month}/${m.year}`,
    Pemasukan: m.income,
    Pengeluaran: m.expense,
  }));

  // Budget items
  const budgetItems = budgetData?.budgets || [];

  // Greeting
  const hour = new Date().getHours();
  const greeting = hour < 12 ? 'Selamat Pagi' : hour < 17 ? 'Selamat Siang' : 'Selamat Malam';

  // ── render ───────────────────────────────────────────────────────
  return (
    <div className="space-y-5 sm:space-y-6 page-enter">

      {/* ── Header ──────────────────────────────────────────────── */}
      <div className="flex items-center justify-between stagger-1">
        <div>
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900 dark:text-gray-100">
            {greeting} 👋
          </h1>
          <p className="text-xs sm:text-sm text-gray-400 dark:text-gray-500 mt-0.5">
            Dashboard · {monthName}
          </p>
        </div>
        <div className={`px-3 py-1.5 rounded-xl text-xs font-semibold flex items-center gap-1.5 ${
          balance >= 0
            ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400'
            : 'bg-red-50 dark:bg-red-900/30 text-red-500 dark:text-red-400'
        }`}>
          {balance >= 0 ? <TrendingUp size={14} /> : <TrendingDown size={14} />}
          Saldo {balance >= 0 ? 'Positif' : 'Negatif'}
        </div>
      </div>

      {/* ── Row 1: 4 Stat Cards ─────────────────────────────────── */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 sm:gap-4 stagger-2">
        {/* Saldo */}
        <div className="stat-gradient-blue rounded-2xl border border-blue-100 dark:border-blue-900/40 p-4 sm:p-5 group card-hover">
          <div className="flex items-center justify-between mb-2">
            <p className="text-[10px] sm:text-xs font-medium text-gray-400 dark:text-gray-500 uppercase tracking-wider">Saldo</p>
            <div className="p-1.5 sm:p-2 rounded-xl bg-blue-500/10 text-blue-600 dark:text-blue-400 transition-transform group-hover:scale-110">
              <Wallet size={16} />
            </div>
          </div>
          <p className="text-lg sm:text-2xl font-bold text-gray-800 dark:text-gray-100 leading-tight">
            <span className="sm:hidden">{fmtC(balance)}</span>
            <span className="hidden sm:inline">{fmt(balance)}</span>
          </p>
          <p className="text-[10px] text-blue-500 dark:text-blue-400 mt-1 font-medium">
            Pemasukan - Pengeluaran
          </p>
        </div>

        {/* Pemasukan */}
        <div className="stat-gradient-green rounded-2xl border border-emerald-100 dark:border-emerald-900/40 p-4 sm:p-5 group card-hover">
          <div className="flex items-center justify-between mb-2">
            <p className="text-[10px] sm:text-xs font-medium text-gray-400 dark:text-gray-500 uppercase tracking-wider">Masuk</p>
            <div className="p-1.5 sm:p-2 rounded-xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 transition-transform group-hover:scale-110">
              <ArrowUpCircle size={16} />
            </div>
          </div>
          <p className="text-lg sm:text-2xl font-bold text-gray-800 dark:text-gray-100 leading-tight">
            <span className="sm:hidden">{fmtC(income)}</span>
            <span className="hidden sm:inline">{fmt(income)}</span>
          </p>
          {analytics?.forecast_income > 0 && (
            <p className="text-[10px] text-emerald-500 dark:text-emerald-400 mt-1 font-medium flex items-center gap-1">
              <TrendingUp size={10} />
              Prediksi: Rp {fmtC(analytics.forecast_income)}
            </p>
          )}
        </div>

        {/* Pengeluaran */}
        <div className="stat-gradient-red rounded-2xl border border-red-100 dark:border-red-900/40 p-4 sm:p-5 group card-hover">
          <div className="flex items-center justify-between mb-2">
            <p className="text-[10px] sm:text-xs font-medium text-gray-400 dark:text-gray-500 uppercase tracking-wider">Keluar</p>
            <div className="p-1.5 sm:p-2 rounded-xl bg-red-500/10 text-red-500 dark:text-red-400 transition-transform group-hover:scale-110">
              <ArrowDownCircle size={16} />
            </div>
          </div>
          <p className="text-lg sm:text-2xl font-bold text-gray-800 dark:text-gray-100 leading-tight">
            <span className="sm:hidden">{fmtC(expense)}</span>
            <span className="hidden sm:inline">{fmt(expense)}</span>
          </p>
          {analytics?.forecast_expense > 0 && (
            <p className="text-[10px] text-red-500 dark:text-red-400 mt-1 font-medium flex items-center gap-1">
              <TrendingUp size={10} />
              Prediksi: Rp {fmtC(analytics.forecast_expense)}
            </p>
          )}
        </div>

        {/* Rasio Tabungan */}
        <div className="stat-gradient-purple rounded-2xl border border-purple-100 dark:border-purple-900/40 p-4 sm:p-5 group card-hover">
          <div className="flex items-center justify-between mb-2">
            <p className="text-[10px] sm:text-xs font-medium text-gray-400 dark:text-gray-500 uppercase tracking-wider">Tabungan</p>
            <div className="p-1.5 sm:p-2 rounded-xl bg-purple-500/10 text-purple-600 dark:text-purple-400 transition-transform group-hover:scale-110">
              <PiggyBank size={16} />
            </div>
          </div>
          <p className="text-lg sm:text-2xl font-bold text-gray-800 dark:text-gray-100 leading-tight">
            {savingsRate.toFixed(1)}%
          </p>
          <p className={`text-[10px] mt-1 font-medium ${
            savingsRate >= 20 ? 'text-emerald-500' :
            savingsRate >= 10 ? 'text-blue-500' :
            savingsRate >= 0  ? 'text-amber-500' : 'text-red-500'
          }`}>
            {savingsRate >= 20 ? '🔥 Luar biasa!' :
             savingsRate >= 10 ? '👍 Bagus' :
             savingsRate >= 0  ? '⚠️ Perlu ditingkatkan' : '🚨 Defisit'}
          </p>
        </div>
      </div>

      {/* ── Row 2: Health Score + AI Insights ───────────────────── */}
      {analytics && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 sm:gap-5 stagger-3">
          {/* Health Score */}
          <div className="bg-white dark:bg-gray-800/90 rounded-2xl border border-gray-100 dark:border-gray-700/60 p-5 sm:p-6 card-hover">
            <div className="flex items-center gap-2 mb-4">
              <ShieldCheck size={18} className="text-indigo-500" />
              <h2 className="text-sm font-semibold text-gray-800 dark:text-gray-100">Kesehatan Keuangan</h2>
            </div>
            <HealthRing score={analytics.health_score} rating={analytics.health_rating} />
            <div className="mt-4 space-y-2">
              <div className="flex items-center justify-between text-xs">
                <span className="text-gray-400 dark:text-gray-500">Rasio Tabungan</span>
                <span className="font-semibold text-gray-600 dark:text-gray-300">{analytics.savings_rate?.toFixed(1) || '0.0'}%</span>
              </div>
              <div className="flex items-center justify-between text-xs">
                <span className="text-gray-400 dark:text-gray-500">Kepatuhan Anggaran</span>
                <span className="font-semibold text-gray-600 dark:text-gray-300">{analytics.budget_adherence?.toFixed(0) || '0'}%</span>
              </div>
            </div>
          </div>

          {/* AI Insights */}
          <div className="lg:col-span-2 bg-white dark:bg-gray-800/90 rounded-2xl border border-gray-100 dark:border-gray-700/60 p-5 sm:p-6 card-hover">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-2">
                <Sparkles size={18} className="text-indigo-500" />
                <h2 className="text-sm font-semibold text-gray-800 dark:text-gray-100">Insight & Tips Keuangan</h2>
              </div>
              <Link
                to="/insights/analytics"
                className="flex items-center gap-1 text-xs text-indigo-500 hover:text-indigo-600 font-medium"
              >
                Selengkapnya <ChevronRight size={13} />
              </Link>
            </div>
            {analytics.insights?.length > 0 ? (
              <div className="space-y-2.5">
                {analytics.insights.slice(0, 4).map((text, i) => (
                  <InsightCard key={i} text={text} index={i} />
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center py-8 text-center">
                <div className="w-12 h-12 rounded-2xl bg-indigo-50 dark:bg-indigo-900/30 flex items-center justify-center mb-3">
                  <Lightbulb size={22} className="text-indigo-400" />
                </div>
                <p className="text-xs text-gray-400 dark:text-gray-500">Tambah lebih banyak transaksi untuk mendapat insight personal</p>
              </div>
            )}
          </div>
        </div>
      )}

      {/* ── Row 3: Area Chart + Donut Chart ─────────────────────── */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-5 stagger-4">

        {/* Area Chart — Tren Bulanan */}
        <div className="bg-white dark:bg-gray-800/90 rounded-2xl border border-gray-100 dark:border-gray-700/60 p-4 sm:p-6 card-hover">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-2">
              <BarChart3 size={18} className="text-indigo-500" />
              <h2 className="text-sm font-semibold text-gray-800 dark:text-gray-100">Tren Bulanan</h2>
            </div>
            <Link
              to="/insights/charts"
              className="flex items-center gap-1 text-xs text-indigo-500 hover:text-indigo-600 font-medium"
            >
              Detail <ChevronRight size={13} />
            </Link>
          </div>

          {monthlyChartData.length > 0 ? (
            <ResponsiveContainer width="100%" height={220}>
              <AreaChart data={monthlyChartData} margin={{ top: 5, right: 10, left: -10, bottom: 0 }}>
                <defs>
                  <linearGradient id="gradIncome" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#10b981" stopOpacity={0.2} />
                    <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
                  </linearGradient>
                  <linearGradient id="gradExpense" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#ef4444" stopOpacity={0.2} />
                    <stop offset="95%" stopColor="#ef4444" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" opacity={0.3} />
                <XAxis dataKey="name" tick={{ fontSize: 10 }} tickLine={false} axisLine={false} />
                <YAxis tick={{ fontSize: 10 }} tickLine={false} axisLine={false} tickFormatter={fmtC} width={48} />
                <Tooltip content={<CustomTooltip />} />
                <Legend wrapperStyle={{ fontSize: 11 }} />
                <Area type="monotone" dataKey="Pemasukan" stroke="#10b981" fill="url(#gradIncome)" strokeWidth={2.5} dot={{ r: 3, fill: '#10b981' }} />
                <Area type="monotone" dataKey="Pengeluaran" stroke="#ef4444" fill="url(#gradExpense)" strokeWidth={2.5} dot={{ r: 3, fill: '#ef4444' }} />
              </AreaChart>
            </ResponsiveContainer>
          ) : (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <div className="w-12 h-12 rounded-2xl bg-gray-50 dark:bg-gray-700/50 flex items-center justify-center mb-3">
                <BarChart3 size={22} className="text-gray-300 dark:text-gray-600" />
              </div>
              <p className="text-xs text-gray-400 dark:text-gray-500">Data tren belum tersedia</p>
            </div>
          )}
        </div>

        {/* Donut Chart — Kategori */}
        <div className="bg-white dark:bg-gray-800/90 rounded-2xl border border-gray-100 dark:border-gray-700/60 p-4 sm:p-6 card-hover">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-2">
              <CreditCard size={18} className="text-indigo-500" />
              <h2 className="text-sm font-semibold text-gray-800 dark:text-gray-100">Pengeluaran per Kategori</h2>
            </div>
            <Link
              to="/insights/tables"
              className="flex items-center gap-1 text-xs text-indigo-500 hover:text-indigo-600 font-medium"
            >
              Detail <ChevronRight size={13} />
            </Link>
          </div>

          {pieData.length > 0 ? (
            <>
              <ResponsiveContainer width="100%" height={180}>
                <PieChart>
                  <Pie
                    data={pieData}
                    cx="50%" cy="50%"
                    innerRadius={50} outerRadius={80}
                    paddingAngle={3}
                    dataKey="value"
                    strokeWidth={0}
                  >
                    {pieData.map((e, i) => <Cell key={i} fill={e.color || DEFAULT_COLORS[i % DEFAULT_COLORS.length]} />)}
                  </Pie>
                  <Tooltip
                    formatter={(v) => [`Rp ${fmt(v)}`, 'Total']}
                    contentStyle={{
                      backgroundColor: 'var(--tooltip-bg, #fff)',
                      border: '1px solid #e5e7eb',
                      borderRadius: '12px',
                      fontSize: 12,
                    }}
                  />
                </PieChart>
              </ResponsiveContainer>
              <div className="grid grid-cols-2 gap-x-4 gap-y-2 mt-2">
                {pieData.slice(0, 6).map((item, i) => (
                  <div key={i} className="flex items-center gap-2 text-xs">
                    <div className="w-2.5 h-2.5 rounded-full shrink-0" style={{ backgroundColor: item.color }} />
                    <span className="text-gray-600 dark:text-gray-300 truncate flex-1">{item.name}</span>
                    <span className="text-gray-400 dark:text-gray-500 shrink-0">{fmtC(item.value)}</span>
                  </div>
                ))}
              </div>
            </>
          ) : (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <div className="w-12 h-12 rounded-2xl bg-gray-50 dark:bg-gray-700/50 flex items-center justify-center mb-3">
                <TrendingUp size={22} className="text-gray-300 dark:text-gray-600" />
              </div>
              <p className="text-xs text-gray-400 dark:text-gray-500">Belum ada pengeluaran bulan ini</p>
            </div>
          )}
        </div>
      </div>

      {/* ── Row 4: Budget Progress + Savings Goals ──────────────── */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-5 stagger-5">

        {/* Budget Progress */}
        <div className="bg-white dark:bg-gray-800/90 rounded-2xl border border-gray-100 dark:border-gray-700/60 p-4 sm:p-6 card-hover">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-2">
              <Target size={18} className="text-indigo-500" />
              <h2 className="text-sm font-semibold text-gray-800 dark:text-gray-100">Anggaran Bulan Ini</h2>
            </div>
            <Link
              to="/budgets"
              className="flex items-center gap-1 text-xs text-indigo-500 hover:text-indigo-600 font-medium"
            >
              Kelola <ChevronRight size={13} />
            </Link>
          </div>

          {budgetItems.length > 0 ? (
            <div className="space-y-4">
              {budgetItems.slice(0, 5).map((b) => (
                <BudgetItem
                  key={b.id}
                  name={b.category_name}
                  color={b.category_color}
                  spent={b.spent}
                  amount={b.amount}
                />
              ))}
              {budgetData && (
                <div className="pt-3 border-t border-gray-100 dark:border-gray-700 flex items-center justify-between text-xs">
                  <span className="text-gray-400 dark:text-gray-500">Total Anggaran</span>
                  <span className="font-bold text-gray-700 dark:text-gray-200">
                    Rp {fmtC(budgetData.total_spent)} / Rp {fmtC(budgetData.total_budget)}
                  </span>
                </div>
              )}
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center py-10 text-center">
              <div className="w-12 h-12 rounded-2xl bg-gray-50 dark:bg-gray-700/50 flex items-center justify-center mb-3">
                <Target size={22} className="text-gray-300 dark:text-gray-600" />
              </div>
              <p className="text-xs text-gray-400 dark:text-gray-500">Belum ada anggaran bulan ini</p>
              <Link to="/budgets" className="text-xs text-indigo-500 hover:text-indigo-600 mt-1.5 font-medium">
                Buat anggaran →
              </Link>
            </div>
          )}
        </div>

        {/* Savings Goals */}
        <div className="bg-white dark:bg-gray-800/90 rounded-2xl border border-gray-100 dark:border-gray-700/60 p-4 sm:p-6 card-hover">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-2">
              <PiggyBank size={18} className="text-indigo-500" />
              <h2 className="text-sm font-semibold text-gray-800 dark:text-gray-100">Target Tabungan</h2>
            </div>
            <Link
              to="/savings"
              className="flex items-center gap-1 text-xs text-indigo-500 hover:text-indigo-600 font-medium"
            >
              Kelola <ChevronRight size={13} />
            </Link>
          </div>

          {savings.length > 0 ? (
            <div className="space-y-2.5">
              {savings.slice(0, 4).map((s) => (
                <SavingsItem
                  key={s.id}
                  name={s.name}
                  current={s.current_amount}
                  target={s.target_amount}
                  color={s.color}
                  deadline={s.deadline}
                />
              ))}
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center py-10 text-center">
              <div className="w-12 h-12 rounded-2xl bg-gray-50 dark:bg-gray-700/50 flex items-center justify-center mb-3">
                <PiggyBank size={22} className="text-gray-300 dark:text-gray-600" />
              </div>
              <p className="text-xs text-gray-400 dark:text-gray-500">Belum ada target tabungan</p>
              <Link to="/savings" className="text-xs text-indigo-500 hover:text-indigo-600 mt-1.5 font-medium">
                Buat target →
              </Link>
            </div>
          )}
        </div>
      </div>

      {/* ── Row 5: Recent Transactions ──────────────────────────── */}
      <div className="bg-white dark:bg-gray-800/90 rounded-2xl border border-gray-100 dark:border-gray-700/60 p-4 sm:p-6 card-hover stagger-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <Activity size={18} className="text-indigo-500" />
            <h2 className="text-sm font-semibold text-gray-800 dark:text-gray-100">Transaksi Terbaru</h2>
          </div>
          <Link
            to="/transactions"
            className="flex items-center gap-1 text-xs text-indigo-500 hover:text-indigo-600 font-medium"
          >
            Lihat semua <ChevronRight size={13} />
          </Link>
        </div>

        {recent.length > 0 ? (
          <div className="divide-y divide-gray-50 dark:divide-gray-700/50">
            {recent.map((t) => (
              <div
                key={t.id}
                className="flex items-center gap-3 py-3 first:pt-0 last:pb-0 group/tx"
              >
                {/* Category icon */}
                <div className={`w-10 h-10 rounded-xl flex items-center justify-center shrink-0 transition-transform group-hover/tx:scale-105 ${
                  t.type === 'income'
                    ? 'bg-emerald-50 dark:bg-emerald-900/30'
                    : 'bg-red-50 dark:bg-red-900/20'
                }`}>
                  {t.type === 'income'
                    ? <ArrowUpCircle size={18} className="text-emerald-500" />
                    : <ArrowDownCircle size={18} className="text-red-400" />}
                </div>

                {/* Description + meta */}
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-medium text-gray-700 dark:text-gray-200 truncate leading-tight">
                    {t.description || t.category_name}
                  </p>
                  <div className="flex items-center gap-1.5 mt-0.5">
                    <span className="text-[10px] px-1.5 py-0.5 rounded-md bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 font-medium">
                      {t.category_name}
                    </span>
                    <span className="text-[10px] text-gray-400 dark:text-gray-500">
                      {new Date(t.date + 'T00:00:00').toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })}
                    </span>
                  </div>
                </div>

                {/* Amount */}
                <span className={`text-sm font-bold shrink-0 ${
                  t.type === 'income'
                    ? 'text-emerald-600 dark:text-emerald-400'
                    : 'text-red-500 dark:text-red-400'
                }`}>
                  {t.type === 'income' ? '+' : '-'}
                  <span className="sm:hidden">{fmtC(t.amount)}</span>
                  <span className="hidden sm:inline">{fmt(t.amount)}</span>
                </span>
              </div>
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center py-10 text-center">
            <div className="w-12 h-12 rounded-2xl bg-gray-50 dark:bg-gray-700/50 flex items-center justify-center mb-3">
              <Wallet size={22} className="text-gray-300 dark:text-gray-600" />
            </div>
            <p className="text-xs text-gray-400 dark:text-gray-500">Belum ada transaksi bulan ini</p>
            <Link to="/transactions" className="text-xs text-indigo-500 hover:text-indigo-600 mt-1.5 font-medium">
              Tambah transaksi →
            </Link>
          </div>
        )}
      </div>

      {/* ── Footer note ─────────────────────────────────────────── */}
      <p className="text-[10px] text-gray-300 dark:text-gray-600 text-center pb-1">
        Data diperbarui secara real-time · {monthName}
      </p>
    </div>
  );
}
