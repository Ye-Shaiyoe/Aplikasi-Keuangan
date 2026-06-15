import { NavLink, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import {
  LayoutDashboard, ArrowLeftRight, Wallet, BarChart3, LogOut, UserCircle, Target,
} from 'lucide-react';

const navItems = [
  { to: '/', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/transactions', label: 'Transaksi', icon: ArrowLeftRight },
  { to: '/budgets', label: 'Anggaran', icon: Wallet },
  { to: '/savings', label: 'Tabungan', icon: Target },
  { to: '/reports', label: 'Laporan', icon: BarChart3 },
];

// Items shown in bottom nav (max 5 for good UX)
const bottomNavItems = navItems;

export default function Layout({ children }) {
  const { user, logout } = useAuth();
  const location = useLocation();

  const pageTitles = {
    '/': 'Dashboard',
    '/transactions': 'Transaksi',
    '/categories': 'Kategori',
    '/savings': 'Tabungan',
    '/reports': 'Laporan',
  };

  const currentTitle = pageTitles[location.pathname] || 'Catatan Keuangan';

  const SidebarContent = ({ onItemClick }) => (
    <>
      <div className="p-6 border-b border-gray-100">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center text-white font-bold text-lg shadow-sm">
            C
          </div>
          <div>
            <h1 className="text-lg font-bold text-gray-800">Catatan Keuangan</h1>
            <p className="text-xs text-gray-400">Personal Finance</p>
          </div>
        </div>
      </div>
      <nav className="flex-1 p-4 space-y-1">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.to === '/'}
            onClick={onItemClick}
            className={({ isActive }) =>
              `flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-all ${
                isActive
                  ? 'bg-blue-50 text-blue-700 shadow-sm'
                  : 'text-gray-500 hover:bg-gray-50 hover:text-gray-700'
              }`
            }
          >
            <item.icon size={20} />
            {item.label}
          </NavLink>
        ))}
      </nav>
      {/* User section at bottom */}
      <div className="p-4 border-t border-gray-100">
        <div className="flex items-center gap-3 px-3 py-3 bg-gray-50 rounded-xl">
          <div className="w-9 h-9 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white text-sm font-bold shadow-sm">
            {user?.name?.charAt(0)?.toUpperCase() || 'U'}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-gray-700 truncate">{user?.name}</p>
            <p className="text-xs text-gray-400 truncate">{user?.email}</p>
          </div>
          <button
            onClick={logout}
            className="p-2 hover:bg-red-50 rounded-lg transition-colors group"
            title="Keluar"
          >
            <LogOut size={18} className="text-gray-400 group-hover:text-red-500" />
          </button>
        </div>
      </div>
    </>
  );

  return (
    <div className="min-h-screen bg-gray-50 flex">
      {/* Sidebar - Desktop */}
      <aside className="hidden md:flex md:flex-col w-64 bg-white border-r border-gray-100 fixed top-0 bottom-0 left-0 z-30">
        <SidebarContent onItemClick={() => {}} />
      </aside>

      {/* Mobile Header */}
      <header className="md:hidden fixed top-0 left-0 right-0 z-50 bg-white/80 backdrop-blur-xl border-b border-gray-100">
        <div className="flex items-center justify-between px-4 py-3">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center text-white font-bold text-sm shadow-sm">
              C
            </div>
            <h1 className="text-base font-bold text-gray-800">{currentTitle}</h1>
          </div>
          <div className="flex items-center gap-1">
            <button
              onClick={logout}
              className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
              title="Keluar"
            >
              <LogOut size={18} className="text-gray-400" />
            </button>
          </div>
        </div>
      </header>

      {/* Bottom Navigation - Mobile */}
      <nav className="md:hidden fixed bottom-0 left-0 right-0 z-50 bg-white border-t border-gray-100 safe-area-pb">
        <div className="grid grid-cols-5 h-16">
          {bottomNavItems.map((item) => {
            const isActive =
              item.to === '/'
                ? location.pathname === '/'
                : location.pathname.startsWith(item.to);
            return (
              <NavLink
                key={item.to}
                to={item.to}
                end={item.to === '/'}
                className={`flex flex-col items-center justify-center gap-0.5 transition-colors ${
                  isActive ? 'text-blue-600' : 'text-gray-400'
                }`}
              >
                <div className={`p-1 rounded-lg transition-all ${isActive ? 'scale-110' : ''}`}>
                  <item.icon size={22} strokeWidth={isActive ? 2.5 : 1.8} />
                </div>
                <span className={`text-[10px] ${isActive ? 'font-semibold' : 'font-medium'}`}>
                  {item.label}
                </span>
                {isActive && (
                  <div className="absolute top-0 w-8 h-0.5 bg-blue-600 rounded-full" />
                )}
              </NavLink>
            );
          })}
        </div>
      </nav>

      {/* Main Content */}
      <main className="flex-1 md:ml-64">
        <div className="pt-14 pb-20 md:pt-0 md:pb-0">
          <div className="max-w-7xl mx-auto px-3 py-4 md:px-8 md:py-8">
            {children}
          </div>
        </div>
      </main>
    </div>
  );
}
