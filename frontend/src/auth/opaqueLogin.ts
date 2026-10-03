import {getRetryAfter, postJSON} from "../util/api.ts";

export type OpaqueLoginResult =
  | { status: "ok" }
  | { status: "wrong" }
  | { status: "error"; error?: string }
  | { status: "ratelimited"; retryAfter: number };

const ratelimited = (res: Response): OpaqueLoginResult => ({ status: "ratelimited", retryAfter: getRetryAfter(res) });
const readError = async (res: Response): Promise<string | undefined> => (await res.json().catch(() => null))?.error;

export async function opaqueLogin(params: {
  userId: number; role: string; password: string;
  identifiers?: { client?: string; server?: string };
}): Promise<OpaqueLoginResult> {
  const { ready, client } = await import("@serenity-kit/opaque")
  await ready;
  try {
    const { clientLoginState, startLoginRequest } =
      client.startLogin({ password: params.password });

    const initRes = await postJSON("/api/login/init", {
      user: params.userId, role: params.role, start_login_request: startLoginRequest,
    });
    if (initRes.status === 429) return ratelimited(initRes);
    if (!initRes.ok) return { status: "error", error: await readError(initRes) };

    const { login_response, login_handle } = await initRes.json();

    const result = client.finishLogin({
      clientLoginState, loginResponse: login_response, password: params.password,
      ...(params.identifiers ? { identifiers: params.identifiers } : {}),
    });
    if (!result) return { status: "wrong" };

    const finRes = await postJSON("/api/login/finish", {
      login_handle, finish_login_request: result.finishLoginRequest,
    });
    if (finRes.status === 429) return ratelimited(finRes);
    if (!finRes.ok) {
      const error = await readError(finRes);
      return { status: error === "WRONG_CREDENTIALS" ? "wrong" : "error", error };
    }
    return { status: "ok" };
  } catch {
    return { status: "error" };
  }
}
