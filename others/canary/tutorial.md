# Kubernetes Gateway API、Flagger、Google Cloud Deploy を使用したカナリア デプロイ

このチュートリアルでは [Gateway API](https://gateway-api.sigs.k8s.io/) に対応した [Flagger](https://flagger.app/) と [Cloud Deploy](https://cloud.google.com/deploy?hl=ja) を利用して、GKE 上のアプリケーションをカナリア デプロイする体験ができます。

<walkthrough-tutorial-duration duration="45"></walkthrough-tutorial-duration>
<walkthrough-tutorial-difficulty difficulty="4"></walkthrough-tutorial-difficulty>

## プロジェクトの選択

ハンズオンを行う Google Cloud プロジェクトを選択して **Start** をクリックしてください。

<walkthrough-project-setup />

## チュートリアルの流れ

前半では

1. GKE クラスタを作ります
1. K8s の Namespace を利用し、dev 環境と prod 環境を作ります
1. カナリア デプロイのための指標を収集する準備をします
1. Flagger をインストールします
1. コンテナイメージの格納場所を作り
1. サンプル アプリケーションをデプロイします

後半では

1. Cloud Deploy を使って成果物をとりまとめ
1. dev 環境へデプロイします
1. prod 環境へデプロイし、カナリア デプロイの様子を確認します

## サンプルコードのダウンロード

```bash
git clone https://github.com/google-cloud-japan/appdev-cicd-handson.git
```

サンプルコードのディレクトリを *ワークスペース* として [Cloud Shell エディタ](https://cloud.google.com/shell/docs/launching-cloud-shell-editor?hl=ja)を起動します。

```bash
cloudshell workspace appdev-cicd-handson/others/canary/sample-resources
```

## 1. サービスの起動

サンプルアプリケーションを利用し、外部から接続できるサービスを起動します。<walkthrough-editor-select-line filePath="src/main.go" startLine="31" endLine="31" startCharacterOffset="0" endCharacterOffset="100">アプリは Go 言語で書かれており、</walkthrough-editor-select-line>環境変数からアプリケーションのバージョンや環境を応答する仕組みをもっています。

## 1.1. API の有効化

gcloud でプロジェクトを設定してください。

```bash
gcloud config set project "<walkthrough-project-id />"
export PROJECT_ID=$(gcloud config get-value project)
```

リージョン・ゾーンには東京を指定しましょう。

```bash
export GOOGLE_CLOUD_REGION=asia-northeast1
export GOOGLE_CLOUD_ZONE=asia-northeast1-a
```

<walkthrough-editor-spotlight spotlightId="menu-terminal-new-terminal">ターミナル</walkthrough-editor-spotlight> を開き、コマンドを実行していきましょう。

Google Cloud では、プロジェクトごとに、利用する **サービスの API** を有効化することでリソースが利用できるようになります。GKE クラスタや Cloud Deploy を利用することから、以下の API を有効化します。

```bash
gcloud services enable compute.googleapis.com container.googleapis.com clouddeploy.googleapis.com artifactregistry.googleapis.com
```

## 1.2. GKE クラスタの作成

クラスタを作成します。以下のコマンドを実行したら 5 分ほど待ちます。

```bash
gcloud container clusters create canary-cluster --zone "${GOOGLE_CLOUD_ZONE}" --machine-type "e2-standard-2" --gateway-api "standard" --enable-managed-prometheus --num-nodes "1" --min-nodes "0" --max-nodes "3" --enable-autoscaling --enable-ip-alias --workload-pool "${PROJECT_ID}.svc.id.goog"
gcloud container clusters get-credentials "canary-cluster" --zone "${GOOGLE_CLOUD_ZONE}"
```

dev 環境と prod 環境を Namespace として用意しましょう。

```bash
kubectl create namespace prod
kubectl create namespace dev
```

作成されたクラスタや、今後作られるさまざまなリソースはクラウドのコンソールからも確認できます。
<walkthrough-menu-navigation sectionId="KUBERNETES_SECTION"></walkthrough-menu-navigation>

例えばワークロードは <walkthrough-spotlight-pointer cssSelector="#cfctest-section-nav-item-workloads">このメニュー</walkthrough-spotlight-pointer> から確認できます。

## 1.3. GMP 関連リソースの作成

[Google Cloud Managed Service for Prometheus (GMP)](https://cloud.google.com/stackdriver/docs/managed-prometheus?hl=ja) を活用して、カナリア デプロイに必要な指標をクエリするためのフロントエンドをデプロイします。

まずはクラウド側にサービスアカウントを作り

```bash
gcloud iam service-accounts create gsa-gmp
gcloud iam service-accounts add-iam-policy-binding "gsa-gmp@${PROJECT_ID}.iam.gserviceaccount.com" --member "serviceAccount:${PROJECT_ID}.svc.id.goog[prod/gmp]" --role "roles/iam.workloadIdentityUser"
gcloud projects add-iam-policy-binding "${PROJECT_ID}" --member "serviceAccount:gsa-gmp@${PROJECT_ID}.iam.gserviceaccount.com" --role "roles/monitoring.viewer"
```

Kubernetes 内にもサービスアカウントを作ります。

```bash
kubectl create -n prod serviceaccount gmp
kubectl annotate -n prod serviceaccount gmp "iam.gke.io/gcp-service-account=gsa-gmp@${PROJECT_ID}.iam.gserviceaccount.com"
```

Prometheus フロントエンドをインストールしましょう。

```bash
sed -i "s/GOOGLE_CLOUD_PROJECT_ID/${PROJECT_ID}/g" gmp-frontend.yaml
kubectl apply -n prod -f gmp-frontend.yaml
```

## 1.4. ロードバランサ向けサブネットの作成

本ワークショップではトラフィック分割をすることから、Envoy ベースのロードバランサを利用します。その準備として、[Virtual Private Cloud (VPC)](https://cloud.google.com/vpc?hl=ja) 内には [ロードバランサ向けプロキシ専用サブネット](https://cloud.google.com/load-balancing/docs/proxy-only-subnets?hl=ja) が存在している必要があります。次のコマンドを使用して作成しましょう。（IP 範囲が使用中の場合は、ネットワーク内の未使用の範囲に変更する必要があります）

```bash
gcloud compute networks subnets create proxy --purpose "REGIONAL_MANAGED_PROXY" --network "default" --region "${GOOGLE_CLOUD_REGION}" --range "10.103.0.0/23" --role "ACTIVE"
```

## 1.5. Flagger のインストール

プログレッシブ デリバリーを実現する [Kubernetes Operator](https://kubernetes.io/ja/docs/concepts/extend-kubernetes/operator/) である [Flagger](https://flagger.app/) をインストールしていきます。

まずは Kubernetes クラスタにリソースを作成し、

```bash
kubectl apply -k github.com/fluxcd/flagger/kustomize/gatewayapi
```

そしてロールアウトの基準となる<walkthrough-editor-select-line filePath="bootstrap.yaml" startLine="42" endLine="44" startCharacterOffset="0" endCharacterOffset="100">指標の定義</walkthrough-editor-select-line>を作成します。

併せて [Gateway API](https://gateway-api.sigs.k8s.io/) を経由して [負荷分散](https://cloud.google.com/load-balancing?hl=ja) のための<walkthrough-editor-select-line filePath="bootstrap.yaml" startLine="22" endLine="22" startCharacterOffset="0" endCharacterOffset="100">リソース</walkthrough-editor-select-line>も生成します。

作成には 5 分ほどかかります。

```bash
kubectl apply -f bootstrap.yaml
kubectl wait --all-namespaces gateways --all --for "condition=READY" --timeout "450s"
```

外部 HTTP(S) ロードバランサの IP アドレスを確認しておきます。

```text
dev_ip="$(kubectl get gateways/app -o jsonpath='{.status.addresses[0].value}' -n dev)"
prod_ip="$(kubectl get gateways/app -o jsonpath='{.status.addresses[0].value}' -n prod)"
cat << EOS

dev  環境 IP アドレス: ${dev_ip}
prod 環境 IP アドレス: ${prod_ip}

EOS
```

## 1.6. Artifact Registry の設定

コンテナを保存するリポジトリを作り、ローカルに Docker credential helper を設定します。

```bash
gcloud artifacts repositories create canary-apps --repository-format "docker" --location "${GOOGLE_CLOUD_REGION}" --description "Docker repository for sample apps"
gcloud auth configure-docker "${GOOGLE_CLOUD_REGION}-docker.pkg.dev" --quiet
```

## 1.7. サンプル アプリケーションのデプロイ

ではいきなりですが、ここで dev 環境へデプロイしてみましょう。

ここからは [Skaffold](https://skaffold.dev/) が大活躍します。Skaffold はもともとローカルでのコンテナ開発のツールとして誕生しましたが、その後、デプロイ・検証までを含めたフィードバックループを高速化するための手段として成熟しつつあります。

<walkthrough-editor-open-file filePath="skaffold.yaml">こちら</walkthrough-editor-open-file>が Skaffold の定義です。アプリケーションをどのようにビルドし、どうやってマニフェストを組み立てるかが環境別に定義されています。

以下 [run コマンド](https://skaffold.dev/docs/references/cli/#skaffold-run) はビルドからデプロイまでを一気に実行するためのコマンドです。kubeconf の接続先が現在は GKE クラスタとなっているため、`--default-repo` オプションを指定し、コンテナが Artifact Registry を経由してデプロイされるよう調整します。

```bash
skaffold run --default-repo "${GOOGLE_CLOUD_REGION}-docker.pkg.dev/${PROJECT_ID}/canary-apps" --status-check
```

`HTTPRoute` 内で <walkthrough-editor-select-line filePath="overlays/dev/gateway.yaml" startLine="20" endLine="20" startCharacterOffset="0" endCharacterOffset="100">hostname として定義されたホスト名</walkthrough-editor-select-line>を指定しつつ、外部からアクセスしてみます。

ただし、デプロイ初回においてはロードバランサーにバックエンドを追加したりヘルスチェックをする過程があるため、正常な結果を返すまで数分かかります。クラウドのコンソールから、そのリソースの状況を確認することもできます。
<walkthrough-menu-navigation sectionId="LOAD_BALANCING_SECTION"></walkthrough-menu-navigation>

```bash
curl -iw'\n\n' -H "Host: app.dev.example.com" "http://${dev_ip}"
```

## 2. Cloud Deploy の利用

ここからは [Cloud Deploy](https://cloud.google.com/deploy?hl=ja) を使ったデプロイをみていきます。Cloud Deploy を使うと

- 複数の環境で同じ成果物を上手に共有して一貫性を確保したり
- デプロイの履歴だけでなく、その時の設定値やメタデータのバージョン管理したり
- サービスの検証やエラー発生時のロールバックを自動化したりできます

## 2.1. サービス アカウントへの権限付与

この例では複雑さを軽減するために、デフォルトのコンピューティング [サービス アカウント](https://cloud.google.com/iam/docs/service-account-overview?hl=ja) を使ったシンプルな構成を使います。（実運用においてはサービス アカウントを別途作成するなど、セキュリティを強化してください）

```bash
project_number="$(gcloud projects list --filter "${PROJECT_ID}" --format 'value(PROJECT_NUMBER)')"
gcloud projects add-iam-policy-binding "${PROJECT_ID}" --member "serviceAccount:${project_number}-compute@developer.gserviceaccount.com" --role "roles/clouddeploy.jobRunner"
gcloud projects add-iam-policy-binding "${PROJECT_ID}" --member "serviceAccount:${project_number}-compute@developer.gserviceaccount.com" --role "roles/container.developer"
```

## 2.2. Cloud Deploy パイプラインの作成

Cloud Deploy の[デリバリー パイプライン](https://cloud.google.com/deploy/docs/managing-delivery-pipeline?hl=ja)を作成します。いわゆる CI/CD のパイプラインではなく、ここでいうパイプラインは

- どんな環境に
- どんな設定を使って
- どんな順序でデプロイするか

を定義したものです。

<walkthrough-editor-open-file filePath="clouddeploy.yaml">実際の定義</walkthrough-editor-open-file>をみてみましょう。`stages` に書かれているものがデプロイ先の環境であり、上から順にデプロイされていきます。環境ごとに設定を変えたい場合は `profiles` を使って指定できます。

以下のコマンドで、この定義をベースに Cloud Deploy のパイプラインを作成します。

```bash
sed -i "s/GOOGLE_CLOUD_PROJECT_ID/${PROJECT_ID}/g" clouddeploy.yaml
gcloud deploy apply --file clouddeploy.yaml --region "${GOOGLE_CLOUD_REGION}"
```

作成したパイプラインをコンソールで確認してみましょう。
<walkthrough-menu-navigation sectionId="CLOUD_DEPLOY_SECTION"></walkthrough-menu-navigation>

## 2.3. コンテナのビルドとプッシュ

CI/CD を実施するために、先ほど `skaffold run` で一気に行なってしまったことを分解していきます。

まずは [build コマンド](https://skaffold.dev/docs/references/cli/#skaffold-build)を使い、改めてコンテナのイメージをビルド、リポジトリへプッシュします。

```bash
skaffold build -p prod --default-repo "${GOOGLE_CLOUD_REGION}-docker.pkg.dev/${PROJECT_ID}/canary-apps" --push --file-output build.out
```

結果として出力される `build.out` にはコンテナ イメージのタグ情報などが格納されます。

```bash
cat build.out | jq
```

## 2.4. Cloud Deploy パイプラインの実行

次に **同時にデプロイしたい成果物の集合** である Cloud Deploy の `リリース` を作成します。リリースの ID を決めて

```text
git_hash=$(git rev-parse --short HEAD 2>/dev/null)
labels="commit-id=${git_hash:-unknown},release-date=$(date -u '+%Y%m%d')"
release_version="git-${git_hash:-unknown}-$(date -u '+%Y%m%d-%H%M%S')"
echo -e "\n${release_version}\n"
```

以下のコマンドを実行すると、リリースの作成に加え、パイプラインの最初のターゲットへのデプロイも始まります。

```bash
gcloud deploy releases create "${release_version}" --region "${GOOGLE_CLOUD_REGION}" --delivery-pipeline "canary-sample" --build-artifacts build.out --annotations "author=$(whoami)" --labels "${labels}"
```

デプロイされるまで 1 分ほど待ちます。コンソールで、dev 環境が緑色になれば OK です。結果をチェックしてみましょう。

```bash
curl -iw'\n\n' -H "Host: app.dev.example.com" "http://${dev_ip}/release"
```

<walkthrough-editor-select-line filePath="src/main.go" startLine="36" endLine="36" startCharacterOffset="0" endCharacterOffset="100">/release</walkthrough-editor-select-line> という API にアクセスし、Skaffold でレンダリングされるマニフェストに `Label` として付与された<walkthrough-editor-select-line filePath="base/deployment.yaml" startLine="37" endLine="40" startCharacterOffset="0" endCharacterOffset="100">Cloud Deploy のリリース ID</walkthrough-editor-select-line> を返しています。

## 2.5. prod 環境へのデプロイ

dev 環境での動作が確認できたので、[作成したリリースを prod 環境に昇格](https://cloud.google.com/deploy/docs/promote-release?hl=ja)します。これにより

- 成果物はまったく同じものを利用しながらも
- Profile として差分管理された設定を使って
- prod 環境にアプリケーションがデプロイされます

```bash
gcloud deploy releases promote --release "${release_version}" --region "${GOOGLE_CLOUD_REGION}" --delivery-pipeline "canary-sample" --quiet
```

デプロイされていく様子や、その結果をコンソールで確認してみてください。
<walkthrough-menu-navigation sectionId="CLOUD_DEPLOY_SECTION"></walkthrough-menu-navigation>

dev 環境同様、負荷分散側のリソースが構成されるまで 3 分ほど待ちます。その後結果をチェックしてみましょう。

```bash
curl -iw'\n\n' -H "Host: app.prod.example.com" "http://${prod_ip}/release"
```

## 3. カナリア デプロイ

さきほどの prod 環境へのデプロイは最初のバージョンだったため、カナリア デプロイは起きませんでした。コードを一部修正し、次のバージョンをデプロイすることでその挙動を確認してみます。

## 3.1. 新しいリリースの作成

メッセージを `Hello GKE!` に変更し

```bash
sed -i "s/World/GKE/g" src/main.go
```

新しいリリース ID を作り

```text
git_hash=$(git rev-parse --short HEAD 2>/dev/null)
release_version="git-${git_hash:-unknown}-$(date -u '+%Y%m%d-%H%M%S')"
echo -e "\n${release_version}\n"
```

リリースを作成、dev 環境へデプロイします。

```bash
gcloud deploy releases create "${release_version}" --region "${GOOGLE_CLOUD_REGION}" --delivery-pipeline "canary-sample" --build-artifacts build.out --annotations "author=$(whoami)" --labels "${labels}"
```

しばらくしたら dev 環境の様子を確認してみましょう。

```bash
curl -iw'\n\n' -H "Host: app.dev.example.com" "http://${dev_ip}/release"
```

## 3.2. カナリアの設定と履歴の確認

リリースを prod 環境へ昇格させます。

```bash
gcloud deploy releases promote --release "${release_version}" --region "${GOOGLE_CLOUD_REGION}" --delivery-pipeline "canary-sample" --quiet
```

prod 環境には Flagger による<walkthrough-editor-select-line filePath="bootstrap.yaml" startLine="71" endLine="72" startCharacterOffset="0" endCharacterOffset="100">カナリア デプロイの設定</walkthrough-editor-select-line>がされていますが、具体的には

- 1 分おきに `success-rate` を確認し
- 成功率が 90% を超えていたら流量を 10% 新しいバージョンに流し
- 50% まで達したら、次は全量を新しいバージョンに流す
- 一方で 1 分以内に成功しなければロールバック

という設定になっています。では実際に状況を確認してみましょう。

```bash
kubectl -n prod describe canary/app
```

Events が少しずつ増えている様子がみてとれるかと思います。実際に応答が変化していく様子も確認してみてください。

```bash
curl -iw'\n\n' -H "Host: app.prod.example.com" "http://${prod_ip}/release"
```

## 4. クリーンアップ

ハンズオンに利用したプロジェクトを削除し、課金を止めます。

```bash
gcloud projects delete ${PROJECT_ID}
```

プロジェクトがそのまま消せない場合は、リソースを個別に削除します。まず外部 HTTP(S) ロードバランサを削除します。

```bash
kubectl delete httproutes --all --all-namespaces
kubectl delete -f bootstrap.yaml
```

Cloud Deploy や GKE そのものを削除します。

```bash
gcloud deploy delivery-pipelines delete canary-sample --force --region "${GOOGLE_CLOUD_REGION}" --quiet
gcloud artifacts repositories delete canary-apps --location "${GOOGLE_CLOUD_REGION}" --quiet
gcloud container clusters delete canary-cluster --zone "${GOOGLE_CLOUD_ZONE}" --quiet
```

Cloud Storage についても以下を参考に削除してください。新規のプロジェクトであればすべて削除して問題ありませんが、そうでない場合は削除対象にお気をつけください。

```bash
gcloud storage buckets list --format json | jq -r '.[].id'
gcloud storage rm -r gs://<id>
```

## これで終わりです

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

すべて完了しました。
