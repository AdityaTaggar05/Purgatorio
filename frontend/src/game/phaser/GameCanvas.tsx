import { useEffect, useRef } from 'react';
import Phaser from 'phaser';
import { TerraceScene } from './scenes/TerraceScene';

export default function GameCanvas() {
  const gameRef = useRef<HTMLDivElement>(null);
  const phaserInstance = useRef<Phaser.Game | null>(null);

  useEffect(() => {
    if (!gameRef.current) return;

    // Phaser Configuration
    const config: Phaser.Types.Core.GameConfig = {
      type: Phaser.AUTO,
      parent: gameRef.current, // Attaches the canvas to this React div
      width: '100%',
      height: '100%',
      scene: [TerraceScene],
      transparent: true,
      pixelArt: true, 
      scale: {
        mode: Phaser.Scale.RESIZE,
        autoCenter: Phaser.Scale.CENTER_BOTH,
      },
    };

    // Initialize the engine
    phaserInstance.current = new Phaser.Game(config);

    // Destroy the game instance on unmount or Fast Refresh
    return () => {
      if (phaserInstance.current) {
        phaserInstance.current.destroy(true);
        phaserInstance.current = null;
      }
    };
  }, []);

  return (
    <div ref={gameRef} className="w-full h-full absolute inset-0 block" />
  );
}
