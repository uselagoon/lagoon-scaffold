# Lagoon Scaffold

*THIS IS NOT A PRODUCTION READY TOOL*

This is a POC tool for Lagoonizing and scaffolding out Lagoon projects.

It is a simple script that points to a number of git repos that contain scaffolding information.
A scaffold contains, minimally, a `lagoon-docker-compose.yml.tmpl`, a `.lagoon.yml.tmpl` file and a `values.yml`.

Given a selected scaffold, we check out the latest commit of the scaffold's branch and clone it into the target directory into a temporary directory (which is removed post lagoonization).
If in interactive mode, we open up the user's editor to show them the values file that is going to be applied to their lagoonization.
This will give them a chance to change values if they need to. Eventually, we hope to make this interactive.

Once they've saved their values file, we recursively search through the cloned scaffolding and apply the values to any `.tmpl` files we find.
We then strip the `.tmpl` from the file name and copy the concretized data to disk.

Once this is done, we copy all the files from the scaffold directory into the target directory.

If there is a `.lagoon/post-message.txt` file, this is shown to the user.

Finally, the temporary directory with the scaffolding is removed.


## Usage example

### Lagoonizing a new Laravel 10 project

Install Laravel (see [installation docs](https://laravel.com/docs/10.x/installation#getting-started-on-linux)).
If your installation is installed at, say, `/home/myaccount/projects/example-app` you can run the following

```
lagoon-init-prot init scaffold --scaffold=laravel --targetdir=/home/myaccount/projects/example-app
```

The command above will ask you to check the values that are going to be used by the installation (you can skip this by
 passing `--no-interaction=true` as a flag and use the defaults).

Your Laravel project should now be ready to be pushed up to Lagoon.
