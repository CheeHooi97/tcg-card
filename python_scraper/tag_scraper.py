import asyncio
from datetime import datetime, timezone

import psycopg

from common import db_config, new_page, unique_id, with_retry


def normalize_set_name(s: str) -> str:
    if s.startswith("Pokémon "):
        return s.replace("Pokémon ", "", 1).strip()
    return s


def split_number(card_number: str) -> tuple[str, str]:
    parts = card_number.split("/")
    if len(parts) == 2:
        return parts[0].strip(), parts[1].strip()
    return card_number.strip(), ""


def save_tag_card(set_name: str, set_link: str, card: dict) -> None:
    now = datetime.now(timezone.utc)
    card_number, set_number = split_number((card.get("cardNumber") or "").strip())
    card_name = (card.get("cardName") or "").strip()
    norm_set = normalize_set_name(set_name)
    with psycopg.connect(**db_config()) as conn:
        with conn.cursor() as cur:
            cur.execute("SELECT id,set_name,card_set,rarity,set_number FROM tag WHERE card_name=%s AND card_number=%s AND set_name=%s LIMIT 1", (card_name, card_number, norm_set))
            row = cur.fetchone()
            tag_id = row[0] if row else unique_id()
            if not row:
                cur.execute("""INSERT INTO tag (id,card_name,card_number,set_number,set_name,card_set,rarity,total,grade_va,grade1,grade1_5,grade2,grade2_5,grade3,grade3_5,grade4,grade4_5,grade5,grade5_5,grade6,grade6_5,grade7,grade7_5,grade8,grade8_5,grade9,grade10,grade10_p,created_date_time,updated_date_time)
                               VALUES (%s,%s,%s,%s,%s,%s,'',%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)""",
                            (tag_id, card_name, card_number, set_number, norm_set, set_name, str(card.get("total", "0")), str(card.get("gradeVA", "0")), str(card.get("grade1", "0")), str(card.get("grade1_5", "0")), str(card.get("grade2", "0")), str(card.get("grade2_5", "0")), str(card.get("grade3", "0")), str(card.get("grade3_5", "0")), str(card.get("grade4", "0")), str(card.get("grade4_5", "0")), str(card.get("grade5", "0")), str(card.get("grade5_5", "0")), str(card.get("grade6", "0")), str(card.get("grade6_5", "0")), str(card.get("grade7", "0")), str(card.get("grade7_5", "0")), str(card.get("grade8", "0")), str(card.get("grade8_5", "0")), str(card.get("grade9", "0")), str(card.get("grade10", "0")), str(card.get("grade10P", "0")), now, now))
            log_set_name = row[1] if row and row[1] else norm_set
            log_card_set = row[2] if row and row[2] else set_name
            log_rarity = row[3] if row and row[3] else ""
            log_set_number = row[4] if row and row[4] else set_number
            cur.execute("""INSERT INTO tag_logging (id,tag_id,card_name,card_number,set_number,set_name,card_set,rarity,description,total,grade_va,grade1,grade1_5,grade2,grade2_5,grade3,grade3_5,grade4,grade4_5,grade5,grade5_5,grade6,grade6_5,grade7,grade7_5,grade8,grade8_5,grade9,grade10,grade10_p,created_date_time,updated_date_time)
                           VALUES (%s,%s,%s,%s,%s,%s,%s,%s,'',%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)""",
                        (unique_id(), tag_id, card_name, card_number, log_set_number, log_set_name, log_card_set, log_rarity, str(card.get("total", "0")), str(card.get("gradeVA", "0")), str(card.get("grade1", "0")), str(card.get("grade1_5", "0")), str(card.get("grade2", "0")), str(card.get("grade2_5", "0")), str(card.get("grade3", "0")), str(card.get("grade3_5", "0")), str(card.get("grade4", "0")), str(card.get("grade4_5", "0")), str(card.get("grade5", "0")), str(card.get("grade5_5", "0")), str(card.get("grade6", "0")), str(card.get("grade6_5", "0")), str(card.get("grade7", "0")), str(card.get("grade7_5", "0")), str(card.get("grade8", "0")), str(card.get("grade8_5", "0")), str(card.get("grade9", "0")), str(card.get("grade10", "0")), str(card.get("grade10P", "0")), now, now))
            cur.execute("SELECT 1 FROM tag_url WHERE url=%s LIMIT 1", (set_link,))
            if cur.fetchone() is None:
                cur.execute("INSERT INTO tag_url (id,url,created_date_time,updated_date_time) VALUES (%s,%s,%s,%s)", (unique_id(), set_link, now, now))
        conn.commit()


