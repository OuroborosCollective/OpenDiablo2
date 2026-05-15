"use client";

import React, { useCallback, useRef } from "react";
import { Engine, Scene } from "react-babylonjs";
import { Vector3, Color3 } from "@babylonjs/core/Maths/math";
import { Mesh } from "@babylonjs/core/Meshes/mesh";

const GameCanvas: React.FC = () => {
  const playerRef = useRef<Mesh | null>(null);

  const onSceneMount = useCallback((sceneEventArgs: any) => {
    const { scene } = sceneEventArgs;

    // Improved scene background
    scene.clearColor = new Color3(0.05, 0.05, 0.1).toColor4();

    // Interaction handling
    scene.onPointerDown = (evt: any) => {
      if (evt.button === 0) { // Left click or touch
        const pickResult = scene.pick(scene.pointerX, scene.pointerY);
        if (pickResult?.hit && pickResult.pickedPoint && playerRef.current) {
          // Visual movement for PoC
          playerRef.current.position.x = pickResult.pickedPoint.x;
          playerRef.current.position.z = pickResult.pickedPoint.z;
        }
      }
    };
  }, []);

  return (
    <div style={{ flex: 1, display: "flex", height: "100%", width: "100%" }}>
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

          {/* Player Representative */}
          <box
            name="player"
            size={1}
            position={new Vector3(0, 0.5, 0)}
            ref={playerRef}
          >
            <standardMaterial name="playerMat" diffuseColor={Color3.Red()} />
          </box>

          {/* World Reference Grid */}
          <ground name="ground1" width={40} height={40} subdivisions={2}>
            <standardMaterial name="groundMat" wireframe={true} />
          </ground>
        </Scene>
      </Engine>
    </div>
  );
};

export default GameCanvas;
