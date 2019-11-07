# CHANGELOG

## v0.24.1 (Thu, 7 Nov 2019)
+ Inventory managment including create, get, list, update, batch update and delete.
+ Example YAML file for inventory batch update.

## v0.23.0 (Tue, 5 Nov 2019)
+ Coupon management command.

## v0.22.0 (Tue, 5 Nov 2019)
+ Address sub command.
+ Reintroduce `pkg/errors` for stacktracing.

## v0.21.1 (Sun, 3 Nov 2019)
+ Fix sysinfo to include stripe and additional app settings.
+ Update Go dependencies using `go get -v all`.
+ Requires API at or above `v0.61.2`.

## v0.21.0 (Fri, 1 Nov 2019)
+ Create, get, list and delete promo rules. Uses interactive prompt.
+ Create a shipping tariff with interactive prompt.

## v0.20.0 (Thu, 31 Oct 2019)
+ Cart management sub command for cart creation and product add/update/list/remove and empty.

## v0.19.1 (Mon 22 Oct 2019)
+ Fix products get, list and delete as API now uses product ids, not skus.

## v0.19.0 (Mon 22 Oct 2019)
+ Remove admins command as it now exists under users.
+ Product to product associations groups create, get, list and delete.
+ Product to product associations list and delete.

## v0.18.0 (Mon 21 Oct 2019)
+ Create, get, list, update and delete webhooks.

## v0.17.0 (Mon 21 Oct 2019)
+ Create and list developer keys.
+ Hide developer key for profiles create sub command.
+ Create and list users with roles.

## v0.16.0 (Mon 21 Oct 2019)
+ Promo rule and shipping tariff management.
+ Go 1.13 error handling with %w.
+ Silence Makefile.

## v0.15.1 (Fri 27 Sep 2019)
+ Prices and Images HTTP responses changed.

## v0.15.0 (Tue 3 Sep 2019)
+ Price lists.
+ Product category relations.

## v0.14.0 (Fri 30 Aug 2019)
+ Brought up to date to match API restructuring.

## v0.13.2 (Mon 1 Jul 2019)
+ API uses integers for pricing not floats.
+ Use `pkg/errors` for better error handling.

## v0.13.1 (Mon 1 Jul 2019)
+ `/assocs` changed to `/associations` in API `v0.50.0`.
+ Fix bug with `GetCatalogAssocs` always returning nil map.
+ `assocs list` displays associations in sorted order. Maps have randomized keys.

## v0.13.0 (Tue 25 Jun 2019)
+ Requires API `v0.49.0` onwards.
+ API `v0.49.0` returns Google GAE Project ID and Firebase separately so `ecom sysinfo` has new output.
+ Makefile for building removing hardcoded `version` string value from `main.go`.

## v0.12.1 (Tue 18 Jun 2019)
+ Fixed associations list.
+ Requires API `v0.46.1`.

## v0.11.0 (Fri 6 Jun 2019)
+ Alter `build.sh` script to do binary builds only.
+ Reverting to `yaml.v2` as `yaml.v3` has a bug causing crashes whilst applying associations.
+ `products apply` accepts single product YAML files or directories of files.
+ `assocs apply` handles 409 HTTP status error messages.

## v0.10.0 (Mon 23 May 2019)
+ Removed injection of version string to use hard coded string instead.
+ Add missing completion sub command.

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
