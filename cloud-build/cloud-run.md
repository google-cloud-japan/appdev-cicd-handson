# Cloud Run で実践する Google Cloud での CI / CD ハンズオン

<walkthrough-watcher-constant key="app" value="cicd-run"></walkthrough-watcher-constant>
<walkthrough-watcher-constant key="region" value="asia-northeast1"></walkthrough-watcher-constant>
<walkthrough-watcher-constant key="github" value="google-cloud-japan/gcp-getting-started-cloudrun/main"></walkthrough-watcher-constant>

## 始めましょう

Cloud Shell をベースにローカル開発、Google Cloud での CI / CD を体験いただくハンズオンです。以下の流れで実際のアプリケーション開発を体験いただきます。

1. ローカルでの開発
1. Cloud Run をベースにした CI / CD
1. 高度なデプロイオプションの利用

<walkthrough-tutorial-duration duration="60"/> 
**所要時間**: 約 60 分

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

gcloud（[Google Cloud の CLI ツール](https://cloud.google.com/sdk/gcloud?hl=ja)
のデフォルト プロジェクトを設定します。

```bash
export PROJECT_ID={{project-id}}
```

```bash
gcloud config set project "${PROJECT_ID}"
```

[Cloud Run](https://cloud.google.com/run?hl=ja) と [Cloud Source Repositories](https://cloud.google.com/source-repositories?hl=ja) を扱える権限があることを担保します。

```bash
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member "user:$(gcloud config get-value core/account)" --role roles/run.admin
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member "user:$(gcloud config get-value core/account)" --role roles/source.admin
```

## 1. ローカルでの開発

[Cloud Shell エディタ](https://cloud.google.com/shell/docs/launching-cloud-shell-editor?hl=ja) は個人ごとに割り当てられる開発環境としてご利用いただけます。このセクションでは、以下のような流れで開発者の個人環境での開発を実施します。

1. サンプルコードのダウンロード
1. アプリの起動まで
1. コードの変更、リアルタイムでの更新
1. デバッグ
1. 個人開発環境へのデプロイ
1. ログの確認

## 1.1. サンプルコードのダウンロード

まず以下のリンクから、Cloud Shell エディタを起動してみましょう。<walkthrough-editor-open-file filePath="README-cloudshell.txt">Cloud Shell エディタを開く</walkthrough-editor-open-file>

1.  ステータス バーから
    <walkthrough-editor-spotlight spotlightId="cloud-code-status-bar">Cloud Code</walkthrough-editor-spotlight> を開き、
1.  <walkthrough-editor-spotlight spotlightId="cloud-code-new-app">New 
    Application</walkthrough-editor-spotlight> > <walkthrough-editor-spotlight spotlightId="cloud-code-new-app-cloud-run">Cloud Run application</walkthrough-editor-spotlight> を選びます。
1.  Cloud Run のサンプル一覧から、**Go: Cloud Run** を選択、
1.  ダウンロード先として任意の場所を選び **Create New Application** をクリックします。
1.  Cloud Shell エディタがリロードされたら、
    <walkthrough-editor-spotlight spotlightId="file-explorer">explorer view</walkthrough-editor-spotlight> でファイルの一覧を確認しましょう。

## 1.2. アプリの起動まで

ではこのアプリケーションを Cloud Shell のローカルで、Cloud Run エミュレータを使い起動してみます。

1.  アプリケーションをビルドし、エミュレータで実行するには
    <walkthrough-editor-spotlight spotlightId="cloud-code-run-on-cloud-run-emulator">Run
    on Cloud Run Emulator</walkthrough-editor-spotlight> を選択します。
1.  **Run/Debug on Cloud Run Emulator** タブが開いたら **Build Settings** の **Bulder** を **Buildpacks** に変更し、**Run** をクリックします。
1.  初回の起動には 5 分ほどかかります。サービスがデプロイされると、
    <walkthrough-editor-spotlight spotlightId="output">Output</walkthrough-editor-spotlight>
    パネルに以下のように表示されます。

    ```terminal
    http://localhost:8080
    Update successful
    ```
1.  Web preview ボタン <walkthrough-web-preview-icon/> を押し、"ポート 8080 でプレビュー" を選んでみましょう。

おめでとうございます！アプリの起動はうまくいきましたね。

## 1.3. コードの変更、リアルタイム更新

コードの書き換えによって、アプリケーションがリアルタイムに更新されることを確認します。

1.  <walkthrough-editor-open-file filePath="index.html">index.html
    </walkthrough-editor-open-file> を開き 
    <walkthrough-editor-select-line filePath="index.html" startLine="32" endLine="68" startCharacterOffset="0" endCharacterOffset="100">
    callout 以下の p 要素</walkthrough-editor-select-line> を削除してみます。
1.  <walkthrough-editor-spotlight spotlightId="output">Output</walkthrough-editor-spotlight>
    パネルに

    ```terminal
    Update initiated
    Deploy started
    ```

    から始まり、最終的にはやはり以下のようなメッセージが表示されます。

    ```terminal
    http://localhost:8080
    Update successful
    ```
1.  Web プレビュー画面をリロードしてみましょう。

画面は更新されましたか？

## 1.4. デバッグ

ローカルでデバッグしてみましょう。

1.  アプリケーションをデバッグ モードで実行するには
    <walkthrough-editor-spotlight spotlightId="cloud-code-debug-on-cloud-run-emulator">Debug
    on Cloud Run Emulator</walkthrough-editor-spotlight> を選択します。
1.  **デバッグ パネル** が開き、デバッガが実際にアタッチされると、ステータス バーの色が変わります。
1.  **THREADS** を見てください。複数のアプリを並行で起動していくと接続ポートが増えていきますので、
    8080 番ポートでのみ開発をするには **デバッグ ツールバー** から不要なスレッドは停止してください。
    ん。
1.  <walkthrough-editor-open-file filePath="main.go">main.go
    </walkthrough-editor-open-file> を開き 
    <walkthrough-editor-select-line filePath="main.go" startLine="62" endLine="62" startCharacterOffset="0" endCharacterOffset="100">
    63 行目</walkthrough-editor-select-line> にブレイク ポイントを設定します。
1.  Web プレビュー <walkthrough-web-preview-icon/> で待機するポート番号に接続先を適宜変更しつつ、
    またはターミナルから `curl` コマンドなどでサービスにアクセスします。

ブレイク ポイントで停止しましたか？

## 1.5. 個人開発環境へのデプロイ

では個人の Cloud Run へ、Cloud Code プラグインを通してデプロイしてみましょう。

1.  ステータスバーから
    <walkthrough-editor-spotlight spotlightId="cloud-code-status-bar">Cloud
    Code</walkthrough-editor-spotlight> を選択し
1.  <walkthrough-editor-spotlight spotlightId="cloud-code-cloud-run-deploy">Deploy
    to Cloud Run</walkthrough-editor-spotlight> を選びます。
1.  もし Google Cloud への API 実行を確認するポップアップが表示されたら承認してください。
1.  もし Cloud Run API の有効化を求められたら、**Enable API** をクリックします。
1.  **Deploy to Cloud Run** タブが開いたら
1.  サービス名は、個人とわかる識別子 + local としておき
1.  リージョンは東京 (asia-northeast1) を選択しましょう。
1.  ここでは誰からでも接続可能となる **Allow unauthenticated invocations** を選択します。
1.  **Build Settings** の **Bulder** を **Buildpacks** に変更します。
1.  **Deploy** をクリックします。
1.  しばらくすると **Deployment completed successfully! URL** が表示されます。
    実際にアクセスしてみましょう。

Cloud Run としてホストされたサービスにアクセスできましたか？

## 1.6. ログの確認

エミュレータやクラウド上の Cloud Run で出力されるログを確認してみます。

1.  ローカルでは <walkthrough-editor-spotlight spotlightId="output">Output
    </walkthrough-editor-spotlight> の右上で、どこからの出力を表示するかを選択できます。
    **Cloud Run: Run/Debug Locally** ではなく **Cloud Run: Run/Debug Locally - Detailed** を選ぶことでエミュレータ内部で出力されたログが確認できます。
1.  クラウド上の Cloud Run については、左側のメニュー
    <walkthrough-editor-spotlight spotlightId="cloud-code-run-icon">Cloud Run
    Explorer</walkthrough-editor-spotlight> で様々な情報が確認できます。
1.  さきほどデプロイしたサービスの上で右クリック、**View Logs** を選びます。
1.  新しいログを確認するためには
    <walkthrough-editor-spotlight spotlightId="cloud-code-logs-viewer-refresh">Logs
    refresh button</walkthrough-editor-spotlight> が利用できます。

ログは確認できましたか？

<walkthrough-footnote>ここまでで、開発者それぞれに与えられた環境での開発フローを見てきました。ここからは、チームとして製品を開発、CI / CD を回す方法を確認していきましょう。</walkthrough-footnote>

## 2. Cloud Run をベースにした CI/CD

ここからは品質向上のため、そしてチームとして開発する上で重要になる CI/CD を織り込む方法をみていきます。

1. git リポジトリの準備
1. チーム開発環境への CD
1. CI によるテストの自動化
1. 本番環境への CD
1. タグでのテスト & カナリアでのロールアウト

## 2.1. git リポジトリの準備

アプリケーション コードを置く git リポジトリとして [Cloud Source Repositories (CSR)](https://cloud.google.com/source-repositories?hl=ja) を利用します。リポジトリを作成し、Cloud Shell からアクセスするための設定を進めます。

1.  <walkthrough-editor-spotlight spotlightId="menu-terminal-new-terminal">ターミナル
    </walkthrough-editor-spotlight> を開き、改めてプロジェクト ID を指定します。

    ```bash
    export PROJECT_ID={{project-id}}
    ```

    API を有効化、git リポジトリを CSR に作成します。

    ```bash
    gcloud services enable sourcerepo.googleapis.com cloudbuild.googleapis.com artifactregistry.googleapis.com
    gcloud source repos create {{app}}
    ```

1.  CSR への認証ヘルパ含め、git クライアントの設定をします。

    ```bash
    git config --global credential.helper gcloud.sh
    git config --global user.name "$(whoami)"
    git config --global user.email "$(gcloud config get-value core/account)"
    ```

1.  コードを git 管理下におき、CSR へ push しましょう。

    ```bash
    git init
    git remote add google "https://source.developers.google.com/p/${PROJECT_ID}/r/{{app}}"
    git checkout -b main
    git add .
    git commit -m 'init'
    git push google main
    ```

## 2.2. チーム開発環境への CD

Cloud Run には [継続的デリバリーを簡単に実施する仕組み](https://cloud.google.com/blog/ja/products/application-development/cloud-run-integrates-with-continuous-deployment) があります。チームで共有する開発環境に対しては、**最新の main ブランチがリリースされ続ける**ように設定してみましょう。

1.  Cloud Run の Web コンソールを開きます。
    <walkthrough-menu-navigation sectionId="CLOUD_RUN_SECTION"></walkthrough-menu-navigation>
1.  <walkthrough-spotlight-pointer spotlightId="run-create-service">サービスの作成
    </walkthrough-spotlight-pointer> ボタンをクリックし作成を開始します。

1.  サービス名は `{{app}}-dev` とし、
    リージョンは `asia-northeast1 (Tokyo)` を選んで `次へ` をクリックします
    [![screenshot](https://raw.githubusercontent.com/{{github}}/images/link_image.png)](https://raw.githubusercontent.com/{{github}}/images/create_a_cloud_run_service.png)
1.  `ソース リポジトリから新しいリビジョンを継続的にデプロイする` をチェックして
    `SET UP WITH CLOUD BUILD` ボタンをクリックします
    [![screenshot](https://raw.githubusercontent.com/{{github}}/images/link_image.png)](https://raw.githubusercontent.com/{{github}}/images/configure_the_first_revision_of_the_service.png)

1.  リポジトリ プロバイダで `Cloud Source Repositories` を、リポジトリは `{{app}}` を選び、
    `次へ` をクリックします

1.  ブランチは `^main$`、Build Type は `Go、Node.js、Python、Java、または .NET Core`
    をチェックし、ビルド コンテキストのディレクトリは `/` のまま
    `保存` をクリック、続けて `次へ` をクリックしましょう

1.  認証の項目で `未認証の呼び出しを許可` をチェックし、`作成` ボタンをクリックします

サービスが作成されたらサイトにアクセスし、git push により環境が更新されるのを確認してみましょう。

```bash
rm -rf Dockerfile img/
git add .
git commit -m 'Delete unnecessary files'
git push google main
```

## 2.3. CI によるテストの自動化

2.2. のままでは、ソフトウェアが壊れてもアクセスするまで気付けません。ここでは git push と同時にテスト実行 + ビルドするステップを自動化してみましょう。

1.  コンテナ レジストリを作ります。

    ```bash
    gcloud artifacts repositories create {{app}} --repository-format=docker --location={{region}} --description="Docker repository for CI/CD hands-on"
    gcloud auth configure-docker {{region}}-docker.pkg.dev
    docker pull alpine:3.14
    docker tag alpine:3.14 {{region}}-docker.pkg.dev/${PROJECT_ID}/{{app}}/app:init
    docker push {{region}}-docker.pkg.dev/${PROJECT_ID}/{{app}}/app:init
    ```

1.  Cloud Build に対して必要な権限を付与します。

    ```bash
    project_number="$( gcloud projects list --filter="${PROJECT_ID}" --format='value(PROJECT_NUMBER)' )"
    gcloud projects add-iam-policy-binding ${PROJECT_ID} --member "serviceAccount:${project_number}@cloudbuild.gserviceaccount.com" --role roles/run.admin
    gcloud iam service-accounts add-iam-policy-binding ${project_number}-compute@developer.gserviceaccount.com --member="serviceAccount:${project_number}@cloudbuild.gserviceaccount.com" --role="roles/iam.serviceAccountUser"
    ```

1.  Cloud Build の設定ファイル、`cloudbuild-ci.yaml` を作ります。

    ```text
    cat << EOF > cloudbuild-ci.yaml
    steps:
    - id: Static Analysis
      name: golangci/golangci-lint:v1.42.0
      args: ['golangci-lint', 'run']
    - id: Build
      name: gcr.io/k8s-skaffold/pack
      entrypoint: pack
      args:
      - build
      - test
      - '--builder=gcr.io/buildpacks/builder:v1'
      - '--path=.'
    tags: ['test']
    EOF
    ```

1.  git push により CI が起動するようトリガーを設定します。

    ```bash
    gcloud beta builds triggers create cloud-source-repositories --name {{app}}-ci --repo={{app}} --branch-pattern='.*' --build-config=cloudbuild-ci.yaml
    ```

1.  Cloud Build コンソールを開きましょう。
    <walkthrough-menu-navigation sectionId="CLOUD_BUILD_SECTION"></walkthrough-menu-navigation>

1.  git push によりビルドが始まることを確認します。

    ```bash
    git add cloudbuild-ci.yaml
    git commit -m 'Add continuous integration'
    git push google main
    ```

## 2.4. 本番環境への CD

本番環境を作り、そこへの継続的デリバリーパイプラインを作成します。ただし開発環境とは以下の点が異なります。

- コンテナ イメージは明示的に保存される
- ユーザからのアクセスは新バージョンには流さない

1.  本番環境想定の Cloud Run サービスを作成しましょう。

    ```bash
    gcloud run deploy {{app}}-prod --image gcr.io/cloudrun/hello --region={{region}} --platform=managed --allow-unauthenticated --quiet
    ```

1.  Cloud Build の設定ファイル、`cloudbuild-cd.yaml` を作ります。

    ```text
    cat << EOF > cloudbuild-cd.yaml
    steps:
    - id: Build
      name: gcr.io/k8s-skaffold/pack
      entrypoint: pack
      args:
      - build
      - '{{region}}-docker.pkg.dev/${PROJECT_ID}/{{app}}/app:\${SHORT_SHA}'
      - '--builder=gcr.io/buildpacks/builder:v1'
      - '--path=.'
    - id: Push
      name: gcr.io/cloud-builders/docker
      args:
      - push
      - '{{region}}-docker.pkg.dev/${PROJECT_ID}/{{app}}/app:\${SHORT_SHA}'
    - id: Deploy
      name: 'gcr.io/google.com/cloudsdktool/cloud-sdk:slim'
      entrypoint: gcloud
      args:
      - run
      - deploy
      - {{app}}-prod
      - '--image={{region}}-docker.pkg.dev/${PROJECT_ID}/{{app}}/app:\${SHORT_SHA}'
      - '--region={{region}}'
      - '--platform=managed'
      - '--allow-unauthenticated'
      - '--no-traffic'
      - '--tag=v\${SHORT_SHA}'
      - '--quiet'
    images:
    - '{{region}}-docker.pkg.dev/${PROJECT_ID}/{{app}}/app:\${SHORT_SHA}'
    tags: ['prod']
    EOF
    ```

1.  **main ブランチへの** git push により CI が起動するようトリガーを設定します。

    ```bash
    gcloud beta builds triggers create cloud-source-repositories --name {{app}}-cd-prod --repo={{app}} --branch-pattern='^main$' --build-config=cloudbuild-cd.yaml
    ```

1.  **main ブランチへの** git push によりデプロイが始まることを確認します。

    ```bash
    sed -ie "s|Congratulations|Congratulations $(whoami)|" index.html
    git add cloudbuild-cd.yaml index.html
    git commit -m 'Add continuous delivery'
    git push google main
    ```

新しいサービスはデプロイされましたが、このサービスのエンドポイントからは新しいバージョンにアクセスできません。ぜひ一度確かめてみてください。

## 2.5. タグでのテスト & カナリアでのロールアウト

2.4. での設定をよくみると、`--tag` というオプションがついています。実はタグをつけてデプロイすると、そのタグがプリフィックスとして付与された特別な URL が払い出され、[新しいバージョンにアクセスできます](https://cloud.google.com/run/docs/rollouts-rollbacks-traffic-migration?hl=ja#tags)。試してみましょう。

1.  タグに付与された URL を取得し、実際にブラウザからアクセスしてみましょう。

    ```bash
    gcloud run services describe {{app}}-prod --region {{region}} --format='value(status.address.url)' | sed -e "s/{{app}}/v$(git rev-parse --short HEAD)---{{app}}/"
    ```

1.  新バージョンをテストし、問題なければユーザからのトラフィックの 10% を振り向けます。

    ```bash
    gcloud run services update-traffic {{app}}-prod --region {{region}} --to-tags "v$(git rev-parse --short HEAD)=10"
    ```

1.  [SLO に関するメトリクス](https://cloud.google.com/architecture/defining-SLOs?hl=ja)に変化がなければ、新バージョンのサービスを 100% ロールアウトします。

    ```bash
    gcloud run services update-traffic {{app}}-prod --region {{region}} --to-tags "v$(git rev-parse --short HEAD)=100"
    ```

タグによる関係者のみのテストや段階的なロールアウトにより、より信頼性を担保しやすい仕組みが実感できましたでしょうか？

## 3. 高度なデプロイオプションの利用

Google Cloud には [Binary Authorization](https://cloud.google.com/binary-authorization?hl=ja) という機能があります。信頼できるコンテナ イメージのみが稼働することを支援する機能で、署名による保護や許可したリポジトリからのみデプロイを許可するといったことが可能です。ここでは後者を実装します。

1. ポリシーの設定
1. BinAuth の挙動確認

## 3.1. ポリシーの設定

今回作成したコンテナ レジストリ以外からのデプロイを拒否するよう、ポリシーの設定を更新します。

1.  Binary Authorization を有効化します。

    ```bash
    gcloud services enable binaryauthorization.googleapis.com
    gcloud beta run services update {{app}}-prod --region {{region}} --binary-authorization=default
    ```

1.  ポリシーの YAML ファイルをエクスポートし、中身を確認してみます。

    ```bash
    gcloud container binauthz policy export > /tmp/policy.yaml
    cat /tmp/policy.yaml
    ```

1.  ポリシーを書き換えます。

    ```text
    cat << EOF > /tmp/policy.yaml
    admissionWhitelistPatterns:
    - namePattern: gcr.io/google_containers/*
    - namePattern: gcr.io/google-containers/*
    - namePattern: k8s.gcr.io/*
    - namePattern: gke.gcr.io/*
    - namePattern: gcr.io/stackdriver-agents/*
    - namePattern: {{region}}-docker.pkg.dev/${PROJECT_ID}/{{app}}/app@*
    globalPolicyEvaluationMode: ENABLE
    defaultAdmissionRule:
      enforcementMode: ENFORCED_BLOCK_AND_AUDIT_LOG
      evaluationMode: ALWAYS_DENY
    name: projects/${PROJECT_ID}/policy
    EOF
    ```

1.  ポリシーを更新します。

    ```bash
    gcloud container binauthz policy import /tmp/policy.yaml
    ```

## 3.2. BinAuth の挙動確認

明示的にリポジトリが許可されていないイメージのデプロイは拒否され、指定したリポジトリのものであればデプロイできる様子を確かめます。

1.  先程は問題なかった hello world コンテナのデプロイが失敗することを確認します。

    ```bash
    gcloud run deploy {{app}}-prod --image gcr.io/cloudrun/hello --region={{region}} --platform=managed --allow-unauthenticated --quiet
    ```

1.  git push からのデプロイは正常に行われる様子をみてみます。

    ```bash
    sed -ie "s|running|running and protected|" index.html
    git add index.html
    git commit -m 'Revised'
    git push google main
    ```

1.  Cloud Build コンソールの履歴をみつつ
    <walkthrough-menu-navigation sectionId="CLOUD_BUILD_SECTION"></walkthrough-menu-navigation>

1.  本番環境へリリースされたら、タグの URL から変更内容を確認してみましょう。

    ```bash
    curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" $(gcloud run services describe {{app}}-prod --region {{region}} --format='value(status.address.url)' | sed -e "s/{{app}}/v$(git rev-parse --short HEAD)---{{app}}/")
    ```

## 4. クリーンアップ

ハンズオンに利用したプロジェクトを削除し、課金を止めます。

```bash
gcloud projects delete ${PROJECT_ID}
```

プロジェクトがそのまま消せない場合は、以下のリソースを個別に削除してください。

- Cloud Run サービス
- Cloud Build のトリガー
- Cloud Source Repositories の git リポジトリ
- Artifact Registry の コンテナ リポジトリ

## これで終わりです

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

すべて完了しました。
