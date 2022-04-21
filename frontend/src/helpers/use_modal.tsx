// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'

// useModal & renderModals are helpers to make it easy for a component to display one or many modals to the user
//
// Call `useModal` with the type parameter of any props that your modal takes in that are dynamic.
// The returned object has a `show` method that takes in the dynamic props you want to render the modal with
// `useModal` will call the passed `modalRenderer` function with the specified dynamic props joined with
// onRequestClose that will handle closing the modal automatically.
//
// Example:
//
// const editUserModal = useModal<{user: User}>(modalProps => (
//   <EditUserModal {...modalProps} staticProp="Some static value" />
// ))
//
// return <>
//   <div>{users.map(user => (
//     <div key={user.id}>
//       <button onClick={editUserModal.show({user})}>Edit email for {user.name}</button>
//     </div>
//   ))}</div>
//
//   {renderModals(editUserModal)}
// </>

type UseModalOutput<ModalProps> = {
  node: React.ReactNode,
  show: (modalProps: ModalProps) => void,
}

type OnRequestClose = {
  onRequestClose: () => void,
}

export function useModal<ModalProps>(
  modalRenderer: (modalProps: ModalProps & OnRequestClose) => React.ReactNode,
  onClose?: () => void
): UseModalOutput<ModalProps> {
  const [modal, setModal] = React.useState<(ModalProps & OnRequestClose) | null>(null)
  const hide = () => { setModal(null) }

  return {
    node: modal == null ? null : modalRenderer(modal),
    show(modalProps: ModalProps) {
      setModal({
        ...modalProps, onRequestClose: () => {
          hide()
          onClose?.()
        }
      })
    },
  }
}

export function renderModals(...modals: Array<UseModalOutput<unknown>>): React.ReactNode {
  for (let modal of modals) {
    if (modal.node != null) return modal.node
  }
  return null
}
