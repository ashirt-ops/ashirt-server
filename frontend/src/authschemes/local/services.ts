import req from 'src/services/data_sources/backend/request_helper'

export async function login(email: string, password: string) {
  await req('POST', '/auth/local/login', { email, password })
}

export async function register(i: {
  firstName: string,
  lastName: string,
  email: string,
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
  userKey: string,
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
  email: string
  password: string,
  confirmPassword: string
}) {
  await req('POST', "/auth/local/link", i)
}
