# Deployment Guide: Render.com + Neon.tech

This guide will walk you through deploying your Go backend to Render.com with Neon.tech PostgreSQL database and GitHub Actions CI/CD.

## Architecture Overview

```
GitHub (Code) → GitHub Actions (CI/CD) → Render.com (Backend API)
                                              ↓
                                         Neon.tech (PostgreSQL Database)
```

## Prerequisites

- GitHub account
- Render.com account (free)
- Neon.tech account (free)

---

## Part 1: Set Up Neon.tech Database

### 1. Create Neon Account

1. Go to https://neon.tech
2. Click **"Sign Up"** and use your GitHub account
3. Verify your email

### 2. Create Database Project

1. Click **"Create a project"**
2. Configure:
   - **Project name**: `work-track-db`
   - **PostgreSQL version**: 16 (latest)
   - **Region**: Choose closest to you (e.g., US East, EU West)
3. Click **"Create project"**

### 3. Get Connection Details

After creation, you'll see the connection details. **Save these**:

```
Host: ep-cool-name-123456.us-east-2.aws.neon.tech
Database: neondb
User: neondb_owner
Password: npg_xxxxxxxxxxxx
```

**Connection String** (you'll need this):
```
postgresql://neondb_owner:npg_xxxxxxxxxxxx@ep-cool-name-123456.us-east-2.aws.neon.tech/neondb?sslmode=require
```

### 4. Run Database Migrations

Using the connection string from Neon, run migrations locally:

```bash
# Navigate to your project
cd /Users/sergey/Documents/Projects/work_track/backend

# Run migrations (replace with your actual connection string)
psql "postgresql://neondb_owner:YOUR_PASSWORD@YOUR_HOST.neon.tech/neondb?sslmode=require" \
  -f migrations/000001_create_users_table.up.sql

psql "postgresql://neondb_owner:YOUR_PASSWORD@YOUR_HOST.neon.tech/neondb?sslmode=require" \
  -f migrations/000002_create_tasks_table.up.sql
```

**Verify migrations worked:**
```bash
# Connect to database
psql "postgresql://neondb_owner:YOUR_PASSWORD@YOUR_HOST.neon.tech/neondb?sslmode=require"

# List tables
\dt

# Should see: users, track_items
# Exit with: \q
```

✅ **Database is now ready!**

---

## Part 2: Push Code to GitHub

### 1. Initialize Git Repository

```bash
cd /Users/sergey/Documents/Projects/work_track/backend

# Initialize git
git init

# Add all files
git add .

# Create initial commit
git commit -m "Initial commit: Go backend with Neon database"
```

### 2. Create GitHub Repository

1. Go to https://github.com/new
2. Configure:
   - **Repository name**: `work-track-backend`
   - **Description**: "Go backend for work tracking application"
   - **Visibility**: Private (recommended) or Public
   - **Do NOT** check "Add README" or ".gitignore" (we already have these)
3. Click **"Create repository"**

### 3. Push to GitHub

```bash
# Add GitHub remote (replace YOUR_USERNAME with your GitHub username)
git remote add origin https://github.com/YOUR_USERNAME/work-track-backend.git

# Push to GitHub
git branch -M main
git push -u origin main
```

✅ **Code is now on GitHub!**

---

## Part 3: Deploy to Render.com

### 1. Create Render Account

1. Go to https://render.com
2. Click **"Get Started"**
3. Sign up with your **GitHub account** (recommended)
4. Authorize Render to access your repositories

### 2. Create Web Service

1. From Render Dashboard, click **"New +"** → **"Web Service"**
2. Click **"Connect account"** if not already connected
3. Find and select your `work-track-backend` repository
4. Click **"Connect"**

### 3. Configure Service

Fill in the following settings:

| Setting | Value |
|---------|-------|
| **Name** | `work-track-api` |
| **Region** | Same as your Neon database region |
| **Branch** | `main` |
| **Runtime** | **Docker** |
| **Instance Type** | **Free** |

### 4. Add Environment Variables

Scroll down to **"Environment Variables"** and add these:

| Key | Value | Notes |
|-----|-------|-------|
| `PORT` | `8080` | |
| `ENV` | `production` | |
| `DB_HOST` | `ep-xxx.neon.tech` | Copy from Neon connection details |
| `DB_PORT` | `5432` | |
| `DB_USER` | `neondb_owner` | Copy from Neon |
| `DB_PASSWORD` | `npg_xxxxx` | Copy from Neon (keep secret!) |
| `DB_NAME` | `neondb` | Copy from Neon |
| `DB_SSLMODE` | `require` | Important for Neon |
| `JWT_SECRET` | Generate random string | See below |
| `ALLOWED_ORIGINS` | `http://localhost:5173` | Update later with frontend URL |

**Generate JWT_SECRET:**
```bash
# On Mac/Linux
openssl rand -base64 32

# Or use online generator: https://generate-secret.vercel.app/32
```

### 5. Deploy

1. Click **"Create Web Service"**
2. Render will:
   - Clone your repository
   - Build Docker image (takes 2-3 minutes)
   - Start your application
3. Watch the logs for any errors

### 6. Get Your API URL

Once deployed, your API will be available at:
```
https://work-track-api.onrender.com
```

(The exact URL will be shown in your Render dashboard)

### 7. Test Your Deployment

```bash
# Health check
curl https://work-track-api.onrender.com/health
# Should return: OK

# Register a test user
curl -X POST https://work-track-api.onrender.com/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Test",
    "last_name": "User",
    "login": "testuser",
    "password": "password123"
  }'

# Should return user data and JWT token
```

✅ **Backend is live!**

---

## Part 4: Set Up GitHub Actions CI/CD

### 1. Get Render Deploy Hook

1. Go to your Render service dashboard
2. Click **"Settings"** tab
3. Scroll to **"Deploy Hook"** section
4. Click **"Create Deploy Hook"**
5. Copy the URL (looks like: `https://api.render.com/deploy/srv-xxxxx?key=yyyyy`)

### 2. Add GitHub Secret

1. Go to your GitHub repository
2. Click **"Settings"** → **"Secrets and variables"** → **"Actions"**
3. Click **"New repository secret"**
4. Configure:
   - **Name**: `RENDER_DEPLOY_HOOK_URL`
   - **Secret**: Paste the Deploy Hook URL
5. Click **"Add secret"**

### 3. How CI/CD Works

Now, every time you push to `main` branch:

1. ✅ GitHub Actions runs tests
2. ✅ Builds Docker image
3. ✅ If tests pass, triggers Render deployment
4. ✅ Render rebuilds and redeploys your app
5. ✅ Changes are live in ~2-3 minutes

### 4. Test CI/CD

```bash
# Make a small change
echo "# Test deployment" >> README.md

# Commit and push
git add README.md
git commit -m "Test CI/CD pipeline"
git push origin main

# Watch GitHub Actions: https://github.com/YOUR_USERNAME/work-track-backend/actions
# Watch Render logs: https://dashboard.render.com
```

✅ **CI/CD is working!**

---

## Part 5: Update for Frontend Integration

### When You Deploy Your Vue 3 Frontend

1. Go to Render dashboard → Your service → **"Environment"**
2. Update `ALLOWED_ORIGINS`:
   ```
   https://your-vue-app.vercel.app,http://localhost:5173
   ```
3. Click **"Save Changes"** (triggers automatic redeploy)

---

## Monitoring & Maintenance

### View Logs

**Render logs:**
1. Go to Render dashboard → Your service
2. Click **"Logs"** tab
3. See real-time logs

**Neon database:**
1. Go to Neon dashboard → Your project
2. Click **"Monitoring"** tab
3. See queries, connections, storage

### Database Backups

**Neon automatic backups:**
- Free tier: 7-day point-in-time restore
- Backups are automatic, no setup needed

**Manual backup:**
```bash
# Backup to file
pg_dump "postgresql://user:pass@host.neon.tech/neondb" > backup.sql

# Restore from file
psql "postgresql://user:pass@host.neon.tech/neondb" < backup.sql
```

### Check Service Health

```bash
# Health endpoint
curl https://work-track-api.onrender.com/health

# Database connection test
curl -X POST https://work-track-api.onrender.com/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Test","last_name":"User","login":"test123","password":"pass123"}'
```

---

## Troubleshooting

### Service Won't Start

**Check Render logs:**
1. Dashboard → Service → Logs
2. Look for errors

**Common issues:**
- ❌ Missing environment variables
- ❌ Wrong database credentials
- ❌ Port mismatch

**Solution:**
- Verify all environment variables are set
- Test database connection locally with same credentials

### Database Connection Failed

**Error:** `connection refused` or `timeout`

**Check:**
1. Neon database is running (check Neon dashboard)
2. `DB_SSLMODE=require` is set
3. Credentials are correct
4. Host includes full domain (e.g., `ep-xxx.us-east-2.aws.neon.tech`)

**Test connection locally:**
```bash
psql "postgresql://user:pass@host.neon.tech/neondb?sslmode=require"
```

### Migrations Not Applied

**Symptoms:** API returns errors about missing tables

**Solution:**
```bash
# Re-run migrations
psql "YOUR_NEON_CONNECTION_STRING" -f migrations/000001_create_users_table.up.sql
psql "YOUR_NEON_CONNECTION_STRING" -f migrations/000002_create_tasks_table.up.sql

# Verify tables exist
psql "YOUR_NEON_CONNECTION_STRING" -c "\dt"
```

### Cold Starts (Free Tier)

**Issue:** First request after 15 minutes takes ~30 seconds

**Why:** Render free tier spins down after inactivity

**Solutions:**
1. **Accept it** - Fine for development
2. **Upgrade** - $7/month for always-on
3. **Keep-alive service** - Use UptimeRobot to ping every 10 minutes

---

## Cost Breakdown

### What's Free Forever:

| Service | Free Tier | Limits |
|---------|-----------|--------|
| **Neon.tech** | 3 GB storage | No expiration! |
| **Render.com** | 750 hours/month | Sleeps after 15 min |
| **GitHub** | Unlimited repos | 2000 Actions minutes/month |

### Total Monthly Cost: **$0** 🎉

---

## Next Steps

- [x] ✅ Database on Neon.tech
- [x] ✅ Backend on Render.com
- [x] ✅ CI/CD with GitHub Actions
- [ ] 🔄 Deploy Vue 3 frontend
- [ ] 🔄 Update CORS settings
- [ ] 🔄 Add custom domain (optional)

---

## Useful Commands

```bash
# Generate JWT secret
openssl rand -base64 32

# Test API health
curl https://your-service.onrender.com/health

# Connect to Neon database
psql "postgresql://user:pass@host.neon.tech/neondb?sslmode=require"

# View database tables
psql "YOUR_CONNECTION_STRING" -c "\dt"

# Force redeploy
curl -X POST "YOUR_DEPLOY_HOOK_URL"

# View recent commits
git log --oneline -5
```

---

## Resources

- [Neon Documentation](https://neon.tech/docs)
- [Render Documentation](https://render.com/docs)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Your API Documentation](./API_DOCUMENTATION.md)

---

## Support

**Neon Issues:**
- Dashboard: https://console.neon.tech
- Docs: https://neon.tech/docs
- Discord: https://discord.gg/neon

**Render Issues:**
- Dashboard: https://dashboard.render.com
- Docs: https://render.com/docs
- Status: https://status.render.com

**Your Backend Issues:**
- Check logs in Render dashboard
- Test database connection
- Verify environment variables
- Review API_DOCUMENTATION.md
