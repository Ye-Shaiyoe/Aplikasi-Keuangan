import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Attach JWT token to every request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle 401 responses (token expired/invalid)
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth
export const authRegister = (data) =>
  api.post('/auth/register', data).then((r) => r.data);

export const authLogin = (data) =>
  api.post('/auth/login', data).then((r) => r.data);

// Categories
export const getCategories = (type) =>
  api.get('/categories', { params: { type } }).then((r) => r.data);

export const createCategory = (data) =>
  api.post('/categories', data).then((r) => r.data);

export const updateCategory = (id, data) =>
  api.put(`/categories/${id}`, data).then((r) => r.data);

export const deleteCategory = (id) =>
  api.delete(`/categories/${id}`).then((r) => r.data);

// Transactions
export const getTransactions = (params) =>
  api.get('/transactions', { params }).then((r) => r.data);

export const getTransaction = (id) =>
  api.get(`/transactions/${id}`).then((r) => r.data);

export const createTransaction = (data) =>
  api.post('/transactions', data).then((r) => r.data);

export const updateTransaction = (id, data) =>
  api.put(`/transactions/${id}`, data).then((r) => r.data);

export const deleteTransaction = (id) =>
  api.delete(`/transactions/${id}`).then((r) => r.data);

// Reports
export const getSummary = (month, year) =>
  api.get('/reports/summary', { params: { month, year } }).then((r) => r.data);

export const getYearlyTrend = (year) =>
  api.get('/reports/yearly-trend', { params: { year } }).then((r) => r.data);

export const getCategoryTrend = (month, year) =>
  api.get('/reports/category-trend', { params: { month, year } }).then((r) => r.data);

// Savings Goals
export const getSavingsGoals = () =>
  api.get('/savings-goals').then((r) => r.data);

export const createSavingsGoal = (data) =>
  api.post('/savings-goals', data).then((r) => r.data);

export const updateSavingsGoal = (id, data) =>
  api.put(`/savings-goals/${id}`, data).then((r) => r.data);

export const deleteSavingsGoal = (id) =>
  api.delete(`/savings-goals/${id}`).then((r) => r.data);

export const depositToSavingsGoal = (id, amount) =>
  api.post(`/savings-goals/${id}/deposit`, { amount }).then((r) => r.data);

export const withdrawFromSavingsGoal = (id, amount) =>
  api.post(`/savings-goals/${id}/withdraw`, { amount }).then((r) => r.data);

// Budgets
export const getBudgets = (month, year) =>
  api.get('/budgets', { params: { month, year } }).then((r) => r.data);

export const getBudgetSummary = (month, year) =>
  api.get('/budgets/summary', { params: { month, year } }).then((r) => r.data);

export const upsertBudget = (data) =>
  api.post('/budgets', data).then((r) => r.data);

export const deleteBudget = (id) =>
  api.delete(`/budgets/${id}`).then((r) => r.data);

export default api;
