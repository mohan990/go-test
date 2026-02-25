# Cloud Run Deployment Guide

## Answer to Your Question

**Does Cloud Run use the GitHub workflow image or my Dockerfile?**

**Answer:** Cloud Run uses a Docker image that was **built using YOUR Dockerfile**. The image can come from either:
1. Manual build: `gcloud builds submit` (uses your Dockerfile)
2. GitHub Actions: Automated workflow (also uses your Dockerfile)

In both cases, **YOUR Dockerfile is used** to build the image.

---

## Current Setup

### Existing GitHub Workflow (`.github/workflows/ci.yml`)
- ✅ Tests your code
- ✅ Builds Docker image (using your Dockerfile)
- ❌ Does NOT push to registry
- ❌ Does NOT deploy to Cloud Run

### New Deployment Workflow (`.github/workflows/deploy-cloudrun.yml`)
- ✅ Builds Docker image using **YOUR Dockerfile**
- ✅ Pushes to Google Container Registry (GCR)
- ✅ Deploys to Cloud Run automatically
- ✅ Triggers on push to `main` branch

---

## Setup Instructions

### Step 1: Create a GCP Service Account

```bash
# Set your project ID
export PROJECT_ID="your-project-id"

# Create service account
gcloud iam service-accounts create github-actions \
  --display-name "GitHub Actions Service Account" \
  --project $PROJECT_ID

# Grant necessary permissions
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/run.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/storage.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"

# Create and download key
gcloud iam service-accounts keys create key.json \
  --iam-account=github-actions@$PROJECT_ID.iam.gserviceaccount.com
```

### Step 2: Add GitHub Secrets

Go to your GitHub repository → Settings → Secrets and variables → Actions

Add these secrets:

1. **GCP_PROJECT_ID**
   - Value: `your-gcp-project-id`

2. **GCP_SA_KEY**
   - Value: Entire contents of `key.json` file (copy/paste all)

### Step 3: Update Workflow Configuration

Edit `.github/workflows/deploy-cloudrun.yml`:

```yaml
env:
  PROJECT_ID: ${{ secrets.GCP_PROJECT_ID }}
  REGION: us-central1  # Change to your region (e.g., us-east1, europe-west1)
  SERVICE_NAME: go-webserver  # Change if you want a different service name
```

### Step 4: Enable Required APIs

```bash
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com
gcloud services enable cloudbuild.googleapis.com
```

### Step 5: Deploy

#### Option A: Automatic Deployment (via GitHub)
```bash
git add .
git commit -m "Add Cloud Run deployment workflow"
git push origin main
```

The workflow will automatically:
1. Build Docker image using YOUR Dockerfile
2. Push to GCR
3. Deploy to Cloud Run

#### Option B: Manual Deployment
```bash
# Build and deploy in one command
gcloud run deploy go-webserver \
  --source . \
  --region us-central1 \
  --allow-unauthenticated
```

This also uses YOUR Dockerfile!

---

## How It Works

### Build Process Flow

```
Your Dockerfile
      ↓
[Build Stage] - Compiles Go code in golang:1.22-alpine
      ↓
[Final Stage] - Creates minimal alpine image with binary
      ↓
Docker Image (~15-20MB)
      ↓
Google Container Registry (gcr.io/PROJECT_ID/SERVICE_NAME)
      ↓
Cloud Run Deployment
```

### Key Points

1. **Always uses YOUR Dockerfile** - Whether you deploy manually or via GitHub Actions
2. **Multi-stage build** - Reduces final image size from 300MB to 15-20MB
3. **Automatic builds** - GitHub Actions builds on every push to main
4. **Health checks** - Cloud Run uses `/healthz`, `/readyz`, `/livez` endpoints
5. **Security** - Runs as non-root user in container

---

## Verify Deployment

### Check Service Status
```bash
gcloud run services list --region us-central1
```

### Get Service URL
```bash
gcloud run services describe go-webserver \
  --region us-central1 \
  --format 'value(status.url)'
```

### Test Endpoints
```bash
# Get the URL
SERVICE_URL=$(gcloud run services describe go-webserver --region us-central1 --format 'value(status.url)')

# Test endpoints
curl $SERVICE_URL/
curl $SERVICE_URL/hello?name=YourName
curl $SERVICE_URL/healthz
curl $SERVICE_URL/readyz
curl $SERVICE_URL/livez
```

---

## Troubleshooting

### Deployment Failed

1. **Check logs:**
```bash
gcloud run services logs read go-webserver --region us-central1
```

2. **Verify image exists:**
```bash
gcloud container images list --repository gcr.io/$PROJECT_ID
```

3. **Check service account permissions:**
```bash
gcloud projects get-iam-policy $PROJECT_ID \
  --flatten="bindings[].members" \
  --filter="bindings.members:serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com"
```

### GitHub Actions Failed

1. Check the Actions tab in your GitHub repository
2. Verify secrets are set correctly
3. Check if APIs are enabled in GCP
4. Verify service account has correct permissions

---

## Cost Optimization

Cloud Run pricing (as of 2024):
- **Free tier:** 2 million requests/month
- **Pay per use:** Only charged when requests are being processed
- **No idle costs:** With this small image, cold starts are fast (~500ms)

### Recommended Settings for Low Traffic
```yaml
--memory 256Mi        # Minimum memory
--cpu 1               # 1 vCPU
--max-instances 10    # Prevent runaway costs
--min-instances 0     # Scale to zero when idle
```

---

## Next Steps

- [ ] Set up GitHub secrets
- [ ] Enable GCP APIs
- [ ] Test manual deployment first
- [ ] Set up automated deployment via GitHub Actions
- [ ] Configure custom domain (optional)
- [ ] Set up monitoring and alerting (optional)
