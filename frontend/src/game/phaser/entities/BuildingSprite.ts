import Phaser from 'phaser';
import type { PlacedBuilding } from '../../../types/building';
import { IsoMath } from '../managers/IsoMath';
import { phaserEvents } from '../events';

const SPRITE_ORIGIN = { x: 0.5, y: 1.0 };
const GF = IsoMath.SUBDIVISIONS;

export class BuildingSprite extends Phaser.GameObjects.Container {
  public buildingData: PlacedBuilding;
  private mainSprite: Phaser.GameObjects.Sprite;
  private selectionRing!: Phaser.GameObjects.Graphics;
  private spriteW: number;
  private spriteH: number;

  public currentHealth: number;
  public maxHealth: number;

  constructor(scene: Phaser.Scene, x: number, y: number, data: PlacedBuilding) {
    super(scene, x, y);
    this.buildingData = data;

    this.maxHealth = data.hp ?? data.size * 500;
    this.currentHealth = this.maxHealth;

    const spriteKey = `building_${data.building_id}`;
    this.spriteW = (IsoMath.TILE_W / GF) * data.size;
    this.spriteH = (IsoMath.TILE_H / GF) * data.size;

    try {
      this.mainSprite = scene.add.sprite(0, 0, spriteKey);
      this.mainSprite.setOrigin(SPRITE_ORIGIN.x, SPRITE_ORIGIN.y);

      const scale = this.spriteW / this.mainSprite.width;
      this.mainSprite.setScale(scale);
      this.mainSprite.setPosition(0, this.spriteH / 2);

      this.add(this.mainSprite);
    } catch (_) {
      console.warn(`Asset binding failed for frame: ${spriteKey}. Using debug fallback.`);
    }

    this.setupInteractions();

    scene.add.existing(this);
  }

  private setupInteractions() {
    const scaledHeight = this.mainSprite ? this.mainSprite.height : this.spriteH;

    this.setInteractive(
      new Phaser.Geom.Rectangle(
        -this.spriteW / 2,
        this.spriteH / 2 - scaledHeight,
        this.spriteW,
        scaledHeight
      ),
      Phaser.Geom.Rectangle.Contains
    );

    this.on('pointerover', () => {
      this.scene.input.setDefaultCursor('pointer');
      this.showSelectionRing();
    });

    this.on('pointerout', () => {
      this.scene.input.setDefaultCursor('default');
      this.hideSelectionRing();
    });

    this.on('pointerdown', () => {
      if (phaserEvents.mode === "none" && phaserEvents.onBuildingClick) {
        phaserEvents.onBuildingClick(this.buildingData);
      }
    });
  }

  private showSelectionRing() {
    if (!this.selectionRing) {
      this.selectionRing = this.scene.add.graphics();
      this.selectionRing.setPosition(this.x, this.y);
      this.selectionRing.lineStyle(3, 0x00ff00, 1.0);
      this.selectionRing.fillStyle(0x00ff00, 0.2);

      const hw = this.spriteW / 2;
      const hh = this.spriteH / 2;

      this.selectionRing.beginPath();
      this.selectionRing.moveTo(0, -hh);
      this.selectionRing.lineTo(hw, 0);
      this.selectionRing.lineTo(0, hh);
      this.selectionRing.lineTo(-hw, 0);
      this.selectionRing.closePath();
      this.selectionRing.strokePath();
    }
    this.selectionRing.setVisible(true);
  }

  private hideSelectionRing() {
    if (this.selectionRing) this.selectionRing.setVisible(false);
  }

  public takeDamage(amount: number) {
    this.currentHealth = Math.max(0, this.currentHealth - amount);

    this.scene.tweens.add({
      targets: this.mainSprite,
      alpha: 0.4,
      duration: 50,
      yoyo: true,
      repeat: 1,
    });

    if (this.currentHealth <= 0) {
      this.handleDestruction();
    }
  }

  private handleDestruction() {
    this.destroy();
  }
}
