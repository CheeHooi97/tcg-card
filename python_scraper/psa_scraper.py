import asyncio
from datetime import datetime, timezone
from typing import Any

import psycopg
from playwright.async_api import TimeoutError as PWTimeoutError

from common import db_config, new_page, unique_id, with_retry


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
    p = {
        "card_name": card_name, "card_number": card_number, "set_number": set_number, "set_name": set_name,
        "description": description, "total": str(card.get("total", "0")),
        "grade1": str(card.get("grade1", "0")), "grade2": str(card.get("grade2", "0")), "grade3": str(card.get("grade3", "0")),
        "grade4": str(card.get("grade4", "0")), "grade5": str(card.get("grade5", "0")), "grade6": str(card.get("grade6", "0")),
        "grade7": str(card.get("grade7", "0")), "grade8": str(card.get("grade8", "0")), "grade9": str(card.get("grade9", "0")),
        "grade10": str(card.get("grade10", "0")),
    }
    with psycopg.connect(**db_config()) as conn:
        with conn.cursor() as cur:
            cur.execute("SELECT id,set_name,set_number,rarity,spec_id,auth FROM psa WHERE card_name=%s AND card_number=%s AND description=%s AND set_name=%s LIMIT 1", (p["card_name"], p["card_number"], p["description"], p["set_name"]))
            row = cur.fetchone()
            psa_id = row[0] if row else unique_id()
            if not row:
                cur.execute("""INSERT INTO psa (id,card_name,card_number,set_number,set_name,rarity,description,spec_id,total,auth,grade1,grade2,grade3,grade4,grade5,grade6,grade7,grade8,grade9,grade10,created_date_time,updated_date_time)
                               VALUES (%s,%s,%s,%s,%s,'',%s,'',%s,'',%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)""",
                            (psa_id,p["card_name"],p["card_number"],p["set_number"],p["set_name"],p["description"],p["total"],p["grade1"],p["grade2"],p["grade3"],p["grade4"],p["grade5"],p["grade6"],p["grade7"],p["grade8"],p["grade9"],p["grade10"],now,now))
            logging_set_name = row[1] if row and row[1] else p["set_name"]
            logging_set_number = row[2] if row and row[2] else p["set_number"]
            logging_rarity = row[3] if row and row[3] else ""
            logging_spec_id = row[4] if row and row[4] else ""
            logging_auth = row[5] if row and row[5] else ""
            cur.execute("""INSERT INTO psa_logging (id,psa_id,card_name,card_number,set_number,set_name,rarity,description,spec_id,total,auth,grade1,grade2,grade3,grade4,grade5,grade6,grade7,grade8,grade9,grade10,created_date_time,updated_date_time)
                           VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)""",
                        (unique_id(), psa_id, p["card_name"], p["card_number"], logging_set_number, logging_set_name, logging_rarity, p["description"], logging_spec_id, p["total"], logging_auth, p["grade1"], p["grade2"], p["grade3"], p["grade4"], p["grade5"], p["grade6"], p["grade7"], p["grade8"], p["grade9"], p["grade10"], now, now))
            cur.execute("SELECT 1 FROM psa_url WHERE url=%s LIMIT 1", (set_link,))
            if cur.fetchone() is None:
                cur.execute("INSERT INTO psa_url (id,set_name,url,created_date_time,updated_date_time) VALUES (%s,%s,%s,%s,%s)", (unique_id(), set_name, set_link, now, now))
        conn.commit()


async def scrape_psa(browser, req_urls: list[str]):
    result, failed = [], []
    inserted_cards, logged_cards = 0, 0
    for year_url in req_urls:
        try:
            async def scrape_year_sets():
                ctx, page = await new_page(browser)
                try:
                    await page.goto(year_url, timeout=60000)
                    await page.wait_for_selector("#tableSets", timeout=45000)
                    return await page.evaluate("""() => { const base="https://www.psacard.com"; return Array.from(document.querySelectorAll("#tableSets tbody tr")).map(tr=>{const a=tr.querySelector("td.text-left a:not([href='#'])"); return a?{name:a.innerText.trim(),link:base+a.getAttribute("href")}:null;}).filter(Boolean);}""")
                finally:
                    await ctx.close()
            sets = await with_retry(f"psa year {year_url}", 3, scrape_year_sets)
            pokemon_sets = [s for s in sets if "pokemon" in (s.get("name", "").lower())]
            summary = {"input": year_url, "totalSets": len(sets), "pokemonSets": len(pokemon_sets), "processedSets": 0, "failedSets": 0}
            for s in pokemon_sets:
                set_name, set_link = s.get("name", "").strip(), s.get("link", "").strip()
                if not set_name or not set_link:
                    continue
                try:
                    async def scrape_set_cards():
                        ctx, page = await new_page(browser)
                        try:
                            await page.goto(set_link, timeout=60000)
                            await page.wait_for_selector("#tablePSA", timeout=45000)
                            return await page.evaluate("""() => { const rows=Array.from(document.querySelectorAll('#tablePSA tbody tr[role="row"]')); const gv=(c)=>{const d=c?c.querySelector('div'):null; return d?d.innerText.trim():"0"}; return rows.map(tr=>{const t=Array.from(tr.querySelectorAll('td')); if(t.length<17||t[2].innerText.includes("TOTAL POPULATION")) return null; const nm=t[2].querySelector('strong'); const cl=t[2].cloneNode(true); const a=cl.querySelector('a'); const s=cl.querySelector('strong'); if(a)a.remove(); if(s)s.remove(); return {cardNumber:t[1].innerText.trim(),cardName:nm?nm.innerText.trim():"",description:cl.innerText.trim(),grade1:gv(t[6]),grade2:gv(t[7]),grade3:gv(t[8]),grade4:gv(t[9]),grade5:gv(t[10]),grade6:gv(t[11]),grade7:gv(t[12]),grade8:gv(t[13]),grade9:gv(t[14]),grade10:gv(t[15]),total:gv(t[16])};}).filter(Boolean);}""")
                        finally:
                            await ctx.close()
                    cards = await with_retry(f"psa set {set_name}", 3, scrape_set_cards)
                    for card in cards:
                        try:
                            await asyncio.to_thread(save_psa_card, set_name, set_link, card)
                            inserted_cards += 1
                            logged_cards += 1
                        except Exception as db_err:
                            failed.append({"input": year_url, "set": set_name, "card": f'{card.get("cardNumber","")} {card.get("cardName","")}'.strip(), "error": f"db save failed: {db_err}"})
                    summary["processedSets"] += 1
                except Exception as set_err:
                    summary["failedSets"] += 1
                    failed.append({"input": year_url, "set": set_name, "error": str(set_err)})
            result.append(summary)
        except (Exception, PWTimeoutError) as e:
            failed.append({"input": year_url, "error": str(e)})
    return {"ok": True, "count": len(result), "failedCount": len(failed), "insertedCards": inserted_cards, "loggedCards": logged_cards, "result": result, "failed": failed}

