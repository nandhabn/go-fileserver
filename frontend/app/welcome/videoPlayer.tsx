import React, { useEffect, useRef } from "react";
import Hls from "hls.js";

interface VideoPlayerProps {
  playlistUrl?: string;
  isFullScreen: boolean;
  onDoubleClick: () => void;
  videoRef?: React.RefObject<HTMLVideoElement| null>;
}

const VideoPlayer: React.FC<VideoPlayerProps> = ({
  playlistUrl,
  onDoubleClick,
  isFullScreen,
  videoRef: externalVideoRef,
}) => {
  const internalVideoRef = useRef<HTMLVideoElement>(null);
  const videoRef = externalVideoRef || internalVideoRef;

  // Initialize HLS or native HLS playback
  useEffect(() => {
    if (!videoRef.current || !playlistUrl) return;

    if (Hls.isSupported()) {
        console.log("HLS is supported");
      const hls = new Hls();
      hls.loadSource(playlistUrl);
      hls.attachMedia(videoRef.current);

      hls.on(Hls.Events.MANIFEST_LOADED, () => {
        videoRef.current?.play();
      });

      return () => hls.destroy(); // Cleanup Hls.js instance on unmount
    } else if (videoRef.current.canPlayType("application/vnd.apple.mpegurl")) {
      // Safari: Native HLS support
      videoRef.current.src = playlistUrl;
      videoRef.current.addEventListener("loadedmetadata", () => {
        videoRef.current?.play();
      });
    }
  }, [playlistUrl, videoRef]);

  return (
    <div>
      <h1>HLS Streaming Player</h1>
      <video
        ref={videoRef}
        id="main-video"
        className={`${
          isFullScreen
            ? "fixed top-0 left-0 w-full h-full z-50"
            : "w-3/4 h-3/4"
        }`}
        controls
        onDoubleClick={onDoubleClick}
      >
        Your browser does not support HLS streaming.
      </video>
    </div>
  );
};

export default VideoPlayer;
