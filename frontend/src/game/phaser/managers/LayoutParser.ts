import type { BaseLayout } from '../../../types/building';

export class LayoutParser {
  static createOccupancyGrid(layout: BaseLayout): string[][] {
    const gridSizeX = layout.grid_w;
    const gridSizeY = layout.grid_h;

    const matrix: string[][] = Array(gridSizeY).fill(null).map(() => Array(gridSizeX).fill('empty'));

    layout.buildings.forEach((building) => {
      for (let y = 0; y < building.size; y++) {
        for (let x = 0; x < building.size; x++) {
          const targetX = building.x + x;
          const targetY = building.y + y;

          if (targetX < gridSizeX && targetY < gridSizeY) {
            matrix[targetY][targetX] = building.building_id;
          }
        }
      }
    });

    return matrix;
  }
}
