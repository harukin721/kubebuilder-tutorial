/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"math/rand"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	amidav1 "github.com/harukin721/kubebuilder-tutorial/amida/api/v1"
	"github.com/sirupsen/logrus"
)

// AmidaReconciler reconciles a Amida object
type AmidaReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=amida.harukin721.github.io,resources=amida,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=amida.harukin721.github.io,resources=amida/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=amida.harukin721.github.io,resources=amida/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Amida object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile

// Reconcile は Amida リソースの状態を監視する
func (r *AmidaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	amida := &amidav1.Amida{}                                     // Amida リソースのインスタンスを作成する
	if err := r.Get(ctx, req.NamespacedName, amida); err != nil { // Amida リソースを取得する
		logrus.Error(err, "unable to fetch Amida")
		return ctrl.Result{}, err // 取得できなかった場合はエラーを返す.再実行してくれる
	}
	// Amida リソースの selects が selectCount と同じ数になったら終了する
	logrus.Info("ここは通る")
	logrus.Info("amida.Spec.Amida.Selects")
	logrus.Info(amida.Spec.Amida.Selects)
	logrus.Info("amida.Spec.Amida.SelectCount")
	logrus.Info(amida.Spec.Amida.SelectCount)
	logrus.Info("amida.Spec.Amida.Results")
	logrus.Info(amida.Spec.Amida.Results)
	// amida リソースの Results が selectCount と同じ数になったら終了する
	if len(amida.Spec.Amida.Results) == amida.Spec.Amida.SelectCount {
		logrus.Info("もういい感じになったので終了する")
		return ctrl.Result{}, nil
	}

	// amida リソースの selects から selectCount の数だけ選択する
	results, err := SelectFromSelects(amida)
	// Go は基本的にこんな感じでエラーハンドリングをする
	if err != nil {
		logrus.Error(err, "unable to select from selects")
		return ctrl.Result{}, err
	}
	// amida リソースの results に結果を追加する
	amida.Spec.Amida.Results = results
	logrus.Info(amida.Spec.Amida.Results)
	// amida リソースを更新する
	if err := r.Update(ctx, amida); err != nil {
		logrus.Error(err, "unable to update Amida")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AmidaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// このコントローラーが監視するリソースを指定する
	return ctrl.NewControllerManagedBy(mgr).
		// Amida リソースを指定する
		For(&amidav1.Amida{}).
		// ログを出力する
		Complete(r)
}

// selectCount は選択する数
func SelectFromSelects(a *amidav1.Amida) ([]string, error) {
	// selects から selectCount の数だけ選択する
	results := []string{}
	for i := 0; i < a.Spec.Amida.SelectCount; i++ { // selectCount の数だけ繰り返す(2)
		selected := rand.Intn(len(a.Spec.Amida.Selects))          // selects の中からランダムに選択する
		results = append(results, a.Spec.Amida.Selects[selected]) // 選択したものを results に追加する
	}
	return results, nil
}

// k8s から amida リソースを取得する
// selects から selectCount の数を取得する
