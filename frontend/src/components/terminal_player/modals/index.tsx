import * as React from 'react'
import Form from 'src/components/form'
import { useForm, useFormField } from 'src/helpers/use_form'
import { TextArea } from 'src/components/input'
import Modal from 'src/components/modal'

export const CreateBookmarkModal = (props: {
    updateEncoding: (bookmarkDesc: string) => string,
    onCancel: () => void,
    onRequestClose: () => void,
    handleSubmit: (b: Blob) => Promise<void>
    initialDescription: string,

}) => {
    const descriptionField = useFormField<string>(props.initialDescription)

    const formComponentProps = useForm({
        fields: [descriptionField],
        onSuccess: () => props.onRequestClose(),
        handleSubmit: () => {
            const newContent = props.updateEncoding(descriptionField.value)
            return props.handleSubmit(new Blob([newContent]))
        },
    })
    return (
        <Modal title="Update Bookmarks" onRequestClose={props.onRequestClose}>
            <Form submitText="Submit" cancelText="Cancel" onCancel={ () => {
                props.onCancel()
                props.onRequestClose()
            }} {...formComponentProps}>
                <TextArea label="Bookmark Description" {...descriptionField} />
            </Form>
        </Modal>
    )
}
