import ErrorDisplay from '../../components/error_display'
import { useLocation } from 'react-router'

export default function NotFound() {
  const { pathname } = useLocation()
  return <ErrorDisplay err={new Error(`404 - The path ${pathname} is invalid`)} />
}
