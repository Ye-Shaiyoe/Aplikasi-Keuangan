import { useState, useEffect } from 'react';
import { Plus, Pencil, Trash2 } from 'lucide-react';
import Modal from '../components/Modal';
import { getCategories, createCategory, updateCategory, deleteCategory } from '../api/client';

const defaultIcons = ['Wallet', 'Briefcase', 'TrendingUp', 'PlusCircle', 'UtensilsCrossed', 'Car', 'ShoppingBag', 'Gamepad2', 'HeartPulse', 'FileText', 'BookOpen', 'Home', 'MoreHorizontal', 'Coffee', 'Gift', 'Plane', 'Smartphone', 'Tv', 'Dumbbell', 'Circle'];

const defaultColors = ['#ef4444', '#f59e0b', '#10b981', '#3b82f6', '#8b5cf6', '#ec4899', '#f97316', '#6366f1', '#14b8a6', '#84cc16', '#06b6d4', '#a855f7'];

export default function Categories() {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [editItem, setEditItem] = useState(null);
  const [form, setForm] = useState({ name: '', type: 'expense', icon: 'Circle', color: '#6b7280' });

  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = () => {
    setLoading(true);
    getCategories().then(setCategories).finally(() => setLoading(false));
  };

  const openAdd = () => {
    setEditItem(null);
    setForm({ name: '', type: 'expense', icon: 'Circle', color: '#6b7280' });
    setModalOpen(true);
  };

  const openEdit = (c) => {
    setEditItem(c);
    setForm({ name: c.name, type: c.type, icon: c.icon, color: c.color });
    setModalOpen(true);
  };

  const handleSave = async () => {
    if (editItem) {
      await updateCategory(editItem.id, form);
    } else {
      await createCategory(form);
    }
    setModalOpen(false);
    loadCategories();
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Hapus kategori ini?')) return;
    await deleteCategory(id);
    loadCategories();
  };

  const incomeCats = categories.filter((c) => c.type === 'income');
  const expenseCats = categories.filter((c) => c.type === 'expense');

  // Mobile: list style / Desktop: grid style
  const CategorySection = ({ items, title, titleColor }) => (
    <div>
      <h3 className={`text-xs sm:text-sm font-semibold uppercase tracking-wider mb-3 ${titleColor}`}>{title}</h3>
      
      {/* Mobile: List view */}
      <div className="sm:hidden space-y-2">
        {items.map((c) => (
          <div key={c.id} className="bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 p-3 flex items-center gap-3 active:bg-gray-50 dark:active:bg-gray-700 transition-colors">
            <div
              className="w-10 h-10 rounded-xl flex items-center justify-center text-white text-sm font-bold shrink-0 shadow-sm"
              style={{ backgroundColor: c.color }}
            >
              {c.name.charAt(0)}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-gray-700 dark:text-gray-200">{c.name}</p>
              <p className="text-xs text-gray-400 dark:text-gray-500 capitalize">{c.type === 'income' ? 'Pemasukan' : 'Pengeluaran'}</p>
            </div>
            <div className="flex items-center gap-1 shrink-0">
              <button onClick={() => openEdit(c)} className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors">
                <Pencil size={16} className="text-gray-400 dark:text-gray-500" />
              </button>
              <button onClick={() => handleDelete(c.id)} className="p-2 hover:bg-red-50 dark:hover:bg-red-900/30 rounded-lg transition-colors">
                <Trash2 size={16} className="text-red-400" />
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Desktop: Grid view */}
      <div className="hidden sm:grid sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
        {items.map((c) => (
          <div key={c.id} className="bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 p-4 hover:shadow-md transition-shadow">
            <div className="flex flex-col items-center text-center gap-2">
              <div
                className="w-10 h-10 rounded-full flex items-center justify-center text-white text-sm font-bold"
                style={{ backgroundColor: c.color }}
              >
                {c.name.charAt(0)}
              </div>
              <p className="text-sm font-medium text-gray-700 dark:text-gray-200">{c.name}</p>
              <div className="flex gap-1">
                <button onClick={() => openEdit(c)} className="p-1 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors">
                  <Pencil size={14} className="text-gray-400 dark:text-gray-500" />
                </button>
                <button onClick={() => handleDelete(c.id)} className="p-1 hover:bg-red-50 dark:hover:bg-red-900/30 rounded-lg transition-colors">
                  <Trash2 size={14} className="text-red-400" />
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );

  return (
    <div className="space-y-5 sm:space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-xl sm:text-2xl font-bold text-gray-800 dark:text-gray-100">Kategori</h1>
        <button
          onClick={openAdd}
          className="flex items-center gap-2 bg-blue-600 text-white px-3.5 sm:px-4 py-2 sm:py-2.5 rounded-xl hover:bg-blue-700 transition-colors text-sm font-medium"
        >
          <Plus size={18} />
          <span className="sm:hidden">Baru</span>
          <span className="hidden sm:inline">Tambah Kategori</span>
        </button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
        </div>
      ) : (
        <div className="space-y-6 sm:space-y-8">
          <CategorySection items={expenseCats} title="Pengeluaran" titleColor="text-red-500" />
          <CategorySection items={incomeCats} title="Pemasukan" titleColor="text-green-500" />
        </div>
      )}

      <Modal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        title={editItem ? 'Edit Kategori' : 'Tambah Kategori'}
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">Nama</label>
            <input
              type="text"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              className="w-full border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-700 dark:text-gray-100"
              placeholder="Nama kategori"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">Tipe</label>
            <div className="flex gap-2.5">
              <button
                type="button"
                onClick={() => setForm({ ...form, type: 'expense' })}
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
                onClick={() => setForm({ ...form, type: 'income' })}
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
            <label className="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-2">Warna</label>
            <div className="flex flex-wrap gap-2">
              {defaultColors.map((color) => (
                <button
                  key={color}
                  type="button"
                  onClick={() => setForm({ ...form, color })}
                  className={`w-9 h-9 sm:w-8 sm:h-8 rounded-full transition-all ${
                    form.color === color ? 'ring-2 ring-offset-2 ring-blue-500 scale-110' : ''
                  }`}
                  style={{ backgroundColor: color }}
                />
              ))}
            </div>
          </div>
          <button
            onClick={handleSave}
            className="w-full bg-blue-600 text-white py-3 rounded-xl hover:bg-blue-700 transition-colors text-sm font-semibold"
          >
            {editItem ? 'Simpan Perubahan' : 'Tambah Kategori'}
          </button>
        </div>
      </Modal>
    </div>
  );
}
