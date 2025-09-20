import { Colors } from '$lib/types/3270';
import { Dimensions } from '$lib/types/3270';
import { Emulators } from '$lib/types/3270';
import { LitElement } from 'lit';
import { Pages } from '$client/pages/root';
import { SignalWatcher } from '@lit-labs/signals';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { defaultColor } from '$lib/types/3270';
import { defaultDimensions } from '$lib/types/3270';
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

  @query('.terminal') terminal!: HTMLCanvasElement;
  @consume({ context: stateContext }) theState!: State;

  override render(): TemplateResult {
    const model = this.theState.model;
    return html`
      <main class="stretcher">
        <section class="emulator">
          <header class="header">
            <md-icon-button
              @click=${(): void => State.theTn3270?.close()}
              title="Disconnect from 3270">
              <app-icon icon="power_settings_new"></app-icon>
            </md-icon-button>

            <article class="controls">
              <md-icon-button
                @click=${(): void => this.theState.increaseFontSize()}
                ?disabled=${model.get().fontSize.actual >=
                model.get().fontSize.max}
                title="Increase text size">
                <app-icon icon="text_increase"></app-icon>
              </md-icon-button>

              <md-icon-button
                @click=${(): void => this.theState.decreaseFontSize()}
                ?disabled=${model.get().fontSize.actual <=
                model.get().fontSize.min}
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
              color: `${Colors[model.get().config.color]}`
            })}>
            <article class="left">
              <app-icon icon="computer">
                ${Emulators[model.get().config.emulator]}
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
    const model = this.theState.model;
    // 👇 only when we transition to the emulator page
    if (model.get().pageNum === Pages.emulator) {
      const fontSpec = `${model.get().fontSize.actual}px Terminal`;
      // 👇 maske sure the font is available
      if (document.fonts.check(fontSpec)) {
        const ctx = this.terminal.getContext('2d');
        if (ctx) {
          // 👇 resize canvas appropraite to font size and dimensions
          ctx.font = fontSpec;
          const metrics = ctx.measureText('A');
          const dims: [number, number] =
            Dimensions[model.get().config.emulator] ??
            defaultDimensions;
          const fontWidth = metrics.width;
          const fontHeight =
            metrics.fontBoundingBoxAscent +
            metrics.fontBoundingBoxDescent;
          this.terminal.width = dims[0] * fontWidth;
          this.terminal.height = dims[1] * fontHeight;
          // 👇 establish terminal font and color
          ctx.font = fontSpec;
          ctx.clearRect(
            0,
            0,
            this.terminal.width,
            this.terminal.height
          );
          ctx.textAlign = 'left';
          ctx.textBaseline = 'top';
          ctx.fillStyle =
            Colors[model.get().config.color] ?? defaultColor;
          // 🔥 TEMPORARY - draw random characters to fill screen
          const chars =
            'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+[{]};:,<.>/?';
          for (let ix = 0, x = 0; ix < dims[0]; ix++, x += fontWidth) {
            for (
              let iy = 0, y = 0;
              iy < dims[1];
              iy++, y += fontHeight
            ) {
              ctx.fillText(
                chars.charAt(Math.floor(Math.random() * chars.length)),
                x,
                y
              );
            }
          }
        }
      }
    }
  }
}
