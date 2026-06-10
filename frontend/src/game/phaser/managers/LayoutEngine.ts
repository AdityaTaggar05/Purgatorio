import Phaser from 'phaser';
import type { BaseLayout } from '../../../types/building';
import { IsoMath } from './IsoMath';
import { BuildingSprite } from '../entities//BuildingSprite';

export class LayoutEngine {
  private scene: Phaser.Scene;
  private buildingsLayer!: Phaser.GameObjects.Layer;
  public activeBuildings: BuildingSprite[] = [];

  private DEPTH_OFFSET = 100

  constructor(scene: Phaser.Scene) {
    this.scene = scene;
    this.buildingsLayer = this.scene.add.layer();
    // Buildings sit safely on top of ground tiles
    this.buildingsLayer.setDepth(this.DEPTH_OFFSET);
  }

  public renderLayout(layout: BaseLayout) {
    this.clearBuildings();

    layout.buildings.forEach((buildingData) => {
      const screenPos = IsoMath.subgridToScreen(
        layout.subgrid_factor,
        buildingData.x,
        buildingData.y,
        buildingData.size
      );

      const buildingInstance = new BuildingSprite(
        this.scene,
        layout.subgrid_factor,
        screenPos.x,
        screenPos.y,
        buildingData
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
