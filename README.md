## Installation

Clone this repository into the directory of your choice

```bash
git clone https://github.com/dallincoons/stool.git
```

Symlink one of the executables found in the `dist` folder to a directory that lies in 
your PATH.

For Mac users, run:
 ```bash
ln -s dist/stool_darwin_amd64 /usr/local/bin
```

For Linux users, run:
 ```bash
ln -s dist/stool_linux_amd64 /usr/local/bin
```

You will need to add a configuration file named `.surgio-tools.yml` to your home directory that contains
your local database credentials. For example:

```yaml
database:
  host: 127.0.0.1
  username: root
  password: secret
  database: dbname
  port: 3306
```

You can create this file by running
```bash
cp surgio-tools.yml.example ~/.surgio-tools.yml
```

But it's up to you to edit the configuration file with the appropriate details.

## Usage

### Clone a database

To clone a database, run
`stool clone --from db_name --to db_clone`

By default this will update any detected Laravel `.env` file that exists in the current working directory, 
to change the `DATABASE=` entry to point to the cloned database. If you wish to not update your `.env`
file, add the flag `--switch false`

 ## Drop a database
 
 To drop a database, run:
 `stool drop --name name_of_db` or `stool drop -n name_of_db`
 
 ## Switch your environment database 
 To change which database your environment is using, run:
 
 `stool switch --name new_db` or `stool switch -n new_db`
