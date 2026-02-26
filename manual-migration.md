
# manual migration with grom
# install atlas cli
from the github release page
https://github.com/ariga/atlas/releases


install the grom atals package 
go install ariga.io/atlas-provider-gorm@latest

and create atlas.hcl

atlas migrate diff users --env dev --var url="mysql://root:root@localhost:3306/movies_test" 

to migrate 
atlas migrate apply --env prod --var url="mysql://root:root@localhost:3306/movies"
OR
atlas migrate apply --dir "file://migrations" --url="mysql://root:root@localhost:3306/movies"

# Instead of handling all models will create a list and bath integrate
# we will create loader to load the models
go get ariga.io/atlas-provider-gorm/gormschema