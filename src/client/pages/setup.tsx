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
import { query } from 'lit/decorators.js';
import { repeat } from 'lit/directives/repeat.js';
import { stateContext } from '$client/state/state';
import { styleMap } from 'lit/directives/style-map.js';

declare global {
  interface HTMLElementTagNameMap {
    'app-setup': Setup;
  }
}

// 📘 327x Setup options

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

        .header {
          font-size: 2rem;
          font-weight: bold;
          text-align: center;
          text-transform: uppercase;
        }

        .palette {
          align-items: stretch;
          display: flex;
          flex-direction: row;
          gap: 2rem;

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
              .monochrome {
                color: #808080;
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
              padding: 1rem;
              td {
                padding: 0;
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
        const color = input.getAttribute('data-color') as string;
        const ix = Number(input.getAttribute('data-color-ix'));
        acc[color]![ix] = input.value;
        return acc;
      }, structuredClone(this.state.model.get().clut));
      this.state.updateCLUT(clut);
    }
  }

  done(): void {
    this.dispatchEvent(new CustomEvent('done'));
  }

  override render(): TemplateResult {
    const clut = this.state.model.get().clut;
    const colors = [
      'blue',
      'red',
      'pink',
      'green',
      'turquoise',
      'yellow',
      'white'
    ];
    const device = this.state.model.get().config.device;
    return html`
      <main class="stretcher">
        <header class="header">Customize ${device} Appearance</header>
        <form @submit=${this.config} name="config">
          <section class="palette">
            <article class="settings">
              <table class="clut">
                <thead>
                  <tr>
                    <td></td>
                    <td>Normal</td>
                    <td>Highlight</td>
                  </tr>
                </thead>

                <tbody>
                  ${repeat(
                    colors,
                    (color) => color,
                    (color) => html`
                      <tr>
                        <td
                          style=${styleMap({
                            color:
                              color !== 'green' && device === '3278'
                                ? '#808080'
                                : 'inherit'
                          })}>
                          ${color}
                        </td>
                        <td>
                          <input
                            data-color=${color}
                            data-color-index="0"
                            type="color"
                            ?disabled=${color !== 'green' &&
                            device === '3278'}
                            .value=${clut[color]![0]} />
                        </td>
                        <td>
                          <input
                            data-color=${color}
                            data-color-ix="1"
                            type="color"
                            ?disabled=${color !== 'green' &&
                            device === '3278'}
                            .value=${clut[color]![1]} />
                        </td>
                      </tr>
                    `
                  )}
                </tbody>
              </table>

              <div class="controls">
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
            </article>

            <article class="settings">
              <table
                class="preview"
                style=${styleMap({
                  'font-size': `${this.state.model.get().config.fontSize}px`
                })}>
                ${repeat(
                  device === '3278' ? ['green'] : colors,
                  (color) => color,
                  (color) => html`
                    <tr
                      style=${styleMap({
                        color: `${clut[color]![0]}`
                      })}>
                      <td>Employee ID :&nbsp;</td>
                      <td>04921</td>
                      <td>&nbsp;&nbsp;&nbsp;&nbsp;</td>
                      <td>Status :&nbsp;</td>
                      <td
                        style=${styleMap({
                          'background-color': `${clut[color]![1]}`,
                          'color': 'black'
                        })}>
                        <b>ACTIVE</b>
                      </td>
                    </tr>
                    <tr
                      style=${styleMap({
                        color: `${clut[color]![0]}`
                      })}>
                      <td
                        style=${styleMap({
                          color: `${clut[color]![1]}`
                        })}>
                        <b>Last Name&nbsp;&nbsp;&nbsp;:&nbsp;</b>
                      </td>
                      <td
                        colspan="3"
                        style=${styleMap({
                          color: `${clut[color]![1]}`
                        })}>
                        Smith
                      </td>
                    </tr>
                  `
                )}
              </table>

              <div class="buttons">
                <md-filled-button>Save</md-filled-button>
                <md-outlined-button
                  @click=${this.restore}
                  type="button">
                  Restore
                </md-outlined-button>
                <md-outlined-button @click=${this.done} type="button">
                  Done
                </md-outlined-button>
              </div>
            </article>
          </section>
        </form>
      </main>
    `;
  }

  restore(evt: Event): void {
    evt.preventDefault();
    this.state.updateCLUT(defaultCLUT);
  }
}
