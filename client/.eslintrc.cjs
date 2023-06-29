/* eslint-env node */
require('@rushstack/eslint-patch/modern-module-resolution')
const prettierConfig = require('./.prettierrc.json')

module.exports = {
  root: true,
  extends: [
    'plugin:vue/vue3-recommended',
    'eslint:recommended',
    '@vue/eslint-config-typescript/recommended',
    'prettier',
    '@vue/eslint-config-prettier',
  ],
  env: {
    'vue/setup-compiler-macros': true,
  },
  rules: {
    semi: ['error', 'never'],
    quotes: ['error', 'single', { avoidEscape: true }],
    'vue/multi-word-component-names': 0,
    '@typescript-eslint/no-explicit-any': 1,
    'prettier/prettier': ['error', prettierConfig],
    '@typescript-eslint/no-non-null-assertion': 0,
  },
}
