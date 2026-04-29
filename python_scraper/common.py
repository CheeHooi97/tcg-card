import asyncio
import os
import random
import time
from pathlib import Path
from typing import Any

from dotenv import load_dotenv

load_dotenv()
load_dotenv(dotenv_path=Path(__file__).resolve().parents[1] / ".env")


async def with_retry(name: str, attempts: int, fn):
    last_err = None
    for i in range(1, attempts + 1):
        try:
            return await fn()
        except Exception as e:
            last_err = e
            print(f"[{name}] attempt {i}/{attempts} failed: {e}")
            if i < attempts:
                await asyncio.sleep(2)
    raise last_err


async def new_page(browser):
    context = await browser.new_context(
        user_agent="Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
        viewport={"width": 1920, "height": 1080},
    )
    return context, await context.new_page()


def unique_id() -> str:
    return f"{int(time.time())}{random.randint(100000, 999999)}"


def db_config() -> dict[str, Any]:
    host = os.getenv("POSTGRES_HOST", "127.0.0.1")
    port = int(os.getenv("POSTGRES_PORT", "5432"))
    user = os.getenv("POSTGRES_USER", "").strip()
    password = os.getenv("POSTGRES_PASSWORD", "")
    dbname = os.getenv("POSTGRES_DATABASE", "").strip()
    sslmode = os.getenv("POSTGRES_SSLMODE", "disable").strip() or "disable"
    if not user or not dbname:
        raise RuntimeError("missing required DB env: POSTGRES_USER / POSTGRES_DATABASE")
    return {"host": host, "port": port, "user": user, "password": password, "dbname": dbname, "sslmode": sslmode}

