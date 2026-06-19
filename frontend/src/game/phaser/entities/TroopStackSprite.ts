import Phaser from 'phaser';
import { IsoMath } from '../managers/IsoMath';

const DEPTH_OFFSET = 101;

export interface TroopStackData {
  id: string;
  troopType: string;
  x: number;
  y: number;
  unitMaxHp: number;
  totalUnits: number;
  hp: number;
  alive: boolean;
}

const GF = IsoMath.SUBDIVISIONS;

export class TroopStackSprite extends Phaser.GameObjects.Container {
  public stackData: TroopStackData;
  private icon!: Phaser.GameObjects.Sprite;
  private hpBar!: Phaser.GameObjects.Graphics;
  private countBadgeBg!: Phaser.GameObjects.Graphics;
  private countBadgeText!: Phaser.GameObjects.Text;
  private iconSize: number;
  private staticPreview: boolean;
  private deathTweenStarted = false;

  constructor(scene: Phaser.Scene, data: TroopStackData, staticPreview = false) {
    const pos = IsoMath.subgridToScreen(data.x, data.y, 1);
    super(scene, pos.x, pos.y);
    this.stackData = data;
    this.staticPreview = staticPreview;
    this.iconSize = (IsoMath.TILE_W / GF) * 0.55;

    const iconKey = `troop_${data.troopType}`;
    try {
      this.icon = scene.add.sprite(0, this.iconSize / 2, iconKey);
      this.icon.setOrigin(0.5, 1.0);
      const scale = this.iconSize / Math.max(this.icon.width, this.icon.height);
      this.icon.setScale(scale);
      this.add(this.icon);
    } catch (_) {
      console.warn(`Asset binding failed for troop frame: ${iconKey}.`);
    }

    this.hpBar = scene.add.graphics();
    this.add(this.hpBar);

    this.countBadgeBg = scene.add.graphics();
    this.countBadgeText = scene.add.text(0, 0, '', {
      fontSize: '34px',
      fontStyle: 'bold',
      color: '#e5e7eb',
    });
    this.countBadgeText.setOrigin(0.5, 0.5);
    this.add(this.countBadgeBg);
    this.add(this.countBadgeText);

    scene.add.existing(this);
    this.refresh();
  }

  public update(data: TroopStackData) {
    const wasAlive = this.stackData.alive;
    const prevX = this.stackData.x;
    const prevY = this.stackData.y;
    this.stackData = data;
    this.refresh();

    // Update container position if (x, y) changed (troop movement)
    if (data.x !== prevX || data.y !== prevY) {
      const pos = IsoMath.subgridToScreen(data.x, data.y, 1);
      this.setPosition(pos.x, pos.y);
      this.setDepth(pos.y + DEPTH_OFFSET);
    }

    if (wasAlive && !data.alive && !this.deathTweenStarted) {
      this.deathTweenStarted = true;
      this.scene.tweens.add({
        targets: this,
        alpha: 0,
        scale: 0.5,
        duration: 400,
        ease: 'Sine.easeIn',
      });
    }
  }

  private refresh() {
    const { hp, unitMaxHp, totalUnits } = this.stackData;
    const maxHp = Math.max(1, unitMaxHp * totalUnits);

    if (this.hpBar) {
      const barW = this.iconSize * 1.1;
      const barH = this.iconSize * 0.12;
      const barY = -this.iconSize - barH - 6;
      const pct = Math.max(0, Math.min(1, hp / maxHp));
      const color = pct > 0.6 ? 0x22c55e : pct > 0.3 ? 0xeab308 : 0xef4444;

      this.hpBar.clear();
      this.hpBar.fillStyle(0x000000, 0.7);
      this.hpBar.fillRect(-barW / 2, barY, barW, barH);
      this.hpBar.fillStyle(color, 1);
      this.hpBar.fillRect(-barW / 2, barY, barW * pct, barH);
    }

    const displayCount = this.staticPreview
      ? totalUnits
      : Math.max(0, Math.min(totalUnits, Math.ceil(hp / Math.max(1, unitMaxHp))));

    if (displayCount > 1) {
      const r = this.iconSize * 0.18;
      const bx = this.iconSize * 0.32;
      const by = this.iconSize * 0.12;
      this.countBadgeBg.clear();
      this.countBadgeBg.fillStyle(0x000000, 0.85);
      this.countBadgeBg.fillCircle(bx, by, r);
      this.countBadgeText.setPosition(bx, by);
      this.countBadgeText.setFontSize(this.iconSize * 0.22);
      this.countBadgeText.setText(String(displayCount));
      this.countBadgeText.setVisible(true);
      this.countBadgeBg.setVisible(true);
    } else {
      this.countBadgeText.setVisible(false);
      this.countBadgeBg.setVisible(false);
    }
  }

  destroy(fromScene?: boolean) {
    this.hpBar?.destroy();
    this.countBadgeBg?.destroy();
    this.countBadgeText?.destroy();
    super.destroy(fromScene);
  }
}
