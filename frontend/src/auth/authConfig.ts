import {readJSON} from "../util/api.ts";

export type PasswordConfig = {
    min_length: number;
    max_length: number;
    hibp_check_enabled: boolean;
};

export type E2EEConfig = {
    argon2_salt_len: number;
    argon2_time: number;
    argon2_memory: number;
    argon2_threads: number;
    argon2_key_len: number;
    recovery_key_word_count: number;
};

export type AuthConfig = {
    password: PasswordConfig;
    e2ee: E2EEConfig;
};

export const AUTH_DEFAULT_CONFIG: AuthConfig = {
  password: {
    min_length: 12,
    max_length: 256,
    hibp_check_enabled: false,
  },
  e2ee: {
    argon2_salt_len: 16,
    argon2_time: 2,
    argon2_memory: 65536,
    argon2_threads: 1,
    argon2_key_len: 32,
    recovery_key_word_count: 24,
  }
};

function normalize_password(raw: unknown): PasswordConfig {
    const value = (raw ?? {}) as Partial<PasswordConfig>;

    return {
        min_length: typeof value.min_length === "number" && value.min_length > 0 ? value.min_length : AUTH_DEFAULT_CONFIG.password.min_length,
        max_length: typeof value.max_length === "number" && value.max_length > 0 ? value.max_length : AUTH_DEFAULT_CONFIG.password.max_length,
        hibp_check_enabled: value.hibp_check_enabled === true,
    };
}

function normalize_e2ee(raw: unknown): E2EEConfig {
    const value = (raw ?? {}) as Partial<E2EEConfig>;

    return {
        argon2_salt_len: typeof value.argon2_salt_len === "number" && value.argon2_salt_len > 0 ? value.argon2_salt_len : AUTH_DEFAULT_CONFIG.e2ee.argon2_salt_len,
        argon2_time: typeof value.argon2_time === "number" && value.argon2_time > 0 ? value.argon2_time : AUTH_DEFAULT_CONFIG.e2ee.argon2_time,
        argon2_memory: typeof value.argon2_memory === "number" && value.argon2_memory > 0 ? value.argon2_memory : AUTH_DEFAULT_CONFIG.e2ee.argon2_memory,
        argon2_threads: typeof value.argon2_threads === "number" && value.argon2_threads > 0 ? value.argon2_threads : AUTH_DEFAULT_CONFIG.e2ee.argon2_threads,
        argon2_key_len: typeof value.argon2_key_len === "number" && value.argon2_key_len > 0 ? value.argon2_key_len : AUTH_DEFAULT_CONFIG.e2ee.argon2_key_len,
        recovery_key_word_count: typeof value.recovery_key_word_count === "number" && value.recovery_key_word_count > 0 ? value.recovery_key_word_count : AUTH_DEFAULT_CONFIG.e2ee.recovery_key_word_count,
    };
}

export async function getPasswordConfig(): Promise<PasswordConfig> {
    try {
        const response = await fetch("/api/v1/config");
        if (!response.ok) return AUTH_DEFAULT_CONFIG.password;

        const json = await readJSON<{ password?: unknown } | null>(response);
        return normalize_password(json?.password);
    } catch {
        return AUTH_DEFAULT_CONFIG.password;
    }
}

export async function getE2EEConfig(): Promise<E2EEConfig> {
    try {
        const response = await fetch("/api/v1/config");
        if (!response.ok) return AUTH_DEFAULT_CONFIG.e2ee;

        const json = await readJSON<{ e2ee?: unknown } | null>(response);
        return normalize_e2ee(json?.e2ee);
    } catch {
        return AUTH_DEFAULT_CONFIG.e2ee;
    }
}

export async function isPasswordExposed(password: string): Promise<boolean | null> {
    // crypto.subtle is undefined on insecure (plain http) origins even though the DOM types say otherwise
    // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
    if (!password || !globalThis.crypto?.subtle) return null;

    try {
        const digest = await globalThis.crypto.subtle.digest("SHA-1", new TextEncoder().encode(password));
        const hash = Array.from(new Uint8Array(digest))
            .map(byte => byte.toString(16).padStart(2, "0"))
            .join("")
            .toUpperCase();

        const response = await fetch(`https://api.pwnedpasswords.com/range/${hash.slice(0, 5)}`);
        if (!response.ok) return null;

        const body = await response.text();
        const suffix = hash.slice(5);

        return body.split("\n").some(line => line.split(":")[0]?.trim() === suffix);
    } catch {
        return null;
    }
}
