import { Component, type ReactNode, type ErrorInfo } from 'react'
import ErrorDisplay from 'src/components/error_display'

type Props = { children: ReactNode }
type State = { error: Error | null }

export default class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null }

  static getDerivedStateFromError(error: Error): State {
    return { error }
  }

  componentDidCatch(_error: Error, _info: ErrorInfo) {
    // Errors are surfaced via the render method
  }

  render() {
    if (this.state.error) {
      return <ErrorDisplay err={this.state.error} />
    }
    return this.props.children
  }
}
