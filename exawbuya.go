//echo "WISH YOU ALL THE LUCK of WORLD *PRANSHU*"; cat ~/Perndrive\ data/1m/buckets.txt | while IFS= read -r line; do aws s3 cp test_file.txt s3://"$line"; done


package main

import (
    "bufio"
    "fmt"
    "os"

    "github.com/robertkrimen/otto"
    "flag"
    "encoding/json"
)

func main() {
    // Define command-line flags
    outputJSON := flag.Bool("j", false, "Output data in JSON format")
    flag.Parse()

    vm := otto.New()

    _, err := vm.Run(javascriptCode)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error initializing JavaScript:", err)
        os.Exit(1)
    }

    scanner := bufio.NewScanner(os.Stdin)
    results := make([]S3Info, 0)

    for scanner.Scan() {
        url := scanner.Text()
        s3Info := getS3Info(vm, url)

        if s3Info != nil {
            results = append(results, *s3Info)
        }
    }

    if scanner.Err() != nil {
        fmt.Fprintln(os.Stderr, "Error reading input:", scanner.Err())
        os.Exit(1)
    }

    // Output the results
    if *outputJSON {
        jsonData, err := json.MarshalIndent(results, "", "    ")
        if err != nil {
            fmt.Fprintln(os.Stderr, "Error encoding JSON:", err)
            os.Exit(1)
        }
        fmt.Println(string(jsonData))
    } else {
        for _, s3Info := range results {
            fmt.Printf("URL: %s\n", s3Info.URL)
            fmt.Printf("Bucket: %s\n", s3Info.Bucket)
            fmt.Printf("Key: %s\n", s3Info.Key)
            fmt.Printf("Region: %s\n", s3Info.Region)
            fmt.Println()
        }
    }
}

type S3Info struct {
    URL     string
    Bucket  string
    Key     string
    Region  string
}

func getS3Info(vm *otto.Otto, url string) *S3Info {
    _, err := vm.Run(fmt.Sprintf("var url = '%s'; var result = s3ParseUrl(url);", url))
    if err != nil {
        return nil
    }

    result, _ := vm.Get("result")

    if !result.IsObject() {
        return nil
    }

    bucket, _ := result.Object().Get("bucket")
    key, _ := result.Object().Get("key")
    region, _ := result.Object().Get("region")

    if bucket.IsString() && key.IsString() && region.IsString() {
        return &S3Info{
            URL:     url,
            Bucket:  bucket.String(),
            Key:     key.String(),
            Region:  region.String(),
        }
    }

    return nil
}

const javascriptCode = `
function s3ParseUrl(url) {
  var decodedUrl = decodeURIComponent(url);

  var result = null;

  // http://s3.amazonaws.com/bucket/key1/key2
  var match = decodedUrl.match(/^https?:\/\/s3.amazonaws.com\/([^\/]+)\/?(.*?)$/);
  if (match) {
    result = {
      bucket: match[1],
      key: match[2],
      region: '',
    };
  }

  // http://s3-aws-region.amazonaws.com/bucket/key1/key2
  match = decodedUrl.match(/^https?:\/\/s3-([^.]+).amazonaws.com\/([^\/]+)\/?(.*?)$/);
  if (match) {
    result = {
      bucket: match[2],
      key: match[3],
      region: match[1],
    };
  }

  // http://bucket.s3.amazonaws.com/key1/key2
  match = decodedUrl.match(/^https?:\/\/([^.]+).s3.amazonaws.com\/?(.*?)$/);
  if (match) {
    result = {
      bucket: match[1],
      key: match[2],
      region: '',
    };
  }

  // http://bucket.s3-aws-region.amazonaws.com/key1/key2 or,
  // http://bucket.s3.aws-region.amazonaws.com/key1/key2
  match = decodedUrl.match(/^https?:\/\/([^.]+).(s3-|s3\.)([^.]+).amazonaws.com\/?(.*?)$/);
  if (match) {
    result = {
      bucket: match[1],
      key: match[4],
      region: match[3],
    };
  }

  // https://s3.us-west-1.amazonaws.com/bucket/
  match = decodedUrl.match(/^https:\/\/s3\.([^/]+)\.amazonaws.com\/([^/]+)\/?$/);
  if (match) {
    result = {
      bucket: match[2],
      key: '',
      region: match[1],
    };
  }

  return result;
}
`

