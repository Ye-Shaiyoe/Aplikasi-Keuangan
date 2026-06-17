export default function Card({ children, className = '' }) {
  return (
    <div className={`bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-5 sm:p-6 transition-colors ${className}`}>
      {children}
    </div>
  );
}

export function StatCard({ title, value, icon: Icon, color = 'blue' }) {
  const colorMap = {
    blue: 'bg-blue-50 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400',
    green: 'bg-green-50 text-green-600 dark:bg-green-900/30 dark:text-green-400',
    red: 'bg-red-50 text-red-600 dark:bg-red-900/30 dark:text-red-400',
    purple: 'bg-purple-50 text-purple-600 dark:bg-purple-900/30 dark:text-purple-400',
  };

  // Compact format for mobile, full format for desktop
  const formatValue = (val) => {
    if (val >= 1000000) {
      return `${(val / 1000000).toFixed(val % 1000000 === 0 ? 0 : 1)}jt`;
    }
    if (val >= 1000) {
      return `${(val / 1000).toFixed(val % 1000 === 0 ? 0 : 0)}rb`;
    }
    return new Intl.NumberFormat('id-ID').format(val);
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-4 sm:p-5 transition-colors">
      <div className="flex items-center justify-between">
        <div className="min-w-0 flex-1">
          <p className="text-xs sm:text-sm text-gray-500 dark:text-gray-400">{title}</p>
          <p className="text-lg sm:text-2xl font-bold text-gray-800 dark:text-gray-100 mt-0.5 sm:mt-1 truncate">
            <span className="sm:hidden">Rp {formatValue(value)}</span>
            <span className="hidden sm:inline">{new Intl.NumberFormat('id-ID').format(value)}</span>
          </p>
        </div>
        {Icon && (
          <div className={`p-2.5 sm:p-3 rounded-xl ${colorMap[color] || colorMap.blue} shrink-0`}>
            <Icon size={20} className="sm:hidden" />
            <Icon size={24} className="hidden sm:block" />
          </div>
        )}
      </div>
    </div>
  );
}
