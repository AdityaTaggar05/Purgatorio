import Phaser from 'phaser';
import { IsoMath } from './IsoMath';

export class CameraManager {
  private cam: Phaser.Cameras.Scene2D.Camera;
  private isPanning = false;

  private mapSize?: number;

  constructor(scene: Phaser.Scene) {
    this.cam = scene.cameras.main;

    this.cam.setBackgroundColor('#111111');
    this.setupControls(scene);
  }

  private setupControls(scene: Phaser.Scene) {
    scene.input.on('pointerdown', (pointer) => {
      if (pointer.middleButtonDown() && this.mapSize) {
        this.centerOnMap(this.mapSize)
      }

      return this.isPanning = true;
    });
    scene.input.on('pointerup', () => this.isPanning = false);

    scene.input.on('pointermove', (pointer: Phaser.Input.Pointer) => {
      if (!this.isPanning) return;
      this.cam.scrollX -= (pointer.x - pointer.prevPosition.x) / this.cam.zoom;
      this.cam.scrollY -= (pointer.y - pointer.prevPosition.y) / this.cam.zoom;
    });

    scene.input.on('wheel', (_: unknown, __: unknown, ___: unknown, deltaY: number) => {
      const newZoom = Phaser.Math.Clamp(this.cam.zoom - deltaY * 0.001, 0.15, 0.6);
      this.cam.setZoom(newZoom);
    });
  }

  public centerOnMap(mapSize: number) {
    if (!this.mapSize) this.mapSize = mapSize;

    const halfMap = mapSize / 2;

    this.cam.setZoom(0.4);
    this.cam.centerOn(0, halfMap * IsoMath.TILE_H);
  }

  public setBoundsFromMap(mapSize: number) {
    const totalMapWidth = mapSize * IsoMath.TILE_W;
    const totalMapHeight = mapSize * IsoMath.TILE_H;

    const horizontalPadding = 1000;
    const verticalPadding = 1600;

    this.cam.setBounds(
      -totalMapWidth / 2 - horizontalPadding,
      -verticalPadding / 2,
      totalMapWidth + (horizontalPadding * 2),
      totalMapHeight + verticalPadding / 2
    );
  }
}
