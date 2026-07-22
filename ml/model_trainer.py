"""
model_trainer.py — ML Model Trainer for Finance App

Trains two models:
1. Category Classifier: Predicts transaction category from description text.
   Uses TF-IDF vectorization + Multinomial Naive Bayes.
2. Expense Forecaster: Predicts next month's expense using Linear Regression
   on monthly aggregate data.

Models are saved as joblib files in the ./models/ directory.
"""

import os
import json
import joblib
import numpy as np
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.naive_bayes import MultinomialNB
from sklearn.linear_model import LinearRegression
from sklearn.pipeline import Pipeline


MODELS_DIR = os.path.join(os.path.dirname(__file__), "models")


def ensure_models_dir():
    """Create the models directory if it doesn't exist."""
    os.makedirs(MODELS_DIR, exist_ok=True)


def train_category_classifier(transactions: list[dict], categories: list[dict]) -> Pipeline:
    """
    Train a text classifier that predicts transaction category_id from description.

    Args:
        transactions: list of dicts with keys: description, category_id, type
        categories: list of dicts with keys: id, name, type
    Returns:
        A trained sklearn Pipeline (TF-IDF + MultinomialNB)
    """
    # Filter out transactions with empty descriptions
    valid = [t for t in transactions if t.get("description", "").strip()]

    if len(valid) < 5:
        raise ValueError("Minimal 5 transaksi dengan deskripsi diperlukan untuk training model kategorisasi.")

    descriptions = [t["description"] for t in valid]
    labels = [t["category_id"] for t in valid]

    # Build pipeline: TF-IDF vectorizer → Naive Bayes classifier
    pipeline = Pipeline([
        ("tfidf", TfidfVectorizer(
            max_features=5000,
            ngram_range=(1, 2),
            stop_words=None,  # Keep Indonesian words
            min_df=1,
            sublinear_tf=True,
        )),
        ("clf", MultinomialNB(alpha=0.1)),
    ])

    pipeline.fit(descriptions, labels)

    # Save model
    ensure_models_dir()
    model_path = os.path.join(MODELS_DIR, "category_classifier.joblib")
    joblib.dump(pipeline, model_path)

    # Save category mapping for later use
    cat_map = {c["id"]: c["name"] for c in categories}
    cat_map_path = os.path.join(MODELS_DIR, "category_map.json")
    with open(cat_map_path, "w", encoding="utf-8") as f:
        json.dump(cat_map, f, ensure_ascii=False)

    return pipeline


def train_expense_forecaster(monthly_data: list[dict]) -> dict:
    """
    Train a simple linear regression to forecast next month's income and expense.

    Args:
        monthly_data: list of dicts with keys: month, year, income, expense
                      (sorted chronologically)
    Returns:
        dict with forecasted income, expense, and model metadata
    """
    if len(monthly_data) < 3:
        raise ValueError("Minimal 3 bulan data historis diperlukan untuk forecasting.")

    # Create feature matrix: sequential index (1, 2, 3, ...)
    X = np.arange(1, len(monthly_data) + 1).reshape(-1, 1)
    y_expense = np.array([m["expense"] for m in monthly_data])
    y_income = np.array([m["income"] for m in monthly_data])

    # Train Linear Regression for expense
    expense_model = LinearRegression()
    expense_model.fit(X, y_expense)

    # Train Linear Regression for income
    income_model = LinearRegression()
    income_model.fit(X, y_income)

    # Forecast next month
    next_idx = np.array([[len(monthly_data) + 1]])
    forecast_expense = max(0, int(expense_model.predict(next_idx)[0]))
    forecast_income = max(0, int(income_model.predict(next_idx)[0]))

    # Save models
    ensure_models_dir()
    joblib.dump(expense_model, os.path.join(MODELS_DIR, "expense_forecaster.joblib"))
    joblib.dump(income_model, os.path.join(MODELS_DIR, "income_forecaster.joblib"))

    # Calculate R² scores for confidence indication
    r2_expense = expense_model.score(X, y_expense)
    r2_income = income_model.score(X, y_income)

    return {
        "forecast_expense": forecast_expense,
        "forecast_income": forecast_income,
        "confidence": {
            "expense_r2": round(r2_expense, 4),
            "income_r2": round(r2_income, 4),
        },
        "data_points": len(monthly_data),
    }


