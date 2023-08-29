package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/udmire/observability-operator/pkg/apps/manifest"
	apps_v1 "k8s.io/api/apps/v1"
	autoscaling_v1 "k8s.io/api/autoscaling/v1"
	batch_v1 "k8s.io/api/batch/v1"
	core_v1 "k8s.io/api/core/v1"
	networking_v1 "k8s.io/api/networking/v1"
	rbac_v1 "k8s.io/api/rbac/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

type Generator struct {
	encoder runtime.Encoder
}

func New() *Generator {
	scheme := runtime.NewScheme()
	gv := core_v1.SchemeGroupVersion
	scheme.AddKnownTypes(gv, &apps_v1.Deployment{}, &apps_v1.DaemonSet{}, &apps_v1.StatefulSet{}, &apps_v1.ReplicaSet{})
	scheme.AddKnownTypes(gv, &batch_v1.Job{}, &batch_v1.CronJob{})
	scheme.AddKnownTypes(gv, &networking_v1.Ingress{})
	scheme.AddKnownTypes(gv, &autoscaling_v1.HorizontalPodAutoscaler{})
	scheme.AddKnownTypes(gv, &rbac_v1.ClusterRole{}, &rbac_v1.ClusterRoleBinding{}, &rbac_v1.Role{}, &rbac_v1.RoleBinding{})
	scheme.AddKnownTypes(gv, &core_v1.ConfigMap{}, &core_v1.Secret{}, &core_v1.ServiceAccount{}, &core_v1.Service{})
	codecs := serializer.NewCodecFactory(scheme)
	encoder := codecs.EncoderForVersion(json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme, scheme, json.SerializerOptions{
		Yaml: true, Pretty: true, Strict: false,
	}), gv)
	return &Generator{
		encoder: encoder,
	}
}

func (g *Generator) Generate(name, version string, manifest *manifest.AppManifests, path string) error {
	if manifest == nil {
		return fmt.Errorf("invalid manifests")
	}

	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("invalid path")
	}

	root, err := os.MkdirTemp(path, name)
	if err != nil {
		return fmt.Errorf("cannot write in given folder")
	}
	// TODO
	writeAppManifestsToGivenFolder(name, manifest, root)

	return nil
}

func writeAppManifestsToGivenFolder(app string, manifests *manifest.AppManifests, root string) error {
	writeManifestsToGivenFolder(app, &manifests.Manifests, root)
	for _, comp := range manifests.CompsMenifests {
		compDir := filepath.Join(root, comp.Name)
		_ = os.Mkdir(filepath.Join(root, comp.Name), os.ModeDir)
		writeManifestsToGivenFolder(comp.Name, &comp.Manifests, compDir)
		//TODO create workloads
	}
	return nil
}

func writeManifestsToGivenFolder(app string, manifests *manifest.Manifests, root string) error {
	//TODO create manifest
	return nil
}
