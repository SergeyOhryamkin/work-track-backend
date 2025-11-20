# GitHub Security Checklist

Before pushing to GitHub, verify these security measures are in place:

## ✅ Files That Are Safe to Commit

These files are already in your repository and are safe to make public:

- ✅ `.env.example` - Template only, no real credentials
- ✅ `.gitignore` - Protects sensitive files
- ✅ `go.mod`, `go.sum` - Dependency files
- ✅ All `.go` source files - No secrets hardcoded
- ✅ `Dockerfile`, `docker-compose.yml` - Configuration templates
- ✅ `Makefile` - Build commands
- ✅ `README.md`, `API_DOCUMENTATION.md`, `DEPLOYMENT.md` - Documentation
- ✅ `migrations/*.sql` - Database schema (no data)
- ✅ `.github/workflows/ci-cd.yml` - CI/CD configuration
- ✅ `render.yaml` - Deployment template
- ✅ `Procfile` - Process configuration

## 🚫 Files That Should NEVER Be Committed

These are protected by `.gitignore`:

- 🚫 `.env` - Contains real credentials (DB passwords, JWT secret)
- 🚫 `bin/` - Compiled binaries
- 🚫 `*.log` - Log files may contain sensitive data
- 🚫 Any files with passwords, API keys, or tokens
- 🚫 Database backup files (`.sql`, `.db`)
- 🚫 SSL certificates (`.pem`, `.key`, `.crt`)

## 🔍 Pre-Commit Security Check

Run these commands before your first commit:

```bash
# 1. Verify .env is NOT tracked
git status | grep -q ".env$" && echo "⚠️  WARNING: .env is tracked!" || echo "✅ .env is ignored"

# 2. Check for accidentally committed secrets
grep -r "password.*=" --include="*.go" --include="*.yml" --include="*.yaml" . || echo "✅ No hardcoded passwords found"

# 3. Verify .gitignore is working
git check-ignore .env && echo "✅ .env is properly ignored" || echo "⚠️  WARNING: .env might be tracked!"

# 4. List what will be committed
git status
```

## 📋 What Your Repository Will Contain

After following the deployment guide, your public GitHub repo will have:

```
work-track-backend/
├── .github/
│   └── workflows/
│       └── ci-cd.yml          # ✅ Safe (uses GitHub Secrets)
├── cmd/
│   └── api/
│       └── main.go            # ✅ Safe (no secrets)
├── internal/
│   ├── config/                # ✅ Safe (reads from env vars)
│   ├── database/              # ✅ Safe
│   ├── handler/               # ✅ Safe
│   ├── middleware/            # ✅ Safe
│   ├── models/                # ✅ Safe
│   ├── repository/            # ✅ Safe
│   ├── service/               # ✅ Safe
│   └── util/                  # ✅ Safe
├── migrations/
│   ├── *.up.sql               # ✅ Safe (schema only)
│   └── *.down.sql             # ✅ Safe
├── .env.example               # ✅ Safe (template only)
├── .gitignore                 # ✅ Safe (protects secrets)
├── API_DOCUMENTATION.md       # ✅ Safe
├── DEPLOYMENT.md              # ✅ Safe
├── Dockerfile                 # ✅ Safe
├── docker-compose.yml         # ✅ Safe (uses env vars)
├── go.mod                     # ✅ Safe
├── go.sum                     # ✅ Safe
├── Makefile                   # ✅ Safe
├── Procfile                   # ✅ Safe
├── QUICKSTART.md              # ✅ Safe
├── README.md                  # ✅ Safe
└── render.yaml                # ✅ Safe (uses env vars)
```

## 🔐 Where Secrets Are Stored (NOT in Git)

Your sensitive data will be stored in:

1. **Local Development**: `.env` file (gitignored)
2. **Render.com**: Environment Variables (in dashboard)
3. **GitHub Actions**: Repository Secrets (encrypted)
4. **Neon.tech**: Database credentials (in dashboard)

## ⚠️ Common Mistakes to Avoid

### ❌ DON'T DO THIS:
```go
// Bad: Hardcoded credentials
const dbPassword = "mypassword123"
const jwtSecret = "supersecret"
```

### ✅ DO THIS INSTEAD:
```go
// Good: Read from environment
password := os.Getenv("DB_PASSWORD")
jwtSecret := os.Getenv("JWT_SECRET")
```

## 🛡️ Additional Security Best Practices

### 1. Use Strong Secrets
```bash
# Generate strong JWT secret
openssl rand -base64 32

# Generate strong password
openssl rand -base64 24
```

### 2. Never Log Sensitive Data
```go
// Bad
log.Printf("User password: %s", password)

// Good
log.Printf("User authenticated: %s", username)
```

### 3. Use GitHub Secrets for CI/CD
- Never put credentials in `.github/workflows/ci-cd.yml`
- Use `${{ secrets.SECRET_NAME }}` instead

### 4. Review Before Pushing
```bash
# See what will be committed
git diff --cached

# See all files that will be added
git status
```

## 📝 First-Time Git Setup Checklist

- [ ] `.gitignore` is in place
- [ ] `.env` file exists locally (for development)
- [ ] `.env` is listed in `.gitignore`
- [ ] No hardcoded secrets in code
- [ ] Reviewed `git status` output
- [ ] Verified `.env` is NOT in the list
- [ ] Ready to commit!

## 🚀 Safe to Proceed

If all checks pass, you're ready to:

```bash
# Initialize git
git init

# Add all files (gitignore will protect .env)
git add .

# Verify .env is NOT staged
git status | grep ".env$" && echo "STOP! .env is staged!" || echo "Safe to commit"

# Create first commit
git commit -m "Initial commit: Go backend for work tracking"

# Push to GitHub (after creating repo)
git remote add origin https://github.com/YOUR_USERNAME/work-track-backend.git
git push -u origin main
```

## 🆘 If You Accidentally Commit Secrets

If you accidentally commit `.env` or secrets:

```bash
# Remove from git but keep local file
git rm --cached .env

# Commit the removal
git commit -m "Remove .env from git"

# Force push (if already pushed to GitHub)
git push -f origin main

# IMPORTANT: Change all exposed secrets immediately!
# - Generate new JWT_SECRET
# - Change database password
# - Rotate any exposed API keys
```

## ✅ Final Verification

Before pushing to GitHub:

```bash
# 1. Check what's being tracked
git ls-files | grep -E "\\.env$|password|secret" && echo "⚠️  Check these files!" || echo "✅ All clear"

# 2. Verify .gitignore is working
git check-ignore .env bin/ *.log && echo "✅ Gitignore working" || echo "⚠️  Check .gitignore"

# 3. Review commit
git log --oneline -1
git show --stat
```

---

## 📚 Resources

- [GitHub Security Best Practices](https://docs.github.com/en/code-security)
- [Git Secrets Tool](https://github.com/awslabs/git-secrets)
- [.gitignore Generator](https://www.toptal.com/developers/gitignore)

---

**Remember**: Once something is pushed to GitHub, assume it's public forever, even if you delete it later!
