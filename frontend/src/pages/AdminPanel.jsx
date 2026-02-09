import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

function AdminPanel() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [currentUser, setCurrentUser] = useState(null);
  const [deleteConfirm, setDeleteConfirm] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    fetchCurrentUser();
    fetchUsers();
  }, []);

  const fetchCurrentUser = async () => {
    try {
      const response = await axios.get('/api/v1/auth/profile');
      setCurrentUser(response.data);

      // Redirect if not admin
      if (response.data.role !== 'admin') {
        navigate('/dashboard');
      }
    } catch (err) {
      console.error('Failed to fetch current user:', err);
      navigate('/dashboard');
    }
  };

  const fetchUsers = async () => {
    try {
      setLoading(true);
      const response = await axios.get('/api/v1/admin/users', {
        params: { limit: 100, offset: 0 }
      });
      setUsers(response.data || []);
      setError(null);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch users');
    } finally {
      setLoading(false);
    }
  };

  const handleRoleChange = async (userId, newRole) => {
    try {
      await axios.put(`/api/v1/admin/users/${userId}/role`, {
        role: newRole
      });

      // Refresh users list
      fetchUsers();
      alert(`User role updated to ${newRole}`);
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to update role');
    }
  };

  const handleDeleteUser = async (userId) => {
    try {
      await axios.delete(`/api/v1/admin/users/${userId}`);

      // Refresh users list
      fetchUsers();
      setDeleteConfirm(null);
      alert('User deleted successfully');
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to delete user');
      setDeleteConfirm(null);
    }
  };

  const getRoleBadgeColor = (role) => {
    switch (role) {
      case 'admin':
        return 'bg-red-500 text-white';
      case 'player':
        return 'bg-blue-500 text-white';
      case 'spectator':
        return 'bg-purple-500 text-white';
      default:
        return 'bg-gray-500 text-white';
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 to-gray-800 flex items-center justify-center">
        <div className="text-white text-2xl">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 to-gray-800 py-8 px-4">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="bg-gradient-to-r from-red-600 to-red-700 rounded-lg p-6 mb-6 shadow-xl">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-bold text-white mb-2">
                🛡️ Admin Panel
              </h1>
              <p className="text-red-100">
                Manage users and system settings
              </p>
            </div>
            <button
              onClick={() => navigate('/dashboard')}
              className="bg-white text-red-600 px-6 py-2 rounded-lg font-semibold hover:bg-red-50 transition-all"
            >
              ← Back to Dashboard
            </button>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="bg-red-900 border-2 border-red-700 rounded-lg p-4 mb-6 text-red-200">
            ❌ {error}
          </div>
        )}

        {/* Stats Card */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
            <div className="text-gray-400 text-sm mb-1">Total Users</div>
            <div className="text-3xl font-bold text-white">{users.length}</div>
          </div>
          <div className="bg-blue-900 rounded-lg p-6 border border-blue-700">
            <div className="text-blue-200 text-sm mb-1">Players</div>
            <div className="text-3xl font-bold text-white">
              {users.filter(u => u.role === 'player').length}
            </div>
          </div>
          <div className="bg-purple-900 rounded-lg p-6 border border-purple-700">
            <div className="text-purple-200 text-sm mb-1">Spectators</div>
            <div className="text-3xl font-bold text-white">
              {users.filter(u => u.role === 'spectator').length}
            </div>
          </div>
          <div className="bg-red-900 rounded-lg p-6 border border-red-700">
            <div className="text-red-200 text-sm mb-1">Admins</div>
            <div className="text-3xl font-bold text-white">
              {users.filter(u => u.role === 'admin').length}
            </div>
          </div>
        </div>

        {/* Users Table */}
        <div className="bg-gray-800 rounded-lg shadow-xl overflow-hidden border border-gray-700">
          <div className="p-6 border-b border-gray-700">
            <h2 className="text-2xl font-bold text-white">User Management</h2>
            <p className="text-gray-400 mt-1">
              {users.length} users registered
            </p>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-900">
                <tr>
                  <th className="px-6 py-4 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                    ID
                  </th>
                  <th className="px-6 py-4 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                    Username
                  </th>
                  <th className="px-6 py-4 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                    Email
                  </th>
                  <th className="px-6 py-4 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                    Role
                  </th>
                  <th className="px-6 py-4 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                    Joined
                  </th>
                  <th className="px-6 py-4 text-right text-xs font-medium text-gray-400 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-700">
                {users.map((user) => (
                  <tr
                    key={user.id}
                    className={`hover:bg-gray-750 transition-colors ${
                      currentUser?.id === user.id ? 'bg-gray-750' : ''
                    }`}
                  >
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">
                      {user.id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        {user.avatar_url && (
                          <img
                            src={user.avatar_url}
                            alt={user.username}
                            className="h-8 w-8 rounded-full mr-3"
                          />
                        )}
                        <div className="text-sm font-medium text-white">
                          {user.username}
                          {currentUser?.id === user.id && (
                            <span className="ml-2 text-xs text-yellow-400">(You)</span>
                          )}
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">
                      {user.email}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <select
                        value={user.role}
                        onChange={(e) => handleRoleChange(user.id, e.target.value)}
                        disabled={currentUser?.id === user.id}
                        className={`px-3 py-1 rounded-full text-xs font-semibold ${getRoleBadgeColor(
                          user.role
                        )} ${
                          currentUser?.id === user.id
                            ? 'opacity-50 cursor-not-allowed'
                            : 'cursor-pointer'
                        }`}
                      >
                        <option value="player">Player</option>
                        <option value="spectator">Spectator</option>
                        <option value="admin">Admin</option>
                      </select>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-400">
                      {new Date(user.created_at).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      {currentUser?.id !== user.id ? (
                        deleteConfirm === user.id ? (
                          <div className="flex justify-end gap-2">
                            <button
                              onClick={() => handleDeleteUser(user.id)}
                              className="bg-red-600 hover:bg-red-700 text-white px-3 py-1 rounded text-xs font-semibold transition-all"
                            >
                              ✓ Confirm
                            </button>
                            <button
                              onClick={() => setDeleteConfirm(null)}
                              className="bg-gray-600 hover:bg-gray-700 text-white px-3 py-1 rounded text-xs font-semibold transition-all"
                            >
                              ✗ Cancel
                            </button>
                          </div>
                        ) : (
                          <button
                            onClick={() => setDeleteConfirm(user.id)}
                            className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg font-semibold transition-all transform hover:scale-105"
                          >
                            🗑️ Delete
                          </button>
                        )
                      ) : (
                        <span className="text-gray-500 text-xs">Cannot delete yourself</span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {users.length === 0 && (
            <div className="text-center py-12 text-gray-400">
              No users found
            </div>
          )}
        </div>

        {/* Warning */}
        <div className="mt-6 bg-yellow-900 border-2 border-yellow-700 rounded-lg p-4">
          <div className="flex items-start">
            <span className="text-2xl mr-3">⚠️</span>
            <div>
              <h3 className="text-yellow-300 font-bold mb-1">Warning</h3>
              <p className="text-yellow-200 text-sm">
                Deleting a user is permanent and will remove all their game data. This action cannot be undone.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default AdminPanel;
