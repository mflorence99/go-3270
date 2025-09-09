import myConfig from 'eslint-config-mflorence99';

export default [
  ...myConfig,
  {
    languageOptions: {
      parserOptions: {
        project: [
          'src/builder/tsconfig.json',
          'src/client/tsconfig.json'
        ]
      }
    }
  },
  {
    ignores: ['eslint.config.mjs']
  }
];
