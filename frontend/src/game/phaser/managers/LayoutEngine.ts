import Phaser from 'phaser';
import { IsoMath } from './IsoMath';
import type { BaseLayout } from '../../../types/building';
import { BuildingSprite } from '../entities/BuildingSprite';

export class LayoutEngine {
  private scene: Phaser.Scene;
  private buildingsLayer!: Phaser.GameObjects.Layer;
  public activeBuildings: BuildingSprite[] = [];

  private DEPTH_OFFSET = 100;

  constructor(scene: Phaser.Scene) {
    this.scene = scene;
    this.buildingsLayer = this.scene.add.layer();
    this.buildingsLayer.setDepth(this.DEPTH_OFFSET);
  }

  public renderLayout(layout: BaseLayout) {
    this.clearBuildings();

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
        building
      );

      buildingInstance.setDepth(screenPos.y + this.DEPTH_OFFSET);

      this.activeBuildings.push(buildingInstance);
      this.buildingsLayer.add(buildingInstance);
    });
  }

  public clearBuildings() {
    this.activeBuildings = [];
    this.buildingsLayer.removeAll(true);
  }
}
