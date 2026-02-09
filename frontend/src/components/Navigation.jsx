import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { FiLogOut, FiUser, FiSettings, FiHelpCircle } from 'react-icons/fi';

function Navigation({ user, onLogout }) {
  const location = useLocation();

  return (
    <nav className="bg-gray-800 border-b border-gray-700 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center space-x-6">
            <h1 className="text-2xl font-bold text-blue-400">🎴 Fasiolas</h1>

            <div className="flex items-center space-x-2">
              <Link
                to="/dashboard"
                className={`px-4 py-2 rounded transition ${
                  location.pathname === '/dashboard'
                    ? 'bg-blue-600 text-white'
                    : 'text-gray-300 hover:bg-gray-700'
                }`}
              >
                Dashboard
              </Link>

              <Link
                to="/trivia"
                className={`flex items-center space-x-2 px-4 py-2 rounded transition ${
                  location.pathname === '/trivia'
                    ? 'bg-indigo-600 text-white'
                    : 'text-gray-300 hover:bg-gray-700'
                }`}
              >
                <FiHelpCircle className="w-4 h-4" />
                <span>Trivia</span>
              </Link>

              {user?.role === 'admin' && (
                <Link
                  to="/admin"
                  className={`flex items-center space-x-2 px-4 py-2 rounded transition ${
                    location.pathname === '/admin'
                      ? 'bg-red-600 text-white'
                      : 'text-gray-300 hover:bg-gray-700'
                  }`}
                >
                  <FiSettings className="w-4 h-4" />
                  <span>Admin Panel</span>
                </Link>
              )}
            </div>
          </div>

          <div className="flex items-center space-x-6">
            <div className="flex items-center space-x-3">
              <FiUser className="w-5 h-5" />
              <span className="text-sm">{user?.username}</span>
              <span className={`px-2 py-1 rounded text-xs font-semibold ${
                user?.role === 'admin' ? 'bg-red-600' :
                user?.role === 'player' ? 'bg-blue-600' :
                'bg-purple-600'
              }`}>
                {user?.role}
              </span>
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
