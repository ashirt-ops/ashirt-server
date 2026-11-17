import { createContext } from 'react'

// Modal popovers must remain in the dialog's accessible subtree and top layer.
export const PopoverPortalContext = createContext<HTMLElement | null>(null)
