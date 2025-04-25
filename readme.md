# File Server

## Overview
This project is a file server application designed to manage and serve video files. It includes features such as downloading videos, tracking watched episodes, and providing a user-friendly frontend interface.

## Features
- Download next 10 episodes or specific episodes.
- Track watched videos and resume playback.
- WebSocket support for real-time updates.
- Frontend interface for managing and playing videos.

## Setup Instructions

### Prerequisites
- Go (version 1.23.2 or higher)
- Node.js and npm
- A modern web browser

### Backend Setup
1. Navigate to the `server` directory:
   ```bash
   cd server
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the server:
   ```bash
   go run main.go
   ```

### Frontend Setup
1. Navigate to the `frontend` directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm start
   ```

## Usage
1. Access the application in your browser at `http://localhost:3000`.
2. Use the interface to browse, play, and manage videos.
3. Download episodes using the "Download Next 10 Videos" button or redownload specific episodes.

## API Endpoints
- `POST /api/videos/download-next-10`: Download the next 10 episodes.
- `POST /api/videos/download`: Download a specific episode.
- `DELETE /api/videos/{videoName}`: Delete a video.
- `GET /api/videos`: Fetch the list of available videos.
- `POST /api/videos/mark-watched`: Mark a video as watched.
- `GET /api/videos/watched`: Retrieve the list of watched videos.
- `GET /api/videos/{videoName}`: Fetch details of a specific video.
- `GET /api/videos/last-watched`: Get the last watched video details.

## Contributing
Contributions are welcome! Please submit a pull request or open an issue for any suggestions or bugs.

## License
This project is licensed under the MIT License.
