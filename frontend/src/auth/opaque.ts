import {type ApiError, getRetryAfter, postJSON} from "../util/api.ts";
import { decodeOpaqueExportKey, deriveWrapKey, unwrapMK, wrapMK, b64encode, type Bytes } from "../util/crypto/keys.ts";

type Identifiers = {client?: string, server?: string}

type OpaqueParams = {
    password: string,
    userId?: number,
    role?: string,
    enrollToken?: string,
    identifiers?: Identifiers,
}

export type ChangePasswordParams = {
    password: string,
    identifiers?: Identifiers,
}

export type RekeyParams = {
    password: string,
    wrapped_mk_password: Bytes,
    crypto_grant_handle: string,
    exportKey: string,
}

type LoginInitResponse = {
    login_response: string,
    login_handle: string,
}

type EnrollInitResponse = {
    enroll_response: string,
}

type ReauthFinishResponse = {
    crypto_grant_handle: string,
}

type Failure = {status: "error"; error?: string | undefined}
type RateLimited = {status: "ratelimited"; retryAfter: number}

export type OpaqueLoginResult = {"status": "ok", exportKey: string} | {status: "wrong"} | Failure | RateLimited
export type OpaqueRegisterResult = {"status": "ok", exportKey: string} | Failure | RateLimited
export type OpaqueReauthResult = {"status": "ok", exportKey: string, crypto_grant_handle: string} | {status: "wrong"} | Failure | RateLimited
export type OpaqueRekeyResult = {"status": "ok"} | {status: "wrong"} | Failure | RateLimited

const rateLimited = (res: Response): RateLimited => ({
    status: "ratelimited",
    retryAfter: getRetryAfter(res),
})

const failure = (error?: string): Failure => ({status: "error", error})

const readError = async (res: Response): Promise<string | undefined> => (await res.json().catch(() => null) as ApiError | null)?.error

const identifierOpts = ({identifiers}: OpaqueParams) =>
    identifiers ? {identifiers} : {};

async function loadOpaque() {
    const {ready, client} = await import("@serenity-kit/opaque")
    await ready
    return client
}

export async function opaqueLogin(params: OpaqueParams): Promise<OpaqueLoginResult> {
    try {
        const client = await loadOpaque()

        const {clientLoginState, startLoginRequest} = client.startLogin({password: params.password})

        const initRes = await postJSON("/api/v1/login/init", {
            user: params.userId,
            role: params.role,
            start_login_request: startLoginRequest,
        })

        if (initRes.status === 429) return rateLimited(initRes)
        if (!initRes.ok) return failure(await readError(initRes))

        const {login_response, login_handle} = await initRes.json() as LoginInitResponse

        const result = client.finishLogin({
            clientLoginState,
            loginResponse: login_response,
            password: params.password,
            ...identifierOpts(params),
        })

        if (!result) return {status: "wrong"}

        const finRes = await postJSON("/api/v1/login/finish", {
            login_handle,
            finish_login_request: result.finishLoginRequest,
        })

        if (finRes.status === 429) return rateLimited(finRes)

        if (!finRes.ok) {
            const error = await readError(finRes)
            return error === "WRONG_CREDENTIALS" ? {status: "wrong"} : failure(error)
        }

        return {
          "status": "ok",
          "exportKey": result.exportKey
        }
    } catch {
        return failure()
    }
}

export async function opaqueRegister(params: OpaqueParams): Promise<OpaqueRegisterResult> {
    try {
        const client = await loadOpaque()

        const {clientRegistrationState, registrationRequest} = client.startRegistration({password: params.password})

        const initRes = await postJSON("/api/v1/enroll/init", {
            enroll_token: params.enrollToken,
            start_enroll_request: registrationRequest,
        })

        if (initRes.status === 429) return rateLimited(initRes)
        if (!initRes.ok) return failure(await readError(initRes))

        const {enroll_response} = await initRes.json() as EnrollInitResponse

        const result = client.finishRegistration({
            clientRegistrationState,
            registrationResponse: enroll_response,
            password: params.password,
            ...identifierOpts(params),
        })

        const finishRes = await postJSON("/api/v1/enroll/finish", {
            enroll_token: params.enrollToken,
            finish_enroll_request: result.registrationRecord,
        })

        if (finishRes.status === 429) return rateLimited(finishRes)
        if (!finishRes.ok) return failure(await readError(finishRes))

        return {
          "status": "ok",
          "exportKey": result.exportKey
        }
    } catch {
        return failure()
    }
}

export async function opaqueReauth(params: ChangePasswordParams): Promise<OpaqueReauthResult> {
    try {
        const client = await loadOpaque()

        const {clientLoginState, startLoginRequest} = client.startLogin({password: params.password})

        const initRes = await postJSON("/api/v1/me/change_password/init", {
            start_login_request: startLoginRequest,
        })

        if (initRes.status === 429) return rateLimited(initRes)
        if (!initRes.ok) return failure(await readError(initRes))

        const {login_response, login_handle} = await initRes.json() as LoginInitResponse

        const result = client.finishLogin({
            clientLoginState,
            loginResponse: login_response,
            password: params.password,
            ...identifierOpts(params),
        })

        if (!result) return {status: "wrong"}

        const finRes = await postJSON("/api/v1/me/change_password/finish", {
            login_handle,
            finish_login_request: result.finishLoginRequest,
        })

        if (finRes.status === 429) return rateLimited(finRes)

        if (!finRes.ok) {
            const error = await readError(finRes);
            return error === "WRONG_CREDENTIALS" ? {status: "wrong"} : failure(error)
        }

        const crypto_grant_handle = (await finRes.json() as ReauthFinishResponse).crypto_grant_handle

        return {
            "status": "ok",
            "exportKey": result.exportKey,
            "crypto_grant_handle": crypto_grant_handle,
        }
    } catch {
        return failure();
    }
}

export async function opaqueRekey(params: RekeyParams): Promise<OpaqueRekeyResult> {
    try {
        const client = await loadOpaque()

        const {clientRegistrationState, registrationRequest} = client.startRegistration({password: params.password})

        const initRes = await postJSON("/api/v1/me/change_password/rekey_init", {
            start_enroll_request: registrationRequest,
        })

        if (initRes.status === 429) return rateLimited(initRes)
        if (!initRes.ok) return failure(await readError(initRes))

        const {enroll_response} = await initRes.json() as EnrollInitResponse

        const result = client.finishRegistration({
            clientRegistrationState,
            registrationResponse: enroll_response,
            password: params.password,
            ...identifierOpts(params),
        })

        const unwrapped_mk = await unwrapMK(params.wrapped_mk_password, await deriveWrapKey(decodeOpaqueExportKey(params.exportKey), "password"))
        const wrapped_mk_password = await wrapMK(unwrapped_mk, await deriveWrapKey(decodeOpaqueExportKey(result.exportKey), "password"))

        const finishRes = await postJSON("/api/v1/me/change_password/rekey_finish", {
            crypto_grant_handle: params.crypto_grant_handle,
            wrapped_mk_password: b64encode(wrapped_mk_password),
            finish_enroll_request: result.registrationRecord,
        })

        if (finishRes.status === 429) return rateLimited(finishRes)
        if (!finishRes.ok) return failure(await readError(finishRes))

        return {
            "status": "ok",
        }
    } catch {
        return failure();
    }
}
