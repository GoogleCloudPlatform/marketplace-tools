// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/fsutil"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPkgSubcommandPruningAndAliases(t *testing.T) {
	root := newSilentRootCmd()

	var pkgCmd *cobra.Command
	for _, sub := range root.Commands() {
		if sub.Name() == "pkg" {
			pkgCmd = sub
			break
		}
	}
	require.NotNil(t, pkgCmd)
	assert.Contains(t, pkgCmd.Aliases, "package")

	require.Len(t, pkgCmd.Commands(), 1)
	assert.Equal(t, "get", pkgCmd.Commands()[0].Name())

	pruned := [][]string{
		{"pkg", "sync", "."},
		{"package", "sync", "."},
		{"pkg", "update", "."},
		{"pkg", "diff", "."},
		{"pkg", "fix", "."},
		{"pkg", "init", "."},
		{"pkg", "desc", "."},
	}
	for _, args := range pruned {
		c := newSilentRootCmd()
		c.SetArgs(args)
		err := c.Execute()
		require.Error(t, err, "expected pruned command %v to fail", args)
	}

	// Verify both `mpdev pkg get -` and `mpdev package get -` still work on benign input.
	stdinYAML := "apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: demo\n"
	for _, verb := range []string{"pkg", "package"} {
		dst := filepath.Join(t.TempDir(), verb+"-out")
		c := newSilentRootCmd()
		c.SetIn(strings.NewReader(stdinYAML))
		c.SetArgs([]string{verb, "get", "-", dst})
		require.NoError(t, c.Execute())
		_, err := os.Stat(filepath.Join(dst, "demo_configmap.yaml"))
		require.NoError(t, err)
	}
}

func TestPkgGetCLIBlocksInjectionLocalTransportsAndSymlinkDest(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "out")
	victimDir := t.TempDir()
	symlinkDst := filepath.Join(t.TempDir(), "symlink-dst")
	require.NoError(t, os.Symlink(victimDir, symlinkDst))

	cases := []struct {
		name    string
		args    []string
		wantErr error
	}{
		{
			name: "upload-pack ref injection",
			args: []string{"pkg", "get", "https://github.com/GoogleCloudPlatform/marketplace-tools.git/examples/tf@--upload-pack=id", dst},
		},
		{
			name: "leading hyphen repo after double-dash",
			args: []string{"pkg", "get", "--", "-oProxyCommand=id:org/repo.git/pkg", dst},
		},
		{
			name: "file scheme rejected in CLI",
			args: []string{"pkg", "get", "file:///tmp/repo.git/pkg", dst},
		},
		{
			name: "local path rejected in CLI",
			args: []string{"pkg", "get", "/tmp/repo.git/pkg", dst},
		},
		{
			name: "ext remote helper rejected in CLI",
			args: []string{"pkg", "get", "ext::sh -c id.git/pkg", dst},
		},
		{
			name: "ssh ProxyCommand host rejected in CLI",
			args: []string{"pkg", "get", "ssh://-oProxyCommand=id/org/repo.git/pkg", dst},
		},
		{
			name:    "symlink destination rejected before preRunE stat",
			args:    []string{"pkg", "get", "https://github.com/GoogleCloudPlatform/marketplace-tools.git/examples/tf", symlinkDst},
			wantErr: fsutil.ErrSymlink,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := newSilentRootCmd()
			c.SetArgs(tc.args)
			err := c.Execute()
			require.Error(t, err)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			}
		})
	}
}

func TestCfgSetDoesNotExecuteStarlarkFunctions(t *testing.T) {
	pkgDir := t.TempDir()
	secretSentinel := "SUPER_SECRET_ENV_TOKEN_98765"
	t.Setenv("MPDEV_TEST_SECRET_TOKEN", secretSentinel)

	kptfileYAML := `apiVersion: kpt.dev/v1alpha1
kind: Kptfile
metadata:
  name: autogen-pkg
openAPI:
  definitions:
    io.k8s.cli.setters.image:
      x-k8s-cli:
        setter:
          name: image
          value: wordpress-old
    io.k8s.cli.setters.projectId:
      x-k8s-cli:
        setter:
          name: projectId
          value: old-project
functions:
  autoRunStarlark: true
  starlarkFunctions:
    - name: exfil
      path: exfil.star
`
	starlarkScript := `def run(r, token):
  for resource in r:
    resource["spec"]["exfiltrated"] = token
run(ctx.resource_list["items"], ctx.environment.get("MPDEV_TEST_SECRET_TOKEN", "ran"))
`
	configYAML := `apiVersion: dev.marketplace.cloud.google.com/v1alpha1
kind: DeploymentManagerTemplate
metadata:
  name: singlevm
spec:
  image: wordpress-old # {"$kpt-set":"image"}
  projectId: old-project # {"$ref":"#/definitions/io.k8s.cli.setters.projectId"}
`

	require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "Kptfile"), []byte(kptfileYAML), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "exfil.star"), []byte(starlarkScript), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "configurations.yaml"), []byte(configYAML), 0644))

	// Test both `cfg set` ($kpt-set shorthand) and `config set` alias ($ref).
	c1 := newSilentRootCmd()
	c1.SetArgs([]string{"cfg", "set", pkgDir, "image", "wordpress-new"})
	require.NoError(t, c1.Execute())

	c2 := newSilentRootCmd()
	c2.SetArgs([]string{"config", "set", pkgDir, "projectId", "clean-project"})
	require.NoError(t, c2.Execute())

	gotCfgBytes, err := os.ReadFile(filepath.Join(pkgDir, "configurations.yaml"))
	require.NoError(t, err)
	gotCfg := string(gotCfgBytes)
	assert.Contains(t, gotCfg, "image: wordpress-new")
	assert.Contains(t, gotCfg, "projectId: clean-project")
	assert.NotContains(t, gotCfg, secretSentinel)
	assert.NotContains(t, gotCfg, "exfiltrated")
}

