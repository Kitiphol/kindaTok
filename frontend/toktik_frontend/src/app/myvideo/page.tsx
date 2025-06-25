'use client';

import React, { useEffect, useState } from 'react';
import { ImEye, ImArrowLeft2, ImArrowRight2 } from 'react-icons/im';
import { BsFillHandThumbsUpFill } from 'react-icons/bs';
import { useRouter } from 'next/navigation';
import { HlsPlayer } from '../components/HlsPlayer/hls'; // Your existing HlsPlayer component
type VideoThumbInfo = {
  videoID: string;
  thumbnailURL: string;
  title: string;
  totalLikeCount: number;
  totalViewCount: number;
  hasLiked: boolean;
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

export default function MyVideoPage() {
  const router = useRouter();

  const [videos, setVideos] = useState<VideoThumbInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Modal and video details states
  const [selectedVideo, setSelectedVideo] = useState<VideoDetailResponse | null>(null);
  const [selectedVideoIndex, setSelectedVideoIndex] = useState<number | null>(null);
  const [modalHasLiked, setModalHasLiked] = useState(false);
  const [newComment, setNewComment] = useState('');
  const [isSendingComment, setIsSendingComment] = useState(false);
  const [playlistLoading, setPlaylistLoading] = useState(false);

  const jwtToken = typeof window !== 'undefined' ? localStorage.getItem('jwtToken') : null;

  // Fetch user's videos on mount
  useEffect(() => {
    const fetchMyVideos = async () => {
      setLoading(true);
      setError(null);

      if (!jwtToken) {
        setError('You must be logged in to view your videos.');
        setLoading(false);
        return;
      }

      try {
        // const res = await fetch('http://localhost:8090/api/videos/user', {
        const res = await fetch('http://localhost:8090/api/videos/user', {
          headers: { Authorization: `Bearer ${jwtToken}` },
        });
        if (!res.ok) throw new Error(`Error fetching videos: ${res.statusText}`);
        const data: VideoThumbInfo[] = await res.json();
        setVideos(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    };

    fetchMyVideos();
  }, [jwtToken]);

  // Delete video handler
  const handleDeleteVideo = async (videoID: string) => {
    if (!jwtToken) return alert('You must be logged in.');

    if (!confirm('Are you sure you want to delete this video?')) return;

    try {
    
        // const res = await fetch(`http://localhost:8090/api/videos/${videoID}`, {
      const res = await fetch(`http://localhost:8090/api/videos/${videoID}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${jwtToken}` },
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || 'Failed to delete video');
      }

      // Remove from UI
      setVideos((prev) => prev.filter((v) => v.videoID !== videoID));

      // Close modal if the deleted video is open
      if (selectedVideoIndex !== null && videos[selectedVideoIndex].videoID === videoID) {
        setSelectedVideo(null);
        setSelectedVideoIndex(null);
      }
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Unknown error');
    }
  };

  // Open video modal & fetch video details + comments + likes + views
  const handleVideoClick = async (videoID: string, idx: number) => {
    if (!jwtToken) {
      alert('You must be logged in to watch videos.');
      return;
    }

    setPlaylistLoading(true);
    setSelectedVideo(null);
    setSelectedVideoIndex(idx);

    try {
      const [vRes, lRes, vwRes, cRes] = await Promise.all([
        // fetch(`http://localhost:8090/api/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`http://localhost:8090/api/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),

        // fetch(`http://localhost:8092/api/likes/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`http://localhost:8092/api/likes/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),

// fetch(`http://localhost:8092/api/views/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`http://localhost:8092/api/views/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),

// fetch(`http://localhost:8092/api/comments/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`http://localhost:8092/api/comments/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
      ]);

      if (!vRes.ok) throw new Error('Failed to load video details');

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
    } catch (err: unknown) {
    console.error('Error loading video details:', err);
      alert('Failed to load video');
      setSelectedVideo(null);
      setSelectedVideoIndex(null);
    } finally {
      setPlaylistLoading(false);
    }
  };

  // Like/unlike toggle
  const toggleLike = async (videoID: string, idx: number) => {
    if (!jwtToken) {
      alert('You must be logged in.');
      return;
    }

    try {
      const currentlyLiked = videos[idx].hasLiked;
      // const res = await fetch(`http://localhost:8092/api/likes/videos/${videoID}`, {
      const res = await fetch(`http://localhost:8092/api/likes/videos/${videoID}`, {
        method: currentlyLiked ? 'DELETE' : 'POST',
        headers: { Authorization: `Bearer ${jwtToken}` },
      });

      if (!res.ok) throw new Error('Toggle like failed');

      const { likes, hasLiked } = await res.json();

      // Update videos list likes
      setVideos((prev) =>
        prev.map((v, i) => (i === idx ? { ...v, totalLikeCount: likes, hasLiked } : v))
      );

      // Update modal likes if open for this video
      if (selectedVideoIndex === idx && selectedVideo) {
        setSelectedVideo({ ...selectedVideo, totalLikeCount: likes });
        setModalHasLiked(hasLiked);
      }
    } catch {
      alert('Could not update like');
    }
  };

  // Send a comment
  const handleSendComment = async () => {
    if (!jwtToken || !newComment.trim() || selectedVideoIndex === null) return;

    try {
      setIsSendingComment(true);

      const videoID = videos[selectedVideoIndex].videoID;

      // const res = await fetch(`http://localhost:8092/api/comments/videos/${videoID}`, {
      const res = await fetch(`http://localhost:8092/api/comments/videos/${videoID}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${jwtToken}`,
        },
        body: JSON.stringify({ content: newComment }),
      });

      if (!res.ok) throw new Error('Failed to post comment');

      const added = await res.json();

      setSelectedVideo((prev) =>
        prev ? { ...prev, comments: [...prev.comments, added] } : prev
      );

      setNewComment('');
    } catch {
      alert('Failed to post comment');
    } finally {
      setIsSendingComment(false);
    }
  };

  // Delete a comment
  const handleDeleteComment = async (commentID: string) => {
    if (!jwtToken || selectedVideoIndex === null) return;
    if (!confirm('Are you sure you want to delete this comment?')) return;

    try {
      const videoID = videos[selectedVideoIndex].videoID;

      // const res = await fetch(`http://localhost:8092/api/comments/videos/${videoID}/${commentID}`, {
      const res = await fetch(`http://localhost:8092/api/comments/videos/${videoID}/${commentID}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${jwtToken}` },
      });

      if (res.ok) {
        setSelectedVideo((prev) =>
          prev ? { ...prev, comments: prev.comments.filter((c) => c.id !== commentID) } : prev
        );
      } else {
        alert('Failed to delete comment');
      }
    } catch {
      alert('An error occurred while deleting the comment.');
    }
  };

  if (loading) return <div className="p-4">Loading your videos...</div>;
  if (error) return <div className="p-4 text-red-600">Error: {error}</div>;
  if (!videos.length) return <div className="p-4">You have no uploaded videos.</div>;

  return (
    <div className="p-4 min-h-screen bg-gray-100">
      <div className="flex justify-end mb-4">
            <button
            onClick={() => router.push('/')}
            className="text-gray-500 hover:text-red-600 text-xl font-bold"
            aria-label="Close"
            >
            ❌
            </button>
        </div>

        <h1 className="text-2xl font-bold mb-4 text-black">My Videos</h1>

        <div className="flex flex-wrap gap-4">
            {videos.map((video, idx) => (
            <div
                key={video.videoID}
                className="relative border rounded p-2 shadow hover:shadow-lg cursor-pointer flex-shrink-0 flex flex-col justify-between bg-white"
                style={{
                    width: '180px',
                    height: '180px',
                }}
                onClick={() => handleVideoClick(video.videoID, idx)}
                >
                {/* Delete button */}
                <button
                    onClick={(e) => {
                    e.stopPropagation();
                    handleDeleteVideo(video.videoID);
                    }}
                    className="absolute top-2 right-2 text-sm text-red-500 hover:text-red-700 z-10"
                    title="Delete video"
                >
                    ✖
                </button>

                <img
                    src={video.thumbnailURL}
                    alt={video.title}
                    className="w-full h-24 object-cover rounded"
                />

                <h2 className="mt-2 font-semibold truncate text-black">{video.title}</h2>

                <div className="flex justify-between mt-1 text-sm text-gray-600 items-center">
                    <span
                    onClick={(e) => {
                        e.stopPropagation();
                        toggleLike(video.videoID, idx);
                    }}
                    className="flex items-center gap-1 cursor-pointer select-none"
                    >
                    <BsFillHandThumbsUpFill color={video.hasLiked ? '#dc2626' : '#6b7280'} />
                    {video.totalLikeCount}
                    </span>
                    <span className="flex items-center gap-1">
                    <ImEye />
                    {video.totalViewCount}
                    </span>
                </div>
                </div>


            ))}
        </div>


      {/* Modal Video Player */}
      {playlistLoading && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white p-6 rounded">Loading video...</div>
        </div>
      )}

      {selectedVideo && selectedVideoIndex !== null && (
        <div className="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded shadow-lg max-w-5xl w-full flex flex-col md:flex-row relative">
            {/* Navigation Buttons */}
            {selectedVideoIndex > 0 && (
              <button
                className="absolute left-2 top-1/2 -translate-y-1/2 text-3xl text-gray-600 hover:text-black"
                onClick={() => handleVideoClick(videos[selectedVideoIndex - 1].videoID, selectedVideoIndex - 1)}
                aria-label="Previous video"
              >
                <ImArrowLeft2 />
              </button>
            )}
            {selectedVideoIndex < videos.length - 1 && (
              <button
                className="absolute right-2 top-1/2 -translate-y-1/2 text-3xl text-gray-600 hover:text-black"
                onClick={() => handleVideoClick(videos[selectedVideoIndex + 1].videoID, selectedVideoIndex + 1)}
                aria-label="Next video"
              >
                <ImArrowRight2 />
              </button>
            )}

            {/* Close modal */}
            <button
              className="absolute top-2 right-2 text-2xl text-gray-600 hover:text-black"
              onClick={() => {
                setSelectedVideo(null);
                setSelectedVideoIndex(null);
                setNewComment('');
              }}
              aria-label="Close video modal"
            >
              &times;
            </button>

            {/* Video player */}
            <div className="md:w-2/3 w-full flex items-center justify-center bg-black rounded">
              <HlsPlayer playlistUrl={selectedVideo.playlist} />
            </div>

            {/* Video info & comments */}
            <div className="md:w-1/3 w-full p-4 flex flex-col">
              <h2 className="text-xl font-bold truncate mb-2">{selectedVideo.title}</h2>

              {/* Like and views */}
              <div className="flex gap-4 items-center mb-4">
                <span
                  onClick={() => toggleLike(videos[selectedVideoIndex].videoID, selectedVideoIndex)}
                  className="cursor-pointer"
                  aria-label={modalHasLiked ? 'Unlike video' : 'Like video'}
                >
                  <BsFillHandThumbsUpFill size={22} color={modalHasLiked ? '#dc2626' : '#6b7280'} />
                </span>
                <span>{selectedVideo.totalLikeCount}</span>
                <ImEye size={18} />
                <span>{selectedVideo.totalViewCount}</span>
              </div>

              {/* Comments section */}
              <div className="flex flex-col flex-grow">
                <h3 className="font-semibold mb-2">Comments</h3>
                <div className="flex-grow overflow-y-auto max-h-48 space-y-2 border border-gray-300 p-2 rounded">
                  {selectedVideo.comments.length ? (
                    selectedVideo.comments.map((c) => (
                      <div
                        key={c.id}
                        className="flex justify-between items-center border-b pb-1"
                      >
                        <span>
                          <strong>{c.username}:</strong> {c.content}
                        </span>
                        {c.isUser && (
                          <button
                            onClick={() => handleDeleteComment(c.id)}
                            className="text-red-600 hover:text-red-800 ml-2"
                            title="Delete comment"
                            aria-label="Delete comment"
                          >
                            🗑️
                          </button>
                        )}
                      </div>
                    ))
                  ) : (
                    <div className="text-gray-400">No comments yet.</div>
                  )}
                </div>

                <textarea
                  value={newComment}
                  onChange={(e) => setNewComment(e.target.value)}
                  placeholder="Add a comment..."
                  className="mt-2 p-2 border border-gray-300 rounded resize-none"
                  rows={2}
                  aria-label="Add a comment"
                />

                <button
                  onClick={handleSendComment}
                  disabled={isSendingComment || !newComment.trim()}
                  className="mt-2 px-4 py-1 bg-blue-600 text-white rounded disabled:opacity-50"
                >
                  {isSendingComment ? 'Sending...' : 'Send'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
