import Phaser from 'phaser';
import { IsoMath } from './IsoMath';
import type { BaseLayout, PlacedBuilding } from '../../../types/building';
import { BuildingSprite } from '../entities/BuildingSprite';
import { phaserEvents } from '../events';

function computeBastionEdges(layout: BaseLayout): Map<string, string | null> {
  const bastions = layout.buildings.filter(b => b.building_id === 'bastion');
  const byKey = new Map<string, PlacedBuilding>();
  for (const b of bastions) {
    byKey.set(`${b.x},${b.y}`, b);
  }

  const result = new Map<string, string | null>();
  if (bastions.length === 0) return result;

  for (const b of bastions) {
    const key = `${b.x},${b.y}`;
    const dx = byKey.has(`${b.x + 1},${b.y}`) ? 1 : byKey.has(`${b.x - 1},${b.y}`) ? -1 : 0;
    const dy = byKey.has(`${b.x},${b.y + 1}`) ? 1 : byKey.has(`${b.x},${b.y - 1}`) ? -1 : 0;

    if (dx === 0 && dy === 0) {
      result.set(key, 'building_bastion-corner');
    } else if (dx !== 0) {
      result.set(key, 'building_bastion-edge-left');
    } else {
      result.set(key, 'building_bastion-edge-right');
    }
  }

  return result;
}

export class LayoutEngine {
  private scene: Phaser.Scene;
  private buildingsLayer!: Phaser.GameObjects.Layer;
  public activeBuildings: BuildingSprite[] = [];
  private interactive: boolean;

  private DEPTH_OFFSET = 100;

  constructor(scene: Phaser.Scene, interactive = true) {
    this.scene = scene;
    this.interactive = interactive;
    this.buildingsLayer = this.scene.add.layer();
    this.buildingsLayer.setDepth(this.DEPTH_OFFSET);

    if (this.interactive) {
      phaserEvents.getActiveBuildings = () => this.activeBuildings;
    }
  }

  public renderLayout(layout: BaseLayout) {
    this.clearBuildings();

    // Recreate the layer to ensure no stale object references
    this.buildingsLayer.destroy();
    this.buildingsLayer = this.scene.add.layer();
    this.buildingsLayer.setDepth(this.DEPTH_OFFSET);

    const edges = computeBastionEdges(layout);

    layout.buildings.forEach((building) => {
      const screenPos = IsoMath.subgridToScreen(
        building.x,
        building.y,
        building.size
      );

      const key = `${building.x},${building.y}`;
      const spriteOverride = edges.get(key) ?? null;
      const buildingData = spriteOverride
        ? { ...building, metadata: { ...building.metadata, sprite_key: spriteOverride } }
        : building;

      const buildingInstance = new BuildingSprite(
        this.scene,
        screenPos.x,
        screenPos.y,
        buildingData,
        this.interactive
      );

      buildingInstance.setDepth(screenPos.y + this.DEPTH_OFFSET);

      this.activeBuildings.push(buildingInstance);
      this.buildingsLayer.add(buildingInstance);
    });
  }

  public clearBuildings() {
    this.activeBuildings.forEach(s => s.destroy(true));
    this.activeBuildings = [];
  }
}
