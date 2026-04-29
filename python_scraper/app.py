from fastapi import FastAPI
from playwright.async_api import async_playwright

from bgs_scraper import scrape_bgs
from cgc_scraper import scrape_cgc
from psa_scraper import scrape_psa
from schemas import BGSRequest, CGCRequest, PSARequest, TAGRequest
from tag_scraper import scrape_tag

app = FastAPI(title="TCG Python Scraper")


@app.get("/health")
async def health():
    return {"ok": True}


@app.post("/scrape/bgs")
async def bgs(req: BGSRequest):
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True, args=["--no-sandbox", "--disable-dev-shm-usage"])
        try:
            return await scrape_bgs(browser, req.url)
        finally:
            await browser.close()


@app.post("/scrape/cgc")
async def cgc(req: CGCRequest):
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True, args=["--no-sandbox", "--disable-dev-shm-usage"])
        try:
            return await scrape_cgc(browser, req.url)
        finally:
            await browser.close()


@app.post("/scrape/psa")
async def psa(req: PSARequest):
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True, args=["--no-sandbox", "--disable-dev-shm-usage"])
        try:
            return await scrape_psa(browser, req.urls)
        finally:
            await browser.close()


@app.post("/scrape/tag")
async def tag(req: TAGRequest):
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True, args=["--no-sandbox", "--disable-dev-shm-usage"])
        try:
            return await scrape_tag(browser, req.urls)
        finally:
            await browser.close()

