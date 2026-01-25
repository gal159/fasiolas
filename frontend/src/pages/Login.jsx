import React, { useState } from 'react';
import axios from 'axios';
import { FiArrowRight } from 'react-icons/fi';

function Login({ onLogin }) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleOAuthLogin = async (provider) => {
    try {
      setLoading(true);
      setError(null);

      // Clear any cached token so user can pick a different account
      localStorage.removeItem('token');

      // Get OAuth URL from backend
      const response = await axios.get(`/api/v1/auth/${provider}`);

      console.log('OAuth response:', response.data);

      // Check if we got a valid URL
      if (response.data && response.data.url) {
        // Redirect to OAuth provider
        window.location.href = response.data.url;
      } else {
        throw new Error('No OAuth URL received from server');
      }
    } catch (err) {
      console.error('OAuth error:', err);
      setError(`Failed to start ${provider} login: ${err.response?.data?.error || err.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-to-br from-gray-900 to-gray-800">
      <div className="w-full max-w-md">
        <div className="bg-gray-800 rounded-lg shadow-2xl p-8 border border-gray-700">
          <div className="text-center mb-8">
            <h1 className="text-4xl font-bold text-blue-400 mb-2">🎴 Fasiolas</h1>
            <p className="text-gray-400">Card Game REST API</p>
          </div>

          {error && (
            <div className="mb-6 p-4 bg-red-900 border border-red-700 rounded text-red-200">
              {error}
            </div>
          )}

          <div className="space-y-4 mb-8">
            <button
              onClick={() => handleOAuthLogin('google')}
              disabled={loading}
              className="w-full flex items-center justify-center space-x-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 px-4 py-3 rounded font-medium transition"
            >
              <span>Google Login</span>
              <FiArrowRight />
            </button>

            <button
              onClick={() => handleOAuthLogin('github')}
              disabled={loading}
              className="w-full flex items-center justify-center space-x-2 bg-gray-700 hover:bg-gray-600 disabled:opacity-50 px-4 py-3 rounded font-medium transition"
            >
              <span>GitHub Login</span>
              <FiArrowRight />
            </button>

            <button
              onClick={() => handleOAuthLogin('discord')}
              disabled={loading}
              className="w-full flex items-center justify-center space-x-2 bg-purple-600 hover:bg-purple-700 disabled:opacity-50 px-4 py-3 rounded font-medium transition"
            >
              <span>Discord Login</span>
              <FiArrowRight />
            </button>
          </div>

          <div className="text-center text-gray-400 text-sm">
            <p>Choose your OAuth provider to login</p>
            <p className="mt-2">Tip: use Incognito or clear site data if it auto-logs you in.</p>
          </div>
        </div>

        <div className="mt-8 bg-gray-800 rounded-lg p-6 border border-gray-700">
          <h3 className="text-lg font-semibold text-blue-400 mb-4">📖 Getting Started</h3>
          <ol className="text-gray-400 text-sm space-y-2">
            <li>1. Click one of the login buttons above</li>
            <li>2. Complete OAuth authentication</li>
            <li>3. You'll be redirected to the game dashboard</li>
            <li>4. Create or join a game</li>
            <li>5. Start playing Fasiolas!</li>
          </ol>
        </div>
      </div>
    </div>
  );
}

export default Login;

