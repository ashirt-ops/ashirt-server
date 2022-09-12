import req from 'src/services/data_sources/backend/request_helper'

export async function login(username: string, password: string) {
  await req('POST', '/auth/local/login', { username, password })
}

export async function requestRecovery(email: string) {
  await req('POST', '/auth/recovery/generateemail', { userEmail: email })
}

export async function register(i: {
  firstName: string,
  lastName: string,
  email: string,
  username: string,
  password: string,
  confirmPassword: string,
}) {
  if (i.password !== i.confirmPassword) {
    throw Error("Passwords do not match")
  }
  await req('POST', '/auth/local/register', i)
}

export async function userResetPassword(i: {
  newPassword: string,
  confirmPassword: string,
}) {
  if (i.newPassword !== i.confirmPassword) {
    throw Error("Passwords do not match")
  }
  await req('POST', '/auth/local/login/resetpassword', i)
}

export async function userChangePassword(i: {
  username: string,
  oldPassword: string,
  newPassword: string,
  confirmPassword: string,
}) {
  if (i.newPassword !== i.confirmPassword) {
    throw Error("Passwords do not match")
  }
  await req('PUT', '/auth/local/password', i)
}

export async function linkLocalAccount(i: {
  username: string
  password: string,
  confirmPassword: string
}) {
  await req('POST', "/auth/local/link", i)
}

export async function totpIsEnabled() {
  return await req('GET', '/auth/local/totp')
}

export async function generateTotpSecret() {
  return await req('GET', '/auth/local/totp/generate')
}

export async function setTotp(data: { secret: string, passcode: string }) {
  return await req('POST', '/auth/local/totp', data)
}

export async function deleteTotp() {
  await req('DELETE', '/auth/local/totp')
}

export async function totpLogin(totpPasscode: string) {
  await req('POST', '/auth/local/login/totp', { totpPasscode })
}
