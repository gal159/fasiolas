import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { FiArrowRight, FiCopy, FiEye } from 'react-icons/fi';

function GameCard({ game, onRefresh, userRole = 'player' }) {
  const navigate = useNavigate();
  const [joining, setJoining] = useState(false);
  const [copied, setCopied] = useState(false);

  const handleJoinGame = async () => {
    try {
      setJoining(true);
      await axios.post('/api/v1/games/join', {
        room_code: game.room_code,
      });
      onRefresh();
      navigate(`/game/${game.id}`);
    } catch (err) {
      alert('Failed to join game: ' + (err.response?.data?.error || err.message));
    } finally {
      setJoining(false);
    }
  };

  const handleWatchGame = async () => {
    try {
      setJoining(true);
      // Spectators don't join, they just navigate to view
      navigate(`/game/${game.id}`);
    } catch (err) {
      alert('Failed to watch game: ' + (err.response?.data?.error || err.message));
    } finally {
      setJoining(false);
    }
  };

  const handleCopyCode = () => {
    navigator.clipboard.writeText(game.room_code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const currentPlayers = game.player_count ?? game.players?.length ?? '?';
  const isFull = currentPlayers !== '?' && currentPlayers >= game.max_players;
  const isSpectator = userRole === 'spectator';

  return (
    <div className="bg-gray-800 rounded-lg p-6 border border-gray-700 hover:border-blue-500 transition">
      <div className="flex justify-between items-start mb-4">
        <div>
          <h3 className="text-xl font-semibold">Game #{game.id}</h3>
          <p className="text-gray-400 text-sm">Room: <code>{game.room_code}</code></p>
        </div>
        <span className="bg-blue-600 px-3 py-1 rounded text-sm">{game.state}</span>
      </div>

      <div className="grid grid-cols-2 gap-4 mb-6">
        <div className="bg-gray-900 rounded p-3">
          <p className="text-gray-400 text-xs">Players</p>
          <p className="text-2xl font-bold text-blue-400">{currentPlayers}/{game.max_players}</p>
        </div>
        <div className="bg-gray-900 rounded p-3">
          <p className="text-gray-400 text-xs">Phase</p>
          <p className="text-2xl font-bold text-green-400">{game.phase}</p>
        </div>
      </div>

      <div className="flex space-x-2">
        <button
          onClick={handleCopyCode}
          className="flex-1 flex items-center justify-center space-x-2 bg-gray-700 hover:bg-gray-600 px-3 py-2 rounded text-sm transition"
        >
          <FiCopy className="w-4 h-4" />
          <span>{copied ? 'Copied' : 'Copy Code'}</span>
        </button>

        {isSpectator ? (
          <button
            onClick={handleWatchGame}
            disabled={joining}
            className="flex-1 flex items-center justify-center space-x-2 bg-purple-600 hover:bg-purple-700 disabled:opacity-50 px-3 py-2 rounded text-sm font-medium transition"
          >
            <FiEye className="w-4 h-4" />
            <span>{joining ? 'Loading...' : 'Watch'}</span>
          </button>
        ) : (
          <button
            onClick={handleJoinGame}
            disabled={joining || isFull}
            className="flex-1 flex items-center justify-center space-x-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 px-3 py-2 rounded text-sm font-medium transition"
          >
            <span>{isFull ? 'Full' : joining ? 'Joining...' : 'Join'}</span>
            <FiArrowRight className="w-4 h-4" />
          </button>
        )}
      </div>
    </div>
  );
}

export default GameCard;

