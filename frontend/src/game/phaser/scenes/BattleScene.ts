import Phaser from "phaser";
import {
  latestDefenderLayout,
  latestBuildingHpByKey,
  latestTroopStacks,
  latestDeploymentZone,
  latestSelectedCells,
  deployClickEnabled,
  troopCatalogIds,
  staticPreviewMode,
  latestSelectedTroop,
} from "../BattleCanvas";
import { battleEvents } from "../battleEvents";
import { CameraManager } from "../managers/CameraManager";
import { IsoMath } from "../managers/IsoMath";
import { LayoutEngine } from "../managers/LayoutEngine";
import { TerrainEngine } from "../managers/TerrainEngine";
import { TroopOverlayEngine } from "../managers/TroopOverlayEngine";
import { preloadBuildingAssets } from "../assetManifest";

const GF = IsoMath.SUBDIVISIONS;
const DEFAULT_GRID = 30;
const PADDING = 6;          // extra subgrid cells on each side for deployment
const PADDING_TILES = Math.ceil(PADDING / GF);

export class BattleScene extends Phaser.Scene {
  private terrain!: TerrainEngine;
  private layoutEngine!: LayoutEngine;
  private troopOverlay!: TroopOverlayEngine;
  private cameraManager!: CameraManager;
  private zoneGraphics!: Phaser.GameObjects.Graphics;
  private troopGhost!: Phaser.GameObjects.Sprite;

  private lastLayoutRef: unknown = undefined;
  private lastGroundDims = "";

  preload() {
    preloadBuildingAssets(this.load);

    troopCatalogIds.forEach((id) => {
      this.load.image(`troop_${id}`, `/assets/${id}.png`);
    });
  }

  create() {
    this.cameraManager = new CameraManager(this);
    this.terrain = new TerrainEngine(this);
    this.layoutEngine = new LayoutEngine(this, false);
    this.troopOverlay = new TroopOverlayEngine(this, staticPreviewMode);

    this.zoneGraphics = this.add.graphics();
    this.zoneGraphics.setDepth(50);

    this.troopGhost = this.add.sprite(0, 0, "");
    this.troopGhost.setVisible(false);
    this.troopGhost.setAlpha(0.6);
    this.troopGhost.setDepth(200);
    this.troopGhost.setOrigin(0.5, 1);

    this.input.on("pointerdown", (pointer: Phaser.Input.Pointer) => {
      if (pointer.middleButtonDown()) return;
      if (!deployClickEnabled || !battleEvents.onCellClick) return;

      const worldPoint = this.cameras.main.getWorldPoint(pointer.x, pointer.y);
      const sub = IsoMath.screenToSubgrid(worldPoint.x, worldPoint.y, 1);
      battleEvents.onCellClick(sub.x, sub.y);
    });

    this.input.on("pointermove", (pointer: Phaser.Input.Pointer) => {
      if (!latestSelectedTroop || !deployClickEnabled) {
        this.troopGhost.setVisible(false);
        return;
      }

      const worldPoint = this.cameras.main.getWorldPoint(pointer.x, pointer.y);
      const sub = IsoMath.screenToSubgrid(worldPoint.x, worldPoint.y, 1);
      const screenPos = IsoMath.subgridToScreen(sub.x, sub.y, 1);
      this.troopGhost.setPosition(screenPos.x, screenPos.y);
      this.troopGhost.setTexture(`troop_${latestSelectedTroop}`);

      const tileW = IsoMath.TILE_W / GF;
      const scale = tileW / this.troopGhost.width;
      this.troopGhost.setScale(scale);
      this.troopGhost.setVisible(true);
    });

    this.renderGround();
    this.renderLayout();
  }

  update() {
    this.renderGround();
    this.renderLayout();
    this.applyBuildingDamage();
    this.troopOverlay.sync(latestTroopStacks);
    this.drawZoneOverlay();
  }

  private renderGround() {
    const gridW = latestDefenderLayout?.grid_w ?? DEFAULT_GRID;
    const gridH = latestDefenderLayout?.grid_h ?? DEFAULT_GRID;
    const dims = `${gridW}x${gridH}`;
    if (dims === this.lastGroundDims) return;
    this.lastGroundDims = dims;

    const tilesW = IsoMath.gridToTiles(gridW);
    const tilesH = IsoMath.gridToTiles(gridH);

    this.terrain.destroyMap();
    this.terrain.generateGroundGrid(tilesW, tilesH, PADDING_TILES);

    const totalTW = tilesW + PADDING_TILES * 2;
    const totalTH = tilesH + PADDING_TILES * 2;
    this.cameraManager.setMapSize(-PADDING_TILES, -PADDING_TILES, totalTW - PADDING_TILES, totalTH - PADDING_TILES);
    const centerPos = IsoMath.tileToScreen(tilesW / 2, tilesH / 2);
    this.cameraManager.centerOnMap(centerPos.x, centerPos.y);
  }

  private renderLayout() {
    if (latestDefenderLayout === this.lastLayoutRef) return;
    this.lastLayoutRef = latestDefenderLayout;

    if (latestDefenderLayout) {
      this.layoutEngine.renderLayout(latestDefenderLayout);
    } else {
      this.layoutEngine.clearBuildings();
    }
  }

  // Building positions/sprites are only (re)created when the snapshot
  // itself changes (see renderLayout). HP, however, changes every tick
  // during playback, so it's applied in place against the existing
  // sprites — same reasoning as TroopOverlayEngine.sync vs full re-render.
  private applyBuildingDamage() {
    this.layoutEngine.activeBuildings.forEach((sprite) => {
      const key = `${sprite.buildingData.x}_${sprite.buildingData.y}`;
      const hp = latestBuildingHpByKey[key];
      if (hp === undefined) return;
      sprite.applyHp(hp);
    });
  }

  private drawZoneOverlay() {
    this.zoneGraphics.clear();
    if (!latestDeploymentZone) return;

    const selectedSet = new Set(latestSelectedCells.map((c) => `${c.x},${c.y}`));

    latestDeploymentZone.forEach(({ x, y }) => {
      const isSelected = selectedSet.has(`${x},${y}`);
      const a = IsoMath.subgridToScreen(x, y, 1);
      const b = IsoMath.subgridToScreen(x + 1, y, 1);
      const c = IsoMath.subgridToScreen(x + 1, y + 1, 1);
      const d = IsoMath.subgridToScreen(x, y + 1, 1);

      this.zoneGraphics.fillStyle(0x22c55e, isSelected ? 0.45 : 0.12);
      this.zoneGraphics.lineStyle(2, 0x22c55e, isSelected ? 0.8 : 0.35);
      this.zoneGraphics.beginPath();
      this.zoneGraphics.moveTo(a.x, a.y);
      this.zoneGraphics.lineTo(b.x, b.y);
      this.zoneGraphics.lineTo(c.x, c.y);
      this.zoneGraphics.lineTo(d.x, d.y);
      this.zoneGraphics.closePath();
      this.zoneGraphics.fillPath();
      this.zoneGraphics.strokePath();
    });
  }
}
