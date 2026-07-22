"""
main.py — FastAPI ML Service for Finance App

Provides REST API endpoints for:
1. POST /api/ml/predict-category — Predict category from transaction description
2. POST /api/ml/forecast — Forecast future income/expense
3. POST /api/ml/train — Train models with provided data
4. GET  /api/ml/health — Health check
"""

import os
import traceback
from typing import Optional

from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field

from model_trainer import (
    train_category_classifier,
    predict_category,
    forecast_expenses,
)

app = FastAPI(
    title="Finance ML Service",
    description="Machine Learning microservice for Catatan Keuangan finance app",
    version="1.0.0",
)

# CORS — allow Go backend and React frontend
app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://localhost:8080",
        "http://localhost:5173",
        "https://umkmkeuangan.vercel.app",
    ],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


# ──────────────── Request/Response Models ────────────────


class PredictCategoryRequest(BaseModel):
    description: str = Field(..., min_length=1, description="Deskripsi transaksi")


class PredictCategoryResponse(BaseModel):
    predicted_category_id: int
    predicted_category_name: str
    confidence: float
    top_predictions: list[dict]


class TransactionData(BaseModel):
    description: str
    category_id: int
    type: str  # "income" or "expense"


class CategoryData(BaseModel):
    id: int
    name: str
    type: str


class TrainRequest(BaseModel):
    transactions: list[TransactionData]
    categories: list[CategoryData]


class MonthlyData(BaseModel):
    month: int
    year: int
    income: int
    expense: int


class ForecastRequest(BaseModel):
    monthly_data: list[MonthlyData]


class ForecastResponse(BaseModel):
    forecast_expense: int
    forecast_income: int
    forecast_balance: int
    confidence: dict
    projections: list[dict]
    data_points: int
    model_type: str


# ──────────────── Endpoints ────────────────


@app.get("/api/ml/health")
async def health_check():
    """Health check endpoint."""
    models_dir = os.path.join(os.path.dirname(__file__), "models")
    has_category_model = os.path.exists(os.path.join(models_dir, "category_classifier.joblib"))
    return {
        "status": "healthy",
        "service": "finance-ml",
        "models": {
            "category_classifier": "ready" if has_category_model else "not_trained",
        },
    }


@app.post("/api/ml/train")
async def train_models(req: TrainRequest):
    """
    Train the category classification model with user's transaction data.
    Called by the Go backend when user initiates training.
    """
    try:
        transactions = [t.model_dump() for t in req.transactions]
        categories = [c.model_dump() for c in req.categories]

        pipeline = train_category_classifier(transactions, categories)
        n_classes = len(pipeline.classes_)

        return {
            "status": "success",
            "message": f"Model kategorisasi berhasil di-train dengan {len(transactions)} transaksi dan {n_classes} kategori.",
            "details": {
                "total_transactions": len(transactions),
                "total_categories": n_classes,
            },
        }
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        traceback.print_exc()
        raise HTTPException(status_code=500, detail=f"Training gagal: {str(e)}")


@app.post("/api/ml/predict-category", response_model=PredictCategoryResponse)
async def predict_category_endpoint(req: PredictCategoryRequest):
    """
    Predict the most likely category for a transaction description.
    Returns top 3 predictions with confidence scores.
    """
    try:
        result = predict_category(req.description)
        return result
    except FileNotFoundError as e:
        raise HTTPException(status_code=404, detail=str(e))
    except Exception as e:
        traceback.print_exc()
        raise HTTPException(status_code=500, detail=f"Prediksi gagal: {str(e)}")


@app.post("/api/ml/forecast", response_model=ForecastResponse)
async def forecast_endpoint(req: ForecastRequest):
    """
    Forecast next months' income and expense based on historical data.
    Returns 3-month projections with confidence metrics.
    If data < 3 months, returns a zero-filled forecast with a 'Rendah' confidence
    instead of raising an error.
    """
    try:
        monthly_data = [m.model_dump() for m in req.monthly_data]

        # Not enough data — return empty forecast instead of 500
        if len(monthly_data) < 3:
            zero_proj = [
                {"month_offset": i, "forecast_expense": 0, "forecast_income": 0, "forecast_balance": 0}
                for i in range(1, 4)
            ]
            return {
                "forecast_expense": 0,
                "forecast_income": 0,
                "forecast_balance": 0,
                "confidence": {"expense_r2": 0.0, "income_r2": 0.0, "level": "Rendah"},
                "projections": zero_proj,
                "data_points": len(monthly_data),
                "model_type": "Tidak cukup data (min. 3 bulan)",
            }

        result = forecast_expenses(monthly_data)
        return result
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        traceback.print_exc()
        raise HTTPException(status_code=500, detail=f"Forecasting gagal: {str(e)}")


if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("ML_PORT", "8000"))
    uvicorn.run("main:app", host="0.0.0.0", port=port, reload=True)
