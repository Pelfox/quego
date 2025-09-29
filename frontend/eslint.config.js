import antfu from '@antfu/eslint-config';

export default antfu({
  formatters: true,
  react: true,
  stylistic: {
    semi: true,
  },
  rules: {
    'react-dom/no-missing-button-type': 'off',
  },
});
