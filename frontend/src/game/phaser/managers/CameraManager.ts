import Phaser from 'phaser';
import { IsoMath } from './IsoMath';

export class CameraManager {
  private cam: Phaser.Cameras.Scene2D.Camera;
  private isPanning = false;

  private mapCenterX = 0;
  private mapCenterY = 0;

  constructor(scene: Phaser.Scene) {
    this.cam = scene.cameras.main;
    this.cam.setBackgroundColor('#111111');
    this.setupControls(scene);
  }

  private setupControls(scene: Phaser.Scene) {
    scene.input.on('pointerdown', (pointer: Phaser.Input.Pointer) => {
      if (pointer.middleButtonDown()) {
        this.centerOnMap();
      }
      this.isPanning = true;
    });
    scene.input.on('pointerup', () => {
      this.isPanning = false;
    });

    scene.input.on('pointermove', (pointer: Phaser.Input.Pointer) => {
      if (!this.isPanning) return;
      this.cam.scrollX -= (pointer.x - pointer.prevPosition.x) / this.cam.zoom;
      this.cam.scrollY -= (pointer.y - pointer.prevPosition.y) / this.cam.zoom;
      this.clampScroll();
    });

    scene.input.on('wheel', (_: unknown, __: unknown, ___: unknown, deltaY: number) => {
      const newZoom = Phaser.Math.Clamp(this.cam.zoom - deltaY * 0.001, 0.05, 0.6);
      this.cam.setZoom(newZoom);
      this.clampScroll();
    });
  }

  private clampScroll() {
    const b = this.cam.getBounds();
    const vw = this.cam.width / this.cam.zoom;
    const vh = this.cam.height / this.cam.zoom;
    this.cam.scrollX = Phaser.Math.Clamp(this.cam.scrollX, b.x, b.x + b.width - vw);
    this.cam.scrollY = Phaser.Math.Clamp(this.cam.scrollY, b.y, b.y + b.height - vh);
  }

  public setMapSize(minTileX: number, minTileY: number, maxTileX: number, maxTileY: number) {
    const corners = [
      IsoMath.tileToScreen(minTileX, minTileY),
      IsoMath.tileToScreen(maxTileX, minTileY),
      IsoMath.tileToScreen(minTileX, maxTileY),
      IsoMath.tileToScreen(maxTileX, maxTileY),
    ];

    const halfW = IsoMath.TILE_W / 2;
    const halfH = IsoMath.TILE_H / 2;
    const padX = 5 * IsoMath.TILE_W;
    const padY = 5 * IsoMath.TILE_H;

    const minX = Math.min(...corners.map((c) => c.x)) - halfW - padX;
    const maxX = Math.max(...corners.map((c) => c.x)) + halfW + padX;
    const minY = Math.min(...corners.map((c) => c.y)) - halfH - padY;
    const maxY = Math.max(...corners.map((c) => c.y)) + halfH + padY;

    this.mapCenterX = (minX + maxX) / 2;
    this.mapCenterY = (minY + maxY) / 2;

    this.cam.setBounds(minX, minY, maxX - minX, maxY - minY);
  }

  public centerOnMap(cx?: number, cy?: number) {
    this.fitZoomToBounds();
    if (cx !== undefined && cy !== undefined) {
      this.cam.centerOn(cx, cy);
    } else {
      this.cam.centerOn(this.mapCenterX, this.mapCenterY);
    }
    this.clampScroll();
  }

  private fitZoomToBounds() {
    const b = this.cam.getBounds();
    if (b.width <= 0 || b.height <= 0) return;
    const fw = this.cam.width;
    const fh = this.cam.height;
    const zoom = Math.min(fw / b.width, fh / b.height, 0.5);
    this.cam.setZoom(Math.max(0.05, zoom));
  }
}
