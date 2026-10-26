export type PasswordConfig = {
    min_length: number,
    max_length: number,
    hibp_check_enabled: boolean,
}

export type E2EEConfig = {
    argon2_salt_len: number,
    argon2_time: number,
    argon2_memory: number,
    argon2_threads: number,
    argon2_key_len: number,
    recovery_key_word_count: number,
}

export type FilesConfig = {
  max_file_size: number,
  allowed_exts: string[]
}

export type Config = {
    password: PasswordConfig,
    e2ee: E2EEConfig,
    files: FilesConfig,
}

export const DEFAULT_CONFIG: Config = {
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
  },
  files: {
    max_file_size: 10240,
    allowed_exts: [".png", ".docx", ".odt", ".7z", ".odp", ".mp4", ".m4a", ".txt", ".epub", ".gz", ".webm", ".gif", ".ogg", ".opus", ".csv", ".mp3", ".zip", ".tgz", ".xls", ".ppt", ".jpg", ".webp", ".wav", ".flac", ".aac", ".pdf", ".doc", ".md", ".tar", ".tar.xz", ".bz2", ".xlsx", ".bmp", ".rtf", ".rar", ".pptx", ".avi", ".mkv", ".tar.gz", ".xz", ".ods", ".jpeg"]
  }
}

export async function getPasswordConfig(): Promise<PasswordConfig> {
    try {
        const response = await fetch("/api/v1/config")
        if (!response.ok) return DEFAULT_CONFIG.password

        const json = await response.json() as {password: PasswordConfig}
        return json.password
    } catch {
        return DEFAULT_CONFIG.password
    }
}

export async function getE2EEConfig(): Promise<E2EEConfig> {
    try {
        const response = await fetch("/api/v1/config")
        if (!response.ok) return DEFAULT_CONFIG.e2ee

        const json = await response.json() as { e2ee: E2EEConfig }
        return json.e2ee
    } catch {
        return DEFAULT_CONFIG.e2ee
    }
}

export async function getFilesConfig(): Promise<FilesConfig> {
    try {
        const response = await fetch("/api/v1/config")
        if (!response.ok) return DEFAULT_CONFIG.files

        const json = await response.json() as { files: FilesConfig }
        return json.files
    } catch {
        return DEFAULT_CONFIG.files
    }
}


export async function isPasswordExposed(password: string): Promise<boolean | null> {
    // crypto.subtle is undefined on insecure (plain http) origins even though the DOM types say otherwise
    // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
    if (!password || !globalThis.crypto?.subtle) return null

    try {
        const digest = await globalThis.crypto.subtle.digest("SHA-1", new TextEncoder().encode(password))
        const hash = Array.from(new Uint8Array(digest))
            .map(byte => byte.toString(16).padStart(2, "0"))
            .join("")
            .toUpperCase()

        const response = await fetch(`https://api.pwnedpasswords.com/range/${hash.slice(0, 5)}`)
        if (!response.ok) return null

        const body = await response.text()
        const suffix = hash.slice(5)

        return body.split("\n").some(line => line.split(":")[0]?.trim() === suffix)
    } catch {
        return null
    }
}
