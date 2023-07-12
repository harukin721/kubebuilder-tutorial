# カスタムコントローラーの動作確認

ref.https://zoetrope.github.io/kubebuilder-training/kubebuilder/kind.html

## kindコマンドを利用してKubernetesクラスターを作成

```
$ kind create cluster
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.26.3) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Not sure what to do next? 😅  Check out https://kind.sigs.k8s.io/docs/user/quick-start/

```

## Webhook用の証明書を発行するためにcert-managerのデプロイ

```
$ kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/latest/download/cert-manager.yaml
namespace/cert-manager created
customresourcedefinition.apiextensions.k8s.io/certificaterequests.cert-manager.io created
customresourcedefinition.apiextensions.k8s.io/certificates.cert-manager.io created
customresourcedefinition.apiextensions.k8s.io/challenges.acme.cert-manager.io created
customresourcedefinition.apiextensions.k8s.io/clusterissuers.cert-manager.io created
customresourcedefinition.apiextensions.k8s.io/issuers.cert-manager.io created
customresourcedefinition.apiextensions.k8s.io/orders.acme.cert-manager.io created
serviceaccount/cert-manager-cainjector created
serviceaccount/cert-manager created
serviceaccount/cert-manager-webhook created
configmap/cert-manager-webhook created
clusterrole.rbac.authorization.k8s.io/cert-manager-cainjector created
clusterrole.rbac.authorization.k8s.io/cert-manager-controller-issuers created
clusterrole.rbac.authorization.k8s.io/cert-manager-controller-clusterissuers created
clusterrole.rbac.authorization.k8s.io/cert-manager-controller-certificates created
clusterrole.rbac.authorization.k8s.io/cert-manager-controller-orders created
clusterrole.rbac.authorization.k8s.io/cert-manager-controller-challenges created
clusterrole.rbac.authorization.k8s.io/cert-manager-controller-ingress-shim created
clusterrole.rbac.authorization.k8s.io/cert-manager-view created
clusterrole.rbac.authorization.k8s.io/cert-manager-edit created
clusterrole.rbac.authorization.k8s.io/cert-manager-controller-approve:cert-manager-io created
clusterrole.rbac.authorization.k8s.io/cert-manager-controller-certificatesigningrequests created
clusterrole.rbac.authorization.k8s.io/cert-manager-webhook:subjectaccessreviews created
clusterrolebinding.rbac.authorization.k8s.io/cert-manager-cainjector created
clusterrolebinding.rbac.authorization.k8s.io/cert-manager-controller-issuers created
clusterrolebinding.rbac.authorization.k8s.io/cert-manager-controller-clusterissuers created
clusterrolebinding.rbac.authorization.k8s.io/cert-manager-controller-certificates created
clusterrolebinding.rbac.authorization.k8s.io/cert-manager-controller-orders created
clusterrolebinding.rbac.authorization.k8s.io/cert-manager-controller-challenges created
clusterrolebinding.rbac.authorization.k8s.io/cert-manager-controller-ingress-shim created
clusterrolebinding.rbac.authorization.k8s.io/cert-manager-controller-approve:cert-manager-io created
clusterrolebinding.rbac.authorization.k8s.io/cert-manager-controller-certificatesigningrequests created
clusterrolebinding.rbac.authorization.k8s.io/cert-manager-webhook:subjectaccessreviews created
role.rbac.authorization.k8s.io/cert-manager-cainjector:leaderelection created
role.rbac.authorization.k8s.io/cert-manager:leaderelection created
role.rbac.authorization.k8s.io/cert-manager-webhook:dynamic-serving created
rolebinding.rbac.authorization.k8s.io/cert-manager-cainjector:leaderelection created
rolebinding.rbac.authorization.k8s.io/cert-manager:leaderelection created
rolebinding.rbac.authorization.k8s.io/cert-manager-webhook:dynamic-serving created
service/cert-manager created
service/cert-manager-webhook created
deployment.apps/cert-manager-cainjector created
deployment.apps/cert-manager created
deployment.apps/cert-manager-webhook created
mutatingwebhookconfiguration.admissionregistration.k8s.io/cert-manager-webhook created
validatingwebhookconfiguration.admissionregistration.k8s.io/cert-manager-webhook created

$ k get po -n cert-manager
NAME                                       READY   STATUS    RESTARTS   AGE
cert-manager-5d77b478-fb8lk                1/1     Running   0          3m40s
cert-manager-cainjector-576655b654-s86r7   1/1     Running   0          3m40s
cert-manager-webhook-795dc979b6-f9twl      1/1     Running   0          3m39s
```

