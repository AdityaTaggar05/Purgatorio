import Phaser from "phaser";
import { latestLayout } from "../GameCanvas";
import { phaserEvents } from "../events";
import { CameraManager } from "../managers/CameraManager";
import { IsoMath } from "../managers/IsoMath";
import { LayoutEngine } from "../managers/LayoutEngine";
import { TerrainEngine } from "../managers/TerrainEngine";

export class TerraceScene extends Phaser.Scene {
  static SUBGRID_LAYER = "subgrid";

  private terrain!: TerrainEngine;
  private layoutEngine!: LayoutEngine;
  private cameraManager!: CameraManager;
  private ghostSprite!: Phaser.GameObjects.Sprite;
  private subgridGraphics!: Phaser.GameObjects.Graphics;
  private initialized = false;
  private subgridDrawn = false;

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

    this.subgridGraphics = this.add.graphics();
    this.subgridGraphics.setDepth(50);
    this.subgridGraphics.setVisible(false);

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

      const size = phaserEvents.placementBuilding.size;
      const sub = IsoMath.screenToSubgrid(worldPoint.x, worldPoint.y, size);

      const clampedX = Phaser.Math.Clamp(sub.x, 0, 29 - size + 1);
      const clampedY = Phaser.Math.Clamp(sub.y, 0, 29 - size + 1);
      phaserEvents.ghostPosition = { x: clampedX, y: clampedY };

      const screenPos = IsoMath.subgridToScreen(clampedX, clampedY, size);
      const tileH = IsoMath.TILE_H / IsoMath.SUBDIVISIONS;
      this.ghostSprite.setPosition(screenPos.x, screenPos.y + (tileH * size) / 2);
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

    const inPlaceMode = phaserEvents.mode === "place";
    if (inPlaceMode && !this.subgridDrawn) {
      this.drawSubgridOverlay();
      this.subgridDrawn = true;
    } else if (!inPlaceMode && this.subgridDrawn) {
      this.subgridGraphics.clear();
      this.subgridGraphics.setVisible(false);
      this.subgridDrawn = false;
    }
  }

  private drawSubgridOverlay() {
    const gf = IsoMath.SUBDIVISIONS;
    const gridW = 30;
    const gridH = 30;

    const toScreen = (sx: number, sy: number) => IsoMath.tileToScreen(sx / gf, sy / gf);

    this.subgridGraphics.clear();

    // Subgrid checkerboard fill
    for (let x = 0; x < gridW; x++) {
      for (let y = 0; y < gridH; y++) {
        if ((x + y) % 2 === 1) continue;
        const a = toScreen(x, y);
        const b = toScreen(x + 1, y);
        const c = toScreen(x + 1, y + 1);
        const d = toScreen(x, y + 1);

        this.subgridGraphics.fillStyle(0xffffff, 0.4);
        this.subgridGraphics.beginPath();
        this.subgridGraphics.moveTo(a.x, a.y);
        this.subgridGraphics.lineTo(b.x, b.y);
        this.subgridGraphics.lineTo(c.x, c.y);
        this.subgridGraphics.lineTo(d.x, d.y);
        this.subgridGraphics.closePath();
        this.subgridGraphics.fillPath();
      }
    }

    // Subgrid lines
    this.subgridGraphics.lineStyle(4, 0xffffff, 0.6);
    for (let x = 0; x <= gridW; x++) {
      const start = toScreen(x, 0);
      const end = toScreen(x, gridH);
      this.subgridGraphics.beginPath();
      this.subgridGraphics.moveTo(start.x, start.y);
      this.subgridGraphics.lineTo(end.x, end.y);
      this.subgridGraphics.strokePath();
    }
    for (let y = 0; y <= gridH; y++) {
      const start = toScreen(0, y);
      const end = toScreen(gridW, y);
      this.subgridGraphics.beginPath();
      this.subgridGraphics.moveTo(start.x, start.y);
      this.subgridGraphics.lineTo(end.x, end.y);
      this.subgridGraphics.strokePath();
    }

    // Ground-tile boundaries (bold)
    const tilesW = IsoMath.gridToTiles(gridW);
    const tilesH = IsoMath.gridToTiles(gridH);
    this.subgridGraphics.lineStyle(2, 0xffffff, 0.15);
    for (let tx = 0; tx <= tilesW; tx++) {
      const x = tx * gf;
      const start = toScreen(x, 0);
      const end = toScreen(x, gridH);
      this.subgridGraphics.beginPath();
      this.subgridGraphics.moveTo(start.x, start.y);
      this.subgridGraphics.lineTo(end.x, end.y);
      this.subgridGraphics.strokePath();
    }
    for (let ty = 0; ty <= tilesH; ty++) {
      const y = ty * gf;
      const start = toScreen(0, y);
      const end = toScreen(gridW, y);
      this.subgridGraphics.beginPath();
      this.subgridGraphics.moveTo(start.x, start.y);
      this.subgridGraphics.lineTo(end.x, end.y);
      this.subgridGraphics.strokePath();
    }

    this.subgridGraphics.setVisible(true);
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
