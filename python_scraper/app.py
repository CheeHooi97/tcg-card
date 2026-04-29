import asyncio
from typing import Any

from fastapi import FastAPI
from pydantic import BaseModel, Field
from playwright.async_api import async_playwright, TimeoutError as PWTimeoutError

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

    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True, args=["--no-sandbox", "--disable-dev-shm-usage"])

        for year_url in req.urls:
            try:
                async def run_one():
                    ctx, page = await new_page(browser)
                    try:
                        await page.goto(year_url, timeout=60000)
                        await page.wait_for_selector("#tableSets", timeout=45000)
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

                sets = await with_retry(f"psa year {year_url}", 3, run_one)
                result.append({"input": year_url, "sets": sets})
            except (Exception, PWTimeoutError) as e:
                failed.append({"input": year_url, "error": str(e)})
                continue

        await browser.close()

    return {"ok": True, "count": len(result), "failedCount": len(failed), "result": result, "failed": failed}


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

