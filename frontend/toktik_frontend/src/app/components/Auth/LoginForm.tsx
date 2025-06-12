// src/components/Auth/LoginForm.tsx
'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation'; 

export default function LoginForm({
  onClose,
  onSuccessLogin,
}: {
  onClose: () => void;
  onSuccessLogin: () => void;
}) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const router = useRouter();

  const handleLogin = async () => {
    if (!username || !password) {
      alert('Please enter both username and password.');
      return;
    }

    // Simulate login
    if (username === 'admin' && password === 'a') {
      alert('Login successful!');
      onSuccessLogin();  // ✅ ONLY here!
    } else {
      alert('Login failed. Try user@example.com / 1234');
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-8 rounded-xl shadow-xl w-full max-w-md">
        <h2 className="text-2xl font-bold mb-6 text-center text-black">Log In</h2>
        <input
          className="w-full p-2 mb-3 border rounded"
          placeholder="username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input
          className="w-full p-2 mb-4 border rounded"
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button
          className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600"
          onClick={handleLogin}
        >
          Log In
        </button>

        
        <button
          className="w-full bg-green-500 text-white py-2 rounded hover:bg-green-600 mt-3"
          onClick={() => router.push('/signup')}
        >
          Sign Up
        </button>


        <button
          className="mt-4 w-full text-sm text-gray-500 hover:underline"
          onClick={onClose}     // ✅ Just closes — doesn’t log in
        >
          Cancel
        </button>
      </div>
    </div>
  );
}
