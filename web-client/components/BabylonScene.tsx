"use client";

import React, { useEffect, useRef, useState } from "react";
import { Engine, Scene, ArcRotateCamera, Vector3, HemisphericLight, MeshBuilder, Mesh, StandardMaterial, Color3 } from "@babylonjs/core";
import { PacketType } from "@/utils/packetTypes";

import { AssetMetadata } from "./AssetSidebar";

interface BabylonSceneProps {
  onAxiomaticUpdate?: (resonance: number, cycle: number) => void;
  onAssetListUpdate?: (assets: AssetMetadata[]) => void;
}

const BabylonScene: React.FC<BabylonSceneProps> = ({ onAxiomaticUpdate, onAssetListUpdate }) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const [status, setStatus] = useState("Connecting...");
  const playersMap = useRef<Map<string, Mesh>>(new Map());
  const targetMarkerRef = useRef<Mesh | null>(null);

  useEffect(() => {
    if (!canvasRef.current) return;

    const engine = new Engine(canvasRef.current, true);
    const scene = new Scene(engine);
    scene.clearColor = new Color3(0.1, 0.1, 0.1).toColor4();

    const camera = new ArcRotateCamera("camera", -Math.PI / 2, Math.PI / 3, 20, new Vector3(0, 0, 0), scene);
    camera.attachControl(canvasRef.current, true);

    new HemisphericLight("light", new Vector3(1, 1, 0), scene);

    // Grid for world reference
    const grid = MeshBuilder.CreateGround("grid", { width: 100, height: 100 }, scene);
    const gridMat = new StandardMaterial("gridMat", scene);
    gridMat.wireframe = true;
    grid.material = gridMat;

    const createPlayerMesh = (id: string, x: number, z: number, isLocal: boolean) => {
      const capsule = MeshBuilder.CreateCapsule(id, { height: 2, radius: 0.5 }, scene);
      capsule.position.set(x, 1, z);
      const mat = new StandardMaterial(id + "Mat", scene);
      mat.diffuseColor = isLocal ? new Color3(0, 1, 0) : new Color3(1, 0, 0);
      capsule.material = mat;
      playersMap.current.set(id, capsule);
      return capsule;
    };

  // Target Marker Mesh
  const marker = MeshBuilder.CreateTorus("targetMarker", { diameter: 1, thickness: 0.1 }, scene);
  marker.isVisible = false;
  const markerMat = new StandardMaterial("markerMat", scene);
  markerMat.emissiveColor = new Color3(0.4, 0.6, 1);
  marker.material = markerMat;
  targetMarkerRef.current = marker;

    // WebSocket Setup
    const ws = new WebSocket("ws://" + window.location.hostname + ":6670/ws");
    socketRef.current = ws;

    ws.onopen = () => {
      setStatus("Connected");
      // Send Connection Request
      const connRequest = {
        packetType: PacketType.PlayerConnectionRequest,
        packetData: btoa(JSON.stringify({
          id: "web-player-" + Math.random().toString(36).substr(2, 9),
          heroName: "WebHero",
          heroType: 0,
          playerState: { x: 0, y: 0 }
        }))
      };
      ws.send(JSON.stringify(connRequest));
    };

    ws.onmessage = (event) => {
      const packet = JSON.parse(event.data);
      if (packet.packetType === PacketType.GenerateMap) {
        // Clear existing players on map change
        playersMap.current.forEach(mesh => mesh.dispose());
        playersMap.current.clear();
      } else if (packet.packetType === PacketType.AddPlayer) {
        const playerData = JSON.parse(atob(packet.packetData));
        if (!playersMap.current.has(playerData.id)) {
          createPlayerMesh(playerData.id, playerData.x / 5, playerData.y / 5, false);
        }
      } else if (packet.packetType === PacketType.UpdateServerInfo) {
        const serverInfo = JSON.parse(atob(packet.packetData));
        // Local player creation handled by server assigning ID
        if (!playersMap.current.has(serverInfo.playerID)) {
          createPlayerMesh(serverInfo.playerID, 0, 0, true);
        }
      } else if (packet.packetType === PacketType.MovePlayer) {
        const moveData = JSON.parse(atob(packet.packetData));
        const mesh = playersMap.current.get(moveData.playerID);
        if (mesh) {
          mesh.position.x = moveData.destX;
          mesh.position.z = moveData.destY;
        }
      } else if (packet.packetType === PacketType.PlayerDisconnectionNotification) {
        const disconnectData = JSON.parse(atob(packet.packetData));
        const mesh = playersMap.current.get(disconnectData.id);
        if (mesh) {
          mesh.dispose();
          playersMap.current.delete(disconnectData.id);
        }
      } else if (packet.packetType === PacketType.AssetList) {
        const assetListData = JSON.parse(atob(packet.packetData));
        if (onAssetListUpdate) {
          onAssetListUpdate(assetListData.assets);
        }
      } else if (packet.packetType === PacketType.AxiomaticStatus) {
        const statusData = JSON.parse(atob(packet.packetData));
        if (onAxiomaticUpdate) {
          onAxiomaticUpdate(statusData.resonance, statusData.cycle);
        }
      }
    };

    ws.onerror = (err) => {
      setStatus("Error: " + err);
    };

    ws.onclose = () => {
      setStatus("Disconnected");
    };

    // Input Handling
    scene.onPointerDown = (evt) => {
      if (evt.button === 0) { // Left click / Touch
        const pickResult = scene.pick(scene.pointerX, scene.pointerY);
        if (pickResult?.hit && pickResult.pickedPoint) {
          const destX = pickResult.pickedPoint.x;
          const destY = pickResult.pickedPoint.z;

          // Visual Feedback
          if (targetMarkerRef.current) {
            targetMarkerRef.current.position.set(destX, 0.1, destY);
            targetMarkerRef.current.isVisible = true;
          }

          if (ws.readyState === WebSocket.OPEN) {
            const movePacket = {
              packetType: PacketType.MovePlayer,
              packetData: btoa(JSON.stringify({
                destX,
                destY
              }))
            };
            ws.send(JSON.stringify(movePacket));
          }
        }
      }
    };

    engine.runRenderLoop(() => {
      if (targetMarkerRef.current?.isVisible) {
        targetMarkerRef.current.rotation.y += 0.05;
        targetMarkerRef.current.scaling.setAll(1 + Math.sin(Date.now() * 0.01) * 0.1);
      }
      scene.render();
    });

    const handleResize = () => {
      engine.resize();
    };

    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
      ws.close();
      engine.dispose();
    };
  }, [onAxiomaticUpdate, onAssetListUpdate]);

  return (
    <div className="relative w-full h-full">
      <canvas ref={canvasRef} className="w-full h-full touch-none" />
      <div className="absolute bottom-4 right-4 bg-black/80 text-xs text-white p-2 rounded border border-white/20">
        Status: {status}
      </div>
    </div>
  );
};

export default BabylonScene;
