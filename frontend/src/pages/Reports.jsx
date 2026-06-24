import { useState, useEffect } from 'react';
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  PieChart, Pie, Cell, Legend, ComposedChart, Line, Area
} from 'recharts';
import Card from '../components/Card';
import { getSummary, getYearlyTrend, getTransactions } from '../api/client';
import { Download, Printer } from 'lucide-react';

const MONTH_NAMES = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des'];

function formatRupiah(n) {
  return new Intl.NumberFormat('id-ID').format(n || 0);
}

function formatCompact(n) {
  if (Math.abs(n) >= 1000000) return `${(n / 1000000).toFixed(n % 1000000 === 0 ? 0 : 1)}jt`;
  if (Math.abs(n) >= 1000) return `${(n / 1000).toFixed(0)}rb`;
  return new Intl.NumberFormat('id-ID').format(n);
}

function CashFlowTooltip({ active, payload, label }) {
  if (!active || !payload?.length) return null;
  return (
    <div className="bg-white/95 dark:bg-gray-800/95 backdrop-blur-xl rounded-xl shadow-xl border border-gray-100 dark:border-gray-700 p-3 text-sm">
      <p className="font-semibold text-gray-700 dark:text-gray-200 mb-1.5">{label}</p>
      {payload.map((entry, i) => (
        <div key={i} className="flex items-center gap-2 py-0.5">
          <div className="w-2.5 h-2.5 rounded-full" style={{ backgroundColor: entry.color }} />
          <span className="text-gray-500 dark:text-gray-400">{entry.name}:</span>
          <span className="font-semibold text-gray-800 dark:text-gray-100">Rp {formatRupiah(entry.value)}</span>
        </div>
      ))}
    </div>
  );
}

