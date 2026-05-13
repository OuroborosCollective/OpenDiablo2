import BabylonScene from "@/components/BabylonScene";

export default function Home() {
  return (
    <main className="relative w-screen h-screen overflow-hidden bg-black">
      <div className="absolute inset-0">
        <BabylonScene />
      </div>
      <div className="absolute top-4 left-4 z-10 p-4 bg-black/50 text-white rounded-lg border border-white/20">
        <h1 className="text-xl font-bold">OpenDiablo2 Mobile</h1>
        <p className="text-sm opacity-80">Migration PoC: Next.js + Babylon.js + Go</p>
      </div>
    </main>
  );
}
