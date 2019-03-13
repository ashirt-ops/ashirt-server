// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import ErrorDisplay from 'src/components/error_display'
import LoadingSpinner from 'src/components/loading_spinner'

// useAsyncComponent is a react hooks helper to make it easy to enable code splitting in a type-safe way
// It takes a function that returns a react component and properly renders a loading spinner in place of
// the component while the javascript bundle is being downloaded
//
// Using this helper reduces the bundle size of the initial javascript bundle served to the client
//
// See https://webpack.js.org/guides/code-splitting/
//
// Example:
// In my_component.tsx:
// export default (props: {some: string}) => {
//   ...
// }
//
// In other_component.tsx:
// const importAsyncMyComponent = () => import('./my_component').then(module => module.default)
//
// () => {
//   const AsyncMyComponent = useAsyncComponent(importAsyncMyComponent)
//   return (
//     <div>
//       <AsyncMyComponent some="props" />
//     </div>
//   )
// }
//
//
// Note that `useAsyncComponent` passes the `getComponentFn` as a dependency to useEffect so components that
// use `useAsyncComponent` should either be called with a static import function (as shown in the example above)
// or should call `React.useCallback` like so:
//
// const AsyncMyComponent = useAsyncComponent(React.useCallback(() => (
//   someImporterFn
// ), [someImporterFn]))

type ComponentOrLoadingOrError<T> = T | React.FunctionComponent<{}>

export function useAsyncComponent<T>(getComponentFn: () => Promise<T>): ComponentOrLoadingOrError<T> {
  // The extra function wrapper in useState and setComponent is required since useState will call the argument
  // if it is a function and in this case we always call setComponent with a React.FunctionComponent
  const [component, setComponent] = React.useState<ComponentOrLoadingOrError<T>>(() => () => <LoadingSpinner />)

  React.useEffect(() => {
    getComponentFn()
      .then(component => setComponent(() => component))
      .catch(err => setComponent(() => () => <ErrorDisplay title="Failed to load frontend module" err={err} />))
  }, [getComponentFn])

  return component
}
