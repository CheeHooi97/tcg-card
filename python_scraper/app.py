import asyncio
import random
import time
from datetime import datetime, timezone
from typing import Any
from pathlib import Path

from fastapi import FastAPI
from pydantic import BaseModel, Field
from playwright.async_api import async_playwright, TimeoutError as PWTimeoutError
import psycopg
from dotenv import load_dotenv

load_dotenv()
load_dotenv(dotenv_path=Path(__file__).resolve().parents[1] / ".env")

app = FastAPI(title="TCG Python Scraper")


class BGSRequest(BaseModel):
    url: list[str] = Field(default_factory=list)


class CGCRequest(BaseModel):
    url: str


class PSARequest(BaseModel):
    urls: list[str] = Field(default_factory=list)


class TAGRequest(BaseModel):
    urls: list[str] = Field(default_factory=list)


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
    import os
    host = os.getenv("POSTGRES_HOST", "127.0.0.1")
    port = int(os.getenv("POSTGRES_PORT", "5432"))
    user = os.getenv("POSTGRES_USER", "").strip()
    password = os.getenv("POSTGRES_PASSWORD", "")
    dbname = os.getenv("POSTGRES_DATABASE", "").strip()
    sslmode = os.getenv("POSTGRES_SSLMODE", "disable").strip() or "disable"

    missing = []
    if not user:
        missing.append("POSTGRES_USER")
    if not dbname:
        missing.append("POSTGRES_DATABASE")
    if missing:
        raise RuntimeError(f"missing required DB env: {', '.join(missing)}")

    return {
        "host": host,
        "port": port,
        "user": user,
        "password": password,
        "dbname": dbname,
        "sslmode": sslmode,
    }


def split_set_number(card_number: str) -> tuple[str, str]:
    parts = card_number.split("/")
    if len(parts) == 2:
        return parts[0].strip(), parts[1].strip()
    return card_number.strip(), ""


def save_psa_card(set_name: str, set_link: str, card: dict[str, Any]) -> None:
    now = datetime.now(timezone.utc)
    card_number, set_number = split_set_number(card.get("cardNumber", ""))
    card_name = (card.get("cardName") or "").strip()
    description = (card.get("description") or "").strip()

    payload = {
        "card_name": card_name,
        "card_number": card_number,
        "set_number": set_number,
        "set_name": set_name,
        "description": description,
        "total": str(card.get("total", "0")),
        "grade1": str(card.get("grade1", "0")),
        "grade2": str(card.get("grade2", "0")),
        "grade3": str(card.get("grade3", "0")),
        "grade4": str(card.get("grade4", "0")),
        "grade5": str(card.get("grade5", "0")),
        "grade6": str(card.get("grade6", "0")),
        "grade7": str(card.get("grade7", "0")),
        "grade8": str(card.get("grade8", "0")),
        "grade9": str(card.get("grade9", "0")),
        "grade10": str(card.get("grade10", "0")),
    }

    with psycopg.connect(**db_config()) as conn:
        with conn.cursor() as cur:
            cur.execute(
                """
                SELECT id, set_name, set_number, rarity, spec_id, auth
                FROM psa
                WHERE card_name=%s AND card_number=%s AND description=%s AND set_name=%s
                LIMIT 1
                """,
                (payload["card_name"], payload["card_number"], payload["description"], payload["set_name"]),
            )
            row = cur.fetchone()

            psa_id = ""
            if row is None:
                psa_id = unique_id()
                cur.execute(
                    """
                    INSERT INTO psa (
                        id, card_name, card_number, set_number, set_name, rarity, description, spec_id, total, auth,
                        grade1, grade2, grade3, grade4, grade5, grade6, grade7, grade8, grade9, grade10,
                        created_date_time, updated_date_time
                    ) VALUES (
                        %s, %s, %s, %s, %s, '', %s, '', %s, '',
                        %s, %s, %s, %s, %s, %s, %s, %s, %s, %s,
                        %s, %s
                    )
                    """,
                    (
                        psa_id, payload["card_name"], payload["card_number"], payload["set_number"], payload["set_name"],
                        payload["description"], payload["total"],
                        payload["grade1"], payload["grade2"], payload["grade3"], payload["grade4"], payload["grade5"],
                        payload["grade6"], payload["grade7"], payload["grade8"], payload["grade9"], payload["grade10"],
                        now, now,
                    ),
                )
            else:
                psa_id = row[0] or ""

            logging_id = unique_id()
            logging_set_name = payload["set_name"]
            logging_set_number = payload["set_number"]
            logging_rarity = ""
            logging_spec_id = ""
            logging_auth = ""

            if row is not None:
                if row[1]:
                    logging_set_name = row[1]
                if row[2]:
                    logging_set_number = row[2]
                if row[3]:
                    logging_rarity = row[3]
                if row[4]:
                    logging_spec_id = row[4]
                if row[5]:
                    logging_auth = row[5]

            cur.execute(
                """
                INSERT INTO psa_logging (
                    id, psa_id, card_name, card_number, set_number, set_name, rarity, description, spec_id, total, auth,
                    grade1, grade2, grade3, grade4, grade5, grade6, grade7, grade8, grade9, grade10,
                    created_date_time, updated_date_time
                ) VALUES (
                    %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s,
                    %s, %s, %s, %s, %s, %s, %s, %s, %s, %s,
                    %s, %s
                )
                """,
                (
                    logging_id, psa_id, payload["card_name"], payload["card_number"], logging_set_number, logging_set_name,
                    logging_rarity, payload["description"], logging_spec_id, payload["total"], logging_auth,
                    payload["grade1"], payload["grade2"], payload["grade3"], payload["grade4"], payload["grade5"],
                    payload["grade6"], payload["grade7"], payload["grade8"], payload["grade9"], payload["grade10"],
                    now, now,
                ),
            )

            cur.execute("SELECT 1 FROM psa_url WHERE url=%s LIMIT 1", (set_link,))
            exists_url = cur.fetchone() is not None
            if not exists_url:
                cur.execute(
                    """
                    INSERT INTO psa_url (id, set_name, url, created_date_time, updated_date_time)
                    VALUES (%s, %s, %s, %s, %s)
                    """,
                    (unique_id(), set_name, set_link, now, now),
                )
        conn.commit()


