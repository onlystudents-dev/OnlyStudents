import {type ApiError, getRetryAfter, postJSON, readJSON} from "../util/api.ts";

type Identifiers = {client?: string, server?: string};

type OpaqueParams = {
    password: string;
    userId?: number;
    role?: string;
    enrollToken?: string;
    identifiers?: Identifiers;
};

type LoginInitResponse = {
    login_response: string;
    login_handle: string;
};

type EnrollInitResponse = {
    enroll_response: string;
};

type Ok = {status: "ok"};
type Failure = {status: "error"; error?: string | undefined};
type RateLimited = {status: "ratelimited"; retryAfter: number};

export type OpaqueLoginResult = Ok | {status: "wrong"} | Failure | RateLimited;
export type OpaqueRegisterResult = Ok | Failure | RateLimited;

const ok: Ok = {status: "ok"};

const rateLimited = (res: Response): RateLimited => ({
    status: "ratelimited",
    retryAfter: getRetryAfter(res),
});

const failure = (error?: string): Failure => ({status: "error", error});

const readError = async (res: Response): Promise<string | undefined> =>
    (await readJSON<ApiError | null>(res).catch(() => null))?.error;

const identifierOpts = ({identifiers}: OpaqueParams) =>
    identifiers ? {identifiers} : {};

async function loadOpaque() {
    const {ready, client} = await import("@serenity-kit/opaque");
    await ready;
    return client;
}

export async function opaqueLogin(params: OpaqueParams): Promise<OpaqueLoginResult> {
    try {
        const client = await loadOpaque();

        const {clientLoginState, startLoginRequest} =
            client.startLogin({password: params.password});

        const initRes = await postJSON("/api/login/init", {
            user: params.userId,
            role: params.role,
            start_login_request: startLoginRequest,
        });

        if (initRes.status === 429) return rateLimited(initRes);
        if (!initRes.ok) return failure(await readError(initRes));

        const {login_response, login_handle} =
            await readJSON<LoginInitResponse>(initRes);

        const result = client.finishLogin({
            clientLoginState,
            loginResponse: login_response,
            password: params.password,
            ...identifierOpts(params),
        });

        if (!result) return {status: "wrong"};

        const finRes = await postJSON("/api/login/finish", {
            login_handle,
            finish_login_request: result.finishLoginRequest,
        });

        if (finRes.status === 429) return rateLimited(finRes);

        if (!finRes.ok) {
            const error = await readError(finRes);
            return error === "WRONG_CREDENTIALS" ? {status: "wrong"} : failure(error);
        }

        return ok;
    } catch {
        return failure();
    }
}

export async function opaqueRegister(params: OpaqueParams): Promise<OpaqueRegisterResult> {
    try {
        const client = await loadOpaque();

        const {clientRegistrationState, registrationRequest} =
            client.startRegistration({password: params.password});

        const initRes = await postJSON("/api/enroll/init", {
            enroll_token: params.enrollToken,
            start_enroll_request: registrationRequest,
        });

        if (initRes.status === 429) return rateLimited(initRes);
        if (!initRes.ok) return failure(await readError(initRes));

        const {enroll_response} =
            await readJSON<EnrollInitResponse>(initRes);

        const result = client.finishRegistration({
            clientRegistrationState,
            registrationResponse: enroll_response,
            password: params.password,
            ...identifierOpts(params),
        });

        const finishRes = await postJSON("/api/enroll/finish", {
            enroll_token: params.enrollToken,
            finish_enroll_request: result.registrationRecord,
        });

        if (finishRes.status === 429) return rateLimited(finishRes);
        if (!finishRes.ok) return failure(await readError(finishRes));

        return ok;
    } catch {
        return failure();
    }
}