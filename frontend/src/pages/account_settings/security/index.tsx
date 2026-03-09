import { useAuthFrontendComponent } from 'src/authschemes'
import { type SupportedAuthenticationScheme, type UserOwnView } from 'src/global_types'

export default function Security(props: { user: UserOwnView | null }) {
  const { user } = props
  if (user == null) return null

  return (
    <>
      {user.authSchemes.map((authScheme) => (
        <AuthSchemeSettings
          key={authScheme.schemeCode}
          authSchemeDetails={authScheme.authDetails}
          authSchemeType={authScheme.schemeType}
          username={authScheme.username}
        />
      ))}
    </>
  )
}

const AuthSchemeSettings = (props: {
  authSchemeDetails?: SupportedAuthenticationScheme
  authSchemeType: string
  username: string
}) => {
  const Settings = useAuthFrontendComponent(
    props.authSchemeType,
    'Settings',
    props.authSchemeDetails,
  )
  return (
    <Settings username={props.username} authFlags={props.authSchemeDetails?.schemeFlags || []} />
  )
}
