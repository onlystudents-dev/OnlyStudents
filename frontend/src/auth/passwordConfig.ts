export type PasswordConfig = {
    minLength: number;
    maxLength: number;
    hibpCheckEnabled: boolean;
};

const DEFAULT_CONFIG: PasswordConfig = {
    minLength: 12,
    maxLength: 256,
    hibpCheckEnabled: false,
};

function normalize(raw: unknown): PasswordConfig {
    const value = (raw ?? {}) as Partial<PasswordConfig>;

    return {
        minLength: typeof value.minLength === "number" && value.minLength > 0 ? value.minLength : DEFAULT_CONFIG.minLength,
        maxLength: typeof value.maxLength === "number" && value.maxLength > 0 ? value.maxLength : DEFAULT_CONFIG.maxLength,
        hibpCheckEnabled: value.hibpCheckEnabled === true,
    };
}

export async function getPasswordConfig(): Promise<PasswordConfig> {
    try {
        const response = await fetch("/api/config");
        if (!response.ok) return DEFAULT_CONFIG;

        const json = await response.json();
        return normalize(json?.password);
    } catch {
        return DEFAULT_CONFIG;
    }
}

export async function isPasswordExposed(password: string): Promise<boolean | null> {
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

        return body.split("\n").some(line => line.split(":")[0].trim() === suffix);
    } catch {
        return null;
    }
}
