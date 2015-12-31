typedef binary bytes;

// A service request, somewhat compatible with HTTP
struct Request {
  // Name of the remote service
  1: string service_id,
  
  // Path to script handling the request on the remote service
  2: string path,

  // Round-trip timeout for response to this request
  3: i32 timeout_m_s,

  // Request payload
  4: bytes body,

  // User-provided HTTP headers
  5: map<string, string> headers,
  
  // Request parameters (Query Params or Form Data). Use only for backward-compatibility
  6: map<string, string> params,

  // Script path to receive the response as a non-blocking callback
  7: string callback,
}

// Service response, somewhat compatible with HTTP
struct Response {
  // Response status code (numeric)
  1: i32 status_code,

  // Human readable response status message
  2: string status,

  // Response body
  3: bytes body,
}

service AgentIO {
  Response SendRequest(1:Request req)
}
