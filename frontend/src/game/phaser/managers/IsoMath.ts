export class IsoMath {
  static TILE_W = 974;
  static TILE_H = 552;
  static SUBDIVISIONS = 3;

  // Converts a tile coordinate to pixel screen space
  static tileToScreen(x: number, y: number): { x: number; y: number } {
    return {
      x: (x - y) * (this.TILE_W / 2),
      y: (x + y) * (this.TILE_H / 2)
    };
  }

  // Converts a subgrid coordinate (0..gridW-1) to pixel screen space
  static subgridToScreen(subX: number, subY: number, buildingSize = 1): { x: number; y: number } {
    const gf = this.SUBDIVISIONS;
    const centerOffset = buildingSize / (gf * 2);
    const effectiveX = centerOffset - 0.5 + (subX / gf);
    const effectiveY = centerOffset - 0.5 + (subY / gf);

    return {
      x: (effectiveX - effectiveY) * (this.TILE_W / 2),
      y: (effectiveX + effectiveY) * (this.TILE_H / 2)
    };
  }

  // Converts API grid dimensions to tile count
  static gridToTiles(size: number): number {
    return Math.ceil(size / this.SUBDIVISIONS);
  }

  // Inverse: screen pixel position → nearest subgrid coordinate
  static screenToSubgrid(screenX: number, screenY: number, buildingSize = 1): { x: number; y: number } {
    const halfW = this.TILE_W / 2;
    const halfH = this.TILE_H / 2;
    const det = 2 * halfW * halfH;
    const tileX = (screenX * halfH + halfW * screenY) / det;
    const tileY = (halfW * screenY - halfH * screenX) / det;
    const gf = this.SUBDIVISIONS;
    const centerOffset = buildingSize / (gf * 2);
    return {
      x: Math.round((tileX + 0.5 - centerOffset) * gf),
      y: Math.round((tileY + 0.5 - centerOffset) * gf),
    };
  }
}
