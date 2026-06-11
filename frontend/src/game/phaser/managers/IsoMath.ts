export class IsoMath {
  static TILE_W = 974;
  static TILE_H = 552;

  // Converts a 10x10 tile coordinate to pixel screen space
  static tileToScreen(x: number, y: number): { x: number; y: number } {
    return {
      x: (x - y) * (this.TILE_W / 2),
      y: (x + y) * (this.TILE_H / 2)
    };
  }

  // Converts a 30x30 subgrid coordinate to pixel screen space
  static subgridToScreen(gridFactor: number, subX: number, subY: number, buildingSize = 1): { x: number; y: number } {
    const centerOffset = buildingSize / (gridFactor * 2);
    const effectiveX = centerOffset - 0.5 + (subX / gridFactor);
    const effectiveY = centerOffset - 0.5 + (subY / gridFactor);

    return {
      x: (effectiveX - effectiveY) * (this.TILE_W / 2),
      y: (effectiveX + effectiveY) * (this.TILE_H / 2)
    };
  }
}
