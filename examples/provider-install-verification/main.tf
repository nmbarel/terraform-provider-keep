terraform {
  required_providers {
    keep = {
      source = "keep/test/keep"
    }
  }
}

provider "keep" {
  api_key = "d4d5ac53-c9fb-48f3-aa5c-e2e5836f4719"
  host_url = "http://localhost:8080"

}

resource "keep_workflows" "test" {
  workflows = [
  {yaml = "actions:\n- name: discord-action\n  provider:\n    config: '{{ providers.discorj }}'\n    type: discord\n    with:\n      content: noam\ndescription: n0am\nid: 581078b1-4f39-41f4-b58c-fdbb6f59e0b5\nowners: []\nservices: []\nsteps: []\ntriggers:\n- filters:\n  - key: source\n    value: r\".*\"\n  type: alert\n"}
  ]
}

output "test_workflows" {
  value = keep_workflows.test
}