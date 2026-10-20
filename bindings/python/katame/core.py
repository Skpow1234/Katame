"""Low-level ctypes wrapper around libkatame.

The library is located via the following resolution order:

1. ``$KATAME_LIB`` environment variable (absolute path).
2. A copy adjacent to this module (``katame/libkatame.{so,dylib,dll}``).
3. ``ctypes.util.find_library('katame')`` (system loader).
"""

from __future__ import annotations

import base64
import ctypes
import ctypes.util
import json
import os
import platform
import sys
from typing import Any, Optional, Sequence

VERSION = "1.0.0"  # aligned with pkg/katame.Version


class KatameError(RuntimeError):
    """Raised when an SDK call returns ``{"ok": false}`` or fails to load."""


def _libname() -> str:
    system = platform.system()
    if system == "Linux":
        return "libkatame.so"
    if system == "Darwin":
        return "libkatame.dylib"
    if system == "Windows":
        return "libkatame.dll"
    raise KatameError(f"unsupported platform: {system}")


def _candidate_paths() -> list[str]:
    paths: list[str] = []
    env = os.environ.get("KATAME_LIB")
    if env:
        paths.append(env)
    here = os.path.dirname(os.path.abspath(__file__))
    paths.append(os.path.join(here, _libname()))
    found = ctypes.util.find_library("katame")
    if found:
        paths.append(found)
    return paths


def _load_library() -> ctypes.CDLL:
    last_err: Optional[Exception] = None
    for path in _candidate_paths():
        if not path:
            continue
        try:
            return ctypes.CDLL(path)
        except OSError as exc:
            last_err = exc
    raise KatameError(
        "libkatame not found. Set KATAME_LIB=/path/to/libkatame."
        f"{_ext_for_msg()} or place the library next to katame/core.py. "
        f"Last loader error: {last_err}"
    )


def _ext_for_msg() -> str:
    return {"Linux": "so", "Darwin": "dylib", "Windows": "dll"}.get(
        platform.system(), "so"
    )


_LIB = _load_library()

_LIB.kt_version.argtypes = []
_LIB.kt_version.restype = ctypes.c_char_p
_LIB.kt_free.argtypes = [ctypes.c_void_p]
_LIB.kt_free.restype = None

for sym in ("kt_protect", "kt_decrypt", "kt_inspect", "kt_sign", "kt_verify", "kt_rewrap"):
    if hasattr(_LIB, sym):
        fn = getattr(_LIB, sym)
        fn.argtypes = [ctypes.c_char_p]
        fn.restype = ctypes.c_void_p


def _call_json(symbol: str, payload: Any) -> dict:
    if isinstance(payload, str):
        encoded = payload.encode("utf-8")
    else:
        encoded = json.dumps(payload, ensure_ascii=False).encode("utf-8")

    fn = getattr(_LIB, symbol)
    ptr = fn(encoded)
    if not ptr:
        raise KatameError(f"{symbol}: null reply from libkatame")

    try:
        raw = ctypes.string_at(ptr).decode("utf-8")
    finally:
        _LIB.kt_free(ptr)

    try:
        reply = json.loads(raw)
    except json.JSONDecodeError as exc:
        raise KatameError(f"{symbol}: invalid JSON reply: {exc}: {raw!r}")

    if not isinstance(reply, dict) or not reply.get("ok", False):
        msg = reply.get("error", "unknown error") if isinstance(reply, dict) else raw
        raise KatameError(f"{symbol}: {msg}")
    return reply


def version() -> str:
    raw = _LIB.kt_version()
    if not raw:
        raise KatameError("kt_version returned null")
    return raw.decode("utf-8")


def protect(
    *,
    input_path: Optional[str] = None,
    plaintext: Optional[bytes] = None,
    output_path: Optional[str] = None,
    key: Optional[bytes] = None,
    password: Optional[str] = None,
    kdf_algo: Optional[str] = None,
    cipher: Optional[str] = None,
    aad: Optional[bytes] = None,
    input_name: Optional[str] = None,
    compress: Optional[str] = None,
    kms_provider: Optional[str] = None,
    kms_key_id: Optional[str] = None,
    recipients: Optional[Sequence[dict[str, str]]] = None,
    split_shares: Optional[int] = None,
    split_threshold: Optional[int] = None,
    hash_algo: Optional[str] = None,
) -> dict:
    """Encrypt a file or in-memory blob into a .katame bundle.

    Supports symmetric key, password, hybrid recipients, KMS wrap, compression,
    and Shamir key splitting (parity with the Go SDK / CLI).
    """
    if (input_path is None) == (plaintext is None):
        raise KatameError("protect: exactly one of input_path or plaintext is required")
    if not output_path and plaintext is None and input_path is None:
        raise KatameError("protect: output_path is required")

    payload: dict[str, Any] = {}
    if output_path is not None:
        payload["output_path"] = output_path
    if input_path is not None:
        payload["input_path"] = input_path
    if plaintext is not None:
        payload["plaintext_b64"] = base64.b64encode(plaintext).decode("ascii")
    if key is not None:
        payload["key_b64"] = base64.b64encode(key).decode("ascii")
    if password is not None:
        payload["password"] = password
    if kdf_algo is not None:
        payload["kdf_algo"] = kdf_algo
    if cipher is not None:
        payload["cipher"] = cipher
    if aad is not None:
        payload["aad_b64"] = base64.b64encode(aad).decode("ascii")
    if input_name is not None:
        payload["input_name"] = input_name
    if compress is not None:
        payload["compress"] = compress
    if kms_provider is not None:
        payload["kms_provider"] = kms_provider
    if kms_key_id is not None:
        payload["kms_key_id"] = kms_key_id
    if recipients is not None:
        payload["recipients"] = list(recipients)
    if split_shares is not None:
        payload["split_shares"] = split_shares
    if split_threshold is not None:
        payload["split_threshold"] = split_threshold
    if hash_algo is not None:
        payload["hash_algo"] = hash_algo
    return _call_json("kt_protect", payload)