@app.get("/health")
async def health():
    return {"ok": True}


@app.post("/scrape/bgs")
async def scrape_bgs(req: BGSRequest):
    result: list[dict[str, Any]] = []
    failed: list[dict[str, Any]] = []

    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True, args=["--no-sandbox", "--disable-dev-shm-usage"])

        for set_name in req.url:
            try:
                async def run_one():
                    ctx, page = await new_page(browser)
                    try:
                        await page.goto("https://www.beckett.com/grading/pop-report", timeout=60000)
                        await page.wait_for_selector("#set_name", timeout=30000)
                        await page.fill("#set_name", set_name)
                        await page.keyboard.press("Enter")
                        await page.wait_for_timeout(2000)
                        links = await page.eval_on_selector_all(
                            'a[href*="/set_match/"]',
                            "els => els.map(a => ({url: a.href, setName: (a.innerText||'').trim()}))",
                        )
                        return links
                    finally:
                        await ctx.close()

                links = await with_retry(f"bgs set {set_name}", 3, run_one)
                result.append({"input": set_name, "links": links})
            except Exception as e:
                failed.append({"input": set_name, "error": str(e)})
                continue

        await browser.close()

    return {"ok": True, "count": len(result), "failedCount": len(failed), "result": result, "failed": failed}


@app.post("/scrape/cgc")
async def scrape_cgc(req: CGCRequest):
    result: list[dict[str, Any]] = []
    failed: list[dict[str, Any]] = []

    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True, args=["--no-sandbox", "--disable-dev-shm-usage"])

        try:
            async def run_root():
                ctx, page = await new_page(browser)
                try:
                    await page.goto(req.url, timeout=60000)
                    await page.wait_for_selector(".ccg-cards", timeout=45000)
                    lists = await page.evaluate(
                        """
                        () => {
                          const baseUrl = "https://www.cgccards.com";
                          return Array.from(document.querySelectorAll(".card.ng-scope a"))
                            .map(a => ({url: baseUrl + a.getAttribute("href")}))
                            .filter(x => x.url);
                        }
                        """
                    )
                    return lists
                finally:
                    await ctx.close()

            lists = await with_retry("cgc list", 3, run_root)
            for item in lists:
                result.append(item)
        except Exception as e:
            failed.append({"input": req.url, "error": str(e)})

        await browser.close()

    return {"ok": True, "count": len(result), "failedCount": len(failed), "result": result, "failed": failed}


