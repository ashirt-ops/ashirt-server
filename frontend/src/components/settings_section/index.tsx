import { type ReactNode } from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

const SettingsSection = (props: {
  children: ReactNode
  className?: string
  title: string
  width?: 'full-width' | 'wide' | 'normal' | 'narrow'
}) => (
  <section className={cx('root', props.width || 'normal', props.className)}>
    <h1 className={cx('title')}>{props.title}</h1>
    {props.children}
  </section>
)
export default SettingsSection
