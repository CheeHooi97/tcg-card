import asyncio
from datetime import datetime, timezone
from typing import Any

import psycopg

from common import db_config, new_page, unique_id, with_retry


def save_bgs_card(card: dict[str, Any], source_url: str) -> None:
    now = datetime.now(timezone.utc)
    p = {
        "card_name": (card.get("cardName") or "").strip(),
        "card_number": (card.get("cardNumber") or "").strip(),
        "set_id": (card.get("setID") or "").strip(),
        "set_name": (card.get("setTitle") or "").strip(),
        "total": str(card.get("totalCount") or "0"),
        "grades": card.get("gradeCounts") or {},
    }
    g = p["grades"]
    with psycopg.connect(**db_config()) as conn:
        with conn.cursor() as cur:
            cur.execute("SELECT id,set_name,set_number,rarity,description FROM bgs WHERE card_name=%s AND card_number=%s AND set_id=%s LIMIT 1", (p["card_name"], p["card_number"], p["set_id"]))
            row = cur.fetchone()
            bgs_id = row[0] if row else unique_id()
            if not row:
                cur.execute(
                    """INSERT INTO bgs (id,card_name,card_number,set_number,set_name,rarity,description,set_id,total,grade1,grade1_5,grade2,grade2_5,grade3,grade3_5,grade4,grade4_5,grade5,grade5_5,grade6,grade6_5,grade7,grade7_5,grade8,grade8_5,grade9,grade9_5,grade10_p,grade10_bl,created_date_time,updated_date_time)
                       VALUES (%s,%s,%s,'',%s,'','',%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)""",
                    (bgs_id, p["card_name"], p["card_number"], p["set_name"], p["set_id"], p["total"], g.get("1", "0"), g.get("1.5", "0"), g.get("2", "0"), g.get("2.5", "0"), g.get("3", "0"), g.get("3.5", "0"), g.get("4", "0"), g.get("4.5", "0"), g.get("5", "0"), g.get("5.5", "0"), g.get("6", "0"), g.get("6.5", "0"), g.get("7", "0"), g.get("7.5", "0"), g.get("8", "0"), g.get("8.5", "0"), g.get("9", "0"), g.get("9.5", "0"), g.get("10P", "0"), g.get("10BL", "0"), now, now),
                )
            log_set_name = row[1] if row and row[1] else p["set_name"]
            log_set_number = row[2] if row and row[2] else ""
            log_rarity = row[3] if row and row[3] else ""
            log_desc = row[4] if row and row[4] else ""
            cur.execute(
                """INSERT INTO bgs_logging (id,bgs_id,card_name,card_number,set_number,set_name,rarity,description,set_id,total,grade1,grade1_5,grade2,grade2_5,grade3,grade3_5,grade4,grade4_5,grade5,grade5_5,grade6,grade6_5,grade7,grade7_5,grade8,grade8_5,grade9,grade9_5,grade10_p,grade10_bl,created_date_time,updated_date_time)
                   VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)""",
                (unique_id(), bgs_id, p["card_name"], p["card_number"], log_set_number, log_set_name, log_rarity, log_desc, p["set_id"], p["total"], g.get("1", "0"), g.get("1.5", "0"), g.get("2", "0"), g.get("2.5", "0"), g.get("3", "0"), g.get("3.5", "0"), g.get("4", "0"), g.get("4.5", "0"), g.get("5", "0"), g.get("5.5", "0"), g.get("6", "0"), g.get("6.5", "0"), g.get("7", "0"), g.get("7.5", "0"), g.get("8", "0"), g.get("8.5", "0"), g.get("9", "0"), g.get("9.5", "0"), g.get("10P", "0"), g.get("10BL", "0"), now, now),
            )
            cur.execute("SELECT 1 FROM bgs_url WHERE url=%s LIMIT 1", (source_url,))
            if cur.fetchone() is None:
                cur.execute("INSERT INTO bgs_url (id,url,created_date_time,updated_date_time) VALUES (%s,%s,%s,%s)", (unique_id(), source_url, now, now))
        conn.commit()


async def scrape_bgs(browser, set_names: list[str]):
    result, failed = [], []
    saved, logged = 0, 0
    for set_name in set_names:
        try:
            async def run_one():
                ctx, page = await new_page(browser)
                try:
                    await page.goto("https://www.beckett.com/grading/pop-report", timeout=60000)
                    await page.wait_for_selector("#set_name", timeout=30000)
                    await page.fill("#set_name", set_name)
                    await page.keyboard.press("Enter")
                    await page.wait_for_timeout(2000)
                    return await page.eval_on_selector_all('a[href*="/set_match/"]', "els => els.map(a => ({url: a.href, setName: (a.innerText||'').trim()}))")
                finally:
                    await ctx.close()
            links = await with_retry(f"bgs set {set_name}", 3, run_one)
            processed = 0
            for link in links:
                u = link.get("url", "")
                if not u:
                    continue
                try:
                    async def scrape_link():
                        ctx, page = await new_page(browser)
                        try:
                            await page.goto(u, timeout=60000)
                            await page.wait_for_selector("tr.rows", timeout=45000)
                            return await page.evaluate(
                                """
                                () => {
                                  const rows = Array.from(document.querySelectorAll("tr.rows"));
                                  return rows.map(r => {
                                    const setTitle = (r.querySelector("input.set_title")||{}).value || "";
                                    const totalVal = (r.querySelector("input.card_total_value")||{}).value || "0";
                                    const tds = Array.from(r.querySelectorAll("td.test"));
                                    const cardName = (tds[0]?.innerText || "").trim();
                                    const cardNumber = (tds[1]?.innerText || "").trim();
                                    const a = r.querySelector("a");
                                    let setID = "";
                                    if (a && a.href) {
                                      const uu = new URL(a.href, location.origin);
                                      setID = uu.searchParams.get("set_id") || "";
                                    }
                                    const gradeCounts = {};
                                    Array.from(r.querySelectorAll("td")).forEach(td => {
                                      const gv = (td.querySelector("input.header_grade")||{}).value || "";
                                      if (!gv) return;
                                      let count = ((td.querySelector("b.popCard a")||{}).innerText || td.innerText || "").trim();
                                      if (!count || count === "-") count = "0";
                                      if (gv === "10") {
                                        const ht = (td.querySelector("input.header_type")||{}).value || "";
                                        if (ht.includes("Black Label")) gradeCounts["10BL"] = count; else gradeCounts["10P"] = count;
                                      } else {
                                        gradeCounts[gv] = count;
                                      }
                                    });
                                    return {cardName, cardNumber, totalCount: totalVal, setTitle, setID, gradeCounts};
                                  });
                                }
                                """
                            )
                        finally:
                            await ctx.close()
                    cards = await with_retry(f"bgs cards {u}", 3, scrape_link)
                    for card in cards:
                        try:
                            await asyncio.to_thread(save_bgs_card, card, u)
                            saved += 1
                            logged += 1
                        except Exception as db_err:
                            failed.append({"input": set_name, "link": u, "card": f'{card.get("cardNumber","")} {card.get("cardName","")}'.strip(), "error": f"db save failed: {db_err}"})
                    processed += 1
                except Exception as link_err:
                    failed.append({"input": set_name, "link": u, "error": str(link_err)})
            result.append({"input": set_name, "links": len(links), "processedLinks": processed})
        except Exception as e:
            failed.append({"input": set_name, "error": str(e)})
    return {"ok": True, "count": len(result), "failedCount": len(failed), "insertedCards": saved, "loggedCards": logged, "result": result, "failed": failed}
