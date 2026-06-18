import type Phaser from 'phaser';

export const BUILDING_ASSETS: Record<string, string> = {
  bastion: '/assets/bastion.png',
  'angel-spire': '/assets/angel-spire.png',
  'lament-basin': '/assets/lament-basin.png',
  barracks: '/assets/barracks.png',
  sanctum: '/assets/sanctum.png',
};

export function preloadBuildingAssets(loader: Phaser.Loader.LoaderPlugin) {
  loader.image('ground-tile', '/assets/ground-tile.png');
  loader.image('ground-tile-edge', '/assets/ground-tile-edge.png');
  Object.entries(BUILDING_ASSETS).forEach(([id, path]) => {
    loader.image(`building_${id}`, path);
  });
}
