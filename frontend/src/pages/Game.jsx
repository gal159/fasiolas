import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import { FiCopy } from 'react-icons/fi';
import GameBoard from '../components/GameBoard';

function Game({ user }) {
  const { id } = useParams();
  const [game, setGame] = useState(null);
  const [players, setPlayers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    const interval = setInterval(fetchGameState, 2000);
    fetchGameState();
    return () => clearInterval(interval);
  }, [id]);

  const fetchGameState = async () => {
    try {
      const response = await axios.get(`/api/v1/games/${id}`);
      setGame(response.data.game);
      setPlayers(response.data.players);
      setError(null);
    } catch (err) {
      console.error('Failed to fetch game state:', err);
      setError('Failed to load game');
    } finally {
      setLoading(false);
    }
  };

  const handleStartGame = async () => {
    try {
      await axios.post(`/api/v1/games/${id}/start`);
      fetchGameState();
    } catch (err) {
      setError('Failed to start game');
    }
  };

  const handleCopyRoomCode = () => {
    navigator.clipboard.writeText(game.room_code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const isPlayerInGame = players.some(
    (p) => p.user?.id === user.id || p.user_id === user.id,
  );

  // Debug logging
  console.log('Current user:', user);
  console.log('Players:', players);
  console.log('Is player in game:', isPlayerInGame);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-red-400">{error}</div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      {/* Game Header */}
      <div className="bg-gray-800 rounded-lg p-6 mb-8 border border-gray-700">
        <div className="flex justify-between items-center">
          <div>
            <h2 className="text-3xl font-bold mb-2">Game #{game.id}</h2>
            <div className="flex items-center space-x-4 text-gray-400">
              <span>Room: <code className="bg-gray-900 px-2 py-1 rounded">{game.room_code}</code></span>
              <button
                onClick={handleCopyRoomCode}
                className="flex items-center space-x-1 hover:text-white transition"
              >
                <FiCopy className="w-4 h-4" />
                <span>{copied ? 'Copied!' : 'Copy'}</span>
              </button>
            </div>
          </div>

          <div className="text-right">
            <p className="text-gray-400 mb-2">Status: <span className="text-blue-400 font-semibold">{game.state}</span></p>
            <p className="text-gray-400">Phase: <span className="text-blue-400 font-semibold">{game.phase}</span></p>
          </div>
        </div>
      </div>

      {/* Game Content */}
      {game.state === 'waiting' ? (
        <div className="bg-gray-800 rounded-lg p-8 border border-gray-700 text-center">
          <h3 className="text-2xl font-semibold mb-6">Waiting for Players</h3>

          <div className="mb-8 grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="bg-gray-900 rounded p-4">
              <p className="text-gray-400 mb-2">Players Joined</p>
              <p className="text-3xl font-bold text-blue-400">{players.length}/{game.max_players}</p>
            </div>
            <div className="bg-gray-900 rounded p-4">
              <p className="text-gray-400 mb-2">Room Code</p>
              <p className="text-2xl font-mono text-green-400">{game.room_code}</p>
            </div>
          </div>

          <div className="mb-8">
            <h4 className="text-lg font-semibold mb-4">Players in Game</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {players.map(player => (
                <div key={player.id} className="bg-gray-900 rounded p-4 text-left">
                  <p className="font-semibold">{player.user?.username || 'Unknown'}</p>
                  <p className="text-gray-400 text-sm">Position: {player.position + 1}</p>
                </div>
              ))}
            </div>
          </div>

          <div className="mt-12 mb-8 flex flex-col items-center">
            <button
              onClick={handleStartGame}
              disabled={players.length < 2}
              className="w-full max-w-md bg-green-600 hover:bg-green-700 disabled:bg-gray-600 disabled:opacity-75 disabled:cursor-not-allowed px-8 py-6 rounded-xl font-bold text-2xl transition-all transform hover:scale-105 shadow-2xl border-2 border-green-500"
            >
              {players.length < 2 ? '⏳ Need 2+ players' : '🎮 START GAME 🎮'}
            </button>
            <p className="text-gray-400 text-sm mt-4">
              {user ? `✓ Logged in as: ${user.username}` : '✗ Not logged in'}
            </p>
            <p className="text-gray-500 text-xs mt-2">
              {players.length}/2 žaidėjai
            </p>
          </div>
        </div>
      ) : (
        <GameBoard game={game} players={players} currentUser={user} onUpdate={fetchGameState} />
      )}
    </div>
  );
}

export default Game;

