## CHANGELOG
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