@app.post("/scrape/psa")
async def scrape_psa(req: PSARequest):
    result: list[dict[str, Any]] = []
    failed: list[dict[str, Any]] = []
    inserted_cards = 0
    logged_cards = 0

    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True, args=["--no-sandbox", "--disable-dev-shm-usage"])

        for year_url in req.urls:
            try:
                async def scrape_year_sets():
                    ctx, page = await new_page(browser)
                    try:
                        await page.goto(year_url, timeout=60000)
                        await page.wait_for_selector("#tableSets", timeout=45000)
                        await page.wait_for_timeout(2000)
                        rows = await page.evaluate(
                            """
                            () => {
                              const base = "https://www.psacard.com";
                              return Array.from(document.querySelectorAll("#tableSets tbody tr"))
                                .map(tr => {
                                  const a = tr.querySelector("td.text-left a:not([href='#'])");
                                  return a ? {name: a.innerText.trim(), link: base + a.getAttribute("href")} : null;
                                })
                                .filter(Boolean);
                            }
                            """
                        )
                        return rows
                    finally:
                        await ctx.close()

                sets = await with_retry(f"psa year {year_url}", 3, scrape_year_sets)
                pokemon_sets = [s for s in sets if "pokemon" in (s.get("name", "").lower())]

                year_summary = {
                    "input": year_url,
                    "totalSets": len(sets),
                    "pokemonSets": len(pokemon_sets),
                    "processedSets": 0,
                    "failedSets": 0,
                }

                for s in pokemon_sets:
                    set_name = s.get("name", "").strip()
                    set_link = s.get("link", "").strip()
                    if not set_name or not set_link:
                        continue
                    try:
                        async def scrape_set_cards():
                            ctx, page = await new_page(browser)
                            try:
                                await page.goto(set_link, timeout=60000)
                                await page.wait_for_selector("#tablePSA", timeout=45000)
                                await page.wait_for_timeout(2000)
                                cards = await page.evaluate(
                                    """
                                    () => {
                                      const rows = Array.from(document.querySelectorAll('#tablePSA tbody tr[role="row"]'));
                                      const getGradeVal = (cell) => {
                                        const d = cell ? cell.querySelector('div') : null;
                                        return d ? d.innerText.trim() : "0";
                                      };
                                      return rows.map(tr => {
                                        const cells = Array.from(tr.querySelectorAll('td'));
                                        if (cells.length < 17 || cells[2].innerText.includes("TOTAL POPULATION")) return null;
                                        const nameElem = cells[2].querySelector('strong');
                                        const cardName = nameElem ? nameElem.innerText.trim() : "";
                                        const clone = cells[2].cloneNode(true);
                                        const a = clone.querySelector('a');
                                        const strong = clone.querySelector('strong');
                                        if (a) a.remove();
                                        if (strong) strong.remove();
                                        return {
                                          cardNumber: cells[1].innerText.trim(),
                                          cardName: cardName,
                                          description: clone.innerText.trim(),
                                          grade1: getGradeVal(cells[6]),
                                          grade2: getGradeVal(cells[7]),
                                          grade3: getGradeVal(cells[8]),
                                          grade4: getGradeVal(cells[9]),
                                          grade5: getGradeVal(cells[10]),
                                          grade6: getGradeVal(cells[11]),
                                          grade7: getGradeVal(cells[12]),
                                          grade8: getGradeVal(cells[13]),
                                          grade9: getGradeVal(cells[14]),
                                          grade10: getGradeVal(cells[15]),
                                          total: getGradeVal(cells[16]),
                                        };
                                      }).filter(Boolean);
                                    }
                                    """
                                )
                                return cards
                            finally:
                                await ctx.close()

                        cards = await with_retry(f"psa set {set_name}", 3, scrape_set_cards)
                        for card in cards:
                            try:
                                await asyncio.to_thread(save_psa_card, set_name, set_link, card)
                                inserted_cards += 1
                                logged_cards += 1
                            except Exception as db_err:
                                failed.append({
                                    "input": year_url,
                                    "set": set_name,
                                    "card": f'{card.get("cardNumber", "")} {card.get("cardName", "")}'.strip(),
                                    "error": f"db save failed: {db_err}",
                                })
                                continue
                        year_summary["processedSets"] += 1
                    except Exception as set_err:
                        year_summary["failedSets"] += 1
                        failed.append({"input": year_url, "set": set_name, "error": str(set_err)})
                        continue
                result.append(year_summary)
            except (Exception, PWTimeoutError) as e:
                failed.append({"input": year_url, "error": str(e)})
                continue

        await browser.close()

    return {
        "ok": True,
        "count": len(result),
        "failedCount": len(failed),
        "insertedCards": inserted_cards,
        "loggedCards": logged_cards,
        "result": result,
        "failed": failed,
    }


@app.post("/scrape/tag")
async def scrape_tag(req: TAGRequest):
    result: list[dict[str, Any]] = []
    failed: list[dict[str, Any]] = []

    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True, args=["--no-sandbox", "--disable-dev-shm-usage"])

        for year_url in req.urls:
            try:
                async def run_one():
                    ctx, page = await new_page(browser)
                    try:
                        await page.goto(year_url, timeout=60000)
                        await page.wait_for_load_state("domcontentloaded")
                        await page.wait_for_timeout(3000)
                        sets = await page.evaluate(
                            """
                            () => {
                              const baseUrl = window.location.origin || "https://my.taggrading.com";
                              const seen = new Set();
                              const out = [];
                              const links = Array.from(document.querySelectorAll('a[href*="/pop-report/"]'));
                              for (const a of links) {
                                const href = a.getAttribute("href");
                                if (!href) continue;
                                const u = new URL(href, baseUrl).href;
                                if (seen.has(u)) continue;
                                seen.add(u);
                                out.push({link: u, name: (a.textContent || "").replace(/\\s+/g, " ").trim()});
                              }
                              return out;
                            }
                            """
                        )
                        return sets
                    finally:
                        await ctx.close()

                sets = await with_retry(f"tag year {year_url}", 3, run_one)
                result.append({"input": year_url, "sets": sets})
            except Exception as e:
                failed.append({"input": year_url, "error": str(e)})
                continue

        await browser.close()

    return {"ok": True, "count": len(result), "failedCount": len(failed), "result": result, "failed": failed}
