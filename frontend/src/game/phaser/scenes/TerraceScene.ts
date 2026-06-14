import Phaser from "phaser";
import { latestLayout } from "../GameCanvas";
import { phaserEvents } from "../events";
import { CameraManager } from "../managers/CameraManager";
import { IsoMath } from "../managers/IsoMath";
import { LayoutEngine } from "../managers/LayoutEngine";
import { TerrainEngine } from "../managers/TerrainEngine";

export class TerraceScene extends Phaser.Scene {
  private terrain!: TerrainEngine;
  private layoutEngine!: LayoutEngine;
  private cameraManager!: CameraManager;
  private ghostSprite!: Phaser.GameObjects.Sprite;
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

    this.ghostSprite = this.add.sprite(0, 0, 'building_bastion');
    this.ghostSprite.setVisible(false);
    this.ghostSprite.setAlpha(0.5);
    this.ghostSprite.setDepth(9999);

    this.input.on('pointerdown', (pointer: Phaser.Input.Pointer) => {
      if (pointer.middleButtonDown()) return;

      const worldPoint = this.cameras.main.getWorldPoint(pointer.x, pointer.y);
      const sub = IsoMath.screenToSubgrid(worldPoint.x, worldPoint.y);

      if (phaserEvents.mode === "place" && phaserEvents.onGridClick) {
        phaserEvents.onGridClick(sub.x, sub.y);
        return;
      }
    });

    this.input.on('pointermove', (pointer: Phaser.Input.Pointer) => {
      if (phaserEvents.mode !== "place" || !phaserEvents.placementBuilding) {
        this.ghostSprite.setVisible(false);
        return;
      }

      const worldPoint = this.cameras.main.getWorldPoint(pointer.x, pointer.y);
      const sub = IsoMath.screenToSubgrid(worldPoint.x, worldPoint.y);

      const size = phaserEvents.placementBuilding.size;
      const clampedX = Phaser.Math.Clamp(sub.x, 0, 29 - size + 1);
      const clampedY = Phaser.Math.Clamp(sub.y, 0, 29 - size + 1);
      phaserEvents.ghostPosition = { x: clampedX, y: clampedY };

      const screenPos = IsoMath.subgridToScreen(clampedX, clampedY, size);
      this.ghostSprite.setPosition(screenPos.x, screenPos.y);
      this.ghostSprite.setTexture(`building_${phaserEvents.placementBuilding.id}`);
      this.ghostSprite.setVisible(true);

      const tileW = IsoMath.TILE_W / IsoMath.SUBDIVISIONS;
      const scale = (tileW * size) / this.ghostSprite.width;
      this.ghostSprite.setScale(scale);
      this.ghostSprite.setOrigin(0.5, 1);
    });

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
