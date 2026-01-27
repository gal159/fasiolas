import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { FiPlus, FiRefreshCw } from 'react-icons/fi';
import { useNavigate } from 'react-router-dom';
import GameCard from '../components/GameCard';

function Dashboard({ user }) {
  const [games, setGames] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [maxPlayers, setMaxPlayers] = useState(4);
  const [joinCode, setJoinCode] = useState('');
  const [joinLoading, setJoinLoading] = useState(false);
  const [joinError, setJoinError] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    fetchGames();
  }, []);

  const fetchGames = async () => {
    try {
      setLoading(true);
      setError(null);
      console.log('📡 Fetching games from /api/v1/games?state=waiting&limit=20');

      const response = await axios.get('/api/v1/games?state=waiting&limit=20');
      console.log('✅ Games response received:', response.status, response.data);

      // Ensure we always set an array
      if (Array.isArray(response.data)) {
        setGames(response.data);
      } else if (response.data === null || response.data === undefined) {
        console.warn('⚠️ Response data was null/undefined, using empty array');
        setGames([]);
      } else {
        console.error('❌ Unexpected response data type:', typeof response.data, response.data);
        setGames([]);
      }
    } catch (err) {
      console.error('❌ Failed to fetch games:', {
        status: err.response?.status,
        statusText: err.response?.statusText,
        errorData: err.response?.data,
        errorMessage: err.message,
        isNetworkError: !err.response
      });

      const errorMessage = err.response?.data?.error || err.message || 'Failed to load games';
      setError('Failed to load games: ' + errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateGame = async (e) => {
    e.preventDefault();
    try {
      console.log('Creating game with maxPlayers:', maxPlayers);
      const response = await axios.post('/api/v1/games', {
        max_players: parseInt(maxPlayers),
      });

      console.log('Game creation response:', response.data);
      console.log('Game ID:', response.data.id);

      if (!response.data || !response.data.id) {
        console.error('ERROR: No game ID in response!', response.data);
        setError('Game created but no ID returned - check console');
        return;
      }

      const gameId = response.data.id;
      console.log('✅ Game created successfully! ID:', gameId);

      setShowCreateForm(false);
      setMaxPlayers(4);
      setError(null);

      // Navigate to game room immediately after creating
      console.log('🎯 Navigating to /game/' + gameId);

      // Use setTimeout to ensure state updates first
      setTimeout(() => {
        navigate(`/game/${gameId}`, { replace: true });
      }, 100);

    } catch (err) {
      console.error('Create game error:', err);
      setError('Failed to create game: ' + (err.response?.data?.error || err.message));
    }
  };

  const handleJoinByCode = async (e) => {
    e.preventDefault();
    const code = joinCode.trim().toUpperCase();
    if (!code) {
      setJoinError('Enter a room code');
      return;
    }
    try {
      setJoinLoading(true);
      setJoinError(null);
      const response = await axios.post('/api/v1/games/join', { room_code: code });
      setJoinCode('');
      await fetchGames();
      navigate(`/game/${response.data.id}`);
    } catch (err) {
      setJoinError(err.response?.data?.error || 'Failed to join game');
    } finally {
      setJoinLoading(false);
    }
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      {/* Header */}
      <div className="flex justify-between items-center mb-12">
        <div>
          <h2 className="text-3xl font-bold mb-2">
            {user?.role === 'spectator' ? '👁️ Watch Games' : 'Games'}
          </h2>
          <p className="text-gray-400">
            {user?.role === 'spectator'
              ? 'Watch other players in action'
              : `Welcome back, ${user?.username}!`}
          </p>
        </div>

        <div className="flex space-x-4">
          <button
            onClick={fetchGames}
            className="flex items-center space-x-2 bg-gray-700 hover:bg-gray-600 px-4 py-2 rounded transition"
          >
            <FiRefreshCw className="w-4 h-4" />
            <span>Refresh</span>
          </button>

          {(user?.role === 'player' || user?.role === 'admin') && (
            <button
              onClick={() => setShowCreateForm(!showCreateForm)}
              className="flex items-center space-x-2 bg-blue-600 hover:bg-blue-700 px-4 py-2 rounded transition"
            >
              <FiPlus className="w-4 h-4" />
              <span>Create Game</span>
            </button>
          )}
        </div>
      </div>

      {/* Join by Room Code - Only show for players and admins */}
      {(user?.role === 'player' || user?.role === 'admin') && (
        <div className="bg-gray-800 rounded-lg p-6 mb-8 border border-gray-700">
          <h3 className="text-xl font-semibold mb-3">Join by Room Code</h3>
          <form onSubmit={handleJoinByCode} className="space-y-3 md:flex md:space-y-0 md:space-x-4">
            <input
              type="text"
              placeholder="Enter code (e.g. ABC123)"
              value={joinCode}
              onChange={(e) => setJoinCode(e.target.value)}
              className="flex-1 bg-gray-700 border border-gray-600 rounded px-3 py-2 text-white"
              maxLength={10}
            />
            <button
              type="submit"
              disabled={joinLoading}
              className="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 px-6 py-2 rounded font-medium transition"
            >
              {joinLoading ? 'Joining...' : 'Join Game'}
            </button>
          </form>
          {joinError && (
            <p className="text-red-300 text-sm mt-2">{joinError}</p>
          )}
        </div>
      )}

      {/* Create Game Form */}
      {showCreateForm && (
        <div className="bg-gray-800 rounded-lg p-6 mb-8 border border-gray-700">
          <h3 className="text-xl font-semibold mb-4">Create New Game</h3>
          <form onSubmit={handleCreateGame} className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-2">Max Players (2-8)</label>
              <select
                value={maxPlayers}
                onChange={(e) => setMaxPlayers(e.target.value)}
                className="w-full bg-gray-700 border border-gray-600 rounded px-3 py-2 text-white"
              >
                {[2, 3, 4, 5, 6, 7, 8].map(n => (
                  <option key={n} value={n}>{n} Players</option>
                ))}
              </select>
            </div>

            <div className="flex space-x-4">
              <button
                type="submit"
                className="flex-1 bg-blue-600 hover:bg-blue-700 px-4 py-2 rounded font-medium transition"
              >
                Create
              </button>
              <button
                type="button"
                onClick={() => setShowCreateForm(false)}
                className="flex-1 bg-gray-700 hover:bg-gray-600 px-4 py-2 rounded font-medium transition"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Error Message */}
      {error && (
        <div className="mb-8 p-4 bg-red-900 border border-red-700 rounded text-red-200">
          {error}
        </div>
      )}

      {/* Games Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {loading ? (
          <div className="col-span-full flex justify-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white"></div>
          </div>
        ) : games.length > 0 ? (
          games.map(game => (
            <GameCard key={game.id} game={game} onRefresh={fetchGames} userRole={user?.role} />
          ))
        ) : (
          <div className="col-span-full text-center py-12">
            <p className="text-gray-400 mb-4">No games available</p>
            <button
              onClick={() => setShowCreateForm(true)}
              className="bg-blue-600 hover:bg-blue-700 px-6 py-2 rounded"
            >
              Create the first game!
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

export default Dashboard;

