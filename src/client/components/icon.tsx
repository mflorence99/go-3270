import { LitElement } from 'lit';
import { TemplateResult } from 'lit';

import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { html } from 'lit';
import { property } from 'lit/decorators.js';

declare global {
  interface HTMLElementTagNameMap {
    'app-icon': Icon;
  }
}

// 📘 display material icon

// 👉 https://marella.me/material-icons/demo/

//  --app-icon-color    any color, default: inherit
//  --app-icon-filter   any filter, default: none
//  --app-icon-size     any size, default: 1em

@customElement('app-icon')
export class Icon extends LitElement {
  static override styles = [
    css`
      :host {
        display: inline-block;
        text-align: center;
        vertical-align: middle;
      }

      .material-icon {
        color: var(--app-icon-color, var(--md-sys-on-background));
        direction: ltr;
        display: inline-block;
        filter: var(--app-icon-filter, none);
        font-family: Material Icons;
        font-feature-settings: 'liga';
        font-size: var(--app-icon-size, 1em);
        font-style: normal;
        font-weight: normal;
        letter-spacing: normal;
        line-height: 1;
        text-rendering: optimizeLegibility;
        text-transform: none;
        white-space: nowrap;
        word-wrap: normal;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
      }
    `
  ];

  @property() icon!: string;

  override render(): TemplateResult {
    return html`
      <i class="material-icon">${this.icon}</i>
    `;
  }
}