def decrypt(
    *,
    input_path: Optional[str] = None,
    bundle: Optional[bytes] = None,
    output_path: Optional[str] = None,
    key: Optional[bytes] = None,
    key_b64: Optional[str] = None,
    password: Optional[str] = None,
    aad: Optional[bytes] = None,
    kms_provider: Optional[str] = None,
    private_key_pem: Optional[str] = None,
    private_key_path: Optional[str] = None,
) -> dict:
    """Decrypt a .katame bundle (symmetric, password, hybrid, or KMS)."""
    if (input_path is None) == (bundle is None):
        raise KatameError("decrypt: exactly one of input_path or bundle is required")

    payload: dict[str, Any] = {}
    if input_path is not None:
        payload["input_path"] = input_path
    if bundle is not None:
        payload["bundle_b64"] = base64.b64encode(bundle).decode("ascii")
    if output_path is not None:
        payload["output_path"] = output_path
    if key is not None:
        payload["key_b64"] = base64.b64encode(key).decode("ascii")
    elif key_b64 is not None:
        payload["key_b64"] = key_b64
    if password is not None:
        payload["password"] = password
    if aad is not None:
        payload["aad_b64"] = base64.b64encode(aad).decode("ascii")
    if kms_provider is not None:
        payload["kms_provider"] = kms_provider
    if private_key_pem is not None:
        payload["private_key_pem"] = private_key_pem
    if private_key_path is not None:
        payload["private_key_path"] = private_key_path
    return _call_json("kt_decrypt", payload)


def inspect(
    path: Optional[str] = None,
    *,
    bundle: Optional[bytes] = None,
    input_path: Optional[str] = None,
) -> dict:
    """Return the parsed manifest of a .katame bundle without decrypting it."""
    if path is not None and bundle is None and input_path is None:
        # Backward-compatible path-only call (raw string still accepted by C ABI).
        return _call_json("kt_inspect", path)
    payload: dict[str, Any] = {}
    if input_path or path:
        payload["input_path"] = input_path or path
    if bundle is not None:
        payload["bundle_b64"] = base64.b64encode(bundle).decode("ascii")
    return _call_json("kt_inspect", payload)


def sign(
    *,
    bundle_path: str,
    private_key_pem: Optional[str] = None,
    private_key_path: Optional[str] = None,
    algo: Optional[str] = None,
) -> dict:
    if (private_key_pem is None) == (private_key_path is None):
        raise KatameError(
            "sign: exactly one of private_key_pem or private_key_path is required"
        )
    payload: dict[str, Any] = {"bundle_path": bundle_path}
    if private_key_pem is not None:
        payload["private_key_pem"] = private_key_pem
    if private_key_path is not None:
        payload["private_key_path"] = private_key_path
    if algo is not None:
        payload["algo"] = algo
    return _call_json("kt_sign", payload)


def verify(
    *,
    bundle_path: str,
    public_key_pem: Optional[str] = None,
    public_key_path: Optional[str] = None,
) -> dict:
    if (public_key_pem is None) == (public_key_path is None):
        raise KatameError(
            "verify: exactly one of public_key_pem or public_key_path is required"
        )
    payload: dict[str, Any] = {"bundle_path": bundle_path}
    if public_key_pem is not None:
        payload["public_key_pem"] = public_key_pem
    if public_key_path is not None:
        payload["public_key_path"] = public_key_path
    return _call_json("kt_verify", payload)


def rewrap(
    *,
    input_path: str,
    output_path: str,
    kms_provider: str,
    from_kms_key_id: str,
    to_kms_key_id: str,
) -> dict:
    """Re-wrap a KMS-protected DEK without decrypting the payload."""
    return _call_json(
        "kt_rewrap",
        {
            "input_path": input_path,
            "output_path": output_path,
            "kms_provider": kms_provider,
            "from_kms_key_id": from_kms_key_id,
            "to_kms_key_id": to_kms_key_id,
        },
    )


if __name__ == "__main__":  # pragma: no cover
    print("libkatame version:", version(), "(python wrapper", VERSION + ")")
    sys.exit(0)
