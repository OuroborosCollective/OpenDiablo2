"use client";

import React, { useEffect, useRef, useState } from "react";
import { Engine, Scene, ArcRotateCamera, Vector3, HemisphericLight, MeshBuilder, Mesh, StandardMaterial, Color3 } from "@babylonjs/core";
import { PacketType } from "@/utils/packetTypes";

import { AssetMetadata } from "@/components/AssetSidebar";

interface BabylonSceneProps {
  onAxiomaticUpdate?: (resonance: number, cycle: number) => void;
  onAssetMetadataUpdate?: (assets: AssetMetadata[]) => void;
  joystickInput?: { x: number; y: number };
}

const BabylonScene: React.FC<BabylonSceneProps> = ({ 
  onAxiomaticUpdate, 
  onAssetMetadataUpdate,
  joystickInput = { x: 0, y: 0 }
}) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const [status, setStatus] = useState("Connecting...");
  const playerRef = useRef<Mesh | null>(null);
  const lastJoystickRef = useRef<{ x: number; y: number }>({ x: 0, y: 0 });

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

    // Player Mesh
    const player = MeshBuilder.CreateCapsule("player", { height: 2, radius: 0.5 }, scene);
    player.position.y = 1;
    const playerMat = new StandardMaterial("playerMat", scene);
    playerMat.diffuseColor = new Color3(1, 0, 0);
    player.material = playerMat;
    playerRef.current = player;

    // Target Marker for mobile feedback
    const marker = MeshBuilder.CreateTorus("marker", { diameter: 1, thickness: 0.1 }, scene);
    marker.isVisible = false;
    const markerMat = new StandardMaterial("markerMat", scene);
    markerMat.diffuseColor = new Color3(0, 1, 0);
    markerMat.emissiveColor = new Color3(0, 0.5, 0);
    marker.material = markerMat;

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
      if (packet.packetType === PacketType.MovePlayer) {
        const moveData = JSON.parse(atob(packet.packetData));
        if (playerRef.current) {
          playerRef.current.position.x = moveData.destX;
          playerRef.current.position.z = moveData.destY;
        }
      } else if (packet.packetType === PacketType.AxiomaticStatus) {
        const statusData = JSON.parse(atob(packet.packetData));
        if (onAxiomaticUpdate) {
          onAxiomaticUpdate(statusData.resonance, statusData.cycle);
        }
      } else if (packet.packetType === PacketType.AssetMetadataList) {
        const assetData = JSON.parse(atob(packet.packetData));
        if (onAssetMetadataUpdate) {
          onAssetMetadataUpdate(assetData.assets || []);
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
      if (evt.button === 0) {
        const pickResult = scene.pick(scene.pointerX, scene.pointerY);
        if (pickResult?.hit && pickResult.pickedPoint) {
          const destX = pickResult.pickedPoint.x;
          const destY = pickResult.pickedPoint.z;

          marker.position.set(destX, 0.1, destY);
          marker.isVisible = true;
          marker.scaling.set(1, 1, 1);

          if (ws.readyState === WebSocket.OPEN) {
            const movePacket = {
              packetType: PacketType.MovePlayer,
              packetData: btoa(JSON.stringify({
                destX,
                destY,
                startX: player.position.x,
                startY: player.position.z
              }))
            };
            ws.send(JSON.stringify(movePacket));
          }
        }
      }
    };

    // Joystick movement handler
    const handleJoystickMovement = () => {
      if (!playerRef.current) return;
      
      const threshold = 0.1; // Minimum joystick movement threshold
      const speed = 0.15; // Movement speed
      
      if (Math.abs(joystickInput.x) > threshold || Math.abs(joystickInput.y) > threshold) {
        const newX = player.position.x + joystickInput.x * speed;
        const newZ = player.position.z + joystickInput.y * speed;
        
        player.position.x = newX;
        player.position.z = newZ;
        
        lastJoystickRef.current = { x: joystickInput.x, y: joystickInput.y };
        
        // Send position update to server
        if (ws.readyState === WebSocket.OPEN) {
          const movePacket = {
            packetType: PacketType.MovePlayer,
            packetData: btoa(JSON.stringify({
              destX: newX,
              destY: newZ,
              startX: player.position.x - joystickInput.x * speed,
              startY: player.position.z - joystickInput.y * speed
            }))
          };
          ws.send(JSON.stringify(movePacket));
        }
      }
    };

    engine.runRenderLoop(() => {
      if (marker.isVisible) {
        marker.rotation.y += 0.05;
        marker.scaling.x *= 0.98;
        marker.scaling.z *= 0.98;
        if (marker.scaling.x < 0.1) marker.isVisible = false;
      }
      
      // Handle joystick input
      handleJoystickMovement();
      
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
  }, [onAxiomaticUpdate, onAssetMetadataUpdate, joystickInput]);

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
