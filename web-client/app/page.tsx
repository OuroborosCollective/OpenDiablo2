"use client";

import React, { useState } from "react";
import BabylonScene from "@/components/BabylonScene";
import AssetSidebar, { AssetMetadata } from "@/components/AssetSidebar";
import { Menu, X } from "lucide-react";

export default function Home() {
  const [resonance, setResonance] = useState(0);
  const [cycle, setCycle] = useState(0);
  const [assets, setAssets] = useState<AssetMetadata[]>([]);
  const [sidebarOpen, setSidebarOpen] = useState(true);

  const handleAxiomaticUpdate = (res: number, cyc: number) => {
    setResonance(res);
    setCycle(cyc);
  };

  return (
    <main className="flex w-screen h-screen overflow-hidden bg-black text-white relative">
      {/* Sidebar Toggle Button for Mobile */}
      <button
        onClick={() => setSidebarOpen(!sidebarOpen)}
        className="fixed top-4 right-4 z-50 p-2 bg-black/60 border border-white/20 rounded-md md:hidden"
      >
        {sidebarOpen ? <X size={20} /> : <Menu size={20} />}
      </button>

      {/* Game View */}
      <div className="relative flex-grow h-full overflow-hidden">
        <div className="absolute inset-0">
          <BabylonScene
            onAxiomaticUpdate={handleAxiomaticUpdate}
            onAssetListUpdate={setAssets}
          />
        </div>
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
      <div className={`${sidebarOpen ? 'translate-x-0' : 'translate-x-full'} fixed right-0 top-0 h-full z-40 transition-transform duration-300 ease-in-out md:relative md:translate-x-0`}>
        <AssetSidebar assets={assets} resonance={resonance} cycle={cycle} />
      </div>
    </main>
  );
}
