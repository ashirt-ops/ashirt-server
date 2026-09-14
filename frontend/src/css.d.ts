// Ambient declaration so side-effect CSS imports (e.g. `import 'pkg/foo.css'`)
// type-check. Required since TypeScript 6.0, which errors (TS2882) on
// side-effect imports that lack type declarations.
declare module '*.css'
