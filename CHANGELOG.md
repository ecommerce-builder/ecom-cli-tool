## CHANGELOG
## v0.9.0 (Thu 23 May 2019)
+ Complete restructure and organized of code into sub commands.
+ Remove bad global deps.
+ Code tidy.

## v0.8.1 (Wed 22 May 2019)
+ Fix for `ecom products get <sku>` govalidate experimental code removed.

## v0.8.0 (Wed 22 May 2019)
+ Severed cross depedency on the API Service codebase.
+ Moved to github.com/ecommerce-builder/ecom-cli-tool public repo.

## v0.7.2 (Wed 22 May 2019)
+ catalog purge shows warning and aborts if associations still exist (409 Conflict returned from API).
+ Modified build.sh to auto copy binaries to the ecom-docs repo public downloads dir.
+ Fixes ecom catalog get output to include the segment and name combined.

## v0.7.1 (Wed 22 May 2019)
+ Fixes broken `ecom catalog apply ...` to update the catalog.
+ Adds endpoints guard to the catalog.yaml file.
+ Requires API Version v0.34.1 or above.

## v0.7.0 (Mon 20 May 2019)
+ `ecom admins remove <uuid>` calls API service OpDeleteAdmin

## v0.6.1 (Mon 20 May 2019)
+ Remove spurious fmt.Printf from admins create.

## v0.6.0 (Mon 20 May 2019)
+ OpCreateAdmin feature added `ecom admins create`.

## v0.5.1 (Mon 20 May 2019)
+ Broken tag fix.

## v0.5.0 (Mon 20 May 2019)
+ admins list operation (requires API v0.33.0 and above).

## v0.4.0 (Fri 17 May 2019)
+ Add Web API Key to sysinfo output.
+ configs becomes profiles.
+ Add new profiles using endpoint and devkey only.
+ No longer asks for Firebase (Google) Web API Keys but instead retrieves from from the API endpoint.

## v0.3.0 (Tue 14 May 2019)
+ `ecom catalog apply <catalog.yaml>` to replace the entire catalog.
+ `ecom catalog purge` to purge the entire catalog.

## v0.2.0 (Mon 13 May 2019)
+ Compatible with API v0.28.4
+ `ecom assocs <assocs.yaml>` apply a new catalog associations YAML file.
+ `ecom assocs list` list current associations.
+ `ecom assocs purge` purge the existing catalog associations.
 
## v0.1.0 (Thu 18 Apr 2019)
+ Initial command line tool providing `apply <product.yaml` feature.
+ Requires API v0.22.0 or above.
