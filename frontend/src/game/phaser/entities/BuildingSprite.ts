import Phaser from 'phaser';
import type { BuildingData } from '../../../types/building';
import { IsoMath } from '../managers/IsoMath';

const SPRITE_ORIGIN = { x: 0.5, y: 1.0 };

export class BuildingSprite extends Phaser.GameObjects.Container {
  public buildingData: BuildingData;
  private mainSprite: Phaser.GameObjects.Sprite;
  private healthBar!: Phaser.GameObjects.Graphics;
  private selectionRing!: Phaser.GameObjects.Graphics;
  private size: number;

  public currentHealth: number;
  public maxHealth: number;

  constructor(scene: Phaser.Scene, gridFactor: number, x: number, y: number, data: BuildingData) {
    super(scene, x, y);
    this.buildingData = data;

    this.maxHealth = data.size * 500;
    this.currentHealth = this.maxHealth;
    this.size = data.size;

    const spriteKey = `building_${data.id.split("_")[0]}`
    const w = IsoMath.TILE_W / gridFactor * data.size;
    const h = IsoMath.TILE_H / gridFactor * data.size;

    try {
      this.mainSprite = scene.add.sprite(0, 0, spriteKey);

      this.mainSprite.setOrigin(SPRITE_ORIGIN.x, SPRITE_ORIGIN.y);

      const scale = w / this.mainSprite.width;
      this.mainSprite.setScale(scale);

      this.mainSprite.setPosition(0, h / 2);

      this.add(this.mainSprite);
    } catch (_) {
      console.warn(`Asset binding failed for frame: ${spriteKey}. Using debug fallback.`);
    }

    this.setupInteractions(gridFactor);

    scene.add.existing(this);
  }

  private setupInteractions(gridFactor: number) {
    const width = this.mainSprite.width || this.buildingData.size * 100;
    const height = this.mainSprite.height || this.buildingData.size * 150;

    this.setInteractive(
      new Phaser.Geom.Rectangle(-width / 2, -height, width, height),
      Phaser.Geom.Rectangle.Contains
    );

    this.on('pointerover', () => {
      this.scene.input.setDefaultCursor('pointer');
      this.showSelectionRing(gridFactor);
    });

    this.on('pointerout', () => {
      this.scene.input.setDefaultCursor('default');
      this.hideSelectionRing();
    });
  }

  private showSelectionRing(gridFactor: number) {
    if (!this.selectionRing) {
      this.selectionRing = this.scene.add.graphics();
      this.selectionRing.setPosition(this.x, this.y);
      this.selectionRing.lineStyle(3, 0x00ff00, 1.0); // Bright neon green debug outline
      this.selectionRing.fillStyle(0x00ff00, 0.2);

      const w = IsoMath.TILE_W / gridFactor * this.size
      const h = IsoMath.TILE_H / gridFactor * this.size

      this.selectionRing.beginPath();
      this.selectionRing.moveTo(0, -h / 2);
      this.selectionRing.lineTo(w / 2, 0);
      this.selectionRing.lineTo(0, h / 2);
      this.selectionRing.lineTo(-w / 2, 0);
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

    // Structural damage hit-flash animation
    this.scene.tweens.add({
      targets: this.mainSprite,
      alpha: 0.4,
      duration: 50,
      yoyo: true,
      repeat: 1
    });

    if (this.currentHealth <= 0) {
      this.handleDestruction();
    }
  }

  private handleDestruction() {
    this.destroy();
  }
}
