import Phaser from 'phaser';
import { IsoMath } from './IsoMath';

export class TerrainEngine {
  private scene: Phaser.Scene;
  private groundLayer!: Phaser.GameObjects.Layer;

  constructor(scene: Phaser.Scene) {
    this.scene = scene;
    this.initLayers();
  }

  private initLayers() {
    this.groundLayer = this.scene.add.layer();
    this.groundLayer.setDepth(0);
  }

  public generateGroundGrid(tilesW: number, tilesH: number, padding: number) {
    const minX = -padding;
    const minY = -padding;
    const maxX = tilesW + padding;
    const maxY = tilesH + padding;

    const tilePositions: Array<{ x: number; y: number; gridX: number; gridY: number; edge: boolean }> = [];

    for (let y = minY; y < maxY; y++) {
      for (let x = minX; x < maxX; x++) {
        const screenPos = IsoMath.tileToScreen(x, y);
        tilePositions.push({
          x: screenPos.x,
          y: screenPos.y,
          gridX: x,
          gridY: y,
          edge: x === maxX - 1 || y === maxY - 1,
        });
      }
    }

    tilePositions.sort((a, b) => a.y - b.y);

    tilePositions.forEach((pos) => {
      const tileKey = pos.edge ? 'ground-tile-edge' : 'ground-tile';
      const tile = this.scene.add.image(pos.x, pos.y, tileKey);

      tile.setOrigin(0.5, 0.45);
      tile.setDepth(pos.gridX + pos.gridY);
      tile.setData('gridCoordinates', { x: pos.gridX, y: pos.gridY });

      this.groundLayer.add(tile);
    });
  }

  public destroyMap() {
    this.groundLayer.removeAll(true);
  }
}
