import * as React from 'react'
import WithLabel from 'src/components/with_label'
import classnames from 'classnames/bind'
import {useDropzone} from 'react-dropzone'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  disabled: boolean,
  label: string,
  onChange: (newValue: File|null) => void,
  value: File | null,
}) => {
  const [err, setErr] = React.useState<Error|null>(null)

  React.useEffect(() => {
    const file = props.value
    if (file == null || file.type === '') {
      setErr(null)
    } else {
      setErr(Error(`Expected a terminal recording, but got ${file.type}.`))
    }
  }, [props.value])

  const {getRootProps, getInputProps, isDragActive} = useDropzone({
    multiple: false,
    noClick: true, // Required in chrome when using a label otherwise multiple file selects open
    onDrop(acceptedFiles: Array<File>) {
      if (acceptedFiles.length === 1) props.onChange(acceptedFiles[0])
    },
  })

  return (
    <WithLabel label={props.label}>
      <div {...getRootProps({
        className: cx('root', {active: isDragActive, disabled: props.disabled}),
      })}>
        <input {...getInputProps({disabled: props.disabled})} />
        <TermRecUploadChildren file={props.value} err={err} />
      </div>
    </WithLabel>
  )
}

const TermRecUploadChildren = (props: {
  err: Error | null,
  file: File | null,
}) => {
  const content = (props.file != null  && props.err == null)
    ? <div className={cx('has-content')}>
        <div>Will Upload: {props.file.name}</div>
      </div>
    : <div className={cx('no-content')}>
        Drag a recording here or <span>Browse for one</span> to upload
        {props.err && <div className={cx('error')}>{props.err.message}</div>}
      </div>

  return content
}
