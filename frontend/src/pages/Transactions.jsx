import { useState, useEffect } from 'react';
import { Plus, Pencil, Trash2, Filter, ChevronDown, ChevronUp, X } from 'lucide-react';
import Modal from '../components/Modal';
import { getTransactions, getCategories, createTransaction, updateTransaction, deleteTransaction } from '../api/client';

export default function Transactions() {
  const [transactions, setTransactions] = useState([]);
  const [categories, setCategories] = useState([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [editItem, setEditItem] = useState(null);
  const [filters, setFilters] = useState({ month: '', year: '', type: '', category_id: '' });
  const [form, setForm] = useState({ category_id: '', amount: '', description: '', date: '', type: 'expense' });
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState('');
  const [showFilters, setShowFilters] = useState(false);

  const limit = 15;

  useEffect(() => {
    getCategories().then(setCategories);
  }, []);

  useEffect(() => {
    setLoading(true);
    const params = { page, limit };
    if (filters.month) params.month = parseInt(filters.month);
    if (filters.year) params.year = parseInt(filters.year);
    if (filters.type) params.type = filters.type;
    if (filters.category_id) params.category_id = parseInt(filters.category_id);

    getTransactions(params).then((res) => {
      setTransactions(res.data || []);
      setTotal(res.total || 0);
    }).finally(() => setLoading(false));
  }, [page, filters]);

  const openAdd = () => {
    setEditItem(null);
    setFormError('');
    setForm({ category_id: '', amount: '', description: '', date: new Date().toISOString().split('T')[0], type: 'expense' });
    setModalOpen(true);
  };

  const openEdit = (t) => {
    setEditItem(t);
    setFormError('');
    setForm({
      category_id: t.category_id,
      amount: t.amount,
      description: t.description,
      date: t.date,
      type: t.type,
    });
    setModalOpen(true);
  };

  const handleSave = async () => {
    if (!form.category_id) { setFormError('Pilih kategori terlebih dahulu'); return; }
    if (!form.amount || parseInt(form.amount) <= 0) { setFormError('Masukkan jumlah yang valid'); return; }
    if (!form.date) { setFormError('Pilih tanggal'); return; }

    setFormError('');
    setSaving(true);

    const payload = {
      ...form,
      amount: parseInt(form.amount),
      category_id: parseInt(form.category_id),
    };

    try {
      if (editItem) {
        const updated = await updateTransaction(editItem.id, payload);
        setTransactions(transactions.map((t) => (t.id === editItem.id ? { ...t, ...updated } : t)));
      } else {
        const created = await createTransaction(payload);
        setTransactions([created, ...transactions]);
      }
      setModalOpen(false);
    } catch (err) {
      const msg = err?.response?.data?.error || 'Gagal menyimpan transaksi';
      setFormError(msg);
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Hapus transaksi ini?')) return;
    await deleteTransaction(id);
    setTransactions(transactions.filter((t) => t.id !== id));
  };

  const totalPages = Math.ceil(total / limit);

  const formatCurrency = (val) => new Intl.NumberFormat('id-ID').format(val || 0);

  const categoryById = {};
  categories.forEach((c) => { categoryById[c.id] = c; });

  const now = new Date();
  const currentYear = now.getFullYear();

  const activeFilterCount = [filters.month, filters.year, filters.type, filters.category_id].filter(Boolean).length;

  const clearFilters = () => {
    setFilters({ month: '', year: '', type: '', category_id: '' });
    setPage(1);
  };

  return (
    <div className="space-y-4 sm:space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl sm:text-2xl font-bold text-gray-800 dark:text-gray-100">Transaksi</h1>
        <button
          onClick={openAdd}
          className="hidden sm:flex items-center gap-2 bg-blue-600 text-white px-4 py-2.5 rounded-xl hover:bg-blue-700 transition-colors text-sm font-medium"
        >
          <Plus size={18} />
          Tambah Transaksi
        </button>
      </div>

      {/* Filter toggle - Mobile */}
      <div className="sm:hidden">
        <button
          onClick={() => setShowFilters(!showFilters)}
          className={`flex items-center gap-2 px-3.5 py-2 rounded-xl text-sm font-medium border transition-colors w-full justify-center ${
            activeFilterCount > 0
              ? 'bg-blue-50 border-blue-200 text-blue-700 dark:bg-blue-900/30 dark:border-blue-800 dark:text-blue-400'
              : 'bg-white border-gray-200 text-gray-600 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-300'
          }`}
        >
          <Filter size={16} />
          Filter
          {activeFilterCount > 0 && (
            <span className="bg-blue-600 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
              {activeFilterCount}
            </span>
          )}
          {showFilters ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
        </button>
      </div>

      {/* Filters */}
      <div className={`${showFilters ? 'block' : 'hidden'} sm:block bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-4 transition-colors`}>
        <div className="hidden sm:flex items-center justify-between mb-3">
          <div className="flex items-center gap-2">
            <Filter size={16} className="text-gray-400 dark:text-gray-500" />
            <span className="text-sm font-medium text-gray-600 dark:text-gray-300">Filter</span>
          </div>
          {activeFilterCount > 0 && (
            <button onClick={clearFilters} className="text-xs text-blue-600 dark:text-blue-400 hover:text-blue-700 font-medium">
              Reset
            </button>
          )}
        </div>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-2.5 sm:gap-3">
          <select
            value={filters.month}
            onChange={(e) => { setFilters({ ...filters, month: e.target.value }); setPage(1); }}
            className="border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-700 dark:text-gray-200"
          >
            <option value="">Bulan</option>
            {Array.from({ length: 12 }, (_, i) => i + 1).map((m) => (
              <option key={m} value={m}>{new Date(2000, m - 1).toLocaleString('id', { month: 'long' })}</option>
            ))}
          </select>
          <select
            value={filters.year}
            onChange={(e) => { setFilters({ ...filters, year: e.target.value }); setPage(1); }}
            className="border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-700 dark:text-gray-200"
          >
            <option value="">Tahun</option>
            {[currentYear - 2, currentYear - 1, currentYear, currentYear + 1].map((y) => (
              <option key={y} value={y}>{y}</option>
            ))}
          </select>
          <select
            value={filters.type}
            onChange={(e) => { setFilters({ ...filters, type: e.target.value }); setPage(1); }}
            className="border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-700 dark:text-gray-200"
          >
            <option value="">Semua Tipe</option>
            <option value="income">Pemasukan</option>
            <option value="expense">Pengeluaran</option>
          </select>
          <select
            value={filters.category_id}
            onChange={(e) => { setFilters({ ...filters, category_id: e.target.value }); setPage(1); }}
            className="border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-700 dark:text-gray-200"
          >
            <option value="">Semua Kategori</option>
            {categories.map((c) => (
              <option key={c.id} value={c.id}>{c.name}</option>
            ))}
          </select>
        </div>
        {/* Mobile clear button */}
        {activeFilterCount > 0 && (
          <button onClick={clearFilters} className="sm:hidden w-full mt-2.5 text-xs text-blue-600 dark:text-blue-400 font-medium py-1.5">
            Reset Filter
          </button>
        )}
      </div>

      {/* Transaction count */}
      <p className="text-xs text-gray-400 dark:text-gray-500">{total} transaksi</p>

      {/* Mobile: Card List */}
      <div className="sm:hidden space-y-2.5">
        {loading ? (
          <div className="flex items-center justify-center h-48">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
          </div>
        ) : transactions.length === 0 ? (
          <div className="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700">
            <p className="text-gray-400 dark:text-gray-500 text-sm">Belum ada transaksi</p>
          </div>
        ) : (
          transactions.map((t) => {
            const cat = categoryById[t.category_id] || {};
            return (
              <div
                key={t.id}
                className="bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 p-3.5 active:bg-gray-50 dark:active:bg-gray-700 transition-colors"
                onClick={() => openEdit(t)}
              >
                <div className="flex items-center gap-3">
                  <div
                    className="w-10 h-10 rounded-xl flex items-center justify-center shrink-0"
                    style={{ backgroundColor: (cat.color || '#6b7280') + '15' }}
                  >
                    <div className="w-3 h-3 rounded-full" style={{ backgroundColor: cat.color || '#6b7280' }} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-700 dark:text-gray-200 truncate">
                      {t.description || t.category_name}
                    </p>
                    <div className="flex items-center gap-1.5 mt-0.5">
                      <span className="text-xs text-gray-400 dark:text-gray-500">{t.category_name}</span>
                      <span className="text-gray-300 dark:text-gray-600">Â·</span>
                      <span className="text-xs text-gray-400 dark:text-gray-500">
                        {new Date(t.date + 'T00:00:00').toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })}
                      </span>
                    </div>
                  </div>
                  <div className="text-right shrink-0">
                    <p className={`text-sm font-bold ${t.type === 'income' ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`}>
                      {t.type === 'income' ? '+' : '-'}{formatCurrency(t.amount)}
                    </p>
                    <button
                      onClick={(e) => { e.stopPropagation(); handleDelete(t.id); }}
                      className="mt-0.5 p-1 hover:bg-red-50 dark:hover:bg-red-900/30 rounded-lg transition-colors ml-auto"
                    >
                      <Trash2 size={14} className="text-red-300 dark:text-red-400" />
                    </button>
                  </div>
                </div>
              </div>
            );
          })
        )}
      </div>

      {/* Desktop: Table */}
      <div className="hidden sm:block bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden transition-colors">
        {loading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
          </div>
        ) : transactions.length === 0 ? (
          <p className="text-gray-400 dark:text-gray-500 text-center py-12">Belum ada transaksi</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50">
                  <th className="text-left px-4 py-3 font-medium text-gray-500 dark:text-gray-400">Tanggal</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-500 dark:text-gray-400">Kategori</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-500 dark:text-gray-400">Deskripsi</th>
                  <th className="text-right px-4 py-3 font-medium text-gray-500 dark:text-gray-400">Jumlah</th>
                  <th className="text-center px-4 py-3 font-medium text-gray-500 dark:text-gray-400">Aksi</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map((t) => {
                  const cat = categoryById[t.category_id] || {};
                  return (
                    <tr key={t.id} className="border-b border-gray-50 dark:border-gray-700/50 hover:bg-gray-50 dark:hover:bg-gray-700/30 transition-colors">
                      <td className="px-4 py-3 text-gray-600 dark:text-gray-300 whitespace-nowrap">
                        {new Date(t.date + 'T00:00:00').toLocaleDateString('id-ID')}
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2">
                          <div className="w-2 h-2 rounded-full" style={{ backgroundColor: cat.color || '#6b7280' }} />
                          <span className="text-gray-700 dark:text-gray-200">{t.category_name}</span>
                        </div>
                      </td>
                      <td className="px-4 py-3 text-gray-600 dark:text-gray-300 max-w-[200px] truncate">
                        {t.description || '-'}
                      </td>
                      <td className={`px-4 py-3 text-right font-semibold whitespace-nowrap ${
                        t.type === 'income' ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'
                      }`}>
                        {t.type === 'income' ? '+' : '-'}Rp {formatCurrency(t.amount)}
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center justify-center gap-2">
                          <button onClick={() => openEdit(t)} className="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors">
                            <Pencil size={16} className="text-gray-400 dark:text-gray-500" />
                          </button>
                          <button onClick={() => handleDelete(t.id)} className="p-1.5 hover:bg-red-50 dark:hover:bg-red-900/30 rounded-lg transition-colors">
                            <Trash2 size={16} className="text-red-400" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 px-4 py-3 transition-colors">
          <span className="text-xs sm:text-sm text-gray-500 dark:text-gray-400">
            {page}/{totalPages}
          </span>
          <div className="flex gap-2">
            <button
              disabled={page <= 1}
              onClick={() => setPage(page - 1)}
              className="px-3 py-1.5 text-sm rounded-lg border border-gray-200 dark:border-gray-600 disabled:opacity-50 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors text-gray-700 dark:text-gray-300"
            >
              <span className="sm:hidden">â†</span>
              <span className="hidden sm:inline">Sebelumnya</span>
            </button>
            <button
              disabled={page >= totalPages}
              onClick={() => setPage(page + 1)}
              className="px-3 py-1.5 text-sm rounded-lg border border-gray-200 dark:border-gray-600 disabled:opacity-50 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors text-gray-700 dark:text-gray-300"
            >
              <span className="sm:hidden">â†’</span>
              <span className="hidden sm:inline">Selanjutnya</span>
            </button>
          </div>
        </div>
      )}

      {/* FAB - Mobile only */}
      <button
        onClick={openAdd}
        className="sm:hidden fixed bottom-20 right-4 z-40 w-14 h-14 bg-blue-600 text-white rounded-2xl shadow-lg shadow-blue-600/30 flex items-center justify-center hover:bg-blue-700 active:scale-95 transition-all"
      >
        <Plus size={26} />
      </button>

      {/* Modal */}
      <Modal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        title={editItem ? 'Edit Transaksi' : 'Tambah Transaksi'}
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-2">Tipe</label>
            <div className="flex gap-2.5">
              <button
                type="button"
                onClick={() => setForm({ ...form, type: 'expense', category_id: '' })}
                className={`flex-1 py-2.5 rounded-xl text-sm font-medium border transition-all ${
                  form.type === 'expense'
                    ? 'bg-red-50 border-red-200 text-red-700 shadow-sm dark:bg-red-900/30 dark:border-red-800 dark:text-red-400'
                    : 'border-gray-200 text-gray-500 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-400 dark:hover:bg-gray-700'
                }`}
              >
                Pengeluaran
              </button>
              <button
                type="button"
                onClick={() => setForm({ ...form, type: 'income', category_id: '' })}
                className={`flex-1 py-2.5 rounded-xl text-sm font-medium border transition-all ${
                  form.type === 'income'
                    ? 'bg-green-50 border-green-200 text-green-700 shadow-sm dark:bg-green-900/30 dark:border-green-800 dark:text-green-400'
                    : 'border-gray-200 text-gray-500 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-400 dark:hover:bg-gray-700'
                }`}
              >
                Pemasukan
              </button>
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">Kategori</label>
            <select
              value={form.category_id}
              onChange={(e) => setForm({ ...form, category_id: e.target.value })}
              className="w-full border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-700 dark:text-gray-200"
              required
            >
              <option value="">Pilih kategori</option>
              {categories
                .filter((c) => c.type === form.type)
                .map((c) => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">Jumlah (Rp)</label>
            <input
              type="number"
              value={form.amount}
              onChange={(e) => setForm({ ...form, amount: e.target.value })}
              className="w-full border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-700 dark:text-gray-100 text-lg font-semibold"
              placeholder="100000"
              required
              inputMode="numeric"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">Tanggal</label>
            <input
              type="date"
              value={form.date}
              onChange={(e) => setForm({ ...form, date: e.target.value })}
              className="w-full border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-700 dark:text-gray-200"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">Deskripsi</label>
            <input
              type="text"
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
              className="w-full border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-700 dark:text-gray-200"
              placeholder="Contoh: Makan siang, Belanja bulanan..."
            />
          </div>
          {formError && (
            <div className="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 text-sm rounded-xl px-4 py-2.5">
              {formError}
            </div>
          )}
          <button
            onClick={handleSave}
            disabled={saving}
            className="w-full bg-blue-600 text-white py-3 rounded-xl hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors text-sm font-semibold"
          >
            {saving ? 'Menyimpan...' : editItem ? 'Simpan Perubahan' : 'Tambah Transaksi'}
          </button>
        </div>
      </Modal>
    </div>
  );
}
