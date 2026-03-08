const tsPlugin = require('@typescript-eslint/eslint-plugin')
const tsParser = require('@typescript-eslint/parser')
const react = require('eslint-plugin-react')
const reactHooks = require('eslint-plugin-react-hooks')
const prettierConfig = require('eslint-config-prettier')

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
        ecmaVersion: 2022,
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
      '@typescript-eslint/no-explicit-any': 'error',
      '@typescript-eslint/consistent-type-imports': ['warn', {
        prefer: 'type-imports',
        fixStyle: 'inline-type-imports',
      }],
      'react/jsx-uses-react': ['error'],
      'react/jsx-uses-vars': ['error'],
      'react-hooks/rules-of-hooks': ['error'],
      'react-hooks/exhaustive-deps': ['error'],
    },
  },
  prettierConfig,
]
