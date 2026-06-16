import Phaser from 'phaser';
import type { PlacedBuilding } from '../../../types/building';
import { IsoMath } from '../managers/IsoMath';
import { phaserEvents } from '../events';

const SPRITE_ORIGIN = { x: 0.5, y: 1.0 };
const GF = IsoMath.SUBDIVISIONS;

export class BuildingSprite extends Phaser.GameObjects.Container {
  public buildingData: PlacedBuilding;
  private mainSprite!: Phaser.GameObjects.Sprite;
  private hoverRing!: Phaser.GameObjects.Graphics;
  private selectedTile!: Phaser.GameObjects.Graphics;
  private upgradeIndicator!: Phaser.GameObjects.Graphics;
  private upgradeTween: Phaser.Tweens.Tween | null = null;
  private spriteW: number;
  private spriteH: number;
  private _selected = false;

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

    this.createRings(scene);

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

  private createRings(scene: Phaser.Scene) {
    const hw = this.spriteW / 2;
    const hh = this.spriteH / 2;

    this.hoverRing = scene.add.graphics();
    this.hoverRing.setPosition(this.x, this.y);
    this.hoverRing.lineStyle(3, 0x00ff00, 1.0);
    this.hoverRing.beginPath();
    this.hoverRing.moveTo(0, -hh);
    this.hoverRing.lineTo(hw, 0);
    this.hoverRing.lineTo(0, hh);
    this.hoverRing.lineTo(-hw, 0);
    this.hoverRing.closePath();
    this.hoverRing.strokePath();
    this.hoverRing.setVisible(false);

    this.selectedTile = scene.add.graphics();
    this.selectedTile.setPosition(this.x, this.y);
    this.selectedTile.fillStyle(0x00ff00, 0.25);
    this.selectedTile.lineStyle(3, 0x00ff00, 0.6);
    this.selectedTile.beginPath();
    this.selectedTile.moveTo(0, -hh);
    this.selectedTile.lineTo(hw, 0);
    this.selectedTile.lineTo(0, hh);
    this.selectedTile.lineTo(-hw, 0);
    this.selectedTile.closePath();
    this.selectedTile.fillPath();
    this.selectedTile.strokePath();
    this.selectedTile.setVisible(false);

    this.upgradeIndicator = scene.add.graphics();
    this.upgradeIndicator.setPosition(this.x, this.y - hh - 4);
    this.upgradeIndicator.fillStyle(0x38bdf8, 0.9);
    this.upgradeIndicator.fillCircle(0, 0, 5);
    this.upgradeIndicator.setVisible(false);

    this.syncUpgradeIndicator();
  }

  set selected(value: boolean) {
    this._selected = value;
    this.selectedTile.setVisible(value);
  }

  get selected(): boolean {
    return this._selected;
  }

  private syncUpgradeIndicator() {
    const upgrading = !!this.buildingData.metadata?.upgrade_ends_at;
    this.upgradeIndicator.setVisible(upgrading);

    if (upgrading && !this.upgradeTween) {
      this.upgradeTween = this.scene.tweens.add({
        targets: this.upgradeIndicator,
        alpha: 0.3,
        duration: 800,
        yoyo: true,
        repeat: -1,
        ease: "Sine.easeInOut",
      });
    } else if (!upgrading && this.upgradeTween) {
      this.upgradeTween.destroy();
      this.upgradeTween = null;
      this.upgradeIndicator.setAlpha(1);
    }
  }

  public updateData(data: PlacedBuilding) {
    this.buildingData = data;
    this.syncUpgradeIndicator();
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
      if (!this._selected) this.hoverRing.setVisible(true);
    });

    this.on('pointerout', () => {
      this.scene.input.setDefaultCursor('default');
      this.hoverRing.setVisible(false);
    });

    this.on('pointerdown', () => {
      if (phaserEvents.mode === "none" && phaserEvents.onBuildingClick) {
        phaserEvents.onBuildingClick(this.buildingData);
      }
    });
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
      this.destroy();
    }
  }

  destroy(fromScene?: boolean) {
    if (this.upgradeTween) {
      this.upgradeTween.destroy();
      this.upgradeTween = null;
    }
    this.hoverRing?.destroy();
    this.selectedTile?.destroy();
    this.upgradeIndicator?.destroy();
    super.destroy(fromScene);
  }
}
