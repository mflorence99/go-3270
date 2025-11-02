import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { defaultCLUT } from '$client/state/state';
import { globals } from '$client/css/globals/shadow-dom';
import { html } from 'lit';
import { ifDefined } from 'lit/directives/if-defined.js';
import { query } from 'lit/decorators.js';
import { repeat } from 'lit/directives/repeat.js';
import { stateContext } from '$client/state/state';
import { styleMap } from 'lit/directives/style-map.js';

declare global {
  interface HTMLElementTagNameMap {
    'app-setup': Setup;
  }
}

// 🟧 327x Setup options

@customElement('app-setup')
export class Setup extends SignalWatcher(LitElement) {
  static override styles = [
    globals,
    css`
      .stretcher {
        align-items: center;
        display: flex;
        flex-direction: column;
        gap: 2rem;
        height: 100%;
        justify-content: center;
        min-width: 800px /* 👈 to prevent wrapping */;

        .setup {
          display: flex;
          flex-direction: column;
          gap: 2rem;
          justify-content: stretch;

          .header {
            font-size: 2rem;
            font-weight: bold;
            text-align: center;
            text-transform: uppercase;
          }

          .config {
            align-items: stretch;
            display: flex;
            flex-direction: row;
            gap: 1rem;
            justify-content: center;

            .settings {
              display: flex;
              flex-direction: column;
              gap: 1rem;
              justify-content: space-between;

              .buttons {
                display: flex;
                flex-direction: row;
                gap: 1rem;
                justify-content: center;
              }

              .clut {
                input[type='color'] {
                  height: 2rem;
                  width: 5rem;
                }
                td {
                  font-weight: bold;
                  padding: 0.25rem;
                  text-transform: capitalize;
                }
              }

              .controls {
                display: flex;
                flex-direction: column;
                justify-content: center;
              }

              .instructions {
                font-weight: bold;
              }

              .preview {
                border: 1px solid var(--md-sys-color-on-surface);
                font-family: Terminal;
                letter-spacing: 0.125ch;
                padding: 0.5rem;
                td {
                  padding: 0;
                }
              }
            }
          }
        }
      }
    `
  ];

  @query('.clut') clut!: HTMLElement;
  @consume({ context: stateContext }) state!: State;

  config(evt: Event): void {
    evt.preventDefault();
    const form = evt.target as HTMLFormElement;
    if (form) {
      const formData = new FormData(form);
      const config = Object.fromEntries(formData.entries());
      this.state.updateConfig(config);
      // 👇 now extract all the colors from the clut
      const inputs: HTMLInputElement[] = Array.from(
        this.clut.querySelectorAll('input[type=color]')
      );
      const clut = inputs.reduce((acc, input) => {
        const color = Number(input.getAttribute('data-color'));
        acc[color]![0] = input.value;
        return acc;
      }, structuredClone(this.state.model.get().clut));
      this.state.updateCLUT(clut);
    }
  }

  override render(): TemplateResult {
    const clut = this.state.model.get().clut;
    const colors = [
      [0xf7, 0xf0], // 👈 FG, BG
      [0xf1, 0xf9], // 👈 blue, deep blue
      [0xf2, 0xfa], // 👈 red, orange
      [0xf3, 0xfb], // 👈 pink, purple
      [0xf4, 0xfc], // 👈 green, pale green
      [0xf5, 0xfd], // 👈 turquoise, pale turquoise
      [0xf6, 0xfe], // 👈 yellow, gray
      [0xf8, 0xff] // 👈 black, white
    ];
    const colorOf = (color: any): string => clut[color]![0];
    const nameOf = (color: any): string => clut[color]![1];
    const device = this.state.model.get().config.device;
    return html`
      <main class="stretcher">
        <section class="setup">
          <header class="header">Customize ${device} Appearance</header>

          <hr />

          <form @submit=${this.config} class="config" name="config">
            <article class="settings">
              <table class="clut">
                <tbody>
                  ${repeat(
                    colors,
                    (color) => color,
                    (color) => html`
                      <tr>
                        <td
                          style=${styleMap({
                            'color':
                              device === '3278' && color[0] !== 0xf4
                                ? '#808080'
                                : 'inherit',
                            'text-align': 'right'
                          })}>
                          ${nameOf(color[0])}
                        </td>
                        <td>
                          <input
                            data-color=${ifDefined(color[0])}
                            type="color"
                            ?disabled=${device === '3278' &&
                            color[0] !== 0xf4}
                            .value=${colorOf(color[0])} />
                        </td>
                        <td>
                          <input
                            data-color=${ifDefined(color[1])}
                            type="color"
                            ?disabled=${device === '3278' &&
                            color[1] !== 0xf4}
                            .value=${colorOf(color[1])} />
                        </td>
                        <td
                          style=${styleMap({
                            color:
                              device === '3278' && color[1] !== 0xf4
                                ? '#808080'
                                : 'inherit'
                          })}>
                          ${nameOf(color[1])}
                        </td>
                      </tr>
                    `
                  )}
                </tbody>
              </table>
            </article>

            <article class="settings">
              <div class="controls">
                <p class="instructions">Sample Fields</p>

                <div
                  class="preview"
                  style=${styleMap({
                    'font-size': `${this.state.model.get().config.fontSize}px`
                  })}>
                  <p style=${styleMap({ color: colorOf(0xf4) })}>
                    Unprotected field
                  </p>

                  <p
                    style=${styleMap({
                      color:
                        device === '3278'
                          ? colorOf(0xf4)
                          : colorOf(0xf2)
                    })}>
                    <b>Unprotected field - highlighted</b>
                  </p>

                  <p
                    style=${styleMap({
                      color:
                        device === '3278'
                          ? colorOf(0xf4)
                          : colorOf(0xf1)
                    })}>
                    Protected field
                  </p>

                  <p
                    style=${styleMap({
                      color:
                        device === '3278'
                          ? colorOf(0xf4)
                          : colorOf(0xf7)
                    })}>
                    <b>Protected field - highlighted</b>
                  </p>
                </div>

                <br />

                <p class="instructions">Select Font Size</p>

                <md-filled-select name="fontSize">
                  ${repeat(
                    [
                      '10',
                      '11',
                      '12',
                      '13',
                      '14',
                      '15',
                      '16',
                      '17',
                      '18',
                      '19',
                      '20'
                    ],
                    (fontSize) => fontSize,
                    (fontSize) => html`
                      <md-select-option
                        ?selected=${this.state.model.get().config
                          .fontSize === fontSize}
                        value=${fontSize}>
                        <div slot="headline">${fontSize}px</div>
                      </md-select-option>
                    `
                  )}
                </md-filled-select>
              </div>

              <div class="buttons">
                <md-filled-button>Save</md-filled-button>
                <md-outlined-button
                  @click=${this.restore}
                  type="button">
                  Restore
                </md-outlined-button>
                <md-outlined-button
                  @click=${(): any =>
                    this.dispatchEvent(new CustomEvent('done'))}
                  type="button">
                  Done
                </md-outlined-button>
              </div>
            </article>
          </form>
        </section>
      </main>
    `;
  }

  restore(evt: Event): void {
    evt.preventDefault();
    this.state.updateCLUT(defaultCLUT);
  }
}
