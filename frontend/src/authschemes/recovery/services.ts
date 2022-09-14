// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import req from 'src/services/data_sources/backend/request_helper'


export async function requestRecovery(email: string) {
  await req('POST', '/auth/recovery/generateemail', { userEmail: email })
}
