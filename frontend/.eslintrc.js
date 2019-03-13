module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  plugins: [
    '@typescript-eslint',
    'react',
    'react-hooks',
  ],
  parserOptions: {
    ecmaVersion: '2018',
    sourceType: 'module',
    ecmaFeatures: { jsx: true },
  },
  extends: [
  ],
  rules: {
    '@typescript-eslint/no-unused-vars': ['error', {args: 'none', varsIgnorePattern: '^_'}],
    'eol-last': ['error', 'always'],
    'no-multiple-empty-lines': ['error', { max: Infinity, maxEOF: 0 }],
    'no-trailing-spaces': ['error'],
    'react/jsx-uses-react': ['error'],
    'react/jsx-uses-vars': ['error'],
    'react-hooks/rules-of-hooks': ['error'],
    'react-hooks/exhaustive-deps': ['warn'],
  }
}