def predict_category(description: str) -> dict:
    """
    Predict category for a transaction description using the trained model.

    Returns:
        dict with predicted category_id, category_name, and confidence score
    """
    model_path = os.path.join(MODELS_DIR, "category_classifier.joblib")
    cat_map_path = os.path.join(MODELS_DIR, "category_map.json")

    if not os.path.exists(model_path):
        raise FileNotFoundError("Model kategorisasi belum di-train. Silakan train terlebih dahulu.")

    pipeline = joblib.load(model_path)

    with open(cat_map_path, "r", encoding="utf-8") as f:
        cat_map = json.load(f)

    # Predict
    predicted_id = pipeline.predict([description])[0]
    probabilities = pipeline.predict_proba([description])[0]
    confidence = float(max(probabilities))

    # Get top 3 predictions
    classes = pipeline.classes_
    top_indices = np.argsort(probabilities)[::-1][:3]
    top_predictions = []
    for idx in top_indices:
        cid = int(classes[idx])
        top_predictions.append({
            "category_id": cid,
            "category_name": cat_map.get(str(cid), f"Kategori {cid}"),
            "confidence": round(float(probabilities[idx]), 4),
        })

    return {
        "predicted_category_id": int(predicted_id),
        "predicted_category_name": cat_map.get(str(predicted_id), f"Kategori {predicted_id}"),
        "confidence": round(confidence, 4),
        "top_predictions": top_predictions,
    }


def forecast_expenses(monthly_data: list[dict]) -> dict:
    """
    Forecast using saved models or retrain on the fly.

    Args:
        monthly_data: list of dicts with keys: month, year, income, expense

    Returns:
        dict with forecasted income, expense, confidence, and multi-month projections
    """
    if len(monthly_data) < 3:
        raise ValueError("Minimal 3 bulan data historis diperlukan untuk forecasting.")

    X = np.arange(1, len(monthly_data) + 1).reshape(-1, 1)
    y_expense = np.array([m["expense"] for m in monthly_data])
    y_income = np.array([m["income"] for m in monthly_data])

    # Train on the fly with current data (always fresh)
    expense_model = LinearRegression()
    expense_model.fit(X, y_expense)

    income_model = LinearRegression()
    income_model.fit(X, y_income)

    # Forecast next 3 months
    projections = []
    for i in range(1, 4):
        next_idx = np.array([[len(monthly_data) + i]])
        proj_expense = max(0, int(expense_model.predict(next_idx)[0]))
        proj_income = max(0, int(income_model.predict(next_idx)[0]))
        projections.append({
            "month_offset": i,
            "forecast_expense": proj_expense,
            "forecast_income": proj_income,
            "forecast_balance": proj_income - proj_expense,
        })

    r2_expense = expense_model.score(X, y_expense)
    r2_income = income_model.score(X, y_income)

    # Determine confidence level text
    avg_r2 = (r2_expense + r2_income) / 2
    if avg_r2 >= 0.7:
        confidence_level = "Tinggi"
    elif avg_r2 >= 0.4:
        confidence_level = "Sedang"
    else:
        confidence_level = "Rendah"

    return {
        "forecast_expense": projections[0]["forecast_expense"],
        "forecast_income": projections[0]["forecast_income"],
        "forecast_balance": projections[0]["forecast_balance"],
        "confidence": {
            "expense_r2": round(r2_expense, 4),
            "income_r2": round(r2_income, 4),
            "level": confidence_level,
        },
        "projections": projections,
        "data_points": len(monthly_data),
        "model_type": "Linear Regression",
    }
