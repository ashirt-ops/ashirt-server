import * as React from 'react'
import Button from 'src/components/button'
import classnames from 'classnames/bind'
import { OIDCInstanceConfig } from '..'

const cx = classnames.bind(require('./stylesheet'))

// use the below

const makeLoginFn = (code: string) => {
  return () => {
    window.location.href = `/web/auth/${code}/login`
  }
}

export const makeLogin = (config: OIDCInstanceConfig) => {
  const loginFn = makeLoginFn(config.code)

  return (_props: {
    authFlags?: Array<string>
  }) => (
    <div>
      <Button className={cx('full-width-button')} primary onClick={loginFn}>Login With {config.name}</Button>
    </div>
  )
}
