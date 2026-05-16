import BabylonScene from "@/components/BabylonScene";
import AssetSidebar from "@/components/AssetSidebar";

export default function Home() {
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
      <AssetSidebar />
    </main>
  );
}
