import Phaser from 'phaser';
import { IsoMath } from './IsoMath';

export class CameraManager {
  private cam: Phaser.Cameras.Scene2D.Camera;
  private isPanning = false;

  private mapW = 30;
  private mapH = 30;

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
    });

    scene.input.on('wheel', (_: unknown, __: unknown, ___: unknown, deltaY: number) => {
      const newZoom = Phaser.Math.Clamp(this.cam.zoom - deltaY * 0.001, 0.1, 0.6);
      this.cam.setZoom(newZoom);
    });
  }

  public setMapSize(tilesW: number, tilesH: number) {
    this.mapW = tilesW;
    this.mapH = tilesH;

    const totalMapWidth = tilesW * IsoMath.TILE_W;
    const totalMapHeight = tilesH * IsoMath.TILE_H;

    const horizontalPadding = 1000;
    const verticalPadding = 1600;

    this.cam.setBounds(
      -totalMapWidth / 2 - horizontalPadding,
      -verticalPadding / 2,
      totalMapWidth + horizontalPadding * 2,
      totalMapHeight + verticalPadding / 2
    );
  }

  public centerOnMap() {
    this.cam.setZoom(0.3);

    const centerPos = IsoMath.tileToScreen(this.mapW / 2, this.mapH / 2);
    this.cam.centerOn(centerPos.x, centerPos.y);
  }
}
