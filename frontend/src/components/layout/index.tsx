import { type ReactNode } from 'react'
import NavBar from './nav_bar'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

const Layout = (props: { children: ReactNode }) => (
  <main className={cx('root')}>
    <NavBar />
    <div className={cx('children')}>{props.children}</div>
  </main>
)
export default Layout
