# atlas.hcl
env "local" {
  src = "ent://internal/infrastructure/persistence/ent/schema"
  url = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
  # 还没有安装docker，先使用本地的
  dev = "postgres://postgres:postgres@localhost:5432/postgres_dev?sslmode=disable"
  
  migration {
    dir = "file://migrations"
    format = atlas
  }
}