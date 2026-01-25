import React from 'react';
import { FiLogOut, FiUser } from 'react-icons/fi';

function Navigation({ user, onLogout }) {
  return (
    <nav className="bg-gray-800 border-b border-gray-700 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <h1 className="text-2xl font-bold text-blue-400">🎴 Fasiolas</h1>
          </div>

          <div className="flex items-center space-x-6">
            <div className="flex items-center space-x-3">
              <FiUser className="w-5 h-5" />
              <span className="text-sm">{user?.username}</span>
              <span className="bg-blue-600 px-2 py-1 rounded text-xs">{user?.role}</span>
            </div>

            <button
              onClick={onLogout}
              className="flex items-center space-x-2 bg-red-600 hover:bg-red-700 px-4 py-2 rounded transition"
            >
              <FiLogOut className="w-4 h-4" />
              <span>Logout</span>
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
}

export default Navigation;