func TestCfgSetDoesNotExecuteAutoRunStarlarkOrFetchURL(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ctx.resource_list['items'][0]['data']['pwned'] = 'yes'\n"))
	}))
	defer srv.Close()

	pkgDir := t.TempDir()
	kptfileYAML := `apiVersion: kpt.dev/v1alpha1
kind: Kptfile
metadata:
  name: remote-star-pkg
openAPI:
  definitions:
    io.k8s.cli.setters.replicas:
      type: integer
      x-k8s-cli:
        setter:
          name: replicas
          value: "1"
functions:
  autoRunStarlark: true
`
	fnConfigYAML := `apiVersion: v1
kind: ConfigMap
metadata:
  name: starlark-fn-cfg
  annotations:
    config.kubernetes.io/function: |
      starlark:
        name: remote-fn
        url: ` + srv.URL + `/payload.star
data:
  mode: test
`
	deployYAML := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
spec:
  replicas: 1 # {"$kpt-set":"replicas"}
`
	require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "Kptfile"), []byte(kptfileYAML), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "fn.yaml"), []byte(fnConfigYAML), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "deployment.yaml"), []byte(deployYAML), 0644))

	c := newSilentRootCmd()
	c.SetArgs([]string{"cfg", "set", pkgDir, "replicas", "3"})
	require.NoError(t, c.Execute())

	assert.Equal(t, int32(0), hits.Load(), "cfg set must not make outbound HTTP requests to starlark.url")
	gotDeploy, err := os.ReadFile(filepath.Join(pkgDir, "deployment.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(gotDeploy), "replicas: 3")
}

func TestCfgSetRejectsExplicitAutoRunTrueAndAllowsFalse(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	pkgDir := t.TempDir()
	kptfileYAML := `apiVersion: kpt.dev/v1alpha1
kind: Kptfile
metadata:
  name: flag-guard-pkg
openAPI:
  definitions:
    io.k8s.cli.setters.tag:
      x-k8s-cli:
        setter:
          name: tag
          value: v1
functions:
  autoRunStarlark: true
`
	fnYAML := `apiVersion: v1
kind: ConfigMap
metadata:
  name: fn
  annotations:
    config.kubernetes.io/function: |
      starlark:
        name: remote
        url: ` + srv.URL + `/evil.star
`
	resYAML := `apiVersion: v1
kind: ConfigMap
metadata:
  name: app
data:
  tag: v1 # {"$kpt-set":"tag"}
`
	require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "Kptfile"), []byte(kptfileYAML), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "fn.yaml"), []byte(fnYAML), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "app.yaml"), []byte(resYAML), 0644))

	for _, badArgs := range [][]string{
		{"cfg", "set", "--auto-run", pkgDir, "tag", "v2"},
		{"cfg", "set", "--auto-run=true", pkgDir, "tag", "v2"},
		{"config", "set", "--auto-run=true", pkgDir, "tag", "v2"},
	} {
		c := newSilentRootCmd()
		c.SetArgs(badArgs)
		err := c.Execute()
		require.Error(t, err, "expected error for %v", badArgs)
		assert.Contains(t, err.Error(), "--auto-run is disabled for security")
	}
	assert.Equal(t, int32(0), hits.Load())

	// Explicit --auto-run=false remains accepted for scripts/muscle memory.
	cOK := newSilentRootCmd()
	cOK.SetArgs([]string{"cfg", "set", "--auto-run=false", pkgDir, "tag", "v2"})
	require.NoError(t, cOK.Execute())
	assert.Equal(t, int32(0), hits.Load())

	gotApp, err := os.ReadFile(filepath.Join(pkgDir, "app.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(gotApp), "tag: v2")

	// Help output must hide --auto-run and never advertise default true.
	var helpBuf bytes.Buffer
	cHelp := newSilentRootCmd()
	cHelp.SetOut(&helpBuf)
	cHelp.SetErr(&helpBuf)
	cHelp.SetArgs([]string{"cfg", "set", "--help"})
	require.NoError(t, cHelp.Execute())
	assert.NotContains(t, helpBuf.String(), "--auto-run")
}

func TestGuardCfgSetFailsClosedWhenFlagMissing(t *testing.T) {
	called := false
	cfgCmd := &cobra.Command{Use: "cfg"}
	setCmd := &cobra.Command{
		Use: "set",
		RunE: func(c *cobra.Command, args []string) error {
			called = true
			return nil
		},
	}
	cfgCmd.AddCommand(setCmd)
	guardCfgSet(cfgCmd)

	cfgCmd.SilenceErrors = true
	cfgCmd.SilenceUsage = true
	cfgCmd.SetArgs([]string{"set", "some-dir", "k", "v"})
	err := cfgCmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing expected --auto-run flag")
	assert.False(t, called)
}

func newSilentRootCmd() *cobra.Command {
	c := GetMain()
	c.SilenceErrors = true
	c.SilenceUsage = true
	c.SetOut(&bytes.Buffer{})
	c.SetErr(&bytes.Buffer{})
	return c
}
