import { Colors } from '$client/types/3270';
import { DataStreamEventDetail } from '$client/pages/connector';
import { Dimensions } from '$client/types/3270';
import { EmulatorContext } from '$client/services/lu3270';
import { Emulators } from '$client/types/3270';
import { LitElement } from 'lit';
import { Lu3270 } from '$client/services/lu3270';
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

  lu3270: Lu3270 | null = null;

  datastream(e: CustomEvent<DataStreamEventDetail>): void {
    const ectx = this.prepare();
    if (ectx) this.lu3270?.outbound(e.detail.bytes);
  }

  prepare(): EmulatorContext | null {
    const fontSpec = `${this.state.model.get().fontSize.actual}px Terminal`;
    const ctx = this.terminal.getContext('2d');
    if (ctx) {
      // 👇 resize canvas appropriate to font size and dimensions
      ctx.font = fontSpec;
      const metrics = ctx.measureText('A');
      const color =
        Colors[this.state.model.get().config.color] ?? defaultColor;
      const dims: [number, number] =
        Dimensions[this.state.model.get().config.emulator] ??
        defaultDimensions;
      const fontWidth = metrics.width;
      const fontHeight =
        metrics.fontBoundingBoxAscent + metrics.fontBoundingBoxDescent;
      const ectx = {
        color,
        ctx,
        dims,
        fontHeight,
        fontSpec,
        fontWidth,
        responder: this.responder.bind(this)
      };
      // 👇 I *think* we only shoukld do this on a delta
      const cx = dims[0] * fontWidth;
      const cy = dims[1] * fontHeight;
      if (
        cx !== this.terminal.offsetWidth ||
        cy !== this.terminal.offsetHeight
      ) {
        this.terminal.width = cx;
        this.terminal.height = cy;
        // 👇 in any event, we must do this each time the
        //    emulator changes
        this.lu3270?.close();
        this.lu3270 = Lu3270.lu3270(ectx);
      }
      ctx.clearRect(0, 0, this.terminal.width, this.terminal.height);
      return ectx;
    } else return null;
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

  responder(bytes: Uint8Array): void {
    this.dispatchEvent(
      new CustomEvent<DataStreamEventDetail>('response', {
        detail: { bytes }
      })
    );
  }

  override updated(): void {
    const ectx = this.prepare();
    if (ectx) this.lu3270?.refresh();
  }
}
