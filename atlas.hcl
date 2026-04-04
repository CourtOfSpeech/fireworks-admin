# atlas.hcl
env "local" {
  src = "ent://internal/ent/schema"
  url = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
  dev = "docker://postgres/16/dev?search_path=public"

  migration {
    dir = "file://migrations"
    format = atlas
  }
}
