import { css } from 'lit';

// 🟦 global styles injected into shadow DOM components

export const globals = css`
  *,
  *::before,
  *::after {
    box-sizing: border-box;
  }

  a {
    color: var(--md-sys-color-primary);
    cursor: pointer;
    text-decoration: none;

    &:hover {
      color: var(--md-sys-color-tertiary);
      text-decoration: none;
    }
  }

  h1,
  h2,
  h3,
  h4,
  h5,
  h6 {
    margin: 0;
  }

  hr {
    margin: 1rem 0;
    width: 100%;
  }

  ol,
  ul {
    list-style-type: none;
  }

  p {
    margin: 0;
  }
`;