## コンテナイメージをビルド

```
$ make docker-build
/Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
test -s /Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/setup-envtest || GOBIN=/Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
go: downloading sigs.k8s.io/controller-runtime/tools/setup-envtest v0.0.0-20230710230511-56a973c46efe
KUBEBUILDER_ASSETS="/Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/k8s/1.26.0-darwin-arm64" go test ./... -coverprofile cover.out
?   	github.com/harukin721/markdown-view	[no test files]
ok  	github.com/harukin721/markdown-view/api/v1	1.138s	coverage: 2.0% of statements
ok  	github.com/harukin721/markdown-view/controllers	0.792s	coverage: 0.0% of statements
docker build -t controller:latest .
[+] Building 160.4s (17/17) FINISHED
 => [internal] load build definition from Dockerfile                                                                                                                                                             0.1s
 => => transferring dockerfile: 1.29kB                                                                                                                                                                           0.0s
 => [internal] load .dockerignore                                                                                                                                                                                0.0s
 => => transferring context: 171B                                                                                                                                                                                0.0s
 => [internal] load metadata for gcr.io/distroless/static:nonroot                                                                                                                                                2.4s
 => [internal] load metadata for docker.io/library/golang:1.19                                                                                                                                                   2.9s
 => [builder 1/9] FROM docker.io/library/golang:1.19@sha256:031338ed4f8477cfa39ac084317f3b5a45d21018279c5597c19a6cb0113e2e18                                                                                   117.8s
 => => resolve docker.io/library/golang:1.19@sha256:031338ed4f8477cfa39ac084317f3b5a45d21018279c5597c19a6cb0113e2e18                                                                                             0.0s
 => => sha256:aaf0f990abf1f411b9c647f08913f64ac9e5dae9eee24868e353beadd4d0ea1c 1.58kB / 1.58kB                                                                                                                   0.0s
 => => sha256:031338ed4f8477cfa39ac084317f3b5a45d21018279c5597c19a6cb0113e2e18 2.36kB / 2.36kB                                                                                                                   0.0s
 => => sha256:427f81f57ee6194897eac6ff466f292e9fc00f512a0cfa7cd849caed4bd662b8 6.88kB / 6.88kB                                                                                                                   0.0s
 => => sha256:42cbebf8bc116ba1aed7897e2d0566bf49da9d5c70be71b6a7cb07805d2f5b7a 49.57MB / 49.57MB                                                                                                                51.6s
 => => sha256:9a0518ec57568a70561f7c04650f9554c88dada973f54d88e36f65b0796d35de 23.57MB / 23.57MB                                                                                                                25.3s
 => => sha256:356172c718acf9930d9465b170864319079e2d2ebac0ddef781d64e85789531e 63.98MB / 63.98MB                                                                                                                59.4s
 => => sha256:a3c1d40c82551fded3ae8435595284e568a5c425b275785f3a78b95ed6f25b15 86.26MB / 86.26MB                                                                                                               102.7s
 => => extracting sha256:42cbebf8bc116ba1aed7897e2d0566bf49da9d5c70be71b6a7cb07805d2f5b7a                                                                                                                        1.3s
 => => sha256:3c9e9da25c4bd8d46fad9ed3a3f731161301dec43cb371c73db8da776d350b86 115.34MB / 115.34MB                                                                                                             114.9s
 => => extracting sha256:9a0518ec57568a70561f7c04650f9554c88dada973f54d88e36f65b0796d35de                                                                                                                        0.5s
 => => extracting sha256:356172c718acf9930d9465b170864319079e2d2ebac0ddef781d64e85789531e                                                                                                                        1.6s
 => => sha256:524bc1e999144800bd0d20cb27e6b8216c0a7ecfe39d256a55e1ebdce669c388 156B / 156B                                                                                                                      72.8s
 => => extracting sha256:a3c1d40c82551fded3ae8435595284e568a5c425b275785f3a78b95ed6f25b15                                                                                                                        2.2s
 => => extracting sha256:3c9e9da25c4bd8d46fad9ed3a3f731161301dec43cb371c73db8da776d350b86                                                                                                                        2.7s
 => => extracting sha256:524bc1e999144800bd0d20cb27e6b8216c0a7ecfe39d256a55e1ebdce669c388                                                                                                                        0.0s
 => [stage-1 1/3] FROM gcr.io/distroless/static:nonroot@sha256:9ecc53c269509f63c69a266168e4a687c7eb8c0cfd753bd8bfcaa4f58a90876f                                                                                 73.5s
 => => resolve gcr.io/distroless/static:nonroot@sha256:9ecc53c269509f63c69a266168e4a687c7eb8c0cfd753bd8bfcaa4f58a90876f                                                                                          0.0s
 => => sha256:9ecc53c269509f63c69a266168e4a687c7eb8c0cfd753bd8bfcaa4f58a90876f 1.46kB / 1.46kB                                                                                                                   0.0s
 => => sha256:6fd06868c905eecb495c751f1e6bb1cda4ce45dcb2ae16ef63b9de4d6532cfb4 1.28kB / 1.28kB                                                                                                                   0.0s
 => => sha256:0b41f743fd4d78cb50ba86dd3b951b51458744109e1f5063a76bc5a792c3d8e7 103.73kB / 103.73kB                                                                                                               1.1s
 => => sha256:fe5ca62666f04366c8e7f605aa82997d71320183e99962fa76b3209fdfbb8b58 21.20kB / 21.20kB                                                                                                                 0.3s
 => => sha256:b02a7525f878e61fc1ef8a7405a2cc17f866e8de222c1c98fd6681aff6e509db 716.49kB / 716.49kB                                                                                                               1.6s
 => => sha256:05810557ec4b4bf01f4df548c06cc915bb29d81cb339495fe1ad2e668434bf8e 1.65kB / 1.65kB                                                                                                                   0.0s
 => => sha256:fcb6f6d2c9986d9cd6a2ea3cc2936e5fc613e09f1af9042329011e43057f3265 317B / 317B                                                                                                                       0.9s
 => => sha256:e8c73c638ae9ec5ad70c49df7e484040d889cca6b4a9af056579c3d058ea93f0 198B / 198B                                                                                                                       1.2s
 => => extracting sha256:0b41f743fd4d78cb50ba86dd3b951b51458744109e1f5063a76bc5a792c3d8e7                                                                                                                        0.1s
 => => sha256:1e3d9b7d145208fa8fa3ee1c9612d0adaac7255f1bbc9ddea7e461e0b317805c 113B / 113B                                                                                                                       1.6s
 => => extracting sha256:fe5ca62666f04366c8e7f605aa82997d71320183e99962fa76b3209fdfbb8b58                                                                                                                        0.0s
 => => extracting sha256:b02a7525f878e61fc1ef8a7405a2cc17f866e8de222c1c98fd6681aff6e509db                                                                                                                        0.1s
 => => extracting sha256:fcb6f6d2c9986d9cd6a2ea3cc2936e5fc613e09f1af9042329011e43057f3265                                                                                                                        0.0s
 => => extracting sha256:e8c73c638ae9ec5ad70c49df7e484040d889cca6b4a9af056579c3d058ea93f0                                                                                                                        0.0s
 => => extracting sha256:1e3d9b7d145208fa8fa3ee1c9612d0adaac7255f1bbc9ddea7e461e0b317805c                                                                                                                        0.0s
 => => sha256:4aa0ea1413d37a58615488592a0b827ea4b2e48fa5a77cf707d0e35f025e613f 385B / 385B                                                                                                                      71.8s
 => => extracting sha256:4aa0ea1413d37a58615488592a0b827ea4b2e48fa5a77cf707d0e35f025e613f                                                                                                                        0.0s
 => => sha256:7c881f9ab25e0d86562a123b5fb56aebf8aa0ddd7d48ef602faf8d1e7cf43d8c 355B / 355B                                                                                                                      72.1s
 => => extracting sha256:7c881f9ab25e0d86562a123b5fb56aebf8aa0ddd7d48ef602faf8d1e7cf43d8c                                                                                                                        0.0s
 => => sha256:5627a970d25e752d971a501ec7e35d0d6fdcd4a3ce9e958715a686853024794a 130.56kB / 130.56kB                                                                                                              73.3s
 => => extracting sha256:5627a970d25e752d971a501ec7e35d0d6fdcd4a3ce9e958715a686853024794a                                                                                                                        0.0s
 => [internal] load build context                                                                                                                                                                                0.1s
 => => transferring context: 85.37kB                                                                                                                                                                             0.1s
 => [builder 2/9] WORKDIR /workspace                                                                                                                                                                             0.6s
 => [builder 3/9] COPY go.mod go.mod                                                                                                                                                                             0.0s
 => [builder 4/9] COPY go.sum go.sum                                                                                                                                                                             0.0s
 => [builder 5/9] RUN go mod download                                                                                                                                                                           15.1s
 => [builder 6/9] COPY main.go main.go                                                                                                                                                                           0.0s
 => [builder 7/9] COPY api/ api/                                                                                                                                                                                 0.0s
 => [builder 8/9] COPY controllers/ controllers/                                                                                                                                                                 0.0s
 => [builder 9/9] RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o manager main.go                                                                                                                      23.5s
 => [stage-1 2/3] COPY --from=builder /workspace/manager .                                                                                                                                                       0.1s
 => exporting to image                                                                                                                                                                                           0.1s
 => => exporting layers                                                                                                                                                                                          0.1s
 => => writing image sha256:b761a3ec3552bd49c5a7352835ef5b97f5d17b76359d58fe42dff7f37502ccf0                                                                                                                     0.0s
 => => naming to docker.io/library/controller:latest                                                                                                                                                             0.0s

```

