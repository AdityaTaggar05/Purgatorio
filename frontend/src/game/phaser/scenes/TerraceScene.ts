import Phaser from "phaser";
import { latestLayout } from "../GameCanvas";
import { CameraManager } from "../managers/CameraManager";
import { IsoMath } from "../managers/IsoMath";
import { LayoutEngine } from "../managers/LayoutEngine";
import { TerrainEngine } from "../managers/TerrainEngine";

export class TerraceScene extends Phaser.Scene {
  private terrain!: TerrainEngine;
  private layoutEngine!: LayoutEngine;
  private cameraManager!: CameraManager;
  private initialized = false;

  preload() {
    this.load.image('ground-tile', '/assets/ground-tile.png');
    this.load.image('ground-tile-edge', '/assets/ground-tile-edge.png');

    this.load.image('building_bastion', '/assets/bastion.png');
    this.load.image('building_angel-spire', '/assets/angel-spire.png');
    this.load.image('building_lament-basin', '/assets/lament-basin.png');
    this.load.image('building_barracks', '/assets/barracks.png');
    this.load.image('building_sanctum', '/assets/sanctum.png');
  }

  create() {
    this.cameraManager = new CameraManager(this);
    this.terrain = new TerrainEngine(this);
    this.layoutEngine = new LayoutEngine(this);
    this.initialized = true;
    this.tryRender();
  }

  update() {
    if (this.initialized && latestLayout) {
      this.tryRender();
    }
  }

  private lastLayoutKey = "";

  private tryRender() {
    if (!latestLayout) return;
    const key = JSON.stringify(latestLayout);
    if (key === this.lastLayoutKey) return;
    this.lastLayoutKey = key;

    const tilesW = IsoMath.gridToTiles(latestLayout.grid_w);
    const tilesH = IsoMath.gridToTiles(latestLayout.grid_h);

    this.cameraManager.setMapSize(tilesW, tilesH);
    this.cameraManager.centerOnMap();

    this.terrain.destroyMap();
    this.terrain.generateGroundGrid(tilesW, tilesH);
    this.layoutEngine.renderLayout(latestLayout);
  }
}
