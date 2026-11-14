"""Type stubs for the Katame Python bindings (M232)."""

from typing import Any, Optional, Sequence

VERSION: str

class KatameError(RuntimeError): ...

def version() -> str: ...

def protect(
    *,
    input_path: Optional[str] = ...,
    plaintext: Optional[bytes] = ...,
    output_path: Optional[str] = ...,
    key: Optional[bytes] = ...,
    password: Optional[str] = ...,
    kdf_algo: Optional[str] = ...,
    cipher: Optional[str] = ...,
    aad: Optional[bytes] = ...,
    input_name: Optional[str] = ...,
    compress: Optional[str] = ...,
    kms_provider: Optional[str] = ...,
    kms_key_id: Optional[str] = ...,
    recipients: Optional[Sequence[dict[str, str]]] = ...,
    split_shares: Optional[int] = ...,
    split_threshold: Optional[int] = ...,
    hash_algo: Optional[str] = ...,
) -> dict[str, Any]: ...

def decrypt(
    *,
    input_path: Optional[str] = ...,
    bundle: Optional[bytes] = ...,
    output_path: Optional[str] = ...,
    key: Optional[bytes] = ...,
    key_b64: Optional[str] = ...,
    password: Optional[str] = ...,
    aad: Optional[bytes] = ...,
    kms_provider: Optional[str] = ...,
    private_key_pem: Optional[str] = ...,
    private_key_path: Optional[str] = ...,
) -> dict[str, Any]: ...

def inspect(
    path: Optional[str] = ...,
    *,
    bundle: Optional[bytes] = ...,
    input_path: Optional[str] = ...,
) -> dict[str, Any]: ...

def sign(
    *,
    bundle_path: str,
    private_key_pem: Optional[str] = ...,
    private_key_path: Optional[str] = ...,
    algo: Optional[str] = ...,
) -> dict[str, Any]: ...

def verify(
    *,
    bundle_path: str,
    public_key_pem: Optional[str] = ...,
    public_key_path: Optional[str] = ...,
) -> dict[str, Any]: ...
