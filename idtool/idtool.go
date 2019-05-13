package main

import id "github.com/koofr/go-idutils"
import "flag"
import "time"
import "fmt"

func main() {
	srcId := flag.Int64("id", 0, "Convert ID to timestamp")
	srcEpoch := flag.Int64("unix", 0, "Convert a Unix timestamp to ID")
	srcTime := flag.String("time", "", "Convert time to ID (allows for various common formats)")
	flag.Parse()

	if srcId != nil && *srcId > 0 {
		tiem := id.IdToTime(*srcId)
		fmt.Printf("%s\n", tiem.Format(time.RFC822))
	}

	if srcEpoch != nil && *srcEpoch > 0 {
		fmt.Printf("%d\n", id.IdEndOfTimestamp(*srcEpoch*1000))
	}

	if srcTime != nil && *srcTime != "" {
		formats := []string{
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			time.RFC822,
			time.RFC822Z,
			time.RFC850,
			time.RFC1123,
			time.RFC1123Z,
			time.RFC3339,
			time.RFC3339Nano,
		}
		for _, f := range formats {
			t, err := time.Parse(f, *srcTime)
			if err == nil {
				fmt.Printf("%d\n", id.IdEndOfTime(t))
				break
			}
		}
	}

}
