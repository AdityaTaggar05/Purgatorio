import Phaser from 'phaser';
import { IsoMath } from './IsoMath';
import { TroopStackSprite, type TroopStackData } from '../entities/TroopStackSprite';

const DEPTH_OFFSET = 101; // one above LayoutEngine's building depth offset

export class TroopOverlayEngine {
  private scene: Phaser.Scene;
  private layer: Phaser.GameObjects.Layer;
  private sprites = new Map<string, TroopStackSprite>();
  private staticPreview: boolean;

  constructor(scene: Phaser.Scene, staticPreview = false) {
    this.scene = scene;
    this.staticPreview = staticPreview;
    this.layer = scene.add.layer();
  }

  public sync(stacks: TroopStackData[]) {
    const seen = new Set<string>();

    stacks.forEach((stack) => {
      seen.add(stack.id);
      const existing = this.sprites.get(stack.id);
      if (existing) {
        existing.update(stack);
        return;
      }

      const sprite = new TroopStackSprite(this.scene, stack, this.staticPreview);
      const screenPos = IsoMath.subgridToScreen(stack.x, stack.y, 1);
      sprite.setDepth(screenPos.y + DEPTH_OFFSET);
      this.layer.add(sprite);
      this.sprites.set(stack.id, sprite);
    });

    for (const [id, sprite] of this.sprites) {
      if (!seen.has(id)) {
        sprite.destroy();
        this.sprites.delete(id);
      }
    }
  }

  public clear() {
    this.sprites.forEach((sprite) => sprite.destroy());
    this.sprites.clear();
  }
}
