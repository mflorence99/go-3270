import { Go3270 } from '$client/types/go3270';
import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
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

// 🟧 Emulate the 3270 emulator

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

  @query('#dpi') dpi!: HTMLDivElement;
  @consume({ context: stateContext }) state!: State;
  @query('#terminal') terminal!: HTMLCanvasElement;

  #go3270: Go3270 | null = null;

  confirm(): void {
    if (this.state.model.get().config.screenshot)
      this.dispatchEvent(new CustomEvent('done'));
    else if (
      window.confirm(
        'Are you sure you want to terminate the 3270 session? You may want to logoff from any open applications before doing so.'
      )
    )
      window.dispatchEvent(new CustomEvent('disconnect'));
  }

  disconnect(): void {
    this.#go3270?.close();
  }

  focussed(focus: boolean): void {
    this.#go3270?.focus(focus);
  }

  keystroke(
    code: string,
    key: string,
    alt: boolean,
    ctrl: boolean,
    shift: boolean
  ): void {
    this.#go3270?.keystroke(code, key, alt, ctrl, shift);
  }

  outbound(chars: Uint8ClampedArray): void {
    this.#go3270?.outbound(chars);
  }

  override render(): TemplateResult {
    return html`
      <main class="stretcher">
        <section class="emulator">
          <header class="header">
            <md-icon-button
              @click=${(): void => this.confirm()}
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
            <canvas class="terminal" id="terminal"></canvas>
          </article>

          <footer
            class="status"
            style=${styleMap({
              'color': `${this.state.model.get().config.device === '3278' ? this.state.model.get().clut[0xf4]![0] : this.state.model.get().clut[0xff]![0]}`,
              'font-size': `${this.state.model.get().config.fontSize}`
            })}>
            <article class="left">
              <app-icon icon="computer">
                ${this.state.model.get().config.device}
              </app-icon>

              <app-icon
                icon="hourglass_empty"
                style=${styleMap({
                  visibility: `${this.state.model.get().status.waiting ? 'visible' : 'hidden'}`
                })}>
                WAIT
              </app-icon>

              <app-icon
                icon="clear"
                style=${styleMap({
                  visibility: `${this.state.model.get().status.error ? 'visible' : 'hidden'}`
                })}>
                ${this.state.model.get().status.message}
              </app-icon>
            </article>

            <article class="right">
              <p
                style=${styleMap({
                  visibility: `${this.state.model.get().status.numeric ? 'visible' : 'hidden'}`
                })}>
                NUM
              </p>

              <p
                style=${styleMap({
                  visibility: `${this.state.model.get().status.protected ? 'visible' : 'hidden'}`
                })}>
                PROT
              </p>

              <p
                style=${styleMap({
                  visibility: `${this.state.model.get().status.cursorAt >= 0 ? 'visible' : 'hidden'}`
                })}>
                ${this.state.cursorAt.get()}
              </p>
            </article>
          </footer>
        </section>
      </main>

      <div class="dpi" id="dpi"></div>
    `;
  }

  // 👇 we rebuild the device emulator as the config changes
  override updated(): void {
    if (this.state.delta.config) {
      // 👇 close any prior handler
      this.#go3270?.close();
      // 👇 construct a new device with its new attributes
      const model = this.state.model.get();
      const bgColor = model.clut[0xf0]![0];
      const dpi = this.dpi.offsetWidth * window.devicePixelRatio;
      const fontSize = Math.round(
        // TODO 🔥 Go "gg" seems to interpret font size differently
        Number(model.config.fontSize) * 0.725
      );
      const monochrome = model.config.device === '3278';
      // 👇 construct a new device with its new attributes
      this.#go3270 = window.NewGo3270?.(
        this.terminal,
        bgColor,
        monochrome,
        model.clut,
        fontSize,
        model.config.dims[0],
        model.config.dims[1],
        dpi,
        model.config.screenshot
      );
    }
  }
}
