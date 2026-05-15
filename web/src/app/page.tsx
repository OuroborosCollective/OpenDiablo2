import GameCanvas from "../components/GameCanvas";
import AssetMetadata from "../components/AssetMetadata";

export default function Home() {
  return (
    <main className="flex flex-1 w-full h-full flex-col md:flex-row bg-black overflow-hidden">
      <div className="flex-1 relative overflow-hidden">
        <GameCanvas />
      </div>
      <AssetMetadata />
    </main>
  );
}
