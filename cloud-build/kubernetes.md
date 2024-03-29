# Kubernetes で実践する Google Cloud での CI / CD ハンズオン

## 始めましょう

Cloud Shell をベースにローカル開発、Google Cloud での CI / CD を体験いただくハンズオンです。以下の流れで実際のアプリケーション開発を体験いただきます。

1. ローカルでの開発
1. Kubernetes をベースにした CI / CD
1. 高度なデプロイオプションの利用

<walkthrough-tutorial-duration duration="60"></walkthrough-tutorial-duration>
<walkthrough-tutorial-difficulty difficulty="2"></walkthrough-tutorial-difficulty>

**前提条件**:

- Google Cloud 上にプロジェクトが作成してある
- プロジェクトの *編集者* 相当の権限をもつユーザーでログインしている
- *プロジェクト IAM 管理者* 相当の権限をもつユーザーでログインしている
- （推奨）Google Chrome を利用している

**[開始]** ボタンをクリックして次のステップに進みます。

## プロジェクトの設定

この手順の中で実際にリソースを構築する対象のプロジェクトを選択してください。

<walkthrough-project-setup></walkthrough-project-setup>

## CLI の初期設定と権限の確認

gcloud（[Google Cloud の CLI ツール](https://cloud.google.com/sdk/gcloud?hl=ja)）のデフォルト プロジェクトを設定します。

```bash
export PROJECT_ID=<walkthrough-project-id/>
```

```bash
gcloud config set project "${PROJECT_ID}"
```

[Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine?hl=ja) と [Cloud Source Repositories](https://cloud.google.com/source-repositories?hl=ja) を扱える権限があることを担保します。

```bash
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member "user:$(gcloud config get-value core/account)" --role roles/container.admin
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member "user:$(gcloud config get-value core/account)" --role roles/source.admin
```

## 1. ローカルでの開発

[Cloud Shell エディタ](https://cloud.google.com/shell/docs/launching-cloud-shell-editor?hl=ja) は個人ごとに割り当てられる開発環境としてご利用いただけます。このセクションでは、以下のような流れで開発者の個人環境での開発を実施します。

1. サンプルコードのダウンロード
1. アプリの起動まで
1. コードの変更、リアルタイムでの更新
1. デバッグ
1. ログの確認

## 1.1. サンプルコードのダウンロード

Java のフレームワークとして [Micronaut](https://micronaut.io) を、コンテナのビルド フレームワークとして [Jib](https://cloud.google.com/blog/ja/products/application-development/using-jib-to-containerize-java-apps) を利用します。サンプルコードをダウンロードしましょう。

1.  コードをダウンロードし

    ```bash
    git clone https://github.com/GoogleContainerTools/jib.git
    rm -rf ~/jib/.git
    ```

1.  サンプルコードのディレクトリを *ワークスペース* として Cloud Shell エディタを起動します。

    ```bash
    cloudshell workspace jib/examples/micronaut
    ```

1.  Hello World を返すコントローラは
    <walkthrough-editor-open-file filePath="jib/examples/micronaut/src/main/groovy/example/micronaut/HelloController.groovy">こちら</walkthrough-editor-open-file>です。

1.  Spock による
    <walkthrough-editor-open-file filePath="jib/examples/micronaut/src/test/groovy/example/micronaut/HelloControllerSpec.groovy">テストコード</walkthrough-editor-open-file>
    があるので
1.  新しい <walkthrough-editor-spotlight spotlightId="menu-terminal-new-terminal">ターミナル
    </walkthrough-editor-spotlight> を開き、テストを実行してみましょう。

    ```bash
    ./gradlew test
    ```

## 1.2. アプリの起動まで

1.  Minikube を起動しましょう。

    ```bash
    minikube start
    ```

1.  結果が以下のように表示されたら

    ```terminal
    🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
    ```

1.  Kubernetes のローカル開発支援ソフトウェア、[Skaffold](https://skaffold.dev/) の設定ファイルを作り

    ```text
    cat << EOF > skaffold.yaml
    apiVersion: skaffold/v3
    kind: Config
    build:
      artifacts:
      - image: app
        jib:
          type: gradle
    manifests:
      kustomize:
        paths: ["k8s/base"]
    deploy:
      kubectl: {}
    profiles:
    - name: local
      patches:
      - op: add
        path: /build/artifacts/0/jib/fromImage
        value: gcr.io/distroless/java:debug
    EOF
    ```

1.  Kubernetes のマニフェストを [kustomize](https://kustomize.io/) ベースで作ります。

    ```text
    mkdir -p k8s/base
    cat << EOF > k8s/base/kustomization.yaml
    resources:
    - web.yaml
    EOF
    cat << EOF > k8s/base/web.yaml
    kind: Deployment
    apiVersion: apps/v1
    metadata:
      name: web-app
    spec:
      selector:
        matchLabels:
          app: web
      template:
        metadata:
          labels:
            app: web
        spec:
          containers:
          - name: main
            image: app
            ports:
            - containerPort: 8080
    ---
    kind: Service
    apiVersion: v1
    metadata:
      name: web-svc
    spec:
      selector:
        app: web
      ports:
      - port: 8080
        name: http
    EOF
    ```

1.  Dry Run で設定内容を確認してみます。

    ```bash
    kubectl apply --dry-run=client --kustomize k8s/base
    ```

1.  Skaffold を使ってアプリケーションを起動してみましょう。

    ```bash
    skaffold dev -p local --auto-build --auto-deploy --auto-sync --port-forward
    ```

1.  Web preview ボタンを押し、"ポート 8080 でプレビュー" を選んでみましょう。
    サンプルアプリは `/hello` で実装されているので、URL に `/hello` を追加しリロードします。
    <walkthrough-web-preview-icon/>

Hello World はうまく返ってきましたか？

## 1.3. コードの変更、リアルタイム更新

コードの書き換えによって、アプリケーションがリアルタイムに更新されることを確認します。

1.  <walkthrough-editor-open-file filePath="jib/examples/micronaut/src/main/groovy/example/micronaut/HelloController.groovy">HelloController.groovy</walkthrough-editor-open-file> を開き

1.  <walkthrough-editor-select-line filePath="jib/examples/micronaut/src/main/groovy/example/micronaut/HelloController.groovy" startLine="15" endLine="15" startCharacterOffset="9" endCharacterOffset="20">Hello World</walkthrough-editor-select-line> を変更してみましょう。

1.  ログが進み、

    ```terminal
    [main] 04:19:01.557 [main] INFO  io.micronaut.runtime.Micronaut - Startup completed in 1999ms. Server Running: http://web-app-xxxxxxxxx-yyyyy:8080
    ```

    のようなメッセージが表示されます。

1.  Web プレビュー画面をリロードしてみましょう。

変更は反映されましたか？

<walkthrough-footnote>ここまでで、開発者それぞれに与えられた環境での開発フローを見てきました。ここからは、チームとして製品を開発、CI / CD を回す方法を確認していきましょう。</walkthrough-footnote>

## 2. Kubernetes をベースにした CI/CD

ここからは品質向上のため、そしてチームとして開発する上で重要になる CI/CD を織り込む方法をみていきます。

1. git リポジトリの準備
1. CI によるテストの自動化
1. 開発環境への CD
1. ログの確認
1. コンテナでターミナルを開く

## 2.1. git リポジトリの準備

アプリケーション コードを置く git リポジトリとして [Cloud Source Repositories (CSR)](https://cloud.google.com/source-repositories?hl=ja) を利用します。リポジトリを作成し、Cloud Shell からアクセスするための設定を進めます。

1.  <walkthrough-editor-spotlight spotlightId="menu-terminal-new-terminal">ターミナル
    </walkthrough-editor-spotlight> を開き、改めてプロジェクト ID を指定します。

    ```bash
    export PROJECT_ID=<walkthrough-project-id/>
    ```

    API を有効化、git リポジトリを CSR に作成します。

    ```bash
    gcloud services enable sourcerepo.googleapis.com cloudbuild.googleapis.com artifactregistry.googleapis.com compute.googleapis.com container.googleapis.com
    gcloud source repos create cicd-gke
    ```

1.  CSR への認証ヘルパ含め、git クライアントの設定をします。

    ```bash
    git config --global credential.helper gcloud.sh
    git config --global user.name "$(whoami)"
    git config --global user.email "$(gcloud config get-value core/account)"
    ```

1.  ignore ファイルを作りつつ

    ```text
    cat << EOF > .gitignore
    .gradle
    **/build/
    !src/**/build/
    !gradle-wrapper.jar
    .gradletasknamecache
    EOF
    ```

1.  コードを git 管理下におき、CSR へ push しましょう。

    ```bash
    git init
    git remote add google "https://source.developers.google.com/p/${PROJECT_ID}/r/cicd-gke"
    git checkout -b main
    git add .
    git commit -m 'init'
    git push google main
    ```

## 2.2. CI によるテストの自動化

git push と同時にテスト実行 + ビルドするステップを自動化してみましょう。

1.  コンテナ レジストリを作ります。

    ```bash
    gcloud artifacts repositories create cicd-gke --repository-format=docker --location=asia-northeast1 --description="Docker repository for CI/CD hands-on"
    gcloud auth configure-docker asia-northeast1-docker.pkg.dev
    docker pull alpine:3.14
    docker tag alpine:3.14 asia-northeast1-docker.pkg.dev/${PROJECT_ID}/cicd-gke/app:init
    docker push asia-northeast1-docker.pkg.dev/${PROJECT_ID}/cicd-gke/app:init
    ```

1.  Cloud Build に対して必要な権限を付与します。

    ```bash
    project_number="$( gcloud projects list --filter="${PROJECT_ID}" --format='value(PROJECT_NUMBER)' )"
    gcloud projects add-iam-policy-binding ${PROJECT_ID} --member "serviceAccount:${project_number}@cloudbuild.gserviceaccount.com" --role roles/container.admin
    gcloud iam service-accounts add-iam-policy-binding ${project_number}-compute@developer.gserviceaccount.com --member="serviceAccount:${project_number}@cloudbuild.gserviceaccount.com" --role="roles/iam.serviceAccountUser"
    ```

1.  Cloud Build の設定ファイル、`cloudbuild-ci.yaml` を作ります。

    ```text
    cat << EOF > cloudbuild-ci.yaml
    steps:
    - id: Test
      name: gcr.io/cloud-builders/gradle
      entrypoint: gradle
      args: ['test']
    tags: ['test']
    EOF
    ```

1.  git push により CI が起動するようトリガーを設定します。

    ```bash
    gcloud builds triggers create cloud-source-repositories --name cicd-gke-ci --repo=cicd-gke --branch-pattern='.*' --build-config=cloudbuild-ci.yaml
    ```

1.  Cloud Build コンソールを開きましょう。
    <walkthrough-menu-navigation sectionId="CLOUD_BUILD_SECTION"></walkthrough-menu-navigation>

1.  git push によりビルドが始まることを確認します。

    ```bash
    git add cloudbuild-ci.yaml
    git commit -m 'Add continuous integration'
    git push google main
    ```

これによりテストが始まります。テストは `Hello World` という応答を期待している一方、先程コントローラを変更したままだとテストは赤くなります。その場合 `Hello World` にもどしておきましょう。

## 2.3. 開発環境への CD

開発環境を作り、そこへの継続的デリバリーパイプラインを作成します。開発環境に対しては **最新の main ブランチがリリースされ続ける**ように設定してみます。

1.  GKE クラスタを作成しましょう。

    ```bash
    gcloud container clusters create "cicd-gke-dev" --zone "asia-northeast1-a" --machine-type "e2-standard-2" --num-nodes=1 --release-channel "stable" --enable-ip-alias --logging "SYSTEM,API_SERVER,WORKLOAD" --workload-pool "${PROJECT_ID}.svc.id.goog" --scopes "cloud-platform" --async
    ```

1.  Skaffold の設定ファイルに開発環境への設定を加えます。

    ```text
    cat << EOF >> skaffold.yaml
    - name: dev
      patches:
      - op: add
        path: /build/tagPolicy
        value:
          gitCommit:
            ignoreChanges: true
      manifests:
        kustomize:
          paths: ["k8s/overlays/dev"]
    EOF
    ```

1.  Kubernetes のマニフェストも base からの差分として定義します。

    ```text
    mkdir -p k8s/overlays/dev
    cat << EOF > k8s/overlays/dev/kustomization.yaml
    bases:
    - ../../base
    patchesStrategicMerge:
    - web.yaml
    EOF
    cat << EOF > k8s/overlays/dev/web.yaml
    kind: Deployment
    apiVersion: apps/v1
    metadata:
      name: web-app
    spec:
      strategy:
        type: RollingUpdate
        rollingUpdate:
          maxUnavailable: 0
          maxSurge: 1
      replicas: 1
      template:
        spec:
          containers:
          - name: main
            env:
            - name: LOG_LEVEL
              value: "info"
    EOF
    ```

1.  Dry Run で設定内容を確認してみます。

    ```bash
    kubectl apply --dry-run=client --kustomize k8s/overlays/dev
    ```

1.  Cloud Build の設定ファイル、`cloudbuild-cd-dev.yaml` を作ります。

    ```text
    cat << EOF > cloudbuild-cd-dev.yaml
    steps:
    - id: Build & Push
      name: gcr.io/cloud-builders/gcloud
      entrypoint: skaffold
      args:
      - build
      - -p
      - dev
      - --default-repo
      - 'asia-northeast1-docker.pkg.dev/${PROJECT_ID}/cicd-gke'
      - --push
      - --file-output=/workspace/build.out
    - id: Render
      name: gcr.io/cloud-builders/gcloud
      entrypoint: skaffold
      args:
      - render
      - -p
      - dev
      - --offline=true
      - --build-artifacts=/workspace/build.out
      - --output=/workspace/resources.yaml
    - id: Deploy
      name: gcr.io/cloud-builders/gke-deploy
      args:
      - run
      - --cluster=cicd-gke-dev
      - --location=asia-northeast1-a
      - --filename=/workspace/resources.yaml
    tags: ['dev']
    EOF
    ```

1.  **main ブランチへの** git push により CD が起動するようトリガーを設定します。

    ```bash
    gcloud builds triggers create cloud-source-repositories --name cicd-gke-cd-dev --repo=cicd-gke --branch-pattern='^main$' --build-config=cloudbuild-cd-dev.yaml
    ```

1.  **main ブランチへの** git push によりデプロイが始まることを確認します。

    ```bash
    git add k8s/overlays/dev skaffold.yaml cloudbuild-cd-dev.yaml src/main/groovy/example/micronaut/HelloController.groovy
    git commit -m 'Add continuous delivery'
    git push google main
    ```

## 2.4. ログの確認

デプロイした開発環境のログを Cloud Shell から確認してみましょう。

1.  クラスタへの接続情報を取得します。

    ```bash
    gcloud container clusters get-credentials "cicd-gke-dev" --zone asia-northeast1-a
    ```

1.  **Ctrl**/**Cmd**+**Shift**+**P** でコマンドパレットを開き、
    **Cloud Code: View Logs** とタイプし、Log Viewer を起動します。

1.  Namespace として `default` を選び、<walkthrough-editor-spotlight spotlightId="cloud-code-logs-viewer-deployment">Deployment</walkthrough-editor-spotlight>
    または
    <walkthrough-editor-spotlight spotlightId="cloud-code-logs-viewer-pod">Pod</walkthrough-editor-spotlight>
    でフィルタリング、目的のログを表示します。

1.  ログは `Streaming` を on にするか、ブラウザを更新するか、
    <walkthrough-editor-spotlight spotlightId="cloud-code-logs-viewer-refresh">更新ボタン</walkthrough-editor-spotlight> で新しいログが確認できます。

## 2.5. コンテナでターミナルを開く

Cloud Code の Kubernetes Explorer では様々な情報が確認できます。ここでは、起動しているコンテナに接続してみましょう。

1.  左側のメニュー
    <walkthrough-editor-spotlight spotlightId="cloud-code-k8s-icon">Kubernetes
    Explorer</walkthrough-editor-spotlight> を開きます。
1.  "cicd-gke-dev" クラスタを選び、*Namespaces > default > Pods* から `web-app` で始まる Pod を探し、
    右クリック、*'Get Terminal'* を選択します。
1.  ps コマンドで、PID 1 で Java プロセスが起動していることを確認してみましょう。

    ```bash
    ps uxw
    ```

## 3. クリーンアップ

ハンズオンに利用したプロジェクトを削除し、課金を止めます。

```bash
gcloud projects delete ${PROJECT_ID}
```

プロジェクトがそのまま消せない場合は、以下のリソースを個別に削除してください。

- GKE クラスタ
- Cloud Build のトリガー
- Cloud Source Repositories の git リポジトリ
- Artifact Registry の コンテナ リポジトリ

## これで終わりです

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

すべて完了しました。
