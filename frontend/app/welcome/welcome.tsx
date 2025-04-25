import { useEffect, useRef, useState, type SyntheticEvent } from "react";
import VideoPlayer from "./videoPlayer";

type Video = {
  name: string;
  path: string;
};

export function Welcome() {
  const [videoList, setVideoList] = useState<Array<Video>>([]);
  const [currentVideo, setCurrentVideo] = useState<Video>();
  const [isFullScreen, setIsFullScreen] = useState(false);
  const [isVideoDeleted, setIsVideoDeleted] = useState(false);
  const [isVideoPlaying, setIsVideoPlaying] = useState(false);
  const [lastWatchedEpisode, setLastWatchedEpisode] = useState<string>("");

  const videoElementRef = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    fetchVideos();
  }, []);

  const showDeleteMessage = () => {
    setIsVideoDeleted(true);
    setTimeout(() => {
      setIsVideoDeleted(false);
    }, 2000);
  };

  useEffect(() => {
    if (videoList.length > 0) {
      const lastEpisode =
        videoList[videoList.length - 1]?.name.match(/\d+/)?.[0];
      if (lastEpisode) {
        setLastWatchedEpisode(lastEpisode);
      }
    }
  }, [videoList]);

  const fetchVideos = async () => {
    const response = await fetch("/api/videos");
    const data = await response.json();
    setVideoList(data);
    if (currentVideo) {
      playVideo(data[0]);
    }
    return data;
  };

  const playVideo = (video: Video) => {
    setCurrentVideo(video);
    setIsVideoPlaying(true);
    if (videoElementRef.current) {
      const videoElement = videoElementRef.current;
      videoElement.src = video.path;
      videoElement.play();
    }
  };

  const skipIntro = () => {
    if (!videoElementRef) {
      return;
    }
    const videoElement = videoElementRef.current;
    if (videoElement) {
      videoElement.currentTime = 30; // Skip the first 30 seconds
    }
  };

  const skipOutro = () => {
    if (!videoElementRef) {
      return;
    }
    const videoElement = videoElementRef.current;
    if (videoElement) {
      videoElement.currentTime = videoElement.duration - 30; // Skip the last 30 seconds
    }
  };

  const nextVideo = () => {
    const currentIndex = videoList.findIndex(
      (video) => video.name === currentVideo?.name
    );
    const nextIndex = (currentIndex + 1) % videoList.length;
    playVideo(videoList[nextIndex]);
  };

  const deleteVideo = async (video?: Video) => {
    let v = currentVideo;
    if (video) {
      v = video as Video;
    }
    if (v === undefined) {
      return;
    }
    const response = await fetch(`/api/videos/${v.name}`, {
      method: "DELETE",
    });
    if (response.ok) {
      const index = videoList.findIndex((video) => video.name === v.name);
      if (index !== -1) {
        videoList.splice(index, 1);
      }
      setVideoList(videoList);
      setCurrentVideo(videoList[index]);
      showDeleteMessage();
    } else {
      console.error("Failed to delete video");
    }
  };

  const downloadVideos = async () => {
    const response = await fetch("/api/videos/download-next-10", {
      method: "POST",
      body: JSON.stringify({ lastDownloaded: lastWatchedEpisode }),
    });
    if (response.ok) {
      alert("Next 10 videos are being downloaded.");
      await fetchVideos();
    } else {
      console.error("Failed to download next 10 videos");
    }
  };

  const redownloadVideo = async (video: Video) => {
    if (video) {
      const response = await fetch(`/api/videos/download`, {
        method: "POST",
        body: JSON.stringify({ episode_id: video.name.match(/\d+/)?.[0] }),
      });
      if (response.ok) {
        alert(`Video ${video.name} is being redownloaded.`);
        fetchVideos(); // Refresh the video list after redownloading
      } else {
        console.error("Failed to redownload video");
      }
    }
  };

  const buttons = [
    {
      id: "skip-intro-button",
      className: "bg-green-500 hover:bg-green-600 text-white py-2 px-4 rounded",
      label: "Skip Intro",
      onClick: skipIntro,
    },
    {
      id: "skip-outro-button",
      className:
        "bg-yellow-500 hover:bg-yellow-600 text-white py-2 px-4 rounded",
      label: "Skip Outro",
      onClick: skipOutro,
    },
    {
      id: "next-button",
      className: "bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded",
      label: "Next",
      onClick: nextVideo,
    },
    {
      id: "delete-button",
      className: "bg-red-500 hover:bg-red-600 text-white py-2 px-4 rounded",
      label: "Delete",
      onClick: (e: SyntheticEvent<HTMLButtonElement>, v?: Video) =>
        deleteVideo(v),
    },
  ];

  return (
    <div className="flex flex-col h-screen">
      <>
        <div className="flex-1 bg-gray-900 text-white h-full p-4">
          <h1 className="text-2xl font-bold mb-4">One Piece Episodes</h1>
          <p className="text-sm mb-4">
            Select an episode to watch. Double click on the video to toggle full
            screen.
          </p>
          {isVideoDeleted && (
            <div className="fixed top-4 right-4 bg-red-500 text-white py-2 px-4 rounded shadow-lg z-50">
              Video deleted successfully.
            </div>
          )}
          {isVideoPlaying && (
            <p className="text-green-500">Now playing: {currentVideo?.name}</p>
          )}
        </div>
      </>
      <div className="flex flex-col md:flex-row h-full">
        <div className="w-64 bg-gray-800 text-white p-4">
          <h2 className="text-lg font-bold mb-4">Episodes</h2>
          <button
            onClick={fetchVideos}
            className="bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded mb-4 w-full"
          >
            Refresh List
          </button>
          <input
            type="text"
            placeholder="Enter last watched episode"
            value={lastWatchedEpisode}
            onChange={(e) => setLastWatchedEpisode(e.target.value)}
            className="w-full p-2 mb-4 rounded bg-gray-700 text-white"
          />
          <button
            onClick={downloadVideos}
            className="bg-purple-500 hover:bg-purple-600 text-white py-2 px-4 rounded mb-4 w-full"
          >
            Download Next 10 Videos
          </button>
          <ul id="video-list" className="space-y-2 overflow-y-auto max-h-96">
            {videoList.map((video, index) => (
              <li
                key={index}
                onClick={() => playVideo(video)}
                className={`flex items-center space-x-4 p-2 hover:bg-gray-700 cursor-pointer rounded ${
                  currentVideo?.name === video.name ? "bg-gray-700" : ""
                }`}
              >
                <span className="text-sm">{video.name.match(/\d+/)}</span>
                <button
                  className="p-2 bg-blue-500 hover:bg-blue-700 text-white rounded cursor-pointer ml-auto"
                  onClick={redownloadVideo.bind(null, video)}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="size-6"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M7.5 7.5h-.75A2.25 2.25 0 0 0 4.5 9.75v7.5a2.25 2.25 0 0 0 2.25 2.25h7.5a2.25 2.25 0 0 0 2.25-2.25v-7.5a2.25 2.25 0 0 0-2.25-2.25h-.75m-6 3.75 3 3m0 0 3-3m-3 3V1.5m6 9h.75a2.25 2.25 0 0 1 2.25 2.25v7.5a2.25 2.25 0 0 1-2.25 2.25h-7.5a2.25 2.25 0 0 1-2.25-2.25v-.75"
                    />
                  </svg>
                </button>
                <button
                  className="p-2 bg-red-500 hover:bg-red-700 text-white rounded cursor-pointer"
                  onClick={() => deleteVideo(video)}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="size-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0"
                    />
                  </svg>
                </button>
              </li>
            ))}
          </ul>
        </div>
        <div className="content">
          <h1 id="video-title">Select an episode</h1>
          {/* <video id="main-video" controls onClick="toggleFullScreen()">
            <source id="video-source" src="" type="video/mp4">
            Your browser does not support the video tag.
        </video> */}
          <VideoPlayer
            videoRef={videoElementRef}
            onDoubleClick={() => setIsFullScreen(!isFullScreen)}
            isFullScreen={isFullScreen}
            playlistUrl={currentVideo?.path}
          />
          <div className="flex justify-center space-x-4 mt-4">
            {buttons.map((button) => (
              <button
                key={button.id}
                id={button.id}
                className={button.className}
                onClick={(e) => button.onClick(e, currentVideo)}
              >
                {button.label}
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
