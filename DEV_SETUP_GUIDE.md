# Complete Deployment Setup Guide (Production + Development)

This guide will help you set up **both** production and development environments from scratch.

## Overview

After this setup, you'll have:
- **Production**: `main` branch → work-track-api.onrender.com (Neon prod DB)
- **Development**: `dev` branch → work-track-api-dev.onrender.com (Neon dev DB)

## Time Estimate: ~30 minutes total

---

## Step 1: Create Databases on Neon.tech

You'll create **two separate databases** - one for production, one for development.

### 1.1 Create Production Database

1. Go to https://console.neon.tech
2. Sign up with GitHub (if you haven't already)
3. Click **"Create a project"**
4. Configure:
   - **Project name**: `work-track-db-prod`
   - **PostgreSQL version**: 16 or 17 (latest available)
   - **Region**: Choose closest to you (e.g., US East, EU West)
5. Click **"Create project"**

### 1.2 Save Production Credentials

⚠️ **IMPORTANT**: Save these immediately! You'll need them later.

```
📝 Production Database Credentials:
Host: ep-xxxxx-prod.region.aws.neon.tech
Database: neondb
User: neondb_owner
Password: npg_xxxxxxxxxxxx

Full Connection String:
postgresql://neondb_owner:npg_xxx@ep-xxx-prod.region.aws.neon.tech/neondb?sslmode=require
```

### 1.3 Create Development Database

1. Still on Neon dashboard, click **"Create a project"** again
2. Configure:
   - **Project name**: `work-track-db-dev`
   - **PostgreSQL version**: 16 or 17 (same as prod)
   - **Region**: Same as prod (for consistency)
3. Click **"Create project"**

### 1.4 Save Development Credentials

```
📝 Development Database Credentials:
Host: ep-xxxxx-dev.region.aws.neon.tech
Database: neondb
User: neondb_owner
Password: npg_xxxxxxxxxxxx

Full Connection String:
postgresql://neondb_owner:npg_xxx@ep-xxx-dev.region.aws.neon.tech/neondb?sslmode=require
```

### 1.5 Run Migrations on BOTH Databases

Navigate to your project:
```bash
cd /Users/sergey/Documents/Projects/work_track/backend
```

**Production Database:**
```bash
psql "postgresql://neondb_owner:npg_moHD7BiTSZp3@ep-bold-cell-agoria02-pooler.c-2.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require" \
  -f migrations/000001_create_users_table.up.sql

psql "postgresql://neondb_owner:npg_moHD7BiTSZp3@ep-bold-cell-agoria02-pooler.c-2.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require" \
  -f migrations/000002_create_tasks_table.up.sql
```

**Development Database:**
```bash
psql "postgresql://neondb_owner:npg_aYLXTRc0pE1o@ep-holy-poetry-agsaixa1-pooler.c-2.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require" \
  -f migrations/000001_create_users_table.up.sql

psql "postgresql://neondb_owner:npg_aYLXTRc0pE1o@ep-holy-poetry-agsaixa1-pooler.c-2.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require" \
  -f migrations/000002_create_tasks_table.up.sql
```

**Verify both databases:**
```bash
# Check prod
psql "postgresql://neondb_owner:npg_moHD7BiTSZp3@ep-bold-cell-agoria02-pooler.c-2.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require" -c "\dt"

# Check dev
psql "postgresql://neondb_owner:npg_aYLXTRc0pE1o@ep-holy-poetry-agsaixa1-pooler.c-2.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require" -c "\dt"

# Both should show: users, track_items
```

✅ **Both databases are ready!**

---

## Step 2: Create Dev Branch

```bash
# Create and switch to dev branch
git checkout -b dev

# Push to GitHub
git push -u origin dev
```

---

## Step 3: Set Up Render Services (Production + Development)

You'll create **two separate web services** on Render.

### 3.1 Create Render Account

1. Go to https://render.com
2. Click **"Get Started"**
3. Sign up with your **GitHub account**
4. Authorize Render to access your repositories

### 3.2 Create Production Service

1. From Render Dashboard, click **"New +"** → **"Web Service"**
2. Find and select your `work-track-backend` repository
3. Click **"Connect"**

**Configure Production Service:**

| Setting | Value |
|---------|-------|
| **Name** | `work-track-api` |
| **Region** | Choose closest to you |
| **Branch** | `main` |
| **Runtime** | **Docker** |
| **Instance Type** | **Free** |

**Add Production Environment Variables:**

| Key | Value | Notes |
|-----|-------|-------|
| `PORT` | `8080` | |
| `ENV` | `production` | |
| `DB_HOST` | `ep-xxx-prod.neon.tech` | From Neon PROD |
| `DB_PORT` | `5432` | |
| `DB_USER` | `neondb_owner` | From Neon PROD |
| `DB_PASSWORD` | `npg_xxx` | From Neon PROD |
| `DB_NAME` | `neondb` | |
| `DB_SSLMODE` | `require` | Critical for Neon |
| `JWT_SECRET` | Generate with command below | |
| `ALLOWED_ORIGINS` | `https://your-frontend-url.com` | Update later |

**Generate JWT_SECRET for production:**
```bash
openssl rand -base64 32
```

Click **"Create Web Service"** and wait for deployment (2-3 minutes).

### 3.3 Create Development Service

1. Click **"New +"** → **"Web Service"** again
2. Select your `work-track-backend` repository
3. Click **"Connect"**

**Configure Development Service:**

| Setting | Value |
|---------|-------|
| **Name** | `work-track-api-dev` |
| **Region** | Same as prod |
| **Branch** | `dev` |
| **Runtime** | **Docker** |
| **Instance Type** | **Free** |

**Add Development Environment Variables:**

| Key | Value | Notes |
|-----|-------|-------|
| `PORT` | `8080` | |
| `ENV` | `development` | **Different from prod** |
| `DB_HOST` | `ep-xxx-dev.neon.tech` | From Neon DEV |
| `DB_PORT` | `5432` | |
| `DB_USER` | `neondb_owner` | From Neon DEV |
| `DB_PASSWORD` | `npg_xxx` | From Neon DEV |
| `DB_NAME` | `neondb` | |
| `DB_SSLMODE` | `require` | |
| `JWT_SECRET` | Generate different one | **Different from prod** |
| `ALLOWED_ORIGINS` | `http://localhost:5173` | Dev frontend |

**Generate different JWT_SECRET for development:**
```bash
openssl rand -base64 32
```

Click **"Create Web Service"** and wait for deployment.

### 3.4 Test Both Deployments

**Production:**
```bash
curl https://work-track-api.onrender.com/health
# Should return: OK
```

**Development:**
```bash
curl https://work-track-api-dev.onrender.com/health
# Should return: OK
```

✅ **Both services are deployed!**

---

## Step 4: Get Deploy Hook URLs

You'll need deploy hooks for both environments.

### 4.1 Production Deploy Hook

1. Go to Render dashboard → `work-track-api` (prod service)
2. Click **"Settings"**
3. Scroll to **"Deploy Hook"**
4. Copy the URL (looks like: `https://api.render.com/deploy/srv-xxxxx?key=yyyyy`)
5. Save as: **PROD_DEPLOY_HOOK_URL**

### 4.2 Development Deploy Hook

1. Go to Render dashboard → `work-track-api-dev`
2. Click **"Settings"**
3. Scroll to **"Deploy Hook"**
4. Copy the URL
5. Save as: **DEV_DEPLOY_HOOK_URL**

---

## Step 5: Configure GitHub Secrets

1. Go to your GitHub repository
2. Click **"Settings"** → **"Secrets and variables"** → **"Actions"**
3. Add/update these secrets:

| Secret Name | Value | Purpose |
|-------------|-------|---------|
| `RENDER_DEPLOY_HOOK_URL` | Prod deploy hook | Deploys main branch |
| `RENDER_DEPLOY_HOOK_URL_DEV` | Dev deploy hook | Deploys dev branch |

---

## Summary

After completing these steps, you'll have:

```
┌─────────────────┐
│  GitHub Repo    │
└────────┬────────┘
         │
    ┌────┴─────┐
    │          │
┌───▼────┐ ┌──▼──────┐
│  main  │ │   dev   │
└───┬────┘ └──┬──────┘
    │         │
┌───▼─────────▼───┐
│ GitHub Actions  │
└───┬─────────┬───┘
    │         │
┌───▼────┐ ┌──▼──────────┐
│ Render │ │ Render Dev  │
│  Prod  │ │             │
└───┬────┘ └──┬──────────┘
    │         │
┌───▼────┐ ┌──▼──────────┐
│ Neon   │ │  Neon Dev   │
│ Prod   │ │             │
└────────┘ └─────────────┘
```

**Next:** Update CI/CD workflow to deploy to the correct environment based on branch.
