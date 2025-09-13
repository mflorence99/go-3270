import { Colors } from '$lib/types/3270';
import { Config } from '$client/state/state';
import { Dimensions } from '$lib/types/3270';
import { LitElement } from 'lit';
import { Models } from '$lib/types/3270';
import { SignalWatcher } from '@lit-labs/signals';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { globals } from '$client/css/globals';
import { html } from 'lit';
import { repeat } from 'lit/directives/repeat.js';
import { stateContext } from '$client/state/state';
import { styleMap } from 'lit/directives/style-map.js';

declare global {
  interface HTMLElementTagNameMap {
    'app-home': Home;
  }
}

// 📘 the whole enchilada

@customElement('app-home')
export class Home extends SignalWatcher(LitElement) {
  static override styles = [
    globals,
    css`
      .home {
        align-items: center;
        display: flex;
        flex-direction: column;
        gap: 2rem;
        height: 100%;
        justify-content: center;
        min-width: 800px;
        width: 100%;

        .header {
          text-align: center;

          .major {
            color: #00cc00;
            font-family: '3270 Font';
            font-size: 12rem;
            font-weight: bold;
            line-height: 1;
            text-shadow: 0.33rem 0.33rem #008000;
          }

          .minor {
            font-size: 2rem;
            text-transform: uppercase;
          }
        }

        .config {
          align-items: stretch;
          display: flex;
          flex-direction: row;
          gap: 1rem;
          justify-content: center;

          .color,
          .connection,
          .emulation {
            display: flex;
            flex-direction: column;
            gap: 0.5rem;
          }

          .color .sample {
            font-family: '3270 Font';
            font-size: larger;
          }

          .connection {
            align-items: center;
            justify-content: space-between;
          }

          .connection .host {
            display: flex;
            flex-direction: row;
            gap: 0.5rem;
          }

          .emulation .dims {
            font-size: smaller;
          }
        }
      }
    `
  ];

  @consume({ context: stateContext }) theState!: State;

  override render(): TemplateResult {
    const model = this.theState.model;
    return html`
      <section class="home">
        <article class="header">
          <header class="major">3270</header>
          <p class="minor">Go-powered 3270 Emulator</p>
        </article>

        <hr />

        <form @submit=${this.#submit} class="config" name="config">
          <article class="connection">
            <div class="host">
              <md-filled-text-field
                label="Hostname or IP"
                name="host"
                style="width: 10rem"
                value=${model.get().config.host}></md-filled-text-field>
              <md-filled-text-field
                label="Port"
                name="port"
                style="width: 5rem"
                value=${model.get().config.port}></md-filled-text-field>
            </div>

            <md-filled-button>Connect</md-filled-button>
          </article>

          <article class="emulation">
            <p>Select 3270 Model to Emulate</p>

            ${repeat(
              Object.entries(Models),
              (emulation) => emulation[0],
              (emulation) => html`
                <label>
                  <md-radio
                    ?checked=${model.get().config.emulation ===
                    emulation[1]}
                    name="emulation"
                    value=${emulation[1]}></md-radio>
                  ${emulation[0]} &mdash;
                  <em class="dims">
                    ${Dimensions[emulation[1]]?.[0]} x
                    ${Dimensions[emulation[1]]?.[1]}
                  </em>
                </label>
              `
            )}
          </article>

          <article class="color">
            <p>Select Default 3270 Color</p>

            ${repeat(
              Object.entries(Colors),
              (color) => color[0],
              (color) => html`
                <label>
                  <md-radio
                    ?checked=${model.get().config.color === color[0]}
                    name="color"
                    value=${color[0]}></md-radio>
                  <span
                    class="sample"
                    style=${styleMap({
                      color: color[1]
                    })}>
                    CUSTOMER NUM: 123456
                  </span>
                </label>
              `
            )}
          </article>
        </form>
      </section>
    `;
  }

  // 👁️ https://dev.to/blikblum/dry-form-handling-with-lit-19f
  #submit(event: Event): void {
    event.preventDefault();
    const form = event.target as HTMLFormElement;
    if (form) {
      const formData = new FormData(form);
      const formValues = Object.fromEntries(formData.entries());
      this.theState.updateConfig(formValues as Config);
    }
  }
}
