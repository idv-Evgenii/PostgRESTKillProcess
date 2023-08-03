# PostgREST
Killing the process for pods PostgREST
Since YA Cloud does not provide a superuser for postgres, accordingly, it is impossible to create a trigger to reset the schema cache so that PostgREST has up-to-date data, the second way to update the data from PostgREST is to delete the process
Above is written a dot that goes to a nearby raised sideCar container and deletes the running process