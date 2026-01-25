import React, { useState, useRef, useEffect } from 'react';
import axios from 'axios';

function GameBoard({ game, players, currentUser, onUpdate }) {
  const [selectedTarget, setSelectedTarget] = useState(null);
  const [error, setError] = useState(null);
  const [drawnCard, setDrawnCard] = useState(null);
  const [canPlaceDrawnCard, setCanPlaceDrawnCard] = useState(false);
  const [waitingForPlacement, setWaitingForPlacement] = useState(false);

  // Use ref to persist state across renders caused by polling
  const isPlacingCardRef = useRef(false);
  const drawnCardRef = useRef(null);

  const currentPlayer = players.find(p => p.user?.id === currentUser.id);
  const isYourTurn = game.current_player_position === currentPlayer?.position;

  // Debug logging
  console.log('GameBoard render - drawnCard:', drawnCard, 'waitingForPlacement:', waitingForPlacement, 'ref:', isPlacingCardRef.current);
  console.log('Current player TopCard:', currentPlayer?.top_card);

  // Helper function to convert suit name to symbol
  const getSuitSymbol = (suit) => {
    const suitMap = {
      'hearts': '♥',
      'diamonds': '♦',
      'clubs': '♣',
      'spades': '♠'
    };
    return suitMap[suit?.toLowerCase()] || suit;
  };

  // Helper function to get suit color
  const getSuitColor = (suit) => {
    const suitLower = suit?.toLowerCase();
    return (suitLower === 'hearts' || suitLower === 'diamonds') ? 'text-red-500' : 'text-gray-900';
  };
  // Position players around the table (you at bottom, others at top)
  const getPlayerPosition = (playerIndex, totalPlayers) => {
    if (!currentPlayer) return null;

    const currentIndex = players.findIndex(p => p.user?.id === currentUser.id);
    const relativeIndex = (playerIndex - currentIndex + totalPlayers) % totalPlayers;

    // Simple layout: you at bottom, others at top
    if (totalPlayers === 2) {
      return relativeIndex === 0 ? 'bottom' : 'top';
    } else {
      // For 3+ players, put current player at bottom, distribute others at top
      return relativeIndex === 0 ? 'bottom' : 'top';
    }
  };

  const handlePlaceCard = async (targetPosition) => {
    try {
      if (!currentPlayer) return;

      // Ensure targetPosition is a valid number
      const position = Number(targetPosition);
      if (isNaN(position)) {
        setError('Invalid target position');
        return;
      }

      console.log('Placing card on position:', position);

      await axios.post(`/api/v1/games/${game.id}/place`, {
        target_player_position: position,
      });

      // Clear all placement-related state and refs
      setSelectedTarget(null);
      setDrawnCard(null);
      setWaitingForPlacement(false);
      isPlacingCardRef.current = false;
      drawnCardRef.current = null;
      setError(null);

      console.log('Card placed, clearing refs');

      onUpdate();
    } catch (err) {
      console.error('Place card error:', err.response?.data);
      setError(err.response?.data?.error || 'Failed to place card');
    }
  };

  const handleDragStart = (e) => {
    if (!drawnCard) return;
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('cardType', 'drawn');
  };

  const handleDragOver = (e) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  };

  const handleDrawCard = async () => {
    try {
      setWaitingForPlacement(true);
      isPlacingCardRef.current = true;

      const res = await axios.post(`/api/v1/games/${game.id}/draw`);
      const drawn = res.data?.card;

      setDrawnCard(drawn || null);
      drawnCardRef.current = drawn || null;
      setError(null);

      console.log('Card drawn:', drawn, 'Setting refs to true');
      // Don't call onUpdate here - it would show the card in the player's pile
      // Wait until the card is placed, then update
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to draw card');
      setWaitingForPlacement(false);
      isPlacingCardRef.current = false;
    }
  };

  // Render a playing card component
  const PlayingCard = ({ card, isTopCard = false, size = 'normal', clickable = false, onClick = null }) => {
    if (!card) return null;

    const sizeClasses = {
      small: 'w-12 h-16 text-sm',
      normal: 'w-20 h-28 text-lg',
      large: 'w-24 h-32 text-2xl',
      deck: 'w-20 h-28 text-lg'
    };

    return (
      <div
        onClick={onClick}
        className={`${sizeClasses[size]} bg-white rounded-lg shadow-lg flex flex-col items-center justify-center border-2 ${isTopCard ? 'border-yellow-400 ring-2 ring-yellow-400' : 'border-gray-300'} relative ${clickable ? 'cursor-pointer hover:shadow-xl hover:scale-105 transition-all' : ''}`}
      >
        <div className={`font-bold ${getSuitColor(card.suit)}`}>
          <div className="text-xs absolute top-1 left-1">{card.rank}</div>
          <div className="text-3xl">{getSuitSymbol(card.suit)}</div>
          <div className="text-xs absolute bottom-1 right-1 rotate-180">{card.rank}</div>
        </div>
      </div>
    );
  };

  // Deck card (face down)
  const DeckCard = ({ count, onClick, disabled }) => (
    <button
      onClick={onClick}
      disabled={disabled}
      className={`relative w-20 h-28 rounded-lg border-2 border-yellow-500 shadow-2xl flex flex-col items-center justify-center transition-all ${
        disabled ? 'opacity-60 cursor-not-allowed bg-blue-700' : 'bg-blue-800 hover:scale-105 hover:-translate-y-2 cursor-pointer'
      }`}
    >
      <div className="text-2xl">🂠</div>
      <div className="text-xs text-white mt-1">Deck</div>
      <div className="text-xs text-yellow-300 font-bold">{count}</div>
      {!disabled && <div className="absolute -top-3 -right-2 bg-yellow-400 text-black text-xs px-2 py-0.5 rounded-full font-bold animate-bounce">TAP</div>}
    </button>
  );

  // Player area with their pile of cards
  const PlayerArea = ({ player, position }) => {
    const isCurrentTurn = game.current_player_position === player.position;
    const isYou = player.user?.id === currentUser.id;
    const isClickable = isYourTurn && !isYou && waitingForPlacement && drawnCard;

    const positionClasses = {
      bottom: 'absolute bottom-0 left-1/2 transform -translate-x-1/2',
      top: 'absolute top-0 left-1/2 transform -translate-x-1/2'
    };

    return (
      <div
        className={`flex flex-col items-center p-4 transition-all`}
        onDragOver={handleDragOver}
        onDrop={(e) => {
          e.preventDefault();
          if (e.dataTransfer.getData('cardType') === 'drawn' && drawnCard && waitingForPlacement) {
            handlePlaceCard(player.position);
          }
        }}
      >
        <div
          className={`bg-gray-800 rounded-lg p-4 border-4 text-center min-w-[200px] transition-all cursor-pointer hover:scale-105 ${
            isCurrentTurn ? 'border-green-400 shadow-lg shadow-green-400/50' :
            isYou ? 'border-blue-400 shadow-lg shadow-blue-400/50' :
            drawnCard && selectedTarget === player.position ? 'border-yellow-400 shadow-lg shadow-yellow-400/50 ring-4 ring-yellow-300' :
            drawnCard && waitingForPlacement ? 'border-gray-600 hover:border-yellow-300 hover:shadow-lg' :
            'border-gray-600'
          }`}
          onClick={() => {
            if (drawnCard && waitingForPlacement) {
              handlePlaceCard(player.position);
            }
          }}
        >
          {/* Player Name */}
          <p className={`font-bold text-sm ${isYou ? 'text-blue-400' : 'text-white'}`}>
            {player.user?.username || 'Unknown'}
            {isYou && ' (You)'}
          </p>

          {isCurrentTurn && <p className="text-xs text-green-400 font-bold">🎯 Playing</p>}

          {/* Top Card Display */}
          <div className="mt-2">
            {player.top_card ? (
              <div className="flex justify-center">
                <PlayingCard card={player.top_card} isTopCard={true} size="normal" />
              </div>
            ) : (
              <div className="w-20 h-28 bg-gray-700 rounded-lg flex items-center justify-center text-gray-500 text-xs border-2 border-gray-600">
                No Cards
              </div>
            )}
          </div>

          {/* Card Count */}
          <p className="text-xs text-gray-400 mt-2">
            📚 {player.card_count} cards
          </p>
        </div>
      </div>
    );
  };

  return (
    <div className="space-y-6">
      {/* Game Info Bar */}
      <div className="bg-gradient-to-r from-gray-800 to-gray-900 rounded-lg p-4 border border-gray-700 shadow-xl">
        <div className="flex justify-between items-center flex-wrap gap-4">
          <div>
            <p className="text-xs text-gray-400">Phase</p>
            <p className="text-xl font-bold text-blue-400">{game.phase}</p>
          </div>
          <div>
            <p className="text-xs text-gray-400">Current Turn</p>
            <p className="text-xl font-bold text-green-400">
              {players.find(p => p.position === game.current_player_position)?.user?.username || 'Unknown'}
            </p>
          </div>
          <div>
            <p className="text-xs text-gray-400">Deck</p>
            <p className="text-xl font-bold text-yellow-400">{game.deck_count} cards</p>
          </div>
          {isYourTurn && (
            <div className="bg-green-900 px-4 py-2 rounded-lg border-2 border-green-400 animate-pulse">
              <p className="text-green-400 font-bold">🎯 YOUR TURN!</p>
            </div>
          )}
        </div>
      </div>

      {/* Error Message */}
      {error && (
        <div className="p-4 bg-red-900 border-2 border-red-700 rounded-lg text-red-200 text-center font-semibold">
          ❌ {error}
        </div>
      )}

      {/* Main Game Table */}
      <div
        className="relative mx-auto bg-gradient-to-br from-green-800 via-green-700 to-green-900 rounded-full shadow-2xl border-8 border-yellow-900 flex items-center justify-center"
        style={{ minHeight: '600px', minWidth: '600px', maxWidth: '100%', aspectRatio: '1' }}
      >
        {/* Table texture overlay */}
        <div className="absolute inset-0 rounded-full opacity-20" style={{
          backgroundImage: 'radial-gradient(circle, rgba(0,0,0,0.1) 1px, transparent 1px)',
          backgroundSize: '20px 20px'
        }}></div>

        {/* Deck in center (clickable) */}
        <div className="relative z-20 flex flex-col items-center justify-center gap-4">
          <DeckCard
            count={game.deck_count}
            onClick={handleDrawCard}
            disabled={!isYourTurn || drawnCard !== null}
          />

          {/* Drawn Card Display */}
          {drawnCard && (
            <div
              className="flex flex-col items-center animate-pulse cursor-move"
              draggable={true}
              onDragStart={handleDragStart}
            >
              <div className="bg-yellow-900/40 rounded-lg p-4 border-4 border-yellow-400 backdrop-blur-sm shadow-2xl">
                <PlayingCard card={drawnCard} size="large" isTopCard={true} />
              </div>
              <p className="text-yellow-300 text-xs sm:text-sm mt-2 font-bold animate-bounce text-center px-2">
                ⬇️ Tap a player to place card ⬇️
              </p>
            </div>
          )}
        </div>

        {/* Players positioned around the circle */}
        <div className="absolute inset-0 pointer-events-none">
          {/* Top players */}
          <div className="absolute top-4 sm:top-8 left-1/2 transform -translate-x-1/2 flex gap-2 sm:gap-4 flex-wrap justify-center pointer-events-auto max-w-[90%]">
            {players
              .map((player, idx) => ({ player, idx }))
              .filter(({ idx }) => {
                const currentIndex = players.findIndex(p => p.user?.id === currentUser.id);
                const relativeIndex = (idx - currentIndex + players.length) % players.length;
                return relativeIndex !== 0;
              })
              .map(({ player, idx }) => (
                <div key={player.id} className="flex flex-col items-center">
                  <PlayerArea player={player} position="top" />
                </div>
              ))}
          </div>

          {/* Bottom player (current user) */}
          <div className="absolute bottom-4 sm:bottom-8 left-1/2 transform -translate-x-1/2 pointer-events-auto">
            {players.find(p => p.user?.id === currentUser.id) && (
              <PlayerArea
                player={players.find(p => p.user?.id === currentUser.id)}
                position="bottom"
              />
            )}
          </div>
        </div>
      </div>

      {/* Action Panel */}
      {currentPlayer && isYourTurn && (
        <div className="bg-gradient-to-r from-blue-900 to-indigo-900 rounded-xl p-6 border-2 border-blue-500 shadow-xl">
          <h3 className="text-xl font-bold mb-4 text-center text-blue-300">🎮 Your Turn</h3>

          {!drawnCard ? (
            <div className="text-center">
              <p className="text-blue-200 mb-4">
                {game.deck_count > 0
                  ? 'Click the deck in the center to draw a card'
                  : 'Deck is empty!'}
              </p>
              <button
                onClick={handleDrawCard}
                disabled={game.deck_count === 0}
                className="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:opacity-50 px-6 py-3 rounded-lg font-bold text-lg transition-all transform hover:scale-105 shadow-lg"
              >
                🎴 Draw Card
              </button>
            </div>
          ) : (
            <div className="text-center">
              <p className="text-yellow-300 font-bold mb-3">
                📌 Card drawn: {drawnCard.rank} {getSuitSymbol(drawnCard.suit)}
              </p>
              <p className="text-blue-200 mb-4">
                Drag the card above to a player or click a player below to place it
              </p>
              {selectedTarget !== null && (
                <div className="bg-green-900 rounded-lg p-3 border-2 border-green-500 mb-4">
                  <p className="text-green-300 font-bold">
                    ✓ Target: {players.find(p => p.position === selectedTarget)?.user?.username}
                  </p>
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {/* Your Hand Info (if you're in the game but not your turn) */}
      {currentPlayer && !isYourTurn && (
        <div className="bg-gray-800 rounded-xl p-4 border-2 border-gray-600">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-400 text-sm">Your Hand</p>
              <p className="text-2xl font-bold text-blue-400">{currentPlayer.card_count} cards</p>
            </div>
            <div>
              <p className="text-gray-400 text-sm mb-2">Your Top Card:</p>
              {currentPlayer.top_card ? (
                <PlayingCard card={currentPlayer.top_card} isTopCard={true} size="normal" />
              ) : (
                <p className="text-gray-500">No cards</p>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default GameBoard;

