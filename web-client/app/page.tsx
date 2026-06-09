"use client";

import React, { useState, useEffect } from "react";
import BabylonScene from "@/components/BabylonScene";
import AssetSidebar, { AssetMetadata } from "@/components/AssetSidebar";
import VirtualJoystick from "@/components/VirtualJoystick";
import MobileHotbar from "@/components/MobileHotbar";
import { ChevronRight, ChevronLeft } from "lucide-react";

export default function Home() {
  const [resonance, setResonance] = useState(0);
  const [cycle, setCycle] = useState(0);
  const [isSidebarOpen, setSidebarOpen] = useState(true);
  const [isMobile, setIsMobile] = useState(false);
  const [joystickInput, setJoystickInput] = useState({ x: 0, y: 0 });

  const [assets, setAssets] = useState<AssetMetadata[]>([]);

  useEffect(() => {
    const checkMobile = () => {
      const mobile = window.innerWidth < 768 || /Android|iPhone|iPad|iPod/i.test(navigator.userAgent);
      setIsMobile(mobile);
      if (mobile) setSidebarOpen(false);
    };

    checkMobile();
    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
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
    // Send to Babylon scene for character movement
    console.log("Joystick input:", x, y);
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

        {/* Sidebar Toggle Button (Mobile Friendly) */}
        <button
          onClick={() => setSidebarOpen(!isSidebarOpen)}
          className={`absolute top-1/2 right-0 transform -translate-y-1/2 z-30 p-2 bg-neutral-900/80 border-l border-y border-white/20 rounded-l-md hover:bg-neutral-800 transition-all ${isSidebarOpen ? "translate-x-0" : "translate-x-0"}`}
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

        {/* Mobile Touch Controls */}
        {isMobile && (
          <>
            <VirtualJoystick
              size={130}
              onMove={handleJoystickMove}
              position="left"
            />
            <MobileHotbar
              slots={6}
              onSkillUse={handleSkillUse}
            />
          </>
        )}

        {/* Mobile Health/Mana Bars */}
        {isMobile && (
          <div className="absolute top-20 left-4 z-20 flex flex-col gap-1">
            <div className="w-32 h-3 bg-red-600/80 rounded-full overflow-hidden border border-red-400/50">
              <div className="w-3/4 h-full bg-gradient-to-r from-red-600 to-red-400" />
            </div>
            <div className="w-32 h-2 bg-blue-600/80 rounded-full overflow-hidden border border-blue-400/50">
              <div className="w-1/2 h-full bg-gradient-to-r from-blue-600 to-blue-400" />
            </div>
          </div>
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
