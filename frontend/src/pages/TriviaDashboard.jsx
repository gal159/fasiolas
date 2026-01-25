import React, { useEffect, useState } from 'react';
import axios from 'axios';

function TriviaDashboard() {
  const [questions, setQuestions] = useState([]);
  const [amount, setAmount] = useState(5);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [selected, setSelected] = useState({}); // tracks selected answer per question

  useEffect(() => {
    fetchTrivia();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const shuffle = (arr) => {
    const copy = [...arr];
    for (let i = copy.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [copy[i], copy[j]] = [copy[j], copy[i]];
    }
    return copy;
  };

  const fetchTrivia = async () => {
    try {
      setLoading(true);
      setError(null);
      const res = await axios.get(`/api/v1/external/trivia`, {
        params: { amount }
      });
      const normalized = (res.data || []).map((q) => {
        const correct = q.correct_answer ?? q.correctAnswer;
        const incorrect = q.incorrect_answers ?? q.incorrectAnswers ?? [];
        const choices = shuffle([...incorrect, correct]);
        return { ...q, correct_answer: correct, incorrect_answers: incorrect, choices };
      });
      setQuestions(normalized);
      setSelected({});
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch trivia');
    } finally {
      setLoading(false);
    }
  };

  const handleAmountChange = (value) => {
    const n = Number(value);
    if (Number.isNaN(n)) return;
    setAmount(Math.min(Math.max(n, 1), 50));
  };

  const handleSelect = (key, choice) => {
    if (selected[key]) return; // lock after first choice
    setSelected((prev) => ({ ...prev, [key]: choice }));
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 to-gray-800 py-8 px-4">
      <div className="max-w-5xl mx-auto space-y-6">
        <div className="bg-gradient-to-r from-blue-600 to-indigo-700 rounded-lg p-6 shadow-xl flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold text-white">🎯 Trivia Dashboard</h1>
            <p className="text-blue-100 mt-1">Pull fresh trivia from Open Trivia DB for players and spectators.</p>
          </div>
          <div className="flex items-center gap-3">
            <label className="text-white text-sm">Questions:</label>
            <input
              type="number"
              min={1}
              max={50}
              value={amount}
              onChange={(e) => handleAmountChange(e.target.value)}
              className="w-20 bg-gray-900 text-white border border-blue-300 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-400"
            />
            <button
              onClick={fetchTrivia}
              className="bg-white text-blue-700 font-semibold px-4 py-2 rounded-lg hover:bg-blue-50 transition"
              disabled={loading}
            >
              {loading ? 'Loading...' : 'Refresh'}
            </button>
          </div>
        </div>

        {error && (
          <div className="bg-red-900 border border-red-700 text-red-100 rounded-lg p-4">
            ❌ {error}
          </div>
        )}

        {!error && questions.length === 0 && !loading && (
          <div className="bg-gray-800 border border-gray-700 text-gray-200 rounded-lg p-6 text-center">
            No trivia loaded yet. Click Refresh to fetch questions.
          </div>
        )}

        <div className="space-y-4">
          {questions.map((q, idx) => {
            const qKey = `${q.question}-${idx}`;
            const correct = q.correct_answer ?? q.correctAnswer;
            const revealedChoice = selected[qKey];

            return (
              <div key={qKey} className="bg-gray-800 border border-gray-700 rounded-lg p-5 shadow">
                <div className="flex justify-between items-start gap-3">
                  <div>
                    <div className="text-xs text-gray-400 uppercase tracking-wide mb-1">
                      {q.category} • {q.difficulty}
                    </div>
                    <div className="text-lg text-white font-semibold">{idx + 1}. {q.question}</div>
                  </div>
                  <div className="text-xs bg-indigo-700 text-white px-3 py-1 rounded-full">
                    {q.type === 'multiple' ? 'Multiple Choice' : 'Trivia'}
                  </div>
                </div>

                <div className="mt-3 grid grid-cols-1 md:grid-cols-2 gap-2">
                  {(q.choices || []).map((choice, cIdx) => {
                    const isCorrect = choice === correct;
                    const isSelected = revealedChoice === choice;
                    const showCorrect = Boolean(revealedChoice);
                    const base = 'px-3 py-2 rounded border text-sm transition cursor-pointer';
                    let stateClass = 'bg-gray-900 border-gray-700 text-gray-200 hover:border-blue-400';

                    if (showCorrect && isCorrect) {
                      stateClass = 'bg-green-900/40 border-green-700 text-green-100';
                    } else if (isSelected) {
                      stateClass = 'bg-blue-900/40 border-blue-600 text-blue-100';
                    }

                    return (
                      <button
                        type="button"
                        key={`${choice}-${cIdx}`}
                        className={`${base} ${stateClass}`}
                        onClick={() => handleSelect(qKey, choice)}
                        disabled={Boolean(revealedChoice)}
                      >
                        {choice}
                      </button>
                    );
                  })}
                </div>

                <div className="mt-3 text-sm text-gray-300">
                  {revealedChoice ? (
                    <span className={revealedChoice === correct ? 'text-green-300' : 'text-red-300'}>
                      {revealedChoice === correct ? '✅ Correct!' : `❌ Incorrect. Correct answer: ${correct}`}
                    </span>
                  ) : (
                    <span className="text-gray-400">Choose an answer to reveal the correct one.</span>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

export default TriviaDashboard;
