import asyncio
from datetime import datetime, timezone

import psycopg

from common import db_config, new_page, unique_id, with_retry


def save_cgc_card(card: dict, set_url: str) -> None:
    now = datetime.now(timezone.utc)
    card_number_raw = (card.get("cardNumber") or "").strip()
    parts = card_number_raw.split("/")
    card_number = parts[0].strip() if parts else card_number_raw
    set_number = parts[1].strip() if len(parts) == 2 else ""
    card_name = (card.get("cardName") or "").strip()
    rarity = ""
    if ")" in card_name:
        idx = card_name.find(")")
        rarity = card_name[idx + 1:].strip()
        card_name = card_name[:idx + 1].strip()
    total = str(card.get("totalGraded") or "0")
    grades = {str((g.get("grade") or "")).strip(): str((g.get("count") or "0")).strip() for g in (card.get("grades") or [])}
    with psycopg.connect(**db_config()) as conn:
        with conn.cursor() as cur:
            set_name = (card.get("setName") or "").strip()
            cur.execute("SELECT id FROM cgc WHERE card_name=%s AND card_number=%s AND set_name=%s AND rarity=%s LIMIT 1", (card_name, card_number, set_name, rarity))
            row = cur.fetchone()
            cgc_id = row[0] if row else unique_id()
            if not row:
                cur.execute("""INSERT INTO cgc (id,card_name,card_number,set_number,set_name,rarity,total,grade1,grade1_5,grade2,grade2_5,grade3,grade3_5,grade4,grade4_5,grade5,grade5_5,grade6,grade6_5,grade7,grade7_5,grade8,grade8_5,grade9,grade9_5,grade10,grade10_p,created_date_time,updated_date_time)
                               VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)""",
                            (cgc_id, card_name, card_number, set_number, set_name, rarity, total, grades.get("1", "0"), grades.get("1.5", "0"), grades.get("2", "0"), grades.get("2.5", "0"), grades.get("3", "0"), grades.get("3.5", "0"), grades.get("4", "0"), grades.get("4.5", "0"), grades.get("5", "0"), grades.get("5.5", "0"), grades.get("6", "0"), grades.get("6.5", "0"), grades.get("7", "0"), grades.get("7.5", "0"), grades.get("8", "0"), grades.get("8.5", "0"), grades.get("9", "0"), grades.get("Mint+ 9.5", "0"), grades.get("Gem Mint 10", "0"), grades.get("Pristine 10", "0"), now, now))
            cur.execute("""INSERT INTO cgc_logging (id,cgc_id,card_name,card_number,set_number,set_name,rarity,total,grade1,grade1_5,grade2,grade2_5,grade3,grade3_5,grade4,grade4_5,grade5,grade5_5,grade6,grade6_5,grade7,grade7_5,grade8,grade8_5,grade9,grade9_5,grade10,grade10_p,grade10_bl,created_date_time,updated_date_time)
                           VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,'',%s,%s)""",
                        (unique_id(), cgc_id, card_name, card_number, set_number, set_name, rarity, total, grades.get("1", "0"), grades.get("1.5", "0"), grades.get("2", "0"), grades.get("2.5", "0"), grades.get("3", "0"), grades.get("3.5", "0"), grades.get("4", "0"), grades.get("4.5", "0"), grades.get("5", "0"), grades.get("5.5", "0"), grades.get("6", "0"), grades.get("6.5", "0"), grades.get("7", "0"), grades.get("7.5", "0"), grades.get("8", "0"), grades.get("8.5", "0"), grades.get("9", "0"), grades.get("Mint+ 9.5", "0"), grades.get("Gem Mint 10", "0"), grades.get("Pristine 10", "0"), now, now))
            cur.execute("SELECT 1 FROM cgc_url WHERE url=%s LIMIT 1", (set_url,))
            if cur.fetchone() is None:
                cur.execute("INSERT INTO cgc_url (id,url,created_date_time,updated_date_time) VALUES (%s,%s,%s,%s)", (unique_id(), set_url, now, now))
        conn.commit()


