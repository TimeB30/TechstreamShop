from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import keygen
import uvicorn
from ahk import AHK
import asyncio

app = FastAPI()

activation_lock = asyncio.Lock()

class ActivationRequest(BaseModel):
    software_id: str
    version: int
    days: int

class ActivationResponse(BaseModel):
    key: str

@app.post("/key", response_model=ActivationResponse)
async def process_key(request: ActivationRequest):
    try:
        async with activation_lock:
            loop = asyncio.get_event_loop()
            result = await loop.run_in_executor(
                None,
                lambda: keygen.get_activator("AHK()", request.software_id, request.version, request.days)
            )
        return ActivationResponse(key=result)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)