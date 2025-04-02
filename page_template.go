package main

const pageTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>One Piece - Your Personal Video Server</title>
    <style>
        body { font-family: Arial, sans-serif; display: flex; height: 100vh; margin: 0; }
        .sidebar { width: 200px; background: #f4f4f4; padding: 10px; height: 100vh; overflow-y: auto; }
        .content { flex: 1; padding: 20px; display: flex; flex-direction: column; }
        #main-video { flex: 1; width: 100%; height: 100%; object-fit: contain; }
        .video-item { cursor: pointer; padding: 5px; border-bottom: 1px solid #ccc; }
        .video-item:hover { background: #ddd; }
        #delete-button { margin-top: 20px; width: 102px; height: 36px; background: red; color: white; border: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="sidebar">
        <h2>One Piece Episodes</h2>
        <button onclick="fetchVideos()">Refresh List</button>
        <ul id="video-list">
        </ul>
    </div>
    <div class="content">
        <h1 id="video-title">Select an episode</h1>
        <video id="main-video" controls ondblclick="toggleFullScreen()">
            <source id="video-source" src="" type="video/mp4">
            Your browser does not support the video tag.
        </video>
        <button id="delete-button" style="display:none;" onclick="deleteVideo()">Delete</button>
    </div>
    <script>
        async function fetchVideos() {
            const response = await fetch('/videos');
            const videos = await response.json();
            const watchedResponse = await fetch('/watched');
            const watchedVideos = await watchedResponse.json();
            const videoList = document.getElementById("video-list");
            videoList.innerHTML = "";
            videos.forEach(video => {
                const li = document.createElement("li");
                li.className = "video-item";
                li.innerText = video.Name + (watchedVideos[video.Path] ? " (Watched)" : "");
                li.onclick = () => playVideo(video.Path, video.Name);
                videoList.appendChild(li);
            });
        }

        async function playVideo(path, name) {
            const videoElement = document.getElementById("main-video");
            const sourceElement = document.getElementById("video-source");

            sourceElement.src = path;
            videoElement.load();
            videoElement.play();

            // Skip first 3 minutes
            videoElement.currentTime = 180;

            videoElement.addEventListener("timeupdate", () => {
                // Skip last 3 minutes
                if (videoElement.duration - videoElement.currentTime <= 180) {
                    videoElement.pause();
                }
            });

            // Mark video as watched
            await fetch('/mark-watched', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ Path: path })
            });

            document.getElementById("video-title").innerText = name;
            localStorage.setItem("lastPlayedVideo", path);
            document.getElementById("delete-button").style.display = "block";
        }

        async function deleteVideo() {
            const videoElement = document.getElementById("main-video");
            const sourceElement = document.getElementById("video-source");
            const videoPath = sourceElement.src;

            if (!videoPath) {
                alert("No video selected to delete.");
                return;
            }

            const confirmDelete = confirm("Are you sure you want to delete this video?");
            if (!confirmDelete) return;

            await fetch('/delete', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ Path: videoPath })
            });

            alert("Video deleted successfully.");
            sourceElement.src = "";
            videoElement.load();
            document.getElementById("video-title").innerText = "Select an episode";
            document.getElementById("delete-button").style.display = "none";
            fetchVideos();
        }

        function toggleFullScreen() {
            const videoElement = document.getElementById("main-video");
            if (!document.fullscreenElement) {
                videoElement.requestFullscreen().catch(err => {
                    console.error("Error attempting to enable full-screen mode:" + err.message);
                });
            } else {
                document.exitFullscreen();
            }
        }

        
        // Automatically load video list on page load
        window.onload = fetchVideos;
    </script>
</body>
</html>
`