export default function Reports() {
  const now = new Date();
  const [month, setMonth] = useState(now.getMonth() + 1);
  const [year, setYear] = useState(now.getFullYear());
  const [summary, setSummary] = useState(null);
  const [loading, setLoading] = useState(true);

  const [yearlyTrend, setYearlyTrend] = useState(null);
  const [exporting, setExporting] = useState(false);

  const handleExportCSV = async () => {
    try {
      setExporting(true);
      const res = await getTransactions({ month, year, limit: 1000 });
      const txs = res.transactions || [];
      if (txs.length === 0) {
        alert('Tidak ada transaksi untuk diekspor pada periode ini.');
        return;
      }

      // Build CSV content
      const headers = ['ID', 'Tanggal', 'Kategori', 'Deskripsi', 'Tipe', 'Jumlah (Rp)'];
      const rows = txs.map(t => [
        t.id,
        t.date,
        t.category_name || '',
        `"${(t.description || '').replace(/"/g, '""')}"`,
        t.type === 'income' ? 'Pemasukan' : 'Pengeluaran',
        t.amount
      ]);

      const csvContent = "data:text/csv;charset=utf-8," 
        + [headers.join(','), ...rows.map(r => r.join(','))].join('\n');
      
      const encodedUri = encodeURI(csvContent);
      const link = document.createElement("a");
      link.setAttribute("href", encodedUri);
      link.setAttribute("download", `Laporan_Keuangan_${month}_${year}.csv`);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch (err) {
      console.error(err);
      alert('Gagal mengekspor data.');
    } finally {
      setExporting(false);
    }
  };

  const handlePrint = () => {
    window.print();
  };

  useEffect(() => {
    setLoading(true);
    Promise.all([getSummary(month, year), getYearlyTrend(year)])
      .then(([s, yt]) => { setSummary(s); setYearlyTrend(yt); })
      .finally(() => setLoading(false));
  }, [month, year]);

  const formatCurrency = formatRupiah;

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
        <h1 className="text-xl sm:text-2xl font-bold text-gray-800 dark:text-gray-100">Laporan</h1>
        <div className="flex gap-2.5 sm:gap-3 print:hidden">
          <select
            value={month}
            onChange={(e) => setMonth(parseInt(e.target.value))}
            className="border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-700 dark:text-gray-200 flex-1 sm:flex-none"
          >
            {Array.from({ length: 12 }, (_, i) => i + 1).map((m) => (
              <option key={m} value={m}>{new Date(2000, m - 1).toLocaleString('id', { month: 'long' })}</option>
            ))}
          </select>
          <select
            value={year}
            onChange={(e) => setYear(parseInt(e.target.value))}
            className="border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-700 dark:text-gray-200 flex-1 sm:flex-none"
          >
            {[year - 2, year - 1, year, year + 1].map((y) => (
              <option key={y} value={y}>{y}</option>
            ))}
          </select>
          <button
            onClick={handleExportCSV}
            disabled={exporting}
            className="flex items-center gap-1.5 border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-2 text-sm bg-white dark:bg-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors disabled:opacity-50"
            title="Ekspor ke CSV"
          >
            <Download size={16} />
            <span className="hidden md:inline">Ekspor CSV</span>
          </button>
          <button
            onClick={handlePrint}
            className="flex items-center gap-1.5 bg-blue-600 text-white rounded-xl px-3 py-2 text-sm hover:bg-blue-700 transition-colors"
            title="Cetak Laporan"
          >
            <Printer size={16} />
            <span className="hidden md:inline">Cetak</span>
          </button>
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

      {/* Cash Flow Chart */}
      <Card>
        <h2 className="text-sm sm:text-lg font-semibold text-gray-800 dark:text-gray-100 mb-3 sm:mb-4">Cash Flow {year}</h2>
        {(() => {
          let cumulative = 0;
          const cfData = yearlyTrend?.months?.map((m) => {
            cumulative += (m.income - m.expense);
            return { name: MONTH_NAMES[m.month - 1], Pemasukan: m.income, Pengeluaran: m.expense, 'Cash Flow': cumulative };
          }) || [];
          if (!cfData.some(d => d.Pemasukan > 0 || d.Pengeluaran > 0)) {
            return <p className="text-gray-400 dark:text-gray-500 text-center py-8 sm:py-12 text-sm">Tidak ada data cash flow</p>;
          }
          return (
            <>
              <div className="flex items-center gap-3 text-xs mb-3">
                <span className="flex items-center gap-1.5 text-gray-600 dark:text-gray-300"><span className="w-3 h-3 rounded bg-green-400" /> Masuk</span>
                <span className="flex items-center gap-1.5 text-gray-600 dark:text-gray-300"><span className="w-3 h-3 rounded bg-red-400" /> Keluar</span>
                <span className="flex items-center gap-1.5 text-gray-600 dark:text-gray-300"><span className="w-3 h-1.5 rounded-full bg-indigo-500" /> Akumulasi</span>
              </div>
              <ResponsiveContainer width="100%" height={220} className="sm:hidden">
                <ComposedChart data={cfData} margin={{ top: 5, right: 5, bottom: 5, left: 5 }}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" vertical={false} />
                  <XAxis dataKey="name" tick={{ fontSize: 10 }} tickLine={false} axisLine={false} />
                  <YAxis tick={{ fontSize: 10 }} tickFormatter={formatCompact} tickLine={false} axisLine={false} />
                  <Tooltip content={<CashFlowTooltip />} />
                  <Bar dataKey="Pemasukan" fill="#10b981" radius={[4, 4, 0, 0]} barSize={14} />
                  <Bar dataKey="Pengeluaran" fill="#ef4444" radius={[4, 4, 0, 0]} barSize={14} />
                  <Area type="monotone" dataKey="Cash Flow" stroke="#6366f1" strokeWidth={2} fill="#6366f1" fillOpacity={0.15} dot={{ r: 2.5, fill: '#6366f1', strokeWidth: 0 }} />
                </ComposedChart>
              </ResponsiveContainer>
              <ResponsiveContainer width="100%" height={300} className="hidden sm:block">
                <ComposedChart data={cfData} margin={{ top: 5, right: 5, bottom: 5, left: 5 }}>
                  <defs>
                    <linearGradient id="rpIncomeGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="0%" stopColor="#10b981" stopOpacity={0.9} />
                      <stop offset="100%" stopColor="#10b981" stopOpacity={0.6} />
                    </linearGradient>
                    <linearGradient id="rpExpenseGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="0%" stopColor="#ef4444" stopOpacity={0.9} />
                      <stop offset="100%" stopColor="#ef4444" stopOpacity={0.6} />
                    </linearGradient>
                    <linearGradient id="rpCashFlowGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="0%" stopColor="#6366f1" stopOpacity={0.3} />
                      <stop offset="100%" stopColor="#6366f1" stopOpacity={0.05} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" vertical={false} />
                  <XAxis dataKey="name" tick={{ fontSize: 11 }} tickLine={false} axisLine={false} />
                  <YAxis tick={{ fontSize: 10 }} tickFormatter={formatCompact} tickLine={false} axisLine={false} />
                  <Tooltip content={<CashFlowTooltip />} />
                  <Bar dataKey="Pemasukan" fill="url(#rpIncomeGrad)" radius={[4, 4, 0, 0]} barSize={20} />
                  <Bar dataKey="Pengeluaran" fill="url(#rpExpenseGrad)" radius={[4, 4, 0, 0]} barSize={20} />
                  <Area type="monotone" dataKey="Cash Flow" stroke="#6366f1" strokeWidth={2.5} fill="url(#rpCashFlowGrad)" dot={{ r: 3, fill: '#6366f1', strokeWidth: 0 }} activeDot={{ r: 5, strokeWidth: 2, stroke: '#fff' }} />
                </ComposedChart>
              </ResponsiveContainer>
            </>
          );
        })()}
      </Card>

      {/* Charts */}
      <div className="grid grid-cols-1 gap-4 sm:gap-6">
        {/* Expense Pie Chart */}
        <Card>
          <h2 className="text-sm sm:text-lg font-semibold text-gray-800 dark:text-gray-100 mb-3 sm:mb-4">Pengeluaran per Kategori</h2>
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
                      <span className="text-xs text-gray-600 dark:text-gray-300 truncate">{item.name}</span>
                    </div>
                    <span className="text-xs font-medium text-gray-700 dark:text-gray-200 ml-2 shrink-0">Rp {formatCompact(item.value)}</span>
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
            <p className="text-gray-400 dark:text-gray-500 text-center py-8 sm:py-12 text-sm">Tidak ada data pengeluaran</p>
          )}
        </Card>

        {/* Bar Chart */}
        <Card>
          <h2 className="text-sm sm:text-lg font-semibold text-gray-800 dark:text-gray-100 mb-3 sm:mb-4">Perbandingan per Kategori</h2>
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
            <p className="text-gray-400 dark:text-gray-500 text-center py-8 sm:py-12 text-sm">Tidak ada data</p>
          )}
        </Card>
      </div>

      {/* Detail by Category - Mobile cards, Desktop table */}
      <Card>
        <h2 className="text-sm sm:text-lg font-semibold text-gray-800 dark:text-gray-100 mb-3 sm:mb-4">Detail per Kategori</h2>
        {summary?.by_category?.length > 0 ? (
          <>
            {/* Mobile: Card list */}
            <div className="sm:hidden space-y-2">
              {summary.by_category.map((c) => (
                <div key={c.category_id} className="flex items-center justify-between py-2.5 border-b border-gray-50 dark:border-gray-700 last:border-0">
                  <div className="flex items-center gap-2.5 min-w-0 flex-1">
                    <div className="w-8 h-8 rounded-lg flex items-center justify-center shrink-0" style={{ backgroundColor: c.category_color + '15' }}>
                      <div className="w-2.5 h-2.5 rounded-full" style={{ backgroundColor: c.category_color }} />
                    </div>
                    <div className="min-w-0">
                      <p className="text-sm font-medium text-gray-700 dark:text-gray-200 truncate">{c.category_name}</p>
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
                    <tr key={c.category_id} className="border-b border-gray-50 dark:border-gray-700">
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2">
                          <div className="w-3 h-3 rounded-full" style={{ backgroundColor: c.category_color }} />
                          <span className="text-gray-700 dark:text-gray-200">{c.category_name}</span>
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
          <p className="text-gray-400 dark:text-gray-500 text-center py-8 text-sm">Belum ada data bulan ini</p>
        )}
      </Card>
    </div>
  );
}
