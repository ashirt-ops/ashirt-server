// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'

import { NavLinkButton } from 'src/components/button'
import Form from 'src/components/form'
import Input from 'src/components/input'
import { useForm, useFormField } from 'src/helpers/use_form'

import { requestRecovery } from './services'
import { RecoverySchemeName } from './constants'

const cx = classnames.bind(require('./stylesheet'))

/**
 * Note that recovery is a special authentication scheme that cannot be disabled. Originally, it was
 * part of local auth, but it has been made to work generically. You cannot link recovery, nor are
 * there other settings. Likewise, you don't really "login" or register, but instead just get
 * recovery codes. As such, while this is an auth scheme, it is not treated like a normal auth scheme,
 * and instead just exposes the two components we need to ensure that recovery works, and is possible
 * for all auth schemes.
 */

/**
 * Returns either the initial recovery component, or the recovery-sent page, depending on the url
 * query string.
 */
export default (props: {
  query: URLSearchParams,
  authFlags?: Array<string>
}) => {
  if (props.query.get('step') === 'recovery-sent') {
    return <AccountRecoveryStarted />
  }
  return <RecoverUserAccount />
}

const RecoverUserAccount = (_: {}) => {
  const emailField = useFormField('')

  const emailForm = useForm({
    fields: [emailField],
    handleSubmit: () => {
      if (emailField.value.trim() == '') {
        return Promise.reject(Error("Please supply a valid email address"))
      }
      return requestRecovery(emailField.value).then(() => window.location.href = `/login/${RecoverySchemeName}?step=recovery-sent`)
    }
  })

  return (<>
    <h2 className={cx('title')}>Find Your Account</h2>
    <Form submitText="Submit" {...emailForm}>
      <Input label="Contact Email" {...emailField} />
    </Form>
  </>)
}

const AccountRecoveryStarted = (_: {}) => (
  <div>
    <div className={cx('messagebox')}>
      You should receive an email shortly with a recovery link.
    </div>
    <NavLinkButton primary className={cx('centered-button')} to={'/login'}>
      Return to Login
    </NavLinkButton>
  </div>
)
