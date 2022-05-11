import { ViewName } from 'src/global_types'

export type NavigateProps = {
  slug: string
  view: ViewName
  query: string
  navTo: (path: string) => void
  searchName?: string
}

export type NavToFunction = (view: ViewName, query: string, searchName?: string) => void

/**
 * navigate provides a consistent way to navigate to the evidence or finding page with a specifed
 * search string
 *
 * @param props.slug The operation to navigate to
 * @param props.view The view (evidence or finding) to navigate to
 * @param props.query The query to use for filtering evidence/findings
 * @param props.navTo The function that actually navigates
 */
export const navigate = (props: NavigateProps) => {
  const { slug, view, query, searchName, navTo } = props

  let path = `/operations/${slug}/${view}`
  if (query != '') {
    path += `?q=${encodeURIComponent(query.trim())}`
    if (searchName) {
      path += `&name=${encodeURIComponent(searchName.trim())}`
    }
  }
  navTo(path)
}

/**
 * mkNavTo constructs a navigate function that is compatible with existing places where
 * the old navigate function was used
 *
 * @param props.slug Which search to use
 * @param props.navTo the function to use to trigger the navigation
 * @returns a function that will navigate to another page based on the provided data
 */
export const mkNavTo = (props: {
  slug: string
  navTo: (path: string) => void
}): NavToFunction => {
  return (view: ViewName, query: string, searchName?: string) => navigate({
    ...props, view, query, searchName
  })
}
