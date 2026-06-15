import { X } from 'lucide-react';
import { useEffect } from 'react';

export default function Modal({ open, onClose, title, children }) {
  // Prevent body scroll when modal is open
  useEffect(() => {
    if (open) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => { document.body.style.overflow = ''; };
  }, [open]);

  if (!open) return null;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 z-[60] bg-black/40 backdrop-blur-sm transition-opacity"
        onClick={onClose}
      />

      {/* Desktop: centered modal / Mobile: bottom sheet */}
      <div className="fixed inset-0 z-[70] flex items-end sm:items-center justify-center pointer-events-none">
        {/* Desktop modal */}
        <div className="hidden sm:block w-full max-w-md pointer-events-auto">
          <div className="bg-white rounded-2xl shadow-2xl mx-4 max-h-[85vh] overflow-hidden flex flex-col">
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-100 shrink-0">
              <h2 className="text-lg font-semibold text-gray-800">{title}</h2>
              <button onClick={onClose} className="p-1.5 hover:bg-gray-100 rounded-lg transition-colors">
                <X size={20} className="text-gray-400" />
              </button>
            </div>
            <div className="p-6 overflow-y-auto">{children}</div>
          </div>
        </div>

        {/* Mobile bottom sheet */}
        <div className="sm:hidden w-full pointer-events-auto">
          <div className="bg-white rounded-t-3xl shadow-2xl max-h-[92vh] overflow-hidden flex flex-col">
            {/* Drag handle */}
            <div className="flex justify-center pt-3 pb-1 shrink-0">
              <div className="w-10 h-1 bg-gray-300 rounded-full" />
            </div>
            <div className="flex items-center justify-between px-5 py-3 border-b border-gray-50 shrink-0">
              <h2 className="text-base font-semibold text-gray-800">{title}</h2>
              <button onClick={onClose} className="p-1.5 hover:bg-gray-100 rounded-lg transition-colors">
                <X size={20} className="text-gray-400" />
              </button>
            </div>
            <div className="px-5 py-4 overflow-y-auto pb-8">{children}</div>
          </div>
        </div>
      </div>
    </>
  );
}
