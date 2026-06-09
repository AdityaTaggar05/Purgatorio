import Phaser from 'phaser';

export class TerraceScene extends Phaser.Scene {
  private tileWidth = 880;
  private tileHeight = 542;
  private mapSize = 10;

  private minZoom = 0.2;
  private maxZoom = 0.6;

  constructor() {
    super({ key: 'TerraceScene' });
  }

  preload() {
    this.load.image('ground-tile', 'assets/ground-tile.png');
  }

  create() {
    this.cameras.main.setBackgroundColor('#111111');
    this.cameras.main.setZoom(0.4);

    this.cameras.main.centerOn(0, this.tileHeight * this.mapSize / 2);

    const groundLayer = this.add.layer();
    //const worldLayer = this.add.layer();

    for (let y = 0; y < this.mapSize; y++) {
      for (let x = 0; x < this.mapSize; x++) {
        const { isoX, isoY } = this.cartesianToIso(x, y);

        if (this.textures.exists('ground-tile')) {
          const tile = this.add.image(isoX, isoY, 'ground-tile');
          tile.setOrigin(0.5, 0.5);
          tile.setSize(880, 542)
          groundLayer.add(tile)
        } else {
          this.drawMockTile(isoX, isoY);
        }
      }
    }

    this.input.on('pointermove', (pointer: Phaser.Input.Pointer) => {
      if (!pointer.isDown) return; // Only execute if mouse click is held down

      this.cameras.main.scrollX -= (pointer.x - pointer.prevPosition.x) / this.cameras.main.zoom;
      this.cameras.main.scrollY -= (pointer.y - pointer.prevPosition.y) / this.cameras.main.zoom;
    });

    this.input.on('wheel', (_: unknown, __: unknown, ___: unknown, deltaY: number) => {
      // Calculate new target zoom factor based on scroll depth
      const zoomFactor = 0.1;
      let newZoom = this.cameras.main.zoom - deltaY * zoomFactor * 0.01;

      // Clamp values so user doesn't zoom out into oblivion or past max resolution
      newZoom = Phaser.Math.Clamp(newZoom, this.minZoom, this.maxZoom);

      this.cameras.main.setZoom(newZoom);
    });
  }

  private cartesianToIso(mapX: number, mapY: number) {
    const isoX = (mapX - mapY) * (this.tileWidth / 2);
    const isoY = (mapX + mapY) * (this.tileHeight / 2);
    return { isoX, isoY };
  }

  private drawMockTile(x: number, y: number) {
    const graphics = this.add.graphics();
    graphics.lineStyle(1, 0x2c2927, 1);
    graphics.fillStyle(0x181716, 1);

    graphics.beginPath();
    graphics.moveTo(x, y - this.tileHeight / 2);                 // Top
    graphics.lineTo(x + this.tileWidth / 2, y);                  // Right
    graphics.lineTo(x, y + this.tileHeight / 2);                 // Bottom
    graphics.lineTo(x - this.tileWidth / 2, y);                  // Left
    graphics.closePath();

    graphics.fillPath();
    graphics.strokePath();
  }
}
