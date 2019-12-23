# go-gae-app-version-rotate
- GAEサービスの最新から指定した世代のバージョンを残し、それ以外の古いバージョンを消すGCFのコードです。
- GAEバージョンが210を超えてデプロイできなくなる事故を防ぐ目的で利用する想定です。
- GCPアーキテクチャとしては、Scheduler → Pub/Sub → GCF を想定しています。

## 注意事項
- あくまでサービス内のバージョンを削除するためのコードなので、サービス自体の削除には未対応です。
- GAEの仕様としてサービス内には必ずトラフィック割り当て済みのバージョンが存在し、そのバージョンは削除不可のため、keep_version_countを0とすることはできません。(リクエストとしては通りますが、GCF内でエラーとなります)
- トラフィック割り当て済みのバージョンが削除対象となった場合、そのバージョンの削除はスキップし、keep_version_countに達するまで処理を続行します。
- `分割トラフィック数 > keep_version_count` だった場合は分割トラフィックを全て残します。

## 環境構築
### 変数表
|変数名|説明|
|---|---|
|YOUR_MANAGE_PROJECT|バージョン管理プロジェクト|
|YOUR_TOPIC_NAME|バージョン管理プロジェクトに作成するPub/Subトピックの名前|
|YOUR_JOB_NAME|バージョン管理プロジェクトに作成するSchedulerジョブの名前|
|YOUR_TARGET_PROJECT|管理対象プロジェクト|
|YOUR_TARGET_PROJECT_ID|管理対象プロジェクトの名前|
|YOUR_TARGET_SERVICE_NAME|管理対象プロジェクトのGAEサービス名|

### バージョン管理プロジェクト
#### APIの有効化
- `App Engine Admin API`を有効にする。

#### Pub/Subトピック作成
```bash
gcloud --project=YOUR_MANAGE_PROJECT \
pubsub topics create YOUR_TOPIC_NAME
```

#### Functionsデプロイ
```bash
gcloud functions deploy GAEAppVersionRotate \
--project=YOUR_MANAGE_PROJECT \
--trigger-topic=YOUR_TOPIC_NAME \
--region=asia-northeast1 \
--runtime=go111 \
--env-vars-file=env.yaml
```

#### Cloud Schedulerの設定
```bash
gcloud --project=YOUR_MANAGE_PROJECT \
scheduler jobs create pubsub YOUR_JOB_NAME \
--schedule="0 */3 * * *" \
--topic=YOUR_TOPIC_NAME \
--time-zone="Asia/Tokyo" \
--message-body='{"project_id": "YOUR_TARGET_PROJECT_ID", "service_name": "YOUR_TARGET_SERVICE_NAME", "keep_version_count": 3}'
```

### 管理対象プロジェクト
#### IAMの追加
- `バージョン管理プロジェクト`の関数で指定されているサービスアカウントを`管理対象プロジェクト`のIAMに追加し、`App Engine サービス管理者`の役割を付与する。

## 動作確認
### 動作確認用GAEアプリのデプロイ
この[リポジトリ](https://github.com/qushot/gae-multi-deploy-service-version)を利用することで指定した個数のサービス&バージョンのデプロイが可能。

### 関数の動作確認
```bash
gcloud --project=YOUR_MANAGE_PROJECT \
pubsub topics publish YOUR_TOPIC_NAME \
--message '{"project_id": "YOUR_TARGET_PROJECT_ID", "service_name": "YOUR_TARGET_SERVICE_NAME", "keep_version_count": 3}'
```
