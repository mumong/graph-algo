# README.md
- [English](README.en.md)
- [Chinese](README.md)
- [Japanses](README.jp.md)


# README.md
- [English](README.en.md)
- [Chinese](README.md)
- [Japanese](README.jp.md)

# 🚀 グラフ検索アルゴリズムに基づく K8s リソーススケジューラ

## 📌 プロジェクト概要
本プロジェクトは、**グラフ検索アルゴリズム**に基づく **Kubernetes リソーススケジューリングシステム** です。**ダイクストラアルゴリズム** と **フロイドアルゴリズム** を使用して **クラスタリソースの探索** を行います。  
プロジェクトでは、**クラスタ内のリソースをグラフのノード** として扱い、**リソース間の関係をグラフのエッジ** として捉え、**有向非巡回グラフ（DAG）** を形成して最適なリソース割り当てを実現します。

### **サポートするリソースタイプ**
- **FPGA**
- **GPU**
- **DPU**
- **カスタムリソース（例: Cola）**

---

## 🎯 プロジェクトのコアロジック
1. **リソースリクエスト** → `curl` を通じてコンテナおよびリソース要件を送信  
2. **リソースグラフの構築** → Kubernetes リソースを解析し、**グラフ構造** を構築  
3. **パス計算** → **ダイクストラアルゴリズム** を使用して **最適なリソースパス** を計算  
4. **リソース割り当て** → 計算された **最適な組み合わせ** を割り当て  

---

## 🛠 API インターフェース
### **最適なリソース割り当ての取得**
**リクエスト方式**：`POST`  
**リクエストパス**：`/getResource/master`  

**リクエスト例**
```bash
curl -X POST http://localhost:8080/getResource/master \
     -H "Content-Type: application/json" \
     -d '[{"name": "container1", "resourceQuantity": {"nvidia.com/gpu": 1, "fpga": 1, "myway5.com/cola": 4}}, 
          {"name": "container2", "resourceQuantity": {"nvidia.com/gpu": 1, "fpga": 1, "myway5.com/cola": 2}}]'
```

最適な組み合わせは、主にさまざまな組み合わせのスコアを計算し、このグループのリソースを取得するための最小パスを表す最小スコアの組み合わせを計算する、ディジェシアン アルゴリズムを通じて取得されます。実際の実行結果を下図に示します。

![ad7122589e06d03b865c84a28ff9e6c](https://github.com/user-attachments/assets/551ed513-10ff-4acc-9a93-ab38d8c3d1b7)

---
## 🛠️ 構築方法
1. `make docker` を実行して Docker イメージをコンパイルおよびビルドします。
2. `bash ./tools/upload.sh` を実行して Docker イメージをローカルの containerd リポジトリにアップロードします。
3. `kubectl apply -f deployment.yaml` を実行してアプリケーションをデプロイします。

