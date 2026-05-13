# OpenDiablo2 Mobile Browser Migration Guide

This project has been updated with a boilerplate for a mobile browser-compatible version using **Next.js**, **Babylon.js**, and **Go**.

## Architecture Overview

### 1. Go Backend (Authoritative Server)
The Go engine now acts as a WebSocket host.
- **WebSocket Gateway**: Located in `d2networking/d2server/websocket_server.go`. It upgrades HTTP connections to WebSockets and bridges them to the `GameServer`.
- **WS Client Connection**: Located in `d2networking/d2server/d2wsclientconnection`. It implements the `ClientConnection` interface, allowing the existing game logic to communicate with web clients seamlessly.
- **Emergent Logic (ARE-Logik)**: A new package in `d2core/d2emergent` provides hooks for the Ouroboros Collective emergent NPC logic.

### 2. Next.js Frontend (Babylon.js Renderer)
A modern web client is located in the `/web-client` directory.
- **Babylon.js Scene**: `web-client/components/BabylonScene.tsx` handles the 3D rendering and input.
- **Touch Interaction**: Tap-to-move is implemented, sending `MovePlayer` packets to the Go server.
- **State Synchronization**: The client listens for server packets and updates the visual representation in real-time.

## How to Run

### Starting the Backend
In your Go initialization (e.g., in `main.go` or a dedicated server start), ensure you call:
```go
server.StartWebSocket(6670)
```

### Starting the Frontend
```bash
cd web-client
# Run: npm install
# Run: npm run dev
```
Visit `http://localhost:3000` on your device.

## Key Features for Mobile
- **Responsive Layout**: Powered by Tailwind CSS.
- **Touch-First Controls**: Normalized touch events for movement.
- **Authoritative Logic**: The client "looks" while the server "thinks," preventing desync and cheating.
