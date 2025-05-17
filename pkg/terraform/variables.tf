variable "project_id" {
  description = "The GCP project where secrets will be created"
  type        = string
}

variable "secrets" {
  description = "Map of secret names to their values"
  type        = map(string)
}
