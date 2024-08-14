import * as React from 'react'
import WithLabel from 'src/components/with_label'
import classnames from 'classnames/bind'
import {useDropzone} from 'react-dropzone'
const cx = classnames.bind(require('./stylesheet'))

const BinaryUpload = (props: {
  disabled: boolean,
  label: string,
  onChange: (newValue: File | null) => void,
  isSupportedFile: (file: File ) => boolean,
  value: File | null,
  error: string,
}) =>  {
  const [err, setErr] = React.useState<Error | null>(null)

  const {value, isSupportedFile, label} = {...props}

  React.useEffect(() => {
    if (!props.error) {
        setErr(null)
    } else {
        setErr(Error(props.error))
    }
  }, [props.error])

  React.useEffect(() => {
    const file = value
    if (file == null || isSupportedFile(file)) {
      if (!props.error) setErr(null)
    } else {
      setErr(Error(`Expected a ${label.toLowerCase()}, but got ${file.type}.`))
    }
  }, [value, isSupportedFile, label, props.error])

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    multiple: false,
    noClick: true, // Required in chrome when using a label otherwise multiple file selects open
    onDrop(acceptedFiles: Array<File>) {
      if (acceptedFiles.length === 1) {
        props.onChange(acceptedFiles[0])
      }
    },
  })

  return (
    <WithLabel label={props.label}>
      <div {...getRootProps({
        className: cx('root', { active: isDragActive, disabled: props.disabled }),
      })}>
        <input {...getInputProps({ disabled: props.disabled })} />
        <BinaryUploadChildren file={props.value} err={err} friendlyFileType={`a ${props.label.toLowerCase()}`} />
      </div>
    </WithLabel>
  )
}

const BinaryUploadChildren = (props: {
  friendlyFileType: string,
  err: Error | null,
  file: File | null,
}) => {
  if (props.file !== null && props.err === null) {
      return (
        <div className={cx('has-content')}>
          <div>Will Upload: {props.file.name}</div>
        </div>
      )
  }

  return (
    <div className={cx('no-content')}>
      Drag {props.friendlyFileType} here or <span>Browse for one</span> to upload
        {props.err && <div className={cx('error')}>{props.err.message}</div>}
    </div>
  )
}

export default BinaryUpload
