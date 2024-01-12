kubebuilder init --domain malek.dev --owner "Malek Zaag" --repo github.com/Malek-Zaag/MyNewOperator

kubebuilder create api --group watchers --version v1beta1 --kind Watcher

make manifests