## kind環境にコンテナイメージをロード

```
$ kind load docker-image controller:latest
Image: "controller:latest" with ID "sha256:b761a3ec3552bd49c5a7352835ef5b97f5d17b76359d58fe42dff7f37502ccf0" not yet present on node "kind-control-plane", loading...
```

## `config/manager/manager.yaml` に `imagePullPolicy: IfNotPresent` 追加

- コンテナイメージのタグ名にlatestを指定した場合、ImagePullPolicyがAlwaysになり、ロードしたコンテナイメージが利用されない場合がある

コミット. 06d48b9fd67a41de7851c780deea1b0156dc9df8

# コントローラーの動作確認

```
$ test -s /Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/kustomize || { curl -Ss "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash -s -- 3.8.7 /Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin; }
Version v3.8.7 does not exist for darwin/arm64, trying darwin/amd64 instead.
{Version:kustomize/v3.8.7 GitCommit:ad092cc7a91c07fdf63a2e4b7f13fa588a39af4f BuildDate:2020-11-11T23:19:38Z GoOs:darwin GoArch:amd64}
kustomize installed to /Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/kustomize
```

## CRDをKubernetesクラスターに適用

```
$ make install
test -s /Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/controller-gen && /Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/controller-gen --version | grep -q v0.11.1 || \
	GOBIN=/Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.11.1
/Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/markdownviews.view.harukin721.github.io created
```

