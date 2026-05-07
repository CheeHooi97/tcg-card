from pydantic import BaseModel, Field


class BGSRequest(BaseModel):
    url: list[str] = Field(default_factory=list)


class CGCRequest(BaseModel):
    url: str


class PSARequest(BaseModel):
    urls: list[str] = Field(default_factory=list)


class TAGRequest(BaseModel):
    urls: list[str] = Field(default_factory=list)


class PriceChartRequest(BaseModel):
    url: str
