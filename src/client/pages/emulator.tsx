import { Colors } from '$client/state/constants';
import { Dimensions } from '$client/state/constants';
import { Emulators } from '$client/state/constants';
import { Go3270 } from '$client/types/go3270';
import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { defaultColor } from '$client/state/constants';
import { defaultDimensions } from '$client/state/constants';
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
      .dpi {
        height: 1in;
        left: -100%;
        position: absolute;
        top: -100%;
        width: 1in;
      }

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

          .wrapper {
            overflow: hidden;

            .terminal {
              scale: 1;
              transform-origin: left top;
            }
          }
        }
      }
    `
  ];

  @query('.dpi') dpi!: HTMLDivElement;
  @consume({ context: stateContext }) state!: State;
  @query('.terminal') terminal!: HTMLCanvasElement;

  go3270: Go3270 | null = null;

  disconnect(): void {
    this.go3270?.close();
  }

  keystroke(
    code: string,
    key: string,
    alt: boolean,
    ctrl: boolean,
    shift: boolean
  ): void {
    this.go3270?.keystroke(code, key, alt, ctrl, shift);
  }

  receiveFromApp(bytes: Uint8ClampedArray): void {
    this.go3270?.receiveFromApp(bytes);
  }

  override render(): TemplateResult {
    return html`
      <main class="stretcher">
        <section class="emulator">
          <header class="header">
            <md-icon-button
              @click=${(): any =>
                window.dispatchEvent(new CustomEvent('disconnect'))}
              title="Disconnect from 3270">
              <app-icon icon="power_settings_new"></app-icon>
            </md-icon-button>

            <article class="controls">
              <md-icon-button title="Get help">
                <app-icon icon="help"></app-icon>
              </md-icon-button>
            </article>
          </header>

          <article class="wrapper">
            <canvas class="terminal"></canvas>
          </article>

          <footer
            class="status"
            style=${styleMap({
              'color': `${Colors[this.state.model.get().config.color]}`,
              'font-size': `${this.state.model.get().config.fontSize}`
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

      <div class="dpi"></div>
    `;
  }

  override updated(): void {
    // 👇 close any prior handler
    this.go3270?.close();
    // 👇 construct a new device with its new attributes
    const color =
      Colors[this.state.model.get().config.color] ?? defaultColor;
    const dims: [number, number] =
      Dimensions[this.state.model.get().config.emulator] ??
      defaultDimensions;
    const dpi = this.dpi.offsetWidth * window.devicePixelRatio;
    const fontSize = Number(this.state.model.get().config.fontSize);
    // 👇 construct a new device with its new attributes
    this.go3270 = window.NewGo3270?.(
      this.terminal,
      color,
      fontSize,
      dims[0],
      dims[1],
      dpi
    );
  }
}
