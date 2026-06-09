"use client";

import React, { useState, useEffect } from "react";
import BabylonScene from "@/components/BabylonScene";
import AssetSidebar, { AssetMetadata } from "@/components/AssetSidebar";
import VirtualJoystick from "@/components/VirtualJoystick";
import MobileHotbar from "@/components/MobileHotbar";
import MobileHUD from "@/components/MobileHUD";
import { getMobileConfig, MOBILE_CONFIG } from "@/utils/mobile";
import { ChevronRight, ChevronLeft } from "lucide-react";

export default function Home() {
  const [resonance, setResonance] = useState(0);
  const [cycle, setCycle] = useState(0);
  const [isSidebarOpen, setSidebarOpen] = useState(true);
  const [isMobile, setIsMobile] = useState(false);
  const [joystickInput, setJoystickInput] = useState({ x: 0, y: 0 });
  const [direction, setDirection] = useState("idle");

  const [assets, setAssets] = useState<AssetMetadata[]>([]);
  
  // Player stats for HUD
  const [playerStats, setPlayerStats] = useState({
    health: 75,
    maxHealth: 100,
    mana: 50,
    maxMana: 100,
    level: 1,
    experience: 35,
  });

  useEffect(() => {
    const config = getMobileConfig();
    setIsMobile(config.isMobile);
    if (config.isMobile) setSidebarOpen(false);
  }, []);

  const handleAxiomaticUpdate = (res: number, cyc: number) => {
    setResonance(res);
    setCycle(cyc);
  };

  const handleAssetMetadataUpdate = (newAssets: AssetMetadata[]) => {
    setAssets(newAssets);
  };

  const handleJoystickMove = (x: number, y: number) => {
    setJoystickInput({ x, y });
  };

  const handleDirectionChange = (dir: string) => {
    setDirection(dir);
    // Could send direction to game server for animation
  };

  const handleSkillUse = (slot: number) => {
    console.log("Skill used:", slot);
    // Send to Axiomatic engine
  };

  return (
    <main className="flex w-screen h-screen overflow-hidden bg-black text-white relative">
      {/* Game View */}
      <div className="relative flex-grow h-full overflow-hidden">
        <div className="absolute inset-0">
          <BabylonScene
            onAxiomaticUpdate={handleAxiomaticUpdate}
            onAssetMetadataUpdate={handleAssetMetadataUpdate}
            joystickInput={joystickInput}
          />
        </div>

        {/* Sidebar Toggle Button */}
        <button
          onClick={() => setSidebarOpen(!isSidebarOpen)}
          className={`absolute top-1/2 right-0 transform -translate-y-1/2 z-30 p-2 bg-neutral-900/80 border-l border-y border-white/20 rounded-l-md hover:bg-neutral-800 transition-all`}
        >
          {isSidebarOpen ? <ChevronRight size={20} /> : <ChevronLeft size={20} />}
        </button>

        {/* Header Badge */}
        <div className="absolute top-4 left-4 z-10 p-4 bg-black/60 text-white rounded-lg border border-white/20 backdrop-blur-md">
          <h1 className="text-xl font-bold tracking-tight bg-gradient-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
            Ouroboros Collective
          </h1>
          <p className="text-[10px] uppercase tracking-widest opacity-70 mt-1 font-semibold">
            ARE-Diablo2-BaalAal
          </p>
        </div>

        {/* Mobile HUD */}
        {isMobile && (
          <MobileHUD
            health={playerStats.health}
            maxHealth={playerStats.maxHealth}
            mana={playerStats.mana}
            maxMana={playerStats.maxMana}
            level={playerStats.level}
            experience={playerStats.experience}
            showDebug={false}
          />
        )}

        {/* Direction Indicator (Mobile) */}
        {isMobile && direction !== "idle" && (
          <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-20 text-4xl opacity-50 pointer-events-none">
            {direction === "up" && "⬆️"}
            {direction === "down" && "⬇️"}
            {direction === "left" && "⬅️"}
            {direction === "right" && "➡️"}
            {direction === "up-left" && "↖️"}
            {direction === "up-right" && "↗️"}
            {direction === "down-left" && "↙️"}
            {direction === "down-right" && "↘️"}
          </div>
        )}

        {/* Mobile Touch Controls */}
        {isMobile && (
          <>
            <VirtualJoystick
              size={MOBILE_CONFIG.joystickSize}
              onMove={handleJoystickMove}
              onDirectionChange={handleDirectionChange}
              position="left"
              sensitivity={1.0}
            />
            <MobileHotbar
              slots={MOBILE_CONFIG.hotbarSlots}
              onSkillUse={handleSkillUse}
            />
          </>
        )}
      </div>

      {/* Asset Sidebar */}
      <div className={`${isSidebarOpen ? "w-80" : "w-0"} transition-all duration-300 overflow-hidden h-full border-l border-white/10`}>
        <AssetSidebar assets={assets} resonance={resonance} cycle={cycle} />
      </div>

      {/* Desktop Quick Actions */}
      {!isMobile && (
        <div className="absolute bottom-4 left-4 z-20 flex gap-2">
          <div className="p-3 bg-blue-500/20 rounded-full border border-blue-400/50 backdrop-blur-md">
            <div className="w-8 h-8 flex items-center justify-center text-blue-400 font-bold">L</div>
          </div>
          <div className="p-3 bg-red-500/20 rounded-full border border-red-400/50 backdrop-blur-md">
            <div className="w-8 h-8 flex items-center justify-center text-red-400 font-bold">R</div>
          </div>
        </div>
      )}
    </main>
  );
}
