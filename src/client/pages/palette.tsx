import { Config } from '$client/state/state';
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

declare global {
  interface HTMLElementTagNameMap {
    'app-palette': Palette;
  }
}

// 📘 define palette for 3270 emulation

@customElement('app-palette')
export class Palette extends SignalWatcher(LitElement) {
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

          .preview,
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

            .controls {
              display: flex;
              flex-direction: column;
              justify-content: center;
            }

            .instructions {
              font-weight: bold;
            }
          }
        }
      }
    `
  ];

  @consume({ context: stateContext }) state!: State;

  config(evt: Event): void {
    evt.preventDefault();
    const form = evt.target as HTMLFormElement;
    if (form) {
      const formData = new FormData(form);
      const config = Object.fromEntries(formData.entries()) as Config;
      this.state.updateConfig(config);
    }
  }

  done(): void {
    this.dispatchEvent(new CustomEvent('done'));
  }

  override render(): TemplateResult {
    return html`
      <main class="stretcher">
        <header class="header">Customize 3270 Appearance</header>
        <form @submit=${this.config} name="config">
          <section class="palette">
            <article class="settings">
              <div
                style="background: pink; width: 400px; height:400px"></div>

              <div class="controls">
                <p class="instructions">Select Font Size</p>

                <md-filled-select name="fontSize">
                  ${repeat(
                    ['6', '7', '8', '9', '10', '11', '12', '13', '14'],
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

            <article class="preview">
              <div
                style="background: coral; width: 400px; height:400px"></div>

              <div class="buttons">
                <md-filled-button>Save</md-filled-button>
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
}
