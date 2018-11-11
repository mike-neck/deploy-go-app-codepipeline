# deploy-go-app-codepipeline

これも golang アプリケーションを codepipeline で lambda にデプロイするサンプル

ビルド
---

```bash
make build
```

すればよい

準備
---

### CloudFormation 用のロール

CodePipeline でパイプラインを作るために、最初にデプロイで使う CloudFormation 用のサービスロールを作る 


1. IAM の画面にて ロールの作成 > AWS サービス > CloudFormation を選択して、ポリシーを付与せずに作成する.
1. (1.)で作成したロールに次のコマンドで出力される json をインラインポリシーとして付与する.
```bash
sed 's/__REGION__/ap-northeast-1/g' cloud-formation-role.json | sed 's/__ACCOUNT_ID__/'$(aws sts get-caller-identity --query 'Account' --output text)'/g'
```

### アーティファクト用の S3 バケット

ビルドした

1. ビルドしたアーティファクトを保存するための S3 バケットを作成する
1. これを Systems Manager のパラメーターストアに入れておき、ビルド時に解決できるようにする
    * Systems Manager のパラメーターストアに (1.) で作成したバケットの名前を `deploy-lambda-example-bucket` という名前で登録する

CodePipeline
---

1. CodePipeline コンソールを開く
1. パイプライン作成をクリックして次の値を入力する
    * パイプライン名 : 任意
    * サービスロール : 新しいサービスロール
    * ロール名 : 任意
    * **Allo AWS CodePipeline...** にチェック
    * アーティファクトストア : デフォルトの場所
1. 次へをクリックしてソースにて次の値を入力する
    * ソースプロバイダー : GitHub
    * GitHub に接続
    * リポジトリー : このリポジトリー
    * ブランチ : master
    * 変更検出オプション : GitHub WebHook
1. 次へをクリックしてビルドに次の値を入力する
    * ビルドプロバイダ : AWS CodeBuild
    * AWS CodeBuild でプロジェクトを作る必要があるので、 **Create Project** をクリックする
1. AWS CodeBuild のプロジェクトを以下の要領で作成する
    * プロジェクト名 : 任意
    * 説明 : 任意
    * イメージ : マネージド型イメージ
    * オペレーティングシステム : Ubuntu
    * ランタイム : Golang
    * ランタイムバージョン : aws/codebuild/golang:1.10
    * イメージバージョン : Always use the latest image for this runtime version
    * 特権付与 : チェックしない
    * サービスロール : 新しいサービスロール
    * ロール名 : 任意
    * Additional configuration : 特に設定しない
    * Buildspec: buildspec.yml
    * **Continue to CodePipeline** にて CodeBuild のプロジェクトを作成
1. CodePipeline に戻り次へをクリックしてデプロイを次の通り入力する. なお、この段階で CloudFormation の stack が作られていないので、 ValidationError が発生する
    * デプロイプロバイダ : AWS CloudFormation
    * アクションモード : 変更セットの作成または置換
    * スタックの名前 : deploy-go-app-stack
    * 変更セット名 : deploy-go-app-cs
    * テンプレート : BuildArtifact::deploy.yml (この段階でビルドアーティファクトの名前がわからないのに入力させるの不親切)
    * テンプレート設定 : 空白
    * 機能 : CAPABILITY_IAM
    * ロール名 : 最初に作成したサービスロール
1. パイプラインを作ると自動でビルドが走るが以下の原因により成功することはない
    * CodeBuild が S3 にオブジェクトを置けない
    * CodeBuild がパラメーターストアの値を取得できない
    * Stack がない
    * Change Set がない
