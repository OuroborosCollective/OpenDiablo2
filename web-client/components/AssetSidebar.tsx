"use client";

import React from "react";
import { Info, Package, Image as ImageIcon, Music, FileJson } from "lucide-react";

export interface AssetMetadata {
  id: string;
  name: string;
  type: "image" | "audio" | "json" | "data";
  size: string;
  path: string;
}

interface AssetSidebarProps {
  assets?: AssetMetadata[];
}

const AssetSidebar: React.FC<AssetSidebarProps> = ({ assets = [] }) => {
  const getIcon = (type: string) => {
    switch (type) {
      case "image": return <ImageIcon size={16} />;
      case "audio": return <Music size={16} />;
      case "json": return <FileJson size={16} />;
      default: return <Package size={16} />;
    }
  };

  return (
    <aside className="w-80 bg-neutral-900 border-l border-white/10 flex flex-col h-full text-white">
      <div className="p-4 border-b border-white/10 flex items-center gap-2">
        <Info size={20} className="text-blue-400" />
        <h2 className="font-semibold text-lg">Asset Metadata</h2>
      </div>
      <div className="flex-grow overflow-y-auto">
        <div className="p-2 space-y-1">
          {assets.length === 0 && (
            <div className="p-8 text-center text-neutral-500 text-sm">
              No assets loaded
            </div>
          )}
          {assets.map((asset) => (
            <div
              key={asset.id}
              className="p-3 rounded-md hover:bg-white/5 transition-colors cursor-pointer group border border-transparent hover:border-white/10"
            >
              <div className="flex items-center gap-3">
                <div className="p-2 bg-white/5 rounded-lg text-neutral-400 group-hover:text-blue-400">
                  {getIcon(asset.type)}
                </div>
                <div className="flex-grow min-w-0">
                  <div className="font-medium text-sm truncate">{asset.name}</div>
                  <div className="text-xs text-neutral-500 flex justify-between">
                    <span>{asset.type.toUpperCase()}</span>
                    <span>{asset.size}</span>
                  </div>
                </div>
              </div>
              <div className="mt-2 text-[10px] text-neutral-600 truncate font-mono">
                {asset.path}
              </div>
            </div>
          ))}
        </div>
      </div>
      <div className="p-4 bg-black/20 border-t border-white/10">
        <div className="text-[10px] text-neutral-500 uppercase tracking-widest font-bold mb-2">System Status</div>
        <div className="flex items-center gap-2">
          <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
          <span className="text-xs text-neutral-300">Engine Responsive</span>
        </div>
      </div>
    </aside>
  );
};

export default AssetSidebar;
