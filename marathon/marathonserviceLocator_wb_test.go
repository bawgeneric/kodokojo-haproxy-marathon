package marathon_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/kodokojo/kodokojo-haproxy-marathon/marathon"
	"testing"
)

func Test_valid_extract_from_json(t *testing.T) {

	jsonInput := GetValidJson()
	marathonServiceLocator := marathon.NewMarathonServiceLocator("http://localhost:8080")
	services := marathonServiceLocator.ExtractServiceFromJson([]byte(jsonInput))

	assert.NotEmpty(t, services)
	assert.Equal(t, 2, len(services))

	for _, service := range services {
		assert.Equal(t, "acme", service.ProjectName)
		assert.Equal(t, 1, len(service.HaProxyHTTPEntries))

		httpEntry := service.HaProxyHTTPEntries[0]
		assert.Equal(t, 1, len(httpEntry.Backends))
		if httpEntry.EntryName == "scm" {
			assert.Equal(t, 1, len(service.HaProxySSHEntries))
		} else if httpEntry.EntryName == "ci" {
			assert.Equal(t, 0, len(service.HaProxySSHEntries))
		}
	}

}

func GetValidJson() string {
	return `{
  "apps": [
    {
      "id": "/acme/scm",
      "cmd": null,
      "args": null,
      "user": null,
      "env": {},
      "instances": 1,
      "cpus": 2,
      "mem": 2048,
      "disk": 0,
      "executor": "",
      "constraints": [
        [
          "type",
          "LIKE",
          "operator"
        ]
      ],
      "uris": [],
      "fetch": [],
      "storeUrls": [],
      "ports": [
        10002,
        10003
      ],
      "requirePorts": false,
      "backoffSeconds": 1,
      "backoffFactor": 1.15,
      "maxLaunchDelaySeconds": 3600,
      "container": {
        "type": "DOCKER",
        "volumes": [],
        "docker": {
          "image": "gitlab/gitlab-ce",
          "network": "BRIDGE",
          "portMappings": [
            {
              "containerPort": 80,
              "hostPort": 0,
              "servicePort": 10002,
              "protocol": "tcp"
            },
            {
              "containerPort": 22,
              "hostPort": 0,
              "servicePort": 10003,
              "protocol": "tcp"
            }
          ],
          "privileged": false,
          "parameters": [
            {
              "key": "env",
              "value": "GITLAB_OMNIBUS_CONFIG=gitlab_rails['gitlab_shell_ssh_port'] = 40022"
            },
            {
              "key": "label",
              "value": "project=acme"
            },
            {
              "key": "label",
              "value": "componentType=scm"
            },
            {
              "key": "label",
              "value": "component=gitlab"
            }
          ],
          "forcePullImage": false
        }
      },
      "healthChecks": [
        {
          "path": "/",
          "protocol": "HTTP",
          "portIndex": 0,
          "gracePeriodSeconds": 180,
          "intervalSeconds": 20,
          "timeoutSeconds": 20,
          "maxConsecutiveFailures": 10,
          "ignoreHttp1xx": false
        }
      ],
      "dependencies": [],
      "upgradeStrategy": {
        "minimumHealthCapacity": 1,
        "maximumOverCapacity": 1
      },
      "labels": {
        "project": "acme",
        "componentType": "scm",
        "component": "gitlab"
      },
      "acceptedResourceRoles": null,
      "ipAddress": null,
      "version": "2016-03-14T14:22:52.747Z",
      "versionInfo": {
        "lastScalingAt": "2016-03-14T14:22:52.747Z",
        "lastConfigChangeAt": "2016-03-14T14:22:52.747Z"
      },
      "tasksStaged": 0,
      "tasksRunning": 1,
      "tasksHealthy": 1,
      "tasksUnhealthy": 0,
      "deployments": [],
      "tasks": [
        {
          "id": "acme_scm.3ddaa47a-e9f0-11e5-b49b-022b20c33e4f",
          "host": "52.50.53.24",
          "ipAddresses": [],
          "ports": [
            43052,
            43053
          ],
          "startedAt": "2016-03-14T14:24:22.396Z",
          "stagedAt": "2016-03-14T14:22:54.384Z",
          "version": "2016-03-14T14:22:52.747Z",
          "slaveId": "f333bb28-7ade-422d-ad16-4a5c6eab957a-S1",
          "appId": "/acme/scm",
          "healthCheckResults": [
            {
              "alive": true,
              "consecutiveFailures": 0,
              "firstSuccess": "2016-03-14T14:25:33.202Z",
              "lastFailure": null,
              "lastSuccess": "2016-03-14T16:34:20.671Z",
              "taskId": "acme_scm.3ddaa47a-e9f0-11e5-b49b-022b20c33e4f"
            }
          ]
        }
      ]
    },
    {
      "id": "/acme/ci",
      "cmd": null,
      "args": null,
      "user": null,
      "env": {},
      "instances": 1,
      "cpus": 1,
      "mem": 2048,
      "disk": 0,
      "executor": "",
      "constraints": [
        [
          "type",
          "LIKE",
          "operator"
        ]
      ],
      "uris": [],
      "fetch": [],
      "storeUrls": [],
      "ports": [
        10004
      ],
      "requirePorts": false,
      "backoffSeconds": 1,
      "backoffFactor": 1.15,
      "maxLaunchDelaySeconds": 3600,
      "container": {
        "type": "DOCKER",
        "volumes": [],
        "docker": {
          "image": "jenkins",
          "network": "BRIDGE",
          "portMappings": [
            {
              "containerPort": 8080,
              "hostPort": 0,
              "servicePort": 10004,
              "protocol": "tcp"
            }
          ],
          "privileged": false,
          "parameters": [
            {
              "key": "env",
              "value": "JENKINS_SLAVE_AGENT_PORT=50000"
            },
            {
              "key": "env",
              "value": "JAVA_OPTS=-Duser.timezone=Europe/Paris"
            },
            {
              "key": "label",
              "value": "project=acme"
            },
            {
              "key": "label",
              "value": "componentType=ci"
            },
            {
              "key": "label",
              "value": "component=jenkins"
            }
          ],
          "forcePullImage": false
        }
      },
      "healthChecks": [
        {
          "path": "/",
          "protocol": "HTTP",
          "portIndex": 0,
          "gracePeriodSeconds": 180,
          "intervalSeconds": 20,
          "timeoutSeconds": 20,
          "maxConsecutiveFailures": 10,
          "ignoreHttp1xx": false
        }
      ],
      "dependencies": [],
      "upgradeStrategy": {
        "minimumHealthCapacity": 1,
        "maximumOverCapacity": 1
      },
      "labels": {
        "project": "acme",
        "componentType": "ci",
        "component": "jenkins"
      },
      "acceptedResourceRoles": null,
      "ipAddress": null,
      "version": "2016-03-14T16:19:47.156Z",
      "versionInfo": {
        "lastScalingAt": "2016-03-14T16:19:47.156Z",
        "lastConfigChangeAt": "2016-03-14T16:19:47.156Z"
      },
      "tasksStaged": 0,
      "tasksRunning": 1,
      "tasksHealthy": 1,
      "tasksUnhealthy": 0,
      "deployments": [],
      "tasks": [
        {
          "id": "acme_ci.9269b9ea-ea00-11e5-b49b-022b20c33e4f",
          "host": "52.50.28.42",
          "ipAddresses": [],
          "ports": [
            46832
          ],
          "startedAt": "2016-03-14T16:19:49.045Z",
          "stagedAt": "2016-03-14T16:19:48.197Z",
          "version": "2016-03-14T16:19:47.156Z",
          "slaveId": "f333bb28-7ade-422d-ad16-4a5c6eab957a-S0",
          "appId": "/acme/ci",
          "healthCheckResults": [
            {
              "alive": true,
              "consecutiveFailures": 0,
              "firstSuccess": "2016-03-14T16:20:27.841Z",
              "lastFailure": null,
              "lastSuccess": "2016-03-14T16:34:28.098Z",
              "taskId": "acme_ci.9269b9ea-ea00-11e5-b49b-022b20c33e4f"
            }
          ]
        }
      ]
    }
  ]
}`
}
