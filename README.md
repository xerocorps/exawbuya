Tool that extracts AWS S3 Buckets NAME from list of URL's

1. Add verbosity to print more information about the processing steps.
2. Add a flag `-o` to specify the output file.
3. Add an option to save only "Buckets" based on a flag `-b`.
4. Refactor the code to handle saving to file directly.
5. Implement the logic to save only buckets if the `-b` flag is set.


This updated version includes verbose output for each URL processed, the option to save only bucket information, and the ability to save the output directly to a file specified by the `-o` flag.
