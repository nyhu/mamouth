package c14

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dutchcoders/goftp"
	"github.com/nyhu/mamouth/entity"
)

type BucketStatus struct {
	status string
	bucket entity.Bucket
}

func Melt(safeID string, start time.Time, end time.Time, batchChannel chan entity.KafkaBatch) error {

	relevantArchives, err := getRelevantArchives(safeID, start, end)
	if err != nil {
		return err
	}
	if err := downloadArchive(relevantArchives, safeID); err != nil {
		return err
	}
	for _, relevantArchive := range relevantArchives {
		var relevantBatchs = getRelevantBatchs(relevantArchive, safeID, start, end)
		for _, relevantBatch := range relevantBatchs {
			batch, err := getBatch(relevantBatch)
			if err != nil {
				return err
			}
			batchChannel <- batch
		}
	}
	return nil
}

/*
** Fetch All Archives and Return Only Those Within Date Range.
 */
func getRelevantArchives(safeID string, start time.Time, end time.Time) ([]entity.Archive, error) {
	allArchives, err := GetAllArchives(safeID)
	if err != nil {
		return nil, err
	}
	return filterRelevantArchives(allArchives, start, end)
}

/*
** List All Batchs in Archive File and Return Only Those Within Date Range.
 */
func getRelevantBatchs(relevantArchive entity.Archive, safeID string, start time.Time, end time.Time) []string {
	var relevantBatchs []string
	archivePath := "/tmp/received/" + safeID + "/" + relevantArchive.Name
	files, err := ioutil.ReadDir(archivePath)
	if err != nil {
		return relevantBatchs
	}
	for _, file := range files {
		fileName := file.Name()
		fileTime, err := decodeFileName(fileName)
		if err != nil {
			return relevantBatchs
		}
		if timeWithinRange(fileTime, start, end) {
			relevantBatchs = append(relevantBatchs, fileName)
		}
	}
	return relevantBatchs

}

/*
** Return Archives List Within Date Range.
 */
func filterRelevantArchives(allArchives []entity.Archive, start time.Time, end time.Time) ([]entity.Archive, error) {
	var relevantArchives []entity.Archive

	for _, archive := range allArchives {
		archiveStart, err := decodeFileName(archive.Name)
		if err != nil {
			return nil, err
		}
		archiveEnd := archiveStart.Add(time.Hour)
		if timeWithinRange(archiveStart, start, end) || timeWithinRange(archiveEnd, start, end) {
			relevantArchives = append(relevantArchives, archive)
		}
	}
	return relevantArchives, nil
}

/*
** Start Unarchiving Remote Archive, Check Job, FTP Transfer to Local Storage
 */
func downloadArchive(relevantArchives []entity.Archive, safeID string) error {
	var attempt = 0
	var maxAttempt = 120
	var isAllDownloaded = false

	var statusArchives []BucketStatus

	for index, archive := range relevantArchives {
		archiveID := archive.Uuid_ref
		_, err := Unarchive(archiveID, safeID)
		if err != nil {
			return err
		}
		statusArchives = append(statusArchives, checkStatus(relevantArchives[index]))
	}

	for !isAllDownloaded || attempt < maxAttempt {
		isAllDownloaded = true
		for index, _ := range relevantArchives {
			if statusArchives[index].status == "PENDING" {
				statusArchives[index] = checkStatus(relevantArchives[index])
				isAllDownloaded = false
			}
			if statusArchives[index].status == "COMPLETED" {
				ftpArchive(statusArchives[index].bucket, relevantArchives[index], safeID, &statusArchives[index])
				statusArchives[index].status = "DOWNLOADING"
				isAllDownloaded = false
			}
		}
		attempt = attempt + 1
		time.Sleep(60 * time.Second)
	}
	if !isAllDownloaded && attempt >= maxAttempt {
		return errors.New("impossible to recover all archives")
	}
	return nil
}

/*
** Start FTP Download, Then Re-Archive Remote Bucket. Update Archive Status Array.
 */
func ftpArchive(bucket entity.Bucket, archive entity.Archive, safeID string, archiveStatus *BucketStatus) {

	credentials := entity.NewCredentials(bucket.Credentials[0].(map[string]interface{}))
	localPath := "/tmp/received/" + safeID
	ensureTopicDirectory(localPath)
	localPath += "/" + archive.Name
	ftp, err := ConnectToBucket(credentials)
	if err != nil {
		return
	}
	GetFromBucket(ftp, localPath)
	Freeze(archive.Uuid_ref)
	archiveStatus.status = "DOWNLOADED"
}

/*
** Check If Archive Is Available For FTP Download.
 */
func checkStatus(archive entity.Archive) BucketStatus {
	archiveID := archive.Uuid_ref
	bucket, err := GetBucket(archiveID)
	if err != nil {
		return BucketStatus{status: "PENDING"}
	}
	return BucketStatus{status: "COMPLETED", bucket: bucket}
}

/*
** Read File and Return Batch Structure
 */
func getBatch(fileName string) (entity.KafkaBatch, error) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		return entity.KafkaBatch{}, err
	}
	var batch entity.KafkaBatch
	err = json.Unmarshal(raw, &batch)
	if err != nil {
		return entity.KafkaBatch{}, err
	}
	return batch, nil
}

/*
func sortBatch(batch []entity.KafkaMessage) []entity.KafkaMessage {
	sort.Slice(batch, func(i, j int) bool {
		return batch[i].Offset < batch[j].Offset
	})
	return batch
}
*/

/*
** Transform Time Timestamp into FileName
 */
func encodeArchiveName(start time.Time) string {
	return fmt.Sprintf("%v", start.Unix())
}

func timeWithinRange(date time.Time, start time.Time, end time.Time) bool {
	return date.After(start) && date.Before(end)
}

/*
** Transform FileName into Timestamp
 */
func decodeFileName(fileName string) (time.Time, error) {
	timestamp, err := strconv.Atoi(fileName)
	if err != nil {
		return time.Now(), err
	}
	now := time.Unix(int64(timestamp), 0)
	return now, nil
}

/*
** FTP Transfer
 */
func GetFromBucket(ftp *goftp.FTP, localPath string) error {

	err := ftp.Walk("/", func(path string, info os.FileMode, err error) error {
		if err != nil {
			return err
		}
		_, err = ftp.Retr(path, func(r io.Reader) error {
			fullPath := localPath + path
			file, err := os.Create(fullPath)
			if err != nil {
				return err
			}
			_, err = io.Copy(file, r)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

/*
** Ensure we can write in Topic Directory
 */
func ensureTopicDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0644)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

/*
** API Call for Unfreezing Archive
 */
func Unarchive(archiveID string, safeID string) (int, error) {

	client := &http.Client{}
	data := entity.UnArchive{Platform: "2", Protocols: []string{"ftp"}}

	encoded, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	url := BASE_URL + "/storage/c14/safe/" + safeID + "/archive/" + archiveID + "/unarchive"
	req, err := http.NewRequest("POST", url, strings.NewReader(string(encoded)))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	return res.StatusCode, nil
}
