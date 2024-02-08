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

data "keep_workflows" "test" {}

output "test_workflows" {
  value = data.keep_workflows.test
}

