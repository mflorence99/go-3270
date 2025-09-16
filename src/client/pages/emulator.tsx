import { Colors } from '$lib/types/3270';
import { Dimensions } from '$lib/types/3270';
import { Emulators } from '$lib/types/3270';
import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { globals } from '$client/css/globals/shadow-dom';
import { html } from 'lit';
import { repeat } from 'lit/directives/repeat.js';
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

          .simulation {
            align-items: start;
            display: grid;
            font-family: '3270 Font';
            justify-content: center;

            .cell {
              height: 1em;
              margin: 0.1px;
              width: 1ch;
            }
          }

          .status {
            border-top: 1px solid currentColor;
            display: flex;
            flex-direction: row;
            font-family: '3270 Font';
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

  @consume({ context: stateContext }) theState!: State;

  override render(): TemplateResult {
    const model = this.theState.model;
    const dims: [number, number] = Dimensions[
      model.get().config.emulator
    ] ?? [80, 24];
    return html`
      <main class="stretcher">
        <section class="emulator">
          <header class="header">
            <h1 class="title">
              ${Emulators[model.get().config.emulator]} at
              ${model.get().config.host}:${model.get().config.port}
            </h1>

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

              <md-icon-button title="Show keyboard">
                <app-icon icon="keyboard"></app-icon>
              </md-icon-button>

              <md-icon-button
                @click=${(): void => State.theTn3270?.close()}
                title="Disconnect from 3270">
                <app-icon icon="power_settings_new"></app-icon>
              </md-icon-button>
              <md-icon-button title="Get help">
                <app-icon icon="help"></app-icon>
              </md-icon-button>
            </article>
          </header>

          <article
            class="simulation"
            style=${styleMap({
              'color': `${Colors[model.get().config.color]}`,
              'font-size': `${model.get().fontSize.actual}px`,
              'grid-template-columns': `repeat(${dims[0]}, auto)`,
              'grid-template-rows': `repeat(${dims[1]}, auto)`
            })}>
            ${repeat(
              Array.from(Array(dims[0] * dims[1]).keys()),
              (item) => item,
              () => {
                const chars =
                  'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+[{]};:,<.>/?';
                return html`
                  <div class="cell">
                    ${chars.charAt(
                      Math.floor(Math.random() * chars.length)
                    )}
                  </div>
                `;
              }
            )}
          </article>

          <footer
            class="status"
            style=${styleMap({
              color: `${Colors[model.get().config.color]}`
            })}>
            <article class="left">
              <app-icon icon="computer">
                ${Emulators[model.get().config.emulator]}
              </app-icon>
              <app-icon icon="hourglass_bottom">WAIT</app-icon>
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
}
