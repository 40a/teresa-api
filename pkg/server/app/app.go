package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/luizalabs/teresa-api/models/storage"
	"github.com/luizalabs/teresa-api/pkg/server/auth"
	"github.com/luizalabs/teresa-api/pkg/server/slug"
	st "github.com/luizalabs/teresa-api/pkg/server/storage"
	"github.com/luizalabs/teresa-api/pkg/server/team"
	"github.com/luizalabs/teresa-api/pkg/server/teresa_errors"
)

type Operations interface {
	Create(user *storage.User, app *App) error
	Logs(user *storage.User, appName string, lines int64, follow bool) (io.ReadCloser, error)
	Info(user *storage.User, appName string) (*Info, error)
	TeamName(appName string) (string, error)
	Get(appName string) (*App, error)
	HasPermission(user *storage.User, appName string) bool
	SetEnv(user *storage.User, appName string, evs []*EnvVar) error
	UnsetEnv(user *storage.User, appName string, evs []string) error
}

type K8sOperations interface {
	NamespaceAnnotation(namespace, annotation string) (string, error)
	NamespaceLabel(namespace, label string) (string, error)
	PodList(namespace string) ([]*Pod, error)
	PodLogs(namespace, podName string, lines int64, follow bool) (io.ReadCloser, error)
	CreateNamespace(app *App, userEmail string) error
	CreateQuota(app *App) error
	CreateSecret(appName, secretName string, data map[string][]byte) error
	CreateAutoScale(app *App) error
	AddressList(namespace string) ([]*Address, error)
	Status(namespace string) (*Status, error)
	AutoScale(namespace string) (*AutoScale, error)
	Limits(namespace, name string) (*Limits, error)
	IsNotFound(err error) bool
	IsAlreadyExists(err error) bool
	SetNamespaceAnnotations(namespace string, annotations map[string]string) error
	DeleteDeployEnvVars(namespace, name string, evNames []string) error
	CreateOrUpdateDeployEnvVars(namespace, name string, evs []*EnvVar) error
}

type AppOperations struct {
	tops team.Operations
	kops K8sOperations
	st   st.Storage
}

const (
	limitsName       = "limits"
	TeresaAnnotation = "teresa.io/app"
	TeresaTeamLabel  = "teresa.io/team"
	TeresaLastUser   = "teresa.io/last-user"
)

func (ops *AppOperations) hasPerm(user *storage.User, team string) bool {
	teams, err := ops.tops.ListByUser(user.Email)
	if err != nil {
		return false
	}
	var found bool
	for _, t := range teams {
		if t.Name == team {
			found = true
			break
		}
	}
	return found
}

func (ops *AppOperations) HasPermission(user *storage.User, appName string) bool {
	teamName, err := ops.TeamName(appName)
	if err != nil {
		return false
	}
	return ops.hasPerm(user, teamName)
}

func (ops *AppOperations) Create(user *storage.User, app *App) error {
	if !ops.hasPerm(user, app.Team) {
		return auth.ErrPermissionDenied
	}

	if err := ops.kops.CreateNamespace(app, user.Email); err != nil {
		if ops.kops.IsAlreadyExists(err) {
			return ErrAlreadyExists
		}
		return teresa_errors.NewInternalServerError(err)
	}

	if err := ops.kops.CreateQuota(app); err != nil {
		return teresa_errors.NewInternalServerError(err)
	}

	secretName := ops.st.K8sSecretName()
	data := ops.st.AccessData()
	if err := ops.kops.CreateSecret(app.Name, secretName, data); err != nil {
		return teresa_errors.NewInternalServerError(err)
	}

	if err := ops.kops.CreateAutoScale(app); err != nil {
		return teresa_errors.NewInternalServerError(err)
	}

	return nil
}