## 各種マニフェストを適用

```
$ make deploy
test -s /Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/controller-gen && /Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/controller-gen --version | grep -q v0.11.1 || \
	GOBIN=/Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.11.1
/Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
cd config/manager && /Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/kustomize edit set image controller=controller:latest
/Users/harukin/ghq/github.com/harukin721/kubebuilder-tutorial/markdown-view/bin/kustomize build config/default | kubectl apply -f -
namespace/markdown-view-system created
customresourcedefinition.apiextensions.k8s.io/markdownviews.view.harukin721.github.io configured
serviceaccount/markdown-view-controller-manager created
role.rbac.authorization.k8s.io/markdown-view-leader-election-role created
clusterrole.rbac.authorization.k8s.io/markdown-view-manager-role created
clusterrole.rbac.authorization.k8s.io/markdown-view-metrics-reader created
clusterrole.rbac.authorization.k8s.io/markdown-view-proxy-role created
rolebinding.rbac.authorization.k8s.io/markdown-view-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/markdown-view-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/markdown-view-proxy-rolebinding created
service/markdown-view-controller-manager-metrics-service created
service/markdown-view-webhook-service created
deployment.apps/markdown-view-controller-manager created
certificate.cert-manager.io/markdown-view-serving-cert created
issuer.cert-manager.io/markdown-view-selfsigned-issuer created
mutatingwebhookconfiguration.admissionregistration.k8s.io/markdown-view-mutating-webhook-configuration created
validatingwebhookconfiguration.admissionregistration.k8s.io/markdown-view-validating-webhook-configuration created
```

