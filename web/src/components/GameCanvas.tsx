"use client";

import React, { useCallback, useRef, useEffect, useState } from "react";
import { Engine, Scene } from "react-babylonjs";
import { Vector3, Color3 } from "@babylonjs/core/Maths/math";
import { Mesh } from "@babylonjs/core/Meshes/mesh";

const GameCanvas: React.FC = () => {
  const playerRef = useRef<Mesh | null>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const [status, setStatus] = useState("Disconnected");
  const [clickPos, setClickPos] = useState<{ x: number, y: number } | null>(null);

  const onSceneMount = useCallback((sceneEventArgs: any) => {
    const { scene } = sceneEventArgs;

    scene.clearColor = new Color3(0.05, 0.05, 0.1).toColor4();

    scene.onPointerDown = (evt: any) => {
      if (evt.button === 0) { // Left click or touch
        const pickResult = scene.pick(scene.pointerX, scene.pointerY);
        if (pickResult?.hit && pickResult.pickedPoint && playerRef.current) {
          const destX = pickResult.pickedPoint.x;
          const destY = pickResult.pickedPoint.z;

          setClickPos({ x: scene.pointerX, y: scene.pointerY });
          setTimeout(() => setClickPos(null), 300);

          if (socketRef.current?.readyState === WebSocket.OPEN) {
            const movePacket = {
              packetType: 3, // MovePlayer
              packetData: btoa(JSON.stringify({
                playerID: "web-player",
                destX,
                destY,
                startX: playerRef.current.position.x,
                startY: playerRef.current.position.z
              }))
            };
            socketRef.current.send(JSON.stringify(movePacket));
          } else {
            playerRef.current.position.x = destX;
            playerRef.current.position.z = destY;
          }
        }
      }
    };
  }, []);

  useEffect(() => {
    const ws = new WebSocket("ws://" + window.location.hostname + ":6670/ws");
    socketRef.current = ws;

    ws.onopen = () => {
      setStatus("Connected");
      const connRequest = {
        packetType: 4, // PlayerConnectionRequest
        packetData: btoa(JSON.stringify({
          id: "web-player-" + Math.random().toString(36).substring(2, 9),
          heroName: "WebHero",
          heroType: 0,
          playerState: { x: 0, y: 0 }
        }))
      };
      ws.send(JSON.stringify(connRequest));
    };

    ws.onmessage = (event) => {
      const packet = JSON.parse(event.data);

      if (packet.packetType === 3) { // MovePlayer
        const moveData = JSON.parse(atob(packet.packetData));
        if (playerRef.current) {
          playerRef.current.position.x = moveData.destX;
          playerRef.current.position.z = moveData.destY;
        }
      } else if (packet.packetType === 13) { // AxiomaticStatus (enum: 13)
        const axData = JSON.parse(atob(packet.packetData));
        window.dispatchEvent(new CustomEvent('axiomatic-status', { detail: axData }));
      }
    };

    ws.onclose = () => setStatus("Disconnected");
    ws.onerror = () => setStatus("Error");

    return () => ws.close();
  }, []);

  return (
    <div style={{ flex: 1, display: "flex", height: "100%", width: "100%", position: "relative" }}>
      <Engine antialias adaptToDeviceRatio canvasId="babylonJS">
        <Scene onSceneMount={onSceneMount}>
          <freeCamera
            name="camera1"
            position={new Vector3(0, 10, -15)}
            setTarget={[Vector3.Zero()]}
          />
          <hemisphericLight
            name="light1"
            intensity={0.8}
            direction={new Vector3(0, 1, 1)}
          />

          <box
            name="player"
            size={1}
            position={new Vector3(0, 0.5, 0)}
            ref={playerRef}
          >
            <standardMaterial name="playerMat" diffuseColor={Color3.Red()} />
          </box>

          <ground name="ground1" width={40} height={40} subdivisions={2}>
            <standardMaterial name="groundMat" wireframe={true} />
          </ground>
        </Scene>
      </Engine>

      {clickPos && (
        <div
          className="absolute w-8 h-8 border-2 border-orange-500 rounded-full animate-ping pointer-events-none"
          style={{ left: clickPos.x - 16, top: clickPos.y - 16 }}
        />
      )}

      <div className="absolute bottom-4 left-4 flex flex-col gap-1 pointer-events-none">
        <div className="bg-black/70 text-[10px] text-white px-2 py-1 rounded border border-white/10 font-mono text-center">
          WS: {status}
        </div>
        <div className="bg-orange-950/70 text-[10px] text-orange-200 px-2 py-1 rounded border border-orange-500/20 font-mono">
          MOBILE_READY: ACTIVE
        </div>
      </div>
    </div>
  );
};

export default GameCanvas;
