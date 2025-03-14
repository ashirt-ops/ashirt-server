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
  error: string,
}) => {
  const [imageDataUriString, setImageDataUriString] = React.useState<string|null>(null)
  const [err, setErr] = React.useState<Error|null>(null)

  React.useEffect(() => {
    if (props.error) {
        setErr(Error(props.error))
    } else {
        setErr(null)
    }
  }, [props.error])

  React.useEffect(() => {
    const file = props.value
    if (file == null) {
      setErr(null)
      setImageDataUriString(null)
    } else if (file.type.startsWith('image/')) {
      setErr(null)
      getImageDataUriFromFile(file).then(setImageDataUriString)
    } else {
      setErr(Error(`Expected an image, but got ${file.type}.`))
      setImageDataUriString(null)
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
        <ImageUploadChildren image={imageDataUriString} err={err} />
      </div>
    </WithLabel>
  )
}

const ImageUploadChildren = (props: {
  err: Error | null,
  image: string | null,
}) => {
  if (props.image == null) {
    return (
      <div className={cx('no-image')}>
        <img src={require('./image.svg')} />
        Drag an image here or <span>Browse for an image</span> to upload
        {props.err && <div className={cx('error')}>{props.err.message}</div>}
      </div>
    )
  }

  return (
    <div className={cx('has-images')}>
      <div className={cx('thumb')} style={{backgroundImage: `url(${props.image})`}} />
    </div>
  )
}

function getImageDataUriFromFile(file: File): Promise<string> {
  return new Promise(res => {
    const reader = new FileReader()
    // @ts-ignore - (https://github.com/Microsoft/TypeScript/issues/299)
    reader.onload = (e: ProgressEvent) => res(e.target.result)
    reader.readAsDataURL(file)
  })
}
