



"use client";

import React, { useState, useEffect } from 'react';
import LoginForm from './Auth/LoginForm';
import Sidebar from './Util/Sidebar';
import { ImEye, ImArrowLeft2, ImArrowRight2 } from 'react-icons/im';
import { BsFillHandThumbsUpFill } from 'react-icons/bs';
// import Hls from 'hls.js';
import {HlsPlayer} from './HlsPlayer/hls';


export type VideoThumbInfo = {
  videoID: string;
  thumbnailURL: string;
  title: string;
  totalLikeCount: number;
  totalViewCount: number;
  hasLiked: boolean;
};

type VideoThumbInfoBackend = {
  videoID: string;
  thumbnailURL: string;
  title: string;
  TotalLikeCount: number;
  TotalViewCount: number;
};

type CommentDTO = {
  id: string;
  content: string;
  username: string;
  isUser: boolean;
};

type VideoDetailResponse = {
  playlist: string;
  title: string;
  totalViewCount: number;
  totalLikeCount: number;
  comments: CommentDTO[];
};

export default function Home() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [showLoginPopup, setShowLoginPopup] = useState(false);
  const [jwtToken, setJwtToken] = useState<string | null>(null);
  const [videos, setVideos] = useState<VideoThumbInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [playlistLoading, setPlaylistLoading] = useState(false);
  const [selectedVideo, setSelectedVideo] = useState<VideoDetailResponse | null>(null);
  const [selectedVideoIndex, setSelectedVideoIndex] = useState<number | null>(null);
  const [modalHasLiked, setModalHasLiked] = useState(false);
  const [newComment, setNewComment] = useState('');
  const [isSendingComment, setIsSendingComment] = useState(false);

  const fetchVideos = async (token: string | null) => {
    setLoading(true);
    try {

      const res = await fetch('http://localhost/api/video/videos');
      // const res = await fetch('http://localhost:8080/api/video/videos');

      // const res = await fetch('http://localhost:8090/api/video/videos');
      const list: VideoThumbInfoBackend[] = res.ok ? await res.json() : [];
      const sorted = [...list].sort((a, b) => a.videoID.localeCompare(b.videoID));
      const thumbs = await Promise.all(
        sorted.map(async (v) => {
          let likes = v.TotalLikeCount;
          let hasLiked = false;
          if (token) {


            const statusRes = await fetch(`http://localhost/api/ws/likes/videos/${v.videoID}`,
            // const statusRes = await fetch(`http://localhost:8080/api/ws/likes/videos/${v.videoID}`,

            // const statusRes = await fetch(`http://localhost:8092/api/ws/likes/videos/${v.videoID}`,

              
              
              { headers: { Authorization: `Bearer ${token}` } }
            );
            if (statusRes.ok) {
              const data = await statusRes.json();
              likes = data.likes;
              hasLiked = data.hasLiked;
            }
          }
          return {
            videoID: v.videoID,
            thumbnailURL: v.thumbnailURL,
            title: v.title,
            totalLikeCount: likes,
            totalViewCount: v.TotalViewCount,
            hasLiked,
          };
        })
      );
      setVideos(thumbs);
    } catch (err) {
      console.error('fetchVideos error:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const token = localStorage.getItem('jwtToken');
    if (token) {
      setIsLoggedIn(true);
      setJwtToken(token);
      fetchVideos(token);
    } else {
      fetchVideos(null);
    }
  }, []);

  useEffect(() => {
    const iv = setInterval(() => fetchVideos(jwtToken), 10000);
    return () => clearInterval(iv);
  }, [jwtToken]);

  const toggleLike = async (videoID: string, idx: number) => {
    if (!jwtToken) {
      setShowLoginPopup(true);
      return;
    }
    try {
      const currentlyLiked = videos[idx].hasLiked;

      const res = await fetch(`http://localhost/api/ws/likes/videos/${videoID}`, {
      // const res = await fetch(`http://localhost:8090/api/ws/likes/videos/${videoID}`, {

      // const res = await fetch(`http://localhost:8092/api/ws/likes/videos/${videoID}`, {
        method: currentlyLiked ? 'DELETE' : 'POST',
        headers: { Authorization: `Bearer ${jwtToken}` },
      });
      if (!res.ok) throw new Error('toggle failed');
      const { likes, hasLiked } = await res.json();
      setVideos((prev) =>
        prev.map((v, i) => (i === idx ? { ...v, totalLikeCount: likes, hasLiked } : v))
      );
      if (selectedVideoIndex === idx && selectedVideo) {
        setSelectedVideo({ ...selectedVideo, totalLikeCount: likes });
        setModalHasLiked(hasLiked);
      }
    } catch (err) {
      console.error('toggleLike error:', err);
      alert('Could not update like');
    }
  };

  const handleVideoClick = async (videoID: string, idx: number) => {
    if (!jwtToken) { setShowLoginPopup(true); return; }
    setPlaylistLoading(true);
    setSelectedVideo(null);
    setSelectedVideoIndex(idx);
    try {
      const [vRes, lRes, vwRes, cRes] = await Promise.all([

        fetch(`http://localhost/video/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`http://localhost/api/ws/likes/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`http://localhost/api/ws/views/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`http://localhost/api/ws/comments/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),


        // fetch(`http://localhost:8090/api/video/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        // fetch(`http://localhost:8080/api/ws/likes/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        // fetch(`http://localhost:8080/api/ws/views/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        // fetch(`http://localhost:8080/api/ws/comments/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),


        // fetch(`http://localhost:8090/api/video/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        // fetch(`http://localhost:8092/api/ws/likes/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        // fetch(`http://localhost:8092/api/ws/views/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        // fetch(`http://localhost:8092/api/ws/comments/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
      ]);
      if (!vRes.ok) throw new Error();
      const vd = await vRes.json();
      const ld = lRes.ok ? await lRes.json() : { likes: 0, hasLiked: false };
      const vw = vwRes.ok ? await vwRes.json() : { views: 0 };
      const cm = cRes.ok ? await cRes.json() : { comments: [] };
      setSelectedVideo({
        playlist: vd.playlist,
        title: vd.title,
        totalViewCount: vw.views,
        totalLikeCount: ld.likes,
        comments: Array.isArray(cm.comments) ? cm.comments : [],
      });

      setModalHasLiked(ld.hasLiked);
    } catch (err) {
      console.error('handleVideoClick error:', err);
      alert('Failed to load video');
    } finally {
      setPlaylistLoading(false);
    }
  };

  const handleSendComment = async () => {
    if (!jwtToken || !newComment.trim() || selectedVideoIndex === null) return;
    try {
      setIsSendingComment(true);
      const videoID = videos[selectedVideoIndex].videoID;

      const res = await fetch(`http://localhost/api/ws/comments/videos/${videoID}`, {
      // const res = await fetch(`http://localhost:8080/api/ws/comments/videos/${videoID}`, {

      // const res = await fetch(`http://localhost:8092/api/ws/comments/videos/${videoID}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${jwtToken}`,
        },
        body: JSON.stringify({ content: newComment }),
      });
      if (!res.ok) throw new Error();
      const added = await res.json();
      setSelectedVideo(prev => prev ? { ...prev, comments: [...prev.comments, added] } : prev);
      setNewComment('');
    } catch (err) {
      console.error('handleSendComment error:', err);
      alert('Failed to post comment');
    } finally {
      setIsSendingComment(false);
    }
  };

  const handleDeleteComment = async (commentID: string) => {
    if (!jwtToken || selectedVideoIndex === null) return;
    if (!confirm("Are you sure you want to delete this comment?")) return;
    try {
      const videoID = videos[selectedVideoIndex].videoID;


      const res = await fetch(`http://localhost/api/ws/comments/videos/${videoID}/${commentID}`, {
      // const res = await fetch(`http://localhost:8080/api/ws/comments/videos/${videoID}/${commentID}`, {

      // const res = await fetch(`http://localhost:8092/api/ws/comments/videos/${videoID}/${commentID}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${jwtToken}` },
      });
      if (res.ok) {
        setSelectedVideo(prev => prev ? { ...prev, comments: prev.comments.filter(c => c.id !== commentID) } : prev);
      } else {
        alert("Failed to delete comment");
      }
    } catch (err) {
      console.error("Error deleting comment:", err);
      alert("An error occurred while deleting.");
    }
  };


  // const HlsPlayer = ({ playlistUrl }: { playlistUrl: string }) => {
  //   const ref = useRef<HTMLVideoElement>(null);
  //   useEffect(() => {
  //     if (ref.current && Hls.isSupported()) {
  //       const hls = new Hls();
  //       hls.loadSource(playlistUrl);
  //       hls.attachMedia(ref.current);
  //       return () => hls.destroy();
  //     } else if (ref.current) {
  //       ref.current.src = playlistUrl;
  //     }
  //   }, [playlistUrl]);
  //   return <video ref={ref} controls autoPlay className="w-full h-80 rounded bg-black" />;
  // };


  

  return (



    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow p-4">
        <h1 className="text-3xl font-bold text-center text-black">TokTik</h1>
      </header>
      <main className="p-6 flex">
        <Sidebar
          isLoggedIn={isLoggedIn}
          onSuccessLogin={(token) => {
            localStorage.setItem('jwtToken', token);
            setJwtToken(token);
            setIsLoggedIn(true);
            fetchVideos(token);
          }}
          onLogout={() => {
            localStorage.removeItem('jwtToken');
            setJwtToken(null);
            setIsLoggedIn(false);
            fetchVideos(null);
          }}
        />
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 pl-[250px] w-full">
          {loading ? (
            <div>Loading...</div>
          ) : (
            videos.map((v, i) => (
              <div
                key={v.videoID}
                onClick={() => handleVideoClick(v.videoID, i)}
                className="bg-gray-300 h-56 rounded-lg shadow-md flex flex-col justify-between p-2 cursor-pointer hover:bg-gray-400"
              >
                <img
                  src={v.thumbnailURL}
                  alt={v.title}
                  className="w-full h-32 object-cover rounded-t"
                />
                <div className="text-lg font-bold text-gray-800 truncate px-1 text-center">
                  {v.title}
                </div>
                <div className="flex gap-4 items-center justify-center text-black">
                  <span
                    onClick={(e: React.MouseEvent) => {
                      e.stopPropagation();
                      toggleLike(v.videoID, i);
                    }}
                    className="cursor-pointer"
                  >
                    <BsFillHandThumbsUpFill
                      size={22}
                      color={v.hasLiked ? '#dc2626' : '#6b7280'}
                    />
                  </span>
                  <span>{v.totalLikeCount}</span>
                  <ImEye size={18} />
                  <span>{v.totalViewCount}</span>
                </div>
              </div>
            ))
          )}
        </div>
      </main>

      {showLoginPopup && <LoginForm onClose={() => setShowLoginPopup(false)} onSuccessLogin={(token) => { localStorage.setItem('jwtToken', token); setJwtToken(token); setIsLoggedIn(true); setShowLoginPopup(false); fetchVideos(token); }} />}

      {playlistLoading && <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center"><div className="bg-white p-6 rounded">Loading...</div></div>}

      {selectedVideo && selectedVideoIndex !== null && (
        <div className="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center">
          <div className="bg-white rounded shadow-lg w-11/12 max-w-4xl p-4 relative flex flex-col md:flex-row">
            {selectedVideoIndex > 0 && <button className="absolute left-2 top-1/2 -translate-y-1/2 text-3xl text-gray-600 hover:text-black" onClick={() => handleVideoClick(videos[selectedVideoIndex-1].videoID, selectedVideoIndex-1)}><ImArrowLeft2 /></button>}
            {selectedVideoIndex < videos.length-1 && <button className="absolute right-2 top-1/2 -translate-y-1/2 text-3xl text-gray-600 hover:text-black" onClick={() => handleVideoClick(videos[selectedVideoIndex+1].videoID, selectedVideoIndex+1)}><ImArrowRight2 /></button>}
            <button className="absolute top-2 right-2 text-2xl text-gray-600 hover:text-black" onClick={() => { setSelectedVideo(null); setSelectedVideoIndex(null); }}>&times;</button>
            <div className="md:w-2/3 w-full flex items-center justify-center"><HlsPlayer playlistUrl={selectedVideo.playlist} /></div>
            <div className="md:w-1/3 w-full p-4 flex flex-col text-black">
              <h2 className="text-xl font-bold truncate align-middle mb-2">{selectedVideo.title}</h2>
              <div className="flex gap-4 items-center mb-4 text-black">
                <span className="cursor-pointer" onClick={() => toggleLike(videos[selectedVideoIndex].videoID, selectedVideoIndex)}>
                  <BsFillHandThumbsUpFill size={22} color={modalHasLiked ? '#dc2626' : '#6b7280'} />
                </span>
                <span>{selectedVideo.totalLikeCount}</span>
                <ImEye size={18} /> <span>{selectedVideo.totalViewCount}</span>
              </div>
              <div className="flex flex-col h-full">
                <div className="flex-grow overflow-y-auto max-h-48 space-y-2">
                  <h3 className="font-semibold mb-2">Comments</h3>
                  {selectedVideo.comments?.length ? (
                    selectedVideo.comments.map((c) => (
                      <div key={c.id} className="border-b pb-1 flex justify-between items-center">
                        <span className="font-bold">{c.username}:</span> {c.content}

                        {c.isUser && (
                          <button onClick={() => handleDeleteComment(c.id)} className="text-red-600 hover:text-red-800 ml-2" title="Delete comment">🗑️</button>
                        )}
                      </div>
                    ))
                  ) : (
                    <div className="text-gray-400">No comments yet.</div>
                  )}
                </div>
                <div className="mt-2">
                  <textarea value={newComment} onChange={(e) => setNewComment(e.target.value)} placeholder="Add a comment..." className="w-full p-2 border border-gray-300 rounded resize-none" rows={2} />
                  <button onClick={handleSendComment} disabled={isSendingComment || !newComment.trim()} className="mt-2 px-4 py-1 bg-blue-600 text-white rounded disabled:opacity-50">
                    {isSendingComment ? 'Sending...' : 'Send'}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

