// 'use client';

// import { useState, useEffect } from 'react';
// import LoginForm from './Auth/LoginForm';
// import Sidebar from './Util/Sidebar';

// export default function Home() {
//   const [isLoggedIn, setIsLoggedIn] = useState(false);
//   const [showLoginPopup, setShowLoginPopup] = useState(false);

//   // Check for JWT token in localStorage on page load
//   useEffect(() => {
//     const token = localStorage.getItem('jwtToken');
//     if (token) {
//       setIsLoggedIn(true);
//     }
//   }, []);

//   const handleVideoClick = (videoNumber: number) => {
//     if (!isLoggedIn) {
//       setShowLoginPopup(true);
//     } else {
//       alert(`Now playing Video ${videoNumber}`);
//     }
//   };

//   return (
//     <div className="min-h-screen bg-gray-100">
//       <header className="bg-white shadow p-4">
//         <h1 className="text-3xl font-bold text-center text-gray-800">TokTik</h1>
//       </header>

//       <main className="p-6">
//         <Sidebar
//           isLoggedIn={isLoggedIn}
//           onSuccessLogin={() => setIsLoggedIn(true)}
//           onLogout={() => setIsLoggedIn(false)}
//         />
//         <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 pl-[250px]">
//           {Array.from({ length: 12 }).map((_, i) => (
//             <div
//               key={i}
//               className="bg-gray-300 h-48 rounded-lg shadow-md flex items-center justify-center text-lg font-semibold text-gray-700 hover:bg-gray-400 cursor-pointer transition"
//               onClick={() => handleVideoClick(i + 1)}
//             >
//               🎬 Video {i + 1}
//             </div>
//           ))}
//         </div>
//       </main>

//       {showLoginPopup && (
//         <LoginForm
//           onClose={() => setShowLoginPopup(false)}
//           onSuccessLogin={() => {
//             setIsLoggedIn(true);
//             setShowLoginPopup(false);
//           }}
//         />
//       )}
//     </div>
//   );
// }









// 'use client';

// import { useState, useEffect } from 'react';
// import LoginForm from './Auth/LoginForm';
// import Sidebar from './Util/Sidebar';

// export default function Home() {
//   const [isLoggedIn, setIsLoggedIn] = useState(false);
//   const [showLoginPopup, setShowLoginPopup] = useState(false);
//   const [jwtToken, setJwtToken] = useState<string | null>(null);
//   const [videos, setVideos] = useState<VideoThumbInfo[]>([]);
//   const [loading, setLoading] = useState(true);

  


//   useEffect(() => {
//     const token = localStorage.getItem('jwtToken');
//     if (token) {
//       setIsLoggedIn(true);
//       setJwtToken(token);
//     }
//   }, []);

//   const handleVideoClick = (videoNumber: number) => {
//     if (!isLoggedIn) {
//       setShowLoginPopup(true);
//     } else {
//       alert(`Now playing Video ${videoNumber}`);
//     }
//   };

//   return (
//     <div className="min-h-screen bg-gray-100">
//       <header className="bg-white shadow p-4">
//         <h1 className="text-3xl font-bold text-center text-gray-800">TokTik</h1>
//       </header>

//       <main className="p-6">
//         <Sidebar
//           isLoggedIn={isLoggedIn}
//           onSuccessLogin={(token: string) => {
//             setIsLoggedIn(true);
//             setJwtToken(token);
//           }}
//           onLogout={() => {
//             setIsLoggedIn(false);
//             setJwtToken(null);
//             localStorage.removeItem('jwtToken');
//             localStorage.removeItem('username');
//           }}
//         />

//         <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 pl-[250px]">
//           {Array.from({ length: 12 }).map((_, i) => (
//             <div
//               key={i}
//               className="bg-gray-300 h-48 rounded-lg shadow-md flex items-center justify-center text-lg font-semibold text-gray-700 hover:bg-gray-400 cursor-pointer transition"
//               onClick={() => handleVideoClick(i + 1)}
//             >
//               🎬 Video {i + 1}
//             </div>
//           ))}
//         </div>
//       </main>

//       {showLoginPopup && (
//         <LoginForm
//           onClose={() => setShowLoginPopup(false)}
//           onSuccessLogin={(token: string) => {
//             setIsLoggedIn(true);
//             setJwtToken(token);
//             setShowLoginPopup(false);
//           }}
//         />
//       )}
//     </div>
//   );
// }










