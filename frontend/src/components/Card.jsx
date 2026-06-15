export default function Card({ children, className = '' }) {
  return (
    <div className={`bg-white rounded-2xl shadow-sm border border-gray-100 p-5 sm:p-6 ${className}`}>
      {children}
    </div>
  );
}

export function StatCard({ title, value, icon: Icon, color = 'blue' }) {
  const colorMap = {
    blue: 'bg-blue-50 text-blue-600',
    green: 'bg-green-50 text-green-600',
    red: 'bg-red-50 text-red-600',
    purple: 'bg-purple-50 text-purple-600',
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
    <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-4 sm:p-5">
      <div className="flex items-center justify-between">
        <div className="min-w-0 flex-1">
          <p className="text-xs sm:text-sm text-gray-500">{title}</p>
          {/* Compact on mobile, full on desktop */}
          <p className="text-lg sm:text-2xl font-bold text-gray-800 mt-0.5 sm:mt-1 truncate">
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
