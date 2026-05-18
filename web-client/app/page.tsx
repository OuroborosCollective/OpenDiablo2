"use client";

import React, { useState } from "react";
import BabylonScene from "@/components/BabylonScene";
import AssetSidebar, { AssetMetadata } from "@/components/AssetSidebar";

export default function Home() {
  const [assets] = useState<AssetMetadata[]>([
    { id: "1", name: "d2data.mpq", type: "data", size: "245 MB", path: "data/d2data.mpq" },
    { id: "2", name: "d2exp.mpq", type: "data", size: "180 MB", path: "data/d2exp.mpq" },
    { id: "3", name: "d2sfx.mpq", type: "audio", size: "520 MB", path: "data/d2sfx.mpq" },
    { id: "4", name: "d2music.mpq", type: "audio", size: "380 MB", path: "data/d2music.mpq" },
    { id: "5", name: "d2video.mpq", type: "data", size: "900 MB", path: "data/d2video.mpq" },
    { id: "6", name: "d2char.mpq", type: "image", size: "110 MB", path: "data/d2char.mpq" },
  ]);

  return (
    <main className="flex w-screen h-screen overflow-hidden bg-black text-white">
      {/* Game View */}
      <div className="relative flex-grow h-full overflow-hidden">
        <div className="absolute inset-0">
          <BabylonScene />
        </div>
        <div className="absolute top-4 left-4 z-10 p-4 bg-black/50 text-white rounded-lg border border-white/20 backdrop-blur-sm">
          <h1 className="text-xl font-bold tracking-tight">OpenDiablo2 Mobile</h1>
          <p className="text-sm opacity-70">Migration PoC: Next.js + Babylon.js + Axiomatic Go</p>
        </div>
      </div>

      {/* Asset Sidebar */}
      <AssetSidebar assets={assets} />
    </main>
  );
}
