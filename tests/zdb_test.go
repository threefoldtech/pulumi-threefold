package test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/go-redis/redis"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestZDB(t *testing.T) {
	mnemonic := os.Getenv("MNEMONIC")
	assert.NotEmpty(t, mnemonic)

	network := os.Getenv("NETWORK")
	if network == "" {
		network = "dev"
	}

	cwd, _ := os.Getwd()

	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Quick:            true,
		SkipRefresh:      true,
		DestroyOnCleanup: true,
		Dir:              path.Join(cwd, "..", "examples/zdb"),
		Config: map[string]string{
			"MNEMONIC": mnemonic,
			"NETWORK":  network,
		},
		ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
			for _, res := range stack.Deployment.Resources {
				if res.Type == "grid:internal:Deployment" {
					assert.NotEmpty(t, res.Outputs["node_deployment_id"])

					zdb := res.Outputs["zdbs"].([]interface{})[0].(map[string]interface{})
					computed := res.Outputs["zdbs_computed"].([]interface{})[0].(map[string]interface{})

					zdbEndpoint := fmt.Sprintf("[%s]:%v", computed["ips"].([]interface{})[1].(string), computed["port"].(float64))
					zdbNamespace := computed["namespace"].(string)
					password := zdb["password"].(string)

					rdb := redis.NewClient(&redis.Options{
						Addr: zdbEndpoint,
					})
					_, err := rdb.Do("SELECT", zdbNamespace, password).Result()
					assert.NoError(t, err)

					_, err = rdb.Set("key1", "val1", 0).Result()
					assert.NoError(t, err)

					res, err := rdb.Get("key1").Result()
					assert.NoError(t, err)
					assert.Equal(t, res, "val1")
				}
			}
		},
	})
}
