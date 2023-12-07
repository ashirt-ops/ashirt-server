import req from 'src/services/data_sources/backend/request_helper'


export async function requestRecovery(email: string) {
  await req('POST', '/auth/recovery/generateemail', { userEmail: email })
}
