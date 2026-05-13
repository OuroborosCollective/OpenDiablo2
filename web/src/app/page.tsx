import GameCanvas from "../components/GameCanvas";

export default function Home() {
  return (
    <main className="flex flex-1 w-full h-full flex-col bg-black overflow-hidden">
      <GameCanvas />
    </main>
  );
}