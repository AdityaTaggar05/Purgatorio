import type { BaseLayout } from "../../../types/building";
import { CameraManager } from "../managers/CameraManager";
import { LayoutEngine } from "../managers/LayoutEngine";
import { TerrainEngine } from "../managers/TerrainEngine";
import Phaser from "phaser";

export class TerraceScene extends Phaser.Scene {
  private terrain!: TerrainEngine;
  private layoutEngine!: LayoutEngine;
  private cameraManager!: CameraManager;

  preload() {
    this.load.image('ground-tile', '/assets/ground-tile.png')
    this.load.image('ground-tile-edge', '/assets/ground-tile-edge.png')

    this.load.image('building_bastion', '/assets/bastion.png');
    this.load.image('building_angel-spire', '/assets/angel-spire.png');
    this.load.image('building_lament-basin', '/assets/lament-basin.png');
  }

  create() {
    this.cameraManager = new CameraManager(this);
    this.terrain = new TerrainEngine(this);
    this.layoutEngine = new LayoutEngine(this);

    const baseLayoutJSON = `
    {
      "user_id": "user_penitent_777",
      "tiles": 10,
      "subgrid_factor": 3,
      "buildings": [
        {
          "id": "bastion_001",
          "x": 0,
          "y": 0,
          "size": 1
        },
        {
          "id": "angel-spire_001",
          "subgridX": 4,
          "subgridY": 12,
          "size": 2
        },
        {
          "id": "lament-basin_001",
          "subgridX": 15,
          "subgridY": 15,
          "size": 2
        }
      ]
    }
    `

    const baseLayout: BaseLayout = JSON.parse(baseLayoutJSON)

    this.cameraManager.centerOnMap(baseLayout.tiles)
    this.cameraManager.setBoundsFromMap(baseLayout.tiles)

    this.terrain.generateGroundGrid(baseLayout.tiles)
    this.layoutEngine.renderLayout(baseLayout);
  }
}