async def scrape_cgc(browser, url: str):
    result, failed = [], []
    saved, logged = 0, 0
    try:
        async def run_root():
            ctx, page = await new_page(browser)
            try:
                await page.goto(url, timeout=60000)
                await page.wait_for_selector(".ccg-cards", timeout=45000)
                return await page.evaluate(
                    """
                    () => {
                      const baseUrl = "https://www.cgccards.com";
                      return Array.from(document.querySelectorAll(".card.ng-scope a"))
                        .map(a => ({url: baseUrl + a.getAttribute("href")}))
                        .filter(x => x.url);
                    }
                    """
                )
            finally:
                await ctx.close()
        lists = await with_retry("cgc list", 3, run_root)
        for li in lists:
            list_url = li.get("url", "")
            if not list_url:
                continue
            try:
                async def scrape_sets():
                    ctx, page = await new_page(browser)
                    try:
                        await page.goto(list_url, timeout=60000)
                        await page.wait_for_selector("tr.ccg-setcounts-table__row", timeout=45000)
                        return await page.evaluate("""() => { const baseUrl="https://www.cgccards.com"; const rows=Array.from(document.querySelectorAll('tr.ccg-setcounts-table__row')); return rows.map(r=>{const a=r.querySelector('.ccg-setcounts-table__name a'); return a?{setUrl:baseUrl+a.getAttribute('href'), setName:a.innerText.trim()}:null;}).filter(Boolean);}""")
                    finally:
                        await ctx.close()
                sets = await with_retry(f"cgc sets {list_url}", 3, scrape_sets)
                processed_sets = 0
                for st in sets:
                    set_url = st.get("setUrl", "")
                    if not set_url:
                        continue
                    try:
                        async def scrape_cards():
                            ctx, page = await new_page(browser)
                            try:
                                await page.goto(set_url, timeout=60000)
                                await page.wait_for_selector("tr.needs-alignment", timeout=45000)
                                return await page.evaluate(
                                    """
                                    () => {
                                      const headerTitle = document.querySelector('.card-list-header__title');
                                      const setName = headerTitle ? headerTitle.innerText.trim() : "";
                                      const headerCells = Array.from(document.querySelectorAll('#tableScroller thead th.ng-binding'));
                                      const gradeLabels = headerCells.map(h => h.innerText.trim()).filter(t => t !== "");
                                      const pinnedRows = Array.from(document.querySelectorAll('.pinned tbody tr.needs-alignment'));
                                      const dataRows = Array.from(document.querySelectorAll('#tableScroller tbody tr.needs-alignment.ng-scope'));
                                      return pinnedRows.map((pRow, index) => {
                                        const numCell = pRow.querySelector('.card-list__cardNumber');
                                        const nameCell = pRow.querySelector('.card-list__name');
                                        const totalCell = pRow.querySelector('.card-list__totalGraded');
                                        const dRow = dataRows[index];
                                        if (!dRow) return null;
                                        const allTds = Array.from(dRow.querySelectorAll('td'));
                                        const grades = gradeLabels.map((label, i) => {
                                          const cell = allTds[1 + i];
                                          const span = cell ? cell.querySelector('span.ng-binding') : null;
                                          const val = span ? span.innerText.trim() : "";
                                          return { grade: label, count: val === "" ? "0" : val };
                                        });
                                        return { setName, cardNumber: numCell ? numCell.innerText.split('\\n')[0].trim() : "", cardName: nameCell ? nameCell.innerText.replace(/\\s+/g, ' ').trim() : "", totalGraded: totalCell ? totalCell.innerText.trim() : "0", grades };
                                      }).filter(Boolean);
                                    }
                                    """
                                )
                            finally:
                                await ctx.close()
                        cards = await with_retry(f"cgc cards {set_url}", 3, scrape_cards)
                        for card in cards:
                            try:
                                await asyncio.to_thread(save_cgc_card, card, set_url)
                                saved += 1
                                logged += 1
                            except Exception as db_err:
                                failed.append({"list": list_url, "set": st.get("setName", ""), "card": f'{card.get("cardNumber","")} {card.get("cardName","")}'.strip(), "error": f"db save failed: {db_err}"})
                        processed_sets += 1
                    except Exception as set_err:
                        failed.append({"list": list_url, "set": st.get("setName", ""), "error": str(set_err)})
                result.append({"list": list_url, "sets": len(sets), "processedSets": processed_sets})
            except Exception as li_err:
                failed.append({"list": list_url, "error": str(li_err)})
    except Exception as e:
        failed.append({"input": url, "error": str(e)})
    return {"ok": True, "count": len(result), "failedCount": len(failed), "insertedCards": saved, "loggedCards": logged, "result": result, "failed": failed}
