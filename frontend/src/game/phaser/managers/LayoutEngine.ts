import Phaser from 'phaser';
import { IsoMath } from './IsoMath';
import type { BaseLayout } from '../../../types/building';
import { BuildingSprite } from '../entities/BuildingSprite';
import { phaserEvents } from '../events';

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

    layout.buildings.forEach((building) => {
      const screenPos = IsoMath.subgridToScreen(
        building.x,
        building.y,
        building.size
      );

      const buildingInstance = new BuildingSprite(
        this.scene,
        screenPos.x,
        screenPos.y,
        building,
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
