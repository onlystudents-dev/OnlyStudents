export type OpaqueLoginResult =
  | { status: "ok" }
  | { status: "wrong" }
  | { status: "error"; error?: string }
  | { status: "ratelimited"; retryAfter: number };

export async function opaqueLogin(params: {
  userId: number; role: string; password: string;
  identifiers?: { client?: string; server?: string };
}): Promise<OpaqueLoginResult> {
  const { ready, client } = await import("@serenity-kit/opaque")
  await ready;
  try {
    const { clientLoginState, startLoginRequest } =
      client.startLogin({ password: params.password });

    const initRes = await fetch("/api/login/init", {
      method: "POST", headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ user: params.userId, role: params.role, start_login_request: startLoginRequest }),
    });
    if (initRes.status === 429) return { status: "ratelimited", retryAfter: Number(initRes.headers.get("Retry-After") ?? "-1") };
    if (!initRes.ok) return { status: "error", error: (await initRes.json().catch(() => null))?.error };

    const { login_response, login_handle } = await initRes.json();

    const result = client.finishLogin({
      clientLoginState, loginResponse: login_response, password: params.password,
      ...(params.identifiers ? { identifiers: params.identifiers } : {}),
    });
    if (!result) return { status: "wrong" };

    const finRes = await fetch("/api/login/finish", {
      method: "POST", headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ login_handle, finish_login_request: result.finishLoginRequest }),
    });
    if (finRes.status === 429) return { status: "ratelimited", retryAfter: Number(finRes.headers.get("Retry-After") ?? "-1") };
    if (!finRes.ok) {
      const error = (await finRes.json().catch(() => null))?.error;
      return { status: error === "WRONG_CREDENTIALS" ? "wrong" : "error", error };
    }
    return { status: "ok" };
  } catch {
    return { status: "error" };
  }
}
