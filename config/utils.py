import os
import hashlib
import secrets
from typing import Optional, Dict, Any
import json
from pathlib import Path


def generate_salt(length: int = 16) -> str:
    return secrets.token_hex(length)


def hash_password(password: str, salt: str) -> str:
    salted_password = password + salt
    return hashlib.sha256(salted_password.encode()).hexdigest()


def read_json_file(file_path: str) -> Optional[Dict[str, Any]]:
    try:
        with open(file_path, 'r') as f:
            return json.load(f)
    except (FileNotFoundError, json.JSONDecodeError):
        return None


def write_json_file(file_path: str, data: Dict[str, Any]) -> bool:
    try:
        with open(file_path, 'w') as f:
            json.dump(data, f, indent=2)
        return True
    except (IOError, TypeError):
        return False


def ensure_directory_exists(dir_path: str) -> bool:
    try:
        Path(dir_path).mkdir(parents=True, exist_ok=True)
        return True
    except OSError:
        return False


def validate_email(email: str) -> bool:
    import re
    pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
    return re.match(pattern, email) is not None


def get_env_variable(key: str, default: Optional[str] = None) -> Optional[str]:
    return os.getenv(key, default)