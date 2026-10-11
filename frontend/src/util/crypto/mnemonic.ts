import {generateMnemonic, mnemonicToSeedSync, validateMnemonic} from '@scure/bip39';
import { wordlist } from '@scure/bip39/wordlists/english.js';

export function generate(length: number) {
  return generateMnemonic(wordlist, length)
}

export function validate(mnemonic: string) {
  return validateMnemonic(mnemonic, wordlist)
}

export function normalize(user_provided: string) {
  return user_provided.trim().toLowerCase().replace(/\s/, '')
}

export function toSeed(mnemonic: string) {
  return mnemonicToSeedSync(mnemonic)
}
