import { Colors } from '$client/types/3270';
import { DataStreamEventDetail } from '$client/pages/connector';
import { Dimensions } from '$client/types/3270';
import { Emulators } from '$client/types/3270';
import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { defaultColor } from '$client/types/3270';
import { defaultDimensions } from '$client/types/3270';
import { globals } from '$client/css/globals/shadow-dom';
import { html } from 'lit';
import { query } from 'lit/decorators.js';
import { stateContext } from '$client/state/state';
import { styleMap } from 'lit/directives/style-map.js';

declare global {
  interface HTMLElementTagNameMap {
    'app-emulator': Emulator;
  }
}

type EmulatorContext = {
  ctx: CanvasRenderingContext2D;
  dims: [number, number];
  fontHeight: number;
  fontSpec: string;
  fontWidth: number;
};

// 📘 emulate the 3270 emulator

@customElement('app-emulator')
export class Emulator extends SignalWatcher(LitElement) {
  static override styles = [
    globals,
    css`
      .stretcher {
        align-items: center;
        display: flex;
        flex-direction: column;
        height: 100%;
        justify-content: center;

        .emulator {
          display: flex;
          flex-direction: column;
          justify-content: stretch;

          .header {
            align-items: center;
            border-bottom: 1px solid currentColor;
            display: flex;
            flex-direction: row;
            gap: 1rem;
            justify-content: space-between;
            margin-bottom: 0.25rem;

            .controls {
              display: flex;
              flex-direction: row;
              justify-content: flex-end;

              --app-icon-color: var(--md-sys-color-primary);
            }

            .title {
              font-size: 1.5rem;
              font-weight: bold;
            }
          }

          .status {
            border-top: 1px solid currentColor;
            display: flex;
            flex-direction: row;
            font-family: Terminal;
            justify-content: space-between;
            margin-top: 0.5rem;

            .left,
            .right {
              display: flex;
              flex-direction: row;
              gap: 1rem;
            }
          }
        }
      }
    `
  ];

  @consume({ context: stateContext }) state!: State;
  @query('.terminal') terminal!: HTMLCanvasElement;

  datastream(e: CustomEvent<DataStreamEventDetail>): void {
    const ectx = this.#prepareCanvas();
    if (ectx) this.#renderCanvas(ectx, e.detail.bytes);
  }

  override render(): TemplateResult {
    return html`
      <main class="stretcher">
        <section class="emulator">
          <header class="header">
            <md-icon-button
              @click=${(): any =>
                this.dispatchEvent(new CustomEvent('disconnect'))}
              title="Disconnect from 3270">
              <app-icon icon="power_settings_new"></app-icon>
            </md-icon-button>

            <article class="controls">
              <md-icon-button
                @click=${(): void => this.state.increaseFontSize()}
                ?disabled=${this.state.model.get().fontSize.actual >=
                this.state.model.get().fontSize.max}
                title="Increase text size">
                <app-icon icon="text_increase"></app-icon>
              </md-icon-button>

              <md-icon-button
                @click=${(): void => this.state.decreaseFontSize()}
                ?disabled=${this.state.model.get().fontSize.actual <=
                this.state.model.get().fontSize.min}
                title="Decrease text size">
                <app-icon icon="text_decrease"></app-icon>
              </md-icon-button>

              <md-icon-button title="Get help">
                <app-icon icon="help"></app-icon>
              </md-icon-button>
            </article>
          </header>

          <canvas class="terminal"></canvas>

          <footer
            class="status"
            style=${styleMap({
              color: `${Colors[this.state.model.get().config.color]}`
            })}>
            <article class="left">
              <app-icon icon="computer">
                ${Emulators[this.state.model.get().config.emulator]}
              </app-icon>
              <app-icon icon="access_time">WAIT</app-icon>
              <app-icon icon="clear">MSG</app-icon>
            </article>

            <article class="right">
              <p>NUM</p>
              <p>PROT</p>
              <p>001/001</p>
            </article>
          </footer>
        </section>
      </main>
    `;
  }

  override updated(): void {
    const ectx = this.#prepareCanvas();
    if (ectx) this.#renderCanvas(ectx);
  }

  #prepareCanvas(): EmulatorContext | null {
    const fontSpec = `${this.state.model.get().fontSize.actual}px Terminal`;
    const ctx = this.terminal.getContext('2d');
    if (ctx) {
      // 👇 resize canvas appropriate to font size and dimensions
      ctx.font = fontSpec;
      const metrics = ctx.measureText('A');
      const dims: [number, number] =
        Dimensions[this.state.model.get().config.emulator] ??
        defaultDimensions;
      const fontWidth = metrics.width;
      const fontHeight =
        metrics.fontBoundingBoxAscent + metrics.fontBoundingBoxDescent;
      const cx = dims[0] * fontWidth;
      const cy = dims[1] * fontHeight;
      if (
        cx !== this.terminal.offsetWidth ||
        cy !== this.terminal.offsetHeight
      ) {
        this.terminal.width = cx;
        this.terminal.height = cy;
      }
      ctx.clearRect(0, 0, this.terminal.width, this.terminal.height);
      return { ctx, dims, fontHeight, fontSpec, fontWidth };
    } else return null;
  }

  // 🔥 TEMPORARY - draw random characters to fill screen
  // 🔥 don't know why we need this
  // 🔥 specifying no-unused-vars doesn't work!
  // eslint-disable-next-line
  #renderCanvas(ectx: EmulatorContext, bytes?: Uint8Array): void {
    const { ctx, dims, fontHeight, fontSpec, fontWidth } = ectx;
    // 👇 establish terminal font and color
    ctx.font = fontSpec;
    ctx.textAlign = 'left';
    ctx.textBaseline = 'top';
    ctx.fillStyle =
      Colors[this.state.model.get().config.color] ?? defaultColor;
    // 👇 will do something with "bytes" or refresh if null

    const chars =
      'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+[{]};:,<.>/?';
    for (let ix = 0, x = 0; ix < dims[0]; ix++, x += fontWidth) {
      for (let iy = 0, y = 0; iy < dims[1]; iy++, y += fontHeight) {
        ctx.fillText(
          chars.charAt(Math.floor(Math.random() * chars.length)),
          x,
          y
        );
      }
    }
  }
}
