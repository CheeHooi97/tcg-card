# Python Scraper Service

## 1) Setup

```bash
cd ~/tcg-card/python_scraper
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
playwright install chromium
```

## 2) Run

```bash
uvicorn app:app --host 0.0.0.0 --port 8010
```

## 3) Endpoints

- `GET /health`
- `POST /scrape/bgs` body: `{"url":["Pokemon"]}`
- `POST /scrape/cgc` body: `{"url":"https://www.cgccards.com/population-report/..." }`
- `POST /scrape/psa` body: `{"urls":["https://www.psacard.com/pop/..."]}`
- `POST /scrape/tag` body: `{"urls":["https://my.taggrading.com/pop-report/..."]}`

All scrape endpoints retry failed items and continue until all inputs are processed.
