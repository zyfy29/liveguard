from dynaconf import Dynaconf

settings = Dynaconf(
    envvar_prefix="DYNACONF", # env var only
)

# `envvar_prefix` = export envvars with `export DYNACONF_FOO=bar`.
# `settings_files` = Load these files in the order.