## コントローラーのPodがRunningになっていることを確認

```
$ k get po -n markdown-view-system
NAME                                                READY   STATUS    RESTARTS   AGE
markdown-view-controller-manager-656694f5b4-65q58   2/2     Running   0          9m27s
```

## コントローラーのログを表示

```
$ k logs  markdown-view-controller-manager-656694f5b4-65q58 -n markdown-view-system -c manager -f
2023-07-12T14:24:42Z	INFO	controller-runtime.metrics	Metrics server is starting to listen	{"addr": "127.0.0.1:8080"}
2023-07-12T14:24:42Z	INFO	controller-runtime.builder	Registering a mutating webhook	{"GVK": "view.harukin721.github.io/v1, Kind=MarkdownView", "path": "/mutate-view-harukin721-github-io-v1-markdownview"}
2023-07-12T14:24:42Z	INFO	controller-runtime.webhook	Registering webhook	{"path": "/mutate-view-harukin721-github-io-v1-markdownview"}
2023-07-12T14:24:42Z	INFO	controller-runtime.builder	Registering a validating webhook	{"GVK": "view.harukin721.github.io/v1, Kind=MarkdownView", "path": "/validate-view-harukin721-github-io-v1-markdownview"}
2023-07-12T14:24:42Z	INFO	controller-runtime.webhook	Registering webhook	{"path": "/validate-view-harukin721-github-io-v1-markdownview"}
2023-07-12T14:24:42Z	INFO	setup	starting manager
2023-07-12T14:24:42Z	INFO	Starting server	{"path": "/metrics", "kind": "metrics", "addr": "127.0.0.1:8080"}
2023-07-12T14:24:42Z	INFO	Starting server	{"kind": "health probe", "addr": "[::]:8081"}
2023-07-12T14:24:42Z	INFO	controller-runtime.webhook.webhooks	Starting webhook server
2023-07-12T14:24:42Z	INFO	controller-runtime.certwatcher	Updated current TLS certificate
I0712 14:24:42.187996       1 leaderelection.go:248] attempting to acquire leader lease markdown-view-system/254c1180.harukin721.github.io...
2023-07-12T14:24:42Z	INFO	controller-runtime.webhook	Serving webhook server	{"host": "", "port": 9443}
2023-07-12T14:24:42Z	INFO	controller-runtime.certwatcher	Starting certificate watcher
I0712 14:24:42.196061       1 leaderelection.go:258] successfully acquired lease markdown-view-system/254c1180.harukin721.github.io
2023-07-12T14:24:42Z	DEBUG	events	markdown-view-controller-manager-656694f5b4-65q58_74712664-c12e-4eaa-a362-80fd51177c5c became leader	{"type": "Normal", "object": {"kind":"Lease","namespace":"markdown-view-system","name":"254c1180.harukin721.github.io","uid":"fe80139e-5192-4bc5-87d4-d65a31d874fe","apiVersion":"coordination.k8s.io/v1","resourceVersion":"3288"}, "reason": "LeaderElection"}
2023-07-12T14:24:42Z	INFO	Starting EventSource	{"controller": "markdownview", "controllerGroup": "view.harukin721.github.io", "controllerKind": "MarkdownView", "source": "kind source: *v1.MarkdownView"}
2023-07-12T14:24:42Z	INFO	Starting Controller	{"controller": "markdownview", "controllerGroup": "view.harukin721.github.io", "controllerKind": "MarkdownView"}
2023-07-12T14:24:42Z	INFO	Starting workers	{"controller": "markdownview", "controllerGroup": "view.harukin721.github.io", "controllerKind": "MarkdownView", "worker count": 1}
^C
```

## サンプルのカスタムリソースを適用

```
$ k apply -f config/samples/view_v1_markdownview.yaml
markdownview.view.harukin721.github.io/markdownview-sample created
```
