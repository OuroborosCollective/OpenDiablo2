"use client";

import React, { useState, useEffect } from "react";
import BabylonScene from "@/components/BabylonScene";
import AssetSidebar, { AssetMetadata } from "@/components/AssetSidebar";
import { ChevronRight, ChevronLeft } from "lucide-react";

export default function Home() {
  const [resonance, setResonance] = useState(0);
  const [cycle, setCycle] = useState(0);
  const [isSidebarOpen, setSidebarOpen] = useState(true);
  const [isMobile, setIsMobile] = useState(false);

  const [assets, setAssets] = useState<AssetMetadata[]>([]);

  useEffect(() => {
    const checkMobile = () => {
      const mobile = window.innerWidth < 768;
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

  return (
    <main className="flex w-screen h-screen overflow-hidden bg-black text-white relative">
      {/* Game View */}
      <div className="relative flex-grow h-full overflow-hidden">
        <div className="absolute inset-0">
          <BabylonScene
            onAxiomaticUpdate={handleAxiomaticUpdate}
            onAssetMetadataUpdate={handleAssetMetadataUpdate}
          />
        </div>

        {/* Sidebar Toggle Button (Mobile Friendly) */}
        <button
          onClick={() => setSidebarOpen(!isSidebarOpen)}
          className={`absolute top-1/2 right-0 transform -translate-y-1/2 z-30 p-2 bg-neutral-900/80 border-l border-y border-white/20 rounded-l-md hover:bg-neutral-800 transition-all ${isSidebarOpen ? "translate-x-0" : "translate-x-0"}`}
        >
          {isSidebarOpen ? <ChevronRight size={20} /> : <ChevronLeft size={20} />}
        </button>

        <div className="absolute top-4 left-4 z-10 p-4 bg-black/60 text-white rounded-lg border border-white/20 backdrop-blur-md">
          <h1 className="text-xl font-bold tracking-tight bg-gradient-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
            OpenDiablo2 Mobile
          </h1>
          <p className="text-[10px] uppercase tracking-widest opacity-70 mt-1 font-semibold">
            Axiomatic Logic Integration
          </p>
        </div>
      </div>

      {/* Asset Sidebar */}
      <div className={`${isSidebarOpen ? "w-80" : "w-0"} transition-all duration-300 overflow-hidden h-full border-l border-white/10`}>
        <AssetSidebar assets={assets} resonance={resonance} cycle={cycle} />
      </div>

      {isMobile && !isSidebarOpen && (
        <div className="absolute bottom-20 left-4 z-20 flex flex-col gap-2">
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