func (ops *AppOperations) Logs(user *storage.User, appName string, lines int64, follow bool) (io.ReadCloser, error) {
	team, err := ops.kops.NamespaceLabel(appName, TeresaTeamLabel)
	if err != nil {
		if ops.kops.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, teresa_errors.NewInternalServerError(err)
	}

	if !ops.hasPerm(user, team) {
		return nil, auth.ErrPermissionDenied
	}

	pods, err := ops.kops.PodList(appName)
	if err != nil {
		return nil, teresa_errors.NewInternalServerError(err)
	}

	r, w := io.Pipe()
	var wg sync.WaitGroup
	for _, pod := range pods {
		wg.Add(1)
		go func(namespace, podName string) {
			defer wg.Done()

			logs, err := ops.kops.PodLogs(namespace, podName, lines, follow)
			if err != nil {
				log.WithError(err).Errorf("streaming logs from pod %s", podName)
				return
			}
			defer logs.Close()

			scanner := bufio.NewScanner(logs)
			for scanner.Scan() {
				fmt.Fprintf(w, "[%s] - %s\n", podName, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				log.WithError(err).Errorf("streaming logs from pod %s", podName)
			}
		}(appName, pod.Name)
	}
	go func() {
		wg.Wait()
		w.Close()
	}()

	return r, nil
}

func (ops *AppOperations) Info(user *storage.User, appName string) (*Info, error) {
	teamName, err := ops.TeamName(appName)
	if err != nil {
		return nil, err
	}

	if !ops.hasPerm(user, teamName) {
		return nil, auth.ErrPermissionDenied
	}

	appMeta, err := ops.Get(appName)
	if err != nil {
		return nil, err
	}

	addr, err := ops.kops.AddressList(appName)
	if err != nil {
		return nil, teresa_errors.NewInternalServerError(err)
	}

	stat, err := ops.kops.Status(appName)
	if err != nil {
		return nil, teresa_errors.NewInternalServerError(err)
	}

	as, err := ops.kops.AutoScale(appName)
	if err != nil {
		return nil, teresa_errors.NewInternalServerError(err)
	}

	lim, err := ops.kops.Limits(appName, limitsName)
	if err != nil {
		return nil, teresa_errors.NewInternalServerError(err)
	}

	info := &Info{
		Team:      teamName,
		Addresses: addr,
		Status:    stat,
		AutoScale: as,
		Limits:    lim,
		EnvVars:   appMeta.EnvVars,
	}
	return info, nil
}

func (ops *AppOperations) TeamName(appName string) (string, error) {
	teamName, err := ops.kops.NamespaceLabel(appName, TeresaTeamLabel)
	if err != nil {
		if ops.kops.IsNotFound(err) {
			return "", ErrNotFound
		}
		return "", teresa_errors.NewInternalServerError(err)
	}
	return teamName, nil
}

func (ops *AppOperations) Get(appName string) (*App, error) {
	an, err := ops.kops.NamespaceAnnotation(appName, TeresaAnnotation)
	if err != nil {
		if ops.kops.IsNotFound(err) {
			return nil, teresa_errors.New(ErrNotFound, err)
		}
		return nil, teresa_errors.NewInternalServerError(err)
	}
	a := new(App)
	if err := json.Unmarshal([]byte(an), a); err != nil {
		err = fmt.Errorf("unmarshal app failed: %v", err)
		return nil, teresa_errors.NewInternalServerError(err)
	}

	return a, nil
}

func (ops *AppOperations) checkPermAndGet(user *storage.User, appName string) (*App, error) {
	team, err := ops.TeamName(appName)
	if err != nil {
		return nil, err
	}

	if !ops.hasPerm(user, team) {
		return nil, auth.ErrPermissionDenied
	}

	return ops.Get(appName)
}

func (ops *AppOperations) saveApp(app *App, lastUser string) error {
	b, err := json.Marshal(app)
	if err != nil {
		return fmt.Errorf("marshal app failed: %v", err)
	}

	anMap := map[string]string{
		TeresaAnnotation: string(b),
		TeresaLastUser:   lastUser,
	}

	return ops.kops.SetNamespaceAnnotations(app.Name, anMap)
}

func (ops *AppOperations) SetEnv(user *storage.User, appName string, evs []*EnvVar) error {
	evNames := make([]string, len(evs))
	for i, _ := range evs {
		evNames[i] = evs[i].Key
	}
	if err := checkForProtectedEnvVars(evNames); err != nil {
		return err
	}

	app, err := ops.checkPermAndGet(user, appName)
	if err != nil {
		return err
	}

	setEnvVars(app, evs)

	if err := ops.saveApp(app, user.Name); err != nil {
		return teresa_errors.NewInternalServerError(err)
	}

	if err = ops.kops.CreateOrUpdateDeployEnvVars(appName, appName, evs); err != nil {
		if ops.kops.IsNotFound(err) {
			return nil
		}
		return teresa_errors.NewInternalServerError(err)
	}
	return nil
}

func (ops *AppOperations) UnsetEnv(user *storage.User, appName string, evNames []string) error {
	if err := checkForProtectedEnvVars(evNames); err != nil {
		return err
	}

	app, err := ops.checkPermAndGet(user, appName)
	if err != nil {
		return err
	}

	unsetEnvVars(app, evNames)

	if err := ops.saveApp(app, user.Name); err != nil {
		return teresa_errors.NewInternalServerError(err)
	}

	if err = ops.kops.DeleteDeployEnvVars(appName, appName, evNames); err != nil {
		if ops.kops.IsNotFound(err) {
			return nil
		}
		return teresa_errors.NewInternalServerError(err)
	}
	return nil
}

func checkForProtectedEnvVars(evsNames []string) error {
	for _, name := range slug.ProtectedEnvVars {
		for _, item := range evsNames {
			if name == item {
				return ErrProtectedEnvVar
			}
		}
	}
	return nil
}

func NewOperations(tops team.Operations, kops K8sOperations, st st.Storage) Operations {
	return &AppOperations{tops: tops, kops: kops, st: st}
}
