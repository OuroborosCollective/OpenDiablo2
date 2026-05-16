"use client";

import React, { useState } from "react";
import BabylonScene from "@/components/BabylonScene";
import AssetSidebar, { AssetMetadata } from "@/components/AssetSidebar";

export default function Home() {
  const [assets] = useState<AssetMetadata[]>([
    { id: "1", name: "hero_sprite.d2s", type: "image", size: "1.2 MB", path: "assets/sprites/hero_sprite.d2s" },
    { id: "2", name: "town_music.mp3", type: "audio", size: "4.5 MB", path: "assets/music/town_music.mp3" },
    { id: "3", name: "item_stats.json", type: "json", size: "45 KB", path: "data/item_stats.json" },
    { id: "4", name: "map_act1.mpq", type: "data", size: "120 MB", path: "data/map_act1.mpq" },
    { id: "5", name: "monster_anim.dcc", type: "image", size: "800 KB", path: "assets/anims/monster_anim.dcc" },
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
