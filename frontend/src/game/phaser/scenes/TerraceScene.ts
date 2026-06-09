import Phaser from 'phaser';

export class TerraceScene extends Phaser.Scene {
  private tileWidth = 974;
  private tileHeight = 552;
  private mapSize = 10;

  private minZoom = 0.2;
  private maxZoom = 0.6;

  constructor() {
    super({ key: 'TerraceScene' });
  }

  preload() {
    this.load.image('ground-tile', 'assets/ground-tile.png');
    this.load.image('ground-tile-edge', 'assets/ground-tile-edge.png');
  }

  create() {
    this.cameras.main.setBackgroundColor('#111111');
    this.cameras.main.setZoom(0.4);

    this.cameras.main.centerOn(0, this.tileHeight * this.mapSize / 2);

    const groundLayer = this.add.layer();
    //const worldLayer = this.add.layer();

    const texture = this.textures.get('ground-tile')
    if (texture) {
      texture.setFilter(Phaser.Textures.FilterMode.LINEAR);
    }

    const tilePositions = [];

    for (let y = 0; y < this.mapSize; y++) {
      for (let x = 0; x < this.mapSize; x++) {
        const { isoX, isoY } = this.cartesianToIso(x, y);
        tilePositions.push({ x: isoX, y: isoY, depth: x + y, edge: y==this.mapSize-1 || x==this.mapSize-1 });
      }
    }

    tilePositions.sort((a, b) => a.depth - b.depth);

    tilePositions.forEach((pos) => {
      let tile;

      if (pos.edge) tile = this.add.image(pos.x, pos.y, 'ground-tile-edge')
      else tile = this.add.image(pos.x, pos.y, 'ground-tile');

      tile.setOrigin(0.5, 0.181);
      tile.setDepth(pos.y);
      groundLayer.add(tile)
    });

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
}
