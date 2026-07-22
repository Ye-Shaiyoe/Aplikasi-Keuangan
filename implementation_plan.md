# Implementation Plan: Adding Machine Learning Features and Python Setup

This plan details how to add a Machine Learning (ML) component to the finance application, including setting up Python, selecting features, and integrating it with the Go backend and React frontend.

## Proposed ML Features

We propose adding one or two of the following ML features:
1. **Transaction Categorization**: Automatically suggest/predict the transaction category based on the description entered by the user (using TF-IDF + Logistic Regression/Naive Bayes from `scikit-learn`).
2. **Expense Forecasting (Predictive Analytics)**: Forecast the user's spending trends for the next few months based on their transaction history (using `pandas` and `scikit-learn` or `statsmodels`).

---

## Architecture

We will introduce a lightweight Python service (using FastAPI) to serve the ML models. The Go backend will communicate with this service via REST API.

```
┌─────────────────┐       REST API      ┌─────────────────┐       REST API      ┌───────────────┐
│   React SPA     │ ──────────────────> │   Go Backend    │ ──────────────────> │ Python FastAPI│
│   (Frontend)    │ <────────────────── │  (Gin Server)   │ <────────────────── │ (ML Service)  │
└─────────────────┘                     └─────────────────┘                     └───────────────┘
```

---

## Python Setup Instructions

### 1. Python Version Recommendation
* **Recommended Version**: **Python 3.10.x** or **Python 3.11.x**.
* Avoid Python 3.12+ for now to prevent dependency compatibility issues with some ML/Data Science packages (like older tensorflow/scikit-learn versions).

### 2. Environment Setup
We will create a virtual environment (`.venv`) inside a new `ml` directory:
```bash
# Navigate to project root
cd finance-app

# Create ml directory
mkdir ml
cd ml

# Create virtual environment
python -m venv .venv

# Activate virtual environment
# On Windows (PowerShell):
.venv\Scripts\Activate.ps1
# On Windows (CMD):
.venv\Scripts\activate.bat
```

### 3. Required Packages (`requirements.txt`)
We will install the following packages:
* `fastapi` and `uvicorn` (for the API server)
* `pandas` and `numpy` (for data manipulation)
* `scikit-learn` (for ML models: classification and regression)
* `pydantic` (for request/response validation)

---

## Proposed Changes

### [ML Component] ✅ DONE

#### [NEW] [main.py](file:///c:/Users/akrom/Documents/finance-app/ml/main.py) ✅
* FastAPI server with 4 endpoints: `/api/ml/health`, `/api/ml/train`, `/api/ml/predict-category`, `/api/ml/forecast`.
* CORS configured for Go backend and React frontend.

#### [NEW] [model_trainer.py](file:///c:/Users/akrom/Documents/finance-app/ml/model_trainer.py) ✅
* `train_category_classifier`: TF-IDF + MultinomialNB pipeline, saves model + category map via joblib.
* `predict_category`: loads trained model, returns top-3 predictions with confidence.
* `forecast_expenses`: Linear Regression on monthly data, returns 3-month projections + R² confidence.

#### [NEW] [requirements.txt](file:///c:/Users/akrom/Documents/finance-app/ml/requirements.txt) ✅
* fastapi, uvicorn, pandas, numpy, scikit-learn, pydantic, joblib, python-dotenv — all installed in `.venv`.

### [Backend Component] ✅ DONE

#### Go ML Handler/Service/Repository/Model ✅
* `handler/ml.go`: MLPredictCategory, MLForecast, MLTrainModel, MLHealth.
* `service/ml.go`: HTTP client proxying to Python ML service (configurable via `ML_SERVICE_URL` env).
* `repository/ml.go`: GetAllTransactionsForML, GetAllCategoriesForML, GetHistoricalMetrics.
* `model/ml.go`: Request/response structs for all endpoints.
* Routes registered in `main.go` under protected `/api/ml/*`.

### [Frontend Component] ✅ DONE

#### [NEW] `src/pages/MLInsights.jsx` ✅
* Service health card showing model status.
* One-click training trigger with status feedback.
* Demo predict-category playground.
* Forecast section: next-month spotlight cards, 3-month area chart, projections table with R² confidence.

#### [MODIFY] `src/pages/Transactions.jsx` ✅
* Debounced ML auto-suggest on description input (triggers after 600ms, min 4 chars).
* Displays suggestion pill with category name + confidence %; one tap to apply.

#### [MODIFY] `src/api/client.js` ✅
* Added `mlPredictCategory`, `mlForecast`, `mlTrain`, `mlHealth`.

#### [MODIFY] `src/App.jsx` ✅
* Added `/insights/ml` route pointing to `MLInsights`.

#### [MODIFY] `src/components/Layout.jsx` ✅
* Added "ML Insights" (Brain icon) to Insight dropdown in both desktop sidebar and mobile nav.

---

## Verification Plan

### Manual Verification ✅ Build passes
1. **Start ML service**: `cd ml && .venv\Scripts\uvicorn main:app --reload --port 8000`
2. **Check Swagger UI**: `http://localhost:8000/docs`
3. **Start Go backend** (sets `ML_SERVICE_URL=http://localhost:8000` in `.env`): `cd backend && go run ./cmd/server`
4. **Start React dev server**: `cd frontend && npm run dev`
5. Go to **Insight → ML Insights** in the sidebar.
6. Click **Train Sekarang** — model trains from your transactions.
7. Go to **Transaksi → Tambah Transaksi**, type a description and watch the ML suggestion appear.
8. Forecast card shows 3-month projections once you have ≥ 3 months of data.
