const tsPlugin = require('@typescript-eslint/eslint-plugin')
const tsParser = require('@typescript-eslint/parser')
const react = require('eslint-plugin-react')
const reactHooks = require('eslint-plugin-react-hooks')

module.exports = [
  {
    files: ['src/**/*.{js,jsx,ts,tsx}'],
    plugins: {
      '@typescript-eslint': tsPlugin,
      'react': react,
      'react-hooks': reactHooks,
    },
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        ecmaVersion: 2018,
        sourceType: 'module',
        ecmaFeatures: { jsx: true },
      },
    },
    rules: {
      '@typescript-eslint/no-unused-vars': ['error', {
        args: 'none',
        varsIgnorePattern: '^_',
        ignoreRestSiblings: true,
      }],
      'eol-last': ['error', 'always'],
      'no-multiple-empty-lines': ['error', { max: Infinity, maxEOF: 0 }],
      'no-trailing-spaces': ['error'],
      'react/jsx-uses-react': ['error'],
      'react/jsx-uses-vars': ['error'],
      'react-hooks/rules-of-hooks': ['error'],
      'react-hooks/exhaustive-deps': ['warn'],
    },
  },
]