'use client';

import { useState, useEffect, useRef } from 'react';
import LoginForm from './Auth/LoginForm';
import Sidebar from './Util/Sidebar';
import { ImEye} from 'react-icons/im';
import { BsFillHandThumbsUpFill } from "react-icons/bs";
import Hls from 'hls.js';

type VideoThumbInfo = {
  videoID: string;
  thumbnailURL: string;
  title: string; // camelCase!
  totalLikeCount: number;
  totalViewCount: number;
};


type VideoThumbInfoBackend = {
  videoID: string;
  thumbnailURL: string;
  title: string;
  TotalLikeCount: number;
  TotalViewCount: number;
};

// ...existing code...
type CommentDTO = {
  id: string;
  content: string;
  username: string;
};

type VideoDetailResponse = {
  playlist: string;
  title: string;
  totalViewCount: number;
  totalLikeCount: number;
  comments: CommentDTO[];
};
// ...existing code...




export default function Home() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [showLoginPopup, setShowLoginPopup] = useState(false);
  const [jwtToken, setJwtToken] = useState<string | null>(null);
  const [videos, setVideos] = useState<VideoThumbInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [playlist, setPlaylist] = useState<string | null>(null);
  const [playlistLoading, setPlaylistLoading] = useState(false);
  const [selectedVideo, setSelectedVideo] = useState<VideoDetailResponse | null>(null);

  // Fetch videos function
  const fetchVideos = async () => {
    setLoading(true);
    try {
      const res = await fetch('http://localhost:8090/api/videos');
      if (res.ok) {
        const data: VideoThumbInfoBackend[] = await res.json();
        console.log("Fetched videos from backend:", data); // <--- Add this
        const mapped: VideoThumbInfo[] = data.map((v) => ({
          videoID: v.videoID,
          thumbnailURL: v.thumbnailURL,
          title: v.title,
          totalLikeCount: v.TotalLikeCount ?? 0,
          totalViewCount: v.TotalViewCount ?? 0,
        }));
        setVideos(mapped);
        console.log("Mapped videos for display:", mapped); // <--- And this
      }
    } catch (err) {
      console.error('Failed to fetch videos:', err);
    }
    setLoading(false);
  };

  useEffect(() => {
    const token = localStorage.getItem('jwtToken');
    if (token) {
      setIsLoggedIn(true);
      setJwtToken(token);
    }
    fetchVideos();
  }, []);

  useEffect(() => {
    const interval = setInterval(() => {
      if (!showLoginPopup) {
        fetchVideos();
      }
    }, 10000);
    return () => clearInterval(interval);
  }, [showLoginPopup]);


  function HlsPlayer({ playlistUrl }: { playlistUrl: string }) {
    const videoRef = useRef<HTMLVideoElement>(null);

    useEffect(() => {
      if (videoRef.current && Hls.isSupported()) {
        const hls = new Hls();
        hls.loadSource(playlistUrl);
        hls.attachMedia(videoRef.current);
        hls.on(Hls.Events.ERROR, (event, data) => {
          console.error("HLS.js error:", data);
        });
        return () => {
          hls.destroy();
        };
      } else if (videoRef.current) {
        videoRef.current.src = playlistUrl;
      }
    }, [playlistUrl]);

    return (
      <video ref={videoRef} controls autoPlay className="w-full h-80 rounded bg-black" />
    );
  }


  // ...existing code...
