# Cloud Deploy 最小構成

GitHub Actions を CI、Cloud Deploy を CD 基盤とするシンプルな CI/CD 環境を作ります。

## 1. GitHub リポジトリを用意します

```bash
echo "# your-project" > README.md
git init
git add README.md
git commit -m "first commit"
git branch -M main
git remote add origin git@github.com:your-org/your-project.git
git push -u origin main
```

## 2. このディレクトリにある設定ファイルをダウンロード

```bash
git clone https://github.com/google-cloud-japan/appdev-cicd-handson.git
cp -r appdev-cicd-handson/cloud-deploy/sample-resources/minimum/. ./
echo "deploy/clouddeploy.yaml" > .gitignore
git checkout README.md
rm -rf appdev-cicd-handson
```

## 3. サンプルのソースコードをダウンロードします

```bash
git clone https://github.com/dart-lang/samples.git
cp -r samples/server/simple/. src
rm -rf samples
```

## 4. Google Cloud に実行環境を作成します

利用する機能を有効化します。

```bash
gcloud services enable cloudresourcemanager.googleapis.com compute.googleapis.com \
    container.googleapis.com serviceusage.googleapis.com stackdriver.googleapis.com \
    monitoring.googleapis.com logging.googleapis.com clouddeploy.googleapis.com \
    cloudbuild.googleapis.com artifactregistry.googleapis.com
```

コンテナのリポジトリを Artifact Registry に作り

```bash
gcloud artifacts repositories create cd-test \
    --repository-format=docker --location=asia-northeast1 \
    --description="Docker repository for CI/CD hands-on"
```

実行環境として GKE クラスタを 1 つ作成します。

```bash
gcloud container clusters create cd-test --zone asia-northeast1-a \
    --release-channel stable --machine-type "e2-standard-4" \
    --num-nodes 1 --preemptible
```

GitHub に渡すサービスアカウントと、鍵を生成します。

```bash
gcloud iam service-accounts create sa-cd-test
PROJECT_ID=$(gcloud config get-value project)
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:sa-cd-test@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/storage.admin"
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:sa-cd-test@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/artifactregistry.writer"
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:sa-cd-test@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/clouddeploy.releaser"
PROJECT_NUMBER="$( gcloud projects list --filter="${PROJECT_ID}" \
    --format='value(PROJECT_NUMBER)' )"
gcloud iam service-accounts add-iam-policy-binding \
    ${PROJECT_NUMBER}-compute@developer.gserviceaccount.com \
    --member="serviceAccount:sa-cd-test@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountUser"
gcloud iam service-accounts keys create credential.json \
    --iam-account=sa-cd-test@${PROJECT_ID}.iam.gserviceaccount.com
cat credential.json
```

## 5. GitHub Actions の Secrets に鍵などを登録します

GitHub から Google Cloud 上のリソースにアクセスするための変数をセットします。

- GOOGLECLOUD_PROJECT_ID: プロジェクト ID
- GOOGLECLOUD_SA_KEY: 4 の最後に出力された JSON 鍵

## 6. Cloud Deploy のパイプラインを作ります

deploy/clouddeploy.yaml を開いて <your-project-id> を 2 ヶ所適切なものに変更した上で、パイプラインを作成します。

```bash
vim deploy/clouddeploy.yaml
gcloud beta deploy apply --file deploy/clouddeploy.yaml --region us-central1
```

## 7. GitHub へ push します

```bash
git add --all
git commit -m "add ci/cd templates"
git push origin main
```

### GitHub Actions の設定

2 でダウンロードした資材の中には GitHub Actions のワークフロー定義が書かれています。  
.github/workflows 以下のファイルのトリガーと実行内容は以下の通りです。

- pr-tests.yaml: PR 作成時にテストとビルドが実行されます
- release.yaml: main ブランチの変更により、Cloud Deploy にリリースが作成されます
- promotion.yaml: prod- で始まるタグを打つと、そのコミットで作られたリリースがプロモーションされます

### Cloud Deploy の設定

2 でダウンロードした資材の中にアプリケーションのビルドやデプロイに関する定義があります。  
各ファイルでの定義内容は以下の通りです。

- skaffold.yaml: ビルド対象は src 以下、デプロイは k8s マニフェストで実施することを定義
- deploy/k8s/dev/web.yaml: prod プロファイルを指定しない限りはこちらがデプロイされる
- deploy/k8s/prod/web.yaml: prod プロファイル指定時にはこちらがデプロイされる

## 8. Cloud Deploy の dev 環境の様子を確認

GitHub Actions の状況や  
Cloud Deploy パイプラインの状況、  
https://console.cloud.google.com/deploy/delivery-pipelines

GKE のワークロードの変化を確認してみてください。  
https://console.cloud.google.com/kubernetes/workload

## 9. プロモーション

画面からもできますが、ここでは git のタグ打ちでプロモーションする様子をみてみます。  
（このフローはリードタイムが伸びるので、望ましいかどうかはよく検討すべきですが・・）

```bash
git tag prod-1.0
git push origin prod-1.0
```

## 10. クリーンアップ

```bash
gcloud beta deploy delivery-pipelines delete minimum-pipeline --force \
    --region us-central1
gcloud container clusters delete cd-test --zone asia-northeast1-a
```
