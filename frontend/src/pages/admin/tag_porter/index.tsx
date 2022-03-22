// Copyright 2022, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import SettingsSection from 'src/components/settings_section'
import { getDefaultTags, mergeDefaultTags } from 'src/services'
import { renderModals, useForm, useModal, useWiredData } from 'src/helpers'
import Modal from 'src/components/modal'
import Button from 'src/components/button'
import Form from 'src/components/form'
import { isError } from 'src/helpers/is_error'
import { isPlainObject } from 'lodash'
import { makeInvisibleDownloadAnchor } from 'src/helpers/invisible_download_anchor'
import { useDropzone } from 'react-dropzone'
import TagList from "src/components/tag_list"

const cx = classnames.bind(require('./stylesheet'))

export const TagPorter = (props: {
  requestReload: () => void
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const wiredTags = useWiredData(getDefaultTags)

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
              <Button className={cx('export-button')} primary onClick={
                () => {
                  const jsonData = JSON.stringify(tags
                    .map(tag => ({
                      name: tag.name,
                      colorName: tag.colorName,
                    })
                    ), null, 2)

                  makeInvisibleDownloadAnchor(new Blob([jsonData], { type: "application/json" }), "tag_export.json")
                }
              }>Export</Button>
            </div>
          </>
        ))}
      </SettingsSection>
      {renderModals(importTagsModal)}
    </>
  )
}

const ImportModal = (props: {
  onRequestClose: () => void,
  requestReload: () => void
}) => {

  const [err, setErr] = React.useState<Error | null>(null)
  const [content, setContent] = React.useState<Array<UpsertTag> | null>(null)
  const [contentFilename, setContentFilename] = React.useState<string | null>(null)

  const formComponentProps = useForm({
    onSuccess: () => { props.onRequestClose() },
    handleSubmit: async () => {
      if (content == null) {
        return
      }
      try {
        await mergeDefaultTags(content)
        props.requestReload()
        return
      }
      catch (err) {
        if (isError(err)) {
          throw new Error(
            (err instanceof SyntaxError)
              ? "Unable to parse data / Unsupported format"
              : err.message
          )
        }
        else {
          // this should never really happen, but is helpful in case something throws a non-error
          throw new Error("Something went wrong")
        }
      }

    }
  })

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    multiple: false,
    onDrop(acceptedFiles: Array<File>) {
      if (acceptedFiles.length === 1) {
        if (acceptedFiles[0].type == 'application/json') {
          setContentFilename(acceptedFiles[0].name)
          acceptedFiles[0].text()
            .then(value => {
              let parsedData = null
              try {
                parsedData = JSON.parse(value)
              }
              catch (err) {
                setErr(err)
                return
              }
              isTagArray(parsedData)
                ? setContent(parsedData)
                : setErr(new Error("Unsupported Format"))
            })
            .catch(err => {
              setErr(err)
            })
        }
        else {
          setErr(new Error("Unsupported Format"))
        }
      }
    },
  })

  return (
    <Modal title="Import Tags" onRequestClose={props.onRequestClose}>
      <div {...getRootProps({
        className: cx('import-upload-area', { active: isDragActive }),
      })}>
        <input {...getInputProps()} />
        <TagImportUpload contentFilename={contentFilename} err={err} />
      </div>
      <Form
        submitText='Import'
        cancelText='Close'
        onCancel={props.onRequestClose}
        {...formComponentProps}
      >
        <div>
          {content !== null && (
            <>
              <div className={cx('tag-list-header')}>These tags will be imported:</div>
              {/* This list will show duplicates, but the backend will filter them out. I expect that duplicates will be rare anyway */}
              <TagList tags={content.map((t, index) => ({
                ...t,
                id: index, // Faking the id, because we don't actually need a real id here
              }))} />
            </>
          )
          }
        </div>
      </Form>
    </Modal>
  )
}

const TagImportUpload = (props: {
  contentFilename: string | null
  err: Error | null,
}) => {
  if (props.contentFilename === null) {
    return (
      <div className={cx('no-content')}>
        Drag the tag_export.json here or <span>Browse for the file</span> to upload
        {props.err && <div className={cx('error')}>{props.err.message}</div>}
      </div>
    )
  }

  return (
    <div className={cx('has-content')}>
      Selected File: "{props.contentFilename}"
    </div>
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
