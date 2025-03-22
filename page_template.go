package main

const pageTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Video Gallery</title>
    <style>
        body { font-family: Arial, sans-serif; display: flex; }
        .sidebar { width: 200px; background: #f4f4f4; padding: 10px; height: 100vh; overflow-y: auto; }
        .content { flex: 1; padding: 20px; }
        video { width: 100%; cursor: pointer; }
        .video-item { cursor: pointer; padding: 5px; border-bottom: 1px solid #ccc; }
        .video-item:hover { background: #ddd; }
    </style>
</head>
<body>
    <div class="sidebar">
        <h2>Video List</h2>
        <button onclick="fetchVideos()">Refresh List</button>
        <ul id="video-list">
        </ul>
    </div>
    <div class="content">
        <h1 id="video-title">Select a video</h1>
        <video id="main-video" controls ondblclick="toggleFullScreen()">
            <source id="video-source" src="" type="video/mp4">
            Your browser does not support the video tag.
        </video>
        <button id="delete-button" style="display:none;" onclick="deleteVideo()">Delete</button>
    </div>
    <script>
        function fetchVideos() {
            fetch('/videos')
                .then(response => response.json())
                .then(videos => {
                    let videoList = document.getElementById("video-list");
                    videoList.innerHTML = "";
                    videos.forEach(video => {
                        let li = document.createElement("li");
                        li.className = "video-item";
                        li.innerText = video.Name;
                        li.onclick = () => playVideo(video.Path, video.Name);
                        videoList.appendChild(li);
                    });

                    // Auto-play the first video or the last played video
                    let lastPlayedVideo = localStorage.getItem("lastPlayedVideo");
                    if (lastPlayedVideo) {
                        let lastVideo = videos.find(video => video.Path === lastPlayedVideo);
                        if (lastVideo) {
                            playVideo(lastVideo.Path, lastVideo.Name);
                        } else if (videos.length > 0) {
                            playVideo(videos[0].Path, videos[0].Name);
                        }
                    } else if (videos.length > 0) {
                        playVideo(videos[0].Path, videos[0].Name);
                    }
                });
        }

        function playVideo(videoPath, videoName) {
            let videoElement = document.getElementById("main-video");
            videoElement.src = videoPath;
            document.getElementById("video-title").innerText = videoName;
            document.getElementById("delete-button").style.display = "block";
            document.getElementById("delete-button").setAttribute("data-path", videoPath);

            // Save the last played video in localStorage
            localStorage.setItem("lastPlayedVideo", videoPath);

            // Skip first 3 minutes and ensure playback starts at 181 seconds if less than 180 seconds
            videoElement.onloadedmetadata = () => {
                if (videoElement.currentTime < 180) {
                    videoElement.currentTime = 181;
                }
                let duration = videoElement.duration;
                videoElement.ontimeupdate = () => {
                    if (videoElement.currentTime >= duration - 180) {
                        playNextVideo(); // Play the next video
                    }
                };
            };
        }

        function playNextVideo() {
            let videoList = document.getElementById("video-list").children;
            let currentVideoTitle = document.getElementById("video-title").innerText;

            for (let i = 0; i < videoList.length; i++) {
                if (videoList[i].innerText === currentVideoTitle && i + 1 < videoList.length) {
                    videoList[i + 1].click(); // Play the next video
                    break;
                }
            }
        }

        function deleteVideo() {
            let videoPath = document.getElementById("delete-button").getAttribute("data-path");
            fetch('/delete', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ path: videoPath })
            }).then(response => {
                if (response.ok) {
                    fetchVideos();
                    document.getElementById("main-video").src = "";
                    document.getElementById("video-title").innerText = "Select a video";
                    document.getElementById("delete-button").style.display = "none";
                } else {
                    alert("Failed to delete video");
                }
            });
        }

        function toggleFullScreen() {
            let video = document.getElementById("main-video");
            if (!document.fullscreenElement) {
                video.requestFullscreen().catch(err => {
                    console.log("Error attempting to enable full-screen mode:" + err.message);
                });
            } else {
                document.exitFullscreen();
            }
        }

        setInterval(fetchVideos, 900000); // Auto refresh every 15 minutes
        fetchVideos(); // Initial load
    </script>
</body>
</html>
`
