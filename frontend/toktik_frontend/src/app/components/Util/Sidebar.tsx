'use client';

import React, { useState } from 'react';
import Button from './Button'; // Adjust the path if needed
import LoginForm from '../Auth/LoginForm';
import UploadForm from './UploadForm'

function Sidebar() {
  const [showLogin, setShowLogin] = useState(false);
  const [showUpload, setShowUpload] = useState(false);

  return (
    <>
      <aside className="fixed top-0 left-0 h-screen w-64 bg-black text-white p-4 flex flex-col gap-4">
        <h2 className="text-xl font-bold">Sidebar</h2>


        <button
          onClick={() => setShowUpload(true)}
          className="flex items-center gap-2 text-white py-2 px-4 rounded text-center bg-transparent hover:bg-white/10 transition"
        >
          <span className="text-2xl">⬆️</span>
          <span>Upload</span>
        </button>
        

        {/* Sign Up button navigates to /Signup */}
        <Button
          label="Sign Up"
          href="/signup"
          className="bg-green-600 hover:bg-green-700 text-white py-2 px-4 rounded text-center"
        />

        {/* Log In opens a modal */}
        <button
          onClick={() => setShowLogin(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded text-center"
        >
          Log In
        </button>
      </aside>

      {/* Login popup rendered conditionally */}
      {showLogin && (
        <LoginForm
          onClose={() => setShowLogin(false)}
          onSuccessLogin={() => setShowLogin(false)}
        />
      )}


      {showUpload && (
        <UploadForm
          onClose={() => setShowUpload(false)}
        />
      )}
    </>
  );
}

export default Sidebar;






// 'use client';

// import React, { useState } from 'react';
// import Button from './Button';
// import LoginForm from '../Auth/LoginForm';
// import UploadForm from './UploadForm';
// import { useAuth } from '../../contexts/AuthContext';

// export default function Sidebar() {
//   const { isLoggedIn, showLogin, setShowLogin, login } = useAuth();
//   const [showUpload, setShowUpload] = useState(false);

//   const handleUploadClick = () => {
//     if (!isLoggedIn) {
//       setShowLogin(true);
//     } else {
//       setShowUpload(true);
//     }
//   };

//   return (
//     <>
//       <aside className="fixed top-0 left-0 h-screen w-64 bg-black text-white p-4 flex flex-col gap-4">
//         <h2 className="text-xl font-bold">Sidebar</h2>

//         <button
//           onClick={handleUploadClick}
//           className="flex items-center gap-2 text-white py-2 px-4 rounded text-center bg-transparent hover:bg-white/10 transition"
//         >
//           <span className="text-2xl">⬆️</span>
//           <span>Upload</span>
//         </button>

//         <Button
//           label="Sign Up"
//           href="/signup"
//           className="bg-green-600 hover:bg-green-700 text-white py-2 px-4 rounded text-center"
//         />

//         <button
//           onClick={() => setShowLogin(true)}
//           className="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded text-center"
//         >
//           Log In
//         </button>
//       </aside>

//       {showLogin && (
//         <LoginForm
//           onClose={() => setShowLogin(false)}
//           onSuccessLogin={login}
//         />
//       )}

//       {showUpload && (
//         <UploadForm onClose={() => setShowUpload(false)} />
//       )}
//     </>
//   );
// }

