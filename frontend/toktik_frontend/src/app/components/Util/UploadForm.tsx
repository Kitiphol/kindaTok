'use client';

import { useState, useCallback, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function UploadModal({ onClose }) {
  const router = useRouter();

  const [videoFile, setVideoFile] = useState<File | null>(null);
  const [videoURL, setVideoURL]     = useState<string | null>(null);
  const [title, setTitle]           = useState('');
  const [uploading, setUploading]   = useState(false);
  const [error, setError]           = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  // 1) HANDLE FILE SELECTION & VALIDATION
  const handleFile = useCallback((file: File) => {
  setError(null);

  if (file.type !== 'video/mp4') {
    setError('Only MP4 allowed.');
    return;
  }

  const url = URL.createObjectURL(file);
  const tempVid = document.createElement('video');

  tempVid.preload = 'metadata';
  tempVid.src = url;

  tempVid.onloadedmetadata = () => {
    if (tempVid.duration > 65) {
      setError('Video must be ≤ 1 minute');
      setVideoFile(null);
      setVideoURL(null);
      URL.revokeObjectURL(url);
    } else {
      setVideoFile(file);
      setVideoURL(url); // do NOT revoke here!
    }
  };

  tempVid.onerror = () => {
    setError('Invalid video file.');
    setVideoFile(null);
    setVideoURL(null);
    URL.revokeObjectURL(url);
  };
}, []);


  // file input
  const onFileChange = (e) => {
    const f = e.target.files?.[0];
    if (f) handleFile(f);
  };

  // drag & drop
  const onDrop = (e) => {
    e.preventDefault();
    const f = e.dataTransfer.files?.[0];
    if (f) handleFile(f);
  };

  // 2) TWO-PHASE UPLOAD
  const handleUpload = async () => {
    if (!videoFile) return setError('Select a video first');
    if (!title.trim()) return setError('Enter a title');

    setUploading(true);
    setError(null);

    try {
      // Phase 1: get presigned URL
      //upload to '/api/uploadUser'
      const res1 = await fetch('/api/uploadUser', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ filename: videoFile.name, title: title.trim() })
      });
      if (!res1.ok) throw new Error('Could not get upload URL');
      const { url: presignedUrl } = await res1.json();

      // Phase 2: PUT the file to S3
      const res2 = await fetch(presignedUrl, {
        method: 'PUT',
        headers: { 'Content-Type': 'video/mp4' },
        body: videoFile
      });
      if (!res2.ok) throw new Error('Upload to storage failed');

      setSuccessMsg('Upload successful!');
      setVideoFile(null);
      setVideoURL(null);
      setTitle('');
    } catch (err: any) {
      setError(err.message || 'Upload failed');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="fixed inset-0 flex items-center justify-center z-50 p-4">
      {/* transparent overlay */}
      <div
        className="absolute inset-0"
        onClick={() => !uploading && router.push('/')}
      />

      <div
        className="relative bg-white rounded-lg p-6 max-w-lg w-full shadow-lg"
        onClick={e => e.stopPropagation()}
      >
        <button
          onClick={onClose}
          className="absolute top-2 right-2 text-gray-600 hover:text-black text-xl"
          disabled={uploading}
        >×</button>

        <h2 className="text-2xl font-semibold mb-4 text-center">Upload Your Video</h2>

        {/* drag & drop + file input */}
        <label
          htmlFor="video-file"
          onDrop={onDrop}
          onDragOver={e => e.preventDefault()}
          className="block mb-4 cursor-pointer border-2 border-dashed border-gray-300 rounded-md p-6 text-center hover:border-blue-500"
        >
          {!videoURL ? (
            <span className="text-gray-500">Select or drag & drop (MP4 ≤ 1 min)</span>
          ) : (
            <video
              key={videoURL}                // force remount when URL changes
              src={videoURL}
              controls
              className="max-w-full max-h-48 rounded"
            />
          )}
          <input
            id="video-file"
            type="file"
            accept="video/mp4"
            onChange={onFileChange}
            className="hidden"
            disabled={uploading}
          />
        </label>

        {/* info */}
        <div className="flex justify-around text-gray-600 text-sm mb-4">
          <div className="flex flex-col items-center">
            <span className="text-2xl">🎥</span><span>MP4 only</span>
          </div>
          <div className="flex flex-col items-center">
            <span className="text-2xl">⏱️</span><span>Max 1 min</span>
          </div>
          <div className="flex flex-col items-center">
            <span className="text-2xl">⏳</span><span>Wait upload</span>
          </div>
        </div>

        {/* title */}
        <div className="mb-4">
          <label htmlFor="video-title" className="block mb-1 text-gray-700">
            Video Title
          </label>
          <input
            id="video-title"
            type="text"
            value={title}
            onChange={e => setTitle(e.target.value)}
            className="w-full rounded border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={uploading}
          />
        </div>

        {/* messages */}
        {error     && <p className="text-red-600 text-center mb-2">{error}</p>}
        {successMsg&& <p className="text-green-600 text-center mb-2">{successMsg}</p>}

        {/* actions */}
        <div className="flex justify-between">
          <button
            onClick={onClose}
            disabled={uploading}
            className="px-4 py-2 border border-gray-400 rounded hover:bg-gray-100 disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            onClick={handleUpload}
            disabled={uploading}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
          >
            {uploading ? 'Uploading…' : 'Upload'}
          </button>
        </div>
      </div>
    </div>
  );
}
