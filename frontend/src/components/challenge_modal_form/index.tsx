// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'

import { useForm, useFormField } from 'src/helpers'
import ModalForm from 'src/components/modal_form'
import Input from 'src/components/input'

const cx = classnames.bind(require('./stylesheet'))

export default <T extends unknown>(props: {
  warningText: string,
  challengeText: string,
  modalTitle: string,
  submitText: string,
  onRequestClose: (success?:boolean) => void,
  handleSubmit: () => Promise<T>
}) => {

  const challenge = useFormField<string>("")

  const formComponentProps = useForm({
    fields: [challenge],
    onSuccess: () => props.onRequestClose(true),
    handleSubmit: () => {
      if (challenge.value !== props.challengeText) {
        return Promise.reject(Error("Challenge text does not match"))
      }
      return props.handleSubmit()
    }
  })

  return <ModalForm submitDanger title={props.modalTitle} submitText={props.submitText} onRequestClose={props.onRequestClose} {...formComponentProps}>
    <em className={cx('warning')}>{props.warningText}</em>
    <p className={cx('challenge-prompt')}>Please enter the following text into the textbox below to continue: </p>
    <em className={cx('challenge')}>{props.challengeText}</em>
    <Input label="Challenge" {...challenge} />
  </ModalForm>
}
