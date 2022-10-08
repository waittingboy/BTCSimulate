package utils

import "time"

func TimeFormat(timeStamp uint64) string {
	unixTime := time.Unix(int64(timeStamp), 0)

	return unixTime.Format("2006-01-02 15:04:05")
}
