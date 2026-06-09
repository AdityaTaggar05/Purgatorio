import Phaser from 'phaser';
import type { BuildingData } from '../../../types/building';
import { IsoMath } from '../managers/IsoMath';

export class BuildingSprite extends Phaser.GameObjects.Container {
  public buildingData: BuildingData;
  private mainSprite: Phaser.GameObjects.Sprite;
  private healthBar!: Phaser.GameObjects.Graphics;
  private selectionRing!: Phaser.GameObjects.Graphics;

  public currentHealth: number;
  public maxHealth: number;

  private debugGlow!: Phaser.GameObjects.Graphics;

  constructor(scene: Phaser.Scene, gridFactor: number, x: number, y: number, data: BuildingData) {
    super(scene, x, y);
    this.buildingData = data;

    this.maxHealth = data.size * 500;
    this.currentHealth = this.maxHealth;

    const spriteKey = `building_${data.id.split("_")[0]}`

    this.debugGlow = scene.add.graphics();
    this.debugGlow.setPosition(x, y);
    this.debugGlow.lineStyle(3, 0x00ff00, 1.0); // Bright neon green debug outline
    this.debugGlow.fillStyle(0x00ff00, 0.2);

    const w = IsoMath.TILE_W / gridFactor * data.size;
    const h = IsoMath.TILE_H / gridFactor * data.size;

    this.debugGlow.beginPath();
    this.debugGlow.moveTo(0, -h / 2);
    this.debugGlow.lineTo(w / 2, 0);
    this.debugGlow.lineTo(0, h / 2);
    this.debugGlow.lineTo(-w / 2, 0);
    this.debugGlow.closePath();
    this.debugGlow.strokePath();
    this.debugGlow.fillPath();

    scene.add.existing(this);

    try {
      this.mainSprite = scene.add.sprite(0, 0, spriteKey);
      this.mainSprite.setOrigin(0.5, 1.0);
      this.add(this.mainSprite);
    } catch (_) {
      console.warn(`Asset binding failed for frame: ${spriteKey}. Using debug fallback.`);
    }

    this.setupInteractions();

    scene.add.existing(this);
  }

  private setupInteractions() {
    const width = this.mainSprite.width || this.buildingData.size * 100;
    const height = this.mainSprite.height || this.buildingData.size * 150;

    this.setInteractive(
      new Phaser.Geom.Rectangle(-width / 2, -height, width, height),
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
  }

  private showSelectionRing() {
    if (!this.selectionRing) {
      this.selectionRing = this.scene.add.graphics();
      this.selectionRing.lineStyle(2, 0xdc2626, 0.8); // Blood red outline
      this.selectionRing.strokeRect(-40, -10, 80, 20);
      this.addAt(this.selectionRing, 0); // Inject beneath the main texture
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
