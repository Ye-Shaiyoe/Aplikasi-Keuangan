# 📖 Technical Documentation — Catatan Keuangan

# 技术文档 — 财务笔记 (Personal Finance Tracker)

<div align="center">

[🌐 Live Demo](https://umkmkeuangan.vercel.app) · [📕 README](./README.md)

</div>

---

## Table of Contents / 目录

1. [Project Overview / 项目概述](#1-project-overview--项目概述)
2. [Architecture / 架构](#2-architecture--架构)
3. [Folder Structure / 目录结构](#3-folder-structure--目录结构)
4. [Database Schema / 数据库模式](#4-database-schema--数据库模式)
5. [API Reference / API 参考](#5-api-reference--api-参考)
6. [Authentication Flow / 认证流程](#6-authentication-flow--认证流程)
7. [Feature Modules / 功能模块](#7-feature-modules--功能模块)
8. [State Management / 状态管理](#8-state-management--状态管理)
9. [Responsive Design / 响应式设计](#9-responsive-design--响应式设计)
10. [Deployment Guide / 部署指南](#10-deployment-guide--部署指南)
11. [Troubleshooting / 故障排除](#11-troubleshooting--故障排除)

---

## 1. Project Overview / 项目概述

**English:**
Catatan Keuangan is a full-stack personal finance web application that allows users to:
- Track daily income and expense transactions with categorization
- Set up monthly budgets per category and monitor spending
- Create savings goals with deposit/withdrawal tracking and visual progress
- Configure recurring transactions (daily/weekly/monthly/yearly) that auto-generate real transactions
- View interactive charts and detailed tables for financial analysis
- Authenticate via email/password or Google OAuth 2.0
- Use the app in light or dark mode with a fully responsive layout

**中文:**
Catatan Keuangan（财务笔记）是一个全栈个人财务 Web 应用，允许用户：
- 追踪每日收入和支出交易，支持分类管理
- 设置按类别的月度预算并监控支出
- 创建储蓄目标，支持存款/取款追踪和可视化进度
- 配置定期交易（每天/每周/每月/每年），自动生成真实交易
- 查看交互式图表和详细表格进行财务分析
- 通过邮箱/密码或 Google OAuth 2.0 登录
- 在浅色或深色模式下使用，完全响应式布局

---

## 2. Architecture / 架构

### High-Level / 高层架构

```
┌─────────────────────────────────────────────────┐
│  Browser (React SPA)                            │
│  - React Router for client-side routing         │
│  - Axios for HTTP with JWT interceptor           │
│  - Recharts for data visualization              │
│  - Tailwind CSS for styling                     │
│  - @react-oauth/google for Google Sign-In       │
└──────────────────────┬──────────────────────────┘
                       │ REST API (JSON)
                       │ Authorization: Bearer <JWT>
                       ▼
┌─────────────────────────────────────────────────┐
│  Go Backend (Gin Framework)                     │
│                                                 │
│  Router → CORS Middleware → Auth Middleware     │
│       → Handler → Service → Repository          │
│                             → PostgreSQL (pgx)  │
└──────────────────────┬──────────────────────────┘
                       │ SQL (pgx)
                       ▼
┌─────────────────────────────────────────────────┐
│  PostgreSQL 16                                  │
│  - users                                        │
│  - categories                                   │
│  - transactions                                 │
│  - budgets                                      │
│  - savings_goals                                │
│  - recurring_transactions                       │
└─────────────────────────────────────────────────┘
```

### Backend Layer Responsibilities / 后端各层职责

| Layer / 层 | File Location / 文件位置 | Responsibility / 职责 |
|-----|-----------|----------------|
| **Handler** | `internal/handler/*.go` | HTTP request parsing, validation, response formatting / HTTP 请求解析、验证、响应格式化 |
| **Service** | `internal/service/*.go` | Business logic, orchestration between repos / 业务逻辑、协调仓储层 |
| **Repository** | `internal/repository/*.go` | Database queries (pgx), raw SQL / 数据库查询（pgx）、原生 SQL |
| **Model** | `internal/model/*.go` | Data structs, request/response DTOs / 数据结构、请求/响应 DTO |
| **Middleware** | `internal/middleware/*.go` | JWT verification, user context injection / JWT 验证、用户上下文注入 |

---

## 3. Folder Structure / 目录结构

```
finance-app/
├── backend/
│   ├── cmd/server/
│   │   └── main.go              # Entry point / 入口
│   ├── internal/
│   │   ├── database/
│   │   │   └── postgres.go      # Connection + migrations / 连接与迁移
│   │   ├── handler/
│   │   │   ├── auth.go          # Register, Login, GoogleLogin, Me / 认证处理器
│   │   │   ├── budget.go        # Budget CRUD / 预算 CRUD
│   │   │   ├── category.go      # Category CRUD / 分类 CRUD
│   │   │   ├── recurring.go     # Recurring transactions / 定期交易
│   │   │   ├── savings_goal.go  # Savings goals / 储蓄目标
│   │   │   └── transaction.go   # Transaction CRUD / 交易 CRUD
│   │   ├── middleware/
│   │   │   └── auth.go          # JWT middleware + token gen / JWT 中间件
│   │   ├── model/
│   │   │   ├── budget.go
│   │   │   ├── category.go
│   │   │   ├── recurring.go
│   │   │   ├── savings_goal.go
│   │   │   ├── transaction.go
│   │   │   └── user.go
│   │   ├── repository/
│   │   │   ├── budget.go
│   │   │   ├── category.go
│   │   │   ├── recurring.go
│   │   │   ├── savings_goal.go
│   │   │   ├── transaction.go
│   │   │   └── user.go
│   │   └── service/
│   │       ├── budget.go
│   │       ├── category.go
│   │       ├── recurring.go
│   │       ├── savings_goal.go
│   │       └── transaction.go
│   ├── .env
│   ├── .env.example
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── public/                  # Static assets / 静态资源
│   ├── src/
│   │   ├── api/
│   │   │   └── client.js        # Axios instance + API functions / Axios 实例与 API 函数
│   │   ├── assets/              # Images / 图片
│   │   ├── components/
│   │   │   ├── Card.jsx         # Reusable card / 可复用卡片
│   │   │   ├── Layout.jsx       # Sidebar + topbar + bottom nav / 布局组件
│   │   │   └── Modal.jsx        # Modal dialog / 模态框
│   │   ├── context/
│   │   │   ├── AuthContext.jsx  # Auth state + JWT management / 认证状态
│   │   │   └── ThemeContext.jsx # Dark mode toggle / 深色模式
│   │   ├── pages/
│   │   │   ├── Budgets.jsx      # Budget management / 预算管理
│   │   │   ├── Categories.jsx   # Category management / 分类管理
│   │   │   ├── Dashboard.jsx    # Main dashboard / 仪表盘
│   │   │   ├── Docs.jsx         # Docs page / 文档页
│   │   │   ├── InsightCharts.jsx # Chart visualization / 图表洞察
│   │   │   ├── InsightTables.jsx # Table visualization / 表格洞察
│   │   │   ├── Login.jsx        # Login + Google Sign-In / 登录
│   │   │   ├── Recurring.jsx    # Recurring transactions / 定期交易
│   │   │   ├── Register.jsx     # Register + Google Sign-Up / 注册
│   │   │   ├── Reports.jsx      # Monthly reports / 月度报告
│   │   │   ├── SavingsGoals.jsx # Savings goals / 储蓄目标
│   │   │   └── Transactions.jsx # Transaction list / 交易列表
│   │   ├── utils/
│   │   │   └── format.js        # Currency/date formatters / 格式化工具
│   │   ├── App.jsx              # Routes + providers / 路由与 Provider
│   │   ├── index.css            # Global styles / 全局样式
│   │   └── main.jsx             # Entry point / 入口
│   ├── .env
│   ├── package.json
│   ├── tailwind.config.js
│   └── vite.config.js
│
├── .gitignore
├── README.md
├── DOKUMENTASI.md               # This file / 本文件
└── start.bat
```

---

## 4. Database Schema / 数据库模式

### `users`
| Column | Type | Constraints |
|--------|------|-------------|
| `id` | SERIAL | PRIMARY KEY |
| `name` | VARCHAR(100) | NOT NULL |
| `email` | VARCHAR(150) | NOT NULL, UNIQUE |
| `password_hash` | VARCHAR(255) | NULLABLE (null for Google-only users) |
| `google_id` | VARCHAR(255) | NULLABLE, UNIQUE (partial index) |
| `created_at` | TIMESTAMP | DEFAULT NOW() |

### `categories`
| Column | Type | Constraints |
|--------|------|-------------|
| `id` | SERIAL | PRIMARY KEY |
| `user_id` | INTEGER | NOT NULL, FK → users(id) ON DELETE CASCADE |
| `name` | VARCHAR(100) | NOT NULL |
| `type` | VARCHAR(10) | NOT NULL, CHECK ('income' or 'expense') |
| `icon` | VARCHAR(50) | DEFAULT '' |
| `color` | VARCHAR(7) | DEFAULT '#6b7280' |
| `created_at` | TIMESTAMP | DEFAULT NOW() |

### `transactions`
| Column | Type | Constraints |
|--------|------|-------------|
| `id` | SERIAL | PRIMARY KEY |
| `user_id` | INTEGER | NOT NULL, FK → users(id) |
| `category_id` | INTEGER | NOT NULL, FK → categories(id) |
| `amount` | NUMERIC(15,2) | NOT NULL |
| `description` | TEXT | DEFAULT '' |
| `date` | DATE | NOT NULL |
| `type` | VARCHAR(10) | NOT NULL, CHECK ('income' or 'expense') |
| `created_at` | TIMESTAMP | DEFAULT NOW() |
| `updated_at` | TIMESTAMP | DEFAULT NOW() |

### `budgets`
| Column | Type | Constraints |
|--------|------|-------------|
| `id` | SERIAL | PRIMARY KEY |
| `user_id` | INTEGER | NOT NULL, FK → users(id) |
| `category_id` | INTEGER | NOT NULL, FK → categories(id) |
| `amount` | NUMERIC(15,2) | NOT NULL |
| `month` | INTEGER | NOT NULL |
| `year` | INTEGER | NOT NULL |
| `created_at` | TIMESTAMP | DEFAULT NOW() |
| `updated_at` | TIMESTAMP | DEFAULT NOW() |

### `savings_goals`
| Column | Type | Constraints |
|--------|------|-------------|
| `id` | SERIAL | PRIMARY KEY |
| `user_id` | INTEGER | NOT NULL, FK → users(id) |
| `name` | VARCHAR(100) | NOT NULL |
| `target_amount` | NUMERIC(15,2) | NOT NULL |
| `current_amount` | NUMERIC(15,2) | DEFAULT 0 |
| `deadline` | DATE | NULLABLE |
| `created_at` | TIMESTAMP | DEFAULT NOW() |
| `updated_at` | TIMESTAMP | DEFAULT NOW() |

### `recurring_transactions`
| Column | Type | Constraints |
|--------|------|-------------|
| `id` | SERIAL | PRIMARY KEY |
| `user_id` | INTEGER | NOT NULL, FK → users(id) |
| `category_id` | INTEGER | NOT NULL, FK → categories(id) |
| `amount` | NUMERIC(15,2) | NOT NULL |
| `description` | TEXT | DEFAULT '' |
| `type` | VARCHAR(10) | NOT NULL, CHECK ('income' or 'expense') |
| `frequency` | VARCHAR(10) | NOT NULL, CHECK ('daily','weekly','monthly','yearly') |
| `start_date` | DATE | NOT NULL |
| `end_date` | DATE | NULLABLE |
| `next_date` | DATE | NOT NULL |
| `last_processed` | DATE | NULLABLE |
| `is_active` | BOOLEAN | DEFAULT TRUE |
| `created_at` | TIMESTAMP | DEFAULT NOW() |
| `updated_at` | TIMESTAMP | DEFAULT NOW() |

---

## 5. API Reference / API 参考

**Base URL / 基础 URL:** `https://<your-backend>/api`

### Authentication / 认证 (Public / 公开)

| Method | Endpoint | Description / 描述 | Body |
|--------|----------|------|------|
| POST | `/auth/register` | Register new user / 注册新用户 | `{name, email, password}` |
| POST | `/auth/login` | Login with email + password / 邮箱密码登录 | `{email, password}` |
| POST | `/auth/google` | Login/Register with Google ID token / Google 登录 | `{credential}` |

**Response / 响应** (all auth endpoints / 所有认证端点):
```json
{
  "token": "eyJhbGciOi...",
  "user": {
    "id": 1,
    "name": "Akrom",
    "email": "user@example.com",
    "google_id": "123456...",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### Protected Endpoints / 受保护端点 (Require `Authorization: Bearer <token>`)

| Resource | Endpoints |
|----------|-----------|
| **Profile / 个人资料** | `GET /auth/me` |
| **Categories / 分类** | `GET/POST /categories`, `GET/PUT/DELETE /categories/:id` |
| **Transactions / 交易** | `GET/POST /transactions`, `GET/PUT/DELETE /transactions/:id` |
| **Budgets / 预算** | `GET /budgets`, `GET /budgets/summary`, `POST /budgets`, `DELETE /budgets/:id` |
| **Savings Goals / 储蓄** | `GET/POST /savings-goals`, `GET/PUT/DELETE /savings-goals/:id`, `POST /savings-goals/:id/deposit`, `POST /savings-goals/:id/withdraw` |
| **Recurring / 定期交易** | `GET/POST /recurring`, `GET/PUT/DELETE /recurring/:id`, `POST /recurring/process` |
| **Reports / 报告** | `GET /reports/summary`, `GET /reports/yearly-trend`, `GET /reports/category-trend` |

---

## 6. Authentication Flow / 认证流程

### Email + Password / 邮箱 + 密码

```
[Client]                              [Backend]                          [PostgreSQL]
   │                                     │                                  │
   │ POST /auth/register                 │                                  │
   │ {name, email, password}             │                                  │
   │ ────────────────────────────────────▶│                                  │
   │                                     │ SELECT by email (check dup)      │
   │                                     │ ─────────────────────────────────▶│
   │                                     │                                  │
   │                                     │ bcrypt hash password             │
   │                                     │ INSERT user                      │
   │                                     │ ─────────────────────────────────▶│
   │                                     │                                  │
   │                                     │ Seed 13 default categories       │
   │                                     │ Generate JWT (72h)               │
   │ ◀────────────────────────────────────│                                  │
   │ {token, user}                       │                                  │
```

### Google OAuth 2.0 (ID Token Flow) / Google OAuth (ID 令牌流程)

```
[Browser]              [Google]              [Backend]              [PostgreSQL]
    │                     │                      │                      │
    │ Click "Sign in      │                      │                      │
    │ with Google"        │                      │                      │
    │ ───────────────────▶│                      │                      │
    │                     │                      │                      │
    │ ◀───────────────────│                      │                      │
    │ Google popup        │                      │                      │
    │ → ID token          │                      │                      │
    │                     │                      │                      │
    │ POST /auth/google   │                      │                      │
    │ {credential}        │                      │                      │
    │ ───────────────────────────────────────────▶│                      │
    │                     │                      │                      │
    │                     │                      │ idtoken.Validate()   │
    │                     │                      │ (checks signature,   │
    │                     │◀─────────────────────│ audience, expiry)    │
    │                     │                      │                      │
    │                     │                      │ Extract sub, email,  │
    │                     │                      │ name, email_verified │
    │                     │                      │                      │
    │                     │                      │ GetUserByGoogleID    │
    │                     │                      │ ─────────────────────▶│
    │                     │                      │                      │
    │                     │                      │ [if not found]       │
    │                     │                      │ GetUserByEmail       │
    │                     │                      │ ─────────────────────▶│
    │                     │                      │ → LinkGoogleID       │
    │                     │                      │                      │
    │                     │                      │ [if still not found] │
    │                     │                      │ CreateUserGoogle     │
    │                     │                      │ + SeedCategories     │
    │                     │                      │ ─────────────────────▶│
    │                     │                      │                      │
    │                     │                      │ GenerateToken (JWT)  │
    │ ◀──────────────────────────────────────────│                      │
    │ {token, user}       │                      │                      │
```

**Account Linking Behavior / 账户关联行为:**
If a user registers with email+password first, then later clicks "Sign in with Google" with the same email, the backend **links the Google account** to the existing user record. Subsequent Google logins work seamlessly. /
如果用户先用邮箱+密码注册，后来用相同邮箱点击"Google 登录"，后端会**将 Google 账户关联**到现有用户记录。后续 Google 登录将无缝工作。

---

## 7. Feature Modules / 功能模块

### 7.1 Recurring Transactions / 定期交易

**English:**
Recurring transactions allow users to automate repetitive income/expenses. The system processes due items in two ways:
1. **Auto-processing on server startup:** `service.ProcessDueRecurring(0)` is called in `main.go` after migrations. It iterates through all users' due items, creates real transactions, and advances `next_date`.
2. **Manual trigger:** The frontend's "Process Now" button calls `POST /api/recurring/process`, which scopes to the current user only.

**Frequency advancement logic / 频率推进逻辑:**
- daily → +1 day
- weekly → +7 days
- monthly → +1 month
- yearly → +1 year

If `end_date` is set and the new `next_date` exceeds it, `is_active` is set to `FALSE` automatically.

### 7.2 Insights / 洞察

Split into two pages / 分为两个页面:
- **Chart Data** (`/insights/charts`) — Recharts visualizations: yearly trend (ComposedChart), cash flow (AreaChart), expense ratio (dual-axis), top categories (horizontal Bar), radar, pie distributions.
- **Table Data** (`/insights/tables`) — Category breakdown with weekly sparklines, full transaction list with totals, daily summary with running balance.

Shared formatters (`MONTH_NAMES`, `formatRupiah`, `formatCompact`) live in `frontend/src/utils/format.js`.

### 7.3 Mobile Navigation / 移动导航

**English:**
The mobile bottom nav uses a 5-tab layout with dropdown popups:
- Dashboard (link) | Transaksi (dropdown: Transaksi, Transaksi Berulang) | Anggaran (link) | Insight (dropdown: Chart Data, Table Data) | Lainnya (dropdown: Tabungan, Laporan, Dokumentasi)

Only one dropdown is open at a time. Popups are anchored `left-0` for left-side tabs (idx<2) and `right-0` for right-side tabs to avoid overflow. The desktop sidebar uses `navItems` directly with collapsible sections.

**中文:**
移动端底部导航使用 5 个标签的下拉弹出布局：
- 仪表盘（链接）| 交易（下拉：交易、定期交易）| 预算（链接）| 洞察（下拉：图表数据、表格数据）| 更多（下拉：储蓄、报告、文档）

一次只能打开一个下拉菜单。左侧标签（idx<2）的弹窗锚定 `left-0`，右侧标签锚定 `right-0` 以避免溢出。

---

## 8. State Management / 状态管理

### `AuthContext`
- Stores `user`, `token`, `isAuthenticated` / 存储用户、令牌、认证状态
- On mount: verifies token via `GET /auth/me` / 挂载时通过 API 验证令牌
- `login(authData)` / `register(authData)` — both save token to localStorage / 都将令牌保存到 localStorage
- `logout()` — clears token and user / 清除令牌和用户

### `ThemeContext`
- Stores `dark` boolean / 存储深色模式标志
- `toggle()` flips and persists to `localStorage` / 切换并持久化到 localStorage
- Applies `class="dark"` to `<html>` element for Tailwind dark mode / 在 `<html>` 元素上应用 dark 类

### Axios Interceptors / Axios 拦截器 (`api/client.js`)
- **Request:** Attaches `Authorization: Bearer <token>` to every request / 每个请求附加 Authorization 头
- **Response:** On 401, clears token and redirects to `/login` / 遇到 401 时清除令牌并重定向

---

## 9. Responsive Design / 响应式设计

| Breakpoint | Layout / 布局 |
|-----------|--------|
| `< md` (mobile) | Top header + content + bottom nav with 5 tabs (3 are dropdowns) / 顶部标题 + 内容 + 底部 5 标签导航 |
| `≥ md` (desktop) | Fixed left sidebar (256px) + content area / 固定左侧边栏 + 内容区 |

**Tailwind dark mode / Tailwind 深色模式:** `class` strategy — toggled via ThemeContext, persisted to localStorage.

---

## 10. Deployment Guide / 部署指南

### Frontend (Vercel)
1. Connect your GitHub repo to Vercel / 将 GitHub 仓库连接到 Vercel
2. Set root directory to `frontend` / 设置根目录为 frontend
3. Add env vars: `VITE_API_URL`, `VITE_GOOGLE_CLIENT_ID` / 添加环境变量
4. Deploy — Vite builds at deploy time, env vars are baked in / 部署时 Vite 构建，环境变量在构建时注入

### Backend (Railway)
1. Create Railway project, link the `backend/` folder / 创建 Railway 项目，链接 backend 文件夹
2. Add PostgreSQL plugin / 添加 PostgreSQL 插件
3. Set env vars: `DATABASE_URL` (from Postgres plugin), `GOOGLE_CLIENT_ID`, `JWT_SECRET` / 设置环境变量
4. Railway auto-detects Go and deploys / Railway 自动检测 Go 并部署

### Google OAuth Console Setup / Google OAuth 控制台配置
1. Visit [Google Cloud Console](https://console.cloud.google.com/)
2. Create/select a project / 创建或选择项目
3. Go to **APIs & Services → Credentials** / 前往 API 与服务 → 凭据
4. Configure **OAuth consent screen** (External, add your email) / 配置 OAuth 同意屏幕
5. **Create Credentials → OAuth client ID → Web application** / 创建凭据 → OAuth 客户端 ID → Web 应用
6. Add to **Authorized JavaScript origins** / 添加到授权的 JavaScript 来源:
   - `http://localhost:5173`
   - `https://umkmkeuangan.vercel.app`
7. Copy Client ID → paste into backend + frontend env vars / 复制客户端 ID → 粘贴到后端和前端环境变量

---

## 11. Troubleshooting / 故障排除

### "server misconfigured" on Google login / Google 登录时显示"服务器配置错误"
**Cause / 原因:** `GOOGLE_CLIENT_ID` is not set on the backend (Railway). / 后端（Railway）未设置 GOOGLE_CLIENT_ID。
**Fix / 修复:** Add `GOOGLE_CLIENT_ID` to Railway service **Variables** tab (make sure it's on the Go backend service, not the Postgres service). / 在 Railway 服务的"变量"选项卡中添加。

### Google button not showing / Google 按钮不显示
**Cause / 原因:** `VITE_GOOGLE_CLIENT_ID` is not set on Vercel, OR deployment hasn't rebuilt. / Vercel 上未设置环境变量，或部署未重建。
**Fix / 修复:** Add env var in Vercel → Settings → Env Vars, then **Redeploy** (uncheck "Use existing Build Cache"). Vite env vars are baked at build time. / 在 Vercel 中添加环境变量，然后重新部署。

### "token google tidak valid" / "Google 令牌无效"
**Cause / 原因:** The Google OAuth client's "Authorized JavaScript origins" doesn't include your domain, or the Client ID on backend doesn't match the frontend. / Google OAuth 客户端的"授权 JavaScript 来源"不包含您的域名，或前后端 Client ID 不匹配。
**Fix / 修复:** Ensure the same Client ID is in both `GOOGLE_CLIENT_ID` (backend) and `VITE_GOOGLE_CLIENT_ID` (frontend), and your domain is in Google Cloud Console. / 确保前后端使用相同的 Client ID，且域名已在 Google Cloud Console 中添加。

### `cannot scan NULL into *string` on recurring transactions / 定期交易出现 NULL 扫描错误
**Cause / 原因:** `end_date` or `last_processed` columns are NULL but scanned into Go `string`. / 这些列为 NULL 但被扫描为 Go 字符串。
**Fix / 修复:** Wrap with `COALESCE(column::text, '') AS column` in queries. Already fixed in current code. / 在查询中使用 COALESCE 包装。当前代码已修复。

### npm install fails with SSL error / npm 安装因 SSL 错误失败
**Fix / 修复:**
```powershell
$env:NODE_OPTIONS='--use-system-ca'
npm install <package>
```

---

<div align="center">

Built with ❤️ using Go, React, PostgreSQL, and Tailwind CSS.
使用 Go、React、PostgreSQL 和 Tailwind CSS，用 ❤️ 构建。

[🌐 Live Demo](https://umkmkeuangan.vercel.app) · [📕 README](./README.md)

</div>
