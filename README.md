# Google Cloud CI/CD ハンズオン

ローカル開発 + CI + Google Cloud への CD を、[Cloud Code](https://cloud.google.com/code?hl=ja) ベースで実施するハンズオンです。

## [Cloud Build](https://cloud.google.com/build?hl=ja)

### Cloud Run 編

以下サービス・ソフトウェアの組み合わせで、ローカル開発からデプロイまで。

- Cloud Build
- Cloud Code
- [Buildpacks](https://github.com/GoogleCloudPlatform/buildpacks)
- [Go](https://golang.org/)

1. 以下をクリックし、Cloud Shell 環境を起動してください。

[![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://console.cloud.google.com/home/dashboard?cloudshell=true)

2. 以下のコマンドをで実行してください。チュートリアルが始まります。

```bash
wget -qO tutorial.md https://raw.githubusercontent.com/google-cloud-japan/appdev-cicd-handson/main/cloud-build/cloud-run.md
teachme tutorial.md
```

### GKE (Kubernetes) 編

以下サービス・ソフトウェアの組み合わせで、ローカル開発からデプロイまで。

- Cloud Build
- Cloud Code
- [Skaffold](https://skaffold.dev/)
- [Kustomize](https://kustomize.io/)
- [Jib](https://github.com/GoogleContainerTools/jib)
- Java ([Micronaut](https://micronaut.io/)) 

1. 以下をクリックし、Cloud Shell 環境を起動してください。

[![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://console.cloud.google.com/home/dashboard?cloudshell=true)

2. 以下のコマンドをで実行してください。チュートリアルが始まります。

```bash
wget -qO tutorial.md https://raw.githubusercontent.com/google-cloud-japan/appdev-cicd-handson/main/cloud-build/kubernetes.md
teachme tutorial.md
```

## [Cloud Deploy](https://cloud.google.com/deploy?hl=ja)

以下サービス・ソフトウェアの組み合わせで CI/CD パイプラインを体験できます。

- Cloud Deploy
- GitHub Actions
- [Skaffold](https://skaffold.dev/)
- [Kustomize](https://kustomize.io/)
- [Dart](https://dart.dev/)

### Cloud Deploy にフォーカス編

（Cloud Deploy が内部的に利用する [Skaffold](https://skaffold.dev/) は Kubernetes のローカル開発環境として大変優れた OSS なこともあり、Cloud Code との組み合わせはぜひ体験いただきたいものの）

ローカル開発は意識せず、**まずはざっと Cloud Deploy を動かしたい方** はこちら

- [Kubernetes のマニフェストをそのまま管理する例](https://github.com/google-cloud-japan/appdev-cicd-handson/tree/main/cloud-deploy/sample-resources/default)
- [Kustomize で環境差異を管理する例](https://github.com/google-cloud-japan/appdev-cicd-handson/tree/main/cloud-deploy/sample-resources/kustomize)

### ローカル開発からデプロイまでの一貫したテスト・ビルド体験編

1. 以下をクリックし、Cloud Shell 環境を起動してください。

[![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://console.cloud.google.com/home/dashboard?cloudshell=true)

2. 以下のコマンドをで実行してください。チュートリアルが始まります。

```bash
wget -qO tutorial.md https://raw.githubusercontent.com/google-cloud-japan/appdev-cicd-handson/main/cloud-deploy/basic.md
teachme tutorial.md
```

## その他

### Skaffold 応用編

- [複数の成果物をまとめる例](https://github.com/google-cloud-japan/appdev-cicd-handson/tree/main/others/sample-resources/multi-apps)
