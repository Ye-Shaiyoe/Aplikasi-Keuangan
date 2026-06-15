import { useState, useEffect } from 'react';
import { Wallet, ArrowDownCircle, ArrowUpCircle } from 'lucide-react';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts';
import { StatCard } from '../components/Card';
import { getSummary, getTransactions } from '../api/client';

export default function Dashboard() {
  const [summary, setSummary] = useState(null);
  const [recent, setRecent] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const now = new Date();
    const month = now.getMonth() + 1;
    const year = now.getFullYear();

    Promise.all([
      getSummary(month, year),
      getTransactions({ limit: 10, month, year }),
    ]).then(([s, t]) => {
      setSummary(s);
      setRecent(t.data || []);
    }).finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
      </div>
    );
  }

  const expenseByCategory = summary?.by_category
    ?.filter((c) => c.total > 0)
    .filter((c) => recent.some((r) => r.category_id === c.category_id && r.type === 'expense'))
    .length > 0
    ? summary.by_category.filter((c) => c.total > 0)
    : summary?.by_category?.slice(0, 5) || [];

  const pieData = expenseByCategory
    .filter((c) => c.total > 0)
    .map((c) => ({
      name: c.category_name,
      value: c.total,
      color: c.category_color,
    }));

  const formatCurrency = (val) =>
    new Intl.NumberFormat('id-ID').format(val || 0);

  const formatCompact = (val) => {
    if (val >= 1000000) return `${(val / 1000000).toFixed(val % 1000000 === 0 ? 0 : 1)}jt`;
    if (val >= 1000) return `${(val / 1000).toFixed(0)}rb`;
    return new Intl.NumberFormat('id-ID').format(val);
  };

  const monthName = new Date().toLocaleString('id', { month: 'long' });

  return (
    <div className="space-y-5 sm:space-y-6">
      {/* Header - mobile has greeting, desktop has title */}
      <div className="hidden sm:block">
        <h1 className="text-2xl font-bold text-gray-800">Dashboard</h1>
      </div>

      {/* Mobile greeting */}
      <div className="sm:hidden">
        <p className="text-gray-400 text-sm">Ringkasan {monthName}</p>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-3 sm:grid-cols-2 lg:grid-cols-3 gap-2.5 sm:gap-4">
        <StatCard
          title="Saldo"
          value={summary?.balance || 0}
          icon={Wallet}
          color="blue"
        />
        <StatCard
          title="Masuk"
          value={summary?.total_income || 0}
          icon={ArrowUpCircle}
          color="green"
        />
        <StatCard
          title="Keluar"
          value={summary?.total_expense || 0}
          icon={ArrowDownCircle}
          color="red"
        />
      </div>

      {/* Mobile: Full-width stacked / Desktop: Side by side */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-6">
        {/* Pie Chart */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-4 sm:p-6">
          <h2 className="text-sm sm:text-lg font-semibold text-gray-800 mb-3 sm:mb-4">Pengeluaran per Kategori</h2>
          {pieData.length > 0 ? (
            <ResponsiveContainer width="100%" height={200} className="sm:hidden">
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  innerRadius={45}
                  outerRadius={75}
                  paddingAngle={3}
                  dataKey="value"
                >
                  {pieData.map((entry, i) => (
                    <Cell key={i} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip
                  formatter={(value) => [`Rp ${formatCurrency(value)}`, 'Total']}
                />
              </PieChart>
            </ResponsiveContainer>
          ) : null}
          {pieData.length > 0 ? (
            <ResponsiveContainer width="100%" height={280} className="hidden sm:block">
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={100}
                  paddingAngle={3}
                  dataKey="value"
                >
                  {pieData.map((entry, i) => (
                    <Cell key={i} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip
                  formatter={(value) => [`Rp ${formatCurrency(value)}`, 'Total']}
                />
              </PieChart>
            </ResponsiveContainer>
          ) : (
            <p className="text-gray-400 text-center py-8 sm:py-12 text-sm">
              Belum ada data pengeluaran bulan ini
            </p>
          )}
          {/* Legend - compact on mobile */}
          {pieData.length > 0 && (
            <div className="grid grid-cols-2 gap-x-3 gap-y-1.5 sm:gap-2 mt-3 sm:mt-4">
              {pieData.map((item, i) => (
                <div key={i} className="flex items-center gap-1.5 sm:gap-2 text-xs sm:text-sm">
                  <div
                    className="w-2.5 h-2.5 sm:w-3 sm:h-3 rounded-full flex-shrink-0"
                    style={{ backgroundColor: item.color }}
                  />
                  <span className="text-gray-600 truncate flex-1">{item.name}</span>
                  <span className="text-gray-400 text-[10px] sm:text-xs hidden sm:inline">
                    Rp {formatCompact(item.value)}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Recent Transactions */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-4 sm:p-6">
          <h2 className="text-sm sm:text-lg font-semibold text-gray-800 mb-3 sm:mb-4">Transaksi Terbaru</h2>
          {recent.length > 0 ? (
            <div className="space-y-2.5 sm:space-y-3">
              {recent.map((t) => (
                <div key={t.id} className="flex items-center justify-between py-2 border-b border-gray-50 last:border-0">
                  <div className="flex items-center gap-2.5 sm:gap-3 min-w-0 flex-1">
                    <div
                      className={`w-9 h-9 sm:w-10 sm:h-10 rounded-xl flex items-center justify-center shrink-0 ${
                        t.type === 'income' ? 'bg-green-50' : 'bg-red-50'
                      }`}
                    >
                      {t.type === 'income' ? (
                        <ArrowUpCircle size={18} className="text-green-500" />
                      ) : (
                        <ArrowDownCircle size={18} className="text-red-500" />
                      )}
                    </div>
                    <div className="min-w-0">
                      <p className="text-sm font-medium text-gray-700 truncate">
                        {t.description || t.category_name}
                      </p>
                      <p className="text-[11px] sm:text-xs text-gray-400 truncate">{t.category_name} · {new Date(t.date + 'T00:00:00').toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })}</p>
                    </div>
                  </div>
                  <span
                    className={`text-sm font-semibold shrink-0 ml-2 ${
                      t.type === 'income' ? 'text-green-600' : 'text-red-600'
                    }`}
                  >
                    {t.type === 'income' ? '+' : '-'}<span className="sm:hidden">{formatCompact(t.amount)}</span><span className="hidden sm:inline">{formatCurrency(t.amount)}</span>
                  </span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-400 text-center py-8 sm:py-12 text-sm">
              Belum ada transaksi bulan ini
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
