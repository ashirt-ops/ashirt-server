import { Component, type ReactNode } from 'react'
import Button from 'src/components/button'
import ErrorDisplay from 'src/components/error_display'

type Props = { children: ReactNode; resetKey?: string }
type State = { error: Error | null }

export default class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null }

  static getDerivedStateFromError(error: unknown): State {
    return { error: error instanceof Error ? error : new Error(String(error)) }
  }

  componentDidUpdate(previousProps: Props) {
    if (this.state.error && previousProps.resetKey !== this.props.resetKey) {
      this.setState({ error: null })
    }
  }

  render() {
    if (this.state.error) {
      return (
        <ErrorDisplay err={this.state.error}>
          <div>
            <Button doNotSubmit onClick={() => window.location.reload()}>
              Reload page
            </Button>
          </div>
        </ErrorDisplay>
      )
    }
    return this.props.children
  }
}
