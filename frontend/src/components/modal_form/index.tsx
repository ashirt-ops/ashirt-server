// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Form from 'src/components/form'
import Modal from 'src/components/modal'
import {Result} from 'src/global_types'

// A convenience wrapper around <Modal><Form>{...}</Form></Modal> since it's such a common UI piece
export default (props: {
  children: React.ReactNode,
  result: Result<string>|null,
  loading: boolean,
  onRequestClose: () => void,
  onSubmit: (e: React.FormEvent) => void,
  submitText?: string,
  title: string,
  submitDanger?:boolean,
}) => (
  <Modal title={props.title} onRequestClose={props.onRequestClose}>
    <Form
      cancelText="Cancel"
      children={props.children}
      result={props.result}
      loading={props.loading}
      onCancel={props.onRequestClose}
      onSubmit={props.onSubmit}
      submitText={props.submitText}
      submitDanger={props.submitDanger}
    />
  </Modal>
)
