data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./pkg/atlas/loader.go" ,
  ]
}

variable "url" {
  type = string
}

env "prod" {
  migration {
    dir = "file://migrations"
  }
  
  url = var.url
}

env "dev" {
  src = data.external_schema.gorm.url
  
  dev = var.url
  
  migration {
    dir = "file://migrations"
  }
  
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}