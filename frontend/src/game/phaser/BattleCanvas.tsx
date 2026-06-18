import { useEffect, useRef } from "react";
import Phaser from "phaser";
import { BattleScene } from "./scenes/BattleScene";
import { battleEvents } from "./battleEvents";
import type { BaseLayout } from "../../types/building";
import type { TroopStackData } from "./entities/TroopStackSprite";

export let latestDefenderLayout: BaseLayout | null = null;
export let latestBuildingHpByKey: Record<string, number> = {};
export let latestTroopStacks: TroopStackData[] = [];
export let latestDeploymentZone: { x: number; y: number }[] | null = null;
export let latestSelectedCells: { x: number; y: number }[] = [];
export let deployClickEnabled = false;
export let latestSelectedTroop: string | null = null;

export let troopCatalogIds: string[] = [];
export let staticPreviewMode = false;

interface BattleCanvasProps {
  layout: BaseLayout | null;
  buildingHpByKey?: Record<string, number>;
  stacks: TroopStackData[];
  deploymentZone?: { x: number; y: number }[] | null;
  selectedCells?: { x: number; y: number }[];
  troopCatalogIds: string[];
  interactiveDeployment?: boolean;
  staticPreview?: boolean;
  selectedTroop?: string | null;
  onCellClick?: (x: number, y: number) => void;
}

export default function BattleCanvas({
  layout,
  buildingHpByKey = {},
  stacks,
  deploymentZone = null,
  selectedCells = [],
  troopCatalogIds: catalogIds,
  interactiveDeployment = false,
  staticPreview = false,
  selectedTroop = null,
  onCellClick,
}: BattleCanvasProps) {
  const gameRef = useRef<HTMLDivElement>(null);
  const phaserInstance = useRef<Phaser.Game | null>(null);

  // Set the boot-time-only values before the Phaser.Game (and therefore
  // BattleScene.preload) is ever constructed.
  troopCatalogIds = catalogIds;
  staticPreviewMode = staticPreview;

  useEffect(() => {
    latestDefenderLayout = layout;
    latestBuildingHpByKey = buildingHpByKey;
    latestTroopStacks = stacks;
    latestDeploymentZone = deploymentZone;
    latestSelectedCells = selectedCells;
    deployClickEnabled = interactiveDeployment;
    latestSelectedTroop = selectedTroop;
  }, [layout, buildingHpByKey, stacks, deploymentZone, selectedCells, interactiveDeployment, selectedTroop]);

  useEffect(() => {
    battleEvents.onCellClick = onCellClick ?? null;
    return () => {
      battleEvents.onCellClick = null;
    };
  }, [onCellClick]);

  useEffect(() => {
    if (!gameRef.current) return;

    const config: Phaser.Types.Core.GameConfig = {
      type: Phaser.AUTO,
      parent: gameRef.current,
      width: "100%",
      height: "100%",
      scene: [BattleScene],
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

  return <div ref={gameRef} className="w-full h-full block" />;
}
