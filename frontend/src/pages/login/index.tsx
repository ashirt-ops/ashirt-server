// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { useLocation, useParams } from 'react-router-dom'
import { getSupportedAuthentications } from 'src/services/auth'
import { useAuthFrontendComponent } from 'src/authschemes'
import { useWiredData } from 'src/helpers'
import { SupportedAuthenticationScheme } from 'src/global_types'
import Button from 'src/components/button'
const cx = classnames.bind(require('./stylesheet'))

import req from 'src/services/data_sources/backend/request_helper'

// This component renders a list of all enabled authscheme login components
// To add a new authentication method add a new authscheme frontend to
// src/auth and ensure it is enabled on the backend
// An optional schemeCode can be provided to only render that auth method
export default () => {
  const { schemeCode: renderOnlyScheme } = useParams<{ schemeCode?: string }>()
  const location = useLocation()
  const query = new URLSearchParams(location.search)
  const wiredAuthSchemes = useWiredData(getSupportedAuthentications)

  return wiredAuthSchemes.render(supportedAuthSchemes => (
    <div className={cx('login')}>
      {supportedAuthSchemes.map((schemeDetails) => {
        const { schemeCode, schemeType } = schemeDetails
        if (renderOnlyScheme != null && schemeCode != renderOnlyScheme) {
          return null
        }
        return (
          <AuthSchemeLogin
            key={schemeCode}
            authSchemeType={schemeType}
            authScheme={schemeDetails}
            query={query}
          />
        )
      })}
      {/* Check if Webauthn is available */}
      {window.PublicKeyCredential &&
        <div>
          <Button onClick={async () => {
            const raw = await beginRegistration({
              email: "joel.smith+temp@originate.com",
              firstName: "Joel",
              lastName: "smith"
            })

            console.log(raw)
            const { publicKey: key } = raw
            const challenge = toByteArray(key.challenge)
            const signed = await navigator.credentials.create({
              publicKey: {
                challenge,
                rp: key.rp,
                user: {
                  ...key.user,
                  id: toByteArray(key.user.id)
                },
                pubKeyCredParams: key.pubKeyCredParams,
                // authenticatorSelection: {},
                // attestation: 'direct'
              }
            })
            console.log("signed credential", signed)
            if (signed == null || signed.type != 'public-key') {
              return; // TODO: not sure what to do here.
            }
            // await finishRegistration(signed)
            // @ts-ignore
            await finishRegistration({
              type: 'public-key',
              id: signed.id,
              // @ts-ignore
              rawId: signed.rawid,
              response: {
                // @ts-ignore
                attestationObject: base64UrlEncode(signed.response.attestationObject),
                // @ts-ignore
                clientDataJSON: base64UrlEncode(signed.response.clientDataJSON),
              },
            })


          }}>Yo</Button>
        </div>
      }
    </div>
  ))
}

const base64UrlEncode = (value: string) => {
  return btoa(value)
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "")
}

const toKeycode = (c: string)=> c.charCodeAt(0)
const toByteArray = (s: string) => Uint8Array.from(s, toKeycode)

const AuthSchemeLogin = (props: {
  authSchemeType: string,
  authScheme: SupportedAuthenticationScheme,
  query: URLSearchParams,
}) => {
  const Login = useAuthFrontendComponent(props.authSchemeType, 'Login', props.authScheme)
  return (
    <div className={cx('auth-scheme-row')}>
      <Login query={props.query} authFlags={props.authScheme.schemeFlags} />
    </div>
  )
}

async function beginRegistration(i: {
  email: string,
  firstName: string,
  lastName: string
}): Promise<WebAuthNCredentialCreationResponse> {
  return await req('POST', '/auth/webauthn/register/begin', i)
}

async function finishRegistration(i: WebAuthNRegisterConfirmation) {
  return await req('POST', '/auth/webauthn/register/end', i)
}

type WebAuthNCredentialCreationResponse = {
  publicKey: {
    challenge: string, //b64 string
    rp: {
      id: string
      name: string
      icon?: string
    }
    user: {
      id: string
      name: string
      displayName: string
      icon?: string
    }
    pubKeyCredParams: Array<{
      type: "public-key"
      alg: number // int. There are a few options here
    }>
    authenticatorSelection?: {
      authenticatorAttachment?: string // "platform" | "cross-platform"
      requireResidentKey?: boolean | null
      residentKey?: string // "required" | "preferred" | "discouraged"
      userVerification?: string // "required" | "preferred" | "discouraged"
    }
    timeout?: number // int
    excludeCredentials?: Array<{
      type: string //"public-key"
      id: string
      transports?: Array<string> // Array<AuthTransport>
    }>
    extensions?: Record<string, unknown>
    attestation?: string // "none" | "indirect" | "direct"
  }
}

type WebAuthNRegisterConfirmation = {
  id: string
  rawId: string // base64
  type: 'public-key'
  response: {
    clientDataJSON: string // base64
    attestationObject: string // base64
  }
}

type AuthTransport = 
  | "usb"
  | "nfc"
  | "ble"
  | "internal"

const AlgES256 = -7
const AlgES384 = -35
const AlgES512 = -36
const AlgRS1 = -65535
const AlgRS256 = -257
const AlgRS384 = -258
const AlgRS512 = -259
const AlgPS256 = -37
const AlgPS384 = -38
const AlgPS512 = -39
const AlgEdDSA = -8
