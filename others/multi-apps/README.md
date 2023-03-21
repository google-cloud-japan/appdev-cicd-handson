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

ローカルで Kubernetes クラスタを起動して

```bash
minikube start
```

サンプル アプリケーションを実行してみます。

```bash
skaffold dev -p local --port-forward
```

データベースにユーザーデータを入れた上で

```bash
mysql -h 127.0.0.1 -P 3306 -u user -pPassw0rd myapp
mysql> CREATE TABLE `users` (`code` varchar(32), `name` varchar(100));
mysql> INSERT INTO `users` VALUES ('001', 'Sato Taro');
```

http://localhost:8080 で挙動を確認できます。


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

GitHub に渡すサービスアカウントを生成します。

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
```

GitHub に安全に権限を渡すため、[Workload Identity 連携](https://cloud.google.com/iam/docs/workload-identity-federation?hl=ja) を設定します。

```bash
gcloud iam workload-identity-pools create "idpool-cicd" --location "global" --display-name "Identity pool for CI/CD services"
idp_id=$( gcloud iam workload-identity-pools describe "idpool-cicd" --location "global" --format "value(name)" )
```

Identity Provider (IdP) を作成します。GitHub リポジトリを一意に識別するための ID を設定し、

```bash
repo=<org-id>/<repo-id>
```

Identity Provider (IdP) を作成します。

```bash
gcloud iam workload-identity-pools providers create-oidc "idp-github" --workload-identity-pool "idpool-cicd" --location "global" --issuer-uri "https://token.actions.githubusercontent.com" --attribute-mapping "google.subject=assertion.sub,attribute.repository=assertion.repository" --display-name "Workload IdP for GitHub"
gcloud iam service-accounts add-iam-policy-binding sa-github@${PROJECT_ID}.iam.gserviceaccount.com --member "principalSet://iam.googleapis.com/${idp_id}/attribute.repository/${repo}" --role "roles/iam.workloadIdentityUser"
gcloud iam workload-identity-pools providers describe "idp-github" --workload-identity-pool "idpool-cicd" --location "global" --format "value(name)"
```

## 4. GitHub Actions の Secrets に鍵などを登録

GitHub から Google Cloud 上のリソースにアクセスするための変数をセットします。

- **GOOGLE_CLOUD_PROJECT**: プロジェクト ID
- **GOOGLE_CLOUD_WORKLOAD_IDP**: 2.1 の最後に出力された IdP ID

## 5. GitHub への push（パイプラインの起動）

```bash
git add --all
git commit -m "add ci/cd templates"
git push origin main
```

### GitHub Actions の設定

2 でダウンロードした資材の中には GitHub Actions のワークフロー定義が書かれています。  
.github/workflows 以下のファイルのトリガーと実行内容は以下の通りです。

- [build-push.yaml](https://github.com/google-cloud-japan/appdev-cicd-handson/blob/main/others/sample-resources/multi-apps/.github/workflows/build-push.yaml): main ブランチの変更により、テスト、コンテナイメージのビルド、結果がコンテナレジストリに push されます

### 成果物の確認

[Artifact Registry コンソール](https://console.cloud.google.com/artifacts)を開いてみましょう。**my-apps/api** と **my-apps/web** というリポジトリに git ハッシュのタグでイメージが確認できます。

## 6. クリーンアップ

```bash
gcloud artifacts repositories delete my-apps --location=asia-northeast1 --quiet
```
