import {generateMnemonic, mnemonicToSeedSync, validateMnemonic} from '@scure/bip39';
import { wordlist } from '@scure/bip39/wordlists/english.js';

const STRENGTH_BITS: Record<number, number> = {
  12: 128,
  15: 160,
  18: 192,
  21: 224,
  24: 256,
}

export function generateMnemonicSecret(length: number) {
  return generateMnemonic(wordlist, STRENGTH_BITS[length] ?? 256)
}

export function validateMnemonicSecret(mnemonic: string) {
  return validateMnemonic(mnemonic, wordlist)
}

export function normalizeMnemonicSecret(user_provided: string) {
  return user_provided.trim().toLowerCase().replace(/\s+/g, '')
}

export function MnemonicSecrettoSeed(mnemonic: string) {
  return mnemonicToSeedSync(mnemonic)
}
