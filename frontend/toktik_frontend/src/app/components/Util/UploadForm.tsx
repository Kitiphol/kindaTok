'use client';

import { useState, useCallback, useRef } from 'react';
import { useRouter } from 'next/navigation';



export default function UploadModal({ onClose }) {
  const router = useRouter();

  const [videoFile, setVideoFile] = useState<File | null>(null);
  const [videoURL, setVideoURL] = useState<string | null>(null);
  const [title, setTitle] = useState('');
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  // Ref for hidden video element to measure duration
  const videoRef = useRef<HTMLVideoElement>(null);

  // When user selects file, check if it's .mp4, and duration <= 60 seconds
  const handleFile = useCallback((file: File) => {
    setError(null);
    setSuccessMsg(null);

    if (file.type !== 'video/mp4') {
      setError('Only MP4 files are allowed.');
      setVideoFile(null);
      return;
    }

    const url = URL.createObjectURL(file);
    setVideoURL(url);

    // Create a video element to check duration
    const tempVideo = document.createElement('video');
    tempVideo.preload = 'metadata';
    tempVideo.src = url;
    tempVideo.onloadedmetadata = () => {
      window.URL.revokeObjectURL(url);
      if (tempVideo.duration > 60) {
        setError('Video length must be 1 minute or less.');
        setVideoFile(null);
        setVideoURL(null);
      } else {
        setVideoFile(file);
      }
    };
    tempVideo.onerror = () => {
      setError('Could not load video metadata. Please select a valid video.');
      setVideoFile(null);
      setVideoURL(null);
    };
  }, []);

  const onFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      handleFile(e.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (!videoFile) {
      setError('Please select a video file.');
      return;
    }
    if (!title.trim()) {
      setError('Please enter a video title.');
      return;
    }
    setError(null);
    setSuccessMsg(null);
    setUploading(true);

    try {
      const formData = new FormData();
      formData.append('video', videoFile);
      formData.append('title', title.trim());

      const res = await fetch('/api/upload', {
        method: 'POST',
        body: formData,
      });

      if (!res.ok) {
        throw new Error(`Upload failed: ${res.statusText}`);
      }

      setSuccessMsg('Upload successful!');
      setVideoFile(null);
      setTitle('');
      setVideoURL(null);
    } catch (err: any) {
      setError(err.message || 'Upload failed.');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div
      className="fixed inset-0 flex items-center justify-center z-50 p-4"
      role="dialog"
      aria-modal="true"
      aria-labelledby="upload-modal-title"
    >
      {/* Transparent overlay */}
      <div
        className="absolute inset-0"
        onClick={() => {
          if (!uploading) router.push('/');
        }}
      />

      {/* Modal form container */}
      <div
        className="relative bg-white rounded-lg p-6 max-w-lg w-full shadow-lg"
        onClick={(e) => e.stopPropagation()}
      >
        <h2
          id="upload-modal-title"
          className="text-2xl font-semibold mb-4 text-center text-gray-900"
        >
          Upload Your Video
        </h2>

        <div className="mb-1">
          <label
            htmlFor="video-file"
            className="block mb-2 cursor-pointer border-2 border-dashed border-gray-300 rounded-md p-6 text-center hover:border-blue-500"
          >
            {!videoURL ? (
              <span className="text-gray-500">
                Select or drag & drop a video file (MP4, ≤1 min)
              </span>
            ) : (
              <video
                src={videoURL}
                controls
                className="max-w-full max-h-48 rounded"
                ref={videoRef}
              />
            )}
            <input
              type="file"
              id="video-file"
              accept="video/mp4"
              onChange={onFileChange}
              className="hidden"
              disabled={uploading}
            />
          </label>
        </div>

        <button
          onClick={onClose}
          className="absolute top-2 right-2 text-gray-600 hover:text-black text-xl"
        >
          ×
        </button>

        {/* Icons & requirements under drag-drop box */}
        <div className="flex justify-around text-gray-600 text-sm mb-4 select-none">
          <div className="flex flex-col items-center gap-1">
            <span
              aria-label="File type"
              title="MP4 format"
              role="img"
              className="text-2xl"
            >
              🎥
            </span>
            <span>MP4 only</span>
          </div>
          <div className="flex flex-col items-center gap-1">
            <span
              aria-label="Max duration"
              title="Max 1 minute"
              role="img"
              className="text-2xl"
            >
              ⏱️
            </span>
            <span>Max 1 min</span>
          </div>
          <div className="flex flex-col items-center gap-1">
            <span
              aria-label="Wait for processing"
              title="Wait for video to finish uploading"
              role="img"
              className="text-2xl"
            >
              ⏳
            </span>
            <span>Wait for upload</span>
          </div>
        </div>

        <div className="mb-4">
          <label htmlFor="video-title" className="block mb-1 text-gray-700">
            Video Title
          </label>
          <input
            id="video-title"
            type="text"
            placeholder="Enter video title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full rounded border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={uploading}
          />
        </div>

        {error && (
          <p className="mb-2 text-red-600 text-center">{error}</p>
        )}
        {successMsg && (
          <p className="mb-2 text-green-600 text-center">{successMsg}</p>
        )}

        <div className="flex justify-between items-center">
          <button
            onClick={onClose}
            className="px-4 py-2 rounded border border-gray-400 hover:bg-gray-100 transition disabled:opacity-50"
            disabled={uploading}
          >
            Cancel
          </button>

          <button
            onClick={handleUpload}
            className={`px-4 py-2 rounded bg-blue-600 text-white hover:bg-blue-700 transition disabled:opacity-50`}
            disabled={uploading}
          >
            {uploading ? 'Uploading...' : 'Upload'}
          </button>
        </div>
      </div>
    </div>
  );
}
