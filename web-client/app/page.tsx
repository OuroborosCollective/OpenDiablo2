"use client";

import React, { useState } from "react";
import BabylonScene from "@/components/BabylonScene";
import AssetSidebar, { AssetMetadata } from "@/components/AssetSidebar";

export default function Home() {
  const [resonance, setResonance] = useState(0);
  const [cycle, setCycle] = useState(0);

  const [assets, setAssets] = useState<AssetMetadata[]>([]);

  const handleAxiomaticUpdate = (res: number, cyc: number) => {
    setResonance(res);
    setCycle(cyc);
  };

  const handleAssetMetadataUpdate = (newAssets: AssetMetadata[]) => {
    setAssets(newAssets);
  };

  return (
    <main className="flex w-screen h-screen overflow-hidden bg-black text-white">
      {/* Game View */}
      <div className="relative flex-grow h-full overflow-hidden">
        <div className="absolute inset-0">
          <BabylonScene
            onAxiomaticUpdate={handleAxiomaticUpdate}
            onAssetMetadataUpdate={handleAssetMetadataUpdate}
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
      <AssetSidebar assets={assets} resonance={resonance} cycle={cycle} />
    </main>
  );
}
