// Copyright 2022, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import SettingsSection from 'src/components/settings_section'
import { getDefaultTags, mergeDefaultTags } from 'src/services'
import { renderModals, useForm, useFormField, useModal, useWiredData } from 'src/helpers'
import Modal from 'src/components/modal'
import Button from 'src/components/button'
import Form from 'src/components/form'
import { Tag } from 'src/global_types'
import { SourcelessCodeblock } from 'src/components/code_block'
import { isError } from 'src/helpers/is_error'
import { isPlainObject } from 'lodash'

const cx = classnames.bind(require('./stylesheet'))

export const TagPorter = (props: {
  requestReload: () => void
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const wiredTags = useWiredData(getDefaultTags)

  const exportTagsModal = useModal<{ tags: Array<Tag> }>(modalProps => (
    <ExportModal {...modalProps} />
  ))

  const importTagsModal = useModal(modalProps => (
    <ImportModal {...modalProps} requestReload={props.requestReload} />
  ))

  React.useEffect(() => {
    props.onReload(wiredTags.reload)
    return () => { props.offReload(wiredTags.reload) }
  })

  return (
    <>
      <SettingsSection title="Import/Export Default Tags">
        <em>You can share your default tags with other AShirt users by import or exporting the default tags here.</em>
        {wiredTags.render(tags => (
          <>
            <div className={cx('button-group')}>
              <Button primary onClick={() => importTagsModal.show({})}>Import</Button>
              <Button className={cx('export-button')} primary onClick={() => exportTagsModal.show({ tags })}>Export</Button>
            </div>
          </>
        ))}
      </SettingsSection>
      {renderModals(exportTagsModal, importTagsModal)}
    </>
  )
}

const ExportModal = (props: {
  tags: Array<Tag>
  onRequestClose: () => void,
}) => {
  const jsonData = JSON.stringify(props.tags
    .map(tag => ({
      name: tag.name,
      colorName: tag.colorName,
    })
    ), null, 2)

  const formComponentProps = useForm({
    onSuccess: () => { props.onRequestClose() },
    handleSubmit: async () => {
      await navigator.clipboard.writeText(jsonData)
    }
  })

  return (
    <Modal title="Export Tags" onRequestClose={props.onRequestClose}>
      <Form submitText="Copy" cancelText="Close" onCancel={props.onRequestClose} {...formComponentProps}>
        <p>Copy the below text and paste it into another AShirt instance</p>
        <SourcelessCodeblock
          className={cx('codeblock')}
          code={jsonData}
          language="json"
        />
      </Form>
    </Modal>
  )
}

const ImportModal = (props: {
  onRequestClose: () => void,
  requestReload: () => void
}) => {

  const codeblockField = useFormField<string>("")
  const formComponentProps = useForm({
    fields: [codeblockField],
    onSuccess: () => { props.onRequestClose() },
    handleSubmit: async () => {
      try {
        const parsedJson = JSON.parse(codeblockField.value) as unknown
        if (!isTagArray(parsedJson)) {
          throw new SyntaxError()
        }

        await mergeDefaultTags(parsedJson)
        props.requestReload()
        return
      }
      catch (err) {
        if (isError(err)) {
          if (err instanceof SyntaxError) {
            throw new Error("Unable to parse data / Unsupported format")
          }
          else {
            throw new Error(err.message)
          }
        }
        else {
          // this should never really happen, but is helpful in case something throws a non-error
          throw new Error("Something went wrong")
        }
      }
    }
  })

  return (
    <Modal title="Import Tags" onRequestClose={props.onRequestClose}>
      <Form submitText="Submit" cancelText="Cancel" onCancel={props.onRequestClose} {...formComponentProps}>
        <p>Enter the serialized tags in the textbox below. Tags will be merged by name.</p>
        <SourcelessCodeblock
          className={cx('codeblock')}
          code={codeblockField.value}
          editable={true}
          onChange={codeblockField.onChange}
          language="json"
        />
      </Form>
    </Modal>
  )
}

type UpsertTag = {
  name: string
  colorName: string
}

const isTag = (t: unknown): t is UpsertTag => {
  if (isPlainObject(t)) {
    const maybeTag = t as Record<string, string>
    return (
      typeof (maybeTag['name']) === 'string' &&
      typeof (maybeTag['colorName']) === 'string'
    )
  }
  return false
}

const isTagArray = (t: unknown): t is Array<UpsertTag> => {
  if (Array.isArray(t) && t.length > 0) {
    return t.map(isTag).every(x => x === true)
  }
  return false
}
