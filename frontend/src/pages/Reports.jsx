import { useState, useEffect } from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend } from 'recharts';
import Card from '../components/Card';
import { getSummary } from '../api/client';

export default function Reports() {
  const now = new Date();
  const [month, setMonth] = useState(now.getMonth() + 1);
  const [year, setYear] = useState(now.getFullYear());
  const [summary, setSummary] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    getSummary(month, year).then(setSummary).finally(() => setLoading(false));
  }, [month, year]);

  const formatCurrency = (val) => new Intl.NumberFormat('id-ID').format(val || 0);

  const formatCompact = (val) => {
    if (val >= 1000000) return `${(val / 1000000).toFixed(val % 1000000 === 0 ? 0 : 1)}jt`;
    if (val >= 1000) return `${(val / 1000).toFixed(0)}rb`;
    return new Intl.NumberFormat('id-ID').format(val);
  };

  const barData = summary?.by_category
    ?.filter((c) => c.total > 0)
    .map((c) => ({
      name: c.category_name,
      Pemasukan: c.category_name.includes('Pemasukan') ? c.total : 0,
      Pengeluaran: c.category_name.includes('Pemasukan') ? 0 : c.total,
      fill: c.category_color,
    })) || [];

  // Separate expense categories for pie
  const expensePie = summary?.by_category
    ?.filter((c) => c.total > 0 && !c.category_name.includes('Pemasukan'))
    .map((c) => ({
      name: c.category_name,
      value: c.total,
      color: c.category_color,
    })) || [];

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
      </div>
    );
  }

  return (
    <div className="space-y-4 sm:space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 sm:gap-4">
        <h1 className="text-xl sm:text-2xl font-bold text-gray-800">Laporan</h1>
        <div className="flex gap-2.5 sm:gap-3">
          <select
            value={month}
            onChange={(e) => setMonth(parseInt(e.target.value))}
            className="border border-gray-200 rounded-xl px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white flex-1 sm:flex-none"
          >
            {Array.from({ length: 12 }, (_, i) => i + 1).map((m) => (
              <option key={m} value={m}>{new Date(2000, m - 1).toLocaleString('id', { month: 'long' })}</option>
            ))}
          </select>
          <select
            value={year}
            onChange={(e) => setYear(parseInt(e.target.value))}
            className="border border-gray-200 rounded-xl px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white flex-1 sm:flex-none"
          >
            {[year - 2, year - 1, year, year + 1].map((y) => (
              <option key={y} value={y}>{y}</option>
            ))}
          </select>
        </div>
      </div>

      {/* Summary Cards - stacked on mobile, grid on desktop */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 sm:gap-4">
        <Card>
          <p className="text-xs sm:text-sm text-gray-500">Total Pemasukan</p>
          <p className="text-lg sm:text-2xl font-bold text-green-600 mt-0.5 sm:mt-1">
            Rp <span className="sm:hidden">{formatCompact(summary?.total_income || 0)}</span><span className="hidden sm:inline">{formatCurrency(summary?.total_income || 0)}</span>
          </p>
        </Card>
        <Card>
          <p className="text-xs sm:text-sm text-gray-500">Total Pengeluaran</p>
          <p className="text-lg sm:text-2xl font-bold text-red-600 mt-0.5 sm:mt-1">
            Rp <span className="sm:hidden">{formatCompact(summary?.total_expense || 0)}</span><span className="hidden sm:inline">{formatCurrency(summary?.total_expense || 0)}</span>
          </p>
        </Card>
        <Card>
          <p className="text-xs sm:text-sm text-gray-500">Saldo</p>
          <p className={`text-lg sm:text-2xl font-bold mt-0.5 sm:mt-1 ${(summary?.balance || 0) >= 0 ? 'text-blue-600' : 'text-red-600'}`}>
            Rp <span className="sm:hidden">{formatCompact(summary?.balance || 0)}</span><span className="hidden sm:inline">{formatCurrency(summary?.balance || 0)}</span>
          </p>
        </Card>
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 gap-4 sm:gap-6">
        {/* Expense Pie Chart */}
        <Card>
          <h2 className="text-sm sm:text-lg font-semibold text-gray-800 mb-3 sm:mb-4">Pengeluaran per Kategori</h2>
          {expensePie.length > 0 ? (
            <>
              <ResponsiveContainer width="100%" height={220} className="sm:hidden">
                <PieChart>
                  <Pie
                    data={expensePie}
                    cx="50%"
                    cy="50%"
                    innerRadius={45}
                    outerRadius={80}
                    paddingAngle={2}
                    dataKey="value"
                  >
                    {expensePie.map((entry, i) => (
                      <Cell key={i} fill={entry.color} />
                    ))}
                  </Pie>
                  <Tooltip formatter={(value) => [`Rp ${formatCurrency(value)}`, 'Total']} />
                </PieChart>
              </ResponsiveContainer>
              {/* Mobile legend below chart */}
              <div className="sm:hidden space-y-2 mt-3">
                {expensePie.map((item, i) => (
                  <div key={i} className="flex items-center justify-between">
                    <div className="flex items-center gap-2 min-w-0 flex-1">
                      <div className="w-3 h-3 rounded-full shrink-0" style={{ backgroundColor: item.color }} />
                      <span className="text-xs text-gray-600 truncate">{item.name}</span>
                    </div>
                    <span className="text-xs font-medium text-gray-700 ml-2 shrink-0">Rp {formatCompact(item.value)}</span>
                  </div>
                ))}
              </div>
              {/* Desktop chart */}
              <ResponsiveContainer width="100%" height={300} className="hidden sm:block">
                <PieChart>
                  <Pie
                    data={expensePie}
                    cx="50%"
                    cy="50%"
                    outerRadius={100}
                    dataKey="value"
                    label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  >
                    {expensePie.map((entry, i) => (
                      <Cell key={i} fill={entry.color} />
                    ))}
                  </Pie>
                  <Tooltip formatter={(value) => [`Rp ${formatCurrency(value)}`, 'Total']} />
                </PieChart>
              </ResponsiveContainer>
            </>
          ) : (
            <p className="text-gray-400 text-center py-8 sm:py-12 text-sm">Tidak ada data pengeluaran</p>
          )}
        </Card>

        {/* Bar Chart */}
        <Card>
          <h2 className="text-sm sm:text-lg font-semibold text-gray-800 mb-3 sm:mb-4">Perbandingan per Kategori</h2>
          {barData.length > 0 ? (
            <ResponsiveContainer width="100%" height={220} className="sm:hidden">
              <BarChart data={barData} margin={{ top: 5, right: 5, bottom: 5, left: 5 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                <XAxis dataKey="name" tick={{ fontSize: 10 }} interval={0} angle={-30} textAnchor="end" height={60} />
                <YAxis tick={{ fontSize: 10 }} tickFormatter={(v) => formatCompact(v)} />
                <Tooltip formatter={(value) => [`Rp ${formatCurrency(value)}`, '']} />
                <Bar dataKey="Pengeluaran" fill="#ef4444" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          ) : null}
          {barData.length > 0 ? (
            <ResponsiveContainer width="100%" height={300} className="hidden sm:block">
              <BarChart data={barData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                <XAxis dataKey="name" tick={{ fontSize: 12 }} />
                <YAxis tick={{ fontSize: 12 }} />
                <Tooltip formatter={(value) => [`Rp ${formatCurrency(value)}`, '']} />
                <Bar dataKey="Pengeluaran" fill="#ef4444" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <p className="text-gray-400 text-center py-8 sm:py-12 text-sm">Tidak ada data</p>
          )}
        </Card>
      </div>

      {/* Detail by Category - Mobile cards, Desktop table */}
      <Card>
        <h2 className="text-sm sm:text-lg font-semibold text-gray-800 mb-3 sm:mb-4">Detail per Kategori</h2>
        {summary?.by_category?.length > 0 ? (
          <>
            {/* Mobile: Card list */}
            <div className="sm:hidden space-y-2">
              {summary.by_category.map((c) => (
                <div key={c.category_id} className="flex items-center justify-between py-2.5 border-b border-gray-50 last:border-0">
                  <div className="flex items-center gap-2.5 min-w-0 flex-1">
                    <div className="w-8 h-8 rounded-lg flex items-center justify-center shrink-0" style={{ backgroundColor: c.category_color + '15' }}>
                      <div className="w-2.5 h-2.5 rounded-full" style={{ backgroundColor: c.category_color }} />
                    </div>
                    <div className="min-w-0">
                      <p className="text-sm font-medium text-gray-700 truncate">{c.category_name}</p>
                      <p className="text-xs text-gray-400">{c.count} transaksi</p>
                    </div>
                  </div>
                  <span className={`text-sm font-bold shrink-0 ml-2 ${
                    c.category_name.includes('Pemasukan') ? 'text-green-600' : 'text-red-600'
                  }`}>
                    Rp {formatCompact(c.total)}
                  </span>
                </div>
              ))}
            </div>

            {/* Desktop: Table */}
            <div className="hidden sm:block overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-100">
                    <th className="text-left px-4 py-3 font-medium text-gray-500">Kategori</th>
                    <th className="text-right px-4 py-3 font-medium text-gray-500">Jumlah Transaksi</th>
                    <th className="text-right px-4 py-3 font-medium text-gray-500">Total</th>
                  </tr>
                </thead>
                <tbody>
                  {summary.by_category.map((c) => (
                    <tr key={c.category_id} className="border-b border-gray-50">
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2">
                          <div className="w-3 h-3 rounded-full" style={{ backgroundColor: c.category_color }} />
                          <span className="text-gray-700">{c.category_name}</span>
                        </div>
                      </td>
                      <td className="px-4 py-3 text-right text-gray-600">{c.count}</td>
                      <td className={`px-4 py-3 text-right font-semibold ${
                        c.category_name.includes('Pemasukan') ? 'text-green-600' : 'text-red-600'
                      }`}>
                        Rp {formatCurrency(c.total)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </>
        ) : (
          <p className="text-gray-400 text-center py-8 text-sm">Belum ada data bulan ini</p>
        )}
      </Card>
    </div>
  );
}