async def scrape_tag(browser, urls: list[str]):
    result, failed = [], []
    saved, logged = 0, 0
    for year_url in urls:
        try:
            async def run_one():
                ctx, page = await new_page(browser)
                try:
                    await page.goto(year_url, timeout=60000)
                    await page.wait_for_load_state("domcontentloaded")
                    await page.wait_for_timeout(3000)
                    return await page.evaluate(
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
                finally:
                    await ctx.close()
            sets = await with_retry(f"tag year {year_url}", 3, run_one)
            processed_sets = 0
            for s in sets:
                set_link = s.get("link", "")
                set_name = s.get("name", "")
                if not set_link:
                    continue
                try:
                    async def scrape_cards():
                        ctx, page = await new_page(browser)
                        try:
                            await page.goto(set_link, timeout=60000)
                            await page.wait_for_selector("tr.MuiTableRow-root", timeout=45000)
                            await page.wait_for_timeout(3000)
                            return await page.evaluate(
                                """
                                () => {
                                  const rows = Array.from(document.querySelectorAll('tbody tr.MuiTableRow-root'));
                                  return rows.map(row => {
                                    const cells = Array.from(row.querySelectorAll('td'));
                                    if (cells.length < 22) return null;
                                    return {
                                      cardNumber: cells[0].innerText.trim(),
                                      cardName: cells[1].innerText.replace(/\\n/g, ' ').trim(),
                                      gradeVA: cells[cells.length - 21].innerText.trim() || "0",
                                      grade1: cells[cells.length - 20].innerText.trim() || "0",
                                      grade1_5: cells[cells.length - 19].innerText.trim() || "0",
                                      grade2: cells[cells.length - 18].innerText.trim() || "0",
                                      grade2_5: cells[cells.length - 17].innerText.trim() || "0",
                                      grade3: cells[cells.length - 16].innerText.trim() || "0",
                                      grade3_5: cells[cells.length - 15].innerText.trim() || "0",
                                      grade4: cells[cells.length - 14].innerText.trim() || "0",
                                      grade4_5: cells[cells.length - 13].innerText.trim() || "0",
                                      grade5: cells[cells.length - 12].innerText.trim() || "0",
                                      grade5_5: cells[cells.length - 11].innerText.trim() || "0",
                                      grade6: cells[cells.length - 10].innerText.trim() || "0",
                                      grade6_5: cells[cells.length - 9].innerText.trim() || "0",
                                      grade7: cells[cells.length - 8].innerText.trim() || "0",
                                      grade7_5: cells[cells.length - 7].innerText.trim() || "0",
                                      grade8: cells[cells.length - 6].innerText.trim() || "0",
                                      grade8_5: cells[cells.length - 5].innerText.trim() || "0",
                                      grade9: cells[cells.length - 4].innerText.trim() || "0",
                                      grade10: cells[cells.length - 3].innerText.trim() || "0",
                                      grade10P: cells[cells.length - 2].innerText.trim() || "0",
                                      total: cells[cells.length - 1].innerText.trim() || "0"
                                    };
                                  }).filter(Boolean);
                                }
                                """
                            )
                        finally:
                            await ctx.close()
                    cards = await with_retry(f"tag cards {set_name}", 3, scrape_cards)
                    for card in cards:
                        try:
                            await asyncio.to_thread(save_tag_card, set_name, set_link, card)
                            saved += 1
                            logged += 1
                        except Exception as db_err:
                            failed.append({"input": year_url, "set": set_name, "card": f'{card.get("cardNumber","")} {card.get("cardName","")}'.strip(), "error": f"db save failed: {db_err}"})
                    processed_sets += 1
                except Exception as set_err:
                    failed.append({"input": year_url, "set": set_name, "error": str(set_err)})
            result.append({"input": year_url, "sets": len(sets), "processedSets": processed_sets})
        except Exception as e:
            failed.append({"input": year_url, "error": str(e)})
    return {"ok": True, "count": len(result), "failedCount": len(failed), "insertedCards": saved, "loggedCards": logged, "result": result, "failed": failed}
