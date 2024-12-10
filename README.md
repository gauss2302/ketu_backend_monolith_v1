# ketu_backend_monolith_v1



### Tips for Db Connection Pooling:

* maxOpenConns:

    Formula: (number of CPU cores * 2) + effective_spindle_count
    For most applications, 25-50 is a good start
    Monitor for connection wait times to adjust


* connMaxLifetime:

    Keep it lower than your database server's timeout
    Typically 15 minutes is good
    If using a connection proxy (like PgBouncer), make it shorter


* maxIdleConns:

    Should be equal to or less than maxOpenConns
    For busy applications, keep it close to maxOpenConns
    For less busy apps, can be lower to save resources


* connMaxIdleTime:

    10-15 minutes is typical
    Lower if memory is a concern
    Higher if you have consistent traffic