const handleVideoClick = async (videoID: string) => {
  if (!isLoggedIn) {
    setShowLoginPopup(true);
    return;
  }
  setPlaylistLoading(true);
  setSelectedVideo(null);
  try {
    const res = await fetch(`http://localhost:8090/api/videos/${videoID}`, {
      headers: {
        Authorization: `Bearer ${jwtToken}`,
      },
    });
    if (res.ok) {
      const data: VideoDetailResponse = await res.json();
      console.log("Fetched video detail:", data); 
      setSelectedVideo(data);
    } else {
      alert("Failed to load video.");
    }
  } catch (err: unknown) {
    if (err instanceof Error) {
      console.error('Error fetching video:', err.message);
    } else {
      console.error('Error fetching video:', err);
    }
  }
  setPlaylistLoading(false);
};
// ...existing code...






  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow p-4">
        <h1 className="text-3xl font-bold text-center text-gray-800">TokTik</h1>
      </header>

      <main className="p-6">
        <Sidebar
          isLoggedIn={isLoggedIn}
          onSuccessLogin={(token: string) => {
            setIsLoggedIn(true);
            setJwtToken(token);
          }}
          onLogout={() => {
            setIsLoggedIn(false);
            setJwtToken(null);
            localStorage.removeItem('jwtToken');
            localStorage.removeItem('username');
          }}
        />

        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 pl-[250px]">
          {loading ? (
            <div>Loading...</div>
          ) : (
            videos.map((video) => (
              <div
                key={video.videoID}
                className="bg-gray-300 h-56 rounded-lg shadow-md flex flex-col items-center justify-between text-lg font-semibold text-gray-700 hover:bg-gray-400 cursor-pointer transition p-2"
                onClick={() => handleVideoClick(video.videoID)}
              >
                <img
                  src={video.thumbnailURL}
                  alt={`Thumbnail for video ${video.title}`}
                  className="w-full h-32 object-cover rounded-t-lg"
                />
                <div className="w-full flex flex-col items-start mt-2 px-2">
                  <span className="font-bold text-base truncate w-full text-center">{video.title}</span>
                  <div className="flex items-center gap-4 mt-1 text-sm">
                    <span title="Likes" className="flex items-center gap-1 ">
                      <span className="inline-block align-middle"><BsFillHandThumbsUpFill size={18} /></span> {video.totalLikeCount}
                    </span>
                    <span title="Views" className="flex items-center gap-1">
                      <span className="inline-block align-middle"><ImEye size={18} /></span> {video.totalViewCount}
                    </span>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </main>

      {showLoginPopup && (
        <LoginForm
          onClose={() => setShowLoginPopup(false)}
          onSuccessLogin={(token: string) => {
            setIsLoggedIn(true);
            setJwtToken(token);
            setShowLoginPopup(false);
          }}
        />
      )}


      
      {playlistLoading && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white p-8 rounded shadow text-center">
            Loading video...
          </div>
        </div>
      )}
      {playlist && (
        <div className="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center z-50">
          <div className="bg-white p-6 rounded shadow max-w-lg w-full">
            <h2 className="text-xl font-bold mb-4">Playlist.m3u8</h2>
            <pre className="overflow-x-auto text-xs bg-gray-100 p-2 rounded">{playlist}</pre>
            <button
              className="mt-4 px-4 py-2 bg-blue-600 text-white rounded"
              onClick={() => setPlaylist(null)}
            >
              Close
            </button>
          </div>
        </div>
      )}


      {selectedVideo && (
        <div className="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center z-50">
          <div className="bg-white rounded shadow max-w-4xl w-full flex flex-col md:flex-row p-4 relative">
            <button
              className="absolute top-2 right-2 text-gray-600 hover:text-black text-2xl"
              onClick={() => setSelectedVideo(null)}
            >
              &times;
            </button>

            {/* Left: Video Player */}
            <div className="md:w-2/3 w-full flex items-center justify-center">
              <HlsPlayer playlistUrl={selectedVideo.playlist} />
            </div>
            
            {/* Right: Info */}
            <div className="md:w-1/3 w-full flex flex-col justify-between p-4">
              <div>
                <h2 className="text-xl font-bold mb-2">{selectedVideo.title}</h2>
                <div className="flex gap-4 mb-4 text-gray-700">
                  <span className="flex items-center gap-1">
                    <BsFillHandThumbsUpFill size={18} /> {selectedVideo.totalLikeCount}
                  </span>
                  <span className="flex items-center gap-1">
                    <ImEye size={18} /> {selectedVideo.totalViewCount}
                  </span>
                </div>
              </div>
              <div className="mt-4">
                <h3 className="font-semibold mb-2">Comments</h3>
                <div className="max-h-48 overflow-y-auto space-y-2">
                  {(selectedVideo.comments ?? []).length === 0 ? (
                    <div className="text-gray-400">No comments yet.</div>
                  ) : (
                    (selectedVideo.comments ?? []).map((c) => (
                      <div key={c.id} className="border-b pb-1">
                        <span className="font-bold">{c.username}:</span> {c.content}
                      </div>
                    ))
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
      )}




    </div>
  );
}