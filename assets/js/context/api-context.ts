import { createContext } from 'react';

interface CommonAPI {
  request: () => Promise<any>
}

async function PaymentAPI() {
  const errorEvent = new CustomEvent("error-event")

  document.addEventListener("error-event", () => {
    console.log("Error Occurred")
  })

  const handleRefreshToken = async () => {
    const res = await fetch("re", {
      method: "POST",
      headers: {}
    })

    document.dispatchEvent(errorEvent)
  }

  const mainRequest = async (): Promise<Response> => {
    const result = await fetch("xxx", {
      method: "POST",
      body: JSON.stringify({})
    })

    return result
  }

  const doRequest = async () => {
    const sessionToken = localStorage.getItem("session-token")
    const refreshToken = localStorage.getItem("refresh-token")
    const result = await mainRequest()
    if (result.status == 401) {
      await handleRefreshToken()
      await mainRequest()
    }
  }
}

const ApiContext = createContext("api-provider")
export default ApiContext