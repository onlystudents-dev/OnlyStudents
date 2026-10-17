import * as argon2 from '@noble/hashes/argon2.js';
import * as hkdf from '@noble/hashes/hkdf.js';
import * as sha2 from '@noble/hashes/sha2.js';
import { AUTH_DEFAULT_CONFIG, type E2EEConfig } from '../../auth/authConfig';
import { normalizeMnemonicSecret } from './mnemonic';
import { concatBytes, utf8ToBytes } from '@noble/hashes/utils.js';

export type Bytes = Uint8Array<ArrayBuffer>;

export function generateMK(): Bytes {
  return randomBytes(32)
}

export function deriveWrapKey(secret: Bytes, purpose: "password" | "recovery"): Promise<CryptoKey> {
  return crypto.subtle.importKey("raw", hkdf.hkdf(sha2.sha256, secret, utf8ToBytes(`onlystudents-e2ee-v1`), utf8ToBytes(`onlystudents-${purpose}-v1`), 32), "AES-GCM", false, ["encrypt", "decrypt"])
}

export function randomBytes(len: number): Bytes {
  return crypto.getRandomValues(new Uint8Array(len))
}

const CHUNK = 0x8000;

export function b64encode(u8: Uint8Array): string {
  let s = "";
  for (let i = 0; i < u8.length; i += CHUNK) {
    s += String.fromCharCode(...u8.subarray(i, i + CHUNK))
  }
  return btoa(s)
}

export function recoveryDigest(mnemonic: string, salt: Bytes, cfg: E2EEConfig): Bytes {
  return argon2.argon2id(normalizeMnemonicSecret(mnemonic), salt, {
    t: cfg.argon2_time,
    m: cfg.argon2_memory,
    p: cfg.argon2_threads,
    dkLen: cfg.argon2_key_len
  })
}

export function paramsToString(cfg: E2EEConfig): string {
    return `t=${cfg.argon2_time},m=${cfg.argon2_memory},p=${cfg.argon2_threads},dkLen=${cfg.argon2_key_len}`
}

export function parseParams(s: string): {t: number, m: number, p: number, dkLen: number} {
    const out: Record<string, number> = {}
    for (const part of s.split(",")) {
        const [k, v] = part.split("=")
        if (typeof k == "undefined") {
          continue
        }
        out[k] = Number(v)
    }
    return {t: out.t ? out.t : AUTH_DEFAULT_CONFIG.e2ee.argon2_time, m: out.m ? out.m : AUTH_DEFAULT_CONFIG.e2ee.argon2_memory, p: out.p ? out.p : AUTH_DEFAULT_CONFIG.e2ee.argon2_threads, dkLen: out.dkLen ? out.dkLen : AUTH_DEFAULT_CONFIG.e2ee.argon2_key_len}
}

function base64UrlDecode(str: string): Bytes {
  const b64 = str.replace(/-/g, '+').replace(/_/g, '/')
  const bin = atob(b64.padEnd(Math.ceil(b64.length / 4) * 4, '='))
  return Uint8Array.from(bin, c => c.charCodeAt(0))
}

export function decodeOpaqueExportKey(exportKey: string): Bytes {
  return base64UrlDecode(exportKey)
}

export async function wrapMK(mkBytes: Bytes, aesKey: CryptoKey): Promise<Bytes> {
  const iv = randomBytes(12)
  return concatBytes(iv,
    new Uint8Array(await crypto.subtle.encrypt({
      "name": "AES-GCM",
      iv,
    },
    aesKey, mkBytes)
  ))
}

export async function unwrapMK(blob: Bytes, aesKey: CryptoKey): Promise<Bytes> {
  const iv = blob.slice(0, 12)

  return new Uint8Array(await crypto.subtle.decrypt({
    "name": "AES-GCM",
    iv,
  }, aesKey, blob.slice(12)))
}
