import { C2Event } from 'src/global_types'

type JsonC2Evidence = {
    c2: string,         // "Cobalt Strike"
    c2Operator: string, // The c2 frameworke user that issued task
    beacon: string,     // Beacon identifier
    externalIP: string,
    internalIP: string,
    hostname: string,
    userContext: string,// The user context that the implant/beacon/agent is running under
    integrity: string,  // "Low", "Medium", "High", "System"
    processName: string,// process image file shortname
    processID: number,  // is 'number' acceptable here? Its a float :/
    command: string,    // The actual command that an operator entered to task a beacon
    result: string,     // The result, if any, that a beacon responded to a tasking with
    metadata?: { [key: string]: string }
}

export const c2eventToBlob = (c2e: C2Event): Blob => {
  const evidence: JsonC2Evidence = {
    c2: c2e.c2,
    c2Operator: c2e.c2Operator,
    beacon: c2e.beacon,
    externalIP: c2e.externalIP,
    internalIP: c2e.internalIP,
    hostname: c2e.hostname,
    userContext: c2e.userContext,
    integrity: c2e.integrity,
    processName: c2e.processName,
    processID: c2e.processID,
    command: c2e.command,
    result: c2e.result,
  }

    evidence.metadata = {}

  return new Blob([JSON.stringify(evidence)])
}
