import asyncio
from datetime import datetime, timezone

import psycopg

from common import db_config, new_page, unique_id, with_retry


async def scrape_sets(browser, year_url: str):
    async def run():
        ctx, page = await new_page(browser)
        try:
            await page.goto(year_url, timeout=60000)
            await page.wait_for_selector(".home-box.all ul", timeout=45000)
            return await page.evaluate(
                """
                () => {
                  const container = document.querySelector('.home-box.all');
                  if (!container) return [];
                  const baseUrl = "https://www.pricecharting.com";
                  return Array.from(container.querySelectorAll('ul li')).map(li => {
                    const a = li.querySelector('a');
                    if (!a) return null;
                    const href = a.getAttribute('href') || '';
                    return {
                      name: a.innerText.trim(),
                      link: href.startsWith('http') ? href : baseUrl + href
                    };
                  }).filter(Boolean);
                }
                """
            )
        finally:
            await ctx.close()
    return await with_retry(f"pc sets {year_url}", 3, run)


async def scrape_cards(browser, set_url: str):
    async def run():
        ctx, page = await new_page(browser)
        try:
            await page.goto(set_url, timeout=60000)
            await page.wait_for_selector("#games_table tbody tr[data-product]", timeout=45000)
            seen = set()
            all_cards = {}
            no_new = 0
            while True:
                batch = await page.evaluate(
                    """
                    () => {
                      const baseUrl = "https://www.pricecharting.com";
                      return Array.from(document.querySelectorAll('#games_table tbody tr[data-product]')).map(row => {
                        const a = row.querySelector('.title a');
                        const href = a ? a.getAttribute('href') : '';
                        return {
                          productId: row.getAttribute('data-product') || "",
                          name: a ? a.innerText.trim() : "",
                          price: row.querySelector('.used_price .js-price')?.innerText.trim() || "N/A",
                          link: a ? (href && href.startsWith('http') ? href : baseUrl + href) : ""
                        };
                      });
                    }
                    """
                )
                new_found = 0
                for b in batch:
                    pid = b.get("productId", "")
                    if pid and pid not in seen:
                        seen.add(pid)
                        all_cards[pid] = b
                        new_found += 1
                if new_found == 0:
                    no_new += 1
                else:
                    no_new = 0
                if no_new > 4:
                    break
                await page.evaluate("window.scrollBy(0, 1800);")
                await page.wait_for_timeout(2500)
            return list(all_cards.values())
        finally:
            await ctx.close()
    return await with_retry(f"pc cards {set_url}", 3, run)


async def scrape_card_detail(browser, detail_url: str):
    async def run():
        ctx, page = await new_page(browser)
        try:
            await page.goto(detail_url, timeout=60000)
            await page.wait_for_selector("#product_details .cover img", timeout=45000)
            await page.wait_for_selector("#price_data", timeout=45000)
            return await page.evaluate(
                """
                () => {
                  const getPrice = (id) => {
                    const el = document.querySelector("#" + id + " .price");
                    return el ? el.innerText.trim().replace("\\n","") : "N/A";
                  };
                  const setAnchor = document.querySelector('#product_name a');
                  let extractedSetName = "";
                  if (setAnchor) {
                    for (const node of setAnchor.childNodes) {
                      if (node.nodeType === Node.TEXT_NODE) extractedSetName += node.textContent;
                    }
                  }
                  const imgEl = document.querySelector('#product_details .cover img');
                  return {
                    setName: extractedSetName.trim(),
                    imageUrl: imgEl ? imgEl.src : "",
                    ungraded: getPrice("used_price"),
                    grade7: getPrice("complete_price"),
                    grade8: getPrice("new_price"),
                    grade9: getPrice("graded_price"),
                    grade9_5: getPrice("box_only_price"),
                    grade10: getPrice("manual_only_price")
                  };
                }
                """
            )
        finally:
            await ctx.close()
    return await with_retry(f"pc detail {detail_url}", 3, run)


def save_card_and_price(card: dict, detail: dict):
    now = datetime.now(timezone.utc)
    with psycopg.connect(**db_config()) as conn:
        with conn.cursor() as cur:
            cur.execute("SELECT id FROM card WHERE name=%s AND set_name=%s LIMIT 1", ((card.get("name") or "").strip(), (detail.get("setName") or "").strip()))
            row = cur.fetchone()
            if row:
                return False
            card_id = unique_id()
            cur.execute(
                """INSERT INTO card (id,name,number,set_number,set_name,rarity,card_type,ungrade,grade7,grade8,grade9,grade9_5,grade10,photo_url,created_date_time,updated_date_time)
                   VALUES (%s,%s,'','',%s,'','Pokemon',%s,%s,%s,%s,%s,%s,%s,%s,%s)""",
                (card_id, (card.get("name") or "").strip(), (detail.get("setName") or "").strip(), str(card.get("price") or "N/A"), str(detail.get("grade7") or "N/A"), str(detail.get("grade8") or "N/A"), str(detail.get("grade9") or "N/A"), str(detail.get("grade9_5") or "N/A"), str(detail.get("grade10") or "N/A"), (detail.get("imageUrl") or "").strip(), now, now),
            )
            cur.execute(
                """INSERT INTO card_price (id,card_id,name,set,ungrade,grade7,grade8,grade9,grade9_5,grade10,created_date_time,updated_date_time)
                   VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)""",
                (unique_id(), card_id, (card.get("name") or "").strip(), (detail.get("setName") or "").strip(), str(detail.get("ungraded") or "N/A"), str(detail.get("grade7") or "N/A"), str(detail.get("grade8") or "N/A"), str(detail.get("grade9") or "N/A"), str(detail.get("grade9_5") or "N/A"), str(detail.get("grade10") or "N/A"), now, now),
            )
        conn.commit()
    return True


async def scrape_pricechart(browser, year_url: str):
    failed = []
    sets = await scrape_sets(browser, year_url)
    inserted = 0
    for y, s in enumerate(sets):
        set_name = s.get("name", "")
        print(f"[pricechart] currently processing set index={y} name={set_name}")
        if y < 152:
            continue
        set_link = s.get("link", "")
        if not set_link:
            continue
        try:
            cards = await scrape_cards(browser, set_link)
        except Exception as e:
            failed.append({"set": s.get("name", ""), "error": str(e)})
            continue
        for z, card in enumerate(cards):
            try:
                detail = await scrape_card_detail(browser, card.get("link", ""))
                created = await asyncio.to_thread(save_card_and_price, card, detail)
                if created:
                    inserted += 1
            except Exception as e:
                failed.append({"set": s.get("name", ""), "card": card.get("name", ""), "error": str(e)})
                continue
    return {"ok": True, "sets": len(sets), "insertedCards": inserted, "failedCount": len(failed), "failed": failed}
