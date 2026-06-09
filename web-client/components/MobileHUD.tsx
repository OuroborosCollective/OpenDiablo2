"use client";

import React, { useState, useEffect } from "react";

interface HUDProps {
  health?: number;
  maxHealth?: number;
  mana?: number;
  maxMana?: number;
  level?: number;
  experience?: number;
  showDebug?: boolean;
}

export default function MobileHUD({
  health = 75,
  maxHealth = 100,
  mana = 50,
  maxMana = 100,
  level = 1,
  experience = 0,
  showDebug = false,
}: HUDProps) {
  const [isLandscape, setIsLandscape] = useState(true);
  const [displayRatio, setDisplayRatio] = useState("16:9");

  useEffect(() => {
    const checkOrientation = () => {
      const width = window.innerWidth;
      const height = window.innerHeight;
      setIsLandscape(width > height);
      
      // Determine display ratio
      const ratio = width / height;
      if (ratio >= 2.0) setDisplayRatio("21:9");
      else if (ratio >= 1.7) setDisplayRatio("16:9");
      else if (ratio >= 1.5) setDisplayRatio("18:9");
      else setDisplayRatio("16:9");
    };

    checkOrientation();
    window.addEventListener("resize", checkOrientation);
    return () => window.removeEventListener("resize", checkOrientation);
  }, []);

  const healthPercent = (health / maxHealth) * 100;
  const manaPercent = (mana / maxMana) * 100;

  // Scale based on screen size
  const scaleFactor = Math.min(window.innerWidth / 400, window.innerHeight / 800);
  const scaledSize = Math.max(0.8, Math.min(1.2, scaleFactor));

  return (
    <>
      {/* Top HUD - Health/Mana/Level */}
      <div 
        className={`absolute top-4 left-4 z-20 flex flex-col gap-2 ${isLandscape ? "" : "flex-row"}`}
        style={{ transform: `scale(${scaledSize})`, transformOrigin: "top left" }}
      >
        {/* Level Badge */}
        <div className="flex items-center gap-2">
          <div className="w-10 h-10 rounded-full bg-gradient-to-br from-yellow-500 to-orange-600 border-2 border-yellow-300 flex items-center justify-center font-bold text-white shadow-lg">
            {level}
          </div>
          {isLandscape && (
            <div className="text-xs text-yellow-400 font-semibold">LVL</div>
          )}
        </div>

        {/* Health Bar */}
        <div className="flex items-center gap-2">
          <div className="text-lg">❤️</div>
          <div className="relative w-32 h-4 bg-black/60 rounded-full border border-red-500/50 overflow-hidden">
            <div 
              className="h-full bg-gradient-to-r from-red-700 to-red-500 transition-all duration-300"
              style={{ width: `${healthPercent}%` }}
            />
            <div className="absolute inset-0 flex items-center justify-center">
              <span className="text-[10px] font-bold text-white drop-shadow-lg">
                {health}/{maxHealth}
              </span>
            </div>
          </div>
        </div>

        {/* Mana Bar */}
        <div className="flex items-center gap-2">
          <div className="text-lg">💙</div>
          <div className="relative w-32 h-3 bg-black/60 rounded-full border border-blue-500/50 overflow-hidden">
            <div 
              className="h-full bg-gradient-to-r from-blue-700 to-blue-500 transition-all duration-300"
              style={{ width: `${manaPercent}%` }}
            />
            <div className="absolute inset-0 flex items-center justify-center">
              <span className="text-[8px] font-bold text-white drop-shadow-lg">
                {mana}/{maxMana}
              </span>
            </div>
          </div>
        </div>

        {/* Experience Bar */}
        <div className="w-32 h-1.5 bg-black/60 rounded-full border border-yellow-500/30 overflow-hidden">
          <div 
            className="h-full bg-gradient-to-r from-yellow-600 to-yellow-400"
            style={{ width: `${experience}%` }}
          />
        </div>
      </div>

      {/* Debug Info (optional) */}
      {showDebug && (
        <div className="absolute top-4 right-4 z-20 bg-black/80 text-xs text-white p-2 rounded border border-white/20 font-mono">
          <div>Display: {displayRatio}</div>
          <div>Orientation: {isLandscape ? "Landscape" : "Portrait"}</div>
          <div>Scale: {(scaledSize * 100).toFixed(0)}%</div>
        </div>
      )}

      {/* Mini Map Placeholder */}
      <div className="absolute top-20 right-4 z-20 w-16 h-16 rounded-lg bg-black/60 border border-white/20 overflow-hidden">
        <div className="w-full h-full bg-gradient-to-br from-neutral-800 to-neutral-900 flex items-center justify-center text-[8px] text-white/50">
          MAP
        </div>
      </div>

      {/* Quest Indicator */}
      <div className="absolute top-20 left-4 z-20 flex items-center gap-2 p-2 bg-black/60 rounded-lg border border-white/20">
        <div className="w-6 h-6 rounded-full bg-purple-600/80 flex items-center justify-center text-xs">!</div>
        <span className="text-[10px] text-white/70">Quest Active</span>
      </div>
    </>
  );
}
