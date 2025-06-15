'use client';

import {
  useState,
  useCallback,
  ChangeEvent,
  DragEvent,
  MouseEvent,
} from 'react';
import { useRouter } from 'next/navigation';

interface UploadModalProps {
  onClose: () => void;
}

export default function UploadModal({ onClose }: UploadModalProps) {
  const router = useRouter();

  const [videoFile, setVideoFile] = useState<File | null>(null);
  const [videoURL, setVideoURL] = useState<string | null>(null);
  const [title, setTitle] = useState('');
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  const handleFile = useCallback((file: File) => {
    setError(null);
    setSuccessMsg(null);

    // basic type check
    if (file.type !== 'video/mp4') {
      setError('Only MP4 files are allowed.');
      return;
    }

    const url = URL.createObjectURL(file);
    const tempVid = document.createElement('video');
    tempVid.preload = 'metadata';
    tempVid.src = url;

    tempVid.onloadedmetadata = () => {
      if (tempVid.duration > 60) {
        setError('Video must be 60 seconds or shorter.');
        URL.revokeObjectURL(url);
      } else {
        setVideoFile(file);
        setVideoURL(url);
      }
    };

    tempVid.onerror = () => {
      setError('Could not load the selected video.');
      URL.revokeObjectURL(url);
    };
  }, []);

  const onFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const f = e.target.files?.[0];
    if (f) handleFile(f);
  };

  const onDrop = (e: DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
    const f = e.dataTransfer.files?.[0];
    if (f) handleFile(f);
  };

  const handleUpload = async () => {
    if (!videoFile) return setError('Please select a video first.');
    if (!title.trim()) return setError('Please enter a title.');

    setUploading(true);
    setProgress(0);
    setError(null);

    try {
      // 1) get presigned URL
      const res1 = await fetch('/api/videos', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          filename: videoFile.name,
          title: title.trim(),
        }),
      });
      if (!res1.ok) throw new Error('Failed to get upload URL');
      const { presignedURL } = await res1.json();

      // 2) upload via XHR so we can track progress
      await new Promise<void>((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.open('PUT', presignedURL);
        xhr.setRequestHeader('Content-Type', 'video/mp4');

        xhr.upload.onprogress = (event) => {
          if (event.lengthComputable) {
            setProgress(Math.round((event.loaded / event.total) * 100));
          }
        };

        xhr.onload = () => {
          if (xhr.status >= 200 && xhr.status < 300) {
            resolve();
          } else {
            reject(new Error('Upload failed with status ' + xhr.status));
          }
        };
        xhr.onerror = () => reject(new Error('Network error during upload'));
        xhr.send(videoFile);
      });

      setSuccessMsg('Upload successful!');
      setVideoFile(null);
      setVideoURL(null);
      setTitle('');
    } catch (err: unknown) {
      if (err instanceof Error) {
        alert(`Error: ${err.message}`);
      } else {
        alert('An unknown error occurred');
      }
    }
  };

  return (
    <div className="fixed inset-0 flex items-center justify-center z-50 p-4">
      {/* backdrop */}
      <div
        className="absolute inset-0 bg-black/50"
        onClick={() => !uploading && router.push('/')}
      />

      <div
        className="relative bg-white rounded-lg p-6 max-w-lg w-full shadow-lg z-10"
        onClick={(e: MouseEvent) => e.stopPropagation()}
      >
        <button
          onClick={onClose}
          disabled={uploading}
          className="absolute top-2 right-2 text-gray-600 hover:text-black text-xl"
        >
          ×
        </button>

        <h2 className="text-2xl font-semibold mb-4 text-center">
          Upload Your Video
        </h2>

        <label
          htmlFor="video-file"
          onDrop={onDrop}
          onDragOver={(e) => e.preventDefault()}
          className="block mb-4 cursor-pointer border-2 border-dashed border-gray-300 rounded-md p-4 text-center hover:border-blue-500"
        >
          {videoURL ? (
            <video
              src={videoURL}
              controls
              className="max-w-full max-h-48 rounded mb-2"
            />
          ) : (
            <span className="text-gray-500">
              Drag & drop or click to select a video
            </span>
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

        {/* requirements bar */}
        <div className="flex justify-around text-gray-600 text-sm mb-4">
          <div className="flex flex-col items-center">
            <span className="text-2xl">🎥</span>
            <span>MP4 only</span>
          </div>
          <div className="flex flex-col items-center">
            <span className="text-2xl">⏱️</span>
            <span>≤ 60 sec</span>
          </div>
        </div>

        {/* title input */}
        <div className="mb-4">
          <label htmlFor="video-title" className="block mb-1 text-gray-700">
            Video Title
          </label>
          <input
            id="video-title"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full rounded border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
            disabled={uploading}
          />
        </div>

        {/* progress */}
        {uploading && (
          <div className="w-full bg-gray-200 rounded h-4 mb-4">
            <div
              className="h-full bg-blue-600 rounded"
              style={{ width: `${progress}%` }}
            />
          </div>
        )}

        {/* feedback */}
        {error && <p className="text-red-600 text-center mb-2">{error}</p>}
        {successMsg && (
          <p className="text-green-600 text-center mb-2">{successMsg}</p>
        )}

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
            {uploading ? `Uploading ${progress}%` : 'Upload'}
          </button>
        </div>
      </div>
    </div>
  );
}
