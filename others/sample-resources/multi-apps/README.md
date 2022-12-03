# 複数の成果物をまとめる構成

GitHub Actions を CI 基盤としつつ、Skaffold で複数の成果物を扱う環境を作ります。

## 1. GitHub リポジトリの用意

```bash
echo "# your-project" > README.md
git init
git add README.md
git commit -m "first commit"
git branch -M main
git remote add origin git@github.com:your-org/your-project.git
git push -u origin main
```

## 2. このディレクトリにある設定ファイル・サンプルコードをダウンロード

```bash
git clone https://github.com/google-cloud-japan/appdev-cicd-handson.git
cp -r appdev-cicd-handson/others/sample-resources/multi-apps/. ./
git checkout README.md
rm -rf appdev-cicd-handson
```

## 3. Google Cloud に実行環境を作成

利用する機能を有効化します。

```bash
gcloud services enable cloudresourcemanager.googleapis.com compute.googleapis.com \
    cloudbuild.googleapis.com artifactregistry.googleapis.com
```

コンテナのリポジトリを Artifact Registry に作ります。

```bash
gcloud artifacts repositories create my-apps \
    --repository-format=docker --location=asia-northeast1 \
    --description="Docker repository for CI/CD hands-on"
```

GitHub に渡すサービスアカウントと、鍵を生成します。

```bash
gcloud iam service-accounts create sa-github
PROJECT_ID=$(gcloud config get-value project)
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:sa-github@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/storage.admin"
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:sa-github@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/artifactregistry.writer"
PROJECT_NUMBER="$( gcloud projects list --filter="${PROJECT_ID}" \
    --format='value(PROJECT_NUMBER)' )"
gcloud iam service-accounts add-iam-policy-binding \
    ${PROJECT_NUMBER}-compute@developer.gserviceaccount.com \
    --member="serviceAccount:sa-github@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountUser"
gcloud iam service-accounts keys create credential.json \
    --iam-account=sa-github@${PROJECT_ID}.iam.gserviceaccount.com
cat credential.json
```

## 4. GitHub Actions の Secrets に鍵などを登録

GitHub から Google Cloud 上のリソースにアクセスするための変数をセットします。

- GOOGLECLOUD_PROJECT_ID: プロジェクト ID
- GOOGLECLOUD_SA_KEY: 4 の最後に出力された JSON 鍵

## 5. GitHub への push（パイプラインの起動）

```bash
git add --all
git commit -m "add ci/cd templates"
git push origin main
```

### GitHub Actions の設定

2 でダウンロードした資材の中には GitHub Actions のワークフロー定義が書かれています。  
.github/workflows 以下のファイルのトリガーと実行内容は以下の通りです。

- [pr-tests.yaml](https://github.com/google-cloud-japan/appdev-cicd-handson/blob/main/others/sample-resources/multi-apps/.github/workflows/pr-tests.yaml): PR 作成時にテストとビルドが実行されます
- [build-push.yaml](https://github.com/google-cloud-japan/appdev-cicd-handson/blob/main/others/sample-resources/multi-apps/.github/workflows/build-push.yaml): main ブランチの変更により、テスト、コンテナイメージのビルド、結果がコンテナレジストリに push されます

## 6. クリーンアップ

```bash
gcloud artifacts repositories delete my-apps --location=asia-northeast1 --quiet
```
