/**
 * TypeScript type definitions for the VaultPack WASM module.
 */

export interface VaultpackError {
  ok: false;
  error: string;
}

export interface VaultpackSuccess {
  ok: true;
  manifest?: unknown;
  bundle?: Uint8Array;
  plaintext?: Uint8Array;
  generatedKey?: Uint8Array;
  shares?: Array<{ index: number; data: Uint8Array }>;
  signature?: Uint8Array;
  algorithm?: string;
  signedAt?: string;
  valid?: boolean;
  [k: string]: unknown;
}

export type VaultpackReply = VaultpackSuccess | VaultpackError;

export interface RecipientOptions {
  publicKeyPEM?: string;
  publicKeyPath?: string;
}

export interface ProtectOptions {
  plaintext: Uint8Array;
  key?: Uint8Array;
  password?: string;
  kdfAlgo?: string;
  cipher?: string;
  inputName?: string;
  aad?: Uint8Array;
  compress?: string;
  kmsProvider?: string;
  kmsKeyId?: string;
  recipients?: RecipientOptions[];
  splitShares?: number;
  splitThreshold?: number;
  hashAlgo?: string;
}

export interface DecryptOptions {
  bundle: Uint8Array;
  key?: Uint8Array;
  password?: string;
  aad?: Uint8Array;
  kmsProvider?: string;
  privateKeyPEM?: string;
  privateKeyPath?: string;
}

export interface SignOptions {
  bundle: Uint8Array;
  privateKeyPEM: string;
  algo?: string;
}

export interface VerifyOptions {
  bundle: Uint8Array;
  publicKeyPEM: string;
}

export interface Vaultpack {
  version(): string;
  protect(opts: ProtectOptions): VaultpackReply;
  decrypt(opts: DecryptOptions): VaultpackReply;
  inspect(bundle: Uint8Array): VaultpackReply;
  sign(opts: SignOptions): VaultpackReply;
  verify(opts: VerifyOptions): VaultpackReply;
}

declare global {
  // eslint-disable-next-line no-var
  var Vaultpack: Vaultpack;
}

export {};
