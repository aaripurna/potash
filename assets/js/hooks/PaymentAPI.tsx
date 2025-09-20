export default function usePaymentAPI() {

  const doRefreshToken = async (): Promise<Response> => {
    const refreshToken = localStorage.getItem("refresh-token")
    const res = await fetch("https://foo.bar", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        refresh_token: refreshToken,
      })
    })

    if (res.status == 200) {
      localStorage.setItem("session-token", "foo")
      localStorage.setItem("refresh-token", "bar-foo")
    }

    return res
  }

  const makeRequest = async (input: RequestInfo | URL, init?: RequestInit): Promise<Response> => {
    const sessionToken = localStorage.getItem("session-token")
    const refreshToken = localStorage.getItem("refresh-token")

    const requestData = init || {}

    const headers: HeadersInit = {}

    requestData.headers ||= headers

    headers['Authorization'] = `Bearer ${sessionToken}`

    requestData.headers = { ...requestData.headers, ...headers }

    const res = await fetch(input, requestData)
    return res
  }

  const doRequest = async (input: RequestInfo | URL, init?: RequestInit): Promise<Response> => {

    const res = await makeRequest(input, init)

    if (res.status == 401) {
      const refreshTokenResponse = await doRefreshToken()

      if (refreshTokenResponse.status != 200) {
        throw new Error("UNAUTHORIZED")
      }

      const res = await makeRequest(input, init)

      return res
    }

    return res
  }

  return {
    doRequest
  }
}

const { doRequest } = usePaymentAPI()

doRequest("https://dsdasdasdas", {
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({})
}).then((e) => {

})