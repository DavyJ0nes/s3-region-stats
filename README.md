# S3 Bucket Location Stats

## Description

Simple script to get the number of buckets that are in each AWS Region.

## Usage

Usage instructions:

```shell
$ AWS_PROFILE=profile s3-region-stats
           Region|   Count|
        eu-west-1|      89|
        us-west-2|      25|
        No Region|      19|
        us-west-1|       8|
        eu-west-2|       6|
     ca-central-1|       5|
     eu-central-1|       4|
   ap-southeast-2|       4|
   ap-southeast-1|       2|
        us-east-2|       1|
        eu-west-3|       1|
   ap-northeast-2|       1|
   ap-northeast-1|       1|
        sa-east-1|       1|
               EU|       1|

3.26s elapsed

# Get trace output
$ AWS_PROFILE=profile s3-region-stats -trace 1> /dev/null 2> output.trace
$ go tool trace output.trace
2018/06/03 17:17:00 Parsing trace...
2018/06/03 17:17:00 Serializing trace...
2018/06/03 17:17:00 Splitting trace...
2018/06/03 17:17:00 Opening browser
```

## TODO

- [ ] Expand to count EC2 instances
- [ ] Improve error handling

## License

[MIT](./LICENSE)
