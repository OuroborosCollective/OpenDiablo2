"use client";

import React from "react";
import { Engine, Scene } from "react-babylonjs";
import { Vector3 } from "@babylonjs/core/Maths/math";

const GameCanvas: React.FC = () => {
  return (
    <div style={{ flex: 1, display: "flex", height: "100%", width: "100%" }}>
      <Engine antialias adaptToDeviceRatio canvasId="babylonJS">
        <Scene>
          <freeCamera
            name="camera1"
            position={new Vector3(0, 5, -10)}
            setTarget={[Vector3.Zero()]}
          />
          <hemisphericLight
            name="light1"
            intensity={0.7}
            direction={new Vector3(0, 1, 0)}
          />
          <box
            name="box1"
            size={2}
            position={new Vector3(0, 1, 0)}
          />
          {/* Ground to see perspective better */}
          <ground name="ground1" width={6} height={6} subdivisions={2} />
        </Scene>
      </Engine>
    </div>
  );
};

export default GameCanvas;