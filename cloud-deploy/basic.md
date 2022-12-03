# Cloud Deploy による継続的デリバリー ハンズオン

<walkthrough-watcher-constant key="region" value="asia-northeast1"></walkthrough-watcher-constant>

## 始めましょう

[Cloud Deploy](https://cloud.google.com/deploy?hl=ja) を使った Google Cloud での継続的デリバリー (CD) を体験いただくハンズオンです。以下の流れで進めます。

1. **Skaffold によるローカルでのテストとビルド**
1. **GitHub Actions による継続的なテストとビルド**
1. **Cloud Deploy による継続的デリバリー**

<walkthrough-tutorial-duration duration="45"/> 
**所要時間**: 約 45 分

**前提条件**:

- Google Cloud 上にプロジェクトが作成してある
- プロジェクトの *編集者* 相当の権限をもつユーザーでログインしている
- *プロジェクト IAM 管理者* 相当の権限をもつユーザーでログインしている
- （推奨）Google Chrome を利用している

**[開始]** ボタンをクリックして次のステップに進みます。

## プロジェクトの設定

この手順の中で実際にリソースを構築する対象のプロジェクトを選択してください。

<walkthrough-project-setup billing=true></walkthrough-project-setup>

## CLI の初期設定と権限の確認

gcloud（[Google Cloud の CLI ツール](https://cloud.google.com/sdk/gcloud?hl=ja)
のデフォルト プロジェクト、Cloud Deploy のデフォルト リージョンを設定します。

```bash
export PROJECT_ID=<walkthrough-project-id/>
```

```bash
gcloud config set project "${PROJECT_ID}"
gcloud config set deploy/region "{{region}}"
```

念のため、[Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine?hl=ja) を扱える権限があることを再確認します。

```bash
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member "user:$(gcloud config get-value core/account)" --role roles/container.admin
```

## 1. Skaffold によるローカルでのテストとビルド

[Cloud Shell エディタ](https://cloud.google.com/shell/docs/launching-cloud-shell-editor?hl=ja) は個人ごとに割り当てられる開発環境としてご利用いただけます。このセクションでは、Cloud Shell エディタを使いつつ、以下のような流れでローカル環境でのテストとビルドを実施します。

1. サンプルコードのダウンロード
1. 単体テストの実行
1. アプリの起動
1. コードの変更、リアルタイムでの更新
1. ログの確認

## 1.1. サンプルコードのダウンロード

プログラミング言語はなんでもよいのですが、ここでは Dart のサンプルコードを利用してみます。

1.  作業ディレクトリを作ります。

    ```bash
    mkdir -p "${HOME}/dart-app"
    ```

1.  コードをダウンロードし

    ```bash
    git clone https://github.com/dart-lang/samples.git
    cp -r samples/server/simple/. "${HOME}/dart-app/src"
    rm -rf samples
    ```

1.  作業ディレクトリを *ワークスペース* として、Cloud Shell エディタを起動します。

    ```bash
    cloudshell workspace dart-app
    ```

## 1.2. 単体テストの実行

新しい <walkthrough-editor-spotlight spotlightId="menu-terminal-new-terminal">ターミナル</walkthrough-editor-spotlight> を開き、テストを実行してみましょう。

```bash
cd "${HOME}/dart-app/src"
dart test
```

Cloud Shell では、dart コマンドを打つとインストールが始まります。

初回ということもありテスト完了まで少々時間を要しますが、期待通り緑色で終わることを確認してください。

## 1.3. アプリの起動まで

1.  ローカル開発、ビルド、デプロイに必要な設定ファイルをダウンロードします。

    ```bash
    cd "${HOME}"
    git clone https://github.com/google-cloud-japan/appdev-cicd-handson.git
    cp -r appdev-cicd-handson/cloud-deploy/sample-resources/kustomize/. "${HOME}/dart-app"
    rm -rf appdev-cicd-handson
    cd "${HOME}/dart-app"
    echo -e ".theia\ncredential.json\ndeploy/clouddeploy.yaml" > .gitignore
    ```

1.  Minikube を起動しましょう。

    ```bash
    minikube start
    ```

1.  結果が以下のように表示されたら

    ```terminal
    🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
    ```

1.  Skaffold を最新（v2 系）にして

    ```bash
    curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
    chmod +x skaffold && sudo mv skaffold /usr/bin
    skaffold version
    ```

1.  アプリケーションをビルドし、デプロイすべく
    <walkthrough-editor-spotlight spotlightId="cloud-code-status-bar">Cloud
    Code</walkthrough-editor-spotlight> のメニューから
    <walkthrough-editor-spotlight spotlightId="cloud-code-run-on-k8s">Run
    on Kubernetes</walkthrough-editor-spotlight> を選択します。

1.  選択肢がポップアップしてきたら *[default]* を選択します。

1.  サービスがデプロイされると、
    <walkthrough-editor-spotlight spotlightId="output">Output</walkthrough-editor-spotlight>
    パネルに以下のように表示されます。

    ```terminal
    Forwarded URL from service web-svc: http://localhost:8080
    Update succeeded
    ```

1.  Web preview ボタン <walkthrough-web-preview-icon/> を押し、
    "ポート 8080 でプレビュー" を選んでみましょう。

アプリの起動はうまくいきましたか？

## 1.4. コードの変更、リアルタイムでの更新

コードの書き換えによって、アプリケーションがリアルタイムに更新されることを確認します。

1.  <walkthrough-editor-open-file filePath="dart-app/src/public/index.html">index.html
    </walkthrough-editor-open-file> を開き、HTML を一部変更してみてください。

1.  <walkthrough-editor-spotlight spotlightId="output">Output</walkthrough-editor-spotlight>
    パネルに

    ```terminal
    Update initiated
    Build started for artifact app
    ```

    から始まり、最終的にはやはり以下のようなメッセージが表示されます。

    ```terminal
    http://localhost:8080
    Update successful
    ```

1.  Web プレビュー画面をリロードしてみましょう。

画面は更新されましたか？

## 1.5. ログの確認

Minikube 上に出力されるログを確認してみます。

1.  ローカルでは <walkthrough-editor-spotlight spotlightId="output">Output
    </walkthrough-editor-spotlight> の右上で、どこからの出力を表示するかを選択できます。
    **Kubernetes: Run/Debug Local** ではなく **Kubernetes: Run/Debug Local - Detailed** を選ぶことでエミュレータ内部で出力されたログが確認できます。

ログは確認できましたか？

<walkthrough-footnote>ここまでローカルの開発環境を見てきました。次に GitHub Actions を使って CI としてテスト・ビルドが行えることを確認していきましょう。</walkthrough-footnote>

## 2. GitHub Actions による継続的なテストとビルド

ここからは品質向上のため、そしてアジリティ高く開発する上で重要になる CI/CD を織り込んでいきます。

1. **Google Cloud のリソース作成**
1. **GitHub での Secrets 登録**
1. **git リポジトリの準備**
1. **CI によるテスト・ビルドの自動化**

## 2.1. Google Cloud のリソース作成

1.  <walkthrough-editor-spotlight spotlightId="menu-terminal-new-terminal">ターミナル
    </walkthrough-editor-spotlight> を開き、この後利用する Google Cloud の機能を有効化します。

    ```bash
    gcloud services enable cloudresourcemanager.googleapis.com compute.googleapis.com container.googleapis.com serviceusage.googleapis.com stackdriver.googleapis.com monitoring.googleapis.com logging.googleapis.com clouddeploy.googleapis.com cloudbuild.googleapis.com artifactregistry.googleapis.com
    ```

1.  コンテナのリポジトリを Artifact Registry に作り

    ```bash
    gcloud artifacts repositories create my-apps --repository-format=docker --location {{region}} --description="Docker repository for CI/CD hands-on"
    ```

1.  実行環境として GKE クラスタを 1 つ作成します。

    ```bash
    gcloud container clusters create-auto my-gke --region {{region}} --release-channel stable
    ```

1.  GitHub に渡すサービスアカウントと、鍵を生成します。

    ```bash
    export PROJECT_ID=<walkthrough-project-id/>
    gcloud iam service-accounts create sa-github
    gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:sa-github@${PROJECT_ID}.iam.gserviceaccount.com" --role="roles/storage.admin"
    gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:sa-github@${PROJECT_ID}.iam.gserviceaccount.com" --role="roles/artifactregistry.writer"
    gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:sa-github@${PROJECT_ID}.iam.gserviceaccount.com" --role="roles/clouddeploy.releaser"
    PROJECT_NUMBER="$( gcloud projects list --filter="${PROJECT_ID}" --format='value(PROJECT_NUMBER)' )"
    gcloud iam service-accounts add-iam-policy-binding ${PROJECT_NUMBER}-compute@developer.gserviceaccount.com --member="serviceAccount:sa-github@${PROJECT_ID}.iam.gserviceaccount.com" --role="roles/iam.serviceAccountUser"
    gcloud iam service-accounts keys create credential.json --iam-account=sa-github@${PROJECT_ID}.iam.gserviceaccount.com
    cat credential.json
    ```

## 2.2. GitHub での Secrets 登録

GitHub から Google Cloud 上のリソースにアクセスするための変数を、**リポジトリの Secret にセット**します。（リポジトリ名などは仮ですが、おおよそ以下の URL でアクセスできる設定画面です https://github.com/your-org/your-repogitory/settings/secrets/actions ）

- **GOOGLECLOUD_PROJECT_ID**: プロジェクト ID
- **GOOGLECLOUD_SA_KEY**: 2.1 の最後に出力された JSON 鍵

## 2.3. git リポジトリの準備

GitHub へアクセスする準備を進めます。

1.  GitHub へアクセスするための SSH 鍵を生成します。

    ```bash
    ssh-keygen -t rsa -b 4096 -C "$(whoami)" -N "" -f "$HOME/.ssh/github_rsa"
    cat "$HOME/.ssh/github_rsa.pub"
    ```

1.  GitHub の設定
    [https://github.com/settings/keys](https://github.com/settings/keys) を開き、
    1 で出力された公開鍵を SSH Key として登録してください。

1.  SSH config に設定を加え、GitHub に接続できることを確認します。

    ```text
    cat << EOF > ~/.ssh/config
    Host github.com
        User $(gcloud config get-value core/account)
        IdentityFile $HOME/.ssh/github_rsa
    EOF
    ssh -T git@github.com
    ```

## 2.4. CI によるテスト・ビルドの自動化

1.  git クライアントの設定をしたら

    ```bash
    git config --global user.name "$(whoami)"
    git config --global user.email "$(gcloud config get-value core/account)"
    ```

1.  コードを git 管理下におき、GitHub リポジトリへ push します。

    ```bash
    git init
    git remote add origin git@github.com:your-org/your-repogitory.git
    git add --all
    git commit -m "add ci/cd templates"
    git branch -M main
    git push -u origin main
    ```

1.  GitHub Actions の実行履歴を確認しましょう。
    3 つの Job が実行されていますが、最後のものが失敗していると思います。

    - Test code: 単体テストは緑
    - Test template: k8s リソースも緑で正常
    - Release: エラーがでており、赤

    詳細を確認いただくと、
    まだ Cloud Deploy のリソースを作成していないことが原因でエラーになっています。

1.  とはいえアプリケーションは正常にビルドされているかと思います。
    Artifact Registry コンソールを開きましょう。
    <walkthrough-menu-navigation sectionId="ARTIFACT_REGISTRY_SECTION"></walkthrough-menu-navigation>
    **my-apps/app** というリポジトリに git ハッシュのタグでイメージが確認できます。

1.  実際にイメージのビルドとプッシュを担当しているのは
    <walkthrough-editor-select-line filePath="dart-app/.github/workflows/release.yaml" startLine="88" endLine="88" startCharacterOffset="0" endCharacterOffset="100">release.yaml</walkthrough-editor-select-line> の 89 行目です。

    **Skaffold でビルドやデプロイをラップしておくことで、実際にビルドの方法が変わっても、デプロイ先が変わっても CI のステップを変更する必要がなくなります。**


## 3. Cloud Deploy による継続的デリバリー

Cloud Deploy を使って GKE へアプリケーションを継続的にデプロイする仕組みをみていきましょう。以下のステップで進めます。

1. **パイプラインの作成**
1. **リリースの作成 & dev 環境へのデプロイ**
1. **デプロイ履歴の確認**
1. **prod 環境へのプロモーション**

## 3.1. パイプラインの作成

まずはパイプラインを作成します。Cloud Deploy における **パイプライン** は、一般的によくみる CI/CD パイプラインのようなステップやタスクといった概念はなく、あくまで **“どこにどういう順序でデプロイするか”** を決めておく単位です。また、デプロイ先の環境のことは **ターゲット** と呼びます。

1.  <walkthrough-editor-open-file filePath="dart-app/deploy/clouddeploy.yaml">deploy/clouddeploy.yaml</walkthrough-editor-open-file> を開き `your-project-id` をご自身のプロジェクト ID に書き換えてください。

1.  パイプラインを作成しましょう。

    ```bash
    gcloud deploy apply --file deploy/clouddeploy.yaml --region {{region}}
    ```

    clouddeploy.yaml は Cloud Deploy のパイプライン定義で、以下のことを宣言しています。

    - **dev** というターゲットがあり、具体的なデプロイ先は `my-gke` という GKE クラスタ
    - **prod** というターゲットもある、デプロイ先は同じく `my-gke` という GKE クラスタ
    - パイプラインでは **dev → prod の順にデプロイ** していく
    - prod については `prod` というプロファイルを利用する

1.  パイプラインの状態を確認しましょう。Cloud Deploy のコンソールを開きます。
    <walkthrough-menu-navigation sectionId="CLOUD_DEPLOY_SECTION"></walkthrough-menu-navigation>

## 3.2. リリースの作成 & dev 環境へのデプロイ

パイプラインもできたので、次にリリースを作成します。

**リリース** は、**一緒にデプロイしたい成果物をまとめる単位** です。Cloud Deploy では、このひとつのリリースを各ターゲットにデプロイする（これを **ロールアウト** と呼びます）ことで、検証したバイナリや設定ができる限り同じ状態で次の環境にデプロイされることを期待していただけます。

先程まさにエラーになったところですが、このリリースの作成は GitHub Actions の最後のステップに組み込まれています。

<walkthrough-editor-select-line filePath="dart-app/.github/workflows/release.yaml" startLine="97" endLine="97" startCharacterOffset="0" endCharacterOffset="300">release.yaml</walkthrough-editor-select-line> の 98 行目です。

3.1 でパイプラインを作りましたので、2.4 で失敗した Actions のジョブ画面を開き、改めて **Re-run jobs** を押してみてください。

## 3.3. デプロイ履歴の確認

GitHub Actions や Cloud Deploy パイプラインの状況を確認してみてください。

1.  GitHub Actions の 3 つのジョブがすべて緑色になるまでお待ち下さい。

1.  パイプラインの状態を確認しましょう。Cloud Deploy のコンソールを開きます。
    <walkthrough-menu-navigation sectionId="CLOUD_DEPLOY_SECTION"></walkthrough-menu-navigation>

1.  パイプラインを選択し **RELEASES** や **ターゲット** タブをひらき
    どんな情報がどう管理されているかもぜひ確認してみてください。

1.  GKE のワークロードの変化も確認してみましょう。GKE のコンソールを開きます。
    <walkthrough-menu-navigation sectionId="KUBERNETES_SECTION"></walkthrough-menu-navigation>

1.  左側のメニューから <walkthrough-spotlight-pointer cssSelector="#cfctest-section-nav-item-workloads">ワークロード</walkthrough-spotlight-pointer> を開きます。
    **web-app** という Deployment が見えていれば成功です！

## 3.4. prod 環境へのプロモーション

画面からもできるのですが、ここでは GitHub Actions に仕込んだ git のタグ打ちでプロモーションする様子をみてみます。

1.  GitHub Actions の定義をみてみましょう。
    <walkthrough-editor-select-line filePath="dart-app/.github/workflows/promotion.yaml" startLine="28" endLine="28" startCharacterOffset="0" endCharacterOffset="300">promotion.yaml</walkthrough-editor-select-line> の 29 行目です。

1.  では実際にプロモーションをしましょう。
    <walkthrough-editor-select-line filePath="dart-app/.github/workflows/promotion.yaml" startLine="5" endLine="5" startCharacterOffset="0" endCharacterOffset="100">6 行目</walkthrough-editor-select-line> を見ると、prod- から始まる名前のタグを打つとこのジョブが起動しそうです。

    ```bash
    git tag prod-1.0
    git push origin prod-1.0
    ```

1.  dev へデプロイしたとき同様、Actions や Cloud Deploy、GKE の各画面から
    prod についても問題なくデプロイができたことを確認してみてください。
    <walkthrough-menu-navigation sectionId="CLOUD_DEPLOY_SECTION"></walkthrough-menu-navigation>

## 4. クリーンアップ

ハンズオンに利用したプロジェクトを削除し、課金を止めます。

```bash
gcloud projects delete ${PROJECT_ID}
```

プロジェクトがそのまま消せない場合は、以下のリソースを個別に削除してください。

```bash
gcloud deploy delivery-pipelines delete kustomize-pipeline --force --region {{region}} --quiet
gcloud artifacts repositories delete my-apps --location {{region}} --quiet
gcloud container clusters delete my-gke --region {{region}} --quiet
```

## これで終わりです

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

すべて完了しました。
