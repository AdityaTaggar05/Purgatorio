import { useEffect, useRef } from 'react';
import Phaser from 'phaser';
import { TerraceScene } from './scenes/TerraceScene';
import type { BaseLayout } from '../../../types/building';

let latestLayout: BaseLayout | null = null;

export function setPhaserLayout(layout: BaseLayout | null) {
  latestLayout = layout;
}

export { latestLayout };

interface GameCanvasProps {
  layout: BaseLayout | null;
}

export default function GameCanvas({ layout }: GameCanvasProps) {
  const gameRef = useRef<HTMLDivElement>(null);
  const phaserInstance = useRef<Phaser.Game | null>(null);

  useEffect(() => {
    latestLayout = layout;
  }, [layout]);

  useEffect(() => {
    if (!gameRef.current) return;

    const config: Phaser.Types.Core.GameConfig = {
      type: Phaser.AUTO,
      parent: gameRef.current,
      width: '100%',
      height: '100%',
      scene: [TerraceScene],
      transparent: true,
      pixelArt: false,
      antialias: true,
      antialiasGL: true,
      scale: {
        mode: Phaser.Scale.RESIZE,
        autoCenter: Phaser.Scale.CENTER_BOTH,
      },
    };

    phaserInstance.current = new Phaser.Game(config);

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
