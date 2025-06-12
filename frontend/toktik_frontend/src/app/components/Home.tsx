'use client';

import { useState } from 'react';
import LoginForm from './Auth/LoginForm';
import Sidebar from './Util/Sidebar';

export default function Home() {
  const [isLoggedIn, setIsLoggedIn] = useState(false); // fake login state
  const [showLoginPopup, setShowLoginPopup] = useState(false);

  const handleVideoClick = (videoNumber: number) => {
    if (!isLoggedIn) {
      setShowLoginPopup(true); // show login popup
    } else {
      alert(`Now playing Video ${videoNumber}`); // show video
    }
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow p-4">
        <h1 className="text-3xl font-bold text-center text-gray-800">TokTik</h1>
      </header>

      <main className="p-6">
        <Sidebar />

        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 pl-[250px]">
          {Array.from({ length: 12 }).map((_, i) => (
            <div
              key={i}
              className="bg-gray-300 h-48 rounded-lg shadow-md flex items-center justify-center text-lg font-semibold text-gray-700 hover:bg-gray-400 cursor-pointer transition"
              onClick={() => handleVideoClick(i + 1)}
            >
              🎬 Video {i + 1}
            </div>
          ))}
        </div>
      </main>



      {showLoginPopup && (
        <LoginForm
          onClose={() => setShowLoginPopup(false)}          // ❌ just closes the popup
          onSuccessLogin={() => {
            setIsLoggedIn(true);                            // ✅ only log in on success
            setShowLoginPopup(false);                       // ✅ close popup after success
          }}
        />

      )}
    </div>
  );
}
