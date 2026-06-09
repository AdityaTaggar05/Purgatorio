import type { BaseLayout } from '../../../types/building';

export class LayoutParser {
  static createOccupancyGrid(layout: BaseLayout): string[][] {
    const gridSize = layout.tiles * layout.subgridFactor;

    const matrix: string[][] = Array(gridSize).fill(null).map(() => Array(gridSize).fill('empty'));

    layout.buildings.forEach((building) => {
      for (let y = 0; y < building.size; y++) {
        for (let x = 0; x < building.size; x++) {
          const targetX = building.x + x;
          const targetY = building.y + y;
          
          if (targetX < gridSize && targetY < gridSize) {
            matrix[targetY][targetX] = building.id;
          }
        }
      }
    });

    return matrix;
  }
}